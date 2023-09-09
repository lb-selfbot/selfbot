package commands

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type SharingCommands struct{}

func (s *SharingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Sharing", "Share the bot with other users")
	router.RegisterCmd(&dgc.Command{
		Name:        "addshareduser",
		Description: "Add a user to the shared users",
		Usage:       "[p]addshareduser <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.addshareduser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removeshareduser",
		Description: "Remove a user from the shared users",
		Usage:       "[p]removeshareduser <user>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.removeshareduser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "listsharedusers",
		Description: "See all shared users",
		Usage:       "[p]listsharedusers",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.listsharedusers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addsharedcategory",
		Description: "Add a command category to be shared",
		Usage:       "[p]addsharedcategory <category name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.addsharedcategory,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removesharedcategory",
		Description: "Remove a shared category",
		Usage:       "[p]removesharedcategory <category name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.removesharedcategory,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "listsharedcategories",
		Description: "See all shared categories",
		Usage:       "[p]listsharedcategories",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.listsharedcategories,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addsharedcommand",
		Description: "Add a command to be shared",
		Usage:       "[p]addsharedcommand <command name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.addsharedcommand,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removesharedcommand",
		Description: "Remove a shared command",
		Usage:       "[p]removesharedcommand <command name>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.removesharedcommand,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "listsharedcommands",
		Description: "See all shared commands",
		Usage:       "[p]listsharedcommands",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.listsharedcommands,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shareall",
		Description: "Share all commands",
		Usage:       "[p]shareall <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.shareall,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blockdangerous",
		Description: "Block users from using dangerous commands",
		Usage:       "[p]blockdangerous <on/off>",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.blockdangerous,
	})
}

func (*SharingCommands) addshareduser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	cfg, _ := config.GetConfig()
	if utils.SliceContains(cfg.Sharing.SharedUsers, user.ID) {
		_ = ctx.RespondText("User is already shared.")
		return
	}
	cfg.Sharing.SharedUsers = append(cfg.Sharing.SharedUsers, user.ID)
	_ = cfg.Save()
}

func (*SharingCommands) removeshareduser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	cfg, _ := config.GetConfig()
	if !utils.SliceContains(cfg.Sharing.SharedUsers, user.ID) {
		_ = ctx.RespondText("User is not shared.")
		return
	}
	cfg.Sharing.SharedUsers = utils.RemoveSliceItem(cfg.Sharing.SharedUsers, user.ID)
	_ = cfg.Save()
}

func (*SharingCommands) listsharedusers(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	userString := ""
	for _, userID := range cfg.Sharing.SharedUsers {
		user, _ := ctx.Session.State.GetUser(userID)
		if user == nil {
			user, _ = ctx.Session.User(userID)
		}
		userString += fmt.Sprintf("%s - %s\n", user.Username, user.ID)
	}
	embed.NewEmbed().
		AddField("Shared Users", userString).
		Send(ctx)
}

func (*SharingCommands) addsharedcategory(ctx *dgc.Ctx) {
	category := ctx.Router.GetCategory(utils.Arguments["<category name>"])
	if category == nil {
		_ = ctx.RespondText("Could not find category.")
		return
	}
	cfg, _ := config.GetConfig()
	if utils.SliceContains(cfg.Sharing.SharedCategories, category.Name) {
		_ = ctx.RespondText("Category is already shared.")
		return
	}
	cfg.Sharing.SharedCategories = append(cfg.Sharing.SharedCategories, category.Name)
	_ = cfg.Save()
}

func (*SharingCommands) removesharedcategory(ctx *dgc.Ctx) {
	category := ctx.Router.GetCategory(utils.Arguments["<category name>"])
	if category == nil {
		_ = ctx.RespondText("Could not find category.")
		return
	}
	cfg, _ := config.GetConfig()
	if !utils.SliceContains(cfg.Sharing.SharedCategories, category.Name) {
		_ = ctx.RespondText("Category is not shared.")
		return
	}
	cfg.Sharing.SharedCategories = utils.RemoveSliceItem(cfg.Sharing.SharedCategories, category.Name)
	_ = cfg.Save()
}

func (*SharingCommands) listsharedcategories(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	embed.NewEmbed().
		AddField("Shared Categories", strings.Join(cfg.Sharing.SharedCategories, "\n")).
		Send(ctx)
}

func (*SharingCommands) addsharedcommand(ctx *dgc.Ctx) {
	command := ctx.Router.GetCmd(utils.Arguments["<command name>"])
	if command == nil {
		_ = ctx.RespondText("Could not find command.")
		return
	}
	cfg, _ := config.GetConfig()
	if utils.SliceContains(cfg.Sharing.SharedCommands, command.Name) {
		_ = ctx.RespondText("Command is already shared.")
		return
	}
	cfg.Sharing.SharedCommands = append(cfg.Sharing.SharedCommands, command.Name)
	_ = cfg.Save()
}

func (*SharingCommands) removesharedcommand(ctx *dgc.Ctx) {
	command := ctx.Router.GetCmd(utils.Arguments["<command name>"])
	if command == nil {
		_ = ctx.RespondText("Could not find command.")
		return
	}
	cfg, _ := config.GetConfig()
	if !utils.SliceContains(cfg.Sharing.SharedCommands, command.Name) {
		_ = ctx.RespondText("Command is not shared.")
		return
	}
	cfg.Sharing.SharedCommands = utils.RemoveSliceItem(cfg.Sharing.SharedCommands, command.Name)
	_ = cfg.Save()
}

func (*SharingCommands) listsharedcommands(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	embed.NewEmbed().
		AddField("Shared Commands", strings.Join(cfg.Sharing.SharedCommands, "\n")).
		Send(ctx)
}

func (*SharingCommands) shareall(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Sharing.ShareAll = toggle == "on"
	_ = cfg.Save()
}

func (*SharingCommands) blockdangerous(ctx *dgc.Ctx) {
	cfg, _ := config.GetConfig()
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}
	cfg.Sharing.BlockDangerousCommands = toggle == "on"
	_ = cfg.Save()
}
