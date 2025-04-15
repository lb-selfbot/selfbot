package protections

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/actioncommands"
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

type ProtectionsCommands struct{ services *services.Services }

func NewProtectionsCommands(s *services.Services) *ProtectionsCommands {
	return &ProtectionsCommands{services: s}
}

func (c *ProtectionsCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Protections",
		Description: "Protect your account",
	}, []*router.Command{
		// Anti-Phishing Protection
		actioncommands.NewBuilder("antiphishing", "Anti Phishing").
			WithValueGetter("Get anti phishing status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.AntiPhishing, nil
			}).
			WithValueSetter("Set anti phishing status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.AntiPhishing = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti phishing", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.AntiPhishing = !cfg.Protections.AntiPhishing
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Anti Friend Request Protection
		actioncommands.NewBuilder("antifriendrequest", "Anti Friend Request Spam").
			WithValueGetter("Get anti friend request status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.AntiFriendRequest, nil
			}).
			WithValueSetter("Set anti friend request status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.AntiFriendRequest = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti friend request", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.AntiFriendRequest = !cfg.Protections.AntiFriendRequest
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Anti Mass DM Protection
		actioncommands.NewBuilder("antimassdm", "Anti Mass DM / DM Spam").
			WithValueGetter("Get anti mass DM status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.AntiMassDM, nil
			}).
			WithValueSetter("Set anti mass DM status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.AntiMassDM = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti mass DM", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.AntiMassDM = !cfg.Protections.AntiMassDM
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Protected Servers List
		actioncommands.NewBuilder("protectedservers", "Protected Servers").
			WithListGetter("Get protected servers", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Protections.ServerProtections.ProtectedServers) == 0 {
					c.services.Embed.Respond(ctx, "No protected servers found.")
					return nil
				}

				var servers []string
				for _, serverID := range cfg.Protections.ServerProtections.ProtectedServers {
					guild, err := ctx.Session.GuildWithFallback(serverID)
					if err != nil {
						servers = append(servers, serverID)
						continue
					}
					servers = append(servers, fmt.Sprintf("%s (%s)", guild.Name, guild.ID))
				}

				c.services.Embed.RespondField(ctx, "Protected Servers", strings.Join(servers, "\n"))
				return nil
			}).
			WithListAdder("Add protected server", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
					return fmt.Errorf("server %s is already in the protected servers list", guild.Name)
				}

				cfg.Protections.ServerProtections.ProtectedServers = append(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added server %s to protected servers.", guild.Name)
				return nil
			}).
			WithListRemover("Remove protected server", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
					return fmt.Errorf("server %s is not in the protected servers list", guild.Name)
				}

				cfg.Protections.ServerProtections.ProtectedServers = fn.Delete(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed server %s from protected servers.", guild.Name)
				return nil
			}).
			WithListClearer("Clear protected servers", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.ProtectedServers = []string{}
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all protected servers.")
				return nil
			}).
			Build(c.services),

		// Anti Invite Protection
		actioncommands.NewBuilder("antiinvite", "Anti Invite (Server Protections)").
			WithValueGetter("Get anti invite status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.ServerProtections.AntiInvite, nil
			}).
			WithValueSetter("Set anti invite status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiInvite = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti invite", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiInvite = !cfg.Protections.ServerProtections.AntiInvite
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Anti Spam Protection
		actioncommands.NewBuilder("antispam", "Anti Spam (Server Protections)").
			WithValueGetter("Get anti spam status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.ServerProtections.AntiSpam, nil
			}).
			WithValueSetter("Set anti spam status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiSpam = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti spam", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiSpam = !cfg.Protections.ServerProtections.AntiSpam
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Anti Upper Protection
		actioncommands.NewBuilder("antiupper", "Anti Upper (Server Protections)").
			WithValueGetter("Get anti upper status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.ServerProtections.AntiUpper, nil
			}).
			WithValueSetter("Set anti upper status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiUpper = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti upper", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiUpper = !cfg.Protections.ServerProtections.AntiUpper
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Anti Delete Protection
		actioncommands.NewBuilder("antidelete", "Anti Delete (Server Protections)").
			WithValueGetter("Get anti delete status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Protections.ServerProtections.AntiDelete, nil
			}).
			WithValueSetter("Set anti delete status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiDelete = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle anti delete", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Protections.ServerProtections.AntiDelete = !cfg.Protections.ServerProtections.AntiDelete
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// DM Control
		{
			Name:        "dm",
			Description: "Open or close DM's from everyone",
			Arguments:   "<open/close>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.dm,
		},

		// Friend Requests Control
		{
			Name:        "friendrequests",
			Description: "Allow or deny new friend requests",
			Arguments:   "<allow/deny>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.friendrequests,
		},
	})
}

func (c *ProtectionsCommands) dm(ctx *router.Ctx) {
	var err error
	var action string

	switch strings.ToLower(ctx.Arguments.Get(0).Raw()) {
	case "open":
		ctx.Session.State.UserSettings.Privacy.RestrictedGuildIds = []uint64{}
		err = ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)

		action = "Opened"
	case "close":
		var guildIDs []uint64
		for _, g := range ctx.Session.State.Guilds {
			id, _ := strconv.Atoi(g.ID)

			guildIDs = append(guildIDs, uint64(id))
		}

		ctx.Session.State.UserSettings.Privacy.RestrictedGuildIds = guildIDs
		err = ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)

		action = "Closed"
	default:
		ctx.RespondText("Invalid <open/close> argument. Use `open` or `close`.")
		return
	}

	if err != nil {
		ctx.RespondText("Error changing settings: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("%s DM's from everyone!", action).
		Send()
}

func (c *ProtectionsCommands) friendrequests(ctx *router.Ctx) {
	var err error
	var action string

	switch strings.ToLower(ctx.Arguments.Get(0).Raw()) {
	case "allow":
		err = c.services.Discord.SetFriendRequestsAllowed(true)
		action = "Allowed"
	case "deny":
		err = c.services.Discord.SetFriendRequestsAllowed(false)
		action = "Denied"
	default:
		ctx.RespondText("Invalid <allow/deny> argument. Use `allow` or `deny`.")
		return
	}

	if err != nil {
		ctx.RespondText("Error changing settings: %v", err)
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("%s friend requests from everyone!", action).
		Send()
}
