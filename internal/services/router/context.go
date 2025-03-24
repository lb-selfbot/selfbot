package router

import (
	"errors"
	"fmt"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/embed"
	"github.com/lb-selfbot/selfbot/internal/fn"
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

	switch cfg.Selfbot.CommandMode {
	case "WebhookEmbed":
		_ = ctx.sendWebhookEmbed(emb)
		return nil
	default:
		maxLen := 2000
		if ctx.Session.State.User.PremiumType == discordgo.UserPremiumTypeNitro {
			maxLen = 4000
		}

		msgs := emb.ToString(embed.EmbedFormat(cfg.Selfbot.CommandMode), maxLen)

		discordMsgs := fn.Map(msgs, func(text string) *discordgo.Message {
			message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, text)
			if err != nil {
				return nil
			}

			return message
		})

		discordMsgs = fn.Filter(discordMsgs, func(msg *discordgo.Message) bool {
			return msg != nil
		})

		if cfg.Selfbot.DeleteResponses {
			go ctx.deleteMessagesAfter(ctx.Session, discordMsgs, cfg.Selfbot.DeleteResponsesDelay)
		}

		return discordMsgs
	}
}

func (ctx *Ctx) sendWebhookEmbed(emb *embed.Embed) error {
	if ctx.Message.GuildID == "" {
		return errors.New("cannot send webhook embed in DMs")
	}

	webhooks, err := ctx.Session.ChannelWebhooks(ctx.Message.ChannelID)
	if err != nil {
		return err
	}

	webhook, ok := fn.First(webhooks, func(wh *discordgo.Webhook) bool {
		return wh.Name == "LightningBot"
	})
	if !ok {
		webhook, err = ctx.Session.WebhookCreate(ctx.Message.ChannelID, "LightningBot", "https://lightning-bot.com/logo.png")
		if err != nil {
			return err
		}
	}

	_, err = ctx.Session.WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Username:  "LightningBot",
		AvatarURL: "https://lightning-bot.com/logo.png",
		Embeds:    []*discordgo.MessageEmbed{emb.ToEmbed()},
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
