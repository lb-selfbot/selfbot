package actions

import (
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// CustomAction is an action for an ActionCommand that
// performs a custom action.
type CustomAction struct {
	ActionName string
	Desc       string
	Arguments  string
	Handler    router.ExecutionHandler
}

func (a *CustomAction) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        a.ActionName,
		Description: a.Desc,
		Arguments:   a.Arguments,
		Handler:     a.Handler,
	}
}
