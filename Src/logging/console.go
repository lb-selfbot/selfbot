package logging

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/common-nighthawk/go-figure"
)

type ConsoleLogging struct{}

func (ConsoleLogging) Debug(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[35;1mDEBUG\x1b[0m] ", t)
	fmt.Println(v...)
}

func (ConsoleLogging) Info(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[32;1mINFO\x1b[0m] ", t)
	fmt.Println(v...)
}

func (ConsoleLogging) Warn(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[33;1mWARNING\x1b[0m] ", t)
	fmt.Println(v...)
}

func (ConsoleLogging) Error(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[31;1mERROR\x1b[0m] ", t)
	fmt.Println(v...)
}

func (ConsoleLogging) Colored(text string, r, g, b int, bold bool) string {
	boldArg := ";1"
	if !bold {
		boldArg = ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%d%sm%s\033[0m", r, g, b, boldArg, text)
}

func (ConsoleLogging) Clear() {
	//fmt.Print("\033[H\033[2J")
}

func (ConsoleLogging) PrintBanner() {
	banner := utils.String.CenterSlices(figure.NewFigure("LightningBot", "stop", true).Slicify())
	message := Console.Colored(banner, 0, 203, 255, true) + "\n\n"

	fmt.Println(message)
}

var Console = ConsoleLogging{}
