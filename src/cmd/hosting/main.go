//go:build hosting

package hosting

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"

	"github.com/LightningDev1/LB-Selfbot/internal/startup"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/extensions"

	"github.com/LightningDev1/LB-Selfbot/selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var session *discordgo.Session

func Main(ctx context.Context) error {
	constants.CONSOLE = true
	constants.HOSTING = true

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	token := os.Getenv("TOKEN")
	debug := os.Getenv("DEBUG") == "true"

	if username == "" || password == "" || token == "" {
		return fmt.Errorf("USERNAME, PASSWORD, and TOKEN environment variables are required")
	}

	if err := initialize(debug); err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	if err := setupConnections(username, password, token, debug); err != nil {
		return fmt.Errorf("setup connections: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	multitokensniper.Init(session, cfg)

	<-ctx.Done()
	return nil
}

func initialize(debug bool) error {
	startup.SetupLogging(debug)
	return nil
}

func setupConnections(username, password, token string, debug bool) error {
	if err := initWebsocket(username, password); err != nil {
		return fmt.Errorf("init websocket: %w", err)
	}

	var err error
	session, err = startup.CreateDiscordSession(token, debug)
	if err != nil {
		return fmt.Errorf("create discord session: %w", err)
	}

	if err := errorhandler.Init(); err != nil {
		logging.UI.Error("Error initializing error handler: " + err.Error())
	}

	_ = setupRouter(session)

	if err := session.Open(); err != nil {
		return fmt.Errorf("open discord connection: %w", err)
	}
	defer session.Close()

	return nil
}

func initWebsocket(username, password string) error {
	websocket.Websocket.LogError = func(err error, closeCode int) {
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

	return websocket.Websocket.Init(username, password)
}

func setupRouter(session *discordgo.Session) *dgc.Router {
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

				if !slices.Contains(cfg.Sharing.SharedUsers, ctx.Event.Author.ID) &&
					!slices.Contains(cfg.Sharing.SharedUsers, ctx.Event.Author.String()) {
					return false
				}

				if slices.Contains(ctx.Command.Flags, commands.CmdFlagDangerous) && cfg.Sharing.BlockDangerousCommands {
					logging.Logging.BlockedDangerousCommand(ctx.Event.Author, ctx.Command.Name)
					return false
				}

				if cfg.Sharing.ShareAll || slices.Contains([]string{"help", "search", "repeat"}, ctx.Command.Name) {
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
				return slices.Contains(cfg.Sharing.SharedCommands, ctx.Command.Name)
			}()

			if !isAllowed {
				return false
			}

			if slices.Contains(ctx.Command.Flags, commands.CmdFlagHostingDangerous) {
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

					if cfg.Selfbot.DeleteCommands && ctx.Event.ID != "0" {
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

	return router
}

func main() {
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create shutdown channel
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		cancel()
	}()

	// Start hosting mode
	if err := Main(ctx); err != nil {
		logging.UI.Error("Error in hosting mode: " + err.Error())
		os.Exit(1)
	}
}
