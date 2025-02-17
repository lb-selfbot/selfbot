package discord

import (
	"fmt"
	"slices"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/data"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/extensions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type DiscordServiceImpl struct {
	authSvc interfaces.AuthService
	uiSvc   interfaces.UIService

	session *discordgo.Session
	router  *dgc.Router
}

func (d *DiscordServiceImpl) Inject(authSvc interfaces.AuthService, uiSvc interfaces.UIService) {
	d.authSvc = authSvc
	d.uiSvc = uiSvc
}

func (d *DiscordServiceImpl) Init(token string, verbose bool) error {
	d.session, _ = discordgo.New(token)

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	d.session.Token = token
	d.session.State.TrackVoice = false
	d.session.State.MaxMessageCount = 50
	d.session.SetIdentify(cfg.GetClient())
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
				// d.uiSvc.BlockedDangerousCommand(ctx.Event.Author, ctx.Command.Name)
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
						_ = d.uiSvc.LogError("Failed to get config")
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

					// d.uiSvc.CommandUsed(ctx.Command.Name)

					d.authSvc.CheckNonce()

					following(ctx)
				}
			},
		},
	})

	eventreactions.RegisterHandlers(d.session, d.router)
	commands.Commands.RegisterCommands(d.session, d.router, false)
	extensions.RegisterExtensions(d.session, d.router)
	commands.Commands.HelpCommand.RegisterHelpCommand(d.router, d.session)

	d.router.Initialize(d.session)

	return nil
}

func (d *DiscordServiceImpl) GetSession() *discordgo.Session {
	return d.session
}

func (d *DiscordServiceImpl) GetRouter() *dgc.Router {
	return d.router
}

func (d *DiscordServiceImpl) Open() error {
	return d.session.Open()
}
