package tui

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

var mainColor = lipgloss.Color("#00cbff")

func reset() {
	time.Sleep(time.Millisecond * 500)

	logging.Console.Clear()
	logging.Console.PrintBanner()
}
