package main

import (
	"encoding/json"
	"fmt"
	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/jpillora/overseer"

	"github.com/LightningDev1/LB-Selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/websocket"
)

func main() {
	overseer.Run(overseer.Config{
		Program: selfbotMain,
	})
}

func selfbotMain(overseer.State) {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	err = websocket.Init(cfg.Account.LightningBot.Username, cfg.Account.LightningBot.Password)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("LightningBot v4.0.0 Alpha")
	rand.Seed(time.Now().UnixNano())
	//ipc.Ipc.Init()
	ipc.Ipc.OnRead = func(s string) {
		fmt.Println(s)
		event := ipc.JsonEvent{}
		err := json.Unmarshal([]byte(s), &event)
		if err != nil {
			return
		}
		switch event.Op {
		case ipc.Ipc.OpCodes.EXIT:
			os.Exit(0)
		case ipc.Ipc.OpCodes.RESTART:
			_ = utils.Process.Restart()
		}
	}

	session, err := discordgo.New(cfg.Account.Discord.Token)
	if err != nil {
		panic(err)
	}

	router := dgc.Create(&dgc.Router{
		PrefixFunc: func() []string {
			cfg, err := config.GetConfig()
			if err != nil {
				return []string{"!"}
			}
			return []string{cfg.Selfbot.CommandPrefix}
		},
		IgnorePrefixCase: true,
		SelfBot:          true,
		Commands:         []*dgc.Command{},
		Middlewares:      []dgc.Middleware{},
	})

	eventhandlers.EventHandlers.RegisterHandlers(session, router)
	commands.Commands.RegisterCommands(session, router)

	router.Initialize(session)

	err = session.Open()
	if err != nil {
		panic(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}
