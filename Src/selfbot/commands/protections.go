package commands

import (
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type ProtectionsCommands struct{}

func (p ProtectionsCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Protections", "Protection commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "antiphishing",
		Description: "Turn anti phishing on or off",
		Usage:       "[p]antiphishing <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antiphishing,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antifriendrequest",
		Description: "Turn anti friend request spam on or off",
		Usage:       "[p]antifriendrequest <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antifriendrequest,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antimassdm",
		Description: "Turn anti mass DM on or off",
		Usage:       "[p]antimassmd <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antimassdm,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "protectedservers",
		Description: "Get a list of protected servers",
		Usage:       "[p]protectedservers",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.protectedservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "protectserver",
		Description: "Protect a server",
		Usage:       "[p]protectserver <server id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.protectserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unprotectserver",
		Description: "Protect a server",
		Usage:       "[p]protectserver <server id>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.unprotectserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antiinvite",
		Description: "Turn anti invite on or off",
		Usage:       "[p]antiinvite <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antiinvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antispam",
		Description: "Turn anti spam on or off",
		Usage:       "[p]antispam <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antispam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antiupper",
		Description: "Turn anti upper on or off",
		Usage:       "[p]antiupper <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antiupper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "antidelete",
		Description: "Turn anti delete on or off",
		Usage:       "[p]antidelete <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.antidelete,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "closedm",
		Description: "Close your DMs",
		Usage:       "[p]closedm",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.closedm,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "opendm",
		Description: "Allow DMs from everyone",
		Usage:       "[p]opendm",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.opendm,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "enablefriendrequests",
		Description: "Enable friend requests",
		Usage:       "[p]enablefriendrequests",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.enablefriendrequests,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "disablefriendrequests",
		Description: "Disable friend requests",
		Usage:       "[p]disablefriendrequests",
		Flags:       []string{"DANGEROUS"},
		Handler:     p.disablefriendrequests,
	})
}

func (p ProtectionsCommands) antiphishing(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.AntiPhishing = toggle == "on"
	_ = cfg.Save()
}

func (p ProtectionsCommands) antifriendrequest(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.AntiFriendRequest = toggle == "on"
	_ = cfg.Save()
}

func (p ProtectionsCommands) antimassdm(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.AntiMassDM = toggle == "on"
	_ = cfg.Save()
}

func (ProtectionsCommands) protectedservers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Could not get config.")
		return
	}

	var servers []string

	for _, s := range cfg.Protections.ServerProtections.ProtectedServers {
		if name := utils.Discord.GetGuildName(ctx.Session, s); name != "No Server" {
			servers = append(servers, "- "+name+" ("+s+")")
		} else {
			servers = append(servers, "- "+s)
		}
	}

	embed.NewEmbed().
		AddField("Protected Servers", strings.Join(servers, "\n")).
		Send(ctx)
}

func (p ProtectionsCommands) protectserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Could not find server.")
		return
	}
	cfg, _ := config.GetConfig()
	if utils.SliceContains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
		_ = ctx.RespondText("Server is already ignored.")
		return
	}
	cfg.Protections.ServerProtections.ProtectedServers = append(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
	_ = cfg.Save()
}

func (p ProtectionsCommands) unprotectserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["<server id>"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Could not find server.")
		return
	}
	cfg, _ := config.GetConfig()
	if !utils.SliceContains(cfg.Protections.ServerProtections.ProtectedServers, guild.ID) {
		_ = ctx.RespondText("Server is already not ignored.")
		return
	}
	cfg.Protections.ServerProtections.ProtectedServers = utils.RemoveSliceItem(cfg.Protections.ServerProtections.ProtectedServers, guild.ID)
	_ = cfg.Save()
}

func (p ProtectionsCommands) antiinvite(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.ServerProtections.AntiInvite = toggle == "on"
	_ = cfg.Save()
}

func (p ProtectionsCommands) antispam(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.ServerProtections.AntiSpam = toggle == "on"
	_ = cfg.Save()
}

func (p ProtectionsCommands) antiupper(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.ServerProtections.AntiUpper = toggle == "on"
	_ = cfg.Save()
}

func (p ProtectionsCommands) antidelete(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg, _ := config.GetConfig()
	cfg.Protections.ServerProtections.AntiDelete = toggle == "on"
	_ = cfg.Save()
}

func (ProtectionsCommands) opendm(ctx *dgc.Ctx) {
	_ = api.ChangeSettings(ctx.Session, "/settings", map[string]any{"restricted_guilds": []any{}})
}

func (ProtectionsCommands) closedm(ctx *dgc.Ctx) {
	var guildIDs []string
	for _, g := range ctx.Session.State.Guilds {
		guildIDs = append(guildIDs, g.ID)
	}

	_ = api.ChangeSettings(ctx.Session, "/settings", map[string]any{"restricted_guilds": guildIDs})
}

func (ProtectionsCommands) enablefriendrequests(ctx *dgc.Ctx) {
	_ = api.AllowFriendRequestsFrom(ctx.Session, []string{"all"})
}

func (ProtectionsCommands) disablefriendrequests(ctx *dgc.Ctx) {
	_ = api.AllowFriendRequestsFrom(ctx.Session, []string{})
}
