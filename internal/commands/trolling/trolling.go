package trolling

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"

	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/random"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type TrollingCommands struct{ services *services.Services }

func NewTrollingCommands(s *services.Services) *TrollingCommands {
	return &TrollingCommands{services: s}
}

func (c *TrollingCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
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

func (c *TrollingCommands) purgehack(ctx *router.Ctx) error {
	return ctx.RespondText("%s", "ﾠ"+strings.Repeat("\n", 1998)+"ﾠ")
}

func (c *TrollingCommands) ghostping(ctx *router.Ctx) error {
	cfg := c.services.Config.GetConfig()

	if cfg.Selfbot.DeleteCommands {
		return nil
	}

	return ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
}

func (c *TrollingCommands) gprole(ctx *router.Ctx) error {
	cfg := c.services.Config.GetConfig()

	if cfg.Selfbot.DeleteCommands {
		return nil
	}

	return ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
}

func (c *TrollingCommands) token(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find user.")
	}

	return ctx.RespondText("%s", base64.StdEncoding.EncodeToString([]byte(user.ID)))
}

func (c *TrollingCommands) trollnitro(ctx *router.Ctx) error {
	return ctx.RespondText("discord\u2024gift/%s %s https://i.imgur.com/0WSpeC1.png", random.String(16), constants.INVISIBLE_TEXT)
}

func (c *TrollingCommands) fakenitro(ctx *router.Ctx) error {
	amount, _ := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 100)

	for range amount {
		ctx.RespondText("https://discord.gift/%s", random.String(16))
	}

	return nil
}

func (c *TrollingCommands) hiddentext(ctx *router.Ctx) error {
	return ctx.RespondText(
		"%s %s %s",
		ctx.Arguments.Get(0).Raw(),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(1).Raw(),
	)
}

func (c *TrollingCommands) hiddeninvite(ctx *router.Ctx) error {
	return ctx.RespondText(
		"%s %s %s",
		ctx.Arguments.RawFrom(1),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(0).Raw(),
	)
}

func (c *TrollingCommands) hiddenping(ctx *router.Ctx) error {
	return ctx.RespondText(
		"%s %s <@%s>",
		ctx.Arguments.RawFrom(1),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(0).Raw(),
	)
}

func (c *TrollingCommands) hiddenroleping(ctx *router.Ctx) error {
	return ctx.RespondText(
		"%s %s <@&%s>",
		ctx.Arguments.RawFrom(1),
		constants.INVISIBLE_TEXT,
		ctx.Arguments.Get(0).Raw(),
	)
}

func (c *TrollingCommands) hiddeneveryone(ctx *router.Ctx) error {
	return ctx.RespondText("%s %s @everyone", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT)
}

func (c *TrollingCommands) hiddenhere(ctx *router.Ctx) error {
	return ctx.RespondText("%s %s @here", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT)
}

func (c *TrollingCommands) voicemove(ctx *router.Ctx) error {
	guild, err := ctx.Session.GuildWithFallback(ctx.Message.GuildID)
	if err != nil {
		return fmt.Errorf("Error getting server: %v", err)
	}

	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 100)
	if err != nil {
		return fmt.Errorf("Argument <amount> must be a number.")
	}

	if conn, ok := ctx.Session.VoiceConnections[ctx.Message.GuildID]; ok {
		_ = conn.Disconnect()
	}

	channels := fn.Filter(guild.Channels, func(channel *discordgo.Channel) bool {
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
		return fmt.Errorf("No voice channels found")
	}

	for range amount {
		channel := random.Choice(channels)
		currentChannel := ctx.Session.VoiceConnections[ctx.Message.GuildID]

		if currentChannel != nil && currentChannel.ChannelID == channel.ID {
			continue
		}

		_, err := ctx.Session.ChannelVoiceJoin(ctx.Message.GuildID, channel.ID, false, true)
		if err != nil {
			return fmt.Errorf("Error joining voice channel: %v", err)
		}

		time.Sleep(time.Second * 1)
	}

	if voice, ok := ctx.Session.VoiceConnections[ctx.Message.GuildID]; ok {
		return voice.Disconnect()
	}
	return nil
}

func (c *TrollingCommands) triggertyping(ctx *router.Ctx) error {
	duration, err := ctx.Arguments.Get(0).Duration()
	duration = fn.Clamp(duration, 1*time.Second, 7*24*time.Hour)
	if err != nil {
		return fmt.Errorf("Failed to parse duration")
	}

	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	return ctx.Session.ChannelTypingDuration(channel.ID, duration)
}

func (c *TrollingCommands) impersonate(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find user.")
	}
	text := ctx.Arguments.RawFrom(1)

	_, err = utils.Discord.Impersonate(ctx.Session, user, text, ctx.Message.ChannelID)
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}
	return nil
}

func (c *TrollingCommands) fakepings(ctx *router.Ctx) error {
	amount, err := ctx.Arguments.Get(0).Int()
	amount = fn.Clamp(amount, 1, 1000000000)
	if err != nil {
		return fmt.Errorf("Argument <amount> must be a number.")
	}

	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 10)
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	if len(messages) == 0 {
		return fmt.Errorf("No messages found in channel")
	}

	err = ctx.Session.MessageUnack(channel.ID, messages[0].ID, amount)
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}
	return nil
}

func (c *TrollingCommands) billionpings(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find channel")
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 10)
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	if len(messages) == 0 {
		return fmt.Errorf("No messages found in channel")
	}

	err = ctx.Session.MessageUnack(channel.ID, messages[0].ID, 1000000000)
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}
	return nil
}
