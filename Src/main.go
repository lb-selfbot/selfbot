package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"syscall"

	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/jpillora/overseer"
	"github.com/rgamba/evtwebsocket"

	"github.com/LightningDev1/LB-Selfbot/commands"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/events"
	"github.com/LightningDev1/LB-Selfbot/ipc"
	"github.com/LightningDev1/LB-Selfbot/utils"
)

func main() {
	overseer.Run(overseer.Config{
		Program: selfbot_main,
	})
}

func selfbot_main(overseer.State) {
	conn := evtwebsocket.Conn{
		// Fires when the connection is established
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected!")
		},
		// Fires when a new message arrives from the server
		OnMessage: func(msg []byte, w *evtwebsocket.Conn) {
			fmt.Printf("New message: %s\n", msg)
		},
		// Fires when an error occurs and connection is closed
		OnError: func(err error) {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		},
		// Ping interval in secs (optional)
		PingIntervalSecs: 5,
		// Ping message to send (optional)
		PingMsg: []byte("PING"),
	}

	err := conn.Dial("ws://127.0.0.1:8000/api/v3/ws/website", "")
	if err != nil {
		fmt.Println(err)
	}
	conn.Send(evtwebsocket.Msg{
		Body: []byte("Message body"),
		Callback: func(resp []byte, w *evtwebsocket.Conn) {
			// This function executes when the server responds
			fmt.Printf("Got response: %s\n", resp)
		},
	})

	fmt.Println("LightningBot v4.0.0 Alpha")
	rand.Seed(time.Now().UnixNano())
	//ipc.Ipc.Init()
	ipc.Ipc.OnRead = func(s string) {
		fmt.Println(s)
		event := ipc.JsonEvent{}
		json.Unmarshal([]byte(s), &event)
		switch event.Op {
		case ipc.Ipc.OpCodes.EXIT:
			os.Exit(0)
		case ipc.Ipc.OpCodes.RESTART:
			// Restart?
			syscall.Exec(os.Args[0], os.Args, os.Environ())
		}
	}

	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	events := events.Events{}

	session, err := discordgo.New(cfg.Account.Discord.Token)
	if err != nil {
		panic(err)
	}

	router := dgc.Create(&dgc.Router{
		Prefixes:         []string{cfg.Selfbot.CommandPrefix},
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
