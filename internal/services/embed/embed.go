package embed

import (
	"errors"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/embed"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/services/variables"
)

// Service handles the creation and sending of embeds.
type Service interface {
	// NewEmbed creates a new embed with the given title and description.
	New(ctx *router.Ctx) *embed.Embed

	// RespondEmbed creates a new embed with the given description and sends
	// it in the given context.
	Respond(ctx *router.Ctx, format string, a ...any) []*discordgo.Message

	// RespondEmbed creates a new embed with the given field and sends
	// it in the given context.
	RespondField(ctx *router.Ctx, name, value string, cond ...bool) []*discordgo.Message
}

type embedService struct {
	configSvc    config.Service
	variablesSvc variables.Service
}

// NewService creates a new instance of the embed service.
func NewService(configSvc config.Service, variablesSvc variables.Service) Service {
	return &embedService{configSvc: configSvc, variablesSvc: variablesSvc}
}

// New creates a new embed with the given title and description.
func (e *embedService) New(ctx *router.Ctx) *embed.Embed {
	cfg := e.configSvc.GetConfig()

	return embed.NewEmbed(e.getEmbedSendFunc(ctx), e.configSvc.GetConfig).
		SetTitle("%s", e.variablesSvc.Parse(cfg.Selfbot.Title)).
		SetSubtext("%s", e.variablesSvc.Parse(cfg.Selfbot.Footer))
}

func (e *embedService) Respond(ctx *router.Ctx, format string, a ...any) []*discordgo.Message {
	return e.New(ctx).SetDescription(format, a...).Send()
}

func (e *embedService) RespondField(ctx *router.Ctx, name, value string, cond ...bool) []*discordgo.Message {
	return e.New(ctx).AddField(name, value, cond...).Send()
}

func (e *embedService) getEmbedSendFunc(ctx *router.Ctx) func(*embed.Embed) []*discordgo.Message {
	return func(emb *embed.Embed) []*discordgo.Message {
		cfg := e.configSvc.GetConfig()

		switch cfg.Selfbot.CommandMode {
		case "WebhookEmbed":
			_ = e.sendWebhookEmbed(ctx, emb)
			return nil
		default:
			maxLen := 2000
			if ctx.Session.State.User.PremiumType == discordgo.UserPremiumTypeNitro {
				maxLen = 4000
			}

			msgs := emb.ToString(cfg.Selfbot.CommandMode, maxLen)

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
				go e.deleteMessagesAfter(ctx.Session, discordMsgs, cfg.Selfbot.DeleteResponsesDelay)
			}

			return discordMsgs
		}
	}
}

func (e *embedService) sendWebhookEmbed(ctx *router.Ctx, emb *embed.Embed) error {
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

func (e *embedService) deleteMessagesAfter(session *discordgo.Session, messages []*discordgo.Message, delay int) {
	time.Sleep(time.Duration(delay) * time.Second)

	for _, message := range messages {
		if err := session.ChannelMessageDelete(message.ChannelID, message.ID); err != nil {
			return
		}
	}
}
