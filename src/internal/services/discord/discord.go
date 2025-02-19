package discord

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/internal/tui"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type DiscordServiceImpl struct {
	authSvc           interfaces.AuthService
	uiSvc             interfaces.UIService
	o11ySvc           interfaces.O11yService
	eventReactionsSvc interfaces.EventReactionsService
	extensionsSvc     interfaces.ExtensionsService

	session *discordgo.Session
	router  *dgc.Router

	verbose bool
}

func (d *DiscordServiceImpl) Inject(
	authSvc interfaces.AuthService,
	uiSvc interfaces.UIService,
	o11ySvc interfaces.O11yService,
	eventReactionsSvc interfaces.EventReactionsService,
	extensionsSvc interfaces.ExtensionsService,
) {
	d.authSvc = authSvc
	d.uiSvc = uiSvc
	d.o11ySvc = o11ySvc
	d.eventReactionsSvc = eventReactionsSvc
	d.extensionsSvc = extensionsSvc
}

func (d *DiscordServiceImpl) discordgoLogger(msgL int, caller int, format string, a ...interface{}) {
	pc, file, line, _ := runtime.Caller(caller + 1)

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]

	msg := fmt.Sprintf(format, a...)

	if strings.Contains(msg, "unknown event") && !d.verbose {
		return
	}

	msg = fmt.Sprintf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)

	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.WriteString(msg)

	if msgL == discordgo.LogError {
		d.uiSvc.LogError("discordgo error: " + msg)
	}
}

func (d *DiscordServiceImpl) Init(token string, verbose bool) error {
	d.verbose = verbose

	discordgo.Logger = d.discordgoLogger

	d.session, _ = discordgo.New(token)

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	d.session.Token = token
	d.session.State.TrackVoice = false
	d.session.State.MaxMessageCount = 50
	d.session.SetIdentify(cfg.GetClient())
	d.session.ErrorHandler = d.o11ySvc.DiscordErrorHandler
	d.session.ForceSubscribeGuildIDs = cfg.Selfbot.ForceSubscribedServers

	if verbose {
		d.session.LogLevel = discordgo.LogDebug
	} else {
		d.session.LogLevel = discordgo.LogWarning
	}

	d.router = dgc.Create(&dgc.Router{
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

			if ctx.Event.Author.ID == d.session.State.User.ID {
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
				d.uiSvc.LogDangerousCmdBlocked(ctx.Event.Author, ctx.Command.Name)
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
						d.uiSvc.LogError("Failed to get config")
						return
					}

					if cfg.Selfbot.DeleteCommands && ctx.Event.ID != "0" {
						go func() {
							_ = ctx.Session.ChannelMessageDelete(ctx.Event.Message.ChannelID, ctx.Event.Message.ID)
						}()
					}

					utils.Arguments = arguments

					if ctx.Command.Name != "repeat" {
						d.eventReactionsSvc.AddUsedCommand(ctx)
					}

					d.uiSvc.LogCmdUsed(ctx.Command.Name)

					d.authSvc.CheckNonce()

					following(ctx)
				}
			},
		},
	})

	d.eventReactionsSvc.Register()

	commands.Commands.RegisterCommands(d.session, d.router, false)

	d.extensionsSvc.Register()
	
	commands.Commands.HelpCommand.RegisterHelpCommand(d.router, d.session)

	d.router.Initialize(d.session)

	return nil
}

func (d *DiscordServiceImpl) GetSession() *discordgo.Session { return d.session }
func (d *DiscordServiceImpl) GetRouter() *dgc.Router         { return d.router }

func (d *DiscordServiceImpl) Open() error {
	err := d.session.Open()
	if err == nil {
		return nil
	}

	// Check if token is invalid
	if strings.Contains(err.Error(), "Authentication failed") {
		if !constants.CONSOLE {
			d.uiSvc.LogError("Your Discord token is invalid! Please re-enter it in the Config tab.")
			return err
		}

		if _, ok := os.LookupEnv("LB_TOKEN"); ok {
			d.uiSvc.LogError("Your Discord token is invalid! Please enter a valid token in the LB_TOKEN environment variable.")
			return err
		}

		token := tui.GetToken("Your Discord token is invalid!")

		cfg, err := config.GetConfig()
		if err != nil {
			d.uiSvc.LogError("Error loading config:", err)
			os.Exit(1)
		}

		cfg.Account.Discord.Token = token
		if err := cfg.Save(); err != nil {
			d.uiSvc.LogError("Save config error:", err)
			os.Exit(1)
		}

		_ = utils.System.Restart()

		return err
	}

	d.uiSvc.LogError("An error has occurred while connecting to Discord: " + err.Error())

	return err
}
