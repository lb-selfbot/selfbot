package commands

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type SharingCommands struct{}

func (s *SharingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Sharing", "Share the bot with other users")

	router.RegisterCmd(&dgc.Command{
		Name:        "sharedusers",
		Description: "Get, add, remove or clear shared users",
		Usage:       "[p]sharedusers <get/add/rm/clear> [user]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.sharedusers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sharedcategories",
		Description: "Get, add, remove or clear shared command categories",
		Usage:       "[p]sharedcategories <get/add/rm/clear/values> [category]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.sharedcategories,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "sharedcommands",
		Description: "Get, add, remove or clear shared commands",
		Usage:       "[p]sharedcommands <get/add/rm/clear/values> [command]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.sharedcommands,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "shareall",
		Description: "Get, set or toggle sharing of all commands",
		Usage:       "[p]shareall <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.shareall,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "blockdangerous",
		Description: "Get, set or toggle blocking sharing of dangerous commands",
		Usage:       "[p]blockdangerous <get/set/toggle> [on/off]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     s.blockdangerous,
	})
}

func (*SharingCommands) sharedusers(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch strings.ToLower(utils.Arguments["<get/add/rm/clear>"]) {
	case "get":
		userString := ""
		for _, userID := range cfg.Sharing.SharedUsers {
			user, ok := utils.Discord.GetUser(ctx.Session, userID)
			if !ok {
				userString += fmt.Sprintf("%s - %s\n", "Unknown", userID)
				continue
			}

			userString += fmt.Sprintf("%s - %s\n", user.Username, userID)
		}

		if userString == "" {
			userString = "No users set."
		}

		embed.NewEmbed().
			AddField("Shared Users", userString).
			Send(ctx)
	case "add":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user.").
				Send(ctx)
			return
		}

		if utils.SliceContains(cfg.Sharing.SharedUsers, user.ID) {
			_ = ctx.RespondText("User is already in the Shared Users list.")
			return
		}

		cfg.Sharing.SharedUsers = append(cfg.Sharing.SharedUsers, user.ID)

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Added %s to the Shared Users list.", user.Mention())).
			Send(ctx)
	case "rm":
		user, ok := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
		if !ok {
			embed.NewEmbed().
				SetDescription("Could not find user.").
				Send(ctx)
			return
		}

		if !utils.SliceContains(cfg.Sharing.SharedUsers, user.ID) {
			_ = ctx.RespondText("User already isn't in the Shared Users list.")
			return
		}

		cfg.Sharing.SharedUsers = utils.RemoveSliceItem(cfg.Sharing.SharedUsers, user.ID)

		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription(fmt.Sprintf("Removed %s from the Shared Users list.", user.Mention())).
			Send(ctx)
	case "clear":
		cfg.Sharing.SharedUsers = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all users from the Shared Users list.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
	}
}

func (*SharingCommands) sharedcategories(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("Shared Categories", strings.Join(cfg.Sharing.SharedCategories, ", ")).
			Send(ctx)
	case "add":
		entry := utils.Arguments["[category]"]

		var found bool
		for _, category := range ctx.Router.Categories {
			if !strings.EqualFold(category.Name, entry) {
				continue
			}

			entry = category.Name
			found = true
			break
		}

		if !found {
			_ = ctx.RespondText("Category not found. Use `[p]sharedcategories values` to get a list of available categories.")
			return
		}

		if utils.SliceContains(cfg.Sharing.SharedCategories, entry) {
			_ = ctx.RespondText("Category is already in the Shared Categories list!")
			return
		}

		cfg.Sharing.SharedCategories = append(cfg.Sharing.SharedCategories, entry)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Added " + entry + " to the Shared Categories list.").
			Send(ctx)
	case "rm":
		entry := utils.Arguments["[category]"]

		var found bool
		for _, category := range ctx.Router.Categories {
			if !strings.EqualFold(category.Name, entry) {
				continue
			}

			entry = category.Name
			found = true
			break
		}

		if !found {
			_ = ctx.RespondText("Category not found. Use `[p]sharedcategories values` to get a list of available categories.")
			return
		}

		if !utils.SliceContains(cfg.Sharing.SharedCategories, entry) {
			_ = ctx.RespondText("Category already isn't in the Shared Categories list!")
			return
		}

		cfg.Sharing.SharedCategories = utils.RemoveSliceItem(cfg.Sharing.SharedCategories, entry)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Removed " + entry + " from the Shared Categories list.").
			Send(ctx)
	case "clear":
		cfg.Sharing.SharedCategories = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all Shared Categories list entries.").
			Send(ctx)
	case "values":
		var categories []string

		for _, category := range ctx.Router.Categories {
			categories = append(categories, category.Name)
		}

		embed.NewEmbed().
			AddField("Available Command Categories", strings.Join(categories, ", ")).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear/values> argument. Use `get`, `add`, `rm`, `clear` or `values`.")
	}
}

func (*SharingCommands) sharedcommands(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config.")
		return
	}

	switch utils.Arguments["<get/add/rm/clear/values>"] {
	case "get":
		embed.NewEmbed().
			AddField("Shared Commands", strings.Join(cfg.Sharing.SharedCommands, ", ")).
			Send(ctx)
	case "add":
		entry := utils.Arguments["[command]"]

		var found bool
		for _, command := range ctx.Router.Commands {
			if !strings.EqualFold(entry, command.Name) {
				continue
			}

			entry = command.Name
			found = true
			break
		}

		if !found {
			_ = ctx.RespondText("Command not found. Use `[p]sharedcommands values` to get a list of available commands.")
			return
		}

		if utils.SliceContains(cfg.Sharing.SharedCommands, entry) {
			_ = ctx.RespondText("Command is already in the Shared Commands list!")
			return
		}

		cfg.Sharing.SharedCommands = append(cfg.Sharing.SharedCommands, entry)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Added " + entry + " to the Shared Commands list.").
			Send(ctx)
	case "rm":
		entry := utils.Arguments["[command]"]

		var found bool
		for _, command := range ctx.Router.Commands {
			if !strings.EqualFold(entry, command.Name) {
				continue
			}

			entry = command.Name
			found = true
			break
		}

		if !found {
			_ = ctx.RespondText("Command not found. Use `[p]sharedcommands values` to get a list of available commands.")
			return
		}

		if !utils.SliceContains(cfg.Sharing.SharedCommands, entry) {
			_ = ctx.RespondText("Command already isn't in the Shared Commands list!")
			return
		}

		cfg.Sharing.SharedCommands = utils.RemoveSliceItem(cfg.Sharing.SharedCommands, entry)
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Removed " + entry + " from the Shared Commands list.").
			Send(ctx)
	case "clear":
		cfg.Sharing.SharedCommands = []string{}
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Cleared all Shared Commands list entries.").
			Send(ctx)
	case "values":
		embed.NewEmbed().
			AddField("Available Commands", "Use `[p]help` or view <https://docs.lightning-bot.com/category/commands> for a list of commands.").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/add/rm/clear/values> argument. Use `get`, `add`, `rm`, `clear` or `values`.")
	}
}

func (*SharingCommands) shareall(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Share All Commands", utils.Misc.FormatBoolOnOff(cfg.Sharing.ShareAll)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Sharing.ShareAll = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " Share All Commands").
			Send(ctx)
	case "toggle":
		cfg.Sharing.ShareAll = !cfg.Sharing.ShareAll
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Sharing.ShareAll) + " Share All Commands").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}

func (*SharingCommands) blockdangerous(ctx *dgc.Ctx) {
	cfg, err := config.GetConfig()
	if err != nil {
		_ = ctx.RespondText("Error getting config")
		return
	}

	switch utils.Arguments["<get/set/toggle>"] {
	case "get":
		embed.NewEmbed().
			AddField("Sharing Block Dangerous", utils.Misc.FormatBoolOnOff(cfg.Sharing.BlockDangerousCommands)).
			Send(ctx)
	case "set":
		toggle := strings.ToLower(utils.Arguments["[on/off]"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}

		cfg.Sharing.BlockDangerousCommands = toggle == "on"
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + toggle + " Sharing Block Dangerous").
			Send(ctx)
	case "toggle":
		cfg.Sharing.BlockDangerousCommands = !cfg.Sharing.BlockDangerousCommands
		_ = cfg.Save()

		embed.NewEmbed().
			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Sharing.BlockDangerousCommands) + " Sharing Block Dangerous").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
	}
}
