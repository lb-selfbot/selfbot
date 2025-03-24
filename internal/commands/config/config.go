package config

import (
	"fmt"
	"time"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type ConfigCommands struct{ services *services.Services }

func NewConfigCommands(s *services.Services) *ConfigCommands {
	return &ConfigCommands{services: s}
}

func (c *ConfigCommands) Register() {
	commands := make([]*router.Command, 0)

	// Command Prefix
	commands = append(commands, router.NewActionCommandBuilder("commandprefix", "Command Prefix").
		WithGetter("Get the current command prefix", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.CommandPrefix, nil
		}).
		WithSetter("Set the command prefix", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.CommandPrefix = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Command Mode
	commands = append(commands, router.NewActionCommandBuilder("commandmode", "Command Mode").
		WithGetter("Get the current command mode", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.CommandMode, nil
		}).
		WithSetter("Set the command mode", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.CommandMode = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available command modes", func() []any {
			return []any{"Codeblock", "Text", "Indent", "WebhookEmbed", "URLEmbed"}
		}).
		Build(c.services.Router))

	// Embed Title
	commands = append(commands, router.NewActionCommandBuilder("title", "Embed Title").
		WithGetter("Get the current embed title", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Title, nil
		}).
		WithSetter("Set the embed title", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Title = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Embed Footer
	commands = append(commands, router.NewActionCommandBuilder("footer", "Embed Footer").
		WithGetter("Get the current embed footer", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Footer, nil
		}).
		WithSetter("Set the embed footer", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Footer = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Delete Responses
	commands = append(commands, router.NewActionCommandBuilder("delresponses", "Delete Responses").
		WithGetter("Get delete responses status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.DeleteResponses, nil
		}).
		WithSetter("Set delete responses status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteResponses = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle delete responses", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Delete Responses Delay
	commands = append(commands, router.NewActionCommandBuilder("delresponsesdelay", "Delete Responses Delay").
		WithGetter("Get delete responses delay", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.DeleteResponsesDelay, nil
		}).
		WithSetter("Set delete responses delay", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteResponsesDelay = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Delete Commands
	commands = append(commands, router.NewActionCommandBuilder("delcommands", "Delete Commands").
		WithGetter("Get delete commands status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.DeleteCommands, nil
		}).
		WithSetter("Set delete commands status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteCommands = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle delete commands", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Client Spoof
	commands = append(commands, router.NewActionCommandBuilder("clientspoof", "Client Spoof").
		WithGetter("Get current client spoof setting", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Client, nil
		}).
		WithSetter("Set client spoof setting", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Client = utils.String.Capitalize(value.(string))
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available clients", func() []any {
			return []any{"Mobile", "Web", "Desktop", "Console"}
		}).
		Build(c.services.Router))

	// Timezone
	commands = append(commands, router.NewActionCommandBuilder("timezone", "Timezone").
		WithGetter("Get current timezone", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Timezone, nil
		}).
		WithSetter("Set timezone", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			timezone := value.(string)
			if _, err := time.LoadLocation(timezone); err != nil {
				return fmt.Errorf("invalid timezone: %v", err)
			}

			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Timezone = timezone
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Embed Color
	commands = append(commands, router.NewActionCommandBuilder("embedcolor", "Embed Color").
		WithGetter("Get current embed color", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.EmbedColor, nil
		}).
		WithSetter("Set embed color", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.EmbedColor = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Embed Image
	commands = append(commands, router.NewActionCommandBuilder("embedimage", "Embed Image").
		WithGetter("Get current embed image", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.EmbedImage, nil
		}).
		WithSetter("Set embed image", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.EmbedImage = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Embed URL
	commands = append(commands, router.NewActionCommandBuilder("embedurl", "Embed URL").
		WithGetter("Get current embed URL", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.EmbedURL, nil
		}).
		WithSetter("Set embed URL", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.EmbedURL = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Notifications
	commands = append(commands, router.NewActionCommandBuilder("notifications", "Notifications").
		WithGetter("Get notifications status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Notifications.Enabled, nil
		}).
		WithSetter("Set notifications status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Notifications.Enabled = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle notifications", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Notifications.Enabled = !cfg.Notifications.Enabled
			return c.services.Config.SaveConfig(cfg)
		}).
		WithExample("Show an example of a notification", func(ctx *router.Ctx) error {
			c.services.UI.ShowNotification("This is an example notification.")
			return nil
		}).
		Build(c.services.Router))

	// Nitro Sniper
	commands = append(commands, router.NewActionCommandBuilder("nitrosniper", "Nitro Sniper").
		WithGetter("Get nitro sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.NitroSniper, nil
		}).
		WithSetter("Set nitro sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.NitroSniper = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle nitro sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Giveaway Sniper
	commands = append(commands, router.NewActionCommandBuilder("giveawaysniper", "Giveaway Sniper").
		WithGetter("Get giveaway sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.GiveawaySniper.Enabled, nil
		}).
		WithSetter("Set giveaway sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.GiveawaySniper.Enabled = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle giveaway sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Giveaway Sniper Sleep
	commands = append(commands, router.NewActionCommandBuilder("giveawaysnipersleep", "Giveaway Sniper Sleep").
		WithGetter("Get giveaway sniper sleep duration", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.GiveawaySniper.Sleep, nil
		}).
		WithSetter("Set giveaway sniper sleep duration", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.GiveawaySniper.Sleep = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Multi-Token Sniper
	commands = append(commands, router.NewActionCommandBuilder("multitokensniper", "Multi-Token Sniper").
		WithGetter("Get multi-token sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.MultiTokenSniper, nil
		}).
		WithSetter("Set multi-token sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.MultiTokenSniper = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle multi-token sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.MultiTokenSniper = !cfg.Snipers.MultiTokenSniper
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Privnote Sniper
	commands = append(commands, router.NewActionCommandBuilder("privnotesniper", "Privnote Sniper").
		WithGetter("Get privnote sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PrivnoteSniper, nil
		}).
		WithSetter("Set privnote sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PrivnoteSniper = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle privnote sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper
	commands = append(commands, router.NewActionCommandBuilder("poketwosniper", "Poketwo Sniper").
		WithGetter("Get poketwo sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Enabled, nil
		}).
		WithSetter("Set poketwo sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Enabled = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Enabled = !cfg.Snipers.PokeTwoSniper.Enabled
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Sleep
	commands = append(commands, router.NewActionCommandBuilder("poketwosleep", "Poketwo Sniper Sleep").
		WithGetter("Get poketwo sniper sleep duration", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Sleep, nil
		}).
		WithSetter("Set poketwo sniper sleep duration", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Sleep = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Random Sleep
	commands = append(commands, router.NewActionCommandBuilder("poketworandomsleep", "Poketwo Sniper Random Sleep").
		WithGetter("Get poketwo sniper random sleep status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.RandomSleep, nil
		}).
		WithSetter("Set poketwo sniper random sleep status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.RandomSleep = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper random sleep", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.RandomSleep = !cfg.Snipers.PokeTwoSniper.RandomSleep
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Timeout
	commands = append(commands, router.NewActionCommandBuilder("poketwotimeout", "Poketwo Sniper Timeout").
		WithGetter("Get poketwo sniper timeout duration", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Timeout, nil
		}).
		WithSetter("Set poketwo sniper timeout duration", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Timeout = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Typing
	commands = append(commands, router.NewActionCommandBuilder("poketwotyping", "Poketwo Sniper Typing").
		WithGetter("Get poketwo sniper typing status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Typing, nil
		}).
		WithSetter("Set poketwo sniper typing status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Typing = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper typing", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Typing = !cfg.Snipers.PokeTwoSniper.Typing
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Catch Rate
	commands = append(commands, router.NewActionCommandBuilder("poketwocatchrate", "Poketwo Sniper Catch Rate").
		WithGetter("Get poketwo sniper catch rate", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.CatchRate, nil
		}).
		WithSetter("Set poketwo sniper catch rate", router.ValueTypeInt, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.CatchRate = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Fake Activity
	commands = append(commands, router.NewActionCommandBuilder("poketwofakeactivity", "Poketwo Sniper Fake Activity").
		WithGetter("Get poketwo sniper fake activity status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.FakeActivity, nil
		}).
		WithSetter("Set poketwo sniper fake activity status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.FakeActivity = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper fake activity", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.FakeActivity = !cfg.Snipers.PokeTwoSniper.FakeActivity
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Poketwo Sniper Language
	commands = append(commands, router.NewActionCommandBuilder("poketwolanguage", "Poketwo Sniper Language").
		WithGetter("Get poketwo sniper language", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.CatchLanguage, nil
		}).
		WithSetter("Set poketwo sniper language", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.CatchLanguage = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available languages", func() []any {
			return []any{"en", "ja", "ja_t", "fr", "de"}
		}).
		Build(c.services.Router))

	// Webhook settings
	// Nitro Sniper Webhook
	commands = append(commands, router.NewActionCommandBuilder("wnitro", "Nitro Sniper Webhook").
		WithGetter("Get nitro sniper webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.NitroSniper, nil
		}).
		WithSetter("Set nitro sniper webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.NitroSniper = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Giveaway Sniper Webhook
	commands = append(commands, router.NewActionCommandBuilder("wgiveawaysniper", "Giveaway Sniper Webhook").
		WithGetter("Get giveaway sniper webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.GiveawaySniper, nil
		}).
		WithSetter("Set giveaway sniper webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.GiveawaySniper = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Selfbot Detection Webhook
	commands = append(commands, router.NewActionCommandBuilder("wselfbotdetection", "Selfbot Detection Webhook").
		WithGetter("Get selfbot detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.SelfbotDetection, nil
		}).
		WithSetter("Set selfbot detection webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.SelfbotDetection = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Staff Detection Webhook
	commands = append(commands, router.NewActionCommandBuilder("wstaffdetection", "Staff Detection Webhook").
		WithGetter("Get staff detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.StaffDetection, nil
		}).
		WithSetter("Set staff detection webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.StaffDetection = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Notifiers Webhook
	commands = append(commands, router.NewActionCommandBuilder("wnotifiers", "Notifiers Webhook").
		WithGetter("Get notifiers webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.Notifiers, nil
		}).
		WithSetter("Set notifiers webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.Notifiers = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Ping Detection Webhook
	commands = append(commands, router.NewActionCommandBuilder("wpingdetection", "Ping Detection Webhook").
		WithGetter("Get ping detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.PingDetection, nil
		}).
		WithSetter("Set ping detection webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.PingDetection = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Misc Webhook
	commands = append(commands, router.NewActionCommandBuilder("wmisc", "Misc Webhook").
		WithGetter("Get misc webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.Misc, nil
		}).
		WithSetter("Set misc webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.Misc = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Sessions Webhook
	commands = append(commands, router.NewActionCommandBuilder("wsessions", "Sessions Webhook").
		WithGetter("Get sessions webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.Sessions, nil
		}).
		WithSetter("Set sessions webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.Sessions = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Ping Self Webhook Mode
	commands = append(commands, router.NewActionCommandBuilder("wpingself", "Ping Self Webhook").
		WithGetter("Get ping self webhook mode", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.PingSelfMode, nil
		}).
		WithSetter("Set ping self webhook mode", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.PingSelfMode = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available ping self webhook modes", func() []any {
			return []any{"Disabled", "Always", "Important"}
		}).
		Build(c.services.Router))

	// Ping Everyone Webhook Mode
	commands = append(commands, router.NewActionCommandBuilder("wpingeveryone", "Ping Everyone Webhook").
		WithGetter("Get ping everyone webhook mode", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.PingEveryoneMode, nil
		}).
		WithSetter("Set ping everyone webhook mode", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.PingEveryoneMode = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available ping everyone webhook modes", func() []any {
			return []any{"Disabled", "Always", "Important"}
		}).
		Build(c.services.Router))

	// Event settings
	// Ban Notifier
	commands = append(commands, router.NewActionCommandBuilder("bannotifier", "Ban Notifier").
		WithGetter("Get ban notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.BanNotifier, nil
		}).
		WithSetter("Set ban notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.BanNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle ban notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.BanNotifier = !cfg.Events.BanNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Staff Detection
	commands = append(commands, router.NewActionCommandBuilder("staffdetection", "Staff Detection").
		WithGetter("Get staff detection status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DiscordStaffDetection, nil
		}).
		WithSetter("Set staff detection status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DiscordStaffDetection = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle staff detection", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DiscordStaffDetection = !cfg.Events.DiscordStaffDetection
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Ping Detection
	commands = append(commands, router.NewActionCommandBuilder("pingdetection", "Ping Detection").
		WithGetter("Get ping detection status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.GhostPingDetection, nil
		}).
		WithSetter("Set ping detection status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.GhostPingDetection = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle ping detection", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.GhostPingDetection = !cfg.Events.GhostPingDetection
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Ticket Notifier
	commands = append(commands, router.NewActionCommandBuilder("ticketnotifier", "Ticket Notifier").
		WithGetter("Get ticket notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.TicketNotifier, nil
		}).
		WithSetter("Set ticket notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.TicketNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle ticket notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.TicketNotifier = !cfg.Events.TicketNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Session Connected
	commands = append(commands, router.NewActionCommandBuilder("sessionconnected", "Session Connected").
		WithGetter("Get session connected status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.SessionConnected, nil
		}).
		WithSetter("Set session connected status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionConnected = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle session connected", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionConnected = !cfg.Events.SessionConnected
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Session Disconnected
	commands = append(commands, router.NewActionCommandBuilder("sessiondisconnected", "Session Disconnected").
		WithGetter("Get session disconnected status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.SessionDisconnected, nil
		}).
		WithSetter("Set session disconnected status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionDisconnected = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle session disconnected", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionDisconnected = !cfg.Events.SessionDisconnected
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Selfbot Detection
	commands = append(commands, router.NewActionCommandBuilder("selfbotdetection", "Selfbot Detection").
		WithGetter("Get selfbot detection status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.SelfbotDetection, nil
		}).
		WithSetter("Set selfbot detection status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SelfbotDetection = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle selfbot detection", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SelfbotDetection = !cfg.Events.SelfbotDetection
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// DM Typing Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmtypingnotifier", "DM Typing Notifier").
		WithGetter("Get DM typing notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMTypingNotifier, nil
		}).
		WithSetter("Set DM typing notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMTypingNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM typing notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMTypingNotifier = !cfg.Events.DMTypingNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// DM Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmnotifier", "DM Notifier").
		WithGetter("Get DM notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMNotifier, nil
		}).
		WithSetter("Set DM notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMNotifier = !cfg.Events.DMNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// DM Edit Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmeditnotifier", "DM Edit Notifier").
		WithGetter("Get DM edit notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMEditNotifier, nil
		}).
		WithSetter("Set DM edit notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMEditNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM edit notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMEditNotifier = !cfg.Events.DMEditNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// DM Delete Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmdeletenotifier", "DM Delete Notifier").
		WithGetter("Get DM delete notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMDeleteNotifier, nil
		}).
		WithSetter("Set DM delete notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMDeleteNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM delete notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMDeleteNotifier = !cfg.Events.DMDeleteNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Message Delete Notifier
	commands = append(commands, router.NewActionCommandBuilder("messagedeletenotifier", "Message Delete Notifier").
		WithGetter("Get message delete notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.MessageDeleteNotifier, nil
		}).
		WithSetter("Set message delete notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.MessageDeleteNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle message delete notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.MessageDeleteNotifier = !cfg.Events.MessageDeleteNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// UI Theme
	commands = append(commands, router.NewActionCommandBuilder("uitheme", "UI Theme").
		WithGetter("Get UI theme", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.Theme, nil
		}).
		WithSetter("Set UI theme", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.Theme = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available UI themes", func() []any {
			return []any{"Dark", "Light"}
		}).
		Build(c.services.Router))

	// Hide to System Tray
	commands = append(commands, router.NewActionCommandBuilder("hidesystemtray", "Hide to System Tray").
		WithGetter("Get hide to system tray status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.HideToSystemTray, nil
		}).
		WithSetter("Set hide to system tray status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.HideToSystemTray = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle hide to system tray", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.HideToSystemTray = !cfg.UI.HideToSystemTray
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Run UI on Startup
	commands = append(commands, router.NewActionCommandBuilder("runuionstartup", "Run UI on Startup").
		WithGetter("Get run UI on startup status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.RunOnStartup, nil
		}).
		WithSetter("Set run UI on startup status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle run UI on startup", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	// Run Bot on Startup
	commands = append(commands, router.NewActionCommandBuilder("runbotonstartup", "Run Bot on Startup").
		WithGetter("Get run bot on startup status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.RunOnStartup, nil
		}).
		WithSetter("Set run bot on startup status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle run bot on startup", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services.Router))

	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Config",
		Description: "Change the bot's settings",
	}, commands)
}
