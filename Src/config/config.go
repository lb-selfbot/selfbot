package config

import (
	"encoding/json"
	"fmt"
	"os"

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
			Sleep  int  `json:"Sleep"`
			Sniper bool `json:"Sniper"`
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
		DiscordStaffDetection    bool `json:"DiscordStaffDetection"`
		GhostPingDetection       bool `json:"GhostPingDetection"`
		TicketDetection          bool `json:"TicketDetection"`
		SelfbotDetection         bool `json:"SelfbotDetection"`
		DMTypingDetection        bool `json:"DMTypingDetection"`
		DMMessageDetection       bool `json:"DMMessageDetection"`
		DMMessageDeleteDetection bool `json:"DMMessageDeleteDetection"`
		DMMessageEditDetection   bool `json:"DMMessageEditDetection"`
		MessageDeleteDetection   bool `json:"MessageDeleteDetection"`
		BanNotifier              bool `json:"BanNotifier"`
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
	}
}

func (c Config) Save() error {
	configFile, err := os.OpenFile(fmt.Sprintf("%s\\Config\\config.json", utils.GetMainFolder()), os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer configFile.Close()
	return json.NewEncoder(configFile).Encode(c)
}

func GetConfig() (Config, error) {
	var config Config
	configFile, err := os.Open(fmt.Sprintf("%s\\Config\\config.json", utils.GetMainFolder()))
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return config, err
	}
	return config, nil
}
