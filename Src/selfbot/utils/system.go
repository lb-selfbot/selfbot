package utils

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
)

type SystemUtils struct{}

func (SystemUtils) OpenURL(url string) error {
	if constants.HOSTING {
		return nil
	}

	if runtime.GOOS == "windows" {
		return exec.Command("explorer", url).Start()
	} else if runtime.GOOS == "darwin" {
		return exec.Command("open", url).Start()
	} else {
		return exec.Command("xdg-open", url).Start()
	}
}

func (SystemUtils) Restart() error {
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

func (SystemUtils) ShutdownPC() error {
	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/s", "/t", "0").Start()
	} else {
		return exec.Command("shutdown", "-h", "now").Start()
	}
}

func (SystemUtils) RestartPC() error {
	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/r", "/t", "0").Start()
	} else {
		return exec.Command("shutdown", "-r", "now").Start()
	}
}

func (SystemUtils) SleepPC() error {
	if runtime.GOOS == "windows" {
		_, _, err := System.RunDllFunc("powrprof.dll", "SetSuspendState", 2, 0, 1, 0)
		return err
	} else if runtime.GOOS == "darwin" {
		return exec.Command("pmset", "displaysleepnow").Start()
	} else {
		return exec.Command("systemctl", "suspend").Start()
	}
}

func (SystemUtils) LockPC() error {
	if runtime.GOOS == "windows" {
		_, _, err := System.RunDllFunc("user32.dll", "LockWorkStation")
		return err
	} else if runtime.GOOS == "darwin" {
		return exec.Command("pmset", "displaysleepnow").Start()
	} else {
		return exec.Command("systemctl", "suspend").Start()
	}
}

// func (SystemUtils) GetIdleTime() (time.Duration, error) {
// 	if runtime.GOOS != "windows" {
// 		return 0, nil
// 	}

// 	var lastInputInfo struct {
// 		cbSize uint32
// 		dwTime uint32
// 	}
// 	lastInputInfo.cbSize = uint32(unsafe.Sizeof(lastInputInfo))

// 	_, _, err := System.RunDllFunc("user32.dll", "GetLastInputInfo", uintptr(unsafe.Pointer(&lastInputInfo)))
// 	if err != nil {
// 		return 0, err
// 	}

// 	tickCount, _, err := System.RunDllFunc("kernel32.dll", "GetTickCount")
// 	if err != nil {
// 		return 0, err
// 	}

// 	tickCount32 := uint32(tickCount)

// 	if tickCount32 < lastInputInfo.dwTime {
// 		const kMaxDWORD = uint32(0xFFFFFFFF)
// 		timeBeforeWrap := kMaxDWORD - lastInputInfo.dwTime
// 		timeAfterWrap := tickCount32
// 		return time.Duration(timeBeforeWrap+timeAfterWrap) * time.Millisecond, nil
// 	} else {
// 		return time.Duration(tickCount32-lastInputInfo.dwTime) * time.Millisecond, nil
// 	}
// }

var System = &SystemUtils{}
