package extensions

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type ExtensionsServiceImpl struct {
	uiSvc      interfaces.UIService
	discordSvc interfaces.DiscordService

	loadedExtensions int
}

func (e *ExtensionsServiceImpl) Inject(uiSvc interfaces.UIService, discordSvc interfaces.DiscordService) {
	e.discordSvc = discordSvc
}

func (e *ExtensionsServiceImpl) clearConsole() {
	e.uiSvc.ClearConsole()
	e.uiSvc.ShowBanner(ui.BannerUI & ui.BannerConsole)
}

func (e *ExtensionsServiceImpl) notification(message string) {
	_ = e.uiSvc.ShowNotification(message)
}

func (e *ExtensionsServiceImpl) logError(message string) {
	_ = e.uiSvc.LogError(message)
}

func (e *ExtensionsServiceImpl) Register() {
	session, router := e.discordSvc.GetSession(), e.discordSvc.GetRouter()

	exports := interp.Exports{
		"github.com/LightningDev1/LB-Extensions/lightningbot": map[string]reflect.Value{
			// function, constant and variable definitions
			"ClearConsole":      reflect.ValueOf(e.clearConsole),
			"Notification":      reflect.ValueOf(e.notification),
			"LogError":          reflect.ValueOf(e.logError),
			"GetConfig":         reflect.ValueOf(config.GetConfig),
			"NewConsoleMessage": reflect.ValueOf(e.uiSvc.NewConsoleMessage),
			"NewEmbed":          reflect.ValueOf(embed.NewEmbed),
			"Version":           reflect.ValueOf(&constants.VERSION).Elem(),

			// type definitions
			"Config": reflect.ValueOf((*config.Config)(nil)),
		},
	}

	i := interp.New(interp.Options{})

	_ = i.Use(stdlib.Symbols)
	_ = i.Use(discordgo.Symbols)
	_ = i.Use(dgc.Symbols)
	_ = i.Use(exports)

	handleError := func(err error) {
		_ = e.uiSvc.LogError("An error has occurred while loading extensions: " + err.Error())
		time.Sleep(5 * time.Second)
	}

	dir := path.Join(utils.IO.GetMainFolder(), "Extensions")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Recover from panics
		defer func() {
			if r := recover(); r != nil {
				handleError(fmt.Errorf("panic loading extension: %v", r))
			}
		}()

		// Handle Walk errors
		if err != nil {
			handleError(err)
			return nil
		}

		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Evaluate extension
		_, err = i.EvalPath(path)
		if err != nil {
			handleError(fmt.Errorf("error evaluating extension: %w", err))
			return nil
		}

		register, err := i.Eval("extensions.Register")
		if err != nil {
			handleError(fmt.Errorf("extension does not have a register function: %w", err))
			return nil
		}

		registerFunc, ok := register.Interface().(func(*discordgo.Session, *dgc.Router))
		if !ok {
			handleError(fmt.Errorf("Register is not a function, but %T", register.Interface()))
			return nil
		}

		registerFunc(session, router)

		e.loadedExtensions++

		return nil
	})
	if err != nil {
		handleError(fmt.Errorf("error loading extensions: %w", err))
	}
}

func (e *ExtensionsServiceImpl) GetLoadedExtensions() int { return e.loadedExtensions }
