package impl

import (
	"strings"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"

	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

var detectables = []Detectable{
	{
		String:  "[lightningbot]",
		Selfbot: "LightningBot",
		Score:   1,
	},
	{
		String:  "[lightningbot free]",
		Selfbot: "LightningBot",
		Score:   1,
	},
	{
		String:  "www.team-luna.org",
		Selfbot: "Luna",
		Score:   2,
	},
	{
		String:  "enabled protections » ",
		Selfbot: "Luna",
		Score:   2,
	},
	{
		String:  "nighty.one",
		Selfbot: "Nighty",
		Score:   1,
	},
	{
		String:  "<> is required",
		Selfbot: "",
		Score:   2,
	},
	{
		String:  "<> = required",
		Selfbot: "",
		Score:   2,
	},
	{
		String:  "> **",
		Selfbot: "",
		Score:   0.5,
	},
}

type Detectable struct {
	String  string
	Selfbot string
	Score   float64
}

type DetectionEvent struct {
	Event   string
	Message *discordgo.Message
}

type DetectionScore struct {
	UserID  string
	Score   float64
	Logged  bool
	Selfbot string
	Events  []DetectionEvent
}

type SelfbotDetection struct {
	uiSvc             interfaces.UIService
	eventReactionsSvc interfaces.EventReactionsService

	DetectionScores sync.Map
}

func NewSelfbotDetection(uiSvc interfaces.UIService, eventReactionsSvc interfaces.EventReactionsService) *SelfbotDetection {
	return &SelfbotDetection{
		uiSvc:             uiSvc,
		eventReactionsSvc: eventReactionsSvc,
		DetectionScores:   sync.Map{},
	}
}

func (s *SelfbotDetection) Type() interfaces.EventReactionType {
	return interfaces.EventReactionSelfbotDetection
}

func (s *SelfbotDetection) Description() string {
	return "Selfbot Detection (default)"
}

func (s *SelfbotDetection) Equals(other interfaces.EventReaction) bool {
	_, ok := other.(*SelfbotDetection)

	return ok
}

func (s *SelfbotDetection) Data() map[string]any {
	return map[string]any{}
}

func (s *SelfbotDetection) ensureUserEntry(userID string) *DetectionScore {
	rawScore, _ := s.DetectionScores.LoadOrStore(userID, &DetectionScore{UserID: userID, Score: 0})
	return rawScore.(*DetectionScore)
}

func (s *SelfbotDetection) logDetection(score *DetectionScore, session *discordgo.Session, message *discordgo.Message) {
	if score.Score < 15 || score.Logged {
		return
	}

	selfbot := utils.Misc.Or(score.Selfbot, "Unknown")
	s.uiSvc.LogSelfbotDetection(message, selfbot)
	score.Logged = true
}

func (s *SelfbotDetection) Handler(session *discordgo.Session, event any) {
	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Debug("SelfbotDetection.Handler: Error getting config: " + err.Error())
		return
	}

	if message, ok := event.(*discordgo.MessageCreate); ok {
		if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == session.State.User.ID {
			return
		}

		s.CreateHandler(session, message, &cfg)
	} else if message, ok := event.(*discordgo.MessageUpdate); ok {
		if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == session.State.User.ID {
			return
		}

		s.UpdateHandler(session, message, &cfg)
	} else if message, ok := event.(*discordgo.MessageDelete); ok {
		if !cfg.Events.SelfbotDetection || message.Author.Bot || message.Author.ID == session.State.User.ID {
			return
		}

		s.DeleteHandler(session, message, &cfg)
	} else {
		utils.Logging.Debug("SelfbotDetection.Handler: Event is not a message create or delete event")
	}
}

func (s *SelfbotDetection) CreateHandler(session *discordgo.Session, message *discordgo.MessageCreate, cfg *config.Config) {
	score := s.ensureUserEntry(message.Author.ID)
	for _, detectable := range detectables {
		if !strings.Contains(strings.ToLower(message.Content), detectable.String) {
			continue
		}

		score.Score += detectable.Score
		score.Events = append(score.Events, DetectionEvent{
			Event:   "detectable",
			Message: message.Message,
		})
		score.Selfbot = utils.Misc.Or(detectable.Selfbot, score.Selfbot)

		s.logDetection(score, session, message.Message)
	}

	s.DetectionScores.Store(message.Author.ID, score)
}

func (s *SelfbotDetection) UpdateHandler(session *discordgo.Session, message *discordgo.MessageUpdate, cfg *config.Config) {
	s.ensureUserEntry(message.Author.ID)
}

func (s *SelfbotDetection) DeleteHandler(session *discordgo.Session, message *discordgo.MessageDelete, cfg *config.Config) {
	score := s.ensureUserEntry(message.Author.ID)
	timeDiff := time.Since(utils.Discord.SnowflakeToTime(message.ID)) - s.eventReactionsSvc.GetTimeDifference()
	if timeDiff < time.Second*1 {
		score.Events = append(score.Events, DetectionEvent{
			Event:   "fast_delete",
			Message: message.Message,
		})
		score.Score += 0.5

		s.logDetection(score, session, message.Message)
	}

	s.DetectionScores.Store(message.Author.ID, score)
}
