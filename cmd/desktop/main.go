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

	"github.com/lb-selfbot/selfbot/internal/assets"
	"github.com/lb-selfbot/selfbot/internal/browser"
	"github.com/lb-selfbot/selfbot/internal/commands"

	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/savedtokens"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/setup"
	"github.com/lb-selfbot/selfbot/internal/system"
	"github.com/lb-selfbot/selfbot/internal/utils"
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
		services.UI.LogError("Error checking for updates: " + err.Error())
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

	_ = system.OpenURL(downloadUrl)

	// Wait for user to press enter
	_, _ = fmt.Scanln()

	os.Exit(0)
}

func checkFirstLaunch(services *services.Services) {
	if !utils.IO.IsFirstLaunch() {
		return
	}

	logging.Info("Welcome to LightningBot! Starting setup...")

	username, password := setup.GetLogin()

	services.UI.ClearConsole()
	services.UI.ShowBanner()

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
	// LightningBot uses ANSI escape sequences which cmd.exe won't parse without enabling virtual terminal processing.
	enableVirtualTerminalProcessing()

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

	if err := services.O11y.Init(); err != nil {
		logging.Error("O11y init error:", err)
		os.Exit(1)
	}

	cmdRegistry := commands.NewCmdRegistry()

	// Handle CLI flags
	constants.CONSOLE = !*isLaunched
	constants.PORTABLE = *isPortable

	// Handle CLI commands
	switch {
	case *showVersion:
		fmt.Println(constants.VERSION[1:])
		return

	case *showFeatureList:
		err := services.Discord.Init("", false)
		if err != nil {
			return
		}

		cmdRegistry.Register(services)

		router := services.Discord.GetRouter()

		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"commands":   router.GetCmds(),
			"categories": router.GetCategories(),
		})
		return

	case *generateDocs:
		err := services.Discord.Init("", false)
		if err != nil {
			return
		}

		cmdRegistry.Register(services)

		router := services.Discord.GetRouter()
		categories := router.GetCategories()

		data := make(map[string]any, len(categories))
		for _, category := range categories {
			data[category.Name] = category.Description
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)
		return

	case *openFolder:
		_ = system.OpenURL(utils.IO.GetMainFolder())
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
		services.UI.LogError("Error creating folders: " + err.Error())
		return
	}

	// Initialize icon
	err := assets.CopyFile(assets.AssetLogo, filepath.Join(utils.IO.GetMainFolder(), "Logo.png"))
	if err != nil {
		logging.Error("Error copying icon file:", err)
		return
	}

	if constants.CONSOLE {
		_ = services.UI.ClearConsole()
		services.UI.ShowBanner()

		checkForUpdates(services)
		checkFirstLaunch(services)
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

	cmdRegistry.Register(services)

	if !constants.CONSOLE {
		<-services.IPC.OnConnect()
	}

	// Wait for auth to be ready
	<-services.Auth.OnReady()
	services.Auth.CheckNonce()

	if err := services.Discord.Open(); err != nil {
		logging.Error("Discord open error:", err)
		os.Exit(1)
	}

	services.Reminder.Start()
	services.Variables.Init()

	go func() {
		err := services.MultiTokenSniper.Init(context.Background())
		if err != nil {
			services.UI.LogError("Error initializing multi-token sniper: " + err.Error())
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
