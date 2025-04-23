package main

import (
	"context"
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

	"github.com/goccy/go-json"

	"github.com/lb-selfbot/selfbot/internal/assets"
	"github.com/lb-selfbot/selfbot/internal/browser"
	"github.com/lb-selfbot/selfbot/internal/commands"
	"github.com/lb-selfbot/selfbot/internal/fs"
	"github.com/lb-selfbot/selfbot/internal/services/authsync"

	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/savedtokens"
	"github.com/lb-selfbot/selfbot/internal/services"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/setup"
	"github.com/lb-selfbot/selfbot/internal/system"
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
	resp, err := http.DefaultClient().Post(
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

	data := resp.Json()
	isLatest := data.Get("is_latest").MustBool()
	latestVersion := data.Get("latest", "version").MustString()

	if isLatest {
		return
	}

	logging.Info("LightningBot version", latestVersion, "is now available!")
	logging.Info("Opening download page...")

	downloadUrl := "https://api.lightning-bot.com/static/executables/"

	switch runtime.GOOS {
	case "darwin":
		downloadUrl += "macos/"
	case "linux":
		downloadUrl += "linux/" + runtime.GOARCH + "/"
	case "windows":
		downloadUrl += "windows/" + runtime.GOARCH + "/"
	default:
		logging.Error("Unsupported OS:", runtime.GOOS)
		os.Exit(1)
	}
	downloadUrl += "selfbot"

	time.Sleep(3 * time.Second)

	_ = system.OpenURL(downloadUrl)

	// Wait for user to press enter
	_, _ = fmt.Scanln()

	os.Exit(0)
}

func checkFirstLaunch(services *services.Services) {
	if !fs.IsFirstLaunch() {
		return
	}

	logging.Info("Welcome to LightningBot! Starting setup...")

	username, password := setup.GetLogin()

	services.Console.Clear()

	token := setup.GetToken()

	cfg := services.Config.GetConfig()

	cfg.Account.LightningBot.Username = username
	cfg.Account.LightningBot.Password = password
	cfg.Account.Discord.Token = token

	if err := services.Config.SaveConfig(cfg); err != nil {
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

		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"commands":   services.Router.GetCommands(),
			"categories": services.Router.GetCategories(),
		})
		return

	case *generateDocs:
		err := services.Discord.Init("", false)
		if err != nil {
			return
		}

		cmdRegistry.Register(services)

		categories := services.Router.GetCategories()

		data := make(map[string]any, len(categories))
		for _, category := range categories {
			data[category.Name] = category.Description
		}

		_ = json.NewEncoder(os.Stdout).Encode(data)
		return

	case *openFolder:
		_ = system.OpenURL(fs.GetMainFolder())
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

	if err := fs.CreateFolders(); err != nil {
		services.UI.LogError("Error creating folders: " + err.Error())
		return
	}

	// Initialize icon
	err := assets.CopyFile(assets.AssetLogo, filepath.Join(fs.GetMainFolder(), "Logo.png"))
	if err != nil {
		logging.Error("Error copying icon file:", err)
		return
	}

	if constants.CONSOLE {
		services.Console.Clear()

		checkForUpdates(services)
		checkFirstLaunch(services)
	} else {
		if err := services.IPC.Connect(); err != nil {
			logging.Error("IPC serve error:", err)
			os.Exit(1)
		}
	}

	cfg := services.Config.GetConfig()

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

	creds := authsync.Credentials{
		Username: username,
		Password: password,
	}

	if err = services.AuthSync.Connect(creds); err != nil {
		os.Exit(1)
	}

	if err := services.Discord.Init(token, *isVerbose); err != nil {
		logging.Error("Discord init error:", err)
		os.Exit(1)
	}

	cmdRegistry.Register(services)

	services.EventReactions.Register(services.Discord.GetSession())

	if !constants.CONSOLE {
		<-services.IPC.GetReadyChan()
	}

	// Wait for auth to be ready
	<-services.AuthSync.GetReadyChan()
	services.AuthSync.CheckNonce()

	if err := services.Discord.Open(); err != nil {
		logging.Error("Discord open error:", err)
		os.Exit(1)
	}

	services.Reminder.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
