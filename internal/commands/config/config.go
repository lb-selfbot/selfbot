package config

import (
	"fmt"
	"time"

	"github.com/lb-selfbot/selfbot/internal/actioncommands"
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
	commands = append(commands, actioncommands.NewBuilder("commandprefix", "Command Prefix").
		WithGetter("Get the current command prefix", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.CommandPrefix, nil
		}).
		WithSetter("Set the command prefix", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.CommandPrefix = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Command Mode
	commands = append(commands, actioncommands.NewBuilder("commandmode", "Command Mode").
		WithGetter("Get the current command mode", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.CommandMode, nil
		}).
		WithSetter("Set the command mode", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.CommandMode = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available command modes", func() []any {
			return []any{"Codeblock", "Text", "Indent", "WebhookEmbed", "URLEmbed"}
		}).
		Build(c.services))

	// Embed Title
	commands = append(commands, actioncommands.NewBuilder("title", "Embed Title").
		WithGetter("Get the current embed title", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Title, nil
		}).
		WithSetter("Set the embed title", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Title = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Embed Footer
	commands = append(commands, actioncommands.NewBuilder("footer", "Embed Footer").
		WithGetter("Get the current embed footer", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Footer, nil
		}).
		WithSetter("Set the embed footer", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Footer = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Delete Responses
	commands = append(commands, actioncommands.NewBuilder("delresponses", "Delete Responses").
		WithGetter("Get delete responses status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.DeleteResponses, nil
		}).
		WithSetter("Set delete responses status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteResponses = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle delete responses", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Delete Responses Delay
	commands = append(commands, actioncommands.NewBuilder("delresponsesdelay", "Delete Responses Delay").
		WithGetter("Get delete responses delay", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.DeleteResponsesDelay, nil
		}).
		WithSetter("Set delete responses delay", actioncommands.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteResponsesDelay = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Delete Commands
	commands = append(commands, actioncommands.NewBuilder("delcommands", "Delete Commands").
		WithGetter("Get delete commands status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.DeleteCommands, nil
		}).
		WithSetter("Set delete commands status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteCommands = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle delete commands", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Client Spoof
	commands = append(commands, actioncommands.NewBuilder("clientspoof", "Client Spoof").
		WithGetter("Get current client spoof setting", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Client, nil
		}).
		WithSetter("Set client spoof setting", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Client = utils.String.Capitalize(value.(string))
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available clients", func() []any {
			return []any{"Mobile", "Web", "Desktop", "Console"}
		}).
		Build(c.services))

	// Timezone
	commands = append(commands, actioncommands.NewBuilder("timezone", "Timezone").
		WithGetter("Get current timezone", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.Timezone, nil
		}).
		WithSetter("Set timezone", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			timezone := value.(string)
			if _, err := time.LoadLocation(timezone); err != nil {
				return fmt.Errorf("invalid timezone: %v", err)
			}

			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.Timezone = timezone
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Embed Color
	commands = append(commands, actioncommands.NewBuilder("embedcolor", "Embed Color").
		WithGetter("Get current embed color", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.EmbedColor, nil
		}).
		WithSetter("Set embed color", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.EmbedColor = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Embed Image
	commands = append(commands, actioncommands.NewBuilder("embedimage", "Embed Image").
		WithGetter("Get current embed image", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.EmbedImage, nil
		}).
		WithSetter("Set embed image", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.EmbedImage = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Embed URL
	commands = append(commands, actioncommands.NewBuilder("embedurl", "Embed URL").
		WithGetter("Get current embed URL", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Selfbot.EmbedURL, nil
		}).
		WithSetter("Set embed URL", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Selfbot.EmbedURL = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Notifications
	commands = append(commands, actioncommands.NewBuilder("notifications", "Notifications").
		WithGetter("Get notifications status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Notifications.Enabled, nil
		}).
		WithSetter("Set notifications status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
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
			return c.services.UI.ShowNotification("This is an example notification.")
		}).
		Build(c.services))

	// Nitro Sniper
	commands = append(commands, actioncommands.NewBuilder("nitrosniper", "Nitro Sniper").
		WithGetter("Get nitro sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.NitroSniper, nil
		}).
		WithSetter("Set nitro sniper status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.NitroSniper = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle nitro sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Giveaway Sniper
	commands = append(commands, actioncommands.NewBuilder("giveawaysniper", "Giveaway Sniper").
		WithGetter("Get giveaway sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.GiveawaySniper.Enabled, nil
		}).
		WithSetter("Set giveaway sniper status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.GiveawaySniper.Enabled = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle giveaway sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Giveaway Sniper Sleep
	commands = append(commands, actioncommands.NewBuilder("giveawaysnipersleep", "Giveaway Sniper Sleep").
		WithGetter("Get giveaway sniper sleep duration", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.GiveawaySniper.Sleep, nil
		}).
		WithSetter("Set giveaway sniper sleep duration", actioncommands.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.GiveawaySniper.Sleep = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Multi-Token Sniper
	commands = append(commands, actioncommands.NewBuilder("multitokensniper", "Multi-Token Sniper").
		WithGetter("Get multi-token sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.MultiTokenSniper, nil
		}).
		WithSetter("Set multi-token sniper status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.MultiTokenSniper = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle multi-token sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.MultiTokenSniper = !cfg.Snipers.MultiTokenSniper
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Privnote Sniper
	commands = append(commands, actioncommands.NewBuilder("privnotesniper", "Privnote Sniper").
		WithGetter("Get privnote sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PrivnoteSniper, nil
		}).
		WithSetter("Set privnote sniper status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PrivnoteSniper = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle privnote sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper
	commands = append(commands, actioncommands.NewBuilder("poketwosniper", "Poketwo Sniper").
		WithGetter("Get poketwo sniper status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Enabled, nil
		}).
		WithSetter("Set poketwo sniper status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Enabled = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Enabled = !cfg.Snipers.PokeTwoSniper.Enabled
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Sleep
	commands = append(commands, actioncommands.NewBuilder("poketwosleep", "Poketwo Sniper Sleep").
		WithGetter("Get poketwo sniper sleep duration", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Sleep, nil
		}).
		WithSetter("Set poketwo sniper sleep duration", actioncommands.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Sleep = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Random Sleep
	commands = append(commands, actioncommands.NewBuilder("poketworandomsleep", "Poketwo Sniper Random Sleep").
		WithGetter("Get poketwo sniper random sleep status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.RandomSleep, nil
		}).
		WithSetter("Set poketwo sniper random sleep status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.RandomSleep = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper random sleep", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.RandomSleep = !cfg.Snipers.PokeTwoSniper.RandomSleep
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Timeout
	commands = append(commands, actioncommands.NewBuilder("poketwotimeout", "Poketwo Sniper Timeout").
		WithGetter("Get poketwo sniper timeout duration", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Timeout, nil
		}).
		WithSetter("Set poketwo sniper timeout duration", actioncommands.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Timeout = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Typing
	commands = append(commands, actioncommands.NewBuilder("poketwotyping", "Poketwo Sniper Typing").
		WithGetter("Get poketwo sniper typing status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.Typing, nil
		}).
		WithSetter("Set poketwo sniper typing status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Typing = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper typing", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.Typing = !cfg.Snipers.PokeTwoSniper.Typing
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Catch Rate
	commands = append(commands, actioncommands.NewBuilder("poketwocatchrate", "Poketwo Sniper Catch Rate").
		WithGetter("Get poketwo sniper catch rate", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.CatchRate, nil
		}).
		WithSetter("Set poketwo sniper catch rate", actioncommands.ValueTypeInt, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.CatchRate = value.(int)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Fake Activity
	commands = append(commands, actioncommands.NewBuilder("poketwofakeactivity", "Poketwo Sniper Fake Activity").
		WithGetter("Get poketwo sniper fake activity status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.FakeActivity, nil
		}).
		WithSetter("Set poketwo sniper fake activity status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.FakeActivity = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle poketwo sniper fake activity", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.FakeActivity = !cfg.Snipers.PokeTwoSniper.FakeActivity
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Poketwo Sniper Language
	commands = append(commands, actioncommands.NewBuilder("poketwolanguage", "Poketwo Sniper Language").
		WithGetter("Get poketwo sniper language", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Snipers.PokeTwoSniper.CatchLanguage, nil
		}).
		WithSetter("Set poketwo sniper language", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Snipers.PokeTwoSniper.CatchLanguage = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available languages", func() []any {
			return []any{"en", "ja", "ja_t", "fr", "de"}
		}).
		Build(c.services))

	// Webhook settings
	// Nitro Sniper Webhook
	commands = append(commands, actioncommands.NewBuilder("wnitro", "Nitro Sniper Webhook").
		WithGetter("Get nitro sniper webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.NitroSniper, nil
		}).
		WithSetter("Set nitro sniper webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.NitroSniper = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Giveaway Sniper Webhook
	commands = append(commands, actioncommands.NewBuilder("wgiveawaysniper", "Giveaway Sniper Webhook").
		WithGetter("Get giveaway sniper webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.GiveawaySniper, nil
		}).
		WithSetter("Set giveaway sniper webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.GiveawaySniper = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Selfbot Detection Webhook
	commands = append(commands, actioncommands.NewBuilder("wselfbotdetection", "Selfbot Detection Webhook").
		WithGetter("Get selfbot detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.SelfbotDetection, nil
		}).
		WithSetter("Set selfbot detection webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.SelfbotDetection = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Staff Detection Webhook
	commands = append(commands, actioncommands.NewBuilder("wstaffdetection", "Staff Detection Webhook").
		WithGetter("Get staff detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.StaffDetection, nil
		}).
		WithSetter("Set staff detection webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.StaffDetection = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Notifiers Webhook
	commands = append(commands, actioncommands.NewBuilder("wnotifiers", "Notifiers Webhook").
		WithGetter("Get notifiers webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.Notifiers, nil
		}).
		WithSetter("Set notifiers webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.Notifiers = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Ping Detection Webhook
	commands = append(commands, actioncommands.NewBuilder("wpingdetection", "Ping Detection Webhook").
		WithGetter("Get ping detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.PingDetection, nil
		}).
		WithSetter("Set ping detection webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.PingDetection = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Misc Webhook
	commands = append(commands, actioncommands.NewBuilder("wmisc", "Misc Webhook").
		WithGetter("Get misc webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.Misc, nil
		}).
		WithSetter("Set misc webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.Misc = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Sessions Webhook
	commands = append(commands, actioncommands.NewBuilder("wsessions", "Sessions Webhook").
		WithGetter("Get sessions webhook", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.Sessions, nil
		}).
		WithSetter("Set sessions webhook", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.Sessions = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Ping Self Webhook Mode
	commands = append(commands, actioncommands.NewBuilder("wpingself", "Ping Self Webhook").
		WithGetter("Get ping self webhook mode", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.PingSelfMode, nil
		}).
		WithSetter("Set ping self webhook mode", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.PingSelfMode = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available ping self webhook modes", func() []any {
			return []any{"Disabled", "Always", "Important"}
		}).
		Build(c.services))

	// Ping Everyone Webhook Mode
	commands = append(commands, actioncommands.NewBuilder("wpingeveryone", "Ping Everyone Webhook").
		WithGetter("Get ping everyone webhook mode", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Webhooks.PingEveryoneMode, nil
		}).
		WithSetter("Set ping everyone webhook mode", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Webhooks.PingEveryoneMode = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available ping everyone webhook modes", func() []any {
			return []any{"Disabled", "Always", "Important"}
		}).
		Build(c.services))

	// Event settings
	// Ban Notifier
	commands = append(commands, actioncommands.NewBuilder("bannotifier", "Ban Notifier").
		WithGetter("Get ban notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.BanNotifier, nil
		}).
		WithSetter("Set ban notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.BanNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle ban notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.BanNotifier = !cfg.Events.BanNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Staff Detection
	commands = append(commands, actioncommands.NewBuilder("staffdetection", "Staff Detection").
		WithGetter("Get staff detection status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DiscordStaffDetection, nil
		}).
		WithSetter("Set staff detection status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DiscordStaffDetection = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle staff detection", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DiscordStaffDetection = !cfg.Events.DiscordStaffDetection
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Ping Detection
	commands = append(commands, actioncommands.NewBuilder("pingdetection", "Ping Detection").
		WithGetter("Get ping detection status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.GhostPingDetection, nil
		}).
		WithSetter("Set ping detection status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.GhostPingDetection = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle ping detection", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.GhostPingDetection = !cfg.Events.GhostPingDetection
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Ticket Notifier
	commands = append(commands, actioncommands.NewBuilder("ticketnotifier", "Ticket Notifier").
		WithGetter("Get ticket notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.TicketNotifier, nil
		}).
		WithSetter("Set ticket notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.TicketNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle ticket notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.TicketNotifier = !cfg.Events.TicketNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Session Connected
	commands = append(commands, actioncommands.NewBuilder("sessionconnected", "Session Connected").
		WithGetter("Get session connected status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.SessionConnected, nil
		}).
		WithSetter("Set session connected status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionConnected = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle session connected", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionConnected = !cfg.Events.SessionConnected
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Session Disconnected
	commands = append(commands, actioncommands.NewBuilder("sessiondisconnected", "Session Disconnected").
		WithGetter("Get session disconnected status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.SessionDisconnected, nil
		}).
		WithSetter("Set session disconnected status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionDisconnected = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle session disconnected", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SessionDisconnected = !cfg.Events.SessionDisconnected
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Selfbot Detection
	commands = append(commands, actioncommands.NewBuilder("selfbotdetection", "Selfbot Detection").
		WithGetter("Get selfbot detection status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.SelfbotDetection, nil
		}).
		WithSetter("Set selfbot detection status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SelfbotDetection = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle selfbot detection", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.SelfbotDetection = !cfg.Events.SelfbotDetection
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// DM Typing Notifier
	commands = append(commands, actioncommands.NewBuilder("dmtypingnotifier", "DM Typing Notifier").
		WithGetter("Get DM typing notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMTypingNotifier, nil
		}).
		WithSetter("Set DM typing notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMTypingNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM typing notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMTypingNotifier = !cfg.Events.DMTypingNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// DM Notifier
	commands = append(commands, actioncommands.NewBuilder("dmnotifier", "DM Notifier").
		WithGetter("Get DM notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMNotifier, nil
		}).
		WithSetter("Set DM notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMNotifier = !cfg.Events.DMNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// DM Edit Notifier
	commands = append(commands, actioncommands.NewBuilder("dmeditnotifier", "DM Edit Notifier").
		WithGetter("Get DM edit notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMEditNotifier, nil
		}).
		WithSetter("Set DM edit notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMEditNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM edit notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMEditNotifier = !cfg.Events.DMEditNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// DM Delete Notifier
	commands = append(commands, actioncommands.NewBuilder("dmdeletenotifier", "DM Delete Notifier").
		WithGetter("Get DM delete notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.DMDeleteNotifier, nil
		}).
		WithSetter("Set DM delete notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMDeleteNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle DM delete notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.DMDeleteNotifier = !cfg.Events.DMDeleteNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Message Delete Notifier
	commands = append(commands, actioncommands.NewBuilder("messagedeletenotifier", "Message Delete Notifier").
		WithGetter("Get message delete notifier status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.Events.MessageDeleteNotifier, nil
		}).
		WithSetter("Set message delete notifier status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.MessageDeleteNotifier = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle message delete notifier", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.Events.MessageDeleteNotifier = !cfg.Events.MessageDeleteNotifier
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// UI Theme
	commands = append(commands, actioncommands.NewBuilder("uitheme", "UI Theme").
		WithGetter("Get UI theme", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.Theme, nil
		}).
		WithSetter("Set UI theme", actioncommands.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.Theme = value.(string)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithValues("List available UI themes", func() []any {
			return []any{"Dark", "Light"}
		}).
		Build(c.services))

	// Hide to System Tray
	commands = append(commands, actioncommands.NewBuilder("hidesystemtray", "Hide to System Tray").
		WithGetter("Get hide to system tray status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.HideToSystemTray, nil
		}).
		WithSetter("Set hide to system tray status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.HideToSystemTray = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle hide to system tray", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.HideToSystemTray = !cfg.UI.HideToSystemTray
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Run UI on Startup
	commands = append(commands, actioncommands.NewBuilder("runuionstartup", "Run UI on Startup").
		WithGetter("Get run UI on startup status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.RunOnStartup, nil
		}).
		WithSetter("Set run UI on startup status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle run UI on startup", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	// Run Bot on Startup
	commands = append(commands, actioncommands.NewBuilder("runbotonstartup", "Run Bot on Startup").
		WithGetter("Get run bot on startup status", func(ctx *router.Ctx) (any, error) {
			cfg := c.services.Config.GetConfig()
			return cfg.UI.RunOnStartup, nil
		}).
		WithSetter("Set run bot on startup status", actioncommands.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = value.(bool)
			return c.services.Config.SaveConfig(cfg)
		}).
		WithToggler("Toggle run bot on startup", func(ctx *router.Ctx) error {
			cfg := c.services.Config.GetConfig()
			cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
			return c.services.Config.SaveConfig(cfg)
		}).
		Build(c.services))

	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Config",
		Description: "Change the bot's settings",
	}, commands)
}
