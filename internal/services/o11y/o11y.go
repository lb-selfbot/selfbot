package o11y

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
)

// O11yService handles observability and error tracking.
type O11yService interface {
	// Init initializes the observability service.
	Init() error

	// CaptureError reports an error to the observability service.
	CaptureError(err any)
}

type o11yService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
}

// NewO11yService creates a new instance of O11yService.
func NewO11yService(eventBroker eventbroker.Service, configSvc config.Service) O11yService {
	svc := &o11yService{eventBroker: eventBroker, configSvc: configSvc}
	svc.subscribe()
	return svc
}

func (s *o11yService) subscribe() {
	s.eventBroker.Subscribe(eventbroker.EventCaptureError, s.CaptureError)
}

func (o *o11yService) Init() error {
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

	if err := sentry.Init(clientOptions); err != nil {
		return err
	}

	sentry.CurrentHub().ConfigureScope(func(scope *sentry.Scope) {
		cfg := o.configSvc.GetConfig()

		scope.SetUser(sentry.User{
			Username: cfg.Account.LightningBot.Username,
		})
	})

	return nil
}

func (o *o11yService) CaptureError(err any) {
	switch err := err.(type) {
	case error:
		sentry.CaptureException(err)
	default:
		sentry.CaptureEvent(&sentry.Event{
			Level:   sentry.LevelError,
			Message: fmt.Sprintf("%v", err),
		})
	}

	o.eventBroker.Publish(eventbroker.EventLogError, fmt.Sprintf("An error has occurred: %v", err))
}
