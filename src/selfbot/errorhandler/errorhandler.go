package errorhandler

import (
	"fmt"
	"runtime/debug"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"

	"github.com/LightningDev1/discordgo"
	"github.com/getsentry/sentry-go"
)

func Init() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "https://d8a09176bdb9454391de5508af97cccd@o4504927797379072.ingest.sentry.io/4504985758662656",
		AttachStacktrace: true,
		SendDefaultPII:   true,
		TracesSampleRate: 1.0,
		Release:          constants.VERSION,
	})
	if err != nil {
		return err
	}

	var environment string
	if constants.HOSTING {
		environment = "Hosting"
	} else if constants.DEBUG {
		environment = "Debug"
	} else {
		environment = "Production"
	}

	if constants.CONSOLE {
		environment += " (Console)"
	} else {
		environment += " (UI)"
	}

	sentry.CurrentHub().ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			Username: cfg.Account.LightningBot.Username,
		})
		scope.SetExtras(map[string]any{
			"version":     constants.VERSION,
			"environment": environment,
			"username":    cfg.Account.LightningBot.Username,
			"start_time":  constants.START_TIME.Format("2006-01-02 15:04:05"),
		})
	})

	return nil
}

func ErrorHandler(session *discordgo.Session, err any) {
	// Sentry error reporting
	switch err := err.(type) {
	case error:
		sentry.CaptureException(err)
	default:
		sentry.CaptureEvent(&sentry.Event{
			Level:   sentry.LevelError,
			Message: fmt.Sprintf("%v", err),
		})
	}

	logging.UI.Error(fmt.Sprintf("A panic has occurred: %s\n\nStacktrace: %s", err, debug.Stack()))
}
