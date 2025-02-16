package startup

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

// InitializeBasics sets up basic requirements for both desktop and hosting modes
func InitializeBasics() error {
	if err := utils.IO.CreateFolders(); err != nil {
		return fmt.Errorf("create folders: %w", err)
	}

	// Initialize logging
	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")
	if err := utils.IO.Write(logPath, ""); err != nil {
		return fmt.Errorf("clear log file: %w", err)
	}

	return nil
}

// SetupDNSResolver configures the DNS resolver to use 1.1.1.1
func SetupDNSResolver() {
	var dialer net.Dialer
	net.DefaultResolver = &net.Resolver{
		PreferGo: false,
		Dial: func(context context.Context, _, _ string) (net.Conn, error) {
			return dialer.DialContext(context, "udp", "1.1.1.1:53")
		},
	}
}

// CreateDiscordSession creates and configures a Discord session
func CreateDiscordSession(token string, debug bool) (*discordgo.Session, error) {
	session, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	session.State.TrackVoice = false
	session.State.MaxMessageCount = 50
	session.SetIdentify(cfg.GetClient())
	session.ForceSubscribeGuildIDs = cfg.Selfbot.ForceSubscribedServers

	if debug {
		session.LogLevel = discordgo.LogDebug
	} else {
		session.LogLevel = discordgo.LogWarning
	}

	return session, nil
}

// SetupLogging configures the logging system
func SetupLogging(debug bool) {
	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")

	discordgo.Logger = func(msgL, caller int, format string, a ...any) {
		pc, file, line, _ := runtime.Caller(caller + 1)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "unknown event") && !debug {
			return
		}

		msg = fmt.Sprintf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)

		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return
		}
		defer f.Close()

		_, _ = f.WriteString(msg)

		if msgL == discordgo.LogError {
			logging.UI.Error("DGC error: " + msg)
		}
	}
}
