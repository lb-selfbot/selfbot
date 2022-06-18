// go:build darwin

package utils

// Should not be used on MacOS.

func GetWindowsFolder(folder int) string {
	return ""
}

func WindowsGetDocumentsFolder() string {
	return ""
}
