package logging

import (
	"fmt"
	"time"

	"github.com/lb-selfbot/selfbot/internal/constants"
)

func Debug(v ...any) {
	if !constants.DEBUG {
		return
	}

	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[35;1mDEBUG\x1b[0m] ", t)
	fmt.Println(v...)
}

func Info(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[32;1mINFO\x1b[0m] ", t)
	fmt.Println(v...)
}

func Warn(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[33;1mWARNING\x1b[0m] ", t)
	fmt.Println(v...)
}

func Error(v ...any) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[31;1mERROR\x1b[0m] ", t)
	fmt.Println(v...)
}

func Colored(text string, r, g, b int, bold bool) string {
	boldArg := ";1"
	if !bold {
		boldArg = ""
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%d%sm%s\033[0m", r, g, b, boldArg, text)
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}
