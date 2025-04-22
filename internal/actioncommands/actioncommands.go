package actioncommands

import (
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actions"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// CommandAction is an action for an ActionCommand.
type CommandAction interface {
	ToSubCommand(services *services.Services) *router.Command
}

// ActionCommand is a command that executes actions.
type ActionCommand struct {
	Name         string
	ReadableName string
	Description  string
	Flags        []router.CommandFlag
	Actions      []CommandAction
}

func (a ActionCommand) ToCommand(services *services.Services) *router.Command {
	cmd := &router.Command{
		Name:        a.Name,
		Description: a.Description,
		Flags:       a.Flags,
	}

	var values []string
	var valueSet *actions.ValueSet

	for _, action := range a.Actions {
		switch action := action.(type) {
		case *actions.ValueGetValues:
			values = action.ValuesGetter()
		case *actions.ValueSet:
			valueSet = action
		}

		cmd.SubCommands = append(cmd.SubCommands, action.ToSubCommand(services))
	}

	if valueSet != nil {
		valueSet.Values = values
	}

	return cmd
}
