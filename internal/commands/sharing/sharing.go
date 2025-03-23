package sharing

import (
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
)

type SharingCommands struct{ services *services.Services }

func NewSharingCommands(s *services.Services) *SharingCommands {
	return &SharingCommands{services: s}
}

func (c *SharingCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Sharing",
		Description: "Share the bot with other users",
	}, []*router.Command{
		// {
		// 	Name:        "sharedusers",
		// 	Description: "Get, add, remove or clear shared users",
		// 	Arguments:   "<get/add/rm/clear> [user]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.sharedusers,
		// },
		// {
		// 	Name:        "sharedcategories",
		// 	Description: "Get, add, remove or clear shared command categories",
		// 	Arguments:   "<get/add/rm/clear/values> [category]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.sharedcategories,
		// },
		// {
		// 	Name:        "sharedcommands",
		// 	Description: "Get, add, remove or clear shared commands",
		// 	Arguments:   "<get/add/rm/clear/values> [command]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.sharedcommands,
		// },
		// {
		// 	Name:        "shareall",
		// 	Description: "Get, set or toggle sharing of all commands",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.shareall,
		// },
		// {
		// 	Name:        "blockdangerous",
		// 	Description: "Get, set or toggle blocking sharing of dangerous commands",
		// 	Arguments:   "<get/set/toggle> [on/off]",
		// 	Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		// 	Handler:     c.blockdangerous,
		// },
	})
}

// func (c *SharingCommands) sharedusers(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch strings.ToLower(utils.Arguments["<get/add/rm/clear>"]) {
// 	case "get":
// 		userString := ""
// 		for _, userID := range cfg.Sharing.SharedUsers {
// 			user, err := utils.Discord.GetUser(ctx.Session, userID)
// 			if err != nil {
// 				userString += fmt.Sprintf("%s - %s\n", "Unknown", userID)
// 				continue
// 			}

// 			userString += fmt.Sprintf("%s - %s\n", user.Username, userID)
// 		}

// 		if userString == "" {
// 			userString = "No users set."
// 		}

// 		c.services.Embed.NewEmbed().
// 			AddField("Shared Users", userString).
// 			Send(ctx)
// 	case "add":
// 		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
// 		if err != nil {
// 			c.services.Embed.NewEmbed().
// 				SetDescription("Could not find user.").
// 				Send(ctx)
// 			return
// 		}

// 		if slices.Contains(cfg.Sharing.SharedUsers, user.ID) {
// 			ctx.RespondText("User is already in the Shared Users list.")
// 			return
// 		}

// 		cfg.Sharing.SharedUsers = append(cfg.Sharing.SharedUsers, user.ID)

// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Added %s to the Shared Users list.", user.Mention()).
// 			Send(ctx)
// 	case "rm":
// 		user, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
// 		if err != nil {
// 			c.services.Embed.NewEmbed().
// 				SetDescription("Could not find user.").
// 				Send(ctx)
// 			return
// 		}

// 		if !slices.Contains(cfg.Sharing.SharedUsers, user.ID) {
// 			ctx.RespondText("User already isn't in the Shared Users list.")
// 			return
// 		}

// 		cfg.Sharing.SharedUsers = utils.RemoveSliceItem(cfg.Sharing.SharedUsers, user.ID)

// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Removed %s from the Shared Users list.", user.Mention()).
// 			Send(ctx)
// 	case "clear":
// 		cfg.Sharing.SharedUsers = []string{}
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Cleared all users from the Shared Users list.").
// 			Send(ctx)
// 	default:
// 		ctx.RespondText("Invalid <get/add/rm/clear> argument. Use `get`, `add`, `rm`, or `clear`.")
// 	}
// }

// func (c *SharingCommands) sharedcategories(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config.")
// 		return
// 	}

// 	switch utils.Arguments["<get/add/rm/clear/values>"] {
// 	case "get":
// 		c.services.Embed.NewEmbed().
// 			AddField("Shared Categories", strings.Join(cfg.Sharing.SharedCategories, ", ")).
// 			Send(ctx)
// 	case "add":
// 		entry := utils.Arguments["[category]"]

// 		var found bool
// 		for _, category := range ctx.Router.GetCategories() {
// 			if !strings.EqualFold(category.Name, entry) {
// 				continue
// 			}

// 			entry = category.Name
// 			found = true
// 			break
// 		}

// 		if !found {
// 			ctx.RespondText("Category not found. Use `[p]sharedcategories values` to get a list of available categories.")
// 			return
// 		}

// 		if slices.Contains(cfg.Sharing.SharedCategories, entry) {
// 			ctx.RespondText("Category is already in the Shared Categories list!")
// 			return
// 		}

// 		cfg.Sharing.SharedCategories = append(cfg.Sharing.SharedCategories, entry)
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Added " + entry + " to the Shared Categories list.").
// 			Send(ctx)
// 	case "rm":
// 		entry := utils.Arguments["[category]"]

// 		var found bool
// 		for _, category := range ctx.Router.GetCategories() {
// 			if !strings.EqualFold(category.Name, entry) {
// 				continue
// 			}

// 			entry = category.Name
// 			found = true
// 			break
// 		}

// 		if !found {
// 			ctx.RespondText("Category not found. Use `[p]sharedcategories values` to get a list of available categories.")
// 			return
// 		}

// 		if !slices.Contains(cfg.Sharing.SharedCategories, entry) {
// 			ctx.RespondText("Category already isn't in the Shared Categories list!")
// 			return
// 		}

// 		cfg.Sharing.SharedCategories = utils.RemoveSliceItem(cfg.Sharing.SharedCategories, entry)
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Removed " + entry + " from the Shared Categories list.").
// 			Send(ctx)
// 	case "clear":
// 		cfg.Sharing.SharedCategories = []string{}
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Cleared all Shared Categories list entries.").
// 			Send(ctx)
// 	case "values":
// 		var categories []string

// 		for _, category := range ctx.Router.GetCategories() {
// 			categories = append(categories, category.Name)
// 		}

// 		c.services.Embed.NewEmbed().
// 			AddField("Available Command Categories", strings.Join(categories, ", ")).
// 			Send(ctx)
// 	default:
// 		ctx.RespondText("Invalid <get/add/rm/clear/values> argument. Use `get`, `add`, `rm`, `clear` or `values`.")
// 	}
// }

// func (c *SharingCommands) sharedcommands(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config.")
// 		return
// 	}

// 	switch utils.Arguments["<get/add/rm/clear/values>"] {
// 	case "get":
// 		c.services.Embed.NewEmbed().
// 			AddField("Shared Commands", strings.Join(cfg.Sharing.SharedCommands, ", ")).
// 			Send(ctx)
// 	case "add":
// 		entry := utils.Arguments["[command]"]

// 		var found bool
// 		for _, command := range ctx.Router.GetCmds() {
// 			if !strings.EqualFold(entry, command.Name) {
// 				continue
// 			}

// 			entry = command.Name
// 			found = true
// 			break
// 		}

// 		if !found {
// 			ctx.RespondText("Command not found. Use `[p]sharedcommands values` to get a list of available commands.")
// 			return
// 		}

// 		if slices.Contains(cfg.Sharing.SharedCommands, entry) {
// 			ctx.RespondText("Command is already in the Shared Commands list!")
// 			return
// 		}

// 		cfg.Sharing.SharedCommands = append(cfg.Sharing.SharedCommands, entry)
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Added " + entry + " to the Shared Commands list.").
// 			Send(ctx)
// 	case "rm":
// 		entry := utils.Arguments["[command]"]

// 		var found bool
// 		for _, command := range ctx.Router.GetCmds() {
// 			if !strings.EqualFold(entry, command.Name) {
// 				continue
// 			}

// 			entry = command.Name
// 			found = true
// 			break
// 		}

// 		if !found {
// 			ctx.RespondText("Command not found. Use `[p]sharedcommands values` to get a list of available commands.")
// 			return
// 		}

// 		if !slices.Contains(cfg.Sharing.SharedCommands, entry) {
// 			ctx.RespondText("Command already isn't in the Shared Commands list!")
// 			return
// 		}

// 		cfg.Sharing.SharedCommands = utils.RemoveSliceItem(cfg.Sharing.SharedCommands, entry)
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Removed " + entry + " from the Shared Commands list.").
// 			Send(ctx)
// 	case "clear":
// 		cfg.Sharing.SharedCommands = []string{}
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Cleared all Shared Commands list entries.").
// 			Send(ctx)
// 	case "values":
// 		c.services.Embed.NewEmbed().
// 			AddField("Available Commands", "Use `[p]help` or view <https://docs.lightning-bot.com/category/commands> for a list of commands.").
// 			Send(ctx)
// 	default:
// 		ctx.RespondText("Invalid <get/add/rm/clear/values> argument. Use `get`, `add`, `rm`, `clear` or `values`.")
// 	}
// }

// func (c *SharingCommands) shareall(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.NewEmbed().
// 			AddField("Share All Commands", utils.Misc.FormatBoolOnOff(cfg.Sharing.ShareAll)).
// 			Send(ctx)
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Sharing.ShareAll = toggle
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Share All Commands").
// 			Send(ctx)
// 	case "toggle":
// 		cfg.Sharing.ShareAll = !cfg.Sharing.ShareAll
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Sharing.ShareAll) + " Share All Commands").
// 			Send(ctx)
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
// 	}
// }

// func (c *SharingCommands) blockdangerous(ctx *router.Ctx) {
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		ctx.RespondText("Error getting config")
// 		return
// 	}

// 	switch utils.Arguments["<get/set/toggle>"] {
// 	case "get":
// 		c.services.Embed.NewEmbed().
// 			AddField("Sharing Block Dangerous", utils.Misc.FormatBoolOnOff(cfg.Sharing.BlockDangerousCommands)).
// 			Send(ctx)
// 	case "set":
// 		toggle, err := ctx.Arguments.Get(0).Bool()
// 		if err != nil {
// 			ctx.RespondText("Invalid value. Use `on` or `off`.")
// 			return
// 		}

// 		cfg.Sharing.BlockDangerousCommands = toggle
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(toggle) + " Sharing Block Dangerous").
// 			Send(ctx)
// 	case "toggle":
// 		cfg.Sharing.BlockDangerousCommands = !cfg.Sharing.BlockDangerousCommands
// 		_ = cfg.Save()

// 		c.services.Embed.NewEmbed().
// 			SetDescription("Turned " + utils.Misc.FormatBoolOnOff(cfg.Sharing.BlockDangerousCommands) + " Sharing Block Dangerous").
// 			Send(ctx)
// 	default:
// 		ctx.RespondText("Invalid <get/set/toggle> argument. Use `get`, `set` or `toggle`.")
// 	}
// }
