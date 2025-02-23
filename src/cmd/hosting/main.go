package hosting

import (
	"context"
	"net"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/LightningDev1/LB-Selfbot/internal/assets"
	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/logging"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"

	"github.com/LightningDev1/LB-Selfbot/internal/commands"

	"github.com/LightningDev1/LB-Selfbot/internal/utils"
)

func main() {
	constants.CONSOLE = true
	constants.HOSTING = true

	if constants.DEBUG {
		logging.Warn("Debug mode is enabled")
	}

	// Use Cloudflare DNS for all DNS queries
	net.DefaultResolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		return net.Dial("udp", "1.1.1.1:53")
	}

	// Set memory limit to 300MB
	debug.SetMemoryLimit(300 * 1024 * 1024)

	// Create services
	services := services.NewServices()

	if err := services.O11y.Init(); err != nil {
		services.UI.LogError("O11y init error:", err)
		os.Exit(1)
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

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	token := os.Getenv("TOKEN")
	debug := os.Getenv("DEBUG") == "true"

	if username == "" || password == "" || token == "" {
		services.UI.LogError("USERNAME, PASSWORD, and TOKEN environment variables are required")
		os.Exit(1)
	}

	creds := interfaces.AuthCredentials{
		Username: username,
		Password: password,
	}

	if err = services.Auth.Connect(context.Background(), creds); err != nil {
		os.Exit(1)
	}

	if err := services.Discord.Init(token, debug); err != nil {
		logging.Error("Discord init error:", err)
		os.Exit(1)
	}

	cmdRegistry := commands.NewCmdRegistry()
	cmdRegistry.Register(services)

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
