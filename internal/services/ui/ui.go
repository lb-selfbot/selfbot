package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/beeep"
	"github.com/common-nighthawk/go-figure"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/embed"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/fs"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/discord"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

// Service handles high-level logging and displaying information to the user.
type Service interface {
	// ClearConsole clears the console.
	ClearConsole()

	// ShowBanner displays the LightningBot banner.
	ShowBanner()

	// ShowNotification displays a notification to the user.
	ShowNotification(format string, e ...any) error

	// LogError logs an error message.
	LogError(format string, e ...any)

	// LogInfo logs an info message.
	LogInfo(format string, e ...any)
}

type uiService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
	discordSvc  discord.Service
}

// NewService creates a new UIService instance.
func NewService(eventBroker eventbroker.Service, configSvc config.Service, discordSvc discord.Service) Service {
	svc := &uiService{eventBroker: eventBroker, configSvc: configSvc, discordSvc: discordSvc}
	svc.subscribe()
	return svc
}

func (ui *uiService) subscribe() {
	ui.eventBroker.Subscribe(eventbroker.EventShowNotification, func(payload any) {
		text, ok := payload.(string)
		if !ok {
			return
		}

		_ = ui.ShowNotification("%s", text)
	})

	ui.eventBroker.Subscribe(eventbroker.EventLogError, func(payload any) {
		text, ok := payload.(string)
		if !ok {
			return
		}

		ui.LogError("%s", text)
	})

	ui.eventBroker.Subscribe(eventbroker.EventLogInfo, func(payload any) {
		text, ok := payload.(string)
		if !ok {
			return
		}

		ui.LogInfo("%s", text)
	})

	ui.eventBroker.Subscribe(eventbroker.EventLogFeature, ui.logFeature)
}

func (ui *uiService) ClearConsole() {
	ui.eventBroker.Publish(eventbroker.EventClearConsole, nil)
}

func (ui *uiService) ShowBanner() {
	ascii := figure.NewFigure("LightningBot", "stop", true).Slicify()
	centered := utils.String.CenterSlices(ascii)
	colored := ui.colored(centered, 0, 203, 255, true) + "\n\n"

	fmt.Println(colored)
}

func (ui *uiService) ShowNotification(format string, e ...any) error {
	cfg := ui.configSvc.GetConfig()
	if !cfg.Notifications.Enabled {
		return nil
	}

	go beeep.Notify(
		"LightningBot",
		fmt.Sprintf(format, e...),
		"LightningBot",
		filepath.Join(fs.GetMainFolder(), "Logo.png"),
	)

	return nil
}

func (ui *uiService) NewConsoleMessage() *ConsoleMessage {
	return NewConsoleMessage(ui.ShowConsoleMessage)
}

func (ui *uiService) ShowConsoleMessage(c *ConsoleMessage) error {
	r, g, b, _ := utils.Misc.HexToRGB(c.Color)

	timeFormat := ui.colored(ui.formattedTime(), 193, 239, 255, false)
	message := ui.colored("[", 238, 238, 238, false) + timeFormat + ui.colored("] ", 238, 238, 238, false) +
		ui.colored(c.Title, r, g, b, false) + "\n"

	for _, info := range c.Info {
		if info.Key == "" {
			message += ui.colored("  "+info.Value, 193, 239, 255, false) + "\n"
			continue
		}

		message += ui.colored("  "+info.Key+": ", 238, 238, 238, false) + ui.colored(info.Value, 193, 239, 255, false) + "\n"
	}

	if constants.CONSOLE {
		fmt.Println(message)
	}

	ui.eventBroker.Publish(eventbroker.EventAppendConsole, message)

	if c.WebhookURL != "" {
		emb := embed.NewEmbed(nil)

		for _, info := range c.Info {
			if info.Key == "" {
				emb.AddField("Extra Info", info.Value, false)
				continue
			}

			emb.AddField(info.Key, info.Value, false)
		}

		cfg := ui.configSvc.GetConfig()

		session := ui.discordSvc.GetSession()
		userID := session.State.User.ID

		var contentBuilder strings.Builder

		switch cfg.Webhooks.PingSelfMode {
		case "Disabled":
			break
		case "Always":
			contentBuilder.WriteString("<@!" + userID + ">\n")
		case "Important":
			if c.Important {
				contentBuilder.WriteString("<@!" + userID + ">\n")
			}
		default:
			cfg.Webhooks.PingSelfMode = "Disabled"
			_ = ui.configSvc.SaveConfig(cfg)
		}

		switch cfg.Webhooks.PingEveryoneMode {
		case "Disabled":
			break
		case "Always":
			contentBuilder.WriteString("@everyone\n")
		case "Important":
			if c.Important {
				contentBuilder.WriteString("@everyone\n")
			}
		default:
			cfg.Webhooks.PingEveryoneMode = "Disabled"
			_ = ui.configSvc.SaveConfig(cfg)
		}

		_, _ = session.Request("POST", c.WebhookURL, map[string]any{
			"content":    contentBuilder.String(),
			"embeds":     []*discordgo.MessageEmbed{emb.ToEmbed()},
			"username":   "LightningBot",
			"avatar_url": "https://lightning-bot.com/logo.png",
		})
	}

	return nil
}

func (ui *uiService) formattedTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}

func (ui *uiService) LogError(format string, e ...any) {
	text := fmt.Sprintf(format, e...)

	// Console
	text = fmt.Sprintf("[%s] [ERROR] %s", ui.formattedTime(), text)

	ui.eventBroker.Publish(eventbroker.EventAppendConsole, ui.colored(text, 255, 100, 100, true))
}

func (ui *uiService) LogInfo(format string, e ...any) {
	text := fmt.Sprintf(format, e...)

	// Console
	text = fmt.Sprintf("[%s] [INFO] %s", ui.formattedTime(), text)

	ui.eventBroker.Publish(eventbroker.EventAppendConsole, ui.colored(text, 193, 239, 255, true))
}

// colored formats the text with ANSI color codes.
func (ui *uiService) colored(text string, r, g, b int, bold bool) string {
	return "\033[38;2;" +
		strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" +
		strconv.Itoa(b) +
		fn.IfElse(bold, ";1", "") +
		"m" +
		text +
		"\033[0m"
}
