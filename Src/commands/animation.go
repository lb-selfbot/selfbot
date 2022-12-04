package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/embed"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/common-nighthawk/go-figure"
)

type AnimationCommands struct{}

func (a AnimationCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Animation", "Send animated messages")
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
		Name:        "fap",
		Description: "Send a fap animation",
		Usage:       "[p]fap",
		Handler:     a.fap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "countup",
		Description: "Count up to a number",
		Usage:       "[p]countup <number>",
		Handler:     a.countup,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "countdown",
		Description: "Count down to a number",
		Usage:       "[p]countdown <number>",
		Handler:     a.countdown,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "asciianimation",
		Description: "Send an animated ascii message",
		Usage:       "[p]asciianimation <text>",
		Handler:     a.asciianimation,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "customanimationhelp",
		Description: "View custom animation help",
		Usage:       "[p]customanimationhelp",
		Handler:     a.customanimationhelp,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "customanimations",
		Description: "View all custom animations",
		Usage:       "[p]customanimations",
		Handler:     a.customanimations,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "customanimation",
		Description: "Invoke a custom animation",
		Usage:       "[p]customanimation <name>",
		Handler:     a.customanimation,
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
	fap = []string{
		"8:punch:=====Đ",
		"8=:punch:====Đ",
		"8==:punch:===Đ",
		"8===:punch:==Đ",
		"8====:punch:=Đ",
		"8=====:punch:Đ",
		"8====:punch:=Đ",
		"8===:punch:==Đ",
		"8==:punch:===Đ",
		"8=:punch:====Đ",
		"8:punch:=====Đ",
		"8=:punch:====Đ",
		"8==:punch:===Đ",
		"8===:punch:==Đ",
		"8====:punch:=Đ",
		"8=====:punch:Đ",
		"8=====Đ:sweat_drops:",
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

	for i := 1; i < len(bomb); i++ {
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

	for i := 1; i < len(fuckYou); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fuckYou[i])
		if err != nil {
			return
		}
	}
}

func (AnimationCommands) fap(ctx *dgc.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, fap[0])
	if err != nil {
		return
	}

	for i := 1; i < len(fap); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fap[i])
		if err != nil {
			return
		}
	}
}

func (AnimationCommands) countup(ctx *dgc.Ctx) {
	number, err := strconv.Atoi(utils.Arguments["<number>"])
	if err != nil && utils.Arguments["<number>"] != "" {
		_ = ctx.RespondText("Argument <number> must be a number.")
		return
	}
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "1")
	if err != nil {
		return
	}

	for i := 2; i <= number; i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fmt.Sprint(i))
		if err != nil {
			return
		}
	}
}

func (AnimationCommands) countdown(ctx *dgc.Ctx) {
	number, err := strconv.Atoi(utils.Arguments["<number>"])
	if err != nil && utils.Arguments["<number>"] != "" {
		_ = ctx.RespondText("Argument <number> must be a number.")
		return
	}
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, fmt.Sprint(number))
	if err != nil {
		return
	}

	for i := number; i > 0; i-- {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, fmt.Sprint(i-1))
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

func (AnimationCommands) customanimationhelp(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Creating an animation", fmt.Sprintf(
			"To create an animation, create a text file in %s. Every line in the file will be a frame of the animation.\n",
			filepath.Join(utils.File.GetMainFolder(), "Animations"),
		)).
		AddField("Using an animation", fmt.Sprintf(
			"To use an animation, use the command `%scustomanimation <filename without .txt>`.\n",
			ctx.Router.GetPrefixes()[0],
		)).
		Send(ctx)
}

func (AnimationCommands) customanimations(ctx *dgc.Ctx) {
	folder := filepath.Join(utils.File.GetMainFolder(), "Animations")
	_ = os.MkdirAll(folder, os.ModePerm)

	customAnimString := ""
	_ = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		nameSplit := strings.Split(info.Name(), ".")
		if len(nameSplit) != 2 {
			return nil
		}
		if nameSplit[1] != "txt" {
			return nil
		}
		customAnimString += " - " + nameSplit[0] + "\n"
		return nil
	})

	embed.NewEmbed().AddField("Custom Animations", customAnimString).Send(ctx)
}

func (AnimationCommands) customanimation(ctx *dgc.Ctx) {
	name := utils.Arguments["<name>"]
	folder := filepath.Join(utils.File.GetMainFolder(), "Animations")
	_ = os.MkdirAll(folder, os.ModePerm)
	path := folder + utils.File.SanitizeFilename(name+".txt")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = ctx.RespondText("Animation does not exist.")
		return
	}
	rawAnimation, ok := utils.File.Read(path)
	if !ok {
		_ = ctx.RespondText("Animation could not be read.")
		return
	}
	animationFrames := strings.Split(rawAnimation, "\n")
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, animationFrames[0])
	if err != nil {
		return
	}
	for i := 1; i < len(animationFrames); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, message.ID, animationFrames[i])
		if err != nil {
			return
		}
	}
}
