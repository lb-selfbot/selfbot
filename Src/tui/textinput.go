package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type textInputModel struct {
	textInput textinput.Model
	prompt    string
	quitting  bool
	exiting   bool
	err       error
}

func NewTextInputModel(prompt, placeholder string, hidden bool) *textInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = listWidth

	if hidden {
		ti.EchoMode = textinput.EchoPassword
		ti.EchoCharacter = '•'
	}

	return &textInputModel{
		textInput: ti,
		prompt:    prompt,
		err:       nil,
	}
}

func (m *textInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.quitting = true

			return m, tea.Quit
		case tea.KeyCtrlC:
			m.exiting = true

			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *textInputModel) View() string {
	if m.exiting {
		os.Exit(0)
	}

	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"\n  %s\n\n%s",
		titleStyle.Render(m.prompt),
		itemStyle.Render(m.textInput.View()),
	) + "\n"
}

func (m *textInputModel) Value() string {
	return m.textInput.Value()
}
