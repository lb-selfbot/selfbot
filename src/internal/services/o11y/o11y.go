package o11y

import (
	"fmt"
	"runtime/debug"

	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/discordgo"
	"github.com/getsentry/sentry-go"
)

type O11yServiceImpl struct {
	uiSvc interfaces.UIService
}

func (o *O11yServiceImpl) Inject(uiSvc interfaces.UIService) {
	o.uiSvc = uiSvc
}

func (o *O11yServiceImpl) Init() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	var environment string

	if constants.HOSTING {
		environment = "hosting"
	} else if constants.DEBUG {
		environment = "debug"
	} else {
		environment = "prod"
	}

	if constants.CONSOLE {
		environment += "-console"
	} else {
		environment += "-ui"
	}

	clientOptions := sentry.ClientOptions{
		Dsn:              "https://d8a09176bdb9454391de5508af97cccd@o4504927797379072.ingest.sentry.io/4504985758662656",
		Debug:            constants.DEBUG,
		AttachStacktrace: true,
		SendDefaultPII:   true,
		TracesSampleRate: 1.0,
		Release:          constants.VERSION,
		Environment:      environment,
		Tags: map[string]string{
			"start_time": constants.START_TIME.Format("2006-01-02 15:04:05"),
		},
	}

	if err = sentry.Init(clientOptions); err != nil {
		return err
	}

	sentry.CurrentHub().ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			Username: cfg.Account.LightningBot.Username,
		})
	})

	return nil
}

func (o *O11yServiceImpl) ReportError(err any) {
	switch err := err.(type) {
	case error:
		sentry.CaptureException(err)
	default:
		sentry.CaptureEvent(&sentry.Event{
			Level:   sentry.LevelError,
			Message: fmt.Sprintf("%v", err),
		})
	}

	o.uiSvc.LogError(fmt.Sprintf("An error has occurred: %v\n\nStacktrace: %v", err, debug.Stack()))
}

func (o *O11yServiceImpl) DiscordErrorHandler(_ *discordgo.Session, err any) {
	o.ReportError(err)
}
