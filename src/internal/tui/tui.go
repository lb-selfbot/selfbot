package tui

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

var (
	mainColor = lipgloss.Color("#00cbff")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(mainColor)

	// Text styles
	titleStyle = lipgloss.NewStyle().
			Foreground(mainColor).
			Bold(true).
			Margin(1)

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	// List styles
	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(mainColor).
				Bold(true)
)

func reset() {
	time.Sleep(time.Millisecond * 500)
	logging.Console.Clear()
	logging.Console.PrintBanner()
}
