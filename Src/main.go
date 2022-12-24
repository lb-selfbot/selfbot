package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"syscall"

	"github.com/LightningDev1/LB-Selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/sessions"
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
		utils.Logging.Error("Error checking for updates:", httpResponse.Error)
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

	_ = utils.Process.OpenURL(downloadUrl)

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
		utils.Logging.Error("Error creating folders: ", err)
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
		utils.Logging.Error("Error saving config:", err)
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	initArgs()

	if constants.DEBUG {
		utils.Logging.Warn("Debug Enabled!")
	}

	session, _ := discordgo.New("")

	initSelfbot(session)

	cfg, err := config.GetConfig()
	if err != nil {
		utils.Logging.Error(err)
		if !constants.DEBUG {
			os.Exit(0)
		}
	}

	if !constants.CONSOLE {
		ipc.IPC.Init()
	}

	err = websocket.Init(cfg.Account.LightningBot.Username, cfg.Account.LightningBot.Password)
	if err != nil {
		utils.Logging.Error(err)
		if !constants.DEBUG {
			os.Exit(0)
		}
	}

	session.Token = cfg.Account.Discord.Token
	session.Identify.Token = cfg.Account.Discord.Token

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

	sessions.StartCheckingSessions(session)

	if !constants.DEBUG {
		if !constants.CONSOLE {
			ipc.IPC.WaitForConnection()
		}
		websocket.WaitForConnection()
	}

	err = session.Open()
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failed") {
			utils.Logging.Error("Invalid token!")
		} else {
			utils.Logging.Error(err)
		}

		if !constants.DEBUG {
			os.Exit(0)
		}
	}

	multitokensniper.Init(cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
