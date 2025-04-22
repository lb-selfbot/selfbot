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
		Handler: func(ctx *router.Ctx) error {
			value, err := a.Getter(ctx)
			if err != nil {
				return fmt.Errorf("error getting %s: %w", a.ReadableName, err)
			}

			valueStr, ok := value.(string)
			if value == nil || (ok && len(valueStr) == 0) {
				services.Embed.New(ctx).
					SetDescription("%s is currently empty", a.ReadableName).
					Send()
				return nil
			}

			// Display the retrieved value
			services.Embed.New(ctx).
				SetDescription("%s is currently set to: %v", a.ReadableName, value).
				Send()
			return nil
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
	Values       []string
	Setter       ValueSetFunc
}

func (a *ValueSet) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "set",
		Description: a.Desc,
		Arguments:   "[value]",
		Handler: func(ctx *router.Ctx) error {
			value, err := actionvalue.Parse(ctx, a.ValueType)
			if err != nil {
				return fmt.Errorf("error parsing value: %w", err)
			}

			if a.ValueType == actionvalue.TypeString && len(a.Values) > 0 {
				valueStr := value.(string)

				matchingValues := fn.Filter(a.Values, func(v string) bool {
					return strings.EqualFold(v, valueStr)
				})
				if len(matchingValues) == 0 {
					return fmt.Errorf("value %q is not valid, valid values are: %s", valueStr, strings.Join(a.Values, ", "))
				}

				value = matchingValues[0]
			}

			if err := a.Setter(ctx, value); err != nil {
				return fmt.Errorf("error setting value: %w", err)
			}

			valueStr, ok := value.(string)
			if value == nil || (ok && len(valueStr) == 0) {
				services.Embed.New(ctx).
					SetDescription("%s set to empty", a.ReadableName).
					Send()
				return nil
			}

			services.Embed.New(ctx).
				SetDescription("%s set to %v", a.ReadableName, value).
				Send()
			return nil
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
		Handler: func(ctx *router.Ctx) error {
			if err := a.Toggler(ctx); err != nil {
				return fmt.Errorf("error toggling %s: %w", a.ReadableName, err)
			}

			services.Embed.New(ctx).
				SetDescription("%s toggled", a.ReadableName).
				Send()
			return nil
		},
	}
}

// ValueGetValuesFunc is a function that retrieves the values of the value.
type ValueGetValuesFunc func() []string

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
		Handler: func(ctx *router.Ctx) error {
			values := a.ValuesGetter()
			services.Embed.New(ctx).
				AddField(a.ReadableName, strings.Join(values, "\n")).
				Send()
			return nil
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
		Handler: func(ctx *router.Ctx) error {
			if err := a.ShowExampleFunc(ctx); err != nil {
				return fmt.Errorf("error showing example of %s: %w", a.ReadableName, err)
			}
			return nil
		},
	}
}
