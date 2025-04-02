package protections

import (
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
		// {
		// 	Name:        "antiphishing",
		// 	Description: "Get, set or toggle Anti Phishing",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antiphishing,
		// },
		// {
		// 	Name:        "antifriendrequest",
		// 	Description: "Get, set or toggle Anti Friend Request Spam",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antifriendrequest,
		// },
		// {
		// 	Name:        "antimassdm",
		// 	Description: "Get, set or toggle Anti Mass DM / DM Spam",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antimassdm,
		// },
		// {
		// 	Name:        "protectedservers",
		// 	Description: "Get, add, remove or clear the protected servers list",
		// 	Arguments:   "<get/add/rm/clear> [server id]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.protectedservers,
		// },
		// {
		// 	Name:        "antiinvite",
		// 	Description: "Get, set or toggle Anti Invite (Server Protections)",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antiinvite,
		// },
		// {
		// 	Name:        "antispam",
		// 	Description: "Get, set or toggle Anti Spam (Server Protections)",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antispam,
		// },
		// {
		// 	Name:        "antiupper",
		// 	Description: "Get, set or toggle Anti Upper (Server Protections)",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antiupper,
		// },
		// {
		// 	Name:        "antidelete",
		// 	Description: "Get, set or toggle Anti Delete (Server Protections)",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.antidelete,
		// },
		// {
		// 	Name:        "dm",
		// 	Description: "Open or close DM's from everyone",
		// 	Arguments:   "<open/close>",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.dm,
		// },
		// {
		// 	Name:        "friendrequests",
		// 	Description: "Allow or deny new friend requests",
		// 	Arguments:   "<allow/deny>",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.friendrequests,
		// },
	})
}

// func (c *ProtectionsCommands) antiphishing(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Phishing", utils.Misc.FormatBoolOnOff(cfg.Protections.AntiPhishing)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.AntiPhishing = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Phishing").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.AntiPhishing = !cfg.Protections.AntiPhishing
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.AntiPhishing) + " Anti Phishing").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) antifriendrequest(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Friend Request Spam", utils.Misc.FormatBoolOnOff(cfg.Protections.AntiFriendRequest)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.AntiFriendRequest = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Friend Request Spam").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.AntiFriendRequest = !cfg.Protections.AntiFriendRequest
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.AntiFriendRequest) + " Anti Friend Request Spam").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) antimassdm(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Mass DM / DM Spam", utils.Misc.FormatBoolOnOff(cfg.Protections.AntiMassDM)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.AntiMassDM = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Mass DM / DM Spam").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.AntiMassDM = !cfg.Protections.AntiMassDM
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.AntiMassDM) + " Anti Mass DM / DM Spam").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) protectedservers(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config.")
// 		return
// 	}

// 	switch utils.Arguments["<get/add/rm/clear>"] {
// 	case "get":
// 		var servers []string

// 		for _, serverID := range cfg.Protections.ServerProtections.ProtectedServers {
// 			guild, err := ctx.Session.GuildWithFallback(serverID)
// 			if err != nil {
// 				servers = append(servers, serverID)
// 				continue
// 			}

// 			servers = append(servers, guild.Name+" ("+guild.ID+")")
// 		}

// 		c.services.Embed.New(ctx).
// 			AddField("Protected Servers List", strings.Join(servers, ", ")).
// 			Send()
// 	case "add":
// 		guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
// 		if err != nil {
// 			ctx.RespondText("Could not find server.")
// 			return
// 		}

// 		if slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
// 			ctx.RespondText("Server %s is already in the Protected Servers List!", guild.Name)
// 			return
// 		}

// 		cfg.Protections.ServerProtections.ProtectedServers = append(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Added %s to the Protected Servers List.", guild.Name).
// 			Send()
// 	case "rm":
// 		guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
// 		if err != nil {
// 			ctx.RespondText("Could not find server.")
// 			return
// 		}

// 		if !slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
// 			ctx.RespondText("Server %s already isn't in the Protected Servers List!", guild.Name)
// 			return
// 		}

// 		cfg.Protections.ServerProtections.ProtectedServers = fn.Delete(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Removed %s from the Protected Servers List.", guild.Name).
// 			Send()
// 	case "clear":
// 		cfg.Protections.ServerProtections.ProtectedServers = []string{}
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Cleared all Protected Servers List entries.").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
// 	}
// }

// func (c *ProtectionsCommands) antiinvite(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Invite (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiInvite)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.ServerProtections.AntiInvite = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Invite (Server Protections)").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.ServerProtections.AntiInvite = !cfg.Protections.ServerProtections.AntiInvite
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiInvite) + " Anti Invite (Server Protections)").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) antispam(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Spam (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiSpam)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.ServerProtections.AntiSpam = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Spam (Server Protections)").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.ServerProtections.AntiSpam = !cfg.Protections.ServerProtections.AntiSpam
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiSpam) + " Anti Spam (Server Protections)").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) antiupper(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Upper (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiUpper)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.ServerProtections.AntiUpper = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Upper (Server Protections)").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.ServerProtections.AntiUpper = !cfg.Protections.ServerProtections.AntiUpper
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiUpper) + " Anti Upper (Server Protections)").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) antidelete(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.New(ctx).
// 			AddField("Anti Delete (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiDelete)).
// 			Send()
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Protections.ServerProtections.AntiDelete = toggle
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Anti Delete (Server Protections)").
// 			Send()
// 	case "toggle":
// 		cfg.Protections.ServerProtections.AntiDelete = !cfg.Protections.ServerProtections.AntiDelete
// 		_ = c.services.Config.SaveConfig(cfg)

// 		c.services.Embed.New(ctx).
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiDelete) + " Anti Delete (Server Protections)").
// 			Send()
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
// 	}
// }

// func (c *ProtectionsCommands) dm(ctx *router.Ctx) {
// 	var err error
// 	var action string

// 	switch strings.ToLower(utils.Arguments["<open/close>"]) {
// 	case "open":
// 		ctx.Session.State.UserSettings.Privacy.RestrictedGuildIds = []uint64{}
// 		err = ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)

// 		action = "Opened"
// 	case "close":
// 		var guildIDs []uint64
// 		for _, g := range ctx.Session.State.Guilds {
// 			id, _ := strconv.Atoi(g.ID)

// 			guildIDs = append(guildIDs, uint64(id))
// 		}

// 		ctx.Session.State.UserSettings.Privacy.RestrictedGuildIds = guildIDs
// 		err = ctx.Session.UserSettingsUpdate(ctx.Session.State.UserSettings, true)

// 		action = "Closed"
// 	default:
// 		ctx.RespondText("Invalid <open/close> argument. Use `open` or `close`.")
// 		return
// 	}

// 	if err != nil {
// 		ctx.RespondText("Error changing settings: %v", err)
// 		return
// 	}

// 	c.services.Embed.New(ctx).
// 		SetDescription(action + " DM's from everyone!").
// 		Send()
// }

// func (c *ProtectionsCommands) friendrequests(ctx *router.Ctx) {
// 	var err error
// 	var action string

// 	switch utils.Arguments["<allow/deny>"] {
// 	case "allow":
// 		err = api.AllowFriendRequests(ctx.Session, true)
// 		action = "Allowed"
// 	case "deny":
// 		err = api.AllowFriendRequests(ctx.Session, false)
// 		action = "Denied"
// 	default:
// 		ctx.RespondText("Invalid <allow/deny> argument. Use `allow` or `deny`.")
// 		return
// 	}

// 	if err != nil {
// 		ctx.RespondText("Error changing settings: %v", err)
// 		return
// 	}

// 	c.services.Embed.New(ctx).
// 		SetDescription(action + " friend requests from everyone!").
// 		Send()
// }
