package assets

import (
	"embed"
	"io"
	"io/fs"
	"os"
)

//go:embed files/*
var files embed.FS

const (
	AssetDefaultConfig = "files/defaultConfig.json"
	AssetLogo          = "files/Logo.png"
)

func Open(name string) (fs.File, error) {
	return files.Open(name)
}

func ReadFile(name string) ([]byte, error) {
	return files.ReadFile(name)
}

func CopyFile(name, dst string) error {
	src, err := Open(name)
	if err != nil {
		return err
	}
	defer src.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, src)
	return err
}
