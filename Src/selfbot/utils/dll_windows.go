package utils

import "syscall"

func (SystemUtils) RunDllFunc(dllName string, funcName string, args ...uintptr) (uintptr, uintptr, error) {
	user32, err := syscall.LoadLibrary(dllName)
	if err != nil {
		return 0, 0, err
	}

	defer syscall.FreeLibrary(user32)

	procAddress, err := syscall.GetProcAddress(user32, funcName)
	if err != nil {
		return 0, 0, err
	}

	r1, r2, winError := syscall.SyscallN(uintptr(procAddress), args...)

	if winError != 0 {
		return r1, r2, err
	}

	return r1, r2, nil
}
