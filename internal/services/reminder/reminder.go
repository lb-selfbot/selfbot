package reminder

import (
	"time"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

// Service handles reminders for the user.
type Service interface {
	// Start initializes the reminder service.
	Start()
}

type reminderService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
}

// NewService creates a new instance of ReminderService.
func NewService(eventBroker eventbroker.Service, configSvc config.Service) Service {
	return &reminderService{eventBroker: eventBroker, configSvc: configSvc}
}

func (r *reminderService) Start() {
	go r.start()
}

func (r *reminderService) start() {
	for {
		cfg := r.configSvc.GetConfig()
		now := time.Now()

		for _, reminder := range cfg.Reminders {
			if reminder.At.After(now) {
				continue
			}

			r.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.Reminder{
				Reminder: reminder,
			})

			cfg.Reminders = fn.Delete(cfg.Reminders, reminder)

			err := r.configSvc.SaveConfig(cfg)
			if err != nil {
				r.eventBroker.Publish(eventbroker.EventLogError, "Failed to save config: "+err.Error())
				continue
			}
		}

		time.Sleep(time.Second * 5)
	}
}
