package discord

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/system"
	"github.com/lb-selfbot/selfbot/internal/tui"
	"github.com/lb-selfbot/selfbot/internal/utils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Service handles Discord-related operations.
type Service interface {
	// Init initializes the Discord session with the provided token.
	Init(token string, verbose bool) error

	// GetSession returns the Discord session.
	GetSession() *discordgo.Session

	// Open opens the Discord session.
	Open() error

	// SetHypesquad sets the Hypesquad house for the session user.
	SetHypesquad(house int) error

	// LeaveHypesquad leaves the Hypesquad house for the session user.
	LeaveHypesquad() error

	// GetTokenInfo retrieves information about the provided token.
	GetTokenInfo(token string) (*jester.Data, error)

	// UpdateSettings updates the settings for the session user.
	UpdateSettings(path string, settings any) error

	// SetAllowedFriendRequests updates the allowed friend requests for the session user.
	SetFriendRequestsAllowed(allowed bool) error
}

type discordService struct {
	eventBroker eventbroker.Service
	configSvc   config.Service

	session *discordgo.Session

	verbose bool
}

// NewService creates a new DiscordService instance.
func NewService(eventBroker eventbroker.Service, configSvc config.Service) Service {
	return &discordService{eventBroker: eventBroker, configSvc: configSvc}
}

func (d *discordService) Init(token string, verbose bool) (err error) {
	d.verbose = verbose
	discordgo.Logger = d.discordgoLogger()

	d.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %w", err)
	}

	cfg := d.configSvc.GetConfig()

	d.session.Token = token
	d.session.State.MaxMessageCount = 50
	d.session.SetIdentify(cfg.GetClient())
	d.session.ErrorHandler = func(session *discordgo.Session, err any) {
		d.eventBroker.Publish(eventbroker.EventCaptureError, err)
	}
	d.session.ForceSubscribeGuildIDs = cfg.Selfbot.ForceSubscribedServers

	if verbose {
		d.session.LogLevel = discordgo.LogDebug
	} else {
		d.session.LogLevel = discordgo.LogWarning
	}

	return nil
}

func (d *discordService) GetSession() *discordgo.Session {
	return d.session
}

func (d *discordService) Open() error {
	err := d.session.Open()
	if err == nil {
		return nil
	}

	// Check if token is invalid
	if strings.Contains(err.Error(), "Authentication failed") {
		if !constants.CONSOLE {
			d.eventBroker.Publish(eventbroker.EventLogError, "Your Discord token is invalid! Please re-enter it in the Config tab.")
			return err
		}

		if _, ok := os.LookupEnv("LB_TOKEN"); ok {
			d.eventBroker.Publish(eventbroker.EventLogError, "Your Discord token is invalid! Please enter a valid token in the LB_TOKEN environment variable.")
			return err
		}

		token := tui.GetToken("Your Discord token is invalid!")

		cfg := d.configSvc.GetConfig()
		cfg.Account.Discord.Token = token

		if err := d.configSvc.SaveConfig(cfg); err != nil {
			d.eventBroker.Publish(eventbroker.EventLogError, "Save config error:"+err.Error())
			os.Exit(1)
		}

		_ = system.RestartExecution()

		return err
	}

	d.eventBroker.Publish(eventbroker.EventLogError, "An error has occurred while connecting to Discord: "+err.Error())

	return err
}

func (d *discordService) discordgoLogger() func(int, int, string, ...any) {
	return func(msgL int, caller int, format string, a ...any) {
		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "unknown event") && !d.verbose {
			return
		}

		pc, file, line, _ := runtime.Caller(caller + 1)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg = fmt.Sprintf("%s:%d %s: %s", file, line, name, msg)

		logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return
		}
		defer f.Close()

		_, _ = f.WriteString(msg)

		if msgL == discordgo.LogError {
			d.eventBroker.Publish(eventbroker.EventLogError, "discordgo error: "+msg)
		}
	}
}

func (d *discordService) SetHypesquad(houseID int) error {
	_, err := d.session.Request(
		"POST",
		"https://discord.com/api/v9/hypesquad/online",
		map[string]int{"house_id": houseID},
	)
	return err
}

func (d *discordService) LeaveHypesquad() error {
	_, err := d.session.Request("DELETE", "https://discord.com/api/v9/hypesquad/online", nil)
	return err
}

func (d *discordService) GetTokenInfo(token string) (*jester.Data, error) {
	resp, err := d.session.RequestWithBucketID(
		"GET",
		discordgo.EndpointUser("@me"),
		nil,
		discordgo.EndpointUser("@me"),
		map[string]string{
			"Authorization": token,
		},
	)
	if err != nil {
		return nil, err
	}

	data, err := jester.NewJson(resp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d *discordService) UpdateSettings(path string, settings any) error {
	_, err := d.session.Request("PATCH", discordgo.EndpointUser("@me")+path, settings)
	return err
}

func (d *discordService) SetFriendRequestsAllowed(allowed bool) error {
	if allowed {
		d.session.State.UserSettings.Privacy.FriendSourceFlags = wrapperspb.UInt32(14)
	} else {
		d.session.State.UserSettings.Privacy.FriendSourceFlags = wrapperspb.UInt32(0)
	}

	return d.session.UserSettingsUpdate(d.session.State.UserSettings, true)
}
