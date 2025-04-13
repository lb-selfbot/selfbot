package actions

import (
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// ListGetFunc is a function that retrieves the value of the list.
type ListGetFunc func(ctx *router.Ctx) error

// ListGet is an action for an ActionCommand that
// gets a list of values.
type ListGet struct {
	ReadableName string
	Desc         string
	Getter       ListGetFunc
}

func (a *ListGet) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "get",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			if err := a.Getter(ctx); err != nil {
				ctx.RespondText("Error getting %s: %v", a.ReadableName, err)
				return
			}
		},
	}
}

// ListAddFunc is a function that adds a value to the list.
type ListAddFunc func(ctx *router.Ctx) error

// ListAdd is an action for an ActionCommand that
// adds a value to the list.
type ListAdd struct {
	ReadableName string
	Desc         string
	Arguments    string
	AddFunc      ListAddFunc
}

func (a *ListAdd) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "add",
		Description: a.Desc,
		Arguments:   a.Arguments,
		Handler: func(ctx *router.Ctx) {
			if err := a.AddFunc(ctx); err != nil {
				ctx.RespondText("Error adding to %s: %v", a.ReadableName, err)
				return
			}
		},
	}
}

// ListRemoveFunc is a function that removes a value from the list.
type ListRemoveFunc func(ctx *router.Ctx) error

// ListRemove is an action for an ActionCommand that
// removes a value from the list.
type ListRemove struct {
	ReadableName string
	Desc         string
	Arguments    string
	RemoveFunc   ListRemoveFunc
}

func (a *ListRemove) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "remove",
		Description: a.Desc,
		Arguments:   a.Arguments,
		Handler: func(ctx *router.Ctx) {
			if err := a.RemoveFunc(ctx); err != nil {
				ctx.RespondText("Error removing from %s: %v", a.ReadableName, err)
				return
			}
		},
	}
}

// ListClearFunc is a function that clears the list.
type ListClearFunc func(ctx *router.Ctx) error

// ListClear is an action for an ActionCommand that
// clears the list.
type ListClear struct {
	ReadableName string
	Desc         string
	ClearFunc    ListClearFunc
}

func (a *ListClear) ToSubCommand(services *services.Services) *router.Command {
	return &router.Command{
		Name:        "clear",
		Description: a.Desc,
		Handler: func(ctx *router.Ctx) {
			if err := a.ClearFunc(ctx); err != nil {
				ctx.RespondText("Error clearing %s: %v", a.ReadableName, err)
				return
			}
		},
	}
}
