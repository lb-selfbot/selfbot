package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"syscall"

	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/jpillora/overseer"

	"github.com/LightningDev/LB-Selfbot/commands"
	"github.com/LightningDev/LB-Selfbot/events"
	"github.com/LightningDev/LB-Selfbot/ipc"
	"github.com/LightningDev/LB-Selfbot/utils"
)

func main() {
	overseer.Run(overseer.Config{
		Program: selfbot_main,
	})
}

func selfbot_main(overseer.State) {
	fmt.Println("LightningBot v4.0.0 Alpha")
	utils.SeedRandom()
	//ipc.Ipc.Init()
	ipc.Ipc.OnRead = func(s string) {
		fmt.Println(s)
		event := ipc.JsonEvent{}
		json.Unmarshal([]byte(s), &event)
		switch event.Op {
		case ipc.Ipc.OpCodes.EXIT:
			os.Exit(0)
		case ipc.Ipc.OpCodes.RESTART:
			// How tf
			os.Exit(0)
		}
	}

	events := events.Events{}

	// Open a simple Discord session
	token := "NzMzOTMzNzk5ODU5MDI4MDA4.YffMNA.UvWJVM3wEgtIGWpfwxhyO_D_PN8"
	session, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}

	// Create a dgc router
	router := dgc.Create(&dgc.Router{
		Prefixes: []string{
			"!",
		},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		SelfBot:          true,
		Commands:         []*dgc.Command{},
		Middlewares:      []dgc.Middleware{},
	})

	router.RegisterMiddleware(func(following dgc.ExecutionHandler) dgc.ExecutionHandler {
		return func(ctx *dgc.Ctx) {
			events.OnCommand(ctx)
			arguments, success := utils.GetArguments(ctx)
			utils.Arguments = arguments
			if success {
				following(ctx)
			}
		}
	})

	commands := commands.Commands{}
	commands.RegisterCommands(router)

	session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		fmt.Println("Logged in as:", event.User.Username)
		fmt.Println(len(router.Commands), "commands")
	})

	router.Initialize(session)

	err = session.Open()
	if err != nil {
		panic(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
}
