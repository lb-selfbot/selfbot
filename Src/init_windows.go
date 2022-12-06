//go:build windows
// +build windows

package main

import (
	"os"
	"runtime"
	"strings"
	"syscall"
)

func init() {
	// Fix ANSI on Windows cmd
	if strings.Contains(runtime.GOOS, "windows") {
		stdout := syscall.Handle(os.Stdout.Fd())

		var originalMode uint32
		syscall.GetConsoleMode(stdout, &originalMode)
		originalMode |= 0x0004

		syscall.MustLoadDLL("kernel32").MustFindProc("SetConsoleMode").Call(uintptr(stdout), uintptr(originalMode))
	}
}
