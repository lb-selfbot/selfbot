package events

import (
	"github.com/LightningDev1/dgc"

	"github.com/LightningDev1/LB-Selfbot/ipc"
)

func (e Events) OnCommand(ctx *dgc.Ctx) {
	ipc.SendNotification("Command " + ctx.Command.Name + " used")
	ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
}
