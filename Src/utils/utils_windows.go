// go:build windows

package utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetWindowsFolder(folder int) string {
	user32DLL := windows.NewLazyDLL("shell32.dll")
	getFolderPath := user32DLL.NewProc("SHGetFolderPathW")
	buffer := [windows.MAX_PATH]uint16{}
	nilPointer := uintptr(unsafe.Pointer(nil))
	getFolderPath.Call(nilPointer, uintptr(folder), nilPointer, 0, uintptr(unsafe.Pointer(&buffer[0])))
	return windows.UTF16ToString(buffer[:])
}

func WindowsGetDocumentsFolder() string {
	return GetWindowsFolder(5)
}
