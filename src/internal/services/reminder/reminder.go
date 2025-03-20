package reminder

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services/eventbroker"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui/featurelogs"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

// ReminderService handles reminders for the user.
type ReminderService interface {
	// Start initializes the reminder service.
	Start()
}

type reminderService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service
}

// NewReminderService creates a new instance of ReminderService.
func NewReminderService(eventBroker eventbroker.Service, configSvc config.Service) ReminderService {
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

			cfg.Reminders = utils.RemoveSliceItem(cfg.Reminders, reminder)

			err := r.configSvc.SaveConfig(cfg)
			if err != nil {
				r.eventBroker.Publish(eventbroker.EventLogError, "Failed to save config: "+err.Error())
				continue
			}
		}

		time.Sleep(time.Second * 5)
	}
}
