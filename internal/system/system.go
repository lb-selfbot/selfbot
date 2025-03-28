package system

import "github.com/lb-selfbot/selfbot/internal/constants"

func OpenURL(url string) error {
	if constants.HOSTING {
		return nil
	}

	return openURL(url)
}

func RestartExecution() error {
	return restartProgram()
}

func Shutdown() error {
	return shutdown()
}

func Restart() error {
	return restart()
}

func Sleep() error {
	return sleep()
}

func Hibernate() error {
	return hibernate()
}

func Lock() error {
	return lock()
}
