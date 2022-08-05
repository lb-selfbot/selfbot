package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"

	"github.com/LightningDev1/LB-Selfbot/ipc"
)

func (e Events) OnCommand(following dgc.ExecutionHandler) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		arguments, success := utils.Discord.GetArguments(ctx)
		if success {
			ipc.SendNotification("Command " + ctx.Command.Name + " used")
			_ = ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
			utils.Arguments = arguments
			following(ctx)
		}
	}
}
