//go:build desktop

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/debug"
	"strings"

	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"

	"github.com/LightningDev1/LB-Selfbot/cmd/desktop"
	"github.com/LightningDev1/LB-Selfbot/internal/browser"
	"github.com/LightningDev1/LB-Selfbot/savedtokens"
	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
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

		commands.Commands.RegisterCommands(session, router, true)

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

		commands.Commands.RegisterCommands(session, router, true)

		data := map[string]any{}

		for _, category := range router.Categories {
			data[category.Name] = category.Commands
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)

		os.Exit(0)
	}

	if strings.Contains(args, "-openfolder") {
		_ = utils.System.OpenURL(utils.IO.GetMainFolder())
		os.Exit(0)
	}

	if strings.Contains(args, "-weblogin") {
		token, err := browser.RunWebLogin(context.Background())
		if err != nil {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": err.Error(),
			})

			os.Exit(1)
		}

		s, err := discordgo.New("")
		if err != nil {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": err.Error(),
			})

			os.Exit(1)
		}

		resp, err := api.GetTokenData(s, token)
		if err != nil {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": err.Error(),
			})

			os.Exit(1)
		}

		tokenData, err := simplejson.NewJson(resp)
		if err != nil {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
				"error": err.Error(),
			})

			os.Exit(1)
		}

		user := discordgo.User{
			ID:       tokenData.Get("id").MustString(),
			Avatar:   tokenData.Get("avatar").MustString(),
			Username: tokenData.Get("username").MustString(),
			Token:    token,
		}

		_ = json.NewEncoder(os.Stdout).Encode(map[string]string{
			"token":    user.Token,
			"username": user.Username,
			"id":       user.ID,
			"avatar":   user.AvatarURL(""),
		})

		os.Exit(0)
	}

	if strings.Contains(args, "-savedtokens") {
		tokens := savedtokens.GetSavedTokens()

		_ = json.NewEncoder(os.Stdout).Encode(tokens)

		os.Exit(0)
	}

	if strings.Contains(args, "-debug") {
		discordgoDebug = true
	}

	if strings.Contains(args, "-pprof") {
		go func() {
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}
}

func init() {
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(context, "udp", "1.1.1.1:53")
		},
	}
}

func main() {
	initArgs()

	debug.SetMemoryLimit(300 * 1 << 20)

	desktop.Main(discordgoDebug)
}
