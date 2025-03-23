package config

import (
	"encoding/json"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type TriggerList struct {
	Type    string   `json:"Type"`
	Users   []string `json:"Users"`
	Servers []string `json:"Servers"`
}

func (t TriggerList) IsTriggered(session *discordgo.Session, userID string, guildID string) bool {
	if t.Type == "Global" {
		return true
	}

	isUserInList := utils.Discord.IsUserInList(session, userID, t.Users)
	isGuildInList := utils.Discord.IsGuildInList(session, guildID, t.Servers)

	switch t.Type {
	case "Whitelist":
		return isUserInList || isGuildInList
	case "Blacklist":
		return !isUserInList && !isGuildInList
	default:
		return false
	}
}

type Activity struct {
	Enabled        bool   `json:"Enabled"`
	PauseForOther  bool   `json:"PauseForOther"`
	ActivityType   string `json:"ActivityType"`
	Type           string `json:"Type"`
	Platform       string `json:"Platform"`
	ClientID       string `json:"ClientID"`
	Name           string `json:"Name"`
	State          string `json:"State"`
	Details        string `json:"Details"`
	Invitable      bool   `json:"Invitable"`
	Time           bool   `json:"Time"`
	StreamingURL   string `json:"StreamingURL"`
	LargeImage     string `json:"LargeImage"`
	LargeImageText string `json:"LargeImageText"`
	SmallImage     string `json:"SmallImage"`
	SmallImageText string `json:"SmallImageText"`
	Button1Label   string `json:"Button1Label"`
	Button1URL     string `json:"Button1URL"`
	Button2Label   string `json:"Button2Label"`
	Button2URL     string `json:"Button2URL"`
}

type Reminder struct {
	Name    string    `json:"Name"`
	Content string    `json:"Content"`
	At      time.Time `json:"At"`

	// Deprecated
	ID *int `json:"ID,omitempty"`
}

type LightningBotConfig struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type DiscordConfig struct {
	Token string `json:"Token"`
}

type AccountConfig struct {
	LightningBot LightningBotConfig `json:"LightningBot"`
	Discord      DiscordConfig      `json:"Discord"`
}

type NotificationsConfig struct {
	Enabled bool `json:"Enabled"`
}

type GiveawaySniperConfig struct {
	Sleep      int      `json:"Sleep"`
	Enabled    bool     `json:"Enabled"`
	CustomBots []string `json:"CustomBots"`
}

type PokeTwoSniperConfig struct {
	Enabled       bool     `json:"Enabled"`
	Channels      []string `json:"Channels"`
	Sleep         int      `json:"Sleep"`
	RandomSleep   bool     `json:"RandomSleep"`
	Timeout       int      `json:"Timeout"`
	Typing        bool     `json:"Typing"`
	CatchRate     int      `json:"CatchRate"`
	FakeActivity  bool     `json:"FakeActivity"`
	CatchLanguage string   `json:"CatchLanguage"`
}

type SnipersConfig struct {
	Blacklist        []string             `json:"Blacklist"`
	ServerBlacklist  []string             `json:"ServerBlacklist"`
	GiveawaySniper   GiveawaySniperConfig `json:"GiveawaySniper"`
	NitroSniper      bool                 `json:"NitroSniper"`
	PrivnoteSniper   bool                 `json:"PrivnoteSniper"`
	PokeTwoSniper    PokeTwoSniperConfig  `json:"PokeTwoSniper"`
	MultiTokenSniper bool                 `json:"MultiTokenSniper"`
	Tokens           []string             `json:"Tokens"`
}

type WebhooksConfig struct {
	NitroSniper      string `json:"NitroSniper"`
	GiveawaySniper   string `json:"GiveawaySniper"`
	SelfbotDetection string `json:"SelfbotDetection"`
	StaffDetection   string `json:"StaffDetection"`
	Notifiers        string `json:"Notifiers"`
	PingDetection    string `json:"PingDetection"`
	Misc             string `json:"Misc"`
	Sessions         string `json:"Sessions"`
	PingSelfMode     string `json:"PingSelfMode"`
	PingEveryoneMode string `json:"PingEveryoneMode"`
}

type EventsConfig struct {
	BanNotifier           bool `json:"BanNotifier"`
	DMNotifier            bool `json:"DMNotifier"`
	DMEditNotifier        bool `json:"DMEditNotifier"`
	DMDeleteNotifier      bool `json:"DMDeleteNotifier"`
	DMTypingNotifier      bool `json:"DMTypingNotifier"`
	MessageDeleteNotifier bool `json:"MessageDeleteNotifier"`
	DiscordStaffDetection bool `json:"DiscordStaffDetection"`
	GhostPingDetection    bool `json:"GhostPingDetection"`
	SelfbotDetection      bool `json:"SelfbotDetection"`
	TicketNotifier        bool `json:"TicketNotifier"`
	SessionConnected      bool `json:"SessionConnected"`
	SessionDisconnected   bool `json:"SessionDisconnected"`
}

type UIConfig struct {
	Theme               string `json:"Theme"`
	HideToSystemTray    bool   `json:"HideToSystemTray"`
	RunOnStartup        bool   `json:"RunOnStartup"`
	RunSelfbotOnStartup bool   `json:"RunSelfbotOnStartup"`
}

type SelfbotConfig struct {
	CommandPrefix          string            `json:"CommandPrefix"`
	CommandMode            string            `json:"CommandMode"`
	Title                  string            `json:"Title"`
	Footer                 string            `json:"Footer"`
	DeleteResponses        bool              `json:"DeleteResponses"`
	DeleteResponsesDelay   int               `json:"DeleteResponsesDelay"`
	DeleteCommands         bool              `json:"DeleteCommands"`
	Client                 string            `json:"Client"`
	Timezone               string            `json:"Timezone"`
	EmbedColor             string            `json:"EmbedColor"`
	EmbedImage             string            `json:"EmbedImage"`
	EmbedURL               string            `json:"EmbedURL"`
	CommandAliases         map[string]string `json:"CommandAliases"`
	ForceSubscribedServers []string          `json:"ForceSubscribedServers"`
}

type ProtectionsConfig struct {
	AntiPhishing      bool `json:"AntiPhishing"`
	AntiFriendRequest bool `json:"AntiFriendRequest"`
	AntiMassDM        bool `json:"AntiMassDM"`
	ServerProtections struct {
		AntiInvite       bool     `json:"AntiInvite"`
		AntiSpam         bool     `json:"AntiSpam"`
		AntiUpper        bool     `json:"AntiUpper"`
		AntiDelete       bool     `json:"AntiDelete"`
		ProtectedServers []string `json:"ProtectedServers"`
	}
}

type AFKConfig struct {
	Enabled         bool        `json:"Enabled"`
	Delay           int         `json:"Delay"`
	RandomDelay     bool        `json:"RandomDelay"`
	OnDM            bool        `json:"OnDM"`
	OnPing          bool        `json:"OnPing"`
	MustBeInactive  bool        `json:"MustBeInactive"`
	InactiveTimeout int         `json:"InactiveTimeout"`
	Reply           bool        `json:"Reply"`
	Typing          bool        `json:"Typing"`
	TriggerList     TriggerList `json:"TriggerList"`
	Message         string      `json:"Message"`
	AI              bool        `json:"AI"`
}

type SharingConfig struct {
	SharedUsers            []string `json:"SharedUsers"`
	SharedCategories       []string `json:"SharedCategories"`
	SharedCommands         []string `json:"SharedCommands"`
	ShareAll               bool     `json:"ShareAll"`
	BlockDangerousCommands bool     `json:"BlockDangerousCommands"`
}

type Config struct {
	Account          AccountConfig       `json:"Account"`
	Notifications    NotificationsConfig `json:"Notifications"`
	Snipers          SnipersConfig       `json:"Snipers"`
	Webhooks         WebhooksConfig      `json:"Webhooks"`
	Events           EventsConfig        `json:"Events"`
	Activity         Activity            `json:"Activity"`
	ActivityPresets  map[string]Activity `json:"ActivityPresets"`
	UI               UIConfig            `json:"UI"`
	Selfbot          SelfbotConfig       `json:"Selfbot"`
	Protections      ProtectionsConfig   `json:"Protections"`
	AFK              AFKConfig           `json:"AFK"`
	Sharing          SharingConfig       `json:"Sharing"`
	SlashCommandJobs json.RawMessage     `json:"SlashCommandJobs"`
	EventReactions   json.RawMessage     `json:"EventReactions"`
	Tags             map[string]string   `json:"Tags"`
	Reminders        []Reminder          `json:"Reminders"`
}

func (c *Config) GetClient() discordgo.Identify {
	switch c.Selfbot.Client {
	case "Mobile":
		return discordgo.IdentifyMobile
	case "Web":
		return discordgo.IdentifyWeb
	case "Desktop":
		return discordgo.IdentifyDiscordClient
	case "Console":
		return discordgo.IdentifyEmbedded
	default:
		return discordgo.IdentifyDiscordClient
	}
}

func (c *Config) GetTimezone() *time.Location {
	location, err := time.LoadLocation(c.Selfbot.Timezone)
	if err != nil {
		location = time.UTC
	}

	// Set the local timezone
	time.Local = location

	return location
}
