package system

import (
	"os"
	"os/exec"
	"syscall"
)

func openURL(url string) error {
	return exec.Command("xdg-open", url).Start()
}

func restartProgram() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	return syscall.Exec(self, os.Args, os.Environ())
}

func shutdownSystem() error {
	return exec.Command("shutdown", "-h", "now").Start()
}

func restartSystem() error {
	return exec.Command("shutdown", "-r", "now").Start()
}

func sleepSystem() error {
	return exec.Command("systemctl", "suspend").Start()
}

func hibernateSystem() error {
	return exec.Command("systemctl", "suspend").Start()
}

func lockSystem() error {
	return exec.Command("systemctl", "suspend").Start()
}
