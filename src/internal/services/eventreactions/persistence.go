package eventreactions

import (
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventreactions/impl"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
)

var (
	ErrInvalidReactionType = errors.New("invalid reaction type")
)

type eventReactionConfig struct {
	Type interfaces.EventReactionType `json:"Type"`
	Data map[string]any               `json:"Data"`
}

func (e *EventReactionsServiceImpl) startLoading() {
	timer := time.NewTimer(5 * time.Second)
	isFirstIteration := true

	for range timer.C {
		if err := e.loadEventReactions(isFirstIteration); err != nil {
			e.uiSvc.LogError("Failed to load event reactions:", err)
		}

		if err := e.LoadCmdAliases(); err != nil {
			e.uiSvc.LogError("Failed to load command aliases:", err)
		}

		if isFirstIteration {
			isFirstIteration = false
		}
	}
}

func (e *EventReactionsServiceImpl) loadEventReactions(isFirstLoad bool) error {
	e.reactionsMu.Lock()
	defer e.reactionsMu.Unlock()

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r == nil {
			return
		}

		e.uiSvc.LogError("Failed to load event reactions. Event reactions have been reset.")

		cfg.EventReactions = []byte("[]")
		_ = cfg.Save()
	}()

	// Load reaction data from config
	var reactionCfgs []eventReactionConfig
	if err = json.Unmarshal(cfg.EventReactions, &reactionCfgs); err != nil {
		panic(err)
	}

	// Add default reactions
	for _, reactionType := range interfaces.DefaultEventReactions {
		reactionCfgs = append(reactionCfgs, eventReactionConfig{Type: reactionType})
	}

	// Load custom reactions along with defaults
	var numLoadedReactions int
	for _, reactionCfg := range reactionCfgs {
		if reaction, err := e.constructReaction(reactionCfg); err == nil {
			e.reactions = append(e.reactions, reaction)
		}
	}

	numCustomReactions := numLoadedReactions - len(interfaces.DefaultEventReactions)
	if isFirstLoad && numCustomReactions > 0 {
		e.uiSvc.LogEventReactions(numCustomReactions)
	}

	return nil
}

func (e *EventReactionsServiceImpl) constructDefaultReactions() []interfaces.EventReaction {
	reactions := make([]interfaces.EventReaction, 0, len(interfaces.DefaultEventReactions))

	for _, reactionType := range interfaces.DefaultEventReactions {
		reaction, err := e.constructReaction(eventReactionConfig{Type: reactionType})
		if err != nil {
			continue
		}

		reactions = append(reactions, reaction)
	}

	return reactions
}

func (e *EventReactionsServiceImpl) constructReaction(reactionCfg eventReactionConfig) (interfaces.EventReaction, error) {
	session := e.discordSvc.GetSession()

	var r interfaces.EventReaction
	var err error

	switch reactionCfg.Type {
	// Default reactions
	case interfaces.EventReactionSnipers:
		r = impl.NewSnipers(e.uiSvc, e, session, "")
	case interfaces.EventReactionAFKResponder:
		r = impl.NewAFK(e.uiSvc, e, e.variablesSvc)
	case interfaces.EventReactionAntiPhishing:
		r = impl.NewPhishing(e.uiSvc)
	case interfaces.EventReactionNotifiers:
		r = impl.NewNotifiers(e.uiSvc)
	case interfaces.EventReactionProtections:
		r = impl.NewProtections(e.uiSvc, e.embedSvc)
	case interfaces.EventReactionSelfbotDetection:
		r = impl.NewSelfbotDetection(e.uiSvc, e)

	// Custom reactions
	case interfaces.EventReactionIgnoreUser, interfaces.EventReactionIgnoreGuild:
		r, err = impl.NewIgnoreFromData(session, reactionCfg.Data)
	case interfaces.EventReactionPingPing, interfaces.EventReactionPingKick, interfaces.EventReactionPingBan, interfaces.EventReactionPingMute, interfaces.EventReactionPingTimeout, interfaces.EventReactionPingInsult, interfaces.EventReactionPingReact, interfaces.EventReactionPingResponse, interfaces.EventReactionPingAIResponse:
		r, err = impl.NewPingFromData(e, session, reactionCfg.Type, reactionCfg.Data)
	case interfaces.EventReactionMirrorGuild, interfaces.EventReactionMirrorChannel:
		r, err = impl.NewMirrorFromData(e, session, reactionCfg.Data)
	case interfaces.EventReactionResponder, interfaces.EventReactionAIResponder, interfaces.EventReactionCopy, interfaces.EventReactionReplyCopy, interfaces.EventReactionMock, interfaces.EventReactionPinSpam, interfaces.EventReactionReactSpam, interfaces.EventReactionDeleteMessages, interfaces.EventReactionKeywordNotifier:
		r, err = impl.NewMessageReactionFromData(session, reactionCfg.Type, reactionCfg.Data)
	case interfaces.EventReactionRemoveReactions:
		r, err = impl.NewRemoveReactionsFromData(session, reactionCfg.Data)
	case interfaces.EventReactionForceDisconnect:
		r, err = impl.NewForceDisconnectFromData(session, reactionCfg.Data)
	case interfaces.EventReactionForceNick:
		r, err = impl.NewForceNickFromData(session, reactionCfg.Data)
	case interfaces.EventReactionForceKick:
		r, err = impl.NewForceKickFromData(session, reactionCfg.Data)
	case interfaces.EventReactionNoLeave:
		r, err = impl.NewNoLeaveFromData(session, reactionCfg.Data)
	case interfaces.EventReactionAntiGroupChat:
		r, err = impl.NewAntiGroupChatFromData(session, reactionCfg.Data)
	case interfaces.EventReactionAntiDeleteMessages:
		r, err = impl.NewAntiDeleteMessagesFromData(session, reactionCfg.Data)

	default:
		return nil, ErrInvalidReactionType
	}

	return r, err
}

func (e *EventReactionsServiceImpl) saveEventReactions() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	reactionCfgs := make([]eventReactionConfig, 0, len(e.reactions)-len(interfaces.DefaultEventReactions))

	for _, reaction := range e.reactions {
		if slices.Contains(interfaces.DefaultEventReactions, reaction.Type()) {
			continue
		}

		reactionCfgs = append(reactionCfgs, eventReactionConfig{
			Type: reaction.Type(),
			Data: reaction.Data(),
		})
	}

	cfg.EventReactions, err = json.Marshal(reactionCfgs)
	if err != nil {
		return err
	}

	return cfg.Save()
}

func (e *EventReactionsServiceImpl) LoadCmdAliases() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	router := e.discordSvc.GetRouter()

	for cmdName, alias := range cfg.Selfbot.CommandAliases {
		cmd := router.GetCmd(cmdName)
		if cmd == nil {
			continue
		}

		cmd.Aliases = append(cmd.Aliases, alias)
	}

	return nil
}
