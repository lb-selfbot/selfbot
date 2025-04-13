package actioncommands

import (
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actions"
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// Builder helps with building an ActionCommand.
type Builder struct {
	cmd ActionCommand
}

// NewBuilder creates a new builder for an ActionCommandSingle.
func NewBuilder(name, readableName string) *Builder {
	return &Builder{
		cmd: ActionCommand{
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

// WithValueGetter adds a getter action to the command.
func (b *Builder) WithValueGetter(desc string, getter actions.ValueGetFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ValueGet{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Getter:       getter,
	})
	return b
}

// WithValueSetter adds a setter action to the command.
func (b *Builder) WithValueSetter(desc string, valueType actionvalue.Type, setter actions.ValueSetFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ValueSet{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ValueType:    valueType,
		Setter:       setter,
	})
	return b
}

// WithValueToggler adds a toggler action to the command.
func (b *Builder) WithValueToggler(desc string, toggler func(ctx *router.Ctx) error) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ValueToggle{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Toggler:      toggler,
	})
	return b
}

// WithValues adds a values action to the command.
func (b *Builder) WithValues(desc string, valuesGetter actions.ValueGetValuesFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ValueGetValues{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ValuesGetter: valuesGetter,
	})
	return b
}

// WithListGetter adds a list getter action to the command.
func (b *Builder) WithListGetter(desc string, getter actions.ListGetFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ListGet{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Getter:       getter,
	})
	return b
}

// WithListAdder adds a list adder action to the command.
func (b *Builder) WithListAdder(desc, args string, adder actions.ListAddFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ListAdd{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Arguments:    args,
		AddFunc:      adder,
	})
	return b
}

// WithListRemover adds a list remover action to the command.
func (b *Builder) WithListRemover(desc, args string, remover actions.ListRemoveFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ListRemove{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		Arguments:    args,
		RemoveFunc:   remover,
	})
	return b
}

// WithListClearer adds a list clearer action to the command.
func (b *Builder) WithListClearer(desc string, clearer actions.ListClearFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ListClear{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ClearFunc:    clearer,
	})
	return b
}

// WithExample adds an example action to the command.
func (b *Builder) WithExample(desc string, example actions.ValueShowExampleFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.ValueShowExample{
		ReadableName:    b.cmd.ReadableName,
		Desc:            desc,
		ShowExampleFunc: example,
	})
	return b
}

// WithEventReactionList adds an event reaction list action to the command.
func (b *Builder) WithEventReactionList(reactionType types.EventReactionType, desc string) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.EventReactionList{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ReactionType: reactionType,
	})
	return b
}

// WithEventReactionClear adds an event reaction clear action to the command.
func (b *Builder) WithEventReactionClear(reactionType types.EventReactionType, desc string) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.EventReactionClear{
		ReadableName: b.cmd.ReadableName,
		Desc:         desc,
		ReactionType: reactionType,
	})
	return b
}

// WithEventReactionToggle adds an event reaction toggle action to the command.
func (b *Builder) WithEventReactionToggle(desc, args string, getReactionFunc actions.GetEventReactionFunc) *Builder {
	b.cmd.Actions = append(b.cmd.Actions, &actions.EventReactionToggle{
		ReadableName:    b.cmd.ReadableName,
		Desc:            desc,
		Arguments:       args,
		GetReactionFunc: getReactionFunc,
	})
	return b
}

// Build creates an ActionCommandSingle from the builder.
func (b *Builder) Build(services *services.Services) *router.Command {
	return b.cmd.ToCommand(services)
}
