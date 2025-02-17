package startup

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/assets"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/extensions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

// Initialize sets up basic requirements for both desktop and hosting modes
func Initialize(verbose bool) error {
	if err := utils.IO.CreateFolders(); err != nil {
		return fmt.Errorf("create folders: %w", err)
	}

	// Initialize icon
	iconFile := filepath.Join(utils.IO.GetMainFolder(), "Logo.png")
	if !utils.IO.Exists(iconFile) && !constants.HOSTING {
		icon, err := assets.Assets.ReadFile("Logo.png")
		if err != nil {
			logging.UI.Error("Error reading icon file: " + err.Error())
		}

		_ = utils.IO.WriteBytes(iconFile, icon)
	}

	// Initialize logging
	logPath := filepath.Join(utils.IO.GetMainFolder(), "LB.log")
	if err := utils.IO.Write(logPath, ""); err != nil {
		return fmt.Errorf("clear log file: %w", err)
	}

	discordgo.Logger = func(msgL, caller int, format string, a ...any) {
		pc, file, line, _ := runtime.Caller(caller + 1)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		if strings.Contains(msg, "unknown event") && !verbose {
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

	return nil
}

// CreateSession creates and configures a Discord session
func CreateSession(token string, verbose bool) (*discordgo.Session, error) {
	session, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	session.Token = token
	session.State.TrackVoice = false
	session.State.MaxMessageCount = 50
	session.SetIdentify(cfg.GetClient())
	session.ForceSubscribeGuildIDs = cfg.Selfbot.ForceSubscribedServers

	if verbose {
		session.LogLevel = discordgo.LogDebug
	} else {
		session.LogLevel = discordgo.LogWarning
	}

	return session, nil
}

// CreateRouter creates a new DGC router
func CreateRouter(session *discordgo.Session) *dgc.Router {
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

					// TODO: Check nonce
					// websocket.Websocket.CheckNonce()

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
