package commands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type AnimationCommands struct{}

func (a AnimationCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "virus",
		Description: "Send an animated virus message",
		Usage:       "[p]virus <virus name>",
		Handler:     a.virus,
	})
}

var (
	virus = []string{
		"`[▓▓▓                    ] / %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓                ] - %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓           ] \\ %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓         ] | %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ] / %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓   ] - %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ] \\ %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] | %s.exe - Compiled`",
		"`[▓▓▓                    ] / %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓           ] - %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ] \\ %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ] | %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] / %s.exe - Downloaded`",
		"`[▓▓▓▓▓▓▓                ] - %s.exe - Injecting`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ] \\ %s.exe - Injecting`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] | %s.exe - Injected`",
		"`Successfully Injected %s.exe`",
	}
)

func (AnimationCommands) virus(ctx *dgc.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, fmt.Sprintf(virus[0], utils.Arguments["<virus name>"]))
	if err != nil {
		return
	}

	for i := 1; i < len(virus); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fmt.Sprintf(virus[i], utils.Arguments["<virus name>"]))
		if err != nil {
			return
		}
	}
}
