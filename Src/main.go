package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"syscall"

	_websocket "github.com/gorilla/websocket"

	"github.com/LightningDev1/LB-Selfbot/assets"
	"github.com/LightningDev1/LB-Selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/rpc"
	"github.com/LightningDev1/LB-Selfbot/setup"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

func initArgs() {
	args := strings.Join(os.Args, " ")

	// constants.Console must be true if "-l" flag is not set
	constants.CONSOLE = !strings.Contains(args, "-l")

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
}

func checkForUpdates() {
	httpResponse := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{
			"branch":  "selfbot",
			"version": constants.VERSION[1:],
		},
	)
	if httpResponse.Error != nil {
		logging.UI.Error("Error checking for updates: " + httpResponse.Error.Error())
		time.Sleep(3 * time.Second)
		return
	}

	json := httpResponse.Json()

	isLatest := json.Get("is_latest").MustBool()
	latestVersion := json.GetPath("latest", "version").MustString()

	if isLatest {
		return
	}

	utils.Logging.Info("LightningBot version", latestVersion, "is now available!")
	utils.Logging.Info("Opening download page...")

	downloadUrl := "https://api.lightning-bot.com/static/executables/"

	if runtime.GOOS == "darwin" {
		downloadUrl += "macos/selfbot"
	} else if runtime.GOOS == "linux" {
		downloadUrl += "linux/"
		downloadUrl += runtime.GOARCH
		downloadUrl += "/selfbot"
	} else if runtime.GOOS == "windows" {
		downloadUrl += "windows/"
		downloadUrl += runtime.GOARCH
		downloadUrl += "/selfbot.exe"
	}

	time.Sleep(3 * time.Second)

	_ = utils.System.OpenURL(downloadUrl)

	// Wait for user to press enter
	_, _ = fmt.Scanln()

	os.Exit(0)
}

func initSelfbot(session *discordgo.Session) {
	if !constants.CONSOLE {
		return
	}

	logging.Console.Clear()
	logging.Logging.PrintBanner(true, false)

	checkForUpdates()

	err := utils.IO.CreateFolders()
	if err != nil {
		logging.UI.Error("Error creating folders: " + err.Error())
	}

	iconFile := filepath.Join(utils.IO.GetMainFolder(), "Logo.png")
	if !utils.IO.Exists(iconFile) {
		icon, err := assets.Assets.ReadFile("Logo.png")
		if err != nil {
			logging.UI.Error("Error reading icon file: " + err.Error())
		}

		_ = utils.IO.WriteBytes(iconFile, icon)
	}

	isFirstLaunch := utils.IO.IsFirstLaunch()

	if !isFirstLaunch {
		return
	}

	utils.Logging.Info("First launch detected, starting setup...")

	username, password := setup.GetLogin()

	logging.Console.Clear()
	logging.Logging.PrintBanner(true, false)

	token := setup.GetToken(session)

	cfg := config.GetDefaultConfig()

	cfg.Account.LightningBot.Username = username
	cfg.Account.LightningBot.Password = password
	cfg.Account.Discord.Token = token

	err = cfg.Save()
	if err != nil {
		logging.UI.Error("Error saving config: " + err.Error())
		os.Exit(1)
	}
}

func onWebsocketError(err error, closeCode int) {
	switch closeCode {
	case 3004, 3010:
		logging.UI.Error("Your LightningBot credentials are invalid! Please re-enter them in the Config tab.")
	case 3011:
		logging.UI.Error("Your free trial has expired! Please purchase a license key to continue using LightningBot.")
	case 3003:
		logging.UI.Error("WS Error 3003: " + err.Error())
	case 3006:
		logging.UI.Error("Another selfbot has connected using your LightningBot account! This instance will now exit.")
	default:
		logging.UI.Error("An unexpected error has occurred while connecting to the websocket: " + err.Error())
	}

	os.Exit(1)
}

func main() {
	initArgs()

	if constants.DEBUG {
		utils.Logging.Warn("Debug Enabled!")
	}

	session, _ := discordgo.New("")

	initSelfbot(session)

	cfg, err := config.GetConfig()
	if err != nil {
		logging.UI.Error("Error loading config: " + err.Error())
		os.Exit(0)
	}

	if !constants.CONSOLE {
		ipc.IPC.Init()
	}

	websocket.Websocket.LogError = onWebsocketError
	err = websocket.Websocket.Init(
		cfg.Account.LightningBot.Username,
		cfg.Account.LightningBot.Password,
	)
	if err != nil {
		if closeError, ok := err.(*_websocket.CloseError); ok {
			onWebsocketError(err, closeError.Code)
		}

		logging.UI.Error("Error connecting to websocket: " + err.Error())
		os.Exit(1)
	}
	session.Token = cfg.Account.Discord.Token
	session.ShouldSubscribeGuilds = false
	session.State.MaxMessageCount = 200
	session.SetIdentify(cfg.GetClient())

	router := dgc.Create(&dgc.Router{
		PrefixFunc: func() []string {
			cfg, err := config.GetConfig()
			if err != nil {
				return []string{"!"}
			}
			return []string{cfg.Selfbot.CommandPrefix}
		},
		IgnorePrefixCase: true,
		IsUserAllowedFunc: func(ctx *dgc.Ctx) bool {
			if ctx.Event.Author.ID == session.State.User.ID {
				return true
			}
			cfg, err := config.GetConfig()
			if err != nil {
				return false
			}

			if !utils.SliceContains(cfg.Sharing.SharedUsers, ctx.Event.Author.ID) &&
				!utils.SliceContains(cfg.Sharing.SharedUsers, ctx.Event.Author.String()) {
				return false
			}

			if cfg.Sharing.ShareAll || utils.SliceContains([]string{"help", "search", "repeat"}, ctx.Command.Name) {
				return true
			}

			categoryName := ""
			if ctx.Command.Category != nil {
				categoryName = ctx.Command.Category.Name
			}
			for _, category := range cfg.Sharing.SharedCategories {
				if strings.EqualFold(category, categoryName) {
					return true
				}
			}
			return utils.SliceContains(cfg.Sharing.SharedCommands, ctx.Command.Name)
		},
		Commands:    []*dgc.Command{},
		Middlewares: []dgc.Middleware{},
	})

	eventhandlers.EventHandlers.RegisterHandlers(session, router)
	commands.Commands.RegisterCommands(session, router)

	router.Initialize(session)

	rpc.StartRPC(session)

	if !constants.CONSOLE {
		ipc.IPC.WaitForConnection()
	}
	websocket.Websocket.WaitUntilReady()

	err = session.Open()
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failed") {
			logging.UI.Error("Your Discord token is invalid! Please enter a valid token in the Config tab.")
		} else {
			logging.UI.Error("An unknown error has occured while connecting to Discord: " + err.Error())
		}

		os.Exit(0)
	}

	multitokensniper.Init(cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
