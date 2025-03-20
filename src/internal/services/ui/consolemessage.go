package ui

import "fmt"

// InfoEntry represents a key-value pair for console messages.
type InfoEntry struct {
	Key   string
	Value string
}

// ConsoleMessage represents a message to be displayed in the console.
type ConsoleMessage struct {
	Title      string
	Color      string
	WebhookURL string
	Important  bool
	Info       []InfoEntry

	showCallback func(*ConsoleMessage) error
}

// NewConsoleMessage creates a new ConsoleMessage instance.
func NewConsoleMessage(showCallback func(*ConsoleMessage) error) *ConsoleMessage {
	return &ConsoleMessage{showCallback: showCallback}
}

// Setter methods with method chaining.
func (c *ConsoleMessage) SetTitle(format string, a ...any) *ConsoleMessage {
	c.Title = fmt.Sprintf(format, a...)
	return c
}

func (c *ConsoleMessage) SetColor(color string) *ConsoleMessage {
	c.Color = color
	return c
}

func (c *ConsoleMessage) SetWebhookURL(webhookURL string) *ConsoleMessage {
	c.WebhookURL = webhookURL
	return c
}

func (c *ConsoleMessage) SetImportant(important bool) *ConsoleMessage {
	c.Important = important
	return c
}

// AddInfo appends key-value pairs to the console message with optional conditions.
func (c *ConsoleMessage) AddInfo(key, value string, conds ...bool) *ConsoleMessage {
	if len(conds) < 1 || conds[0] {
		c.Info = append(c.Info, InfoEntry{Key: key, Value: value})
	}
	return c
}

func (c *ConsoleMessage) Show() error {
	return c.showCallback(c)
}
