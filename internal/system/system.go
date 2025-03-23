package system

import "github.com/LightningDev1/LB-Selfbot/internal/constants"

func OpenURL(url string) error {
	if constants.HOSTING {
		return nil
	}

	return openURL(url)
}

func RestartExecution() error {
	return restartProgram()
}

func ShutdownSystem() error {
	return shutdownSystem()
}

func RestartSystem() error {
	return restartSystem()
}

func SleepSystem() error {
	return sleepSystem()
}

func HibernateSystem() error {
	return hibernateSystem()
}

func LockSystem() error {
	return lockSystem()
}
