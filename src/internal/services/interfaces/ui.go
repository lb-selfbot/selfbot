package interfaces

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/discordgo"
)

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
	Inject(authSvc AuthService, ipcSvc IPCService, discordSvc DiscordService, embedSvc EmbedService)

	ANSIColored(text string, r, g, b int, bold bool) string
	ShowBanner(flags int) error
	ShowNotification(format string, e ...any) error

	LogError(format string, e ...any)
	LogInfo(format string, e ...any)

	AppendConsole(ansi string) error
	ClearConsole() error
	SetConsole(ansi string) error
	NewConsoleMessage() *ConsoleMessage
	ShowConsoleMessage(c *ConsoleMessage) error

	LogCmdUsed(cmdName string)
	LogReady(event *discordgo.Ready, numCmds int)
	LogNitroSniper(msg *discordgo.Message, elapsed time.Duration, status, code, account string)
	LogNitroSniperFake(msg *discordgo.Message, code, account string)
	LogPrivnoteSniper(msg *discordgo.Message, elapsed time.Duration, status, path, account string)
	LogGiveawaySniperFound(msg *discordgo.Message, prize, host, endsIn, account string)
	LogGiveawaySniperWon(msg *discordgo.Message, prize, account string)
	LogPokeTwoSniper(msg *discordgo.Message, elapsed time.Duration, confidence float64, pokemon, received, iv, account string)
	LogBanNotifier(guildID string)
	LogDMNotifier(msg *discordgo.Message)
	LogDMEditNotifier(msg *discordgo.Message, before, after string)
	LogDMDeleteNotifier(msg *discordgo.Message, attachments []string)
	LogDMTypingNotifier(username string)
	LogMessageDeleteNotifier(msg *discordgo.Message, attachments []string)
	LogDiscordStaffDetection(msg *discordgo.Message)
	LogGhostPingDetection(msg *discordgo.Message)
	LogSelfbotDetection(msg *discordgo.Message, selfbot string)
	LogTicketNotifier(channel *discordgo.Channel)
	LogKeywordNotifier(msg *discordgo.Message, keyword string)
	LogPhishingDetection(msg *discordgo.Message, domain string)
	LogMassDMBlocked()
	LogFriendRequestSpamBlocked()
	LogAntiDelete(msg *discordgo.Message)
	LogAFKMessage(msg *discordgo.Message, response string)
	LogMultiTokenSniper(numTokens, numGuilds int)
	LogSessionConnected(id, os, client string)
	LogSessionDisconnected(id, os, client string)
	LogSCJobs(numJobs int)
	LogEventReactions(numReactions int)
	LogDangerousCmdBlocked(user *discordgo.User, cmd string)
	LogDangerousHostingCmdBlocked(cmd string)
	LogReminder(reminder config.Reminder)
}
