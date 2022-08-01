// go:build darwin

package utils

// Should not be used on macOS.

func GetWindowsFolder(folder int) string {
	return ""
}

func WindowsGetDocumentsFolder() string {
	return ""
}
