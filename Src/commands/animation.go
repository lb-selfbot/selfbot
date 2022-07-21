package commands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/common-nighthawk/go-figure"
)

type AnimationCommands struct{}

func (a AnimationCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "virus",
		Description: "Send an animated virus message",
		Usage:       "[p]virus <virus name>",
		Handler:     a.virus,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "100",
		Description: "Count to 100",
		Usage:       "[p]100",
		Handler:     a._100,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bomb",
		Description: "Send an animated bomb message",
		Usage:       "[p]bomb",
		Handler:     a.bomb,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fuckyou",
		Description: "Send an animated fuck you message",
		Usage:       "[p]fuckyou",
		Handler:     a.fuckyou,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "asciianimation",
		Description: "Send an animated ascii message",
		Usage:       "[p]asciianimation <text>",
		Handler:     a.asciianimation,
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
	bomb = []string{
		":bomb: ---------------- :fire:",
		":bomb: --------------- :fire:",
		":bomb: -------------- :fire:",
		":bomb: ------------- :fire:",
		":bomb: ------------ :fire:",
		":bomb: ----------- :fire:",
		":bomb: ---------- :fire:",
		":bomb: --------- :fire:",
		":bomb: -------- :fire:",
		":bomb: ------- :fire:",
		":bomb: ------ :fire:",
		":bomb: ----- :fire:",
		":bomb: ---- :fire:",
		":bomb: --- :fire:",
		":bomb: -- :fire:",
		":bomb: - :fire:",
		":bomb:  :fire:",
		":boom:",
	}
	fuckYou = []string{
		"**F**",
		"**FU**",
		"**FUC**",
		"**FUCK**",
		"**FUCK Y**",
		"**FUCK YO**",
		"**FUCK YOU**",
		"**FUCK YOU!**",
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

func (AnimationCommands) _100(ctx *dgc.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "1")
	if err != nil {
		return
	}

	for i := 2; i <= 100; i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fmt.Sprint(i))
		if err != nil {
			return
		}
	}
}

func (AnimationCommands) bomb(ctx *dgc.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, bomb[0])
	if err != nil {
		return
	}

	for i := 1; i < len(virus); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, bomb[i])
		if err != nil {
			return
		}
	}
}

func (AnimationCommands) fuckyou(ctx *dgc.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, fuckYou[0])
	if err != nil {
		return
	}

	for i := 1; i < len(virus); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fuckYou[i])
		if err != nil {
			return
		}
	}
}

func (AnimationCommands) asciianimation(ctx *dgc.Ctx) {
	text := utils.Arguments["<text>"]
	firstCharAscii := figure.NewFigure(string(text[0]), "", false).String()
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, fmt.Sprintf("```%s```", firstCharAscii))
	if err != nil {
		return
	}
	currentText := string(text[0])
	for index, char := range text {
		if index == 0 {
			continue
		}
		time.Sleep(time.Second * 1)
		currentText = currentText + string(char)
		ascii := figure.NewFigure(currentText, "", false).String()
		_, err := ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fmt.Sprintf("```%s```", ascii))
		if err != nil {
			return
		}
	}
}

// Custom Animations
