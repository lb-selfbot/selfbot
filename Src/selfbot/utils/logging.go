package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/robert-nix/ansihtml"
)

type LoggingUtils struct{}

func (LoggingUtils) Debug(v ...any) {
	if !constants.DEBUG {
		return
	}

	t := time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[35;1mDEBUG\x1b[0m] ", t)
	fmt.Println(v...)
}

func (LoggingUtils) Info(v ...any) {
	t := time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[32;1mINFO\x1b[0m] ", t)
	fmt.Println(v...)
}

func (LoggingUtils) Warn(v ...any) {
	t := time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[33;1mWARNING\x1b[0m] ", t)
	fmt.Println(v...)
}

func (LoggingUtils) Error(v ...any) {
	t := time.Now().Local().Format("2006-01-02 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[31;1mERROR\x1b[0m] ", t)
	fmt.Println(v...)
}

func (LoggingUtils) AnsiToHtml(text string) string {
	html := string(ansihtml.ConvertToHTML([]byte(text)))

	formattedHtml := Map(strings.Split(html, "\n"), func(line string) string {
		if strings.Contains(line, "ALIGN_CENTER") {
			line = strings.ReplaceAll(line, "ALIGN_CENTER", "")
			line = "<div style=\"text-align:center;\">" + line + "</div>"
		}

		return line
	})

	return strings.Join(formattedHtml, "\n")
}

var Logging = &LoggingUtils{}
