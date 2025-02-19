package eventreactions

import (
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/eventreactions/impl"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
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

		if err := e.loadCmdAliases(); err != nil {
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

	switch reactionCfg.Type {
	// Default reactions
	case interfaces.EventReactionSnipers:
		return impl.NewSnipers(e.uiSvc, e, session, ""), nil
		// return NewSnipers(e.discordSvc, reactionCfg.Data), nil
	}

	return nil, ErrInvalidReactionType
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

func (e *EventReactionsServiceImpl) loadCmdAliases() error {
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
