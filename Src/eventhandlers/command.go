package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

func (e *Events) OnCommand(following dgc.ExecutionHandler) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		arguments, success := utils.Discord.GetArguments(ctx)
		if !success {
			return
		}

		_ = ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)

		utils.Arguments = arguments

		e.Settings.UsedCommands = append(e.Settings.UsedCommands, ctx)

		logging.Logging.CommandUsed(ctx.Command.Name)

		following(ctx)
	}
}
