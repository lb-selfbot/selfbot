package system

import (
	"os"
	"os/exec"
)

func openURL(url string) error {
	return exec.Command("explorer", url).Start()
}

func restartProgram() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}

	cmd := exec.Command(self, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	_ = cmd.Run()

	// If the program still exists after it should
	// have restarted, exit with an error.
	os.Exit(1)
	return nil
}

func shutdownSystem() error {
	return exec.Command("shutdown", "/s", "/t", "0").Start()
}

func restartSystem() error {
	return exec.Command("shutdown", "/r", "/t", "0").Start()
}

func sleepSystem() error {
	return exec.Command("rundll32", "powrprof.dll,SetSuspendState", "0", "1", "0").Start()
}

func hibernateSystem() error {
	return exec.Command("rundll32", "powrprof.dll,SetSuspendState", "1", "1", "0").Start()
}

func lockSystem() error {
	return exec.Command("rundll32", "user32.dll,LockWorkStation").Start()
}
