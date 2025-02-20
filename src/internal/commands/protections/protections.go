package protections

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/dgc"
)

func (p ProtectionsCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Protections", "Protection commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "antiphishing",
		Description: "Get, set or toggle Anti Phishing",
		Usage:       "[p]antiphishing <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antiphishing,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "antifriendrequest",
		Description: "Get, set or toggle Anti Friend Request Spam",
		Usage:       "[p]antifriendrequest <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antifriendrequest,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "antimassdm",
		Description: "Get, set or toggle Anti Mass DM / DM Spam",
		Usage:       "[p]antimassdm <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antimassdm,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "protectedservers",
		Description: "Get, add, remove or clear the protected servers list",
		Usage:       "[p]protectedservers <get/add/rm/clear> [server id]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.protectedservers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "antiinvite",
		Description: "Get, set or toggle Anti Invite (Server Protections)",
		Usage:       "[p]antiinvite <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antiinvite,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "antispam",
		Description: "Get, set or toggle Anti Spam (Server Protections)",
		Usage:       "[p]antispam <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antispam,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "antiupper",
		Description: "Get, set or toggle Anti Upper (Server Protections)",
		Usage:       "[p]antiupper <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antiupper,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "antidelete",
		Description: "Get, set or toggle Anti Delete (Server Protections)",
		Usage:       "[p]antidelete <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.antidelete,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "dm",
		Description: "Open or close DM's from everyone",
		Usage:       "[p]dm <open/close>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.dm,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "friendrequests",
		Description: "Allow or deny new friend requests",
		Usage:       "[p]friendrequests <allow/deny>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     p.friendrequests,
	})
}

func (p ProtectionsCommands) antiphishing(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (p ProtectionsCommands) antifriendrequest(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (p ProtectionsCommands) antimassdm(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ProtectionsCommands) protectedservers(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
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
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
		if err != nil {
			_ = ctx.RespondText("Could not find server.")
			return
		}

		if slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
			_ = ctx.RespondText(fmt.Sprintf("Server %s is already in the Protected Servers List!", guild.Name))
			return
		}

		cfg.Protections.ServerProtections.ProtectedServers = append(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added %s to the Protected Servers List.", guild.Name)).
			Send(ctx)
	case "rm":
		guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
		if err != nil {
			_ = ctx.RespondText("Could not find server.")
			return
		}

		if !slices.Contains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
			_ = ctx.RespondText(fmt.Sprintf("Server %s already isn't in the Protected Servers List!", guild.Name))
			return
		}

		cfg.Protections.ServerProtections.ProtectedServers = utils.RemoveSliceItem(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed %s from the Protected Servers List.", guild.Name)).
			Send(ctx)
	case "clear":
		cfg.Protections.ServerProtections.ProtectedServers = []string{}
		_ = cfg.Save()

		c.services.Embed.NewEmbed().
			SetDescription("Cleared all Protected Servers List entries.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (p ProtectionsCommands) antiinvite(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (p ProtectionsCommands) antispam(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (p ProtectionsCommands) antiupper(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (p ProtectionsCommands) antidelete(ctx *router.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
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
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
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
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get` or `set` or `toggle`.")
	}
}

func (ProtectionsCommands) dm(ctx *router.Ctx) {
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
		_ = ctx.RespondText("Invalid <open/close> argument. Use `open` or `close`.")
		return
	}

	if err != nil {
		_ = ctx.RespondText("Error changing settings: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(action + " DM's from everyone!").
		Send(ctx)
}

func (ProtectionsCommands) friendrequests(ctx *router.Ctx) {
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
		_ = ctx.RespondText("Invalid <allow/deny> argument. Use `allow` or `deny`.")
		return
	}

	if err != nil {
		_ = ctx.RespondText("Error changing settings: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetDescription(action + " friend requests from everyone!").
		Send(ctx)
}

type ProtectionsCommands struct{ services *services.Services }

func NewProtectionsCommands(s *services.Services) *ProtectionsCommands {
	return &ProtectionsCommands{services: s}
}
func (c *ProtectionsCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

}
