package utils

import "syscall"

func (ProcessUtils) RunDllFunc(dllName string, funcName string, args ...uintptr) error {
	user32, err := syscall.LoadLibrary(dllName)
	if err != nil {
		return err
	}

	defer syscall.FreeLibrary(user32)

	lockWorkStation, err := syscall.GetProcAddress(user32, funcName)
	if err != nil {
		return err
	}

	_, _, err = syscall.SyscallN(uintptr(lockWorkStation), args...)

	return err
}
