package extensions

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/constants"

	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/utils"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// ExtensionsService handles the loading and interop of extensions.
type ExtensionsService interface {
	// GetNumExtensions returns the number of loaded extensions.
	GetNumExtensions() int
}

type extensionsService struct {
	eventBroker eventbroker.Service

	numLoadedExtensions int
}

// NewExtensionsService creates a new instance of ExtensionsService.
func NewExtensionsService(eventBroker eventbroker.Service) ExtensionsService {
	svc := &extensionsService{eventBroker: eventBroker}
	svc.subscribe()
	return svc
}

func (e *extensionsService) subscribe() {
	e.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session := payload.(*discordgo.Session)
		if session.State.User.ID != session.State.User.ID {
			return
		}

		e.load(session)
	})
}

func (e *extensionsService) load(session *discordgo.Session) {
	exports := interp.Exports{
		"github.com/LightningDev1/LB-Extensions/lightningbot": map[string]reflect.Value{
			// function, constant and variable definitions
			"ClearConsole":      reflect.ValueOf(e.clearConsole),
			"Notification":      reflect.ValueOf(e.notification),
			"LogError":          reflect.ValueOf(e.logError),
			"GetConfig":         reflect.ValueOf(config.GetConfig),
			"NewConsoleMessage": reflect.ValueOf(e.uiSvc.NewConsoleMessage),
			"NewEmbed":          reflect.ValueOf(e.embedSvc.NewEmbed),
			"Version":           reflect.ValueOf(&constants.VERSION).Elem(),

			// type definitions
			"Config": reflect.ValueOf((*config.Config)(nil)),
		},
	}

	i := interp.New(interp.Options{})

	_ = i.Use(stdlib.Symbols)
	_ = i.Use(discordgo.Symbols)
	// TODO: add router symbols
	// _ = i.Use(router.Symbols)
	_ = i.Use(exports)

	handleError := func(err error) {
		e.eventBroker.Publish(eventbroker.EventLogError, "An error has occurred while loading extensions: "+err.Error())
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

		registerFunc, ok := register.Interface().(func(*discordgo.Session, router.Router))
		if !ok {
			handleError(fmt.Errorf("Register is not a function, but %T", register.Interface()))
			return nil
		}

		registerFunc(session, cmdRouter)

		e.numLoadedExtensions++

		return nil
	})
	if err != nil {
		handleError(fmt.Errorf("error loading extensions: %w", err))
	}
}

func (e *extensionsService) GetNumExtensions() int {
	return e.numLoadedExtensions
}
