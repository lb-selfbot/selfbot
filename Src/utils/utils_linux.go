// go:build linux

package utils

// Should not be used on Linux.

func GetWindowsFolder(folder int) string {
	return ""
}

func WindowsGetDocumentsFolder() string {
	return ""
}
