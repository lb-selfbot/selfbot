package commands

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/random"
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
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.purgehack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ghostping",
		Description: "Ghostping a user",
		Usage:       "[p]ghostping <user>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.ghostping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gprole",
		Description: "Ghostping a role",
		Usage:       "[p]gprole <role>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.gprole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "token",
		Description: "Send the first part of someones token",
		Usage:       "[p]token <user>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.token,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trollnitro",
		Description: "Send a fake nitro link with an image",
		Usage:       "[p]trollnitro",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.trollnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakenitro",
		Description: "Generate fake nitro codes",
		Usage:       "[p]fakenitro <amount>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.fakenitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddentext",
		Description: "Hide text behind other text",
		Usage:       "[p]hiddentext <\"text\"> <\"hidden\">",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.hiddentext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeninvite",
		Description: "Hide an invite behind text",
		Usage:       "[p]hiddeninvite <invite> <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.hiddeninvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenping",
		Description: "Hide a ping behind text",
		Usage:       "[p]hiddenping <userid> <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.hiddenping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenroleping",
		Description: "Hide a role ping behind text",
		Usage:       "[p]hiddenroleping <roleid> <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.hiddenroleping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeneveryone",
		Description: "Hide an everyone ping behind text",
		Usage:       "[p]hiddeneveryone <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.hiddeneveryone,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenhere",
		Description: "Hide a here ping behind text",
		Usage:       "[p]hiddenhere <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.hiddenhere,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "voicemove",
		Description: "Move to random voice channels",
		Usage:       "[p]voicemove <amount>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.voicemove,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggertyping",
		Description: "Fake \"User is typing\" message",
		Usage:       "[p]triggertyping <seconds> [channel]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.triggertyping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "impersonate",
		Description: "Impersonate a user using a webhook",
		Usage:       "[p]impersonate <user> <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.impersonate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakepings",
		Description: "Get fake pinged in a channel",
		Usage:       "[p]fakepings <amount> [channel]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.fakepings,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "billionpings",
		Description: "Get a billion pings in a channel",
		Usage:       "[p]billionpings [channel]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     t.billionpings,
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

	_ = ctx.RespondText(base64.StdEncoding.EncodeToString([]byte(user.ID)))
}

func (TrollingCommands) trollnitro(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("discord\u2024gift/%s %s https://i.imgur.com/0WSpeC1.png", random.String(16), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) fakenitro(ctx *dgc.Ctx) {
	amount, err := ctx.Arguments.Get(0).AsInt()
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}
	for i := 0; i < amount; i++ {
		_ = ctx.RespondText("https://discord.gift/" + random.String(16))
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

		channel := random.Choice(channels)
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

func (TrollingCommands) fakepings(ctx *dgc.Ctx) {
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	if amount < 1 || amount > 1000000000 {
		_ = ctx.RespondText("Amount must be between 1 and 1000000000")
		return
	}

	channel, success := utils.Discord.TryGetChannelFromMentionWithDefault(ctx, ctx.Arguments.Get(1))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 10)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	if len(messages) == 0 {
		_ = ctx.RespondText("No messages found in channel")
		return
	}

	err = ctx.Session.MessageUnack(channel.ID, messages[0].ID, amount)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
}

func (TrollingCommands) billionpings(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMentionWithDefault(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 10)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	if len(messages) == 0 {
		_ = ctx.RespondText("No messages found in channel")
		return
	}

	err = ctx.Session.MessageUnack(channel.ID, messages[0].ID, 1000000000)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
}
