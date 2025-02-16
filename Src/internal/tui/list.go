package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	listWidth  = 60
	listHeight = 14
)

// List item
type ListItem struct {
	string
	Index int
}

func (i ListItem) FilterValue() string { return "" }

// Item delegate
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ListItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.string)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// List model
type listModel struct {
	list     list.Model
	Choice   ListItem
	quitting bool
}

func NewListModel(items []list.Item, title string) *listModel {
	model := &listModel{}

	model.list = list.New(items, itemDelegate{}, listWidth, listHeight)
	model.list.Title = title
	model.list.SetShowStatusBar(false)
	model.list.SetFilteringEnabled(false)
	model.list.SetShowHelp(false)
	model.list.Styles.Title = titleStyle
	model.list.Styles.FilterPrompt = inputStyle
	model.list.Styles.FilterCursor = inputStyle
	model.list.SetShowFilter(false)

	return model
}

func (m *listModel) Init() tea.Cmd {
	return nil
}

func (m *listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(ListItem)
			if ok {
				m.Choice = i
			} else {
				m.quitting = true
			}
			return m, tea.Quit

		case "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *listModel) View() string {
	if m.quitting {
		os.Exit(0)
	}

	if m.Choice.Index == 0 {
		return baseStyle.Render("\n" + m.list.View())
	}

	return ""
}
