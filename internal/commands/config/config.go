package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/lb-selfbot/selfbot/internal/actioncommands"
	"github.com/lb-selfbot/selfbot/internal/actioncommands/actionvalue"
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
			WithValues("List available command modes", func() []any {
				return []any{"Codeblock", "Text", "Indent", "WebhookEmbed", "URLEmbed"}
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
			WithValues("List available clients", func() []any {
				return []any{"Mobile", "Web", "Desktop", "Console"}
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

		// TODO: Force Subscribed Servers

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

		// TODO: Custom Bots

		// TODO: Sniper Blacklist

		// TODO: Server Blacklist

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

		// TODO: Sniper Tokens

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

		// TODO: Poketwo Channels

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
			WithValues("List available languages", func() []any {
				return []any{"en", "ja", "ja_t", "fr", "de"}
			}).
			Build(c.services),

		// TODO: Setup Webhooks

		// TODO: Webhook all

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
			WithValues("List available ping self webhook modes", func() []any {
				return []any{"Disabled", "Always", "Important"}
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
			WithValues("List available ping everyone webhook modes", func() []any {
				return []any{"Disabled", "Always", "Important"}
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
			WithValues("List available UI themes", func() []any {
				return []any{"Dark", "Light"}
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
