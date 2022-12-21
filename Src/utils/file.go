package utils

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type FileUtils struct{}

func (*FileUtils) SanitizeFilename(filename string) string {
	blacklist := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", "\x00"}
	reserved := []string{
		"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5",
		"COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5",
		"LPT6", "LPT7", "LPT8", "LPT9",
	}
	for _, c := range blacklist {
		filename = strings.Replace(filename, c, "", -1)
	}
	filename = norm.NFKD.String(filename)
	filename = strings.TrimSuffix(filename, ". ")
	filename = strings.TrimSpace(filename)
	isDots := true
	for _, c := range filename {
		// Remove all characters in filename below code point 32
		if c < ' ' {
			filename = strings.Replace(filename, string(c), "", -1)
		}
		if isDots && c != '.' {
			isDots = false
		}
	}
	if isDots {
		filename = "__" + filename
	}
	if SliceContains(reserved, filename) {
		filename = "__" + filename
	}
	if len(filename) > 255 {
		filename = filename[:255]
	}
	return filename
}

func (*FileUtils) Read(file string) (string, bool) {
	fileContents, err := os.ReadFile(filepath.FromSlash(file))
	if err != nil {
		return "", false
	}
	return string(fileContents), true
}

func (*FileUtils) WriteBytes(filePath string, data []byte) error {
	file, err := os.Create(filepath.FromSlash(filePath))
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		_ = file.Close()
		return err
	}
	_ = file.Close()
	return nil
}

func (f *FileUtils) Write(filePath, content string) error {
	return f.WriteBytes(filePath, []byte(content))
}

func (*FileUtils) GetMainFolder() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, "LightningBot")
}

var File = &FileUtils{}
