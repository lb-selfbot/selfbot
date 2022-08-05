package eventhandlers

import (
	"fmt"
	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func (e *Events) OnMessage(_ *discordgo.Session, m *discordgo.MessageCreate) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	e.checkNitroSniper(m.Message, &cfg)
	e.checkPrivnoteSniper(m.Message, &cfg)
	e.checkGiveawaySniper(m.Message, &cfg)
	e.checkIgnore(m.Message)
	e.checkPhising(m.Message, &cfg)
}

func (e *Events) OnMessageEdit(_ *discordgo.Session, m *discordgo.MessageUpdate) {
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	e.checkNitroSniper(m.Message, &cfg)
	e.checkPrivnoteSniper(m.Message, &cfg)
	e.checkGiveawaySniper(m.Message, &cfg)
	e.checkIgnore(m.Message)
	e.checkPhising(m.Message, &cfg)
}

var nitroRegex = regexp.MustCompile(`(discord\.gift|discord\.com/gifts|discordapp\.com/gifts)/(\w{0,25})`)

func (e *Events) checkNitroSniper(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Snipers.NitroSniper {
		return
	}
	nitroLinks := nitroRegex.FindAllString(message.Content, -1)
	for _, nitroLink := range nitroLinks {
		nitroLinkSplit := strings.Split(nitroLink, "/")
		nitroCode := nitroLinkSplit[len(nitroLinkSplit)-1]
		if utils.SliceContains(e.Settings.alreadyCheckedNitroCodes, nitroCode) {
			return
		}
		e.Settings.alreadyCheckedNitroCodes = append(e.Settings.alreadyCheckedNitroCodes, nitroCode)
		if len(nitroCode) < 15 || len(nitroCode) > 26 {
			// TODO: Fake Code
			return
		}
		// TODO: check if author is spamming codes
		start := time.Now()
		result, err := http.Post("https://discord.com/api/v10/entitlements/gift-codes/{nitroCode}/redeem", "", http.GetDiscordHeaders(cfg.Account.Discord.Token))
		if err != nil {
			return
		}
		elapsed := time.Since(start)
		fmt.Println("Elapsed:", elapsed)
		if strings.Contains(result, "This gift has been redeemed already.") {
			// TODO: Already Redeemed
		} else if strings.Contains(result, "Already purchased") {
			// TODO: Already Purchased
		} else if strings.Contains(result, "subscription_plan") {
			// TODO: Redeemed
		} else if strings.Contains(result, "Unknown Gift Code") {
			// TODO: Unknown Code
		} else {
			// TODO: Unknown Status (json message)
		}
	}
}

var privnoteRegex = regexp.MustCompile(`privnote.com/([a-zA-Z\d_#]{18})`)

func (e *Events) checkPrivnoteSniper(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Snipers.PrivnoteSniper {
		return
	}
	privnoteLinks := privnoteRegex.FindAllString(message.Content, -1)
	for _, privnoteLink := range privnoteLinks {
		privnoteLink = "https://" + privnoteLink
		// TODO: try to open
	}
}

func (e *Events) checkGiveawaySniper(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Snipers.GiveawaySniper.Sniper {
		return
	}
	// TODO: giveaway sniper
}

func (e *Events) checkIgnore(message *discordgo.Message) {
	for _, user := range e.Settings.IgnoredUsers {
		if message.Author.ID == user.ID {
			_ = e.Session.MessageAck(message.ChannelID, message.ID)
		}
	}
	for _, guild := range e.Settings.IgnoredGuilds {
		if message.GuildID == guild.ID {
			_ = e.Session.GuildAck(message.GuildID)
		}
	}
}

func (e *Events) checkPhising(message *discordgo.Message, cfg *config.Config) {
	if !cfg.Protections.AntiPhising {
		return
	}

	for _, part := range strings.Split(message.Content, " ") {
		link, err := url.Parse(part)
		if err == nil {
			hostname := link.Hostname()
			if utils.SliceContains(e.Settings.phisingDomains, hostname) ||
				utils.SliceContains(e.Settings.suspiciousDomains, hostname) {
				// TODO: log phishing
				fmt.Println("Phishing domain detected:", hostname, "in", message.Content)
				_ = e.Session.MessageAck(message.ChannelID, message.ID)
			}
		}
	}
}
