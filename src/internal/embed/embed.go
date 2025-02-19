package embed

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/config"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

// https://github.com/Clinet/discordgo-embed

// Embed ...
type Embed struct {
	*discordgo.MessageEmbed
	FooterStr string
	Subtext   string
}

var (
	LimitTitle       = 256
	LimitDescription = 2048
	LimitFieldValue  = 1024
	LimitFieldName   = 256
	LimitField       = 25
	LimitFooter      = 2048
)

// SetTitle ...
func (e *Embed) SetTitle(name string) *Embed {
	e.Title = name
	return e
}

// SetDescription [desc]
func (e *Embed) SetDescription(description string) *Embed {
	if len(description) > 2048 {
		description = description[:2048]
	}
	e.Description = description
	return e
}

// SetSubtext ...
func (e *Embed) SetSubtext(text string) *Embed {
	e.Subtext = text
	return e
}

// AddField [name] [value]
func (e *Embed) AddField(name, value string, flags ...bool) *Embed {
	inline := len(flags) > 0
	bypass := len(flags) > 1
	fields := make([]*discordgo.MessageEmbedField, 0)

	if len(name) > LimitFieldName {
		name = name[:LimitFieldName]
	}

	if len(value) > LimitFieldValue && !bypass {
		i := LimitFieldValue
		extended := false
		for i = LimitFieldValue; i < len(value); {
			if i != LimitFieldValue && !extended {
				name += " (extended)"
				extended = true
			}
			if value[i] == []byte(" ")[0] || value[i] == []byte("\n")[0] || value[i] == []byte("-")[0] {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:   name,
					Value:  value[i-LimitFieldValue : i],
					Inline: inline,
				})
			} else {
				fields = append(fields, &discordgo.MessageEmbedField{
					Name:   name,
					Value:  value[i-LimitFieldValue:i-1] + "-",
					Inline: inline,
				})
				i--
			}

			if (i + LimitFieldValue) < len(value) {
				i += LimitFieldValue
			} else {
				break
			}
		}
		if i < len(value) {
			name += " (extended)"
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   name,
				Value:  value[i:],
				Inline: inline,
			})
		}
	} else {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  value,
			Inline: inline,
		})
	}

	e.Fields = append(e.Fields, fields...)

	return e
}

// SetFooter [Text] [iconURL]
func (e *Embed) SetFooter(args ...string) *Embed {
	iconURL := ""
	text := ""
	proxyURL := ""

	switch {
	case len(args) > 2:
		proxyURL = args[2]
		fallthrough
	case len(args) > 1:
		iconURL = args[1]
		fallthrough
	case len(args) > 0:
		text = args[0]
	case len(args) == 0:
		return e
	}

	e.Footer = &discordgo.MessageEmbedFooter{
		IconURL:      iconURL,
		Text:         text,
		ProxyIconURL: proxyURL,
	}

	e.FooterStr = text

	return e
}

// SetImage ...
func (e *Embed) SetImage(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Image = &discordgo.MessageEmbedImage{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetThumbnail ...
func (e *Embed) SetThumbnail(args ...string) *Embed {
	var URL string
	var proxyURL string

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		URL = args[0]
	}
	if len(args) > 1 {
		proxyURL = args[1]
	}
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL:      URL,
		ProxyURL: proxyURL,
	}
	return e
}

// SetAuthor ...
func (e *Embed) SetAuthor(args ...string) *Embed {
	var (
		name     string
		iconURL  string
		URL      string
		proxyURL string
	)

	if len(args) == 0 {
		return e
	}
	if len(args) > 0 {
		name = args[0]
	}
	if len(args) > 1 {
		iconURL = args[1]
	}
	if len(args) > 2 {
		URL = args[2]
	}
	if len(args) > 3 {
		proxyURL = args[3]
	}

	e.Author = &discordgo.MessageEmbedAuthor{
		Name:         name,
		IconURL:      iconURL,
		URL:          URL,
		ProxyIconURL: proxyURL,
	}

	return e
}

// SetURL ...
func (e *Embed) SetURL(URL string) *Embed {
	e.URL = URL
	return e
}

// SetColor ...
func (e *Embed) SetColor(clr int) *Embed {
	e.Color = clr
	return e
}

// InlineAllFields sets all fields in the embed to be inline
func (e *Embed) InlineAllFields() *Embed {
	for _, v := range e.Fields {
		v.Inline = true
	}
	return e
}

// Truncate truncates any embed value over the character limit.
func (e *Embed) Truncate() *Embed {
	e.TruncateDescription()
	e.TruncateFields()
	e.TruncateFooter()
	e.TruncateTitle()
	return e
}

// MakeFieldInline adds last field as InLine
func (e *Embed) MakeFieldInline() *Embed {
	length := len(e.Fields) - 1
	e.Fields[length].Inline = true
	return e
}

// TruncateFields truncates fields that are too long
func (e *Embed) TruncateFields() *Embed {
	if len(e.Fields) > 25 {
		e.Fields = e.Fields[:LimitField]
	}

	for _, v := range e.Fields {

		if len(v.Name) > LimitFieldName {
			v.Name = v.Name[:LimitFieldName]
		}

		if len(v.Value) > LimitFieldValue {
			v.Value = v.Value[:LimitFieldValue]
		}

	}
	return e
}

// TruncateDescription ...
func (e *Embed) TruncateDescription() *Embed {
	if len(e.Description) > LimitDescription {
		e.Description = e.Description[:LimitDescription]
	}
	return e
}

// TruncateTitle ...
func (e *Embed) TruncateTitle() *Embed {
	if len(e.Title) > LimitTitle {
		e.Title = e.Title[:LimitTitle]
	}
	return e
}

// TruncateFooter ...
func (e *Embed) TruncateFooter() *Embed {
	if e.FooterStr != "" && len(e.FooterStr) > LimitFooter {
		e.FooterStr = e.FooterStr[:LimitFooter]
		e.Footer.Text = e.Footer.Text[:LimitFooter]
	}
	return e
}

func (e *Embed) ToJson() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func (e *Embed) ToCodeblock() []string {
	embedString := ""
	var embedList []string
	if e.Author != nil {
		embedString += fmt.Sprintf("```fix\n%s\n```\n", e.Author.Name)
	}
	if e.Title != "" {
		embedString += fmt.Sprintf("```ini\n[%s]\n```", e.Title)
	}
	if e.Description != "" {
		if len(embedString+e.Description) > 2000 {
			embedList = append(embedList, embedString)
			embedString = ""
		}
		embedString += fmt.Sprintf("```md\n%s\n```", e.Description)
	}
	if len(e.Fields) > 0 {
		mdAdded := false
		for _, v := range e.Fields {
			fieldText := fmt.Sprintf("# %s\n%s\n\n", v.Name, v.Value)
			if len(embedString+fieldText) > 2000 {
				if mdAdded {
					embedString += "```"
					mdAdded = false
				}
				embedList = append(embedList, embedString)
				embedString = ""
			}
			if !mdAdded {
				embedString += "```md\n"
				mdAdded = true
			}
			embedString += fieldText
		}
		embedString += "```"
	}
	if e.Subtext != "" {
		embedString += fmt.Sprintf("```ansi\n%s\n```", e.Subtext)
	}
	if e.FooterStr != "" {
		footerString := fmt.Sprintf("```css\n[%s]```", e.FooterStr)
		if len(embedString+footerString) > 2000 {
			embedList = append(embedList, embedString)
			embedString = ""
		}
		embedString += footerString
	}
	if embedString != "" {
		embedList = append(embedList, embedString)
	}
	return embedList
}

func (e *Embed) ToText() []string {
	str := ""
	if e.Author != nil {
		str += fmt.Sprintf("%s\n\n", e.Author.Name)
	}
	if e.Title != "" {
		str += fmt.Sprintf("```ini\n[%s]```\n\n", e.Title)
	}
	if e.Description != "" {
		str += fmt.Sprintf("%s\n\n", e.Description)
	}
	if len(e.Fields) > 0 {
		for _, v := range e.Fields {
			if v.Inline {
				v.Value = "```" + v.Value + "```"
			}
			str += fmt.Sprintf("> **%s**\n%s\n\n", v.Name, v.Value)
		}
	}
	if e.Subtext != "" {
		str += fmt.Sprintf("```ansi\n%s\n```", e.Subtext)
	}
	if e.FooterStr != "" {
		str += fmt.Sprintf("\n```css\n[%s]```\n", e.FooterStr)
	}
	var strings_ []string
	newStr := ""
	for _, line := range strings.Split(str, "\n") {
		if len(newStr+line)+1 > 2000 {
			strings_ = append(strings_, newStr)
			newStr = ""
		}
		newStr += line + "\n"
	}
	strings_ = append(strings_, newStr)
	return strings_
}

func (e *Embed) ToIndent() []string {
	str := ""
	if e.Author != nil {
		str += fmt.Sprintf("%s\n\n", e.Author.Name)
	}
	if e.Title != "" {
		str += fmt.Sprintf("```ini\n[%s]```\n\n", e.Title)
	}
	if e.Description != "" {
		str += fmt.Sprintf("%s\n\n", e.Description)
	}
	if len(e.Fields) > 0 {
		for _, v := range e.Fields {
			if v.Inline {
				v.Value = "```" + v.Value + "```"
			}
			str += fmt.Sprintf("**%s**\n%s\n\n", v.Name, v.Value)
		}
	}
	if e.Subtext != "" {
		str += fmt.Sprintf("```ansi\n%s\n```", e.Subtext)
	}
	if e.FooterStr != "" {
		str += fmt.Sprintf("\n```css\n[%s]```", e.FooterStr)
	}

	splitted := strings.Split(str, "\n")

	var strings_ []string
	var newStr string

	for _, line := range splitted {
		if len(newStr+line)+1 > 2000 {
			strings_ = append(strings_, newStr)
			newStr = ""
		}

		newStr += "> " + line + "\n"
	}

	return append(strings_, newStr)
}

func (e *Embed) ToURLEmbed() []string {
	// Generate description
	description := ""
	if e.Description != "" {
		description = e.Description
	}
	if len(e.Fields) > 0 {
		for _, v := range e.Fields {
			description += fmt.Sprintf("\n\n• %s\n%s", v.Name, v.Value)
		}
	}
	if e.Subtext != "" {
		e.Subtext = constants.ANSI_REGEX.ReplaceAllString(e.Subtext, "")
		description += fmt.Sprintf("\n\n%s", e.Subtext)
	}
	if e.FooterStr != "" {
		description += fmt.Sprintf("\n\n%s", e.FooterStr)
	}

	descriptions := []string{}

	// Split descriptions into multiple embeds
	var current string
	for _, line := range strings.Split(description, "\n") {
		if len(current+line)+1 > 350 {
			descriptions = append(descriptions, current)
			current = ""
		}

		current += line + "\n"
	}

	descriptions = append(descriptions, current)

	cfg, _ := config.GetConfig()

	// Generate embeds
	embeds := []api.URLEmbed{}

	embeds = append(embeds, api.URLEmbed{
		Title:       e.Title,
		Description: descriptions[0],
		URL:         utils.Misc.Or(cfg.Selfbot.EmbedURL, "https://lightning-bot.com/"),
		Color:       utils.Misc.Or(cfg.Selfbot.EmbedColor, "3990c1"),
		Image:       utils.Misc.Or(cfg.Selfbot.EmbedImage, "https://lightning-bot.com/logo.png"),
	})

	if e.Author != nil {
		embeds[0].AuthorName = e.Author.Name
		embeds[0].AuthorURL = e.Author.URL
	}

	for _, description := range descriptions[1:] {
		embeds = append(embeds, api.URLEmbed{
			Description: description,
			Color:       utils.Misc.Or(cfg.Selfbot.EmbedColor, "3990c1"),
		})
	}

	urls := utils.Map(embeds, func(embed api.URLEmbed) string {
		url, err := api.GenerateURLEmbed(embed)
		if err != nil {
			return ""
		}

		return url
	})

	// Generate messages
	var batches [][]string
	batchSize := 5

	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}
		batch := urls[i:end]
		batches = append(batches, batch)
	}

	return utils.Map(batches, func(batch []string) string {
		return constants.INVISIBLE_TEXT + " " + strings.Join(batch, " ")
	})
}

func (e *Embed) ToEmbed() *discordgo.MessageEmbed {
	if e.Subtext != "" {
		e.AddField("Subtext", e.Subtext)
	}

	e.Truncate()
	e.SetColor(0x3990c1)

	return e.MessageEmbed
}

func (e *Embed) ToString(commandMode string) []string {
	switch commandMode {
	case "Codeblock":
		return e.ToCodeblock()
	case "Text":
		return e.ToText()
	case "Indent":
		return e.ToIndent()
	case "URLEmbed":
		return e.ToURLEmbed()
	case "load":
		cfg, _ := config.GetConfig()

		return e.ToString(utils.Misc.Or(cfg.Selfbot.CommandMode, "Codeblock"))
	default:
		return e.ToCodeblock()
	}
}

func (e *Embed) Send(ctx *dgc.Ctx) []*discordgo.Message {
	cfg, _ := config.GetConfig()

	commandMode := utils.Misc.Or(cfg.Selfbot.CommandMode, "Codeblock")

	if commandMode == "WebhookEmbed" {
		err := e.sendWebhookEmbed(ctx)
		if err == nil {
			return nil
		}

		// If the webhook embed fails, fallback to the default embed
		commandMode = "Codeblock"
	}

	messages := []*discordgo.Message{}

	for _, text := range e.ToString(commandMode) {
		message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, text)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	if e.Image != nil && e.Image.URL != "" {
		message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, e.Image.URL)
		if err == nil {
			messages = append(messages, message)
		}
	}

	if cfg.Selfbot.DeleteResponses {
		go deleteMessagesAfter(ctx.Session, messages, cfg.Selfbot.DeleteResponsesDelay)
	}

	return messages
}

func (e *Embed) sendWebhookEmbed(ctx *dgc.Ctx) error {
	if ctx.Event.GuildID == "" {
		return errors.New("cannot send webhook embed in DMs")
	}

	webhooks, err := ctx.Session.ChannelWebhooks(ctx.Event.ChannelID)
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
		webhook, err = ctx.Session.WebhookCreate(ctx.Event.ChannelID, "LightningBot", "https://lightning-bot.com/logo.png")
		if err != nil {
			return err
		}
	}

	_, err = ctx.Session.WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Username:  "LightningBot",
		AvatarURL: "https://lightning-bot.com/logo.png",
		Embeds: []*discordgo.MessageEmbed{
			e.ToEmbed(),
		},
	})

	return err
}

func (e *Embed) SendWebhook(userID, webhookURL string, important bool) *http.HttpResponse {
	cfg, err := config.GetConfig()
	if err != nil {
		return &http.HttpResponse{Error: err}
	}

	content := ""

	switch cfg.Webhooks.PingSelfMode {
	case "Disabled":
		break
	case "Always":
		content += "<@!" + userID + ">\n"
	case "Important":
		if important {
			content += "<@!" + userID + ">\n"
		}
	default:
		cfg.Webhooks.PingSelfMode = "Disabled"
		_ = cfg.Save()
	}

	switch cfg.Webhooks.PingEveryoneMode {
	case "Disabled":
		break
	case "Always":
		content += "@everyone\n"
	case "Important":
		if important {
			content += "@everyone\n"
		}
	default:
		cfg.Webhooks.PingEveryoneMode = "Disabled"
		_ = cfg.Save()
	}

	return http.Post(webhookURL, map[string]any{
		"content":    content,
		"embeds":     []any{e},
		"username":   "LightningBot Webhook Notifications",
		"avatar_url": "https://lightning-bot.com/logo.png",
	})
}

func deleteMessagesAfter(session *discordgo.Session, messages []*discordgo.Message, delay int) {
	time.Sleep(time.Duration(delay) * time.Second)

	for _, message := range messages {
		_ = session.ChannelMessageDelete(message.ChannelID, message.ID)
	}
}
