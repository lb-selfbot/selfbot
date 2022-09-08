package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/utils"
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
			Token    string `json:"Token"`
			Password string `json:"Password"`
		} `json:"Discord"`
	} `json:"Account"`
	Notifications struct {
		Enabled bool   `json:"Enabled"`
		Type    string `json:"Type"`
	} `json:"Notifications"`
	Snipers struct {
		Blacklist      []string `json:"Blacklist"`
		GiveawaySniper struct {
			Sleep      int      `json:"Sleep"`
			Sniper     bool     `json:"Sniper"`
			CustomBots []string `json:"CustomBots"`
		} `json:"GiveawaySniper"`
		NitroSniper      bool     `json:"NitroSniper"`
		PrivnoteSniper   bool     `json:"PrivnoteSniper"`
		MultiTokenSniper bool     `json:"MultiTokenSniper"`
		Tokens           []string `json:"Tokens"`
		RedeemToken      string   `json:"RedeemToken"`
	} `json:"Snipers"`
	Webhooks struct {
		NitroSniper      string `json:"NitroSniper"`
		GiveawaySniper   string `json:"GiveawaySniper"`
		PrivnoteSniper   string `json:"PrivnoteSniper"`
		SelfbotDetection string `json:"SelfbotDetection"`
		StaffDetection   string `json:"StaffDetection"`
		Notifiers        string `json:"Notifiers"`
		PingDetection    string `json:"PingDetection"`
		Misc             string `json:"Misc"`
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
	} `json:"Events"`
	RPC RPCData `json:"RPC"`
	UI  struct {
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

func (c Config) Save() error {
	configMutex.Lock()
	defer configMutex.Unlock()
	configFile, err := os.Create(fmt.Sprintf("%s\\Config\\config.json", utils.File.GetMainFolder()))
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
	// Check if the config is already loaded in the last 2 seconds
	if lastConfig != nil && time.Since(lastConfigTime) < time.Second*10 {
		return *lastConfig, nil
	}
	configFile, err := os.Open(fmt.Sprintf("%s\\Config\\config.json", utils.File.GetMainFolder()))
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
