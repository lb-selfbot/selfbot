package console

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/common-nighthawk/go-figure"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/services/console/consoledefs"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/unsafeutil"
	"github.com/lb-selfbot/selfbot/internal/utils"
	"github.com/robert-nix/ansihtml"
	"golang.org/x/term"
)

const (
	commandAlignCenter = "ALIGN_CENTER"

	clearConsoleANSI = "\033[H\033[2J"
)

// Service handles the output console state.
type Service interface {
	GetHTMLConsole() string
	GetRawConsole() string

	Set(ansi string)
	Append(ansi string)
	Clear()
}

type consoleService struct {
	eventBroker eventbroker.Service

	htmlConsole string
	rawConsole  string
	mu          sync.RWMutex
}

// NewService creates a new ConsoleService instance.
func NewService(eventBroker eventbroker.Service) Service {
	svc := &consoleService{eventBroker: eventBroker}
	svc.subscribe()
	return svc
}

func (c *consoleService) subscribe() {
	c.eventBroker.Subscribe(eventbroker.EventSetConsole, func(payload any) {
		console, ok := payload.(consoledefs.Set)
		if !ok {
			return
		}

		c.mu.Lock()
		defer c.mu.Unlock()

		c.htmlConsole = console.HTML
		c.rawConsole = console.ANSI

		if constants.CONSOLE {
			c.clearAndPrintBanner()
			fmt.Println(c.processCommands(console.ANSI, commandAlignCenter, c.alignCenterANSIReplacer))
		}
	})

	c.eventBroker.Subscribe(eventbroker.EventAppendConsole, func(payload any) {
		console, ok := payload.(consoledefs.Append)
		if !ok {
			return
		}

		c.mu.Lock()
		defer c.mu.Unlock()

		c.htmlConsole += console.HTML
		c.rawConsole += console.ANSI

		if constants.CONSOLE {
			fmt.Println(c.processCommands(console.ANSI, commandAlignCenter, c.alignCenterANSIReplacer))
		}
	})

	c.eventBroker.Subscribe(eventbroker.EventClearConsole, func(any) {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.htmlConsole = ""
		c.rawConsole = ""

		if constants.CONSOLE {
			c.clearAndPrintBanner()
		}
	})
}

func (c *consoleService) clearAndPrintBanner() {
	ascii := figure.NewFigure("LightningBot", "stop", true).Slicify()
	centered := utils.String.CenterSlices(ascii)
	colored := logging.Colored(centered, 0, 203, 255, true) + "\n\n"

	fmt.Print(clearConsoleANSI)
	fmt.Println(colored)
}

func (c *consoleService) processCommands(console string, command string, replacer func(string) string) string {
	var builder strings.Builder
	start := 0
	for i := 0; i <= len(console); i++ {
		// Find line breaks or end of string
		if i == len(console) || console[i] == '\n' {
			line := console[start:i]
			if strings.HasPrefix(line, command) {
				line = replacer(line[len(command):])
			}
			builder.WriteString(line)
			if i < len(console) {
				builder.WriteByte('\n')
			}
			start = i + 1
		}
	}
	return builder.String()
}

func (c *consoleService) alignCenterHTMLReplacer(line string) string {
	return `<div style="text-align:center;">` + line + `</div>`
}

func (c *consoleService) alignCenterANSIReplacer(line string) string {
	// Calculate the length of the line without ANSI escape codes
	length := len(constants.ANSI_REGEX.ReplaceAllString(line, ""))

	// Get the terminal width
	fd := int(os.Stdout.Fd())
	w, _, err := term.GetSize(fd)
	if err != nil {
		return line
	}

	// Calculate the space count
	spaceCount := (w - length) / 2
	spaceCount = int(math.Max(float64(spaceCount), 0))

	return strings.Repeat(" ", spaceCount) + line
}

func (c *consoleService) ansiToHTML(ansi string) string {
	// Convert ANSI to HTML without copying the data
	html := unsafeutil.ByteSliceToString(ansihtml.ConvertToHTML(unsafeutil.StringToByteSlice(ansi)))

	return c.processCommands(html, commandAlignCenter, c.alignCenterHTMLReplacer)
}

func (c *consoleService) GetHTMLConsole() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.htmlConsole
}

func (c *consoleService) GetRawConsole() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.rawConsole
}

func (c *consoleService) Set(ansi string) {
	html := c.ansiToHTML(ansi)
	c.eventBroker.Publish(eventbroker.EventSetConsole, consoledefs.Set{ANSI: ansi, HTML: html})
}

func (c *consoleService) Append(ansi string) {
	html := "\n" + c.ansiToHTML(ansi)
	c.eventBroker.Publish(eventbroker.EventAppendConsole, consoledefs.Append{ANSI: ansi, HTML: html})
}

func (c *consoleService) Clear() {
	c.eventBroker.Publish(eventbroker.EventClearConsole, nil)
}
