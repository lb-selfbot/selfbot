package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/browser"
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/logging"
	"github.com/LightningDev1/LB-Selfbot/internal/savedtokens"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/services/ui"
	"github.com/LightningDev1/LB-Selfbot/selfbot/assets"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/errorhandler"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/setup"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var (
	// CLI flags
	isLaunched = flag.Bool("l", false, "Launched by UI application")
	isPortable = flag.Bool("portable", false, "Portable mode")
	isVerbose  = flag.Bool("verbose", false, "Verbose mode")
	isPprof    = flag.Bool("pprof", false, "Enable pprof")

	// CLI commands
	showVersion     = flag.Bool("version", false, "Show version information")
	showFeatureList = flag.Bool("featurelist", false, "Show feature list")
	generateDocs    = flag.Bool("docgen", false, "Generate documentation")
	openFolder      = flag.Bool("openfolder", false, "Open data folder")
	runWebLogin     = flag.Bool("weblogin", false, "Run web login")
	getSavedTokens  = flag.Bool("savedtokens", false, "Get saved tokens")
)

func checkForUpdates(services *services.Services) {
	resp, err := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/check-version",
		map[string]string{
			"branch":  "selfbot",
			"version": constants.VERSION[1:],
		},
	)
	if err != nil {
		_ = services.UI.LogError("Error checking for updates: " + err.Error())
		time.Sleep(3 * time.Second)
		return
	}

	isLatest := resp.Get("is_latest").MustBool()
	latestVersion := resp.GetPath("latest", "version").MustString()

	if isLatest {
		return
	}

	logging.Info("LightningBot version", latestVersion, "is now available!")
	logging.Info("Opening download page...")

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

func checkFirstLaunch() {
	if !utils.IO.IsFirstLaunch() {
		return
	}

	logging.Info("Welcome to LightningBot! Starting setup...")

	username, password := setup.GetLogin()

	logging.Console.Clear()
	logging.Console.PrintBanner()

	token := setup.GetToken()

	cfg := config.GetDefaultConfig()

	cfg.Account.LightningBot.Username = username
	cfg.Account.LightningBot.Password = password
	cfg.Account.Discord.Token = token

	if err := cfg.Save(); err != nil {
		logging.Error("Save config error:", err)
		os.Exit(1)
	}
}

func main() {
	if constants.DEBUG {
		logging.Warn("Debug mode is enabled")
	}

	flag.Parse()

	// Use Cloudflare DNS for all DNS queries
	net.DefaultResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		return net.Dial("udp", "1.1.1.1:53")
	}

	// Set memory limit to 300MB
	debug.SetMemoryLimit(300 * 1024 * 1024)

	// Create services
	services := services.NewServices()

	// Handle CLI flags
	constants.CONSOLE = !*isLaunched
	constants.PORTABLE = *isPortable

	// Handle CLI commands
	switch {
	case *showVersion:
		fmt.Println(constants.VERSION[1:])
		return

	case *showFeatureList:
		router := dgc.Create(&dgc.Router{})
		commands.Commands.RegisterCommands(nil, router, true)
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"commands":   router.Commands,
			"categories": router.Categories,
		})
		return

	case *generateDocs:
		router := dgc.Create(&dgc.Router{})
		commands.Commands.RegisterCommands(nil, router, true)

		data := make(map[string]any, len(router.Categories))
		for _, category := range router.Categories {
			data[category.Name] = category.Description
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)
		return

	case *openFolder:
		_ = utils.System.OpenURL(utils.IO.GetMainFolder())
		return

	case *runWebLogin:
		handleError := func(err error) {
			_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
				"error": err.Error(),
			})
		}

		token, err := browser.RunWebLogin(context.Background())
		if err != nil {
			handleError(err)
			return
		}

		session, err := discordgo.New(token)
		if err != nil {
			handleError(err)
			return
		}
		session.SetIdentify(discordgo.IdentifyDiscordClient)
		user, err := session.User("@me")
		if err != nil {
			handleError(err)
			return
		}

		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"token":    token,
			"username": user.Username,
			"id":       user.ID,
			"avatar":   user.AvatarURL(""),
		})
		return

	case *getSavedTokens:
		tokens := savedtokens.GetSavedTokens()
		_ = json.NewEncoder(os.Stdout).Encode(tokens)
		return
	}

	if err := utils.IO.CreateFolders(); err != nil {
		_ = services.UI.LogError("Error creating folders: " + err.Error())
		return
	}

	// Initialize icon
	iconFile := filepath.Join(utils.IO.GetMainFolder(), "Logo.png")
	if !utils.IO.Exists(iconFile) && !constants.HOSTING {
		icon, err := assets.Assets.ReadFile("Logo.png")
		if err != nil {
			logging.Error("Error reading icon file:", err)
			return
		}

		_ = utils.IO.WriteBytes(iconFile, icon)
	}

	if constants.CONSOLE {
		_ = services.UI.ClearConsole()
		_ = services.UI.ShowBanner(ui.BannerConsole)

		checkForUpdates(services)
		checkFirstLaunch()
	} else {
		if err := services.IPC.Serve(); err != nil {
			logging.Error("IPC serve error:", err)
			os.Exit(1)
		}
	}

	cfg, err := config.GetConfig()
	if err != nil {
		logging.Error("Load config error:", err)
		os.Exit(1)
	}

	username, ok := os.LookupEnv("LB_USERNAME")
	if !ok {
		username = cfg.Account.LightningBot.Username
	}
	password, ok := os.LookupEnv("LB_PASSWORD")
	if !ok {
		password = cfg.Account.LightningBot.Password
	}
	token, ok := os.LookupEnv("LB_TOKEN")
	if !ok {
		token = cfg.Account.Discord.Token
	}

	creds := interfaces.AuthCredentials{
		Username: username,
		Password: password,
	}

	if err = services.Auth.Connect(context.Background(), creds); err != nil {
		os.Exit(1)
	}

	if err := services.Discord.Init(token, *isVerbose); err != nil {
		logging.Error("Discord init error:", err)
		os.Exit(1)
	}

	if !constants.CONSOLE {
		<-services.IPC.OnConnect()
	}

	// Wait for auth to be ready
	<-services.Auth.OnReady()
	services.Auth.CheckNonce()

	if err := errorhandler.Init(); err != nil {
		logging.Error("Error handler init error:", err)
		os.Exit(1)
	}

	if err := services.Discord.Open(); err != nil {
		logging.Error("Discord open error:", err)
		os.Exit(1)
	}

	services.Reminder.Start()

	go func() {
		err := services.MultiTokenSniper.Init(context.Background())
		if err != nil {
			_ = services.UI.LogError("Error initializing multi-token sniper: " + err.Error())
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
