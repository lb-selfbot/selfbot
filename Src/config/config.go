package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/assets"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type RPCButton struct {
	Label string `json:"Label"`
	URL   string `json:"URL"`
}

type RPCData struct {
	Enabled    bool        `json:"Enabled"`
	ClientID   string      `json:"ClientID"`
	State      string      `json:"State"`
	Details    string      `json:"Details"`
	Time       bool        `json:"Time"`
	LargeImage string      `json:"LargeImage"`
	LargeText  string      `json:"LargeText"`
	SmallImage string      `json:"SmallImage"`
	SmallText  string      `json:"SmallText"`
	Invitable  bool        `json:"Invitable"`
	Buttons    []RPCButton `json:"Buttons"`
}

type Config struct {
	Account struct {
		LightningBot struct {
			Username string `json:"Username"`
			Password string `json:"Password"`
		} `json:"LightningBot"`
		Discord struct {
			Token string `json:"Token"`
		} `json:"Discord"`
	} `json:"Account"`
	Notifications struct {
		Enabled bool `json:"Enabled"`
	} `json:"Notifications"`
	Snipers struct {
		Blacklist       []string `json:"Blacklist"`
		ServerBlacklist []string `json:"ServerBlacklist"`
		GiveawaySniper  struct {
			Sleep      int      `json:"Sleep"`
			Enabled    bool     `json:"Enabled"`
			CustomBots []string `json:"CustomBots"`
		} `json:"GiveawaySniper"`
		NitroSniper      bool     `json:"NitroSniper"`
		MultiTokenSniper bool     `json:"MultiTokenSniper"`
		Tokens           []string `json:"Tokens"`
	} `json:"Snipers"`
	Webhooks struct {
		NitroSniper      string `json:"NitroSniper"`
		GiveawaySniper   string `json:"GiveawaySniper"`
		SelfbotDetection string `json:"SelfbotDetection"`
		StaffDetection   string `json:"StaffDetection"`
		Notifiers        string `json:"Notifiers"`
		PingDetection    string `json:"PingDetection"`
		Misc             string `json:"Misc"`
		Sessions         string `json:"Sessions"`
	} `json:"Webhooks"`
	Events struct {
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
	} `json:"Events"`
	RPC      RPCData `json:"RPC"`
	Activity struct {
		Enabled        bool   `json:"Enabled"`
		ActivityType   string `json:"ActivityType"`
		Type           string `json:"Type"`
		ClientID       string `json:"ClientID"`
		Name           string `json:"Name"`
		State          string `json:"State"`
		Details        string `json:"Details"`
		Invitable      bool   `json:"Invitable"`
		Time           bool   `json:"Time"`
		LargeImage     string `json:"LargeImage"`
		LargeImageText string `json:"LargeImageText"`
		SmallImage     string `json:"SmallImage"`
		SmallImageText string `json:"SmallImageText"`
		Button1Label   string `json:"Button1Label"`
		Button1URL     string `json:"Button1URL"`
		Button2Label   string `json:"Button2Label"`
		Button2URL     string `json:"Button2URL"`
	} `json:"Activity"`
	UI struct {
		Theme            string `json:"Theme"`
		HideToSystemTray bool   `json:"HideToSystemTray"`
	} `json:"UI"`
	Selfbot struct {
		CommandPrefix        string `json:"CommandPrefix"`
		CommandMode          string `json:"CommandMode"`
		Title                string `json:"Title"`
		Footer               string `json:"Footer"`
		DeleteResponses      bool   `json:"DeleteResponses"`
		DeleteResponsesDelay int    `json:"DeleteResponsesDelay"`
		DeleteCommands       bool   `json:"DeleteCommands"`
		MobileFriendly       bool   `json:"MobileFriendly"`
		Client               string `json:"Client"`
	}
	Protections struct {
		AntiPhising       bool `json:"AntiPhising"`
		AntiFriendRequest bool `json:"AntiFriendRequest"`
		AntiMassDM        bool `json:"AntiMassDM"`
		ServerProtections struct {
			AntiInvite       bool     `json:"AntiInvite"`
			AntiSpam         bool     `json:"AntiSpam"`
			AntiUpper        bool     `json:"AntiUpper"`
			AntiDelete       bool     `json:"AntiDelete"`
			ProtectedServers []string `json:"ProtectedServers"`
		}
	} `json:"Protections"`
	AFK struct {
		Enabled        bool     `json:"Enabled"`
		OnPing         bool     `json:"OnPing"`
		OnDM           bool     `json:"OnDM"`
		Message        string   `json:"Message"`
		Reply          bool     `json:"Reply"`
		Typing         bool     `json:"Typing"`
		Delay          int      `json:"Delay"`
		IgnoredServers []string `json:"IgnoredServers"`
		IgnoredUsers   []string `json:"IgnoredUsers"`
	} `json:"AFK"`
	Sharing struct {
		SharedUsers      []string `json:"SharedUsers"`
		SharedCategories []string `json:"SharedCategories"`
		SharedCommands   []string `json:"SharedCommands"`
		ShareAll         bool     `json:"ShareAll"`
	}
}

var (
	lastConfig     *Config
	lastConfigTime time.Time
	configMutex    = &sync.Mutex{}
)

func (c Config) GetClient() discordgo.Identify {
	switch c.Selfbot.Client {
	case "Mobile":
		return discordgo.IdentifyMobile
	case "Web":
		return discordgo.IdentifyWeb
	case "Desktop":
		return discordgo.IdentifyDiscordClient
	}

	return discordgo.IdentifyDiscordClient
}

func (c Config) Save() error {
	configMutex.Lock()
	defer configMutex.Unlock()

	configFile, err := os.OpenFile(filepath.Join(utils.IO.GetMainFolder(), "Config", "config.json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer configFile.Close()

	configBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	_, err = configFile.Write(configBytes)
	if err != nil {
		return err
	}

	lastConfig = &c
	lastConfigTime = time.Now()

	return nil
}

func GetConfig() (config Config, err error) {
	configMutex.Lock()
	defer configMutex.Unlock()

	// Check if the config is already loaded in the last 10 seconds
	if lastConfig != nil && time.Since(lastConfigTime) < time.Second*10 {
		return *lastConfig, nil
	}

	configFile, err := os.Open(filepath.Join(utils.IO.GetMainFolder(), "Config", "config.json"))
	if err != nil {
		return config, err
	}

	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, err
	}

	lastConfig = &config
	lastConfigTime = time.Now()

	return config, nil
}

func SaveConfigFromUI(config *Config) {
	configMutex.Lock()
	defer configMutex.Unlock()

	lastConfig = config
	lastConfigTime = time.Now()
}

func GetDefaultConfig() Config {
	var defaultConfig Config

	defaultConfigBytes, err := assets.Assets.ReadFile("defaultConfig.json")
	if err != nil {
		return defaultConfig
	}

	_ = json.Unmarshal(defaultConfigBytes, &defaultConfig)

	return defaultConfig
}
