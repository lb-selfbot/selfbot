package hosting

import (
	"context"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	_ "net/http/pprof"

	"github.com/lb-selfbot/selfbot/internal/assets"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/fs"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/authsync"

	"github.com/lb-selfbot/selfbot/internal/commands"
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

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	token := os.Getenv("TOKEN")
	debug := os.Getenv("DEBUG") == "true"

	if username == "" || password == "" || token == "" {
		services.UI.LogError("USERNAME, PASSWORD, and TOKEN environment variables are required")
		os.Exit(1)
	}

	creds := authsync.Credentials{
		Username: username,
		Password: password,
	}

	if err = services.AuthSync.Connect(creds); err != nil {
		os.Exit(1)
	}

	if err := services.Discord.Init(token, debug); err != nil {
		logging.Error("Discord init error:", err)
		os.Exit(1)
	}

	cmdRegistry := commands.NewCmdRegistry()
	cmdRegistry.Register(services)

	services.EventReactions.Register(services.Discord.GetSession())

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
