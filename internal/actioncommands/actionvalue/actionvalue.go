package actionvalue

import (
	"fmt"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// Type defines the type of a value for an ActionCommand.
type Type int

const (
	TypeString Type = iota
	TypeInt
	TypePercentage
	TypeBool
	TypeDuration
)

// Parse parses a value from the context based on its type.
func Parse(ctx *router.Ctx, valueType Type) (any, error) {
	switch valueType {
	case TypeString:
		return ctx.Arguments.Raw(), nil
	case TypeInt:
		return ctx.Arguments.Get(0).Int()
	case TypePercentage:
		value, err := ctx.Arguments.Get(0).Int()
		if err != nil {
			return nil, err
		}
		return fn.Clamp(value, 0, 100), nil
	case TypeBool:
		return ctx.Arguments.Get(0).Bool()
	case TypeDuration:
		duration, err := ctx.Arguments.Get(0).Duration()
		if err != nil {
			return nil, err
		}
		return int(duration.Seconds()), nil
	default:
		return nil, fmt.Errorf("unknown value type: %d", valueType)
	}
}
