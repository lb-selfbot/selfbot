package console

import "sync"

// Service handles the output console state.
type Service interface {
	GetHTMLConsole() string
	SetHTMLConsole(html string)

	GetRawConsole() string
	SetRawConsole(ansi string)
}

type consoleService struct {
	htmlConsole string
	rawConsole  string
	mu          sync.Mutex
}

// NewConsoleService creates a new ConsoleService instance.
func NewConsoleService() Service {
	return &consoleService{}
}

func (c *consoleService) GetHTMLConsole() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.htmlConsole
}

func (c *consoleService) SetHTMLConsole(html string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.htmlConsole = html
}

func (c *consoleService) GetRawConsole() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.rawConsole
}

func (c *consoleService) SetRawConsole(ansi string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rawConsole = ansi
}
