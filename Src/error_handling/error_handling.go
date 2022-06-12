package error_handling

import (
	"fmt"

	"github.com/lus/dgc"
)

type CommandError int64

const (
	MissingRequiredArgument CommandError = 0
	ServerOnlyCommand       CommandError = 1
	InvalidArgument         CommandError = 2
)

func OnError(ctx *dgc.Ctx, err CommandError, text string) {
	switch err {
	case MissingRequiredArgument:
		ctx.RespondText("Missing required argument: " + text)
	case ServerOnlyCommand:
		ctx.RespondText(fmt.Sprintf("The `%s` is only available in a server.", text))
	}
}
