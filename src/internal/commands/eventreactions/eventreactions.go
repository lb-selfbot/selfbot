package eventreactions

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventreactions/impl"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"

	"github.com/LightningDev1/discordgo"
)

type EventReactionsCommands struct{ services *services.Services }

func NewEventReactionsCommands(s *services.Services) *EventReactionsCommands {
	return &EventReactionsCommands{services: s}
}
func (c *EventReactionsCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "EventReactions",
		Description: "Event Reaction commands",
	}, []*router.Command{
		{
			Name:        "listreactions",
			Description: "List all active event reactions",
			Arguments:   "listreactions",
			Handler:     c.listreactions,
		},
		{
			Name:        "clearreactions",
			Description: "Clear all active event reactions",
			Arguments:   "clearreactions",
			Handler:     c.clearreactions,
		},
		{
			Name:        "ignoreuser",
			Description: "Ignore messages from a user",
			Arguments:   "<toggle/list/clear> [@user]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.ignoreuser,
		},
		{
			Name:        "ignoreserver",
			Description: "Ignore messages from a server",
			Arguments:   "<toggle/list/clear> [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.ignoreserver,
		},
		{
			Name:        "pingping",
			Description: "Reply to pings by pinging the user",
			Arguments:   "<toggle/list/clear> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingping,
		},
		{
			Name:        "pingkick",
			Description: "Reply to pings by kicking the user",
			Arguments:   "<toggle/list/clear> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingkick,
		},
		{
			Name:        "pingban",
			Description: "Reply to pings by banning the user",
			Arguments:   "<toggle/list/clear> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingban,
		},
		{
			Name:        "pingmute",
			Description: "Reply to pings by muting the user",
			Arguments:   "<toggle/list/clear> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingmute,
		},
		{
			Name:        "pingtimeout",
			Description: "Reply to pings by timing the user out",
			Arguments:   "<toggle/list/clear> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingtimeout,
		},
		{
			Name:        "pinginsult",
			Description: "Reply to pings by insulting the user",
			Arguments:   "<toggle/list/clear> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pinginsult,
		},
		{
			Name:        "pingreact",
			Description: "Reply to pings by adding a reaction",
			Arguments:   "<toggle/list/clear> [reaction] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingreact,
		},
		{
			Name:        "pingresponse",
			Description: "Reply to pings with a custom response",
			Arguments:   "<toggle/list/clear> [\"response\"] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingresponse,
		},
		{
			Name:        "pingairesponse",
			Description: "Reply to pings with an AI response",
			Arguments:   "<toggle/list/clear> [\"system\"] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pingairesponse,
		},
		{
			Name:        "mirrorchannel",
			Description: "Mirror messages from one channel to another",
			Arguments:   "<toggle/list/clear> [source channel] [dest channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.mirrorchannel,
		},
		{
			Name:        "mirrorserver",
			Description: "Mirror messages from one server to another",
			Arguments:   "<toggle/list/clear> [source server] [dest server]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.mirrorserver,
		},
		{
			Name:        "responder",
			Description: "Reply to messages with a custom response",
			Arguments:   "<toggle/list/clear> [\"response\"] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.responder,
		},
		{
			Name:        "airesponder",
			Description: "Reply to messages with an AI response",
			Arguments:   "<toggle/list/clear> [\"system\"] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.airesponder,
		},
		{
			Name:        "copy",
			Description: "Copy a user's messages",
			Arguments:   "<toggle/list/clear> [@user] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.copy,
		},
		{
			Name:        "replycopy",
			Description: "Replycopy a user's messages",
			Arguments:   "<toggle/list/clear> [@user] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.replycopy,
		},
		{
			Name:        "mock",
			Description: "Mock a user's messages",
			Arguments:   "<toggle/list/clear> <@user> [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.mock,
		},
		{
			Name:        "pinspam",
			Description: "Pin all of a user's messages",
			Arguments:   "<toggle/list/clear> [@user] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.pinspam,
		},
		{
			Name:        "unpinspam",
			Description: "Unpin messages pinned by pinspam",
			Arguments:   "unpinspam",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.unpinspam,
		},
		{
			Name:        "reactspam",
			Description: "React to all of a user's messages",
			Arguments:   "<toggle/list/clear> [@user] [reaction] [global/server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.reactspam,
		},
		{
			Name:        "deletemessages",
			Description: "Delete all of a user's messages",
			Arguments:   "<toggle/list/clear> [@user] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.deletemessages,
		},
		{
			Name:        "keywordnotifier",
			Description: "Get notified when someone mentions a keyword",
			Arguments:   "<toggle/list/clear> [\"keyword\"] [global/server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.keywordnotifier,
		},
		{
			Name:        "removereactions",
			Description: "Delete all of a user's reactions",
			Arguments:   "<toggle/list/clear> [@user] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.removereactions,
		},
		{
			Name:        "forcedisconnect",
			Description: "Disconnect a user from voice channels",
			Arguments:   "<toggle/list/clear> [@user] [server/channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.forcedisconnect,
		},
		{
			Name:        "forcenick",
			Description: "Force change a user's nickname",
			Arguments:   "<toggle/list/clear> [@user] [nickname] [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.forcenick,
		},
		{
			Name:        "forcekick",
			Description: "Force kick a user",
			Arguments:   "<toggle/list/clear> [@user] [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.forcekick,
		},
		{
			Name:        "noleave",
			Description: "Prevent a user from leaving a group chat",
			Arguments:   "<toggle/list/clear> [@user] [channel id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.noleave,
		},
		{
			Name:        "antigroupchat",
			Description: "Prevent users from adding you to group chats",
			Arguments:   "<on/off> [silent leave: on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antigroupchat,
		},
		{
			Name:        "antideletemessages",
			Description: "Prevent users from deleting your messages",
			Arguments:   "<toggle/list/clear> [channel]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antideletemessages,
		},
		{
			Name:        "afkresponder",
			Description: "Get, set or toggle the AFK Responder",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkresponder,
		},
		{
			Name:        "afkonping",
			Description: "Get, set or toggle if AFK message should be sent on ping",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkonping,
		},
		{
			Name:        "afkondm",
			Description: "Get, set or toggle if AFK message should be sent on DM",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkondm,
		},
		{
			Name:        "afkdelay",
			Description: "Get or set the delay before sending AFK message",
			Arguments:   "<get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkdelay,
		},
		{
			Name:        "afkrandomdelay",
			Description: "Get, set or toggle if AFK delay is random",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkrandomdelay,
		},
		{
			Name:        "afkinactive",
			Description: "Get, set or toggle if AFK messages are only sent when inactive",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkinactive,
		},
		{
			Name:        "afkinactivetimeout",
			Description: "Get or set the timeout before being considered inactive",
			Arguments:   "<get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkinactivetimeout,
		},
		{
			Name:        "afktyping",
			Description: "Get, set or toggle if the AFK Responder should show typing",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afktyping,
		},
		{
			Name:        "afkai",
			Description: "Get, set or toggle if the AFK Responder should use AI",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkai,
		},
		{
			Name:        "afkmessage",
			Description: "Get or set the AFK message (or instructions for AI if enabled)",
			Arguments:   "<get/set> [value]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afkmessage,
		},
		{
			Name:        "afktriggertype",
			Description: "Get or set the AFK Triggerlist type",
			Arguments:   "<get/set/values> [type]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afktriggertype,
		},
		{
			Name:        "afktriggerusers",
			Description: "Get, add, remove or clear users that trigger AFK message",
			Arguments:   "<get/add/rm/clear> [user]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afktriggerusers,
		},
		{
			Name:        "afktriggerservers",
			Description: "Get, add, remove or clear servers that trigger AFK message",
			Arguments:   "<get/add/rm/clear> [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.afktriggerservers,
		},
	})
}

// Utils

func (c *EventReactionsCommands) toggleReaction(ctx *router.Ctx, reaction interfaces.EventReaction) {
	// Delete the raction if it already exists
	if slices.Contains(c.services.EventReactions.GetAllReactions(), reaction) {
		c.services.EventReactions.RemoveReactionsByType(reaction.Type())
		c.services.Embed.NewEmbed().
			SetDescription("Removed event reaction: " + reaction.Description()).
			Send(ctx)
		return
	}

	c.services.EventReactions.AddReaction(reaction)

	c.services.Embed.NewEmbed().
		SetDescription("Added event reaction: " + reaction.Description()).
		Send(ctx)
}

func (c *EventReactionsCommands) mapReactions(ctx *router.Ctx, reactions []interfaces.EventReaction, errText, title string) {
	if len(reactions) == 0 {
		c.services.Embed.NewEmbed().
			SetDescription(errText).
			Send(ctx)
		return
	}

	descriptions := utils.Map(reactions, func(reaction interfaces.EventReaction) string {
		return reaction.Description()
	})

	c.services.Embed.NewEmbed().
		SetDescription(title + ":\n\n" + strings.Join(descriptions, "\n")).
		Send(ctx)
}

func (c *EventReactionsCommands) getGuildOrChannel(ctx *router.Ctx, arg *router.Argument) (*discordgo.Guild, *discordgo.Channel, error) {
	guild, err := arg.Guild()
	if err == nil {
		return guild, nil, nil
	}

	channel, err := arg.Channel()
	if err == nil {
		return nil, channel, nil
	}

	if ctx.Message.GuildID != "" {
		guild, err = utils.Discord.GetGuild(ctx.Session, ctx.Message.GuildID)
		return guild, nil, err
	}

	channel, err = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
	return nil, channel, err
}

// Commands

func (c *EventReactionsCommands) listreactions(ctx *router.Ctx) {
	reactions := c.services.EventReactions.GetAllReactions()
	c.mapReactions(ctx, reactions, "There are no active event reactions", "Active event reactions")
}

func (c *EventReactionsCommands) clearreactions(ctx *router.Ctx) {
	c.services.EventReactions.RemoveAllReactions()
	c.services.Embed.NewEmbed().
		SetDescription("Cleared all event reactions").
		Send(ctx)
}

// Ignore User

func (c *EventReactionsCommands) ignoreuser(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewIgnoreUser(user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionIgnoreUser)

		c.mapReactions(ctx, reactions, "There are no ignored users.", "Ignored users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionIgnoreUser)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all ignored users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

// Ignore Server

func (c *EventReactionsCommands) ignoreserver(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guildID := utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID)

		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewIgnoreGuild(guild))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionIgnoreGuild)

		c.mapReactions(ctx, reactions, "There are no ignored servers.", "Ignored servers")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionIgnoreGuild)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all ignored servers.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

// Ping reactions

func (c *EventReactionsCommands) pingReaction(ctx *router.Ctx, reactionType interfaces.EventReactionType, reaction string, i int) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(i))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		eventReaction := impl.NewPingReaction(c.services.EventReactions, reactionType, guild, channel)

		eventReaction.System = utils.Arguments["[\"system\"]"]
		eventReaction.Response = utils.Arguments["[\"response\"]"]

		eventReaction.Reaction = utils.Arguments["[reaction]"]
		eventReaction.Reaction = strings.ReplaceAll(eventReaction.Reaction, "<", "")
		eventReaction.Reaction = strings.ReplaceAll(eventReaction.Reaction, ">", "")

		if eventReaction.Reaction == "" && reactionType == interfaces.EventReactionPingReact {
			c.services.Embed.NewEmbed().
				SetDescription("Reaction cannot be empty!").
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, eventReaction)
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(reactionType)

		c.mapReactions(ctx, reactions, fmt.Sprintf("There are no ping --> %s reactions.", reaction), fmt.Sprintf("Ping --> %s reactions", reaction))
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(reactionType)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all ping --> %s reactions.", reaction).
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) pingping(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingPing, "ping", 1)
}

func (c *EventReactionsCommands) pingkick(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingKick, "kick", 1)
}

func (c *EventReactionsCommands) pingban(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingBan, "ban", 1)
}

func (c *EventReactionsCommands) pingmute(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingMute, "mute", 1)
}

func (c *EventReactionsCommands) pingtimeout(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingTimeout, "timeout", 1)
}

func (c *EventReactionsCommands) pinginsult(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingInsult, "insult", 1)
}

func (c *EventReactionsCommands) pingreact(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingReact, "react", 2)
}

func (c *EventReactionsCommands) pingresponse(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingResponse, "response", 2)
}

func (c *EventReactionsCommands) pingairesponse(ctx *router.Ctx) {
	c.pingReaction(ctx, interfaces.EventReactionPingAIResponse, "AI response", 2)
}

// Mirror channel/server

func (c *EventReactionsCommands) mirrorchannel(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		srcChannel, err := ctx.Arguments.Get(1).Channel()
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find source channel").
				Send(ctx)
			return
		}

		destChannel, err := ctx.Arguments.Get(2).Channel()
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find destination channel").
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMirrorChannel(c.services.EventReactions, srcChannel, destChannel))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionMirrorChannel)

		c.mapReactions(ctx, reactions, "There are no mirrored channels.", "Mirrored channels")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionMirrorChannel)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all mirrored channels.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) mirrorserver(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		srcGuild, srcErr := utils.Discord.GetGuild(ctx.Session, utils.Arguments["[source server]"])
		destGuild, destErr := utils.Discord.GetGuild(ctx.Session, utils.Arguments["[dest server]"])

		if srcErr != nil || destErr != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server").
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMirrorGuild(c.services.EventReactions, srcGuild, destGuild))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionMirrorGuild)

		c.mapReactions(ctx, reactions, "There are no mirrored servers.", "Mirrored servers")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionMirrorGuild)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all mirrored servers.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

// Message Reactions

func (c *EventReactionsCommands) responder(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		reaction := impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionResponder, guild, channel, nil)
		reaction.Response = utils.Arguments["[\"response\"]"]
		c.toggleReaction(ctx, reaction)
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionResponder)

		c.mapReactions(ctx, reactions, "There are no active responders.", "Responders")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionResponder)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all responders.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) airesponder(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		reaction := impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionAIResponder, guild, channel, nil)
		reaction.System = utils.Arguments["[\"system\"]"]
		c.toggleReaction(ctx, reaction)
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionAIResponder)

		c.mapReactions(ctx, reactions, "There are no active AI responders.", "AI Responders")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionAIResponder)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all AI responders.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) copy(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionCopy, guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionCopy)

		c.mapReactions(ctx, reactions, "There are no copied users.", "Copied users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionCopy)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all copied users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) replycopy(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionReplyCopy, guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionReplyCopy)

		c.mapReactions(ctx, reactions, "There are no replycopied users.", "Replycopied users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionReplyCopy)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all replycopied users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) mock(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionMock, guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionMock)

		c.mapReactions(ctx, reactions, "There are no mocked users.", "Mocked users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionMock)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all mocked users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) pinspam(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionPinSpam, guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionPinSpam)

		c.mapReactions(ctx, reactions, "There are no pinspammed users.", "Pinspammed users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionPinSpam)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all pinspammed users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) unpinspam(ctx *router.Ctx) {
	for _, message := range c.services.EventReactions.GetPinSpamMessages() {
		err := ctx.Session.ChannelMessageUnpin(message.ChannelID, message.ID)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not unpin message: " + err.Error()).
				Send(ctx)
			return
		}
	}

	c.services.EventReactions.RemovePinSpamMessages()

	c.services.Embed.NewEmbed().
		SetDescription("Unpinned all messages.").
		Send(ctx)
}

func (c *EventReactionsCommands) reactspam(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		arg := strings.ToLower(utils.Arguments["[global/server/channel id]"])
		global := arg == "global" || arg == ""

		var guild *discordgo.Guild
		var channel *discordgo.Channel

		if !global {
			guild, channel, err = c.getGuildOrChannel(ctx, ctx.Arguments.Get(3))
			if err != nil {
				c.services.Embed.NewEmbed().
					SetDescription("Could not find server/channel: " + err.Error()).
					Send(ctx)
				return
			}
		}

		reaction := impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionReactSpam, guild, channel, user)
		reaction.Global = global
		reaction.Reaction = utils.Arguments["[reaction]"]
		reaction.Reaction = strings.ReplaceAll(reaction.Reaction, "<", "")
		reaction.Reaction = strings.ReplaceAll(reaction.Reaction, ">", "")

		if reaction.Reaction == "" {
			ctx.RespondText("Reaction cannot be empty!")
			return
		}

		c.toggleReaction(ctx, reaction)
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionReactSpam)

		c.mapReactions(ctx, reactions, "There are no reactspammed users.", "Reactspammed users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionReactSpam)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all reactspammed users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) deletemessages(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionDeleteMessages, guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionDeleteMessages)

		c.mapReactions(ctx, reactions, "There are no deletemessage reactions.", "Deletemessage reactions")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionDeleteMessages)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all deletemessage reactions.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) keywordnotifier(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		arg := strings.ToLower(utils.Arguments["[global/server/channel id]"])
		global := arg == "global" || arg == ""

		var guild *discordgo.Guild
		var channel *discordgo.Channel
		var err error

		if !global {
			guild, channel, err = c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
			if err != nil {
				c.services.Embed.NewEmbed().
					SetDescription("Could not find server/channel: " + err.Error()).
					Send(ctx)
				return
			}
		}

		reaction := impl.NewMessageReaction(c.services.UI, c.services.EventReactions, interfaces.EventReactionKeywordNotifier, guild, channel, nil)
		reaction.Keyword = utils.Arguments["[\"keyword\"]"]
		reaction.Global = global

		c.toggleReaction(ctx, reaction)
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionKeywordNotifier)

		c.mapReactions(ctx, reactions, "There are no keyword notifiers.", "Keyword notifiers")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionKeywordNotifier)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all keyword notifiers.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) removereactions(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewRemoveReactions(guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionRemoveReactions)

		c.mapReactions(ctx, reactions, "There are no users with removed reactions.", "Users with removed reactions")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionRemoveReactions)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all users with removed reactions.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) forcedisconnect(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewForceDisconnect(guild, channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionForceDisconnect)

		c.mapReactions(ctx, reactions, "There are no force disconnected users.", "Force disconnected users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionForceDisconnect)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all force disconnected users").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) forcenick(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guildID := utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID)

		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server: " + err.Error()).
				Send(ctx)
			return
		}

		nickname := utils.Arguments["[nickname]"]

		err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, user.ID, nickname)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not set nickname: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewForceNick(guild, user, nickname))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionForceNick)

		c.mapReactions(ctx, reactions, "There are no force nicked users.", "Force nicked users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionForceNick)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all force nicked users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) forcekick(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		guildID := utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID)

		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server: " + err.Error()).
				Send(ctx)
			return
		}

		err = ctx.Session.GuildMemberDelete(ctx.Message.GuildID, user.ID)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not kick user: " + err.Error()).
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewForceKick(guild, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionForceKick)

		c.mapReactions(ctx, reactions, "There are no force kicked users.", "Force kicked users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionForceKick)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all force kicked users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) noleave(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return

		}

		relationships := utils.Filter(ctx.Session.State.Relationships, func(r *discordgo.Relationship) bool {
			return r.Type == 1 && r.ID == user.ID
		})

		if len(relationships) == 0 {
			c.services.Embed.NewEmbed().
				SetDescription("The user must be a friend.").
				Send(ctx)
			return
		}

		channel, err := ctx.Arguments.Get(2).Channel(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find channel").
				Send(ctx)
			return
		}

		if channel.Type != discordgo.ChannelTypeGroupDM {
			c.services.Embed.NewEmbed().
				SetDescription("The channel must be a group DM.").
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewNoLeave(channel, user))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionNoLeave)

		c.mapReactions(ctx, reactions, "There are no noleave users.", "Noleave users")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionNoLeave)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all noleave users.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) antigroupchat(ctx *router.Ctx) {
	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	isSilent, _ := ctx.Arguments.Get(1).Bool()

	reaction := impl.NewAntiGroupChat(isSilent)

	// Remoe to reset
	c.services.EventReactions.RemoveReactionsByType(reaction.Type())

	if toggle {
		c.services.EventReactions.AddReaction(reaction)

		c.services.Embed.NewEmbed().
			SetDescription("Turned on the anti group chat reaction!").
			Send(ctx)
	} else {
		c.services.Embed.NewEmbed().
			SetDescription("Turned off the anti group chat reaction!").
			Send(ctx)
	}
}

func (c *EventReactionsCommands) antideletemessages(ctx *router.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find channel").
				Send(ctx)
			return
		}

		c.toggleReaction(ctx, impl.NewAntiDeleteMessages(channel))
	case "list":
		reactions := c.services.EventReactions.GetReactionsByType(interfaces.EventReactionAntiDeleteMessages)

		c.mapReactions(ctx, reactions, "There are no anti delete messages.", "Anti Delete Messages")
	case "clear":
		c.services.EventReactions.RemoveReactionsByType(interfaces.EventReactionAntiDeleteMessages)

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all anti delete messages.").
			Send(ctx)
	default:
		c.services.Embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

// AFK

func (c *EventReactionsCommands) afkresponder(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK Responder", utils.Misc.FormatBoolOnOff(cfg.AFK.Enabled)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.Enabled = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK Responder").
			Send(ctx)
	case "toggle":
		cfg.AFK.Enabled = !cfg.AFK.Enabled
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.Enabled) + " AFK Responder").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkonping(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK On Ping", utils.Misc.FormatBoolOnOff(cfg.AFK.OnPing)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.OnPing = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK On Ping").
			Send(ctx)
	case "toggle":
		cfg.AFK.OnPing = !cfg.AFK.OnPing
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.OnPing) + " AFK On Ping").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkondm(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK On DM", utils.Misc.FormatBoolOnOff(cfg.AFK.OnDM)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.OnDM = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK On DM").
			Send(ctx)
	case "toggle":
		cfg.AFK.OnDM = !cfg.AFK.OnDM
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.OnDM) + " AFK On DM").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkdelay(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		formatted := utils.Misc.FormatDuration(time.Duration(cfg.AFK.Delay) * time.Second)

		c.services.Embed.NewEmbed().
			AddField("AFK Delay", formatted).
			Send(ctx)
	case "set":
		delay, err := utils.Misc.ParseDuration(utils.Arguments["[value]"])
		if err != nil {
			ctx.RespondText("Invalid duration: " + err.Error())
			return
		}

		cfg.AFK.Delay = int(delay.Seconds())
		_ = cfg.Save()

		formatted := utils.Misc.FormatDuration(time.Duration(cfg.AFK.Delay) * time.Second)

		c.services.Embed.NewEmbed().
			SetDescription("AFK Delay set to %s.", formatted).
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (c *EventReactionsCommands) afkrandomdelay(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK Random Delay", utils.Misc.FormatBoolOnOff(cfg.AFK.RandomDelay)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.RandomDelay = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK Random Delay").
			Send(ctx)
	case "toggle":
		cfg.AFK.RandomDelay = !cfg.AFK.RandomDelay
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.RandomDelay) + " AFK Random Delay").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkinactive(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK Must Be Inactive", utils.Misc.FormatBoolOnOff(cfg.AFK.MustBeInactive)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.MustBeInactive = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK Must Be Inactive").
			Send(ctx)
	case "toggle":
		cfg.AFK.MustBeInactive = !cfg.AFK.MustBeInactive
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.MustBeInactive) + " AFK Must Be Inactive").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkinactivetimeout(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		formatted := utils.Misc.FormatDuration(time.Duration(cfg.AFK.InactiveTimeout) * time.Second)

		c.services.Embed.NewEmbed().
			AddField("AFK Inactive Timeout", formatted).
			Send(ctx)
	case "set":
		delay, err := utils.Misc.ParseDuration(utils.Arguments["[value]"])
		if err != nil {
			ctx.RespondText("Invalid duration: " + err.Error())
			return
		}

		cfg.AFK.InactiveTimeout = int(delay.Seconds())
		_ = cfg.Save()

		formatted := utils.Misc.FormatDuration(time.Duration(cfg.AFK.InactiveTimeout) * time.Second)

		c.services.Embed.NewEmbed().
			SetDescription("AFK Inactive Timeout set to %s.", formatted).
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (c *EventReactionsCommands) afktyping(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK Typing", utils.Misc.FormatBoolOnOff(cfg.AFK.Typing)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.Typing = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK Typing").
			Send(ctx)
	case "toggle":
		cfg.AFK.Typing = !cfg.AFK.Typing
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.Typing) + " AFK Typing").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkai(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK AI Response", utils.Misc.FormatBoolOnOff(cfg.AFK.AI)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.AFK.AI = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " AFK AI Response").
			Send(ctx)
	case "toggle":
		cfg.AFK.AI = !cfg.AFK.AI
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.AFK.AI) + " AFK AI Response").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (c *EventReactionsCommands) afkmessage(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK Message", cfg.AFK.Message).
			Send(ctx)
	case "set":
		value := ctx.Arguments.RawFrom(1)

		cfg.AFK.Message = value
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("AFK Message set to " + value).
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set> argument. Use `get` or `set`.")
	}
}

func (c *EventReactionsCommands) afktriggertype(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/values>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("AFK Trigger Type", cfg.AFK.TriggerList.Type).
			Send(ctx)
	case "set":
		value := strings.ToLower(utils.Arguments["[type]"])

		switch value {
		case "global":
			cfg.AFK.TriggerList.Type = "Global"
		case "whitelist":
			cfg.AFK.TriggerList.Type = "Whitelist"
		case "blacklist":
			cfg.AFK.TriggerList.Type = "Blacklist"
		default:
			ctx.RespondText("Invalid trigger type. Use `[p]afktriggertype values` to get a list of available types.")
			return
		}

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("AFK Trigger Type set to " + cfg.Selfbot.Client).
			Send(ctx)
	case "values":
		c.services.Embed.NewEmbed().
			SetDescription("Available AFK Trigger Type values: `Global`, `Whitelist`, `Blacklist`.").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/values> argument. Use `get` or `set` or `values`.")
	}
}

func (c *EventReactionsCommands) afktriggerusers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch strings.ToLower(utils.Arguments["<get/add/rm/clear>"]) {
	case "get":
		userString := ""
		for _, userID := range cfg.AFK.TriggerList.Users {
			user, err := utils.Discord.GetUser(ctx.Session, userID)
			if err != nil {
				userString += fmt.Sprintf("%s - %s\n", "Unknown", userID)
				continue
			}

			userString += fmt.Sprintf("%s - %s\n", user.Username, userID)
		}

		if userString == "" {
			userString = "No users set."
		}

		c.services.Embed.NewEmbed().
			AddField("AFK Triggerlist Users", userString).
			Send(ctx)
	case "add":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user.").
				Send(ctx)
			return

		}

		if slices.Contains(cfg.AFK.TriggerList.Users, user.ID) {
			ctx.RespondText("User is already in the AFK Triggerlist.")
			return
		}

		cfg.AFK.TriggerList.Users = append(cfg.AFK.TriggerList.Users, user.ID)

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Added %s to the AFK Triggerlist.", user.Mention()).
			Send(ctx)
	case "rm":
		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find user.").
				Send(ctx)
			return

		}

		if !slices.Contains(cfg.AFK.TriggerList.Users, user.ID) {
			ctx.RespondText("User already isn't in the AFK Triggerlist.")
			return
		}

		cfg.AFK.TriggerList.Users = utils.RemoveSliceItem(cfg.AFK.TriggerList.Users, user.ID)

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Removed %s from the AFK Triggerlist.", user.Mention()).
			Send(ctx)
	case "clear":
		cfg.AFK.TriggerList.Users = []string{}
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all users from the AFK Triggerlist.").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (c *EventReactionsCommands) afktriggerservers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch strings.ToLower(utils.Arguments["<get/add/rm/clear>"]) {
	case "get":
		var servers []string

		for _, serverID := range cfg.AFK.TriggerList.Servers {
			guild, err := utils.Discord.GetGuild(ctx.Session, serverID)
			if err != nil {
				servers = append(servers, serverID)
				continue
			}

			servers = append(servers, guild.Name+" ("+guild.ID+")")
		}

		c.services.Embed.NewEmbed().
			AddField("AFK Triggerlist Servers", strings.Join(servers, ", ")).
			Send(ctx)
	case "add":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server.").
				Send(ctx)
			return
		}

		if slices.Contains(cfg.AFK.TriggerList.Servers, guild.ID) {
			ctx.RespondText("Server %s is already in the AFK Triggerlist!", guild.Name)
			return
		}

		cfg.AFK.TriggerList.Servers = append(cfg.AFK.TriggerList.Servers, guild.ID)

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Added %s to the AFK Triggerlist.", guild.Name).
			Send(ctx)
	case "rm":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
		if err != nil {
			c.services.Embed.NewEmbed().
				SetDescription("Could not find server.").
				Send(ctx)
			return
		}

		if !slices.Contains(cfg.AFK.TriggerList.Servers, guild.ID) {
			ctx.RespondText("Server %s already isn't in the AFK Triggerlist!", guild.Name)
			return
		}

		cfg.AFK.TriggerList.Servers = utils.RemoveSliceItem(cfg.AFK.TriggerList.Servers, guild.ID)

		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Removed %s from the AFK Triggerlist.", guild.Name).
			Send(ctx)
	case "clear":
		cfg.AFK.TriggerList.Servers = []string{}
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all servers from the AFK Triggerlist.").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}
