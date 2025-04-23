package tui

import (
	"context"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lb-selfbot/selfbot/internal/browser"
	"github.com/lb-selfbot/selfbot/internal/logging"
	"github.com/lb-selfbot/selfbot/internal/savedtokens"
)

func GetToken(clear func(), title ...string) string {
	title = append(title, "How would you like to login to Discord?")

	items := []list.Item{
		ListItem{"Login with a token", 1},
		ListItem{"Login with a web browser", 2},
		ListItem{"Login with a saved account from the Discord app", 3},
	}

	model := NewListModel(items, strings.Join(title, " "))

	if _, err := tea.NewProgram(model).Run(); err != nil {
		logging.Error("Error running program:", err)
		os.Exit(1)
	}

	var token string

	switch model.Choice.Index {
	case 1:
		token = loginWithToken()
	case 2:
		token = loginWithBrowser(clear)
	case 3:
		token = loginWithSavedToken(clear)
	}

	if token == "" {
		return GetToken(clear, "You must enter a token!")
	}

	clear()

	return token
}

func loginWithToken() string {
	model := NewTextInputModel("Please enter your Discord token:", "", true)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		logging.Error("Error running program:", err)
		os.Exit(1)
	}

	return model.Value()
}

func loginWithBrowser(clear func()) string {
	model := NewLoadingModel("Please login to Discord in the browser window.")
	program := tea.NewProgram(model)

	go func() {
		_, _ = program.Run()
	}()

	token, err := browser.RunWebLogin(context.Background())
	if err != nil {
		model.Quit()
		program.Quit()
		clear()

		logging.Error("Error logging in:", err, "\n")

		return ""
	}

	model.Quit()
	program.Quit()
	clear()

	return token
}

func loginWithSavedToken(clear func()) string {
	model := NewLoadingModel("Loading saved accounts...")
	program := tea.NewProgram(model)

	go func() {
		_, _ = program.Run()
	}()

	savedTokens := savedtokens.GetSavedTokens()
	if len(savedTokens) == 0 {
		model.Quit()
		program.Quit()
		clear()

		logging.Error("You have no saved accounts!")

		return ""
	}

	model.Quit()
	program.Quit()
	clear()

	items := make([]list.Item, len(savedTokens))

	for i, token := range savedTokens {
		items[i] = ListItem{token.Username + " - " + token.ID, i + 1}
	}

	accListModel := NewListModel(items, "Please select a saved account:")
	if _, err := tea.NewProgram(accListModel).Run(); err != nil {
		logging.Error("Error running program:", err)
		os.Exit(1)
	}

	return savedTokens[accListModel.Choice.Index-1].Token
}
