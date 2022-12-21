package utils

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

type ProcessUtils struct{}

func (ProcessUtils) OpenURL(url string) error {
	if runtime.GOOS == "windows" {
		return exec.Command("explorer", url).Start()
	} else if runtime.GOOS == "darwin" {
		return exec.Command("open", url).Start()
	} else {
		return exec.Command("xdg-open", url).Start()
	}
}

func (ProcessUtils) Restart() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	env := os.Environ()
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}
	return syscall.Exec(self, os.Args, env)
}

func (ProcessUtils) ShutdownPC() error {
	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/s", "/t", "0").Start()
	} else {
		return exec.Command("shutdown", "-h", "now").Start()
	}
}

func (ProcessUtils) RestartPC() error {
	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/r", "/t", "0").Start()
	} else {
		return exec.Command("shutdown", "-r", "now").Start()
	}
}

var Process = &ProcessUtils{}
