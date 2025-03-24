package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lb-selfbot/selfbot/internal/logging"
)

var stdInReader = bufio.NewReader(os.Stdin)

func PrintPrompt(prompt string) {
	t := time.Now().Local().Format("02/01/2006 15:04:05")
	fmt.Printf("[\x1b[34;1m%s\x1b[0m] [\x1b[35;1mINPUT\x1b[0m] %s: ", t, prompt)
}

func GetString(prompt string) string {
	PrintPrompt(prompt)

	input, err := stdInReader.ReadString('\n')
	if err != nil {
		logging.Error("Error reading input:", err)
		return ""
	}

	return strings.TrimSpace(input)
}
