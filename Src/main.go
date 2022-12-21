package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"syscall"

	"github.com/LightningDev1/LB-Selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/consolelogging"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/multitokensniper"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

func initArgs() {
	args := strings.Join(os.Args, " ")

	// constants.Console must be true if "-l" flag is not set
	constants.CONSOLE = !strings.Contains(args, "-l")

	if strings.Contains(args, "-v") {
		fmt.Println(constants.VERSION_STRING[1:])
		os.Exit(0)
	}
}

func main() {
	initArgs()

	if constants.DEBUG {
		consolelogging.Warn("Debug Enabled!")
	}

	//startParams := api.GetStartParams()
	rand.Seed(time.Now().UnixNano())
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	if !constants.CONSOLE {
		ipc.IPC.Init()
	}

	err = websocket.Init(cfg.Account.LightningBot.Username, cfg.Account.LightningBot.Password+"a")
	if err != nil {
		consolelogging.Error(err)
		if !constants.DEBUG {
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

			if !utils.SliceContains(cfg.Sharing.SharedUsers, ctx.Event.Author.ID) &&
				!utils.SliceContains(cfg.Sharing.SharedUsers, ctx.Event.Author.String()) {
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

	if !constants.DEBUG {
		if !constants.CONSOLE {
			ipc.IPC.WaitForConnection()
		}
		websocket.WaitForConnection()
	}

	err = session.Open()
	if err != nil {
		if strings.Contains(err.Error(), "Authentication failed") {
			consolelogging.Error("Invalid token!")
		} else {
			consolelogging.Error(err)
		}

		if !constants.DEBUG {
			os.Exit(0)
		}
	}

	multitokensniper.Init(cfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
