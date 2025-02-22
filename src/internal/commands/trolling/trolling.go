package trolling

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/dgc"
)

func (t TrollingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Trolling", "Trolling commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "purgehack",
		Description: "Purge a chat without permissions",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.purgehack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ghostping",
		Description: "Ghostping a user",
		Arguments:   "<user>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.ghostping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gprole",
		Description: "Ghostping a role",
		Arguments:   "<role>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.gprole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "token",
		Description: "Send the first part of someones token",
		Arguments:   "<user>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.token,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trollnitro",
		Description: "Send a fake nitro link with an image",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.trollnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakenitro",
		Description: "Generate fake nitro codes",
		Arguments:   "<amount>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.fakenitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddentext",
		Description: "Hide text behind other text",
		Arguments:   "<\"text\"> <\"hidden\">",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.hiddentext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeninvite",
		Description: "Hide an invite behind text",
		Arguments:   "<invite> <text>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.hiddeninvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenping",
		Description: "Hide a ping behind text",
		Arguments:   "<userid> <text>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.hiddenping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenroleping",
		Description: "Hide a role ping behind text",
		Arguments:   "<roleid> <text>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.hiddenroleping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeneveryone",
		Description: "Hide an everyone ping behind text",
		Arguments:   "<text>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.hiddeneveryone,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenhere",
		Description: "Hide a here ping behind text",
		Arguments:   "<text>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.hiddenhere,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "voicemove",
		Description: "Move to random voice channels",
		Arguments:   "<amount>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.voicemove,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggertyping",
		Description: "Fake \"User is typing\" message",
		Arguments:   "<seconds> [channel]",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.triggertyping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "impersonate",
		Description: "Impersonate a user using a webhook",
		Arguments:   "<user> <text>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.impersonate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakepings",
		Description: "Get fake pinged in a channel",
		Arguments:   "<amount> [channel]",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.fakepings,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "billionpings",
		Description: "Get a billion pings in a channel",
		Arguments:   "[channel]",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     t.billionpings,
	})
}

func (TrollingCommands) purgehack(ctx *router.Ctx) {
	_ = ctx.RespondText("ﾠ" + strings.Repeat("\n", 1998) + "ﾠ")
}

func (TrollingCommands) ghostping(*router.Ctx) {
	// noop
}

func (TrollingCommands) gprole(*router.Ctx) {
	// noop
}

func (TrollingCommands) token(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	_ = ctx.RespondText(base64.StdEncoding.EncodeToString([]byte(user.ID)))
}

func (TrollingCommands) trollnitro(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("discord\u2024gift/%s %s https://i.imgur.com/0WSpeC1.png", random.String(16), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) fakenitro(ctx *router.Ctx) {
	amount, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}
	for i := 0; i < amount; i++ {
		_ = ctx.RespondText("https://discord.gift/" + random.String(16))
	}
}

func (TrollingCommands) hiddentext(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"%s %s %s",
		utils.Arguments["<\"text\">"],
		constants.INVISIBLE_TEXT,
		utils.Arguments["<\"hidden\">"],
	))
}

func (TrollingCommands) hiddeninvite(ctx *router.Ctx) {
	text := ctx.Arguments.RawFrom(1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.INVISIBLE_TEXT, utils.Arguments["<invite>"]))
}

func (TrollingCommands) hiddenping(ctx *router.Ctx) {
	text := ctx.Arguments.RawFrom(1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s <@%s>", text, constants.INVISIBLE_TEXT, utils.Arguments["<userid>"]))
}

func (TrollingCommands) hiddenroleping(ctx *router.Ctx) {
	text := ctx.Arguments.RawFrom(1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s <@&%s>", text, constants.INVISIBLE_TEXT, utils.Arguments["<roleid>"]))
}

func (TrollingCommands) hiddeneveryone(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("%s %s @everyone", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) hiddenhere(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("%s %s @here", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) voicemove(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Message.GuildID)
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	if conn, ok := ctx.Session.VoiceConnections[ctx.Message.GuildID]; ok {
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
		currentChannel := ctx.Session.VoiceConnections[ctx.Message.GuildID]

		if currentChannel != nil && currentChannel.ChannelID == channel.ID {
			continue
		}

		_, err := ctx.Session.ChannelVoiceJoin(ctx.Message.GuildID, channel.ID, false, true)
		if err != nil {
			_ = ctx.RespondText("Error joining voice channel: " + err.Error())
			return
		}

		time.Sleep(time.Second * 1)

		i++
	}

	if voice, ok := ctx.Session.VoiceConnections[ctx.Message.GuildID]; ok {
		_ = voice.Disconnect()
	}
}

func (TrollingCommands) triggertyping(ctx *router.Ctx) {
	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	channel, err := ctx.Arguments.Get(1).Channel()
	if err != nil {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	utils.Discord.SleepWithTypingFix(ctx.Session, channel.ID, time.Duration(seconds)*time.Second)
}

func (TrollingCommands) impersonate(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	text := ctx.Arguments.RawFrom(1)

	_, err := utils.Discord.Impersonate(ctx.Session, user, text, ctx.Message.ChannelID)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
}

func (TrollingCommands) fakepings(ctx *router.Ctx) {
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

func (TrollingCommands) billionpings(ctx *router.Ctx) {
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

type TrollingCommands struct{ services *services.Services }

func NewTrollingCommands(s *services.Services) *TrollingCommands {
	return &TrollingCommands{services: s}
}
func (c *TrollingCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

}
