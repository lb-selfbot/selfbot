package interfaces

type InfoEntry struct {
	Key   string
	Value string
}

type ConsoleMessage struct {
	uiSvc UIService

	Title      string
	Color      string
	WebhookURL string
	Important  bool
	Info       []InfoEntry
}

// NewConsoleMessage creates a new ConsoleMessage instance.
func NewConsoleMessage(uiSvc UIService) *ConsoleMessage {
	return &ConsoleMessage{uiSvc: uiSvc}
}

// Setter methods with method chaining.
func (c *ConsoleMessage) SetTitle(title string) *ConsoleMessage {
	c.Title = title
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

// AddInfo appends key-value pairs to the Info slice to maintain order.
func (c *ConsoleMessage) AddInfo(key, value string) *ConsoleMessage {
	c.Info = append(c.Info, InfoEntry{Key: key, Value: value})
	return c
}

// AddInfoConditional appends key-value pairs only if the condition is met.
func (c *ConsoleMessage) AddInfoConditional(key, value string, condition bool) *ConsoleMessage {
	if condition {
		c.AddInfo(key, value)
	}
	return c
}

func (c *ConsoleMessage) Show() error {
	return c.uiSvc.ShowConsoleMessage(c)
}

type UIService interface {
	Inject(authSvc AuthService, ipcSvc IPCService)

	ANSIColored(text string, r, g, b int, bold bool) string
	ShowBanner(flags int) error
	ShowNotification(message string) error

	LogError(format string, e ...any) error
	LogInfo(format string, e ...any) error

	AppendConsole(ansi string) error
	ClearConsole() error
	SetConsole(ansi string) error
	NewConsoleMessage() *ConsoleMessage
	ShowConsoleMessage(c *ConsoleMessage) error
}
