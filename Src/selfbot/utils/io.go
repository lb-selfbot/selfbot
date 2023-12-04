package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"golang.org/x/text/unicode/norm"
)

type IOUtils struct {
	MainFolderOverride string
}

func (*IOUtils) SanitizeFilename(filename string) string {
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

func (*IOUtils) Read(file string) (string, bool) {
	fileContents, err := os.ReadFile(filepath.FromSlash(file))
	if err != nil {
		return "", false
	}
	return string(fileContents), true
}

func (*IOUtils) WriteBytes(filePath string, data []byte) error {
	file, err := os.Create(filepath.FromSlash(filePath))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (i *IOUtils) Write(filePath, content string) error {
	return i.WriteBytes(filePath, []byte(content))
}

func (*IOUtils) GetMainFolder() string {
	if IO.MainFolderOverride != "" {
		return IO.MainFolderOverride
	}

	var mainPath string
	var err error

	if constants.PORTABLE {
		mainPath, err = os.Executable()
		if err != nil {
			return "./LightningBot"
		}

		mainPath = filepath.Dir(mainPath)
	} else {
		var err error
		mainPath, err = os.UserConfigDir()
		if err != nil {
			return "./LightningBot"
		}
	}

	return filepath.Join(mainPath, "LightningBot")
}

func (*IOUtils) CreateFolders() error {
	mainFolder := IO.GetMainFolder()

	folders := []string{
		filepath.Join(mainFolder, "Config"),
		filepath.Join(mainFolder, "Animations"),
		filepath.Join(mainFolder, "Dumping"),
		filepath.Join(mainFolder, "Backups"),
		filepath.Join(mainFolder, "Downloads"),
		filepath.Join(mainFolder, "Extensions"),
		filepath.Join(mainFolder, "Executables"),
		filepath.Join(mainFolder, "Executables", "windows"),
		filepath.Join(mainFolder, "Executables", "windows", "amd64"),
		filepath.Join(mainFolder, "Executables", "windows", "386"),
		filepath.Join(mainFolder, "Executables", "linux"),
		filepath.Join(mainFolder, "Executables", "linux", "amd64"),
		filepath.Join(mainFolder, "Executables", "linux", "386"),
		filepath.Join(mainFolder, "Executables", "linux", "arm64"),
		filepath.Join(mainFolder, "Executables", "linux", "arm"),
		filepath.Join(mainFolder, "Executables", "macos"),
	}

	for _, folder := range folders {
		err := os.MkdirAll(folder, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*IOUtils) Exists(path string) bool {
	// Check if file exists
	_, err := os.Stat(filepath.FromSlash(path))

	return !os.IsNotExist(err)
}

func (*IOUtils) IsFirstLaunch() bool {
	// Check if config file exists
	return !IO.Exists(filepath.Join(IO.GetMainFolder(), "Config", "config.json"))
}

var IO = &IOUtils{}
