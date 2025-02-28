package config

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

type ConfigCommands struct{ services *services.Services }

func NewConfigCommands(s *services.Services) *ConfigCommands {
	return &ConfigCommands{services: s}
}

func (c *ConfigCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	commands := make([]*router.Command, 0)

	// Command Prefix
	commands = append(commands, router.NewActionCommandBuilder("commandprefix", "Command Prefix").
		WithGetter("Get the current command prefix", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.CommandPrefix, nil
		}).
		WithSetter("Set the command prefix", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.CommandPrefix = value.(string)
			return cfg.Save()
		}).
		Build())

	// Command Mode
	commands = append(commands, router.NewActionCommandBuilder("commandmode", "Command Mode").
		WithGetter("Get the current command mode", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.CommandMode, nil
		}).
		WithSetter("Set the command mode", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.CommandMode = value.(string)
			return cfg.Save()
		}).
		WithValues("List available command modes", func() []any {
			return []any{"Codeblock", "Text", "Indent", "WebhookEmbed", "URLEmbed"}
		}).
		Build())

	// Embed Title
	commands = append(commands, router.NewActionCommandBuilder("title", "Embed Title").
		WithGetter("Get the current embed title", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.Title, nil
		}).
		WithSetter("Set the embed title", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.Title = value.(string)
			return cfg.Save()
		}).
		Build())

	// Embed Footer
	commands = append(commands, router.NewActionCommandBuilder("footer", "Embed Footer").
		WithGetter("Get the current embed footer", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.Footer, nil
		}).
		WithSetter("Set the embed footer", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.Footer = value.(string)
			return cfg.Save()
		}).
		Build())

	// Delete Responses
	commands = append(commands, router.NewActionCommandBuilder("delresponses", "Delete Responses").
		WithGetter("Get delete responses status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.DeleteResponses, nil
		}).
		WithSetter("Set delete responses status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.DeleteResponses = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle delete responses", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.DeleteResponses = !cfg.Selfbot.DeleteResponses
			return cfg.Save()
		}).
		Build())

	// Delete Responses Delay
	commands = append(commands, router.NewActionCommandBuilder("delresponsesdelay", "Delete Responses Delay").
		WithGetter("Get delete responses delay", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.DeleteResponsesDelay, nil
		}).
		WithSetter("Set delete responses delay", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.DeleteResponsesDelay = value.(int)
			return cfg.Save()
		}).
		Build())

	// Delete Commands
	commands = append(commands, router.NewActionCommandBuilder("delcommands", "Delete Commands").
		WithGetter("Get delete commands status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.DeleteCommands, nil
		}).
		WithSetter("Set delete commands status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.DeleteCommands = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle delete commands", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.DeleteCommands = !cfg.Selfbot.DeleteCommands
			return cfg.Save()
		}).
		Build())

	// Client Spoof
	commands = append(commands, router.NewActionCommandBuilder("clientspoof", "Client Spoof").
		WithGetter("Get current client spoof setting", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.Client, nil
		}).
		WithSetter("Set client spoof setting", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.Client = utils.String.Capitalize(value.(string))
			return cfg.Save()
		}).
		WithValues("List available clients", func() []any {
			return []any{"Mobile", "Web", "Desktop", "Console"}
		}).
		Build())

	// Timezone
	commands = append(commands, router.NewActionCommandBuilder("timezone", "Timezone").
		WithGetter("Get current timezone", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.Timezone, nil
		}).
		WithSetter("Set timezone", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			timezone := value.(string)
			if _, err := time.LoadLocation(timezone); err != nil {
				return fmt.Errorf("invalid timezone: %v", err)
			}

			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.Timezone = timezone
			return cfg.Save()
		}).
		Build())

	// Embed Color
	commands = append(commands, router.NewActionCommandBuilder("embedcolor", "Embed Color").
		WithGetter("Get current embed color", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.EmbedColor, nil
		}).
		WithSetter("Set embed color", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.EmbedColor = value.(string)
			return cfg.Save()
		}).
		Build())

	// Embed Image
	commands = append(commands, router.NewActionCommandBuilder("embedimage", "Embed Image").
		WithGetter("Get current embed image", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.EmbedImage, nil
		}).
		WithSetter("Set embed image", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.EmbedImage = value.(string)
			return cfg.Save()
		}).
		Build())

	// Embed URL
	commands = append(commands, router.NewActionCommandBuilder("embedurl", "Embed URL").
		WithGetter("Get current embed URL", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Selfbot.EmbedURL, nil
		}).
		WithSetter("Set embed URL", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Selfbot.EmbedURL = value.(string)
			return cfg.Save()
		}).
		Build())

	// Notifications
	commands = append(commands, router.NewActionCommandBuilder("notifications", "Notifications").
		WithGetter("Get notifications status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Notifications.Enabled, nil
		}).
		WithSetter("Set notifications status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Notifications.Enabled = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle notifications", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Notifications.Enabled = !cfg.Notifications.Enabled
			return cfg.Save()
		}).
		WithExample("Show an example of a notification", func(ctx *router.Ctx) error {
			c.services.UI.ShowNotification("This is an example notification.")
			return nil
		}).
		Build())

	// Nitro Sniper
	commands = append(commands, router.NewActionCommandBuilder("nitrosniper", "Nitro Sniper").
		WithGetter("Get nitro sniper status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.NitroSniper, nil
		}).
		WithSetter("Set nitro sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.NitroSniper = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle nitro sniper", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.NitroSniper = !cfg.Snipers.NitroSniper
			return cfg.Save()
		}).
		Build())

	// Giveaway Sniper
	commands = append(commands, router.NewActionCommandBuilder("giveawaysniper", "Giveaway Sniper").
		WithGetter("Get giveaway sniper status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.GiveawaySniper.Enabled, nil
		}).
		WithSetter("Set giveaway sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.GiveawaySniper.Enabled = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle giveaway sniper", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.GiveawaySniper.Enabled = !cfg.Snipers.GiveawaySniper.Enabled
			return cfg.Save()
		}).
		Build())

	// Giveaway Sniper Sleep
	commands = append(commands, router.NewActionCommandBuilder("giveawaysnipersleep", "Giveaway Sniper Sleep").
		WithGetter("Get giveaway sniper sleep duration", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.GiveawaySniper.Sleep, nil
		}).
		WithSetter("Set giveaway sniper sleep duration", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.GiveawaySniper.Sleep = value.(int)
			return cfg.Save()
		}).
		Build())

	// Multi-Token Sniper
	commands = append(commands, router.NewActionCommandBuilder("multitokensniper", "Multi-Token Sniper").
		WithGetter("Get multi-token sniper status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.MultiTokenSniper, nil
		}).
		WithSetter("Set multi-token sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.MultiTokenSniper = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle multi-token sniper", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.MultiTokenSniper = !cfg.Snipers.MultiTokenSniper
			return cfg.Save()
		}).
		Build())

	// Privnote Sniper
	commands = append(commands, router.NewActionCommandBuilder("privnotesniper", "Privnote Sniper").
		WithGetter("Get privnote sniper status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PrivnoteSniper, nil
		}).
		WithSetter("Set privnote sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PrivnoteSniper = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle privnote sniper", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PrivnoteSniper = !cfg.Snipers.PrivnoteSniper
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper
	commands = append(commands, router.NewActionCommandBuilder("poketwosniper", "Poketwo Sniper").
		WithGetter("Get poketwo sniper status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.Enabled, nil
		}).
		WithSetter("Set poketwo sniper status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.Enabled = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle poketwo sniper", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.Enabled = !cfg.Snipers.PokeTwoSniper.Enabled
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Sleep
	commands = append(commands, router.NewActionCommandBuilder("poketwosleep", "Poketwo Sniper Sleep").
		WithGetter("Get poketwo sniper sleep duration", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.Sleep, nil
		}).
		WithSetter("Set poketwo sniper sleep duration", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.Sleep = value.(int)
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Random Sleep
	commands = append(commands, router.NewActionCommandBuilder("poketworandomsleep", "Poketwo Sniper Random Sleep").
		WithGetter("Get poketwo sniper random sleep status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.RandomSleep, nil
		}).
		WithSetter("Set poketwo sniper random sleep status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.RandomSleep = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle poketwo sniper random sleep", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.RandomSleep = !cfg.Snipers.PokeTwoSniper.RandomSleep
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Timeout
	commands = append(commands, router.NewActionCommandBuilder("poketwotimeout", "Poketwo Sniper Timeout").
		WithGetter("Get poketwo sniper timeout duration", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.Timeout, nil
		}).
		WithSetter("Set poketwo sniper timeout duration", router.ValueTypeDuration, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.Timeout = value.(int)
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Typing
	commands = append(commands, router.NewActionCommandBuilder("poketwotyping", "Poketwo Sniper Typing").
		WithGetter("Get poketwo sniper typing status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.Typing, nil
		}).
		WithSetter("Set poketwo sniper typing status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.Typing = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle poketwo sniper typing", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.Typing = !cfg.Snipers.PokeTwoSniper.Typing
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Catch Rate
	commands = append(commands, router.NewActionCommandBuilder("poketwocatchrate", "Poketwo Sniper Catch Rate").
		WithGetter("Get poketwo sniper catch rate", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.CatchRate, nil
		}).
		WithSetter("Set poketwo sniper catch rate", router.ValueTypeInt, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.CatchRate = value.(int)
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Fake Activity
	commands = append(commands, router.NewActionCommandBuilder("poketwofakeactivity", "Poketwo Sniper Fake Activity").
		WithGetter("Get poketwo sniper fake activity status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.FakeActivity, nil
		}).
		WithSetter("Set poketwo sniper fake activity status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.FakeActivity = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle poketwo sniper fake activity", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.FakeActivity = !cfg.Snipers.PokeTwoSniper.FakeActivity
			return cfg.Save()
		}).
		Build())

	// Poketwo Sniper Language
	commands = append(commands, router.NewActionCommandBuilder("poketwolanguage", "Poketwo Sniper Language").
		WithGetter("Get poketwo sniper language", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Snipers.PokeTwoSniper.CatchLanguage, nil
		}).
		WithSetter("Set poketwo sniper language", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Snipers.PokeTwoSniper.CatchLanguage = value.(string)
			return cfg.Save()
		}).
		WithValues("List available languages", func() []any {
			return []any{"en", "ja", "ja_t", "fr", "de"}
		}).
		Build())

	// Webhook settings
	// Nitro Sniper Webhook
	commands = append(commands, router.NewActionCommandBuilder("wnitro", "Nitro Sniper Webhook").
		WithGetter("Get nitro sniper webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.NitroSniper, nil
		}).
		WithSetter("Set nitro sniper webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.NitroSniper = value.(string)
			return cfg.Save()
		}).
		Build())

	// Giveaway Sniper Webhook
	commands = append(commands, router.NewActionCommandBuilder("wgiveawaysniper", "Giveaway Sniper Webhook").
		WithGetter("Get giveaway sniper webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.GiveawaySniper, nil
		}).
		WithSetter("Set giveaway sniper webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.GiveawaySniper = value.(string)
			return cfg.Save()
		}).
		Build())

	// Selfbot Detection Webhook
	commands = append(commands, router.NewActionCommandBuilder("wselfbotdetection", "Selfbot Detection Webhook").
		WithGetter("Get selfbot detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.SelfbotDetection, nil
		}).
		WithSetter("Set selfbot detection webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.SelfbotDetection = value.(string)
			return cfg.Save()
		}).
		Build())

	// Staff Detection Webhook
	commands = append(commands, router.NewActionCommandBuilder("wstaffdetection", "Staff Detection Webhook").
		WithGetter("Get staff detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.StaffDetection, nil
		}).
		WithSetter("Set staff detection webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.StaffDetection = value.(string)
			return cfg.Save()
		}).
		Build())

	// Notifiers Webhook
	commands = append(commands, router.NewActionCommandBuilder("wnotifiers", "Notifiers Webhook").
		WithGetter("Get notifiers webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.Notifiers, nil
		}).
		WithSetter("Set notifiers webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.Notifiers = value.(string)
			return cfg.Save()
		}).
		Build())

	// Ping Detection Webhook
	commands = append(commands, router.NewActionCommandBuilder("wpingdetection", "Ping Detection Webhook").
		WithGetter("Get ping detection webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.PingDetection, nil
		}).
		WithSetter("Set ping detection webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.PingDetection = value.(string)
			return cfg.Save()
		}).
		Build())

	// Misc Webhook
	commands = append(commands, router.NewActionCommandBuilder("wmisc", "Misc Webhook").
		WithGetter("Get misc webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.Misc, nil
		}).
		WithSetter("Set misc webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.Misc = value.(string)
			return cfg.Save()
		}).
		Build())

	// Sessions Webhook
	commands = append(commands, router.NewActionCommandBuilder("wsessions", "Sessions Webhook").
		WithGetter("Get sessions webhook", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.Sessions, nil
		}).
		WithSetter("Set sessions webhook", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.Sessions = value.(string)
			return cfg.Save()
		}).
		Build())

	// Ping Self Webhook Mode
	commands = append(commands, router.NewActionCommandBuilder("wpingself", "Ping Self Webhook").
		WithGetter("Get ping self webhook mode", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.PingSelfMode, nil
		}).
		WithSetter("Set ping self webhook mode", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.PingSelfMode = value.(string)
			return cfg.Save()
		}).
		WithValues("List available ping self webhook modes", func() []any {
			return []any{"Disabled", "Always", "Important"}
		}).
		Build())

	// Ping Everyone Webhook Mode
	commands = append(commands, router.NewActionCommandBuilder("wpingeveryone", "Ping Everyone Webhook").
		WithGetter("Get ping everyone webhook mode", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Webhooks.PingEveryoneMode, nil
		}).
		WithSetter("Set ping everyone webhook mode", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Webhooks.PingEveryoneMode = value.(string)
			return cfg.Save()
		}).
		WithValues("List available ping everyone webhook modes", func() []any {
			return []any{"Disabled", "Always", "Important"}
		}).
		Build())

	// Event settings
	// Ban Notifier
	commands = append(commands, router.NewActionCommandBuilder("bannotifier", "Ban Notifier").
		WithGetter("Get ban notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.BanNotifier, nil
		}).
		WithSetter("Set ban notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.BanNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle ban notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.BanNotifier = !cfg.Events.BanNotifier
			return cfg.Save()
		}).
		Build())

	// Staff Detection
	commands = append(commands, router.NewActionCommandBuilder("staffdetection", "Staff Detection").
		WithGetter("Get staff detection status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.DiscordStaffDetection, nil
		}).
		WithSetter("Set staff detection status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DiscordStaffDetection = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle staff detection", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DiscordStaffDetection = !cfg.Events.DiscordStaffDetection
			return cfg.Save()
		}).
		Build())

	// Ping Detection
	commands = append(commands, router.NewActionCommandBuilder("pingdetection", "Ping Detection").
		WithGetter("Get ping detection status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.GhostPingDetection, nil
		}).
		WithSetter("Set ping detection status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.GhostPingDetection = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle ping detection", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.GhostPingDetection = !cfg.Events.GhostPingDetection
			return cfg.Save()
		}).
		Build())

	// Ticket Notifier
	commands = append(commands, router.NewActionCommandBuilder("ticketnotifier", "Ticket Notifier").
		WithGetter("Get ticket notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.TicketNotifier, nil
		}).
		WithSetter("Set ticket notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.TicketNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle ticket notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.TicketNotifier = !cfg.Events.TicketNotifier
			return cfg.Save()
		}).
		Build())

	// Session Connected
	commands = append(commands, router.NewActionCommandBuilder("sessionconnected", "Session Connected").
		WithGetter("Get session connected status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.SessionConnected, nil
		}).
		WithSetter("Set session connected status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.SessionConnected = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle session connected", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.SessionConnected = !cfg.Events.SessionConnected
			return cfg.Save()
		}).
		Build())

	// Session Disconnected
	commands = append(commands, router.NewActionCommandBuilder("sessiondisconnected", "Session Disconnected").
		WithGetter("Get session disconnected status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.SessionDisconnected, nil
		}).
		WithSetter("Set session disconnected status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.SessionDisconnected = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle session disconnected", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.SessionDisconnected = !cfg.Events.SessionDisconnected
			return cfg.Save()
		}).
		Build())

	// Selfbot Detection
	commands = append(commands, router.NewActionCommandBuilder("selfbotdetection", "Selfbot Detection").
		WithGetter("Get selfbot detection status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.SelfbotDetection, nil
		}).
		WithSetter("Set selfbot detection status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.SelfbotDetection = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle selfbot detection", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.SelfbotDetection = !cfg.Events.SelfbotDetection
			return cfg.Save()
		}).
		Build())

	// DM Typing Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmtypingnotifier", "DM Typing Notifier").
		WithGetter("Get DM typing notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.DMTypingNotifier, nil
		}).
		WithSetter("Set DM typing notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMTypingNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle DM typing notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMTypingNotifier = !cfg.Events.DMTypingNotifier
			return cfg.Save()
		}).
		Build())

	// DM Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmnotifier", "DM Notifier").
		WithGetter("Get DM notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.DMNotifier, nil
		}).
		WithSetter("Set DM notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle DM notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMNotifier = !cfg.Events.DMNotifier
			return cfg.Save()
		}).
		Build())

	// DM Edit Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmeditnotifier", "DM Edit Notifier").
		WithGetter("Get DM edit notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.DMEditNotifier, nil
		}).
		WithSetter("Set DM edit notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMEditNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle DM edit notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMEditNotifier = !cfg.Events.DMEditNotifier
			return cfg.Save()
		}).
		Build())

	// DM Delete Notifier
	commands = append(commands, router.NewActionCommandBuilder("dmdeletenotifier", "DM Delete Notifier").
		WithGetter("Get DM delete notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.DMDeleteNotifier, nil
		}).
		WithSetter("Set DM delete notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMDeleteNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle DM delete notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.DMDeleteNotifier = !cfg.Events.DMDeleteNotifier
			return cfg.Save()
		}).
		Build())

	// Message Delete Notifier
	commands = append(commands, router.NewActionCommandBuilder("messagedeletenotifier", "Message Delete Notifier").
		WithGetter("Get message delete notifier status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.Events.MessageDeleteNotifier, nil
		}).
		WithSetter("Set message delete notifier status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.MessageDeleteNotifier = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle message delete notifier", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.Events.MessageDeleteNotifier = !cfg.Events.MessageDeleteNotifier
			return cfg.Save()
		}).
		Build())

	// UI Theme
	commands = append(commands, router.NewActionCommandBuilder("uitheme", "UI Theme").
		WithGetter("Get UI theme", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.UI.Theme, nil
		}).
		WithSetter("Set UI theme", router.ValueTypeString, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.Theme = value.(string)
			return cfg.Save()
		}).
		WithValues("List available UI themes", func() []any {
			return []any{"Dark", "Light"}
		}).
		Build())

	// Hide to System Tray
	commands = append(commands, router.NewActionCommandBuilder("hidesystemtray", "Hide to System Tray").
		WithGetter("Get hide to system tray status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.UI.HideToSystemTray, nil
		}).
		WithSetter("Set hide to system tray status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.HideToSystemTray = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle hide to system tray", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.HideToSystemTray = !cfg.UI.HideToSystemTray
			return cfg.Save()
		}).
		Build())

	// Run UI on Startup
	commands = append(commands, router.NewActionCommandBuilder("runuionstartup", "Run UI on Startup").
		WithGetter("Get run UI on startup status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.UI.RunOnStartup, nil
		}).
		WithSetter("Set run UI on startup status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.RunOnStartup = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle run UI on startup", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
			return cfg.Save()
		}).
		Build())

	// Run Bot on Startup
	commands = append(commands, router.NewActionCommandBuilder("runbotonstartup", "Run Bot on Startup").
		WithGetter("Get run bot on startup status", func(ctx *router.Ctx) (any, error) {
			cfg, err := config.GetConfig()
			if err != nil {
				return nil, err
			}
			return cfg.UI.RunOnStartup, nil
		}).
		WithSetter("Set run bot on startup status", router.ValueTypeBool, func(ctx *router.Ctx, value any) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.RunOnStartup = value.(bool)
			return cfg.Save()
		}).
		WithToggler("Toggle run bot on startup", func(ctx *router.Ctx) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			cfg.UI.RunOnStartup = !cfg.UI.RunOnStartup
			return cfg.Save()
		}).
		Build())

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Config",
		Description: "Change the bot's settings",
	}, commands)
}
