package reminder

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
)

func StartReminder() {
	for {
		cfg, err := config.GetConfig()
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}

		now := time.Now()

		for _, reminder := range cfg.Reminders {
			// Check if reminder has already passed
			if reminder.At.Before(now) {
				logging.Logging.Reminder(reminder)

				err = deleteReminder(&cfg, reminder)
				if err != nil {
					utils.Logging.Error("Failed to delete reminder: ", err)
					continue
				}
			}
		}

		time.Sleep(time.Second * 5)
	}
}

func deleteReminder(cfg *config.Config, reminder config.Reminder) error {
	for i, r := range cfg.Reminders {
		if r == reminder {
			cfg.Reminders = append(cfg.Reminders[:i], cfg.Reminders[i+1:]...)
			break
		}
	}

	return cfg.Save()
}
