//go:build desktop

package desktop

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"syscall"

	_websocket "github.com/gorilla/websocket"

	"github.com/LightningDev1/LB-Selfbot/selfbot/assets"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/extensions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/selfbot/rpc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/setup"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/LB-Selfbot/tui"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

func checkForUpdates() {
	resp, err := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{
			"branch":  "selfbot",
			"version": constants.VERSION[1:],
		},
	)
	if err != nil {
		logging.UI.Error("Error checking for updates: " + err.Error())
		time.Sleep(3 * time.Second)
		return
	}

	isLatest := resp.Get("is_latest").MustBool()
	latestVersion := resp.GetPath("latest", "version").MustString()

	if isLatest {
		return
	}

	utils.Logging.Info("LightningBot version", latestVersion, "is now available!")
	utils.Logging.Info("Opening download page...")

	downloadUrl := "https://api.lightning-bot.com/static/executables/"

	if runtime.GOOS == "darwin" {
		downloadUrl += "macos/selfbot"
	} else if runtime.GOOS == "linux" {
		downloadUrl += "linux/" + runtime.GOARCH + "/selfbot"
	} else if runtime.GOOS == "windows" {
		downloadUrl += "windows/" + runtime.GOARCH + "/selfbot.exe"
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
		if !constants.CONSOLE {
			logging.UI.Error("Your LightningBot credentials are invalid! Please re-enter them in the Config tab.")
			break
		}

		utils.Logging.Error("Your LightningBot credentials are invalid!")

		username, password := setup.GetLogin()

		cfg, err := config.GetConfig()
		if err != nil {
			logging.UI.Error("Error loading config: " + err.Error())
			os.Exit(0)
		}

		cfg.Account.LightningBot.Username = username
		cfg.Account.LightningBot.Password = password

		err = cfg.Save()
		if err != nil {
			logging.UI.Error("Error saving config: " + err.Error())
			os.Exit(1)
		}

		utils.System.Restart()
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

func Main(discordgoDebug bool) {
	if constants.DEBUG {
		utils.Logging.Warn("Debug Enabled!")
	}

	session, _ := discordgo.New("")
	session.State.TrackVoice = false
	session.ErrorHandler = errorhandler.ErrorHandler

	initSelfbot(session)

	cfg, err := config.GetConfig()
	if err != nil {
		logging.UI.Error("Error loading config: " + err.Error())
		os.Exit(0)
	}

	if !constants.CONSOLE {
		ipc.IPC.Init()
	}

	username := cfg.Account.LightningBot.Username
	password := cfg.Account.LightningBot.Password
	token := cfg.Account.Discord.Token

	// Check for env variables
	if envUsername, ok := os.LookupEnv("LB_USERNAME"); ok {
		username = envUsername
	}
	if envPassword, ok := os.LookupEnv("LB_PASSWORD"); ok {
		password = envPassword
	}
	if envToken, ok := os.LookupEnv("LB_TOKEN"); ok {
		token = envToken
	}

	websocket.Websocket.LogError = onWebsocketError
	err = websocket.Websocket.Init(
		username,
		password,
	)
	if err != nil {
		if closeError, ok := err.(*_websocket.CloseError); ok {
			onWebsocketError(err, closeError.Code)
		}

		logging.UI.Error("Error connecting to websocket: " + err.Error())
		os.Exit(1)
	}

	session.Token = token
	session.State.MaxMessageCount = 200
	session.SetIdentify(cfg.GetClient())

	session.LogLevel = discordgo.LogWarning

	if discordgoDebug {
		session.LogLevel = discordgo.LogDebug
	}

	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")

	// Clear log file
	_ = utils.IO.Write(logPath, "")

	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		pc, file, line, _ := runtime.Caller(caller + 1)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "unknown event") && !discordgoDebug {
			// Ignore unknown event errors in release mode
			return
		}

		msg = fmt.Sprintf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)

		// Write log to file
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return
		}

		_, _ = f.WriteString(msg)

		f.Close()

		if msgL == discordgo.LogError {
			utils.Logging.Error("DGC error:", msg)
		}
	}

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
			if ctx.Event.Author == nil {
				return false
			}

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

			if utils.SliceContains(ctx.Command.Flags, commands.CmdFlagDangerous) && cfg.Sharing.BlockDangerousCommands {
				logging.Logging.BlockedDangerousCommand(ctx.Event.Author, ctx.Command.Name)
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
		Commands: []*dgc.Command{},
		Middlewares: []dgc.Middleware{
			func(following dgc.ExecutionHandler) dgc.ExecutionHandler {
				return func(ctx *dgc.Ctx) {
					arguments, success := utils.Discord.GetArguments(ctx)
					if !success {
						return
					}

					cfg, err := config.GetConfig()
					if err != nil {
						logging.UI.Error("Failed to get config")
						return
					}

					if cfg.Selfbot.DeleteCommands {
						go func() {
							_ = ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
						}()
					}

					utils.Arguments = arguments

					if ctx.Command.Name != "repeat" {
						eventreactions.UsedCommands.Add(ctx)
					}

					data.UsedCommands++

					logging.Logging.CommandUsed(ctx.Command.Name)

					websocket.Websocket.CheckNonce()

					following(ctx)
				}
			},
		},
	})

	eventreactions.RegisterHandlers(session, router)
	commands.Commands.RegisterCommands(session, router, false)
	extensions.RegisterExtensions(session, router)
	commands.Commands.HelpCommand.RegisterHelpCommand(router, session)

	router.Initialize(session)

	rpc.StartRPC(session)

	if !constants.CONSOLE {
		ipc.IPC.WaitForConnection()
	}
	websocket.Websocket.WaitUntilReady()

	websocket.Websocket.CheckNonce()

	err = errorhandler.Init()
	if err != nil {
		logging.UI.Error("Error initializing error handler: " + err.Error())
	}

	err = session.Open()
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failed") {
			if !constants.CONSOLE {
				logging.UI.Error("Your Discord token is invalid! Please enter a valid token in the Config tab.")
				os.Exit(0)
			}

			if _, ok := os.LookupEnv("LB_TOKEN"); ok {
				logging.UI.Error("Your Discord token is invalid! Please enter a valid token in the LB_TOKEN environment variable.")
				os.Exit(0)
			}

			token := tui.GetToken("Your Discord token is invalid!")

			cfg, err := config.GetConfig()
			if err != nil {
				logging.UI.Error("Error loading config: " + err.Error())
				os.Exit(0)
			}

			cfg.Account.Discord.Token = token

			err = cfg.Save()
			if err != nil {
				logging.UI.Error("Error saving config: " + err.Error())
				os.Exit(1)
			}

			utils.System.Restart()
		}

		logging.UI.Error("An unknown error has occured while connecting to Discord: " + err.Error())
	}

	multitokensniper.Init(session, cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
