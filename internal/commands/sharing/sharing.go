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

				users := fn.Map(cfg.Sharing.SharedUsers, func(userID string) string {
					user, err := ctx.Session.UserWithFallback(userID)
					if err != nil {
						return fmt.Sprintf("%s (Unknown)", userID)
					}
					return fmt.Sprintf("%s (%s)", user.Username, userID)
				})

				c.services.Embed.RespondField(ctx, "Shared Users", strings.Join(users, "\n"))
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

				categories := ctx.Router.GetCategories()
				filtered := fn.Filter(categories, func(category *router.Category) bool {
					return strings.EqualFold(category.Name, categoryName)
				})
				if len(filtered) == 0 {
					return fmt.Errorf("category %s not found. Use the 'values' subcommand to see available categories", categoryName)
				}
				actualName := filtered[0].Name

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

				categories := ctx.Router.GetCategories()
				filtered := fn.Filter(categories, func(category *router.Category) bool {
					return strings.EqualFold(category.Name, categoryName)
				})
				if len(filtered) == 0 {
					return fmt.Errorf("category %s not found. Use the 'values' subcommand to see available categories", categoryName)
				}
				actualName := filtered[0].Name

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
			WithValues("Get available categories", func() []string {
				return fn.Map(c.services.Router.GetCategories(), func(category *router.Category) string {
					return category.Name
				})
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

				commands := fn.Filter(ctx.Router.GetCommands(), func(command *router.Command) bool {
					return strings.EqualFold(command.Name, commandName)
				})

				if len(commands) == 0 {
					return fmt.Errorf("command %s not found", commandName)
				}

				actualName := commands[0].Name
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

				commands := fn.Filter(ctx.Router.GetCommands(), func(command *router.Command) bool {
					return strings.EqualFold(command.Name, commandName)
				})

				if len(commands) == 0 {
					return fmt.Errorf("command %s not found", commandName)
				}

				actualName := commands[0].Name
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
			WithValues("Get available commands", func() []string {
				return []string{"Use the `help` command or view <https://docs.lightning-bot.com/category/commands> for a list of available commands."}
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
	})
}
