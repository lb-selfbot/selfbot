package protections

import (
	"slices"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

type ProtectionsCommands struct{ services *services.Services }

func NewProtectionsCommands(s *services.Services) *ProtectionsCommands {
	return &ProtectionsCommands{services: s}
}
func (c *ProtectionsCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Protections",
		Description: "Protect your account",
	}, []*router.Command{
		{
			Name:        "antiphishing",
			Description: "Get, set or toggle Anti Phishing",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antiphishing,
		},
		{
			Name:        "antifriendrequest",
			Description: "Get, set or toggle Anti Friend Request Spam",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antifriendrequest,
		},
		{
			Name:        "antimassdm",
			Description: "Get, set or toggle Anti Mass DM / DM Spam",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antimassdm,
		},
		{
			Name:        "protectedservers",
			Description: "Get, add, remove or clear the protected servers list",
			Arguments:   "<get/add/rm/clear> [server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.protectedservers,
		},
		{
			Name:        "antiinvite",
			Description: "Get, set or toggle Anti Invite (Server Protections)",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antiinvite,
		},
		{
			Name:        "antispam",
			Description: "Get, set or toggle Anti Spam (Server Protections)",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antispam,
		},
		{
			Name:        "antiupper",
			Description: "Get, set or toggle Anti Upper (Server Protections)",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antiupper,
		},
		{
			Name:        "antidelete",
			Description: "Get, set or toggle Anti Delete (Server Protections)",
			Arguments:   "<get/set/toggle> [on/off]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.antidelete,
		},
		{
			Name:        "dm",
			Description: "Open or close DM's from everyone",
			Arguments:   "<open/close>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.dm,
		},
		{
			Name:        "friendrequests",
			Description: "Allow or deny new friend requests",
			Arguments:   "<allow/deny>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.friendrequests,
		},
	})

}

func (c *ProtectionsCommands) antiphishing(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Phishing", utils.Misc.FormatBoolOnOff(cfg.Protections.AntiPhishing)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.AntiPhishing = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Phishing").
			Send(ctx)
	case "toggle":
		cfg.Protections.AntiPhishing = !cfg.Protections.AntiPhishing
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.AntiPhishing) + " Anti Phishing").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) antifriendrequest(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Friend Request Spam", utils.Misc.FormatBoolOnOff(cfg.Protections.AntiFriendRequest)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.AntiFriendRequest = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Friend Request Spam").
			Send(ctx)
	case "toggle":
		cfg.Protections.AntiFriendRequest = !cfg.Protections.AntiFriendRequest
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.AntiFriendRequest) + " Anti Friend Request Spam").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) antimassdm(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Mass DM / DM Spam", utils.Misc.FormatBoolOnOff(cfg.Protections.AntiMassDM)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.AntiMassDM = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Mass DM / DM Spam").
			Send(ctx)
	case "toggle":
		cfg.Protections.AntiMassDM = !cfg.Protections.AntiMassDM
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.AntiMassDM) + " Anti Mass DM / DM Spam").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) protectedservers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear>"] {
	case "get":
		var servers []string

		for _, serverID := range cfg.Protections.ServerProtections.ProtectedServers {
			guild, err := utils.Discord.GetGuild(ctx.Session, serverID)
			if err != nil {
				servers = append(servers, serverID)
				continue
			}

			servers = append(servers, guild.Name+" ("+guild.ID+")")
		}

		c.services.Embed.NewEmbed().
			AddField("Protected Servers List", strings.Join(servers, ", ")).
			Send(ctx)
	case "add":
		guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
		if err != nil {
			ctx.RespondText("Could not find server.")
			return
		}

		if slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
			ctx.RespondText("Server %s is already in the Protected Servers List!", guild.Name)
			return
		}

		cfg.Protections.ServerProtections.ProtectedServers = append(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Added %s to the Protected Servers List.", guild.Name).
			Send(ctx)
	case "rm":
		guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
		if err != nil {
			ctx.RespondText("Could not find server.")
			return
		}

		if !slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
			ctx.RespondText("Server %s already isn't in the Protected Servers List!", guild.Name)
			return
		}

		cfg.Protections.ServerProtections.ProtectedServers = utils.RemoveSliceItem(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Removed %s from the Protected Servers List.", guild.Name).
			Send(ctx)
	case "clear":
		cfg.Protections.ServerProtections.ProtectedServers = []string{}
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all Protected Servers List entries.").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (c *ProtectionsCommands) antiinvite(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Invite (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiInvite)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.ServerProtections.AntiInvite = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Invite (Server Protections)").
			Send(ctx)
	case "toggle":
		cfg.Protections.ServerProtections.AntiInvite = !cfg.Protections.ServerProtections.AntiInvite
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiInvite) + " Anti Invite (Server Protections)").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) antispam(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Spam (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiSpam)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.ServerProtections.AntiSpam = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Spam (Server Protections)").
			Send(ctx)
	case "toggle":
		cfg.Protections.ServerProtections.AntiSpam = !cfg.Protections.ServerProtections.AntiSpam
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiSpam) + " Anti Spam (Server Protections)").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) antiupper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Upper (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiUpper)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.ServerProtections.AntiUpper = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Upper (Server Protections)").
			Send(ctx)
	case "toggle":
		cfg.Protections.ServerProtections.AntiUpper = !cfg.Protections.ServerProtections.AntiUpper
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiUpper) + " Anti Upper (Server Protections)").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) antidelete(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		c.services.Embed.NewEmbed().
			AddField("Anti Delete (Server Protections)", utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiDelete)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Protections.ServerProtections.AntiDelete = toggle == "on"
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + toggle + " Anti Delete (Server Protections)").
			Send(ctx)
	case "toggle":
		cfg.Protections.ServerProtections.AntiDelete = !cfg.Protections.ServerProtections.AntiDelete
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Protections.ServerProtections.AntiDelete) + " Anti Delete (Server Protections)").
			Send(ctx)
	default:
		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (c *ProtectionsCommands) dm(ctx *router.Ctx) {
	var err error
	var action string

	switch strings.ToLower(utils.Arguments["<open/close>"]) {
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
		ctx.RespondText("Error changing settings: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(action + " DM's from everyone!").
		Send(ctx)
}

func (c *ProtectionsCommands) friendrequests(ctx *router.Ctx) {
	var err error
	var action string

	switch utils.Arguments["<allow/deny>"] {
	case "allow":
		err = api.AllowFriendRequests(ctx.Session, true)
		action = "Allowed"
	case "deny":
		err = api.AllowFriendRequests(ctx.Session, false)
		action = "Denied"
	default:
		ctx.RespondText("Invalid <allow/deny> argument. Use `allow` or `deny`.")
		return
	}

	if err != nil {
		ctx.RespondText("Error changing settings: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(action + " friend requests from everyone!").
		Send(ctx)
}
