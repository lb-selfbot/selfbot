package actioncommands

import (
	"fmt"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// CommandAction is an action for an ActionCommand.
type CommandAction interface {
	ToSubCommand(services *services.Services) *router.Command
}

// CommandActionGet is an action for an ActionCommand that
// gets a value.
type CommandActionGet struct {
	ReadableName string
	Desc         string
	Getter       func(ctx *router.Ctx) (any, error)
}

func (a *CommandActionGet) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "get",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			value, err := a.Getter(ctx)
			if err != nil {
				ctx.RespondText("Error getting value: %v", err)
				return
			}

			// Display the retrieved value
			services.Embed.New(ctx).
				SetDescription("%s", fmt.Sprintf("%s is currently set to: %v", a.ReadableName, value)).
				Send()
		},
	}
}

// CommandActionSet is an action for an ActionCommand that
// sets a value.
type CommandActionSet struct {
	ReadableName string
	Desc         string
	ValueType    ActionValueType
	Setter       func(ctx *router.Ctx, value any) error
}

func (a *CommandActionSet) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "set",
		Description: a.Desc,
		Arguments:   "[value]",
		Handler: func(ctx *router.Ctx) {
			value, err := parseValue(ctx, a.ValueType)
			if err != nil {
				ctx.RespondText("Error parsing value: %v", err)
				return
			}

			if err := a.Setter(ctx, value); err != nil {
				ctx.RespondText("Error setting value: %v", err)
				return
			}

			services.Embed.New(ctx).
				SetDescription("%s", fmt.Sprintf("%s set to %v", a.ReadableName, value)).
				Send()
		},
	}
}

// CommandActionToggle is an action for an ActionCommand that
// toggles a boolean value.
type CommandActionToggle struct {
	ReadableName string
	Desc         string
	Toggler      func(ctx *router.Ctx) error
}

func (a *CommandActionToggle) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "toggle",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			if err := a.Toggler(ctx); err != nil {
				ctx.RespondText("Error toggling value: %v", err)
				return
			}

			services.Embed.New(ctx).
				SetDescription("%s", fmt.Sprintf("%s toggled", a.ReadableName)).
				Send()
		},
	}
}

// CommandActionValues is an action for an ActionCommand that
// lists available values.
type CommandActionValues struct {
	ReadableName string
	Desc         string
	ValuesGetter func() []any
}

func (a *CommandActionValues) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "values",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			values := a.ValuesGetter()
			stringValues := fn.Map(values, func(value any) string {
				return fmt.Sprint(value)
			})

			services.Embed.New(ctx).
				AddField(a.ReadableName, strings.Join(stringValues, "\n")).
				Send()
		},
	}
}

// CommandActionExample is an action for an ActionCommand that
// demonstrates a feature.
type CommandActionExample struct {
	ReadableName string
	Desc         string
	Example      func(ctx *router.Ctx) error
}

func (a *CommandActionExample) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "example",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			if err := a.Example(ctx); err != nil {
				ctx.RespondText("Error demonstrating example: %v", err)
				return
			}
		},
	}
}

// ActionValueType defines the type of a value for an ActionCommand.
type ActionValueType int

const (
	ValueTypeString ActionValueType = iota
	ValueTypeInt
	ValueTypePercentage
	ValueTypeBool
	ValueTypeDuration
)

// parseValue parses a value from the context based on its type.
func parseValue(ctx *router.Ctx, valueType ActionValueType) (any, error) {
	switch valueType {
	case ValueTypeString:
		return ctx.Arguments.Raw(), nil
	case ValueTypeInt:
		return ctx.Arguments.Get(0).Int()
	case ValueTypePercentage:
		value, err := ctx.Arguments.Get(0).Int()
		if err != nil {
			return nil, err
		}
		return fn.Clamp(value, 0, 100), nil
	case ValueTypeBool:
		return ctx.Arguments.Get(0).Bool()
	case ValueTypeDuration:
		duration, err := ctx.Arguments.Get(0).Duration()
		if err != nil {
			return nil, err
		}
		return int(duration.Seconds()), nil
	default:
		return nil, fmt.Errorf("unknown value type: %d", valueType)
	}
}

// ActionCommandSingle is a command that executes actions
// on a single value.
type ActionCommandSingle struct {
	Name         string
	ReadableName string
	Description  string
	Flags        []router.CommandFlag
	Actions      []CommandAction
}

func (a ActionCommandSingle) ToCommand(services *services.Services) *router.Command {
	cmd := &router.Command{
		Name:        a.Name,
		Description: a.Description,
		Flags:       a.Flags,
	}

	for _, action := range a.Actions {
		cmd.SubCommands = append(cmd.SubCommands, action.ToSubCommand(services))
	}

	return cmd
}

// Builder helps with building an ActionCommandSingle.
type Builder struct {
	cmd ActionCommandSingle
}

// NewBuilder creates a new builder for an ActionCommandSingle.
func NewBuilder(name, readableName string) *Builder {
	return &Builder{
		cmd: ActionCommandSingle{
			Name:         name,
			ReadableName: readableName,
			Flags:        []router.CommandFlag{router.CmdFlagDangerous},
		},
	}
}

// WithDescription adds a description to the command.
func (b *Builder) WithDescription(description string) *Builder {
	b.cmd.Description = description
	return b
}

// WithFlags adds flags to the command.
func (b *Builder) WithFlags(flags ...router.CommandFlag) *Builder {
	b.cmd.Flags = flags
	return b
}

// WithGetter adds a getter action to the command.
func (b *Builder) WithGetter(desc string, getter func(ctx *router.Ctx) (any, error)) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &CommandActionGet{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Getter:       getter,
	})
	return b
}

// WithSetter adds a setter action to the command.
func (b *Builder) WithSetter(desc string, valueType ActionValueType, setter func(ctx *router.Ctx, value any) error) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &CommandActionSet{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ValueType:    valueType,
		Setter:       setter,
	})
	return b
}

// WithToggler adds a toggler action to the command.
func (b *Builder) WithToggler(desc string, toggler func(ctx *router.Ctx) error) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &CommandActionToggle{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Toggler:      toggler,
	})
	return b
}

// WithValues adds a values action to the command.
func (b *Builder) WithValues(desc string, valuesGetter func() []any) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &CommandActionValues{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ValuesGetter: valuesGetter,
	})
	return b
}

// WithExample adds an example action to the command.
func (b *Builder) WithExample(desc string, example func(ctx *router.Ctx) error) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &CommandActionExample{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Example:      example,
	})
	return b
}

// Build creates an ActionCommandSingle from the builder.
func (b *Builder) Build(services *services.Services) *router.Command {
	return b.cmd.ToCommand(services)
}
