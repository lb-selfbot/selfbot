package events

import (
	"github.com/LightningDev1/dgc"

	"github.com/LightningDev/LB-Selfbot/utils"
)

func (e Events) OnCommand(ctx *dgc.Ctx) {
	utils.SendNotification("Command " + ctx.Command.Name + " used")
	ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
}
