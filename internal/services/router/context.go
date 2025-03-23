package router

import (
	"errors"
	"fmt"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/embed"
	"github.com/lb-selfbot/selfbot/internal/services/config"
)

// Ctx represents the context for command execution.
type Ctx struct {
	Router    Service
	configSvc config.Service

	Command   *Command
	Arguments *Arguments

	Session *discordgo.Session
	Message *discordgo.Message
}

func (ctx *Ctx) NewEmbed() *embed.Embed {
	return embed.NewEmbed(ctx.sendEmbed)
}

func (ctx *Ctx) sendEmbed(emb *embed.Embed) []*discordgo.Message {
	cfg := ctx.configSvc.GetConfig()

	commandMode := cfg.Selfbot.CommandMode
	if commandMode == "WebhookEmbed" {
		if err := ctx.sendWebhookEmbed(emb); err == nil {
			return nil
		}
	}

	var messages []*discordgo.Message

	for _, text := range emb.ToString(commandMode) {
		message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, text)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if emb.Image != nil && emb.Image.URL != "" {
		message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, emb.Image.URL)
		if err == nil {
			messages = append(messages, message)
		}
	}

	if cfg.Selfbot.DeleteResponses {
		go ctx.deleteMessagesAfter(ctx.Session, messages, cfg.Selfbot.DeleteResponsesDelay)
	}

	return messages
}

func (ctx *Ctx) sendWebhookEmbed(emb *embed.Embed) error {
	if ctx.Message.GuildID == "" {
		return errors.New("cannot send webhook embed in DMs")
	}

	webhooks, err := ctx.Session.ChannelWebhooks(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	var webhook *discordgo.Webhook

	for _, wh := range webhooks {
		if wh.Name != "LightningBot" {
			continue
		}

		webhook = wh
		break
	}

	if webhook == nil {
		webhook, err = ctx.Session.WebhookCreate(ctx.Message.ChannelID, "LightningBot", "https://lightning-bot.com/logo.png")
		if err != nil {
			return err
		}
	}

	_, err = ctx.Session.WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Username:  "LightningBot",
		AvatarURL: "https://lightning-bot.com/logo.png",
		Embeds: []*discordgo.MessageEmbed{
			emb.ToEmbed(),
		},
	})

	return err
}

func (ctx *Ctx) deleteMessagesAfter(session *discordgo.Session, messages []*discordgo.Message, delay int) {
	time.Sleep(time.Duration(delay) * time.Second)

	for _, message := range messages {
		if err := session.ChannelMessageDelete(message.ChannelID, message.ID); err != nil {
			return
		}
	}
}

func (ctx *Ctx) RespondText(format string, e ...any) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(format, e...))
	return err
}
