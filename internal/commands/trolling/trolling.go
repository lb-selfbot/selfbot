package trolling

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/internal/utils2"
)

type TrollingCommands struct{ services *services.Services }

func NewTrollingCommands(s *services.Services) *TrollingCommands {
	return &TrollingCommands{services: s}
}

func (c *TrollingCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Trolling",
		Description: "Trolling commands",
	}, []*router.Command{
		{
			Name:        "purgehack",
			Description: "Purge a chat without permissions",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.purgehack,
		},
		{
			Name:        "ghostping",
			Description: "Ghostping a user",
			Arguments:   "<user>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.ghostping,
		},
		{
			Name:        "gprole",
			Description: "Ghostping a role",
			Arguments:   "<role>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.gprole,
		},
		{
			Name:        "token",
			Description: "Send the first part of someones token",
			Arguments:   "<user>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.token,
		},
		{
			Name:        "trollnitro",
			Description: "Send a fake nitro link with an image",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.trollnitro,
		},
		{
			Name:        "fakenitro",
			Description: "Generate fake nitro codes",
			Arguments:   "[amount]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.fakenitro,
		},
		{
			Name:        "hiddentext",
			Description: "Hide text behind other text",
			Arguments:   "<\"text\"> <\"hidden\">",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.hiddentext,
		},
		{
			Name:        "hiddeninvite",
			Description: "Hide an invite behind text",
			Arguments:   "<invite> <text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.hiddeninvite,
		},
		{
			Name:        "hiddenping",
			Description: "Hide a ping behind text",
			Arguments:   "<user id> <text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.hiddenping,
		},
		{
			Name:        "hiddenroleping",
			Description: "Hide a role ping behind text",
			Arguments:   "<role id> <text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.hiddenroleping,
		},
		{
			Name:        "hiddeneveryone",
			Description: "Hide an everyone ping behind text",
			Arguments:   "<text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.hiddeneveryone,
		},
		{
			Name:        "hiddenhere",
			Description: "Hide a here ping behind text",
			Arguments:   "<text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.hiddenhere,
		},
		{
			Name:        "voicemove",
			Description: "Move to random voice channels",
			Arguments:   "<amount>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.voicemove,
		},
		{
			Name:        "triggertyping",
			Description: "Fake \"User is typing\" message",
			Arguments:   "<duration> [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.triggertyping,
		},
		{
			Name:        "impersonate",
			Description: "Impersonate a user using a webhook",
			Arguments:   "<user> <text>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.impersonate,
		},
		{
			Name:        "fakepings",
			Description: "Get fake pinged in a channel",
			Arguments:   "<amount> [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.fakepings,
		},
		{
			Name:        "billionpings",
			Description: "Get a billion pings in a channel",
			Arguments:   "[channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.billionpings,
		},
	})
}

func (c *TrollingCommands) purgehack(ctx *router.Ctx) {
	ctx.RespondText("ﾠ" + strings.Repeat("\n", 1998) + "ﾠ")
}

func (c *TrollingCommands) ghostping(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if cfg.Selfbot.DeleteCommands {
		return
	}

	_ = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
}

func (c *TrollingCommands) gprole(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	if cfg.Selfbot.DeleteCommands {
		return
	}

	_ = ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
}

func (c *TrollingCommands) token(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find user.")
		return
	}

	ctx.RespondText(base64.StdEncoding.EncodeToString([]byte(user.ID)))
}

func (c *TrollingCommands) trollnitro(ctx *router.Ctx) {
	ctx.RespondText("discord\u2024gift/%s %s https://i.imgur.com/0WSpeC1.png", random.String(16), constants.INVISIBLE_TEXT)
}

func (c *TrollingCommands) fakenitro(ctx *router.Ctx) {
	amount, _ := ctx.Arguments.Get(0).Int()
	amount = utils2.Clamp(amount, 1, 100)

	for range amount {
		ctx.RespondText("https://discord.gift/%s", random.String(16))
	}
}

func (c *TrollingCommands) hiddentext(ctx *router.Ctx) {
	ctx.RespondText(
		"%s %s %s",
		ctx.Arguments.Get(0).Raw(),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(1).Raw(),
	)
}

func (c *TrollingCommands) hiddeninvite(ctx *router.Ctx) {
	ctx.RespondText(
		"%s %s %s",
		ctx.Arguments.RawFrom(1),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(0).Raw(),
	)
}

func (c *TrollingCommands) hiddenping(ctx *router.Ctx) {
	ctx.RespondText(
		"%s %s <@%s>",
		ctx.Arguments.RawFrom(1),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(0).Raw(),
	)
}

func (c *TrollingCommands) hiddenroleping(ctx *router.Ctx) {
	ctx.RespondText(
		"%s %s <@&%s>",
		ctx.Arguments.RawFrom(1),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(0).Raw(),
	)
}

func (c *TrollingCommands) hiddeneveryone(ctx *router.Ctx) {
	ctx.RespondText("%s %s @everyone", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT)
}

func (c *TrollingCommands) hiddenhere(ctx *router.Ctx) {
	ctx.RespondText("%s %s @here", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT)
}

func (c *TrollingCommands) voicemove(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Message.GuildID)
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return
	}

	amount, err := ctx.Arguments.Get(0).Int()
	amount = utils2.Clamp(amount, 1, 100)
	if err != nil {
		ctx.RespondText("Argument <amount> must be a number.")
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
		ctx.RespondText("No voice channels found")
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
			ctx.RespondText("Error joining voice channel: %v", err)
			return
		}

		time.Sleep(time.Second * 1)

		i++
	}

	if voice, ok := ctx.Session.VoiceConnections[ctx.Message.GuildID]; ok {
		_ = voice.Disconnect()
	}
}

func (c *TrollingCommands) triggertyping(ctx *router.Ctx) {
	duration, err := ctx.Arguments.Get(0).Duration()
	duration = utils2.Clamp(duration, 1*time.Second, 7*24*time.Hour)
	if err != nil {
		ctx.RespondText("Failed to parse duration")
		return
	}

	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}
	utils.Discord.SleepWithTypingFix(ctx.Session, channel.ID, duration)
}

func (c *TrollingCommands) impersonate(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find user.")
		return
	}
	text := ctx.Arguments.RawFrom(1)

	_, err = utils.Discord.Impersonate(ctx.Session, user, text, ctx.Message.ChannelID)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}
}

func (c *TrollingCommands) fakepings(ctx *router.Ctx) {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = utils2.Clamp(amount, 1, 1000000000)
	if err != nil {
		ctx.RespondText("Argument <amount> must be a number.")
		return
	}

	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 10)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	if len(messages) == 0 {
		ctx.RespondText("No messages found in channel")
		return
	}

	err = ctx.Session.MessageUnack(channel.ID, messages[0].ID, amount)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}
}

func (c *TrollingCommands) billionpings(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Could not find channel")
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 10)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	if len(messages) == 0 {
		ctx.RespondText("No messages found in channel")
		return
	}

	err = ctx.Session.MessageUnack(channel.ID, messages[0].ID, 1000000000)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}
}
