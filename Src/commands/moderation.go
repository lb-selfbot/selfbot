package commands

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type ModerationCommands struct{}

func (m ModerationCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "channelinfo",
		Description: "Get info about a channel",
		Usage:       "[p]channelinfo <#channel>",
		Handler:     m.channelinfo,
	})
}

func (ModerationCommands) channelinfo(ctx *dgc.Ctx) {
	channel, success := utils.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find channel")
		return
	}

	embed.NewEmbed().
		AddField("Name", channel.Name).
		AddField("ID", channel.ID).
		AddField("Topic", channel.Topic).
		AddField("Category ID", channel.ParentID).
		AddField("Is NSFW", utils.FormatBool(channel.NSFW)).
		AddField("Created At", utils.SnowflakeToTime(channel.ID).Format("Monday, 02. January 2006 @ 15:04:05")).
		AddField("Slowmode Delay", fmt.Sprint(channel.RateLimitPerUser)+" seconds").
		AddField("Position", fmt.Sprint(channel.Position)).
		AddField("Last Message ID", channel.LastMessageID).
		Send(ctx)
}
