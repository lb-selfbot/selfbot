package ui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/beeep"
	"github.com/common-nighthawk/go-figure"
	"github.com/robert-nix/ansihtml"
)

type UIServiceImpl struct {
	authSvc    interfaces.AuthService
	ipcSvc     interfaces.IPCService
	discordSvc interfaces.DiscordService
	embedSvc   interfaces.EmbedService
}

func (ui *UIServiceImpl) Inject(authSvc interfaces.AuthService, ipcSvc interfaces.IPCService, discordSvc interfaces.DiscordService, embedSvc interfaces.EmbedService) {
	ui.authSvc = authSvc
	ui.ipcSvc = ipcSvc
	ui.discordSvc = discordSvc
	ui.embedSvc = embedSvc
}

func (*UIServiceImpl) ANSIColored(text string, r, g, b int, bold bool) string {
	boldArg := ";1"
	if !bold {
		boldArg = ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%d%sm%s\033[0m", r, g, b, boldArg, text)
}

func (ui *UIServiceImpl) ShowBanner(flags int) error {
	if flags&BannerConsole != 0 {
		ascii := figure.NewFigure("LightningBot", "stop", true).Slicify()
		centered := utils.String.CenterSlices(ascii)
		colored := ui.ANSIColored(centered, 0, 203, 255, true) + "\n\n"

		fmt.Println(colored)
	}

	if flags&BannerUI != 0 {
		banner := `<image Source="assets/Banner.png" Width="650" Height="150" Centered />`
		return ui.SetConsole(banner)
	}

	return nil
}

func (ui *UIServiceImpl) ShowNotification(format string, e ...any) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	if !cfg.Notifications.Enabled {
		return nil
	}

	go beeep.Notify(
		"LightningBot",
		fmt.Sprintf(format, e...),
		"LightningBot",
		filepath.Join(utils.IO.GetMainFolder(), "Logo.png"),
	)

	return nil
}

func (ui *UIServiceImpl) formattedTime() string {
	return time.Now().Local().Format("2006-01-02 15:04:05")
}

func (ui *UIServiceImpl) LogError(format string, e ...any) {
	text := fmt.Sprintf(format, e...)

	// IPC
	if err := ui.ipcSvc.SendError(text); err != nil {
		return
	}

	// Console
	text = fmt.Sprintf("[%s] [ERROR] %s", ui.formattedTime(), text)
	_ = ui.AppendConsole(ui.ANSIColored(text, 255, 100, 100, true))
}

func (ui *UIServiceImpl) LogInfo(format string, e ...any) {
	text := fmt.Sprintf(format, e...)

	// Console
	text = fmt.Sprintf("[%s] [INFO] %s", ui.formattedTime(), text)
	_ = ui.AppendConsole(ui.ANSIColored(text, 193, 239, 255, true))
}

func (ui *UIServiceImpl) ansiToHTML(ansi string) string {
	html := ansihtml.ConvertToHTML([]byte(ansi))

	var result strings.Builder

	for _, line := range strings.Split(string(html), "\n") {
		if strings.Contains(line, "ALIGN_CENTER") {
			line = strings.ReplaceAll(line, "ALIGN_CENTER", "")
			line = "<div style=\"text-align:center;\">" + line + "</div>"
		}
		result.WriteString(line + "\n")
	}

	return result.String()
}

func (ui *UIServiceImpl) AppendConsole(ansi string) error {
	html := "\n" + ui.ansiToHTML(ansi)

	if err := ui.ipcSvc.AppendConsole(ansi, html); err != nil {
		return err
	}

	return ui.authSvc.AppendConsole(html)
}

func (ui *UIServiceImpl) ClearConsole() error {
	if err := ui.ipcSvc.ClearConsole(); err != nil {
		return err
	}

	return ui.authSvc.ClearConsole()
}

func (ui *UIServiceImpl) SetConsole(ansi string) error {
	html := ui.ansiToHTML(ansi)

	if err := ui.ipcSvc.SetConsole(ansi, html); err != nil {
		return err
	}

	return ui.authSvc.SetConsole(html)
}

func (ui *UIServiceImpl) NewConsoleMessage() *interfaces.ConsoleMessage {
	return interfaces.NewConsoleMessage(ui)
}

func (ui *UIServiceImpl) ShowConsoleMessage(c *interfaces.ConsoleMessage) error {
	r, g, b, _ := utils.Misc.HexToRGB(c.Color)

	timeFormat := ui.ANSIColored(ui.formattedTime(), 193, 239, 255, false)
	message := ui.ANSIColored("[", 238, 238, 238, false) + timeFormat + ui.ANSIColored("] ", 238, 238, 238, false) +
		ui.ANSIColored(c.Title, r, g, b, false) + "\n"

	for _, info := range c.Info {
		if info.Key == "" {
			message += ui.ANSIColored("  "+info.Value, 193, 239, 255, false) + "\n"
			continue
		}

		message += ui.ANSIColored("  "+info.Key+": ", 238, 238, 238, false) + ui.ANSIColored(info.Value, 193, 239, 255, false) + "\n"
	}

	if constants.CONSOLE {
		fmt.Println(message)
	}

	if err := ui.AppendConsole(message); err != nil {
		return err
	}

	if c.WebhookURL != "" {
		emb := ui.embedSvc.NewEmbed().SetTitle(c.Title)

		for _, info := range c.Info {
			if info.Key == "" {
				emb.AddField("Extra Info", info.Value, false)
				continue
			}

			emb.AddField(info.Key, info.Value, false)
		}

		_ = emb.SendWebhook(c.WebhookURL, c.Important)
	}

	return nil
}
