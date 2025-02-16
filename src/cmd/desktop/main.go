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

	"github.com/LightningDev1/LB-Selfbot/internal/browser"
	"github.com/LightningDev1/LB-Selfbot/internal/startup"
	"github.com/LightningDev1/LB-Selfbot/internal/tui"
	"github.com/LightningDev1/LB-Selfbot/savedtokens"
	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/extensions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/selfbot/setup"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
)

var session *discordgo.Session

// ProcessArgs processes command line arguments and returns false if the program should exit
func ProcessArgs() bool {
	args := strings.Join(os.Args, " ")

	// constants.Console must be true if "-l" flag is not set
	constants.CONSOLE = !strings.Contains(args, "-l")
	constants.PORTABLE = strings.Contains(args, "-portable")

	if strings.Contains(args, "-v") {
		fmt.Println(constants.VERSION[1:])
		return false
	}

	if strings.Contains(args, "-featurelist") {
		session, _ := discordgo.New("")
		router := setupRouter(session)
		commands.Commands.RegisterCommands(session, router, true)

		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"categories": router.Categories,
			"commands":   router.Commands,
		})
		return false
	}

	if strings.Contains(args, "-docgen") {
		session, _ := discordgo.New("")
		router := setupRouter(session)
		commands.Commands.RegisterCommands(session, router, true)

		data := map[string]any{}
		for _, category := range router.Categories {
			data[category.Name] = category.Commands
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)
		return false
	}

	if strings.Contains(args, "-openfolder") {
		_ = utils.System.OpenURL(utils.IO.GetMainFolder())
		return false
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

		return false
	}

	if strings.Contains(args, "-savedtokens") {
		tokens := savedtokens.GetSavedTokens()

		_ = json.NewEncoder(os.Stdout).Encode(tokens)

		return false
	}

	if strings.Contains(args, "-debug") {
		constants.DEBUG = true
	}

	if strings.Contains(args, "-pprof") {
		go func() {
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}

	return true
}

// Entry point for desktop mode
func Main(ctx context.Context) error {
	if !ProcessArgs() {
		return nil
	}

	if err := initialize(); err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	username, password, token, err := getCredentials()
	if err != nil {
		return fmt.Errorf("get credentials: %w", err)
	}

	if err := setupConnections(username, password, token); err != nil {
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

func initialize() error {
	if err := startup.InitializeBasics(); err != nil {
		return fmt.Errorf("initialize basics: %w", err)
	}

	startup.SetupDNSResolver()
	startup.SetupLogging(constants.DEBUG)

	if err := handleFirstLaunch(); err != nil {
		return fmt.Errorf("first launch setup: %w", err)
	}

	return nil
}

func setupConnections(username, password, token string) error {
	if !constants.CONSOLE {
		ipc.IPC.Init()
	}

	if err := initWebsocket(username, password); err != nil {
		return fmt.Errorf("init websocket: %w", err)
	}

	var err error
	session, err = startup.CreateDiscordSession(token, constants.DEBUG)
	if err != nil {
		if isAuthError(err) {
			return handleAuthError(token)
		}
		return fmt.Errorf("create discord session: %w", err)
	}

	if err := errorhandler.Init(); err != nil {
		logging.UI.Error("Error initializing error handler: " + err.Error())
	}

	if !constants.CONSOLE {
		ipc.IPC.WaitForConnection()
	}
	websocket.Websocket.WaitUntilReady()
	websocket.Websocket.CheckNonce()

	_ = setupRouter(session)

	if err := session.Open(); err != nil {
		if isAuthError(err) {
			return handleAuthError(token)
		}
		return fmt.Errorf("open discord connection: %w", err)
	}
	defer session.Close()

	return nil
}

func handleFirstLaunch() error {
	if !utils.IO.IsFirstLaunch() {
		return nil
	}

	utils.Logging.Info("First launch detected, starting setup...")

	username, password := setup.GetLogin()
	token := setup.GetToken(nil)

	cfg := config.GetDefaultConfig()
	cfg.Account.LightningBot.Username = username
	cfg.Account.LightningBot.Password = password
	cfg.Account.Discord.Token = token

	return cfg.Save()
}

func getCredentials() (username, password, token string, err error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", "", "", fmt.Errorf("load config: %w", err)
	}

	username = cfg.Account.LightningBot.Username
	password = cfg.Account.LightningBot.Password
	token = cfg.Account.Discord.Token

	// Check for env variables
	if envUsername := os.Getenv("LB_USERNAME"); envUsername != "" {
		username = envUsername
	}
	if envPassword := os.Getenv("LB_PASSWORD"); envPassword != "" {
		password = envPassword
	}
	if envToken := os.Getenv("LB_TOKEN"); envToken != "" {
		token = envToken
	}

	return username, password, token, nil
}

func handleAuthError(token string) error {
	if !constants.CONSOLE {
		return fmt.Errorf("discord token is invalid")
	}

	if _, exists := os.LookupEnv("LB_TOKEN"); exists {
		return fmt.Errorf("discord token in LB_TOKEN environment variable is invalid")
	}

	newToken := tui.GetToken("Your Discord token is invalid!")

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	cfg.Account.Discord.Token = newToken
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	utils.System.Restart()
	return nil
}

func isAuthError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Authentication failed")
}

func initWebsocket(username, password string) error {
	websocket.Websocket.LogError = func(err error, closeCode int) {
		switch closeCode {
		case 3004, 3010:
			if !constants.CONSOLE {
				logging.UI.Error("Your LightningBot credentials are invalid! Please re-enter them in the Config tab.")
				os.Exit(1)
				return
			}

			utils.Logging.Error("Your LightningBot credentials are invalid!")
			username, password := setup.GetLogin()

			cfg, err := config.GetConfig()
			if err != nil {
				logging.UI.Error("Error loading config: " + err.Error())
				os.Exit(1)
				return
			}

			cfg.Account.LightningBot.Username = username
			cfg.Account.LightningBot.Password = password

			if err := cfg.Save(); err != nil {
				logging.UI.Error("Error saving config: " + err.Error())
				os.Exit(1)
				return
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

	Main(context.Background())
}
