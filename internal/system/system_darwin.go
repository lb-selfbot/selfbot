package system

import (
	"os"
	"os/exec"
	"syscall"
)

func openURL(url string) error {
	return exec.Command("open", url).Start()
}

func restartProgram() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	return syscall.Exec(self, os.Args, os.Environ())
}

func shutdown() error {
	return exec.Command("shutdown", "-h", "now").Start()
}

func restart() error {
	return exec.Command("shutdown", "-r", "now").Start()
}

func sleep() error {
	return exec.Command("pmset", "displaysleepnow").Start()
}

func hibernate() error {
	return exec.Command("pmset", "displaysleepnow").Start()
}

func lock() error {
	return exec.Command("pmset", "displaysleepnow").Start()
}
