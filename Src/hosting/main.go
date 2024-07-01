//go:build hosting

package hosting

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	_websocket "github.com/gorilla/websocket"

	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/extensions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/selfbot/rpc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

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

func Main() {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	token := os.Getenv("TOKEN")
	mainFolder := os.Getenv("MAIN_FOLDER")
	debug := os.Getenv("DEBUG") == "true"

	constants.CONSOLE = true
	constants.HOSTING = true

	// Setup instance
	session, _ := discordgo.New(token)
	session.State.TrackVoice = false
	session.ErrorHandler = errorhandler.ErrorHandler

	utils.IO.MainFolderOverride = mainFolder

	err := utils.IO.CreateFolders()
	if err != nil {
		logging.UI.Error("Error creating folders: " + err.Error())
		os.Exit(1)
	}

	if utils.IO.IsFirstLaunch() {
		err = config.GetDefaultConfig().Save()
		if err != nil {
			logging.UI.Error("Error creating config: " + err.Error())
			os.Exit(1)
		}
	}

	cfg, err := config.GetConfig()
	if err != nil {
		logging.UI.Error("Error loading config: " + err.Error())
		os.Exit(1)
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
	}

	session.SetIdentify(cfg.GetClient())
	session.State.MaxMessageCount = 200
	session.LogLevel = discordgo.LogWarning

	if debug {
		session.LogLevel = discordgo.LogDebug
	}

	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")

	// Remove log file if it exists
	_ = utils.IO.Write(logPath, "")

	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		pc, file, line, _ := runtime.Caller(caller + 1)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "unknown event") && !debug {
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
			logging.UI.Error("DGC error: " + msg)
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
			isAllowed := func() bool {
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
			}()

			if !isAllowed {
				return false
			}

			if utils.SliceContains(ctx.Command.Flags, commands.CmdFlagHostingDangerous) {
				logging.Logging.BlockedHostingDangerousCommand(ctx.Command.Name)
				return false
			}

			return true
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
	commands.Commands.RegisterCommands(session, router, true)
	extensions.RegisterExtensions(session, router)

	router.Initialize(session)

	rpc.StartRPC(session)

	err = errorhandler.Init()
	if err != nil {
		logging.UI.Error("Error initializing error handler: " + err.Error())
	}

	err = session.Open()
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failed") {
			logging.UI.Error("Invalid token provided")
		} else {
			logging.UI.Error("Error opening connection: " + err.Error())
		}

		os.Exit(1)
	}

	multitokensniper.Init(session, cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
