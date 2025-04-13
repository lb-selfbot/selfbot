package actions

import (
	"slices"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// EventReactionList is an action for an ActionCommand that
// lists available event reactions.
type EventReactionList struct {
	ReadableName string
	Desc         string
	ReactionType types.EventReactionType
}

func (a *EventReactionList) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "list",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			reactions := services.EventReactions.GetReactionsByType(a.ReactionType)
			if len(reactions) == 0 {
				services.Embed.New(ctx).
					SetDescription("No %s reactions found", a.ReadableName).
					Send()
				return
			}

			stringReactions := fn.Map(reactions, func(reaction types.EventReaction) string {
				return reaction.Description()
			})

			services.Embed.New(ctx).
				AddField(a.ReadableName, strings.Join(stringReactions, "\n")).
				Send()
		},
	}
}

// EventReactionClear is an action for an ActionCommand that
// clears all event reactions.
type EventReactionClear struct {
	ReadableName string
	Desc         string
	ReactionType types.EventReactionType
}

func (a *EventReactionClear) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "clear",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			services.EventReactions.RemoveReactionsByType(a.ReactionType)

			services.Embed.New(ctx).
				SetDescription("Cleared all %s reactions", a.ReadableName).
				Send()
		},
	}
}

// GetEventReactionFunc is a function that retrieves an event reaction.
type GetEventReactionFunc func(ctx *router.Ctx) (types.EventReaction, error)

// EventReactionToggle is an action for an ActionCommand that
// toggles an event reaction.
type EventReactionToggle struct {
	ReadableName    string
	Desc            string
	Arguments       string
	GetReactionFunc GetEventReactionFunc
}

func (a *EventReactionToggle) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "toggle",
		Description: a.Desc,
		Arguments:   a.Arguments,
		Handler: func(ctx *router.Ctx) {
			reaction, err := a.GetReactionFunc(ctx)
			if err != nil {
				ctx.RespondText("Error getting reaction: %v", err)
				return
			}

			if slices.Contains(services.EventReactions.GetAllReactions(), reaction) {
				services.EventReactions.RemoveReaction(reaction)
				services.Embed.New(ctx).
					SetDescription("Removed event reaction: %s", reaction.Description()).
					Send()
			} else {
				services.EventReactions.AddReaction(reaction)
				services.Embed.New(ctx).
					SetDescription("Added event reaction: %s", reaction.Description()).
					Send()
			}
		},
	}
}
