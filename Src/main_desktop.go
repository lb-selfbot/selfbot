//go:build desktop

package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/desktop"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/qrlogin"
)

var discordgoDebug = false

func initArgs() {
	args := strings.Join(os.Args, " ")

	// constants.Console must be true if "-l" flag is not set
	constants.CONSOLE = !strings.Contains(args, "-l")
	constants.PORTABLE = strings.Contains(args, "-portable")

	if strings.Contains(args, "-v") {
		fmt.Println(constants.VERSION[1:])
		os.Exit(0)
	}

	if strings.Contains(args, "-featurelist") {
		session, _ := discordgo.New("")
		router := dgc.Create(&dgc.Router{})

		commands.Commands.RegisterCommands(session, router)

		data := map[string]any{
			"categories": router.Categories,
			"commands":   router.Commands,
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)

		os.Exit(0)
	}

	if strings.Contains(args, "-docgen") {
		session, _ := discordgo.New("")
		router := dgc.Create(&dgc.Router{})

		commands.Commands.RegisterCommands(session, router)

		data := map[string]any{}

		for _, category := range router.Categories {
			data[category.Name] = category.Commands
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)

		os.Exit(0)
	}

	if strings.Contains(args, "-qrlogin") {
		login := qrlogin.Login{
			ImageReadyCallback: func(_ image.Image, raw []byte) error {
				_, err := os.Stdout.Write(raw)
				return err
			},
			TokenCallback: func(token string) error {
				_, err := os.Stdout.WriteString(token)
				return err
			},
		}

		if err := login.Start(); err != nil {
			fmt.Println("ERR:", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	if strings.Contains(args, "-debug") {
		discordgoDebug = true
	}
}

func main() {
	initArgs()

	desktop.Main(discordgoDebug)
}
