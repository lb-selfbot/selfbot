package events

import (
	"github.com/lus/dgc"

	"github.com/LightningDev/LBGo/utils"
)

func (e Events) OnCommand(ctx *dgc.Ctx) {
	utils.SendNotification("Command " + ctx.Command.Name + " used")
	ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
}
