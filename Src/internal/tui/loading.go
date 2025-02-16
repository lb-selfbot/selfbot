package tui

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loadingModel struct {
	spinner  spinner.Model
	text     string
	quitting bool
	exiting  bool
	err      error
}

func NewLoadingModel(text string) *loadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(mainColor)
	return &loadingModel{spinner: s, text: text}
}

func (m *loadingModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *loadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.exiting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m *loadingModel) View() string {
	if m.exiting {
		os.Exit(0)
	}
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s %s\n\n", m.spinner.View(), titleStyle.Render(m.text))
	if m.quitting {
		return ""
	}
	return str
}

func (m *loadingModel) Quit() {
	m.quitting = true
}
