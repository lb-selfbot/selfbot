package featurelogs

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/config"
	"github.com/LightningDev1/discordgo"
)

// CmdUsed has the data for a command used log.
type CmdUsed struct {
	Name string
}

// Ready has the data for a ready event log.
type Ready struct {
	User    *discordgo.User
	Guilds  []*discordgo.Guild
	NumCmds int
}

// NitroSniper has the data for a nitro sniper log.
type NitroSniper struct {
	Message *discordgo.Message
	Elapsed time.Duration
	Status  string
	Code    string
	Account string
}

// NitroSniperFake has the data for a fake nitro sniper log.
type NitroSniperFake struct {
	Message *discordgo.Message
	Code    string
	Account string
}

// PrivnoteSniper has the data for a privnote sniper log.
type PrivnoteSniper struct {
	Message *discordgo.Message
	Elapsed time.Duration
	Status  string
	Path    string
	Account string
}

// GiveawaySniperFound has the data for a found giveaway log.
type GiveawaySniperFound struct {
	Message *discordgo.Message
	Prize   string
	Host    string
	EndsIn  string
	Account string
}

// GiveawaySniperWon has the data for a won giveaway log.
type GiveawaySniperWon struct {
	Message *discordgo.Message
	Prize   string
	Account string
}

// PokeTwoSniper has the data for a PokéTwo sniper log.
type PokeTwoSniper struct {
	Message    *discordgo.Message
	Elapsed    time.Duration
	Confidence float64
	Pokemon    string
	Received   string
	IV         string
	Account    string
}

// BanNotifier has the data for a ban notification.
type BanNotifier struct {
	GuildID string
}

// DMNotifier has the data for a DM notification.
type DMNotifier struct {
	Message *discordgo.Message
}

// DMEditNotifier has the data for an edited DM notification.
type DMEditNotifier struct {
	Message *discordgo.Message
	Before  string
	After   string
}

// DMDeleteNotifier has the data for a deleted DM notification.
type DMDeleteNotifier struct {
	Message     *discordgo.Message
	Attachments []string
}

// DMTypingNotifier has the data for a DM typing notification.
type DMTypingNotifier struct {
	Username string
}

// MessageDeleteNotifier has the data for a message delete notification.
type MessageDeleteNotifier struct {
	Message     *discordgo.Message
	Attachments []string
}

// DiscordStaffDetection has the data for a Discord staff detection.
type DiscordStaffDetection struct {
	Message *discordgo.Message
}

// GhostPingDetection has the data for a ghost ping detection.
type GhostPingDetection struct {
	Message *discordgo.Message
}

// SelfbotDetection has the data for a selfbot detection.
type SelfbotDetection struct {
	Message *discordgo.Message
	Selfbot string
}

// TicketNotifier has the data for a ticket notification.
type TicketNotifier struct {
	Channel *discordgo.Channel
}

// KeywordNotifier has the data for a keyword notification.
type KeywordNotifier struct {
	Message *discordgo.Message
	Keyword string
}

// PhishingDetection has the data for a phishing detection.
type PhishingDetection struct {
	Message *discordgo.Message
	Domain  string
}

// MassDMBlocked represents a mass DM block event.
type MassDMBlocked struct{}

// FriendRequestSpamBlocked represents a friend request spam block event.
type FriendRequestSpamBlocked struct{}

// AntiDelete has the data for an anti-delete event.
type AntiDelete struct {
	Message *discordgo.Message
}

// AFKMessage has the data for an AFK message response.
type AFKMessage struct {
	Message  *discordgo.Message
	Response string
}

// MultiTokenSniper has the data for a multi-token sniper log.
type MultiTokenSniper struct {
	NumTokens int
	NumGuilds int
}

// SessionConnected has the data for a session connection log.
type SessionConnected struct {
	ID     string
	OS     string
	Client string
}

// SessionDisconnected has the data for a session disconnection log.
type SessionDisconnected struct {
	ID     string
	OS     string
	Client string
}

// SCJobs has the data for slash command jobs log.
type SCJobs struct {
	NumJobs int
}

// EventReactions has the data for event reactions log.
type EventReactions struct {
	NumReactions int
}

// DangerousCmdBlocked has the data for a blocked dangerous command.
type DangerousCmdBlocked struct {
	User *discordgo.User
	Cmd  string
}

// DangerousHostingCmdBlocked has the data for a blocked dangerous hosting command.
type DangerousHostingCmdBlocked struct {
	Cmd string
}

// Reminder has the data for a reminder log.
type Reminder struct {
	Reminder config.Reminder
}
