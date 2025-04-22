package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/actioncommands"
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type ConfigCommands struct{ services *services.Services }

func NewConfigCommands(s *services.Services) *ConfigCommands {
	return &ConfigCommands{services: s}
}

func (c *ConfigCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Config",
		Description: "Configuration commands",
	}, []*router.Command{
		// Command Prefix
		actioncommands.NewBuilder("commandprefix", "Command Prefix").
			WithValueGetter("Get the current command prefix", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.CommandPrefix, nil
			}).
			WithValueSetter("Set the command prefix", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.CommandPrefix = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Command Mode
		actioncommands.NewBuilder("commandmode", "Command Mode").
			WithValueGetter("Get the current command mode", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.CommandMode, nil
			}).
			WithValueSetter("Set the command mode", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.CommandMode = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available command modes", func() []string {
				return []string{"Codeblock", "Text", "Indent", "WebhookEmbed", "URLEmbed"}
			}).
			Build(c.services),

		// Embed Title
		actioncommands.NewBuilder("title", "Embed Title").
			WithValueGetter("Get the current embed title", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.Title, nil
			}).
			WithValueSetter("Set the embed title", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.Title = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Embed Footer
		actioncommands.NewBuilder("footer", "Embed Footer").
			WithValueGetter("Get the current embed footer", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.Footer, nil
			}).
			WithValueSetter("Set the embed footer", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.Footer = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Delete Responses
		actioncommands.NewBuilder("delresponses", "Delete Responses").
			WithValueGetter("Get delete responses status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.DeleteResponses, nil
			}).
			WithValueSetter("Set delete responses status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.DeleteResponses = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle delete responses", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Delete Responses Delay
		actioncommands.NewBuilder("delresponsesdelay", "Delete Responses Delay").
			WithValueGetter("Get delete responses delay", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.DeleteResponsesDelay, nil
			}).
			WithValueSetter("Set delete responses delay", actionvalue.TypeDuration, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.DeleteResponsesDelay = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Delete Commands
		actioncommands.NewBuilder("delcommands", "Delete Commands").
			WithValueGetter("Get delete commands status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.DeleteCommands, nil
			}).
			WithValueSetter("Set delete commands status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.DeleteCommands = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle delete commands", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Client Spoof
		actioncommands.NewBuilder("clientspoof", "Client Spoof").
			WithValueGetter("Get current client spoof setting", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.Client, nil
			}).
			WithValueSetter("Set client spoof setting", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.Client = utils.String.Capitalize(value.(string))
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available clients", func() []string {
				return []string{"Mobile", "Web", "Desktop", "Console"}
			}).
			Build(c.services),

		// Timezone
		actioncommands.NewBuilder("timezone", "Timezone").
			WithValueGetter("Get current timezone", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.Timezone, nil
			}).
			WithValueSetter("Set timezone", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				timezone := value.(string)
				if _, err := time.LoadLocation(timezone); err != nil {
					return fmt.Errorf("invalid timezone: %v", err)
				}

				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.Timezone = timezone
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Embed Color
		actioncommands.NewBuilder("embedcolor", "Embed Color").
			WithValueGetter("Get current embed color", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.EmbedColor, nil
			}).
			WithValueSetter("Set embed color", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.EmbedColor = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Embed Image
		actioncommands.NewBuilder("embedimage", "Embed Image").
			WithValueGetter("Get current embed image", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.EmbedImage, nil
			}).
			WithValueSetter("Set embed image", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.EmbedImage = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Embed URL
		actioncommands.NewBuilder("embedurl", "Embed URL").
			WithValueGetter("Get current embed URL", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Selfbot.EmbedURL, nil
			}).
			WithValueSetter("Set embed URL", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.EmbedURL = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Force Subscribed Servers
		actioncommands.NewBuilder("forcesubs", "Force Subscribed Servers").
			WithListGetter("Get force subscribed servers", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Selfbot.ForceSubscribedServers) == 0 {
					c.services.Embed.Respond(ctx, "No force subscribed servers found.")
					return nil
				}

				servers := fn.Map(cfg.Selfbot.ForceSubscribedServers, func(guildID string) string {
					guild, err := ctx.Session.GuildWithFallback(guildID)
					if err != nil {
						return fmt.Sprintf("%s (Unknown)", guildID)
					}

					return fmt.Sprintf("%s (%s)", guild.Name, guild.ID)
				})

				c.services.Embed.RespondField(ctx, "Force Subscribed Servers", strings.Join(servers, "\n"))
				return nil
			}).
			WithListAdder("Add force subscribed server", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Selfbot.ForceSubscribedServers, guild.ID) {
					return fmt.Errorf("server %s (%s) is already in the force subscribed servers list", guild.Name, guild.ID)
				}

				cfg.Selfbot.ForceSubscribedServers = append(cfg.Selfbot.ForceSubscribedServers, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added server %s (%s) to force subscribed servers.", guild.Name, guild.ID)
				return nil
			}).
			WithListRemover("Remove force subscribed server", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Selfbot.ForceSubscribedServers, guild.ID) {
					return fmt.Errorf("server %s (%s) is not in the force subscribed servers list", guild.Name, guild.ID)
				}

				cfg.Selfbot.ForceSubscribedServers = fn.Delete(cfg.Selfbot.ForceSubscribedServers, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed server %s (%s) from force subscribed servers.", guild.Name, guild.ID)
				return nil
			}).
			WithListClearer("Clear force subscribed servers", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.ForceSubscribedServers = make([]string, 0)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all force subscribed servers.")
				return nil
			}).
			Build(c.services),

		// Aliases
		actioncommands.NewBuilder("aliases", "Command Aliases").
			WithListGetter("Get command aliases", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				prefix := cfg.Selfbot.CommandPrefix

				if len(cfg.Selfbot.CommandAliases) == 0 {
					c.services.Embed.Respond(ctx, "No command aliases found.")
					return nil
				}

				aliases := make([]string, 0, len(cfg.Selfbot.CommandAliases))

				for alias, command := range cfg.Selfbot.CommandAliases {
					aliases = append(aliases, fmt.Sprintf("%s%s -> %s%s", prefix, alias, prefix, command))
				}

				c.services.Embed.RespondField(ctx, "Command Aliases", strings.Join(aliases, "\n"))
				return nil
			}).
			WithListAdder("Add command alias", "<alias> <command>", func(ctx *router.Ctx) error {
				alias, cmd := ctx.Arguments.Get(0).Raw(), ctx.Arguments.Get(1).Raw()
				cfg := c.services.Config.GetConfig()

				if c.services.Router.GetCommand(cmd) == nil {
					return fmt.Errorf("command %s not found", cmd)
				}

				cfg.Selfbot.CommandAliases[alias] = cmd
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Alias %s -> %s added.", alias, cmd)

				return c.services.EventReactions.LoadCmdAliases()
			}).
			WithListRemover("Remove command alias", "<alias>", func(ctx *router.Ctx) error {
				alias := ctx.Arguments.Get(0).Raw()
				cfg := c.services.Config.GetConfig()

				cmd, ok := cfg.Selfbot.CommandAliases[alias]
				if !ok {
					return fmt.Errorf("alias %s not found", alias)
				}

				delete(cfg.Selfbot.CommandAliases, alias)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Alias %s -> %s removed.", alias, cmd)

				return c.services.EventReactions.LoadCmdAliases()
			}).
			WithListClearer("Clear command aliases", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Selfbot.CommandAliases = make(map[string]string)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all command aliases.")

				return c.services.EventReactions.LoadCmdAliases()
			}).
			Build(c.services),

		// Notifications
		actioncommands.NewBuilder("notifications", "Notifications").
			WithValueGetter("Get notifications status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Notifications.Enabled, nil
			}).
			WithValueSetter("Set notifications status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Notifications.Enabled = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle notifications", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Notifications.Enabled = !cfg.Notifications.Enabled
				return c.services.Config.SaveConfig(cfg)
			}).
			WithExample("Show an example of a notification", func(ctx *router.Ctx) error {
				return c.services.UI.ShowNotification("This is an example notification.")
			}).
			Build(c.services),

		// Nitro Sniper
		actioncommands.NewBuilder("nitrosniper", "Nitro Sniper").
			WithValueGetter("Get nitro sniper status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.NitroSniper, nil
			}).
			WithValueSetter("Set nitro sniper status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.NitroSniper = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle nitro sniper", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Giveaway Sniper
		actioncommands.NewBuilder("giveawaysniper", "Giveaway Sniper").
			WithValueGetter("Get giveaway sniper status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.GiveawaySniper.Enabled, nil
			}).
			WithValueSetter("Set giveaway sniper status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.GiveawaySniper.Enabled = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle giveaway sniper", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Giveaway Sniper Sleep
		actioncommands.NewBuilder("giveawaysnipersleep", "Giveaway Sniper Sleep").
			WithValueGetter("Get giveaway sniper sleep duration", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.GiveawaySniper.Sleep, nil
			}).
			WithValueSetter("Set giveaway sniper sleep duration", actionvalue.TypeDuration, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.GiveawaySniper.Sleep = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Custom Bots
		actioncommands.NewBuilder("custombots", "Custom Giveaway Sniper Bots").
			WithListGetter("Get custom bots", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Snipers.GiveawaySniper.CustomBots) == 0 {
					c.services.Embed.Respond(ctx, "No custom bots found.")
					return nil
				}

				bots := fn.Map(cfg.Snipers.GiveawaySniper.CustomBots, func(botID string) string {
					user, err := ctx.Session.UserWithFallback(botID)
					if err != nil {
						return fmt.Sprintf("%s (Unknown)", botID)
					}

					return fmt.Sprintf("%s (%s)", user.Username, user.ID)
				})

				c.services.Embed.RespondField(ctx, "Custom Giveaway Sniper Bots", strings.Join(bots, "\n"))
				return nil
			}).
			WithListAdder("Add custom bot", "<@bot>", func(ctx *router.Ctx) error {
				bot, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find bot: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID) {
					return fmt.Errorf("bot %s (%s) is already in the custom bots list", bot.Username, bot.ID)
				}

				cfg.Snipers.GiveawaySniper.CustomBots = append(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added bot %s (%s) to custom bots.", bot.Username, bot.ID)
				return nil
			}).
			WithListRemover("Remove custom bot", "<@bot>", func(ctx *router.Ctx) error {
				bot, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find bot: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID) {
					return fmt.Errorf("bot %s (%s) is not in the custom bots list", bot.Username, bot.ID)
				}

				cfg.Snipers.GiveawaySniper.CustomBots = fn.Delete(cfg.Snipers.GiveawaySniper.CustomBots, bot.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed bot %s (%s) from custom bots.", bot.Username, bot.ID)
				return nil
			}).
			WithListClearer("Clear custom bots", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.GiveawaySniper.CustomBots = make([]string, 0)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all custom bots.")
				return nil
			}).
			Build(c.services),

		// Sniper Blacklist
		actioncommands.NewBuilder("sniperblacklist", "Sniper Blacklist").
			WithListGetter("Get sniper blacklist", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Snipers.Blacklist) == 0 {
					c.services.Embed.Respond(ctx, "No sniper blacklist entries found.")
					return nil
				}

				c.services.Embed.RespondField(ctx, "Sniper Blacklist", strings.Join(cfg.Snipers.Blacklist, "\n"))
				return nil
			}).
			WithListAdder("Add to sniper blacklist", "<value>", func(ctx *router.Ctx) error {
				value := ctx.Arguments.Get(0).Raw()

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Snipers.Blacklist, value) {
					return fmt.Errorf("value '%s' is already in the sniper blacklist", value)
				}

				cfg.Snipers.Blacklist = append(cfg.Snipers.Blacklist, value)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added '%s' to sniper blacklist.", value)
				return nil
			}).
			WithListRemover("Remove from sniper blacklist", "<value>", func(ctx *router.Ctx) error {
				value := ctx.Arguments.Get(0).Raw()

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Snipers.Blacklist, value) {
					return fmt.Errorf("value '%s' is not in the sniper blacklist", value)
				}

				cfg.Snipers.Blacklist = fn.Delete(cfg.Snipers.Blacklist, value)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed '%s' from sniper blacklist.", value)
				return nil
			}).
			WithListClearer("Clear sniper blacklist", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.Blacklist = make([]string, 0)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared sniper blacklist.")
				return nil
			}).
			Build(c.services),

		// Server Blacklist
		actioncommands.NewBuilder("serverblacklist", "Server Blacklist").
			WithListGetter("Get server blacklist", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Snipers.ServerBlacklist) == 0 {
					c.services.Embed.Respond(ctx, "No server blacklist entries found.")
					return nil
				}

				servers := fn.Map(cfg.Snipers.ServerBlacklist, func(guildID string) string {
					guild, err := ctx.Session.GuildWithFallback(guildID)
					if err != nil {
						return fmt.Sprintf("%s (Unknown)", guildID)
					}

					return fmt.Sprintf("%s (%s)", guild.Name, guild.ID)
				})

				c.services.Embed.RespondField(ctx, "Server Blacklist", strings.Join(servers, "\n"))
				return nil
			}).
			WithListAdder("Add to server blacklist", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Snipers.ServerBlacklist, guild.ID) {
					return fmt.Errorf("server %s (%s) is already in the server blacklist", guild.Name, guild.ID)
				}

				cfg.Snipers.ServerBlacklist = append(cfg.Snipers.ServerBlacklist, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added server %s (%s) to server blacklist.", guild.Name, guild.ID)
				return nil
			}).
			WithListRemover("Remove from server blacklist", "[server id]", func(ctx *router.Ctx) error {
				guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
				if err != nil {
					return fmt.Errorf("could not find server: %v", err)
				}

				cfg := c.services.Config.GetConfig()

				if !slices.Contains(cfg.Snipers.ServerBlacklist, guild.ID) {
					return fmt.Errorf("server %s (%s) is not in the server blacklist", guild.Name, guild.ID)
				}

				cfg.Snipers.ServerBlacklist = fn.Delete(cfg.Snipers.ServerBlacklist, guild.ID)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed server %s (%s) from server blacklist.", guild.Name, guild.ID)
				return nil
			}).
			WithListClearer("Clear server blacklist", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.ServerBlacklist = make([]string, 0)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared server blacklist.")
				return nil
			}).
			Build(c.services),

		// Sniper Tokens
		actioncommands.NewBuilder("snipertokens", "Sniper Tokens").
			WithListGetter("Get sniper tokens", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				if len(cfg.Snipers.Tokens) == 0 {
					c.services.Embed.Respond(ctx, "No sniper tokens found.")
					return nil
				}

				// Create a masked list of tokens for security
				tokens := fn.MapIndexed(cfg.Snipers.Tokens, func(i int, token string) string {
					if len(token) <= 12 {
						return strconv.Itoa(i+1) + ". ••••••••••••"
					}
					return strconv.Itoa(i+1) + ". " + token[:6] + "••••••••••••" + token[len(token)-6:]
				})

				c.services.Embed.RespondField(ctx, "Sniper Tokens", strings.Join(tokens, "\n"))
				return nil
			}).
			WithListAdder("Add sniper token", "<token>", func(ctx *router.Ctx) error {
				token := ctx.Arguments.Get(0).Raw()
				if token == "" {
					return fmt.Errorf("please provide a token to add")
				}

				cfg := c.services.Config.GetConfig()

				if slices.Contains(cfg.Snipers.Tokens, token) {
					return fmt.Errorf("token is already in the sniper tokens list")
				}

				cfg.Snipers.Tokens = append(cfg.Snipers.Tokens, token)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Added token to sniper tokens.")
				return nil
			}).
			WithListRemover("Remove sniper token", "<index>", func(ctx *router.Ctx) error {
				index, err := ctx.Arguments.Get(0).Int()
				if err != nil {
					return fmt.Errorf("invalid index: %v", err)
				}

				cfg := c.services.Config.GetConfig()
				index--

				if index < 0 || index >= len(cfg.Snipers.Tokens) {
					return fmt.Errorf("index out of range (0-%d)", len(cfg.Snipers.Tokens)-1)
				}

				cfg.Snipers.Tokens = slices.Delete(cfg.Snipers.Tokens, index, index+1)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Removed token from sniper tokens.")
				return nil
			}).
			WithListClearer("Clear sniper tokens", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.Tokens = make([]string, 0)
				if err := c.services.Config.SaveConfig(cfg); err != nil {
					return fmt.Errorf("error saving config: %v", err)
				}

				c.services.Embed.Respond(ctx, "Cleared all sniper tokens.")
				return nil
			}).
			Build(c.services),

		// Multi-Token Sniper
		actioncommands.NewBuilder("multitokensniper", "Multi-Token Sniper").
			WithValueGetter("Get multi-token sniper status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.MultiTokenSniper, nil
			}).
			WithValueSetter("Set multi-token sniper status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.MultiTokenSniper = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle multi-token sniper", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.MultiTokenSniper = !cfg.Snipers.MultiTokenSniper
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Privnote Sniper
		actioncommands.NewBuilder("privnotesniper", "Privnote Sniper").
			WithValueGetter("Get privnote sniper status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PrivnoteSniper, nil
			}).
			WithValueSetter("Set privnote sniper status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PrivnoteSniper = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle privnote sniper", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper
		actioncommands.NewBuilder("poketwosniper", "Poketwo Sniper").
			WithValueGetter("Get poketwo sniper status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.Enabled, nil
			}).
			WithValueSetter("Set poketwo sniper status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.Enabled = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle poketwo sniper", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.Enabled = !cfg.Snipers.PokeTwoSniper.Enabled
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Sleep
		actioncommands.NewBuilder("poketwosleep", "Poketwo Sniper Sleep").
			WithValueGetter("Get poketwo sniper sleep duration", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.Sleep, nil
			}).
			WithValueSetter("Set poketwo sniper sleep duration", actionvalue.TypeDuration, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.Sleep = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Random Sleep
		actioncommands.NewBuilder("poketworandomsleep", "Poketwo Sniper Random Sleep").
			WithValueGetter("Get poketwo sniper random sleep status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.RandomSleep, nil
			}).
			WithValueSetter("Set poketwo sniper random sleep status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.RandomSleep = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle poketwo sniper random sleep", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.RandomSleep = !cfg.Snipers.PokeTwoSniper.RandomSleep
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Timeout
		actioncommands.NewBuilder("poketwotimeout", "Poketwo Sniper Timeout").
			WithValueGetter("Get poketwo sniper timeout duration", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.Timeout, nil
			}).
			WithValueSetter("Set poketwo sniper timeout duration", actionvalue.TypeDuration, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.Timeout = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Typing
		actioncommands.NewBuilder("poketwotyping", "Poketwo Sniper Typing").
			WithValueGetter("Get poketwo sniper typing status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.Typing, nil
			}).
			WithValueSetter("Set poketwo sniper typing status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.Typing = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle poketwo sniper typing", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.Typing = !cfg.Snipers.PokeTwoSniper.Typing
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Catch Rate
		actioncommands.NewBuilder("poketwocatchrate", "Poketwo Sniper Catch Rate").
			WithValueGetter("Get poketwo sniper catch rate", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.CatchRate, nil
			}).
			WithValueSetter("Set poketwo sniper catch rate", actionvalue.TypeInt, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.CatchRate = value.(int)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Fake Activity
		actioncommands.NewBuilder("poketwofakeactivity", "Poketwo Sniper Fake Activity").
			WithValueGetter("Get poketwo sniper fake activity status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.FakeActivity, nil
			}).
			WithValueSetter("Set poketwo sniper fake activity status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.FakeActivity = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle poketwo sniper fake activity", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.FakeActivity = !cfg.Snipers.PokeTwoSniper.FakeActivity
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Poketwo Sniper Channels

		// Poketwo Sniper Language
		actioncommands.NewBuilder("poketwolanguage", "Poketwo Sniper Language").
			WithValueGetter("Get poketwo sniper language", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Snipers.PokeTwoSniper.CatchLanguage, nil
			}).
			WithValueSetter("Set poketwo sniper language", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Snipers.PokeTwoSniper.CatchLanguage = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available languages", func() []string {
				return []string{"en", "ja", "ja_t", "fr", "de"}
			}).
			Build(c.services),

		// Setup Webhooks
		{
			Name:        "setupwebhooks",
			Description: "Setup all webhooks automatically",
			Arguments:   "[server id]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.setupwebhooks,
		},

		// Webhook all
		{
			Name:        "wall",
			Description: "Set all webhooks to the same URL",
			Arguments:   "[webhook url]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.wall,
		},

		// Nitro Sniper Webhook
		actioncommands.NewBuilder("wnitro", "Nitro Sniper Webhook").
			WithValueGetter("Get nitro sniper webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.NitroSniper, nil
			}).
			WithValueSetter("Set nitro sniper webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.NitroSniper = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Giveaway Sniper Webhook
		actioncommands.NewBuilder("wgiveawaysniper", "Giveaway Sniper Webhook").
			WithValueGetter("Get giveaway sniper webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.GiveawaySniper, nil
			}).
			WithValueSetter("Set giveaway sniper webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.GiveawaySniper = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Selfbot Detection Webhook
		actioncommands.NewBuilder("wselfbotdetection", "Selfbot Detection Webhook").
			WithValueGetter("Get selfbot detection webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.SelfbotDetection, nil
			}).
			WithValueSetter("Set selfbot detection webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.SelfbotDetection = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Staff Detection Webhook
		actioncommands.NewBuilder("wstaffdetection", "Staff Detection Webhook").
			WithValueGetter("Get staff detection webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.StaffDetection, nil
			}).
			WithValueSetter("Set staff detection webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.StaffDetection = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Notifiers Webhook
		actioncommands.NewBuilder("wnotifiers", "Notifiers Webhook").
			WithValueGetter("Get notifiers webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.Notifiers, nil
			}).
			WithValueSetter("Set notifiers webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.Notifiers = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Ping Detection Webhook
		actioncommands.NewBuilder("wpingdetection", "Ping Detection Webhook").
			WithValueGetter("Get ping detection webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.PingDetection, nil
			}).
			WithValueSetter("Set ping detection webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.PingDetection = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Misc Webhook
		actioncommands.NewBuilder("wmisc", "Misc Webhook").
			WithValueGetter("Get misc webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.Misc, nil
			}).
			WithValueSetter("Set misc webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.Misc = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Sessions Webhook
		actioncommands.NewBuilder("wsessions", "Sessions Webhook").
			WithValueGetter("Get sessions webhook", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.Sessions, nil
			}).
			WithValueSetter("Set sessions webhook", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.Sessions = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Ping Self Webhook Mode
		actioncommands.NewBuilder("wpingself", "Ping Self Webhook").
			WithValueGetter("Get ping self webhook mode", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.PingSelfMode, nil
			}).
			WithValueSetter("Set ping self webhook mode", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.PingSelfMode = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available ping self webhook modes", func() []string {
				return []string{"Disabled", "Always", "Important"}
			}).
			Build(c.services),

		// Ping Everyone Webhook Mode
		actioncommands.NewBuilder("wpingeveryone", "Ping Everyone Webhook").
			WithValueGetter("Get ping everyone webhook mode", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Webhooks.PingEveryoneMode, nil
			}).
			WithValueSetter("Set ping everyone webhook mode", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Webhooks.PingEveryoneMode = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available ping everyone webhook modes", func() []string {
				return []string{"Disabled", "Always", "Important"}
			}).
			Build(c.services),

		// Event settings
		// Ban Notifier
		actioncommands.NewBuilder("bannotifier", "Ban Notifier").
			WithValueGetter("Get ban notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.BanNotifier, nil
			}).
			WithValueSetter("Set ban notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.BanNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle ban notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.BanNotifier = !cfg.Events.BanNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Staff Detection
		actioncommands.NewBuilder("staffdetection", "Staff Detection").
			WithValueGetter("Get staff detection status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.DiscordStaffDetection, nil
			}).
			WithValueSetter("Set staff detection status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DiscordStaffDetection = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle staff detection", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DiscordStaffDetection = !cfg.Events.DiscordStaffDetection
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Ping Detection
		actioncommands.NewBuilder("pingdetection", "Ping Detection").
			WithValueGetter("Get ping detection status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.GhostPingDetection, nil
			}).
			WithValueSetter("Set ping detection status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.GhostPingDetection = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle ping detection", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.GhostPingDetection = !cfg.Events.GhostPingDetection
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Ticket Notifier
		actioncommands.NewBuilder("ticketnotifier", "Ticket Notifier").
			WithValueGetter("Get ticket notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.TicketNotifier, nil
			}).
			WithValueSetter("Set ticket notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.TicketNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle ticket notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.TicketNotifier = !cfg.Events.TicketNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Session Connected
		actioncommands.NewBuilder("sessionconnected", "Session Connected").
			WithValueGetter("Get session connected status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.SessionConnected, nil
			}).
			WithValueSetter("Set session connected status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.SessionConnected = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle session connected", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.SessionConnected = !cfg.Events.SessionConnected
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Session Disconnected
		actioncommands.NewBuilder("sessiondisconnected", "Session Disconnected").
			WithValueGetter("Get session disconnected status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.SessionDisconnected, nil
			}).
			WithValueSetter("Set session disconnected status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.SessionDisconnected = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle session disconnected", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.SessionDisconnected = !cfg.Events.SessionDisconnected
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Selfbot Detection
		actioncommands.NewBuilder("selfbotdetection", "Selfbot Detection").
			WithValueGetter("Get selfbot detection status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.SelfbotDetection, nil
			}).
			WithValueSetter("Set selfbot detection status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.SelfbotDetection = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle selfbot detection", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.SelfbotDetection = !cfg.Events.SelfbotDetection
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// DM Typing Notifier
		actioncommands.NewBuilder("dmtypingnotifier", "DM Typing Notifier").
			WithValueGetter("Get DM typing notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.DMTypingNotifier, nil
			}).
			WithValueSetter("Set DM typing notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMTypingNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle DM typing notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMTypingNotifier = !cfg.Events.DMTypingNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// DM Notifier
		actioncommands.NewBuilder("dmnotifier", "DM Notifier").
			WithValueGetter("Get DM notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.DMNotifier, nil
			}).
			WithValueSetter("Set DM notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle DM notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMNotifier = !cfg.Events.DMNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// DM Edit Notifier
		actioncommands.NewBuilder("dmeditnotifier", "DM Edit Notifier").
			WithValueGetter("Get DM edit notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.DMEditNotifier, nil
			}).
			WithValueSetter("Set DM edit notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMEditNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle DM edit notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMEditNotifier = !cfg.Events.DMEditNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// DM Delete Notifier
		actioncommands.NewBuilder("dmdeletenotifier", "DM Delete Notifier").
			WithValueGetter("Get DM delete notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.DMDeleteNotifier, nil
			}).
			WithValueSetter("Set DM delete notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMDeleteNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle DM delete notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.DMDeleteNotifier = !cfg.Events.DMDeleteNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Message Delete Notifier
		actioncommands.NewBuilder("messagedeletenotifier", "Message Delete Notifier").
			WithValueGetter("Get message delete notifier status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.Events.MessageDeleteNotifier, nil
			}).
			WithValueSetter("Set message delete notifier status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.MessageDeleteNotifier = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle message delete notifier", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.Events.MessageDeleteNotifier = !cfg.Events.MessageDeleteNotifier
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// UI Theme
		actioncommands.NewBuilder("uitheme", "UI Theme").
			WithValueGetter("Get UI theme", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.UI.Theme, nil
			}).
			WithValueSetter("Set UI theme", actionvalue.TypeString, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.Theme = value.(string)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValues("List available UI themes", func() []string {
				return []string{"Dark", "Light"}
			}).
			Build(c.services),

		// Hide to System Tray
		actioncommands.NewBuilder("hidetosystray", "Hide to System Tray").
			WithValueGetter("Get hide to system tray status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.UI.HideToSystemTray, nil
			}).
			WithValueSetter("Set hide to system tray status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.HideToSystemTray = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle hide to system tray", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.HideToSystemTray = !cfg.UI.HideToSystemTray
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Run UI on Startup
		actioncommands.NewBuilder("runuionstartup", "Run UI on Startup").
			WithValueGetter("Get run UI on startup status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.UI.RunOnStartup, nil
			}).
			WithValueSetter("Set run UI on startup status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.RunOnStartup = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle run UI on startup", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),

		// Run Bot on Startup
		actioncommands.NewBuilder("runbotonstartup", "Run Bot on Startup").
			WithValueGetter("Get run bot on startup status", func(ctx *router.Ctx) (any, error) {
				cfg := c.services.Config.GetConfig()
				return cfg.UI.RunOnStartup, nil
			}).
			WithValueSetter("Set run bot on startup status", actionvalue.TypeBool, func(ctx *router.Ctx, value any) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.RunOnStartup = value.(bool)
				return c.services.Config.SaveConfig(cfg)
			}).
			WithValueToggler("Toggle run bot on startup", func(ctx *router.Ctx) error {
				cfg := c.services.Config.GetConfig()
				cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
				return c.services.Config.SaveConfig(cfg)
			}).
			Build(c.services),
	})
}

func createWebhookChannel(ctx *router.Ctx, serverID, categoryID, name string) (string, error) {
	channel, err := ctx.Session.GuildChannelCreateComplex(serverID, discordgo.GuildChannelCreateData{
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
		Name:     name,
	})
	if err != nil {
		return "", err
	}

	webhook, err := ctx.Session.WebhookCreate(channel.ID, name, "")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.Token), nil
}

func (c *ConfigCommands) setupwebhooks(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Could not find server: %v", err)
	}

	category, err := ctx.Session.GuildChannelCreate(guild.ID, "LightningBot", discordgo.ChannelTypeGuildCategory)
	if err != nil {
		return fmt.Errorf("Failed to create category: %v", err)
	}

	cfg := c.services.Config.GetConfig()
	cfg.Webhooks.NitroSniper, _ = createWebhookChannel(ctx, guild.ID, category.ID, "nitrosniper")
	cfg.Webhooks.GiveawaySniper, _ = createWebhookChannel(ctx, guild.ID, category.ID, "giveawaysniper")
	cfg.Webhooks.SelfbotDetection, _ = createWebhookChannel(ctx, guild.ID, category.ID, "selfbotdetection")
	cfg.Webhooks.StaffDetection, _ = createWebhookChannel(ctx, guild.ID, category.ID, "staffdetection")
	cfg.Webhooks.Notifiers, _ = createWebhookChannel(ctx, guild.ID, category.ID, "notifiers")
	cfg.Webhooks.PingDetection, _ = createWebhookChannel(ctx, guild.ID, category.ID, "pingdetection")
	cfg.Webhooks.Misc, _ = createWebhookChannel(ctx, guild.ID, category.ID, "misc")
	cfg.Webhooks.Sessions, _ = createWebhookChannel(ctx, guild.ID, category.ID, "sessions")

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("Failed to save config: %v", err)
	}

	c.services.Embed.Respond(ctx, "Webhooks setup successfully.")
	return nil
}

func (c *ConfigCommands) wall(ctx *router.Ctx) error {
	cfg := c.services.Config.GetConfig()

	webhookURL := ctx.Arguments.Get(0).Raw()
	cfg.Webhooks.NitroSniper = webhookURL
	cfg.Webhooks.GiveawaySniper = webhookURL
	cfg.Webhooks.SelfbotDetection = webhookURL
	cfg.Webhooks.StaffDetection = webhookURL
	cfg.Webhooks.Notifiers = webhookURL
	cfg.Webhooks.PingDetection = webhookURL
	cfg.Webhooks.Misc = webhookURL
	cfg.Webhooks.Sessions = webhookURL

	if err := c.services.Config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("Error saving config: %v", err)
	}

	c.services.Embed.Respond(ctx, "All webhooks set to %s", webhookURL)
	return nil
}
