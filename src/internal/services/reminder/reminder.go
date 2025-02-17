package reminder

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
)

type ReminderServiceImpl struct {
	uiSvc interfaces.UIService
}

func (r *ReminderServiceImpl) Inject(uiSvc interfaces.UIService) {
	r.uiSvc = uiSvc
}

func (r *ReminderServiceImpl) Start() {
	go r.start()
}

func (r *ReminderServiceImpl) start() {
	for {
		cfg, err := config.GetConfig()
		if err != nil {
			_ = r.uiSvc.LogError("Failed to get config: ", err)
			time.Sleep(time.Second * 5)
			continue
		}

		now := time.Now()

		for _, reminder := range cfg.Reminders {
			if reminder.At.After(now) {
				continue
			}

			// r.uiSvc.LogReminder(reminder)

			cfg.Reminders = utils.RemoveSliceItem(cfg.Reminders, reminder)

			err = cfg.Save()
			if err != nil {
				_ = r.uiSvc.LogError("Failed to save config: ", err)
				continue
			}
		}

		time.Sleep(time.Second * 5)

	}
}
