package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/encodings"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type TrollingCommands struct{}

func (t TrollingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Trolling", "Trolling commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "purgehack",
		Description: "Purge a chat without permissions",
		Usage:       "[p]purgehack",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.purgehack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ghostping",
		Description: "Ghostping a user",
		Usage:       "[p]ghostping <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.ghostping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gprole",
		Description: "Ghostping a role",
		Usage:       "[p]gprole <role>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.gprole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "token",
		Description: "Send the first part of someones token",
		Usage:       "[p]token <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.token,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trollnitro",
		Description: "Send a fake nitro link with an image",
		Usage:       "[p]trollnitro",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.trollnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakenitro",
		Description: "Generate fake nitro codes",
		Usage:       "[p]fakenitro <amount>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.fakenitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddentext",
		Description: "Hide text behind other text",
		Usage:       "[p]hiddentext <\"text\"> <\"hidden\">",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.hiddentext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeninvite",
		Description: "Hide an invite behind text",
		Usage:       "[p]hiddeninvite <invite> <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.hiddeninvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenping",
		Description: "Hide a ping behind text",
		Usage:       "[p]hiddenping <userid> <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.hiddenping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenroleping",
		Description: "Hide a role ping behind text",
		Usage:       "[p]hiddenroleping <roleid> <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.hiddenroleping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeneveryone",
		Description: "Hide an everyone ping behind text",
		Usage:       "[p]hiddeneveryone <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.hiddeneveryone,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenhere",
		Description: "Hide a here ping behind text",
		Usage:       "[p]hiddenhere <text>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.hiddenhere,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "voicemove",
		Description: "Move to random voice channels",
		Usage:       "[p]voicemove <amount>",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.voicemove,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggertyping",
		Description: "Fake \"User is typing\" message",
		Usage:       "[p]triggertyping <seconds> [channel]",
		Flags:       []string{"DANGEROUS"},
		Handler:     t.triggertyping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "impersonate",
		Description: "Impersonate a user using a webhook",
		Usage:       "[p]impersonate <user> <text>",
		Handler:     t.impersonate,
	})
}

func (TrollingCommands) purgehack(ctx *dgc.Ctx) {
	_ = ctx.RespondText("ﾠ" + strings.Repeat("\n", 1998) + "ﾠ")
}

func (TrollingCommands) ghostping(*dgc.Ctx) {
	// noop
}

func (TrollingCommands) gprole(*dgc.Ctx) {
	// noop
}

func (TrollingCommands) token(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	_ = ctx.RespondText(encodings.Base64(user.ID, false))
}

func (TrollingCommands) trollnitro(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("discord\u2024gift/%s %s https://i.imgur.com/0WSpeC1.png", utils.Random.String(16), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) fakenitro(ctx *dgc.Ctx) {
	amount, err := ctx.Arguments.Get(0).AsInt()
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}
	for i := 0; i < amount; i++ {
		_ = ctx.RespondText("https://discord.gift/" + utils.Random.String(16))
	}
}

func (TrollingCommands) hiddentext(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"%s %s %s",
		utils.Arguments["<\"text\">"],
		constants.INVISIBLE_TEXT,
		utils.Arguments["<\"hidden\">"],
	))
}

func (TrollingCommands) hiddeninvite(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.INVISIBLE_TEXT, utils.Arguments["<invite>"]))
}

func (TrollingCommands) hiddenping(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s <@%s>", text, constants.INVISIBLE_TEXT, utils.Arguments["<userid>"]))
}

func (TrollingCommands) hiddenroleping(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s <@&%s>", text, constants.INVISIBLE_TEXT, utils.Arguments["<roleid>"]))
}

func (TrollingCommands) hiddeneveryone(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("%s %s @everyone", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) hiddenhere(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("%s %s @here", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) voicemove(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Event.GuildID)
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	if conn, ok := ctx.Session.VoiceConnections[ctx.Event.GuildID]; ok {
		_ = conn.Disconnect()
	}

	channels := utils.Filter(guild.Channels, func(channel *discordgo.Channel) bool {
		if channel.Type != discordgo.ChannelTypeGuildVoice {
			return false
		}

		perms, err := ctx.Session.UserChannelPermissions(ctx.Session.State.User.ID, channel.ID)
		if err != nil {
			return false
		}

		if perms&discordgo.PermissionVoiceConnect == 0 {
			return false
		}

		return true
	})

	if len(channels) == 0 {
		_ = ctx.RespondText("No voice channels found")
		return
	}

	i := 0

	for {
		if i >= amount {
			break
		}

		channel := utils.RandomChoice(channels)
		currentChannel := ctx.Session.VoiceConnections[ctx.Event.GuildID]

		if currentChannel != nil && currentChannel.ChannelID == channel.ID {
			continue
		}

		_, err := ctx.Session.ChannelVoiceJoin(ctx.Event.GuildID, channel.ID, false, true)
		if err != nil {
			_ = ctx.RespondText("Error joining voice channel: " + err.Error())
			return
		}

		time.Sleep(time.Second * 1)

		i++
	}

	if voice, ok := ctx.Session.VoiceConnections[ctx.Event.GuildID]; ok {
		_ = voice.Disconnect()
	}
}

func (TrollingCommands) triggertyping(ctx *dgc.Ctx) {
	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	utils.Discord.SleepWithTypingFix(ctx.Session, channel.ID, time.Duration(seconds)*time.Second)
}

func (TrollingCommands) impersonate(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	_, err := utils.Discord.Impersonate(ctx.Session, user, text, ctx.Event.ChannelID)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
}
