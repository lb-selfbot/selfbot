package actions

import (
	"fmt"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// ValueGetFunc is a function that retrieves the value of the value.
type ValueGetFunc func(ctx *router.Ctx) (any, error)

// ValueGet is an action for an ActionCommand that
// gets a value.
type ValueGet struct {
	ReadableName string
	Desc         string
	Getter       ValueGetFunc
}

func (a *ValueGet) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "get",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			value, err := a.Getter(ctx)
			if err != nil {
				ctx.RespondText("Error getting value: %v", err)
				return
			}

			valueStr, ok := value.(string)
			if value == nil || (ok && len(valueStr) == 0) {
				services.Embed.New(ctx).
					SetDescription("%s is currently empty", a.ReadableName).
					Send()
				return
			}

			// Display the retrieved value
			services.Embed.New(ctx).
				SetDescription("%s is currently set to: %v", a.ReadableName, value).
				Send()
		},
	}
}

// ValueSetFunc is a function that sets the value of the value.
type ValueSetFunc func(ctx *router.Ctx, value any) error

// ValueSet is an action for an ActionCommand that
// sets a value.
type ValueSet struct {
	ReadableName string
	Desc         string
	ValueType    actionvalue.Type
	Setter       ValueSetFunc
}

func (a *ValueSet) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "set",
		Description: a.Desc,
		Arguments:   "[value]",
		Handler: func(ctx *router.Ctx) {
			value, err := actionvalue.Parse(ctx, a.ValueType)
			if err != nil {
				ctx.RespondText("Error parsing value: %v", err)
				return
			}

			if err := a.Setter(ctx, value); err != nil {
				ctx.RespondText("Error setting value: %v", err)
				return
			}

			valueStr, ok := value.(string)
			if value == nil || (ok && len(valueStr) == 0) {
				services.Embed.New(ctx).
					SetDescription("%s set to empty", a.ReadableName).
					Send()
				return
			}

			services.Embed.New(ctx).
				SetDescription("%s set to %v", a.ReadableName, value).
				Send()
		},
	}
}

// ValueToggleFunc is a function that toggles the value of the value.
type ValueToggleFunc func(ctx *router.Ctx) error

// ValueToggle is an action for an ActionCommand that
// toggles a boolean value.
type ValueToggle struct {
	ReadableName string
	Desc         string
	Toggler      ValueToggleFunc
}

func (a *ValueToggle) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "toggle",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			if err := a.Toggler(ctx); err != nil {
				ctx.RespondText("Error toggling value: %v", err)
				return
			}

			services.Embed.New(ctx).
				SetDescription("%s toggled", a.ReadableName).
				Send()
		},
	}
}

// ValueGetValuesFunc is a function that retrieves the values of the value.
type ValueGetValuesFunc func() []any

// ValueGetValues is an action for an ActionCommand that
// lists available values.
type ValueGetValues struct {
	ReadableName string
	Desc         string
	ValuesGetter ValueGetValuesFunc
}

func (a *ValueGetValues) ToSubCommand(services *services.Services) *router.Command {
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

// ValueShowExampleFunc is a function that demonstrates a feature.
type ValueShowExampleFunc func(ctx *router.Ctx) error

// ValueShowExample is an action for an ActionCommand that
// demonstrates a feature.
type ValueShowExample struct {
	ReadableName    string
	Desc            string
	ShowExampleFunc ValueShowExampleFunc
}

func (a *ValueShowExample) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "example",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			if err := a.ShowExampleFunc(ctx); err != nil {
				ctx.RespondText("Error demonstrating example: %v", err)
				return
			}
		},
	}
}
