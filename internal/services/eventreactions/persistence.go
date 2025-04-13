package eventreactions

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/impl"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

var ErrInvalidReactionType = errors.New("invalid reaction type")

type eventReactionConfig struct {
	Type types.EventReactionType `json:"Type"`
	Data map[string]any          `json:"Data"`
}

func (e *eventReactionsService) startLoading(session *discordgo.Session) {
	timer := time.NewTimer(5 * time.Second)
	isFirstIteration := true

	for range timer.C {
		if err := e.loadEventReactions(session, isFirstIteration); err != nil {
			e.eventBroker.Publish(eventbroker.EventLogError, fmt.Sprintf("Failed to load event reactions: %v", err))
		}

		if err := e.LoadCmdAliases(); err != nil {
			e.eventBroker.Publish(eventbroker.EventLogError, fmt.Sprintf("Failed to load command aliases: %v", err))
		}

		if isFirstIteration {
			isFirstIteration = false
		}
	}
}

func (e *eventReactionsService) loadEventReactions(session *discordgo.Session, isFirstLoad bool) error {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	cfg := e.configSvc.GetConfig()

	reset := func() {
		e.eventBroker.Publish(eventbroker.EventLogError, "Failed to load event reactions. Event reactions have been reset.")

		cfg := e.configSvc.GetConfig()
		cfg.EventReactions = []byte("[]")
		_ = e.configSvc.SaveConfig(cfg)
	}

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		logging.Debug("Recovered from panic while loading event reactions", r)

		reset()
	}()

	// Load reaction data from config
	var reactionCfgs []eventReactionConfig
	if err := json.Unmarshal(cfg.EventReactions, &reactionCfgs); err != nil {
		panic(err)
	}

	// Add default reactions
	for _, reactionType := range types.DefaultReactions {
		reactionCfgs = append(reactionCfgs, eventReactionConfig{Type: reactionType})
	}

	// Load custom reactions along with defaults
	var numLoadedReactions int
	for _, reactionCfg := range reactionCfgs {
		if reaction, err := e.constructReaction(session, reactionCfg); err == nil {
			e.reactions = append(e.reactions, reaction)
		}
	}

	numCustomReactions := numLoadedReactions - len(types.DefaultReactions)
	if isFirstLoad && numCustomReactions > 0 {
		e.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.EventReactions{
			NumReactions: numCustomReactions,
		})
	}

	return nil
}

func (e *eventReactionsService) constructReaction(session *discordgo.Session, reactionCfg eventReactionConfig) (types.EventReaction, error) {
	var r types.EventReaction
	var err error

	switch reactionCfg.Type {
	// Default reactions
	case types.Snipers:
		r = impl.NewSnipers(e.eventBroker, e.configSvc, session, e.GetCountryCode, e.GetPaymentSrcID, "")
	case types.AFKResponder:
		r = impl.NewAFK(e.eventBroker, e.configSvc, e.variablesSvc, e.IsReactionMessage)
	case types.AntiPhishing:
		r = impl.NewPhishing(e.eventBroker, e.configSvc)
	case types.Notifiers:
		r = impl.NewNotifiers(e.eventBroker, e.configSvc)
	case types.Protections:
		r = impl.NewProtections(e.eventBroker, e.configSvc, e.discordSvc)
	case types.SelfbotDetection:
		r = impl.NewSelfbotDetection(e.eventBroker, e.configSvc)

	// Custom reactions
	case types.IgnoreUser, types.IgnoreGuild:
		r, err = impl.NewIgnoreFromData(session, reactionCfg.Data)
	case types.PingPing, types.PingKick, types.PingBan, types.PingMute, types.PingTimeout, types.PingInsult, types.PingReact, types.PingResponse, types.PingAIResponse:
		r, err = impl.NewPingFromData(e.eventBroker, e.configSvc, session, reactionCfg.Type, reactionCfg.Data)
	case types.MirrorGuild, types.MirrorChannel:
		r, err = impl.NewMirrorFromData(e.eventBroker, session, reactionCfg.Data)
	case types.Responder, types.AIResponder, types.Copy, types.ReplyCopy, types.Mock, types.PinSpam, types.ReactSpam, types.DeleteMessages, types.KeywordNotifier:
		r, err = impl.NewMessageReactionFromData(e.eventBroker, e.configSvc, session, reactionCfg.Type, reactionCfg.Data)
	case types.RemoveReactions:
		r, err = impl.NewRemoveReactionsFromData(session, reactionCfg.Data)
	case types.ForceDisconnect:
		r, err = impl.NewForceDisconnectFromData(session, reactionCfg.Data)
	case types.ForceNick:
		r, err = impl.NewForceNickFromData(session, reactionCfg.Data)
	case types.ForceKick:
		r, err = impl.NewForceKickFromData(session, reactionCfg.Data)
	case types.NoLeave:
		r, err = impl.NewNoLeaveFromData(session, reactionCfg.Data)
	case types.AntiGroupChat:
		r, err = impl.NewAntiGroupChatFromData(session, reactionCfg.Data)
	case types.AntiDeleteMessages:
		r, err = impl.NewAntiDeleteMessagesFromData(e.configSvc, session, reactionCfg.Data)

	default:
		return nil, ErrInvalidReactionType
	}

	return r, err
}

func (e *eventReactionsService) saveEventReactions() error {
	reactionCfgs := make([]eventReactionConfig, 0, len(e.reactions)-len(types.DefaultReactions))

	for _, reaction := range e.reactions {
		if slices.Contains(types.DefaultReactions, reaction.Type()) {
			continue
		}

		reactionCfgs = append(reactionCfgs, eventReactionConfig{
			Type: reaction.Type(),
			Data: reaction.Data(),
		})
	}

	cfg := e.configSvc.GetConfig()

	var err error
	cfg.EventReactions, err = json.Marshal(reactionCfgs)
	if err != nil {
		return err
	}

	return e.configSvc.SaveConfig(cfg)
}

func (e *eventReactionsService) LoadCmdAliases() error {
	cfg := e.configSvc.GetConfig()

	for alias, cmdName := range cfg.Selfbot.CommandAliases {
		cmd := e.routerSvc.GetCommand(cmdName)
		if cmd == nil {
			continue
		}

		cmd.Aliases = append(cmd.Aliases, alias)
	}

	return nil
}
