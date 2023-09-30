package eventhandlers

import (
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/dgc"
)

func (e *Events) OnCommand(following dgc.ExecutionHandler) dgc.ExecutionHandler {
	return func(ctx *dgc.Ctx) {
		arguments, success := utils.Discord.GetArguments(ctx)
		if !success {
			return
		}

		cfg, err := config.GetConfig()
		if err != nil {
			logging.UI.Error("Failed to get config")
			return
		}

		if cfg.Selfbot.DeleteCommands {
			_ = ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
		}

		utils.Arguments = arguments

		e.Settings.UsedCommands = append(e.Settings.UsedCommands, ctx)

		logging.Logging.CommandUsed(ctx.Command.Name)

		websocket.Websocket.CheckNonce()

		following(ctx)
	}
}
