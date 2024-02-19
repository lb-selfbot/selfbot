package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventhandlers/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type EventReactionsCommands struct{}

func (e *EventReactionsCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("EventReactions", "Event Reaction commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "listreactions",
		Description: "List all active event reactions",
		Usage:       "[p]listreactions",
		Handler:     e.listreactions,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clearreactions",
		Description: "Clear all active event reactions",
		Usage:       "[p]clearreactions",
		Handler:     e.clearreactions,
	})

	// Ignore user
	router.RegisterCmd(&dgc.Command{
		Name:        "ignoreuser",
		Description: "Ignore messages from a user",
		Usage:       "[p]ignoreuser <toggle/list/clear> [@user]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.ignoreuser,
	})

	// Ignore server
	router.RegisterCmd(&dgc.Command{
		Name:        "ignoreserver",
		Description: "Ignore messages from a server",
		Usage:       "[p]ignoreserver <toggle/list/clear> [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.ignoreserver,
	})

	// Ping reactions
	router.RegisterCmd(&dgc.Command{
		Name:        "pingping",
		Description: "Reply to pings by pinging the user",
		Usage:       "[p]pingping <toggle/list/clear> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingkick",
		Description: "Reply to pings by kicking the user",
		Usage:       "[p]pingkick <toggle/list/clear> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingkick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingban",
		Description: "Reply to pings by banning the user",
		Usage:       "[p]pingban <toggle/list/clear> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingban,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingmute",
		Description: "Reply to pings by muting the user",
		Usage:       "[p]pingmute <toggle/list/clear> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingmute,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingtimeout",
		Description: "Reply to pings by timing the user out",
		Usage:       "[p]pingtimeout <toggle/list/clear> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingtimeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pinginsult",
		Description: "Reply to pings by insulting the user",
		Usage:       "[p]pinginsult <toggle/list/clear> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pinginsult,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingresponse",
		Description: "Reply to pings with a custom response",
		Usage:       "[p]pingresponse <toggle/list/clear> [\"response\"] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingresponse,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingairesponse",
		Description: "Reply to pings with an AI response",
		Usage:       "[p]pingairesponse <toggle/list/clear> [\"system\"] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pingairesponse,
	})

	// Mirror channel/server
	router.RegisterCmd(&dgc.Command{
		Name:        "mirrorchannel",
		Description: "Mirror messages from one channel to another",
		Usage:       "[p]mirrorchannel <toggle/list/clear> [source channel] [dest channel]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.mirrorchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mirrorserver",
		Description: "Mirror messages from one server to another",
		Usage:       "[p]mirrorserver <toggle/list/clear> [source server] [dest server]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.mirrorserver,
	})

	// Message Reactions
	router.RegisterCmd(&dgc.Command{
		Name:        "responder",
		Description: "Reply to messages with a custom response",
		Usage:       "[p]responder <toggle/list/clear> [\"response\"] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.responder,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "airesponder",
		Description: "Reply to messages with an AI response",
		Usage:       "[p]airesponder <toggle/list/clear> [\"system\"] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.airesponder,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "copy",
		Description: "Copy a user's messages",
		Usage:       "[p]copy <toggle/list/clear> [@user] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.copy,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "replycopy",
		Description: "Replycopy a user's messages",
		Usage:       "[p]replycopy <toggle/list/clear> [@user] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.replycopy,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mock",
		Description: "Mock a user's messages",
		Usage:       "[p]mock <toggle/list/clear> <@user> [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.mock,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pinspam",
		Description: "Pin all of a user's messages",
		Usage:       "[p]pinspam <toggle/list/clear> [@user] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.pinspam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unpinspam",
		Description: "Unpin messages pinned by pinspam",
		Usage:       "[p]unpinspam",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.unpinspam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reactspam",
		Description: "React to all of a user's messages",
		Usage:       "[p]reactspam <toggle/list/clear> [@user] [reaction] [global/server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.reactspam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletemessages",
		Description: "Delete all of a user's messages",
		Usage:       "[p]deletemessages <toggle/list/clear> [@user] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.deletemessages,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "keywordnotifier",
		Description: "Get notified when someone mentions a keyword",
		Usage:       "[p]keywordnotifier <toggle/list/clear> [\"keyword\"] [global/server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.keywordnotifier,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removereactions",
		Description: "Delete all of a user's reactions",
		Usage:       "[p]removereactions <toggle/list/clear> [@user] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.removereactions,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "forcedisconnect",
		Description: "Disconnect a user from voice channels",
		Usage:       "[p]forcedisconnect <toggle/list/clear> [@user] [server/channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.forcedisconnect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "forcenick",
		Description: "Force change a user's nickname",
		Usage:       "[p]forcenick <toggle/list/clear> [@user] [nickname] [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.forcenick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "forcekick",
		Description: "Force kick a user",
		Usage:       "[p]forcekick <toggle/list/clear> [@user] [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.forcekick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "noleave",
		Description: "Prevent a user from leaving a group chat",
		Usage:       "[p]noleave <toggle/list/clear> [@user] [channel id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.noleave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antigroupchat",
		Description: "Prevent users from adding you to group chats",
		Usage:       "[p]antigroupchat <on/off> [silent leave: on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.antigroupchat,
	})

	// AFK

	router.RegisterCmd(&dgc.Command{
		Name:        "getafkenabled",
		Description: "Get the current AFK responder setting",
		Usage:       "[p]getafkenabled",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkenabled,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkenabled",
		Description: "Turn the AFK responder on or off",
		Usage:       "[p]afkenabled <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkenabled,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkonping",
		Description: "Get the current AFK on ping setting",
		Usage:       "[p]getafkonping",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkonping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkonping",
		Description: "Set if an AFK message should be sent on ping",
		Usage:       "[p]afkonping <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkonping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkondm",
		Description: "Get the current AFK on DM setting",
		Usage:       "[p]getafkondm",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkondm,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkondm",
		Description: "Set if an AFK message should be sent on DM",
		Usage:       "[p]afkondm <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkondm,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkdelay",
		Description: "Get the current AFK delay setting",
		Usage:       "[p]getafkdelay",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkdelay",
		Description: "Set delay before sending AFK message",
		Usage:       "[p]afkdelay <seconds>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkrandomdelay",
		Description: "Get the current AFK random delay setting",
		Usage:       "[p]getafkrandomdelay",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkrandomdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkrandomdelay",
		Description: "Set the AFK random delay setting",
		Usage:       "[p]afkrandomdelay <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkrandomdelay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkmustbeinactive",
		Description: "Get the current must be inactive setting",
		Usage:       "[p]getafkmustbeinactive",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkmustbeinactive,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkmustbeinactive",
		Description: "Only send AFK messages when inactive",
		Usage:       "[p]afkmustbeinactive <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkmustbeinactive,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkinactivetimeout",
		Description: "Get the current inactive timeout",
		Usage:       "[p]getafkinactivetimeout",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkinactivetimeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkinactivetimeout",
		Description: "Minimum inactive time before being AFK",
		Usage:       "[p]afkinactivetimeout <seconds>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkinactivetimeout,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafktyping",
		Description: "Get the current AFK typing setting",
		Usage:       "[p]getafktyping",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafktyping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afktyping",
		Description: "Set if the bot should show typing before sending message",
		Usage:       "[p]afktyping <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afktyping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkairesponse",
		Description: "Get the current AFK AI response setting",
		Usage:       "[p]getafkairesponse",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkairesponse,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkairesponse",
		Description: "AI responder will use AI to generate response",
		Usage:       "[p]afkairesponse <on/off>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkairesponse,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafkmessage",
		Description: "Get the current AFK message",
		Usage:       "[p]getafkmessage",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafkmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afkmessage",
		Description: "Set the AFK message",
		Usage:       "[p]afkmessage <text>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afkmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "getafktriggertype",
		Description: "Get the current AFK trigger type",
		Usage:       "[p]getafktriggertype",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.getafktriggertype,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afktriggertype",
		Description: "Set the AFK trigger type",
		Usage:       "[p]afktriggertype <type>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afktriggertype,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afktriggerusers",
		Description: "Add/remove/list users that trigger AFK message",
		Usage:       "[p]afktriggerusers <add/rm/list> [@user]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afktriggerusers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "afktriggerservers",
		Description: "Add/remove/list servers that trigger AFK message",
		Usage:       "[p]afktriggerservers <add/rm/list> [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     e.afktriggerservers,
	})
}

// Utils

func (EventReactionsCommands) toggleReaction(ctx *dgc.Ctx, reaction eventreactions.EventReaction) {
	// Delete the raction if it already exists
	if eventreactions.Exists(reaction) {
		eventreactions.Remove(reaction)
		embed.NewEmbed().
			SetDescription("Removed event reaction: " + reaction.Description()).
			Send(ctx)
		return
	}

	eventreactions.Add(reaction)

	embed.NewEmbed().
		SetDescription("Added event reaction: " + reaction.Description()).
		Send(ctx)
}

func (EventReactionsCommands) mapReactions(ctx *dgc.Ctx, reactions []eventreactions.EventReaction, errText, title string) {
	if len(reactions) == 0 {
		embed.NewEmbed().
			SetDescription(errText).
			Send(ctx)
		return
	}

	descriptions := utils.Map(reactions, func(reaction eventreactions.EventReaction) string {
		return reaction.Description()
	})

	embed.NewEmbed().
		SetDescription(title + ":\n\n" + strings.Join(descriptions, "\n")).
		Send(ctx)
}

func (EventReactionsCommands) getGuildOrChannel(ctx *dgc.Ctx, arg *dgc.Argument) (guild *discordgo.Guild, channel *discordgo.Channel, err error) {
	givenID := utils.Misc.Or(arg.AsChannelMentionID(), arg.Raw())

	if givenID == "" {
		// Use the current server/channel
		if ctx.Event.GuildID != "" {
			guild, err = utils.Discord.GetGuild(ctx.Session, ctx.Event.GuildID)
			if err != nil {
				return
			}
		} else {
			channel, err = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
			if err != nil {
				return
			}
		}
	} else {
		// Use the given server/channel
		guild, err = utils.Discord.GetGuild(ctx.Session, givenID)
		if err != nil {
			channel, err = utils.Discord.GetChannel(ctx.Session, givenID)
			if err != nil {
				return
			}
		}
	}

	return
}

// Commands

func (e EventReactionsCommands) listreactions(ctx *dgc.Ctx) {
	reactions := eventreactions.ActiveReactions

	e.mapReactions(ctx, reactions, "There are no active event reactions", "Active event reactions")
}

func (e EventReactionsCommands) clearreactions(ctx *dgc.Ctx) {
	eventreactions.Clear()

	embed.NewEmbed().
		SetDescription("Cleared all event reactions").
		Send(ctx)
}

// Ignore User

func (e EventReactionsCommands) ignoreuser(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewIgnoreUser(user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.IGNORE_USER)

		e.mapReactions(ctx, reactions, "There are no ignored users.", "Ignored users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.IGNORE_USER)

		embed.NewEmbed().
			SetDescription("Cleared all ignored users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

// Ignore Server

func (e EventReactionsCommands) ignoreserver(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guildID := utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID)

		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewIgnoreGuild(guild))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.IGNORE_GUILD)

		e.mapReactions(ctx, reactions, "There are no ignored servers.", "Ignored servers")
	case "clear":
		eventreactions.RemoveAll(eventreactions.IGNORE_GUILD)

		embed.NewEmbed().
			SetDescription("Cleared all ignored servers.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

// Ping reactions

func (e EventReactionsCommands) pingReaction(ctx *dgc.Ctx, reactionType eventreactions.EventReactionType, reaction string, i int) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(i))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		eventReaction := eventreactions.NewPingReaction(reactionType, guild, channel)

		eventReaction.System = utils.Arguments["[\"system\"]"]
		eventReaction.Response = utils.Arguments["[\"response\"]"]

		e.toggleReaction(ctx, eventReaction)
	case "list":
		reactions := eventreactions.GetAll(reactionType)

		e.mapReactions(ctx, reactions, fmt.Sprintf("There are no ping --> %s reactions.", reaction), fmt.Sprintf("Ping --> %s reactions", reaction))
	case "clear":
		eventreactions.RemoveAll(reactionType)

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Cleared all ping --> %s reactions.", reaction)).
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) pingping(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_PING, "ping", 1)
}

func (e EventReactionsCommands) pingkick(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_KICK, "kick", 1)
}

func (e EventReactionsCommands) pingban(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_BAN, "ban", 1)
}

func (e EventReactionsCommands) pingmute(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_MUTE, "mute", 1)
}

func (e EventReactionsCommands) pingtimeout(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_TIMEOUT, "timeout", 1)
}

func (e EventReactionsCommands) pinginsult(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_INSULT, "insult", 1)
}

func (e EventReactionsCommands) pingresponse(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_RESPONSE, "response", 2)
}

func (e EventReactionsCommands) pingairesponse(ctx *dgc.Ctx) {
	e.pingReaction(ctx, eventreactions.PING_AI_RESPONSE, "AI response", 2)
}

// Mirror channel/server

func (e EventReactionsCommands) mirrorchannel(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		srcChannel, srcOk := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
		destChannel, destOk := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(2))
		if !srcOk || !destOk {
			embed.NewEmbed().
				SetDescription("Could not find channel").
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMirrorChannel(srcChannel, destChannel))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.MIRROR_CHANNEL)

		e.mapReactions(ctx, reactions, "There are no mirrored channels.", "Mirrored channels")
	case "clear":
		eventreactions.RemoveAll(eventreactions.MIRROR_CHANNEL)

		embed.NewEmbed().
			SetDescription("Cleared all mirrored channels.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) mirrorserver(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		srcGuild, srcErr := utils.Discord.GetGuild(ctx.Session, utils.Arguments["[source server]"])
		destGuild, destErr := utils.Discord.GetGuild(ctx.Session, utils.Arguments["[dest server]"])

		if srcErr != nil || destErr != nil {
			embed.NewEmbed().
				SetDescription("Could not find server").
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMirrorGuild(srcGuild, destGuild))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.MIRROR_CHANNEL)

		e.mapReactions(ctx, reactions, "There are no mirrored servers.", "Mirrored servers")
	case "clear":
		eventreactions.RemoveAll(eventreactions.MIRROR_CHANNEL)

		embed.NewEmbed().
			SetDescription("Cleared all mirrored servers.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

// Message Reactions

func (e EventReactionsCommands) responder(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		reaction := eventreactions.NewMessageReaction(eventreactions.RESPONDER, guild, channel, nil)
		reaction.Response = utils.Arguments["[\"response\"]"]
		e.toggleReaction(ctx, reaction)
	case "list":
		reactions := eventreactions.GetAll(eventreactions.RESPONDER)

		e.mapReactions(ctx, reactions, "There are no active responders.", "Responders")
	case "clear":
		eventreactions.RemoveAll(eventreactions.RESPONDER)

		embed.NewEmbed().
			SetDescription("Cleared all responders.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) airesponder(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		reaction := eventreactions.NewMessageReaction(eventreactions.AI_RESPONDER, guild, channel, nil)
		reaction.System = utils.Arguments["[\"system\"]"]
		e.toggleReaction(ctx, reaction)
	case "list":
		reactions := eventreactions.GetAll(eventreactions.AI_RESPONDER)

		e.mapReactions(ctx, reactions, "There are no active AI responders.", "AI Responders")
	case "clear":
		eventreactions.RemoveAll(eventreactions.AI_RESPONDER)

		embed.NewEmbed().
			SetDescription("Cleared all AI responders.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) copy(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMessageReaction(eventreactions.COPY, guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.COPY)

		e.mapReactions(ctx, reactions, "There are no copied users.", "Copied users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.COPY)

		embed.NewEmbed().
			SetDescription("Cleared all copied users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) replycopy(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMessageReaction(eventreactions.REPLY_COPY, guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.REPLY_COPY)

		e.mapReactions(ctx, reactions, "There are no replycopied users.", "Replycopied users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.REPLY_COPY)

		embed.NewEmbed().
			SetDescription("Cleared all replycopied users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) mock(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMessageReaction(eventreactions.MOCK, guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.MOCK)

		e.mapReactions(ctx, reactions, "There are no mocked users.", "Mocked users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.MOCK)

		embed.NewEmbed().
			SetDescription("Cleared all mocked users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) pinspam(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMessageReaction(eventreactions.PIN_SPAM, guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.PIN_SPAM)

		e.mapReactions(ctx, reactions, "There are no pinspammed users.", "Pinspammed users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.PIN_SPAM)

		embed.NewEmbed().
			SetDescription("Cleared all pinspammed users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) unpinspam(ctx *dgc.Ctx) {
	eventreactions.PinSpamMessages.ForEach(func(message *discordgo.Message) {
		err := ctx.Session.ChannelMessageUnpin(message.ChannelID, message.ID)
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not unpin message: " + err.Error()).
				Send(ctx)
			return
		}
	})

	eventreactions.PinSpamMessages.Clear()

	embed.NewEmbed().
		SetDescription("Unpinned all messages.").
		Send(ctx)
}

func (e EventReactionsCommands) reactspam(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		arg := strings.ToLower(utils.Arguments["[global/server/channel id]"])
		global := arg == "global" || arg == ""

		var guild *discordgo.Guild
		var channel *discordgo.Channel
		var err error

		if !global {
			guild, channel, err = e.getGuildOrChannel(ctx, ctx.Arguments.Get(3))
			if err != nil {
				embed.NewEmbed().
					SetDescription("Could not find server/channel: " + err.Error()).
					Send(ctx)
				return
			}
		}

		reaction := eventreactions.NewMessageReaction(eventreactions.REACT_SPAM, guild, channel, user)
		reaction.Global = global
		reaction.Reaction = utils.Arguments["[reaction]"]
		reaction.Reaction = strings.ReplaceAll(reaction.Reaction, "<", "")
		reaction.Reaction = strings.ReplaceAll(reaction.Reaction, ">", "")

		if reaction.Reaction == "" {
			_ = ctx.RespondText("Reaction cannot be empty!")
			return
		}

		e.toggleReaction(ctx, reaction)
	case "list":
		reactions := eventreactions.GetAll(eventreactions.REACT_SPAM)

		e.mapReactions(ctx, reactions, "There are no reactspammed users.", "Reactspammed users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.REACT_SPAM)

		embed.NewEmbed().
			SetDescription("Cleared all reactspammed users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) deletemessages(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewMessageReaction(eventreactions.DELETE_MESSAGES, guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.DELETE_MESSAGES)

		e.mapReactions(ctx, reactions, "There are no deletemessage reactions.", "Deletemessage reactions")
	case "clear":
		eventreactions.RemoveAll(eventreactions.DELETE_MESSAGES)

		embed.NewEmbed().
			SetDescription("Cleared all deletemessage reactions.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) keywordnotifier(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		arg := strings.ToLower(utils.Arguments["[global/server/channel id]"])
		global := arg == "global" || arg == ""

		var guild *discordgo.Guild
		var channel *discordgo.Channel
		var err error

		if !global {
			guild, channel, err = e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
			if err != nil {
				embed.NewEmbed().
					SetDescription("Could not find server/channel: " + err.Error()).
					Send(ctx)
				return
			}
		}

		reaction := eventreactions.NewMessageReaction(eventreactions.KEYWORD_NOTIFIER, guild, channel, nil)
		reaction.Keyword = utils.Arguments["[\"keyword\"]"]
		reaction.Global = global

		e.toggleReaction(ctx, reaction)
	case "list":
		reactions := eventreactions.GetAll(eventreactions.KEYWORD_NOTIFIER)

		e.mapReactions(ctx, reactions, "There are no keyword notifiers.", "Keyword notifiers")
	case "clear":
		eventreactions.RemoveAll(eventreactions.KEYWORD_NOTIFIER)

		embed.NewEmbed().
			SetDescription("Cleared all keyword notifiers.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) removereactions(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewRemoveReactions(guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.REMOVE_REACTIONS)

		e.mapReactions(ctx, reactions, "There are no users with removed reactions.", "Users with removed reactions")
	case "clear":
		eventreactions.RemoveAll(eventreactions.REMOVE_REACTIONS)

		embed.NewEmbed().
			SetDescription("Cleared all users with removed reactions.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) forcedisconnect(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guild, channel, err := e.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server/channel: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewForceDisconnect(guild, channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.FORCE_DISCONNECT)

		e.mapReactions(ctx, reactions, "There are no force disconnected users.", "Force disconnected users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.FORCE_DISCONNECT)

		embed.NewEmbed().
			SetDescription("Cleared all force disconnected users").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument").
			Send(ctx)
	}
}

func (e EventReactionsCommands) forcenick(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guildID := utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID)

		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server: " + err.Error()).
				Send(ctx)
			return
		}

		nickname := utils.Arguments["[nickname]"]

		err = ctx.Session.GuildMemberNickname(ctx.Event.GuildID, user.ID, nickname)
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not set nickname: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewForceNick(guild, user, nickname))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.FORCE_NICK)

		e.mapReactions(ctx, reactions, "There are no force nicked users.", "Force nicked users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.FORCE_NICK)

		embed.NewEmbed().
			SetDescription("Cleared all force nicked users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

func (e EventReactionsCommands) forcekick(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		guildID := utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID)

		guild, err := utils.Discord.GetGuild(ctx.Session, guildID)
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server: " + err.Error()).
				Send(ctx)
			return
		}

		err = ctx.Session.GuildMemberDelete(ctx.Event.GuildID, user.ID)
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not kick user: " + err.Error()).
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewForceKick(guild, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.FORCE_KICK)

		e.mapReactions(ctx, reactions, "There are no force kicked users.", "Force kicked users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.FORCE_KICK)

		embed.NewEmbed().
			SetDescription("Cleared all force kicked users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

func (e EventReactionsCommands) noleave(ctx *dgc.Ctx) {
	switch strings.ToLower(utils.Arguments["<toggle/list/clear>"]) {
	case "toggle":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user").
				Send(ctx)
			return
		}

		relationships := utils.Filter(ctx.Session.State.Relationships, func(r *discordgo.Relationship) bool {
			return r.Type == 1 && r.ID == user.ID
		})

		if len(relationships) == 0 {
			embed.NewEmbed().
				SetDescription("The user must be a friend.").
				Send(ctx)
			return
		}

		arg := ctx.Arguments.Get(2)
		givenID := utils.Misc.Or(arg.AsChannelMentionID(), arg.Raw())

		channel, err := utils.Discord.GetChannel(ctx.Session, utils.Misc.Or(givenID, ctx.Event.ChannelID))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find channel: " + err.Error()).
				Send(ctx)
			return
		}

		if channel.Type != discordgo.ChannelTypeGroupDM {
			embed.NewEmbed().
				SetDescription("The channel must be a group DM.").
				Send(ctx)
			return
		}

		e.toggleReaction(ctx, eventreactions.NewNoLeave(channel, user))
	case "list":
		reactions := eventreactions.GetAll(eventreactions.NO_LEAVE)

		e.mapReactions(ctx, reactions, "There are no noleave users.", "Noleave users")
	case "clear":
		eventreactions.RemoveAll(eventreactions.NO_LEAVE)

		embed.NewEmbed().
			SetDescription("Cleared all noleave users.").
			Send(ctx)
	default:
		embed.NewEmbed().
			SetDescription("Invalid toggle/list/clear argument.").
			Send(ctx)
	}
}

func (e EventReactionsCommands) antigroupchat(ctx *dgc.Ctx) {
	var toggle bool
	var silentToggle bool

	rawToggle := strings.ToLower(utils.Arguments["<on/off>"])
	if rawToggle != "on" && rawToggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	toggle = rawToggle == "on"

	rawSilentToggle := strings.ToLower(utils.Arguments["[silent leave: on/off]"])
	silentToggle = rawSilentToggle == "on"

	reaction := eventreactions.NewAntiGroupChat(silentToggle)

	if toggle {
		// Turn on (first remove to reset silent argument)
		eventreactions.Remove(reaction)
		eventreactions.Add(reaction)

		embed.NewEmbed().
			SetDescription("Turned on the anti group chat reaction!").
			Send(ctx)
	} else {
		// Turn off
		eventreactions.Remove(reaction)

		embed.NewEmbed().
			SetDescription("Turned off the anti group chat reaction!").
			Send(ctx)
	}
}

// AFK

func (EventReactionsCommands) getafkenabled(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK responder is currently %s", utils.Misc.FormatBoolOnOff(cfg.AFK.Enabled))).
		Send(ctx)
}

func (EventReactionsCommands) afkenabled(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.Enabled = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " the AFK responder.").
		Send(ctx)
}

func (EventReactionsCommands) getafkonping(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK responder is currently %s on pings", utils.Misc.FormatBoolOnOff(cfg.AFK.OnPing))).
		Send(ctx)
}

func (EventReactionsCommands) afkonping(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.OnPing = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " the AFK responder on pings.").
		Send(ctx)
}

func (EventReactionsCommands) getafkondm(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK responder is currently %s on DMs", utils.Misc.FormatBoolOnOff(cfg.AFK.OnDM))).
		Send(ctx)
}

func (EventReactionsCommands) afkondm(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.OnDM = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " the AFK responder on DMs.").
		Send(ctx)
}

func (EventReactionsCommands) getafkdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK delay is currently %d seconds", cfg.AFK.Delay)).
		Send(ctx)
}

func (EventReactionsCommands) afkdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	delay, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	cfg.AFK.Delay = delay
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK delay set to %d seconds", delay)).
		Send(ctx)
}

func (EventReactionsCommands) getafkrandomdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK random delay is currently %s", utils.Misc.FormatBoolOnOff(cfg.AFK.RandomDelay))).
		Send(ctx)
}

func (EventReactionsCommands) afkrandomdelay(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.RandomDelay = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " random delay.").
		Send(ctx)
}

func (EventReactionsCommands) getafkmustbeinactive(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK must be inactive is currently %s", utils.Misc.FormatBoolOnOff(cfg.AFK.MustBeInactive))).
		Send(ctx)
}

func (EventReactionsCommands) afkmustbeinactive(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.MustBeInactive = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " must be inactive.").
		Send(ctx)
}

func (EventReactionsCommands) getafkinactivetimeout(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK inactive timeout is currently %d seconds", cfg.AFK.InactiveTimeout)).
		Send(ctx)
}

func (EventReactionsCommands) afkinactivetimeout(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	timeout, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	cfg.AFK.InactiveTimeout = timeout
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK inactive timeout set to %d seconds", timeout)).
		Send(ctx)
}

func (EventReactionsCommands) getafktyping(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK typing is currently %s", utils.Misc.FormatBoolOnOff(cfg.AFK.Typing))).
		Send(ctx)
}

func (EventReactionsCommands) afktyping(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.Typing = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " typing.").
		Send(ctx)
}

func (EventReactionsCommands) getafkairesponse(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription(fmt.Sprintf("AFK AI response is currently %s", utils.Misc.FormatBoolOnOff(cfg.AFK.Typing))).
		Send(ctx)
}

func (EventReactionsCommands) afkairesponse(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	cfg.AFK.AI = toggle == "on"
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("Turned " + toggle + " AI response.").
		Send(ctx)
}

func (EventReactionsCommands) getafkmessage(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription("AFK message is currently " + cfg.AFK.Message).
		Send(ctx)
}

func (EventReactionsCommands) afkmessage(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	cfg.AFK.Message = ctx.Arguments.Raw()
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("AFK message set to " + cfg.AFK.Message).
		Send(ctx)
}

func (EventReactionsCommands) getafktriggertype(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	embed.NewEmbed().
		SetDescription("AFK trigger type is currently " + cfg.AFK.TriggerList.Type).
		Send(ctx)
}

func (EventReactionsCommands) afktriggertype(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	arg := strings.ToLower(utils.Arguments["<type>"])
	if !utils.SliceContains([]string{"global", "whitelist", "blacklist"}, arg) {
		_ = ctx.RespondText("Invalid value. Use `Global`, `Whitelist`, or `Blacklist`.")
		return
	}

	cfg.AFK.TriggerList.Type = utils.String.Capitalize(arg)
	_ = cfg.Save()

	embed.NewEmbed().
		SetDescription("AFK trigger type set to " + cfg.AFK.Message).
		Send(ctx)
}

func (EventReactionsCommands) afktriggerusers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch strings.ToLower(utils.Arguments["<add/rm/list>"]) {
	case "add":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user.").
				Send(ctx)
			return
		}

		cfg.AFK.TriggerList.Users = append(cfg.AFK.TriggerList.Users, user.ID)

		err = cfg.Save()
		if err != nil {
			embed.NewEmbed().
				SetDescription("Error saving config: " + err.Error()).
				Send(ctx)
			return
		}

	case "rm":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user.").
				Send(ctx)
			return
		}

		cfg.AFK.TriggerList.Users = utils.RemoveSliceItem(cfg.AFK.TriggerList.Users, user.ID)

		err = cfg.Save()
		if err != nil {
			embed.NewEmbed().
				SetDescription("Error saving config: " + err.Error()).
				Send(ctx)
			return
		}
	case "list":
		var users []string

		for _, id := range cfg.AFK.TriggerList.Users {
			user, ok := utils.Discord.GetUser(ctx.Session, id)
			if ok {
				users = append(users, "- "+user.Username+" ("+id+")")
			} else {
				users = append(users, "- "+id)
			}
		}

		embed.NewEmbed().
			AddField("AFK Trigger Users", strings.Join(users, "\n")).
			Send(ctx)

	default:
		embed.NewEmbed().
			SetDescription("Invalid add/rm/list argument.").
			Send(ctx)
	}
}

func (EventReactionsCommands) afktriggerservers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch strings.ToLower(utils.Arguments["<add/rm/list>"]) {
	case "add":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server.").
				Send(ctx)
			return
		}

		cfg.AFK.TriggerList.Servers = append(cfg.AFK.TriggerList.Servers, guild.ID)

		err = cfg.Save()
		if err != nil {
			embed.NewEmbed().
				SetDescription("Error saving config: " + err.Error()).
				Send(ctx)
			return
		}

	case "rm":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
		if err != nil {
			embed.NewEmbed().
				SetDescription("Could not find server.").
				Send(ctx)
			return
		}

		cfg.AFK.TriggerList.Servers = utils.RemoveSliceItem(cfg.AFK.TriggerList.Servers, guild.ID)

		err = cfg.Save()
		if err != nil {
			embed.NewEmbed().
				SetDescription("Error saving config: " + err.Error()).
				Send(ctx)
			return
		}
	case "list":
		var servers []string

		for _, s := range cfg.AFK.TriggerList.Servers {
			if name := utils.Discord.GetGuildName(ctx.Session, s); name != "No Server" {
				servers = append(servers, "- "+name+" ("+s+")")
			} else {
				servers = append(servers, "- "+s)
			}
		}

		embed.NewEmbed().
			AddField("AFK Trigger Servers", strings.Join(servers, "\n")).
			Send(ctx)

	default:
		embed.NewEmbed().
			SetDescription("Invalid add/rm/list argument.").
			Send(ctx)
	}
}
