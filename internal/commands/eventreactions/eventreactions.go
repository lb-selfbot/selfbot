package eventreactions

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/actioncommands"
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/impl"
	"github.com/lb-selfbot/selfbot/internal/services/eventreactions/types"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type EventReactionsCommands struct{ services *services.Services }

func NewEventReactionsCommands(s *services.Services) *EventReactionsCommands {
	return &EventReactionsCommands{services: s}
}

func (c *EventReactionsCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "EventReactions",
		Description: "Event Reaction commands",
	}, []*router.Command{
		{
			Name:        "listreactions",
			Description: "List all active event reactions",
			Handler:     c.listreactions,
		},
		{
			Name:        "clearreactions",
			Description: "Clear all active event reactions",
			Handler:     c.clearreactions,
		},

		actioncommands.NewBuilder("ignoreuser", "Ignored Users").
			WithEventReactionToggle("Ignore messages from a user", "<user>", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, err
				}
				return impl.NewIgnoreUser(user), nil
			}).
			WithEventReactionList(types.IgnoreUser, "List all ignored users").
			WithEventReactionClear(types.IgnoreUser, "Clear all ignored users").
			Build(c.services),

		actioncommands.NewBuilder("ignoreserver", "Ignored Servers").
			WithEventReactionToggle("Ignore messages from a server", "[server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, err
				}
				return impl.NewIgnoreGuild(user), nil
			}).
			WithEventReactionList(types.IgnoreUser, "List all ignored servers").
			WithEventReactionClear(types.IgnoreUser, "Clear all ignored servers").
			Build(c.services),

		actioncommands.NewBuilder("pingping", "Ping → Ping").
			WithEventReactionToggle("Reply to pings by pinging the user", "[channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(0))
				if err != nil {
					return nil, err
				}
				return impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingPing, guild, channel), nil
			}).
			WithEventReactionList(types.PingPing, "List all ping → ping reactions").
			WithEventReactionClear(types.PingPing, "Clear all ping → ping reactions").
			Build(c.services),

		actioncommands.NewBuilder("pingkick", "Ping → Kick").
			WithEventReactionToggle("Reply to pings by kicking the user", "[channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(0))
				if err != nil {
					return nil, err
				}
				return impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingKick, guild, channel), nil
			}).
			WithEventReactionList(types.PingKick, "List all ping → kick reactions").
			WithEventReactionClear(types.PingKick, "Clear all ping → kick reactions").
			Build(c.services),

		actioncommands.NewBuilder("pingban", "Ping → Ban").
			WithEventReactionToggle("Reply to pings by banning the user", "[channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(0))
				if err != nil {
					return nil, err
				}
				return impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingBan, guild, channel), nil
			}).
			WithEventReactionList(types.PingBan, "List all ping → ban reactions").
			WithEventReactionClear(types.PingBan, "Clear all ping → ban reactions").
			Build(c.services),

		actioncommands.NewBuilder("pingtimeout", "Ping → Timeout").
			WithEventReactionToggle("Reply to pings by timing the user out", "[channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(0))
				if err != nil {
					return nil, err
				}
				return impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingTimeout, guild, channel), nil
			}).
			WithEventReactionList(types.PingTimeout, "List all ping → timeout reactions").
			WithEventReactionClear(types.PingTimeout, "Clear all ping → timeout reactions").
			Build(c.services),

		actioncommands.NewBuilder("pinginsult", "Ping → Insult").
			WithEventReactionToggle("Reply to pings by insulting the user", "[channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(0))
				if err != nil {
					return nil, err
				}
				return impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingInsult, guild, channel), nil
			}).
			WithEventReactionList(types.PingInsult, "List all ping → insult reactions").
			WithEventReactionClear(types.PingInsult, "Clear all ping → insult reactions").
			Build(c.services),

		actioncommands.NewBuilder("pingreact", "Ping → React").
			WithEventReactionToggle("Reply to pings by adding a reaction", "[reaction] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				reaction := ctx.Arguments.Get(0).Raw()
				reaction = strings.ReplaceAll(reaction, "<", "")
				reaction = strings.ReplaceAll(reaction, ">", "")

				if reaction == "" {
					return nil, fmt.Errorf("reaction cannot be empty")
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				pingReaction := impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingReact, guild, channel)
				pingReaction.Reaction = reaction
				return pingReaction, nil
			}).
			WithEventReactionList(types.PingReact, "List all ping → react reactions").
			WithEventReactionClear(types.PingReact, "Clear all ping → react reactions").
			Build(c.services),

		actioncommands.NewBuilder("pingresponse", "Ping → Response").
			WithEventReactionToggle("Reply to pings with a custom response", "[\"response\"] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				response := ctx.Arguments.Get(0).Raw()

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
				if err != nil {
					return nil, err
				}

				pingReaction := impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingResponse, guild, channel)
				pingReaction.Response = response
				return pingReaction, nil
			}).
			WithEventReactionList(types.PingResponse, "List all ping → response reactions").
			WithEventReactionClear(types.PingResponse, "Clear all ping → response reactions").
			Build(c.services),

		actioncommands.NewBuilder("pingairesponse", "Ping → AI Response").
			WithEventReactionToggle("Reply to pings with an AI response", "[\"system\"] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				system := ctx.Arguments.Get(0).Raw()

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				pingReaction := impl.NewPingReaction(c.services.EventBroker, c.services.Config, types.PingAIResponse, guild, channel)
				pingReaction.System = system
				return pingReaction, nil
			}).
			WithEventReactionList(types.PingAIResponse, "List all ping → AI response reactions").
			WithEventReactionClear(types.PingAIResponse, "Clear all ping → AI response reactions").
			Build(c.services),

		actioncommands.NewBuilder("mirrorchannel", "Mirror Channel").
			WithEventReactionToggle("Mirror messages from one channel to another", "[source channel] [dest channel]", func(ctx *router.Ctx) (types.EventReaction, error) {
				sourceChannel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find source channel: %v", err)
				}

				destChannel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find destination channel: %v", err)
				}

				return impl.NewMirrorChannel(c.services.EventBroker, sourceChannel, destChannel), nil
			}).
			WithEventReactionList(types.MirrorChannel, "List all mirrored channels").
			WithEventReactionClear(types.MirrorChannel, "Clear all mirrored channels").
			Build(c.services),

		actioncommands.NewBuilder("mirrorserver", "Mirror Server").
			WithEventReactionToggle("Mirror messages from one server to another", "[source server] [dest server]", func(ctx *router.Ctx) (types.EventReaction, error) {
				sourceGuild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find source server: %v", err)
				}

				destGuild, err := ctx.Arguments.Get(1).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find destination server: %v", err)
				}

				return impl.NewMirrorGuild(c.services.EventBroker, sourceGuild, destGuild), nil
			}).
			WithEventReactionList(types.MirrorGuild, "List all mirrored servers").
			WithEventReactionClear(types.MirrorGuild, "Clear all mirrored servers").
			Build(c.services),

		actioncommands.NewBuilder("responder", "Responder").
			WithEventReactionToggle("Reply to messages with a custom response", "[\"response\"] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				response := ctx.Arguments.Get(0).Raw()

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				reaction := impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.Responder, guild, channel, nil)
				reaction.Response = response
				return reaction, nil
			}).
			WithEventReactionList(types.Responder, "List all responders").
			WithEventReactionClear(types.Responder, "Clear all responders").
			Build(c.services),

		actioncommands.NewBuilder("airesponder", "AI Responder").
			WithEventReactionToggle("Reply to messages with an AI response", "[\"system\"] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				system := ctx.Arguments.Get(0).Raw()

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				reaction := impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.AIResponder, guild, channel, nil)
				reaction.System = system
				return reaction, nil
			}).
			WithEventReactionList(types.AIResponder, "List all AI responders").
			WithEventReactionClear(types.AIResponder, "Clear all AI responders").
			Build(c.services),

		actioncommands.NewBuilder("copy", "Copy User").
			WithEventReactionToggle("Copy a user's messages", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.Copy, guild, channel, user), nil
			}).
			WithEventReactionList(types.Copy, "List all copied users").
			WithEventReactionClear(types.Copy, "Clear all copied users").
			Build(c.services),

		actioncommands.NewBuilder("replycopy", "Reply Copy User").
			WithEventReactionToggle("Reply copy a user's messages", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.ReplyCopy, guild, channel, user), nil
			}).
			WithEventReactionList(types.ReplyCopy, "List all reply copied users").
			WithEventReactionClear(types.ReplyCopy, "Clear all reply copied users").
			Build(c.services),

		actioncommands.NewBuilder("mock", "Mock User").
			WithEventReactionToggle("Mock a user's messages", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.Mock, guild, channel, user), nil
			}).
			WithEventReactionList(types.Mock, "List all mocked users").
			WithEventReactionClear(types.Mock, "Clear all mocked users").
			Build(c.services),

		actioncommands.NewBuilder("pinspam", "Pin Spam").
			WithEventReactionToggle("Pin all of a user's messages", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.PinSpam, guild, channel, user), nil
			}).
			WithEventReactionList(types.PinSpam, "List all pinspammed users").
			WithEventReactionClear(types.PinSpam, "Clear all pinspammed users").
			WithCustomAction("undo", "Unpin all messages pinned by pinspam", "", c.unpinspam).
			Build(c.services),

		actioncommands.NewBuilder("reactspam", "React Spam").
			WithEventReactionToggle("React to all of a user's messages", "[@user] [reaction] [global/server/channel id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				reaction := ctx.Arguments.Get(1).Raw()
				reaction = strings.ReplaceAll(reaction, "<", "")
				reaction = strings.ReplaceAll(reaction, ">", "")

				if reaction == "" {
					return nil, fmt.Errorf("reaction cannot be empty")
				}

				arg := strings.ToLower(ctx.Arguments.Get(2).Raw())
				global := arg == "global" || arg == ""

				var guild *discordgo.Guild
				var channel *discordgo.Channel

				if !global {
					guild, channel, err = c.getGuildOrChannel(ctx, ctx.Arguments.Get(2))
					if err != nil {
						return nil, err
					}
				}

				messageReaction := impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.ReactSpam, guild, channel, user)
				messageReaction.Global = global
				messageReaction.Reaction = reaction
				return messageReaction, nil
			}).
			WithEventReactionList(types.ReactSpam, "List all reactspammed users").
			WithEventReactionClear(types.ReactSpam, "Clear all reactspammed users").
			Build(c.services),

		actioncommands.NewBuilder("deletemessages", "Delete Messages").
			WithEventReactionToggle("Delete all of a user's messages", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.DeleteMessages, guild, channel, user), nil
			}).
			WithEventReactionList(types.DeleteMessages, "List all delete message reactions").
			WithEventReactionClear(types.DeleteMessages, "Clear all delete message reactions").
			Build(c.services),

		actioncommands.NewBuilder("keywordnotifier", "Keyword Notifier").
			WithEventReactionToggle("Get notified when someone mentions a keyword", "[\"keyword\"] [global/server/channel id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				keyword := ctx.Arguments.Get(0).Raw()

				arg := strings.ToLower(ctx.Arguments.Get(1).Raw())
				global := arg == "global" || arg == ""

				var guild *discordgo.Guild
				var channel *discordgo.Channel
				var err error

				if !global {
					guild, channel, err = c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
					if err != nil {
						return nil, err
					}
				}

				reaction := impl.NewMessageReaction(c.services.EventBroker, c.services.Config, types.KeywordNotifier, guild, channel, nil)
				reaction.Keyword = keyword
				reaction.Global = global

				return reaction, nil
			}).
			WithEventReactionList(types.KeywordNotifier, "List all keyword notifiers").
			WithEventReactionClear(types.KeywordNotifier, "Clear all keyword notifiers").
			Build(c.services),

		actioncommands.NewBuilder("removereactions", "Remove Reactions").
			WithEventReactionToggle("Delete all of a user's reactions", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewRemoveReactions(guild, channel, user), nil
			}).
			WithEventReactionList(types.RemoveReactions, "List all users with removed reactions").
			WithEventReactionClear(types.RemoveReactions, "Clear all users with removed reactions").
			Build(c.services),

		actioncommands.NewBuilder("forcedisconnect", "Force Disconnect").
			WithEventReactionToggle("Disconnect a user from voice channels", "[@user] [channel/server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guild, channel, err := c.getGuildOrChannel(ctx, ctx.Arguments.Get(1))
				if err != nil {
					return nil, err
				}

				return impl.NewForceDisconnect(guild, channel, user), nil
			}).
			WithEventReactionList(types.ForceDisconnect, "List all force disconnected users").
			WithEventReactionClear(types.ForceDisconnect, "Clear all force disconnected users").
			Build(c.services),

		actioncommands.NewBuilder("forcenick", "Force Nick").
			WithEventReactionToggle("Force change a user's nickname", "[@user] [nickname] [server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				nickname := ctx.Arguments.Get(1).Raw()

				guildID := utils.Misc.Or(ctx.Arguments.Get(2).Raw(), ctx.Message.GuildID)
				guild, err := ctx.Session.GuildWithFallback(guildID)
				if err != nil {
					return nil, fmt.Errorf("could not find server: %v", err)
				}

				err = ctx.Session.GuildMemberNickname(ctx.Message.GuildID, user.ID, nickname)
				if err != nil {
					return nil, fmt.Errorf("could not set nickname: %v", err)
				}

				return impl.NewForceNick(guild, user, nickname), nil
			}).
			WithEventReactionList(types.ForceNick, "List all force nicked users").
			WithEventReactionClear(types.ForceNick, "Clear all force nicked users").
			Build(c.services),

		actioncommands.NewBuilder("forcekick", "Force Kick").
			WithEventReactionToggle("Force kick a user", "[@user] [server id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				guildID := utils.Misc.Or(ctx.Arguments.Get(1).Raw(), ctx.Message.GuildID)
				guild, err := ctx.Session.GuildWithFallback(guildID)
				if err != nil {
					return nil, fmt.Errorf("could not find server: %v", err)
				}

				err = ctx.Session.GuildMemberDelete(ctx.Message.GuildID, user.ID)
				if err != nil {
					return nil, fmt.Errorf("could not kick user: %v", err)
				}

				return impl.NewForceKick(guild, user), nil
			}).
			WithEventReactionList(types.ForceKick, "List all force kicked users").
			WithEventReactionClear(types.ForceKick, "Clear all force kicked users").
			Build(c.services),

		actioncommands.NewBuilder("noleave", "No Leave").
			WithEventReactionToggle("Prevent a user from leaving a group chat", "[@user] [channel id]", func(ctx *router.Ctx) (types.EventReaction, error) {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find user: %v", err)
				}

				relationships := fn.Filter(ctx.Session.State.Relationships, func(r *discordgo.Relationship) bool {
					return r.Type == 1 && r.ID == user.ID
				})

				if len(relationships) == 0 {
					return nil, fmt.Errorf("the user must be a friend")
				}

				channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find channel: %v", err)
				}

				if channel.Type != discordgo.ChannelTypeGroupDM {
					return nil, fmt.Errorf("the channel must be a group DM")
				}

				return impl.NewNoLeave(channel, user), nil
			}).
			WithEventReactionList(types.NoLeave, "List all noleave users").
			WithEventReactionClear(types.NoLeave, "Clear all noleave users").
			Build(c.services),

		{
			Name:        "antigroupchat",
			Description: "Prevent users from adding you to group chats",
			Arguments:   "<on/off> [silent leave: on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler: func(ctx *router.Ctx) error {
				toggle, err := ctx.Arguments.Get(0).Bool()
				if err != nil {
					return fmt.Errorf("Invalid value. Use `on` or `off`.")
				}

				isSilent, _ := ctx.Arguments.Get(1).Bool()

				reaction := impl.NewAntiGroupChat(isSilent)

				// Remove to reset
				c.services.EventReactions.RemoveReactionsByType(reaction.Type())

				if toggle {
					c.services.EventReactions.AddReaction(reaction)

					c.services.Embed.New(ctx).
						SetDescription("Turned on the anti group chat reaction!").
						Send()
				} else {
					c.services.Embed.New(ctx).
						SetDescription("Turned off the anti group chat reaction!").
						Send()
				}

				return nil
			},
		},

		actioncommands.NewBuilder("antideletemessages", "Anti Delete Messages").
			WithEventReactionToggle("Prevent users from deleting your messages", "[channel]", func(ctx *router.Ctx) (types.EventReaction, error) {
				channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return nil, fmt.Errorf("could not find channel: %v", err)
				}

				return impl.NewAntiDeleteMessages(c.services.Config, channel), nil
			}).
			WithEventReactionList(types.AntiDeleteMessages, "List all anti delete messages").
			WithEventReactionClear(types.AntiDeleteMessages, "Clear all anti delete messages").
			Build(c.services),

		// AFK Related Commands
		actioncommands.NewBuilder("afkresponder", "AFK Responder").
			WithValueGetter("Get AFK responder status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.Enabled, nil
			}).
			WithValueSetter("Set AFK responder status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.Enabled = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK responder", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.Enabled = !cfg.AFK.Enabled
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkonping", "AFK On Ping").
			WithValueGetter("Get AFK on ping status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.OnPing, nil
			}).
			WithValueSetter("Set AFK on ping status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.OnPing = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK on ping", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.OnPing = !cfg.AFK.OnPing
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkondm", "AFK On DM").
			WithValueGetter("Get AFK on DM status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.OnDM, nil
			}).
			WithValueSetter("Set AFK on DM status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.OnDM = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK on DM", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.OnDM = !cfg.AFK.OnDM
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkdelay", "AFK Delay").
			WithValueGetter("Get AFK delay", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.Delay, nil
			}).
			WithValueSetter("Set AFK delay", actionvalue.TypeDuration, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.Delay = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkrandomdelay", "AFK Random Delay").
			WithValueGetter("Get AFK random delay status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.RandomDelay, nil
			}).
			WithValueSetter("Set AFK random delay status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.RandomDelay = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK random delay", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.RandomDelay = !cfg.AFK.RandomDelay
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkinactive", "AFK Inactive").
			WithValueGetter("Get AFK inactive status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.MustBeInactive, nil
			}).
			WithValueSetter("Set AFK inactive status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.MustBeInactive = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK inactive", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.MustBeInactive = !cfg.AFK.MustBeInactive
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkinactivetimeout", "AFK Inactive Timeout").
			WithValueGetter("Get AFK inactive timeout", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.InactiveTimeout, nil
			}).
			WithValueSetter("Set AFK inactive timeout", actionvalue.TypeDuration, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.InactiveTimeout = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afktyping", "AFK Typing").
			WithValueGetter("Get AFK typing status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.Typing, nil
			}).
			WithValueSetter("Set AFK typing status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.Typing = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK typing", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.Typing = !cfg.AFK.Typing
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkai", "AFK AI").
			WithValueGetter("Get AFK AI status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.AI, nil
			}).
			WithValueSetter("Set AFK AI status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.AI = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle AFK AI", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.AI = !cfg.AFK.AI
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afkmessage", "AFK Message").
			WithValueGetter("Get AFK message", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.Message, nil
			}).
			WithValueSetter("Set AFK message", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.Message = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		actioncommands.NewBuilder("afktriggertype", "AFK Trigger Type").
			WithValueGetter("Get AFK trigger type", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.AFK.TriggerList.Type, nil
			}).
			WithValueSetter("Set AFK trigger type", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				valueStr := strings.ToLower(value.(string))

				switch valueStr {
				case "global":
					cfg.AFK.TriggerList.Type = "Global"
				case "whitelist":
					cfg.AFK.TriggerList.Type = "Whitelist"
				case "blacklist":
					cfg.AFK.TriggerList.Type = "Blacklist"
				default:
					return fmt.Errorf("invalid trigger type. Use `global`, `whitelist`, or `blacklist`")
				}

				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available trigger types", func() []string {
				return []string{"Global", "Whitelist", "Blacklist"}
			}).
			Build(c.services),

		actioncommands.NewBuilder("afktriggerusers", "AFK Trigger Users").
			WithListGetter("Get AFK trigger users", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.AFK.TriggerList.Users) == 0 {
					c.services.Embed.Respond(ctx, "No AFK trigger users found.")
					return nil
				}

				users := fn.Map(cfg.AFK.TriggerList.Users, func(userID string) string {
					user, err := ctx.Session.UserWithFallback(userID)
					if err != nil {
						return userID
					}
					return fmt.Sprintf("%s - %s", user.DisplayName, userID)
				})

				c.services.Embed.RespondField(ctx, "AFK Triggerlist Users", strings.Join(users, "\n"))
				return nil
			}).
			WithListAdder("Add AFK trigger user", "<user>", func(ctx *router.Ctx) error {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find user: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.AFK.TriggerList.Users, user.ID) {
					return fmt.Errorf("user %s is already in the AFK Triggerlist", user.Username)
				}

				cfg.AFK.TriggerList.Users = append(cfg.AFK.TriggerList.Users, user.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added %s to the AFK Triggerlist.", user.Mention())
				return nil
			}).
			WithListRemover("Remove AFK trigger user", "<user>", func(ctx *router.Ctx) error {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find user: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.AFK.TriggerList.Users, user.ID) {
					return fmt.Errorf("user %s is not in the AFK Triggerlist", user.Username)
				}

				cfg.AFK.TriggerList.Users = fn.Delete(cfg.AFK.TriggerList.Users, user.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed %s from the AFK Triggerlist.", user.Mention())
				return nil
			}).
			WithListClearer("Clear AFK trigger users", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.TriggerList.Users = []string{}
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all users from the AFK Triggerlist.")
				return nil
			}).
			Build(c.services),

		actioncommands.NewBuilder("afktriggerservers", "AFK Trigger Servers").
			WithListGetter("Get AFK trigger servers", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.AFK.TriggerList.Servers) == 0 {
					c.services.Embed.Respond(ctx, "No AFK trigger servers found.")
					return nil
				}

				servers := fn.Map(cfg.AFK.TriggerList.Servers, func(serverID string) string {
					guild, err := ctx.Session.GuildWithFallback(serverID)
					if err != nil {
						return serverID
					}

					return fmt.Sprintf("%s - %s", guild.Name, guild.ID)
				})

				c.services.Embed.RespondField(ctx, "AFK Triggerlist Servers", strings.Join(servers, "\n"))
				return nil
			}).
			WithListAdder("Add AFK trigger server", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.AFK.TriggerList.Servers, guild.ID) {
					return fmt.Errorf("server %s is already in the AFK Triggerlist", guild.Name)
				}

				cfg.AFK.TriggerList.Servers = append(cfg.AFK.TriggerList.Servers, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added server %s to the AFK Triggerlist.", guild.Name)
				return nil
			}).
			WithListRemover("Remove AFK trigger server", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.AFK.TriggerList.Servers, guild.ID) {
					return fmt.Errorf("server %s is not in the AFK Triggerlist", guild.Name)
				}

				cfg.AFK.TriggerList.Servers = fn.Delete(cfg.AFK.TriggerList.Servers, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed server %s from the AFK Triggerlist.", guild.Name)
				return nil
			}).
			WithListClearer("Clear AFK trigger servers", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.AFK.TriggerList.Servers = []string{}
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all servers from the AFK Triggerlist.")
				return nil
			}).
			Build(c.services),
	})
}

// Utils

func (c *EventReactionsCommands) getGuildOrChannel(ctx *router.Ctx, arg *router.Argument) (guild *discordgo.Guild, channel *discordgo.Channel, err error) {
	// Check if the argument is a guild
	guild, err = arg.Guild(router.ArgumentFlagNoFetch)
	if err == nil {
		return guild, nil, nil
	}

	// Check if the argument is a channel
	channel, err = arg.Channel(router.ArgumentFlagNoFetch)
	if err == nil {
		return nil, channel, nil
	}

	// Fallback to channel if DM or guild if not
	if utils.Discord.IsDM(ctx.Session, ctx.Message) {
		channel, err = arg.Channel(router.ArgumentFlagFallbackDefault)
		return nil, channel, err
	} else {
		guild, err = arg.Guild(router.ArgumentFlagFallbackDefault)
		return guild, nil, err
	}
}

// Commands

func (c *EventReactionsCommands) listreactions(ctx *router.Ctx) error {
	reactions := c.services.EventReactions.GetAllReactions()
	if len(reactions) == 0 {
		c.services.Embed.Respond(ctx, "No active event reactions found.")
	}

	descriptions := fn.Map(reactions, func(reaction types.EventReaction) string {
		return reaction.Description()
	})

	c.services.Embed.RespondField(ctx, "Active Event Reactions", strings.Join(descriptions, "\n"))
	return nil
}

func (c *EventReactionsCommands) clearreactions(ctx *router.Ctx) error {
	c.services.EventReactions.RemoveAllReactions()
	c.services.Embed.Respond(ctx, "Cleared all event reactions")
	return nil
}

func (c *EventReactionsCommands) unpinspam(ctx *router.Ctx) error {
	for _, message := range c.services.EventReactions.GetPinSpamMessages() {
		err := ctx.Session.ChannelMessageUnpin(message.ChannelID, message.ID)
		if err != nil {
			return fmt.Errorf("could not unpin message: %v", err)
		}
	}

	c.services.EventReactions.ClearPinSpamMessages()
	c.services.Embed.Respond(ctx, "Unpinned all messages.")
	return nil
}
