package sharing

import (
	"fmt"
	"slices"
	"strings"

	"github.com/lb-selfbot/selfbot/internal/actioncommands"
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

type SharingCommands struct{ services *services.Services }

func NewSharingCommands(s *services.Services) *SharingCommands {
	return &SharingCommands{services: s}
}

func (c *SharingCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Sharing",
		Description: "Share the bot with other users",
	}, []*router.Command{
		// Shared Users
		actioncommands.NewBuilder("sharedusers", "Shared Users").
			WithListGetter("Get shared users", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Sharing.SharedUsers) == 0 {
					c.services.Embed.Respond(ctx, "No shared users found.")
					return nil
				}

				userStrings := make([]string, 0, len(cfg.Sharing.SharedUsers))
				for _, userID := range cfg.Sharing.SharedUsers {
					user, err := ctx.Session.UserWithFallback(userID)
					if err != nil {
						userStrings = append(userStrings, fmt.Sprintf("%s (Unknown)", userID))
						continue
					}
					userStrings = append(userStrings, fmt.Sprintf("%s (%s)", user.Username, userID))
				}

				c.services.Embed.RespondField(ctx, "Shared Users", strings.Join(userStrings, "\n"))
				return nil
			}).
			WithListAdder("Add shared user", "<user>", func(ctx *router.Ctx) error {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find user: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Sharing.SharedUsers, user.ID) {
					return fmt.Errorf("user %s is already in the shared users list", user.Username)
				}

				cfg.Sharing.SharedUsers = append(cfg.Sharing.SharedUsers, user.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added user %s to shared users.", user.Username)
				return nil
			}).
			WithListRemover("Remove shared user", "<user>", func(ctx *router.Ctx) error {
				user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find user: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Sharing.SharedUsers, user.ID) {
					return fmt.Errorf("user %s is not in the shared users list", user.Username)
				}

				cfg.Sharing.SharedUsers = fn.Delete(cfg.Sharing.SharedUsers, user.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed user %s from shared users.", user.Username)
				return nil
			}).
			WithListClearer("Clear shared users", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.SharedUsers = []string{}
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all shared users.")
				return nil
			}).
			Build(c.services),

		// Shared Categories
		actioncommands.NewBuilder("sharedcategories", "Shared Categories").
			WithListGetter("Get shared categories", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Sharing.SharedCategories) == 0 {
					c.services.Embed.Respond(ctx, "No shared command categories found.")
					return nil
				}

				c.services.Embed.RespondField(ctx, "Shared Categories", strings.Join(cfg.Sharing.SharedCategories, "\n"))
				return nil
			}).
			WithListAdder("Add shared category", "<category>", func(ctx *router.Ctx) error {
				categoryName := ctx.Arguments.Get(0).Raw()

				var found bool
				var actualName string
				for _, category := range ctx.Router.GetCategories() {
					if !strings.EqualFold(category.Name, categoryName) {
						continue
					}

					actualName = category.Name
					found = true
					break
				}

				if !found {
					return fmt.Errorf("category %s not found. Use the 'values' subcommand to see available categories", categoryName)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Sharing.SharedCategories, actualName) {
					return fmt.Errorf("category %s is already in the shared categories list", actualName)
				}

				cfg.Sharing.SharedCategories = append(cfg.Sharing.SharedCategories, actualName)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added category %s to shared categories.", actualName)
				return nil
			}).
			WithListRemover("Remove shared category", "<category>", func(ctx *router.Ctx) error {
				categoryName := ctx.Arguments.Get(0).Raw()

				var found bool
				var actualName string
				for _, category := range ctx.Router.GetCategories() {
					if !strings.EqualFold(category.Name, categoryName) {
						continue
					}

					actualName = category.Name
					found = true
					break
				}

				if !found {
					return fmt.Errorf("category %s not found. Use the 'values' subcommand to see available categories", categoryName)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Sharing.SharedCategories, actualName) {
					return fmt.Errorf("category %s is not in the shared categories list", actualName)
				}

				cfg.Sharing.SharedCategories = fn.Delete(cfg.Sharing.SharedCategories, actualName)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed category %s from shared categories.", actualName)
				return nil
			}).
			WithListClearer("Clear shared categories", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.SharedCategories = []string{}
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all shared categories.")
				return nil
			}).
			Build(c.services),

		// Shared Commands
		actioncommands.NewBuilder("sharedcommands", "Shared Commands").
			WithListGetter("Get shared commands", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Sharing.SharedCommands) == 0 {
					c.services.Embed.Respond(ctx, "No shared commands found.")
					return nil
				}

				c.services.Embed.RespondField(ctx, "Shared Commands", strings.Join(cfg.Sharing.SharedCommands, "\n"))
				return nil
			}).
			WithListAdder("Add shared command", "<command>", func(ctx *router.Ctx) error {
				commandName := ctx.Arguments.Get(0).Raw()

				var found bool
				var actualName string
				for _, command := range ctx.Router.GetCommands() {
					if !strings.EqualFold(command.Name, commandName) {
						continue
					}

					actualName = command.Name
					found = true
					break
				}

				if !found {
					return fmt.Errorf("command %s not found", commandName)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Sharing.SharedCommands, actualName) {
					return fmt.Errorf("command %s is already in the shared commands list", actualName)
				}

				cfg.Sharing.SharedCommands = append(cfg.Sharing.SharedCommands, actualName)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added command %s to shared commands.", actualName)
				return nil
			}).
			WithListRemover("Remove shared command", "<command>", func(ctx *router.Ctx) error {
				commandName := ctx.Arguments.Get(0).Raw()

				var found bool
				var actualName string
				for _, command := range ctx.Router.GetCommands() {
					if !strings.EqualFold(command.Name, commandName) {
						continue
					}

					actualName = command.Name
					found = true
					break
				}

				if !found {
					return fmt.Errorf("command %s not found", commandName)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Sharing.SharedCommands, actualName) {
					return fmt.Errorf("command %s is not in the shared commands list", actualName)
				}

				cfg.Sharing.SharedCommands = fn.Delete(cfg.Sharing.SharedCommands, actualName)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed command %s from shared commands.", actualName)
				return nil
			}).
			WithListClearer("Clear shared commands", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.SharedCommands = []string{}
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all shared commands.")
				return nil
			}).
			Build(c.services),

		// Share All
		actioncommands.NewBuilder("shareall", "Share All").
			WithValueGetter("Get share all status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Sharing.ShareAll, nil
			}).
			WithValueSetter("Set share all status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.ShareAll = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle share all", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.ShareAll = !cfg.Sharing.ShareAll
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Block Dangerous
		actioncommands.NewBuilder("blockdangerous", "Block Dangerous Commands").
			WithValueGetter("Get block dangerous commands status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Sharing.BlockDangerousCommands, nil
			}).
			WithValueSetter("Set block dangerous commands status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.BlockDangerousCommands = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle block dangerous commands", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Sharing.BlockDangerousCommands = !cfg.Sharing.BlockDangerousCommands
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Values subcommands registered as separate commands
		{
			Name:        "sharedcategories-values",
			Description: "List available command categories",
			Handler: func(ctx *router.Ctx) {
				var categories []string
				for _, category := range ctx.Router.GetCategories() {
					categories = append(categories, category.Name)
				}

				c.services.Embed.RespondField(ctx, "Available Categories", strings.Join(categories, "\n"))
			},
		},
		{
			Name:        "sharedcommands-values",
			Description: "Show help for available commands",
			Handler: func(ctx *router.Ctx) {
				c.services.Embed.RespondField(ctx, "Available Commands", "Use `help` or view <https://docs.lightning-bot.com/category/commands> for a list of commands.")
			},
		},
	})
}
