package eventhandlers

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/selfbot/eventhandlers/eventreactions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/snipers"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/websocket"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnMessage(_ *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Message == nil || m.Author == nil {
		utils.Logging.Debug("Message: Author or Message is nil: ", m.Author, m.Message)
		return
	}

	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply && m.Type != discordgo.MessageTypeChatInputCommand {
		return
	}

	eventreactions.TimeDifference = time.Since(utils.Discord.SnowflakeToTime(m.ID))

	e.Session.State.MessageAdd(m.Message)

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	e.Settings.AlreadyCheckedNitroCodes, e.Settings.SentNitroCodes = snipers.CheckNitroSniper(
		e.Session,
		m.Message,
		cfg,
		e.Settings.CountryCode,
		e.Settings.PaymentSourceID,
		e.Settings.AlreadyCheckedNitroCodes,
		e.Settings.SentNitroCodes,
		"",
	)

	snipers.CheckPrivnoteSniper(e.Session, m.Message, &cfg, "")

	snipers.CheckGiveawaySniper(e.Session, m.Message, &cfg, "")

	websocket.Websocket.CheckNonce()
}

func (e *Events) OnMessageEdit(_ *discordgo.Session, m *discordgo.MessageUpdate) {
	message, err := e.Session.State.Message(m.Message.ChannelID, m.Message.ID)
	if err == nil {
		m.Message = message
	}

	if m.Message == nil || m.Author == nil {
		utils.Logging.Debug("MessageEdit: Message or Author is nil: ", m.Message)
		return
	}

	if m.BeforeUpdate == nil {
		m.BeforeUpdate = message
	}

	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply && m.Type != discordgo.MessageTypeChatInputCommand {
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	e.Settings.AlreadyCheckedNitroCodes, e.Settings.SentNitroCodes = snipers.CheckNitroSniper(
		e.Session,
		m.Message,
		cfg,
		e.Settings.CountryCode,
		e.Settings.PaymentSourceID,
		e.Settings.AlreadyCheckedNitroCodes,
		e.Settings.SentNitroCodes,
		"",
	)

	snipers.CheckPrivnoteSniper(e.Session, m.Message, &cfg, "")

	snipers.CheckGiveawaySniperEdit(e.Session, m.Message, &cfg, "")
}

func (e *Events) OnMessageDelete(_ *discordgo.Session, m *discordgo.MessageDelete) {
	if m.BeforeDelete != nil {
		m.Message = m.BeforeDelete
	}

	if m.Message == nil {
		utils.Logging.Debug("MessageDelete: Message is nil: ", m.Message)
		return
	}

	if m.Author == nil {
		utils.Logging.Debug("MessageDelete: Author is nil: ", m.Author)
		return
	}

	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
		return
	}
}
