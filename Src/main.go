package main

import (
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"syscall"

	"github.com/LightningDev1/LB-Selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/jpillora/overseer"
)

func main() {
	overseer.Run(overseer.Config{
		Program: selfbotMain,
	})
}

func selfbotMain(s overseer.State) {
	if constants.Debug {
		logging.Warn("Debug Enabled!")
	}
	//startParams := api.GetStartParams()
	rand.Seed(time.Now().UnixNano())
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	ipc.IPC.Init()

	err = websocket.Init(cfg.Account.LightningBot.Username, cfg.Account.LightningBot.Password)
	if err != nil {
		logging.Error(err)
		if !constants.Debug {
			os.Exit(0)
		}
	}

	session, _ := discordgo.New(cfg.Account.Discord.Token)

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
			if ctx.Event.Author.ID == session.State.User.ID {
				return true
			}
			cfg, err := config.GetConfig()
			if err != nil {
				return false
			}

			if !utils.SliceContains(cfg.Sharing.SharedUsers, ctx.Event.Author.ID) {
				return false
			}

			if cfg.Sharing.ShareAll || utils.SliceContains([]string{"help", "search", "repeat"}, ctx.Command.Name) {
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
			return utils.SliceContains(cfg.Sharing.SharedCommands, ctx.Command.Name)
		},
		Commands:    []*dgc.Command{},
		Middlewares: []dgc.Middleware{},
	})

	eventhandlers.EventHandlers.RegisterHandlers(session, router)
	commands.Commands.RegisterCommands(session, router)

	router.Initialize(session)

	if !constants.Debug {
		ipc.IPC.WaitForConnection()
		websocket.WaitForConnection()
	}

	err = session.Open()
	if err != nil {
		logging.Error(err)
	}

	multitokensniper.Init(cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
