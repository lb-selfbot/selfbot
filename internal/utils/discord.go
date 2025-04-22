package utils

import (
	"fmt"
	"io"
	"math"
	nethttp "net/http"
	"os"
	"path"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fs"
	"github.com/lb-selfbot/selfbot/internal/http"
)

type DiscordUtils struct{}

func (DiscordUtils) GetChannelMessages(session *discordgo.Session, channelID string, limit int) ([]*discordgo.Message, error) {
	channel, err := session.State.Channel(channelID)
	if err == nil && len(channel.Messages) >= limit {
		return channel.Messages, nil
	}

	// The maximum number of messages that can be retrieved at once is 100.
	if limit <= 100 {
		return session.ChannelMessages(channelID, limit, "", "", "")
	}

	amountOfRequests := int(math.Ceil(float64(limit) / 100))
	var messages []*discordgo.Message

	for range amountOfRequests {
		lastID := ""
		if len(messages) > 0 {
			lastID = messages[len(messages)-1].ID
		}
		newMessages, err := session.ChannelMessages(channelID, int(math.Min(float64(limit), 100)), lastID, "", "")
		if err != nil {
			return nil, err
		}
		limit -= len(newMessages)
		messages = append(messages, newMessages...)
	}

	return messages, nil
}

func (DiscordUtils) GetChannelMessagesRange(session *discordgo.Session, channelID, fromID, toID string) ([]*discordgo.Message, error) {
	messages, err := session.ChannelMessages(channelID, 100, toID, fromID, "")
	if err != nil {
		return nil, err
	}

	messages = append(messages, &discordgo.Message{ID: fromID})

	if len(messages) < 100 {
		return messages, nil
	}

	lastID := messages[len(messages)-1].ID

	for {
		newMessages, err := session.ChannelMessages(channelID, 100, toID, lastID, "")
		if err != nil {
			return nil, err
		}

		if len(newMessages) == 0 {
			break
		}

		messages = append(messages, newMessages...)
		lastID = messages[len(messages)-1].ID
	}

	return messages, nil
}

func (DiscordUtils) GetChannelMessage(session *discordgo.Session, channelID, ID string) (*discordgo.Message, error) {
	messages, err := session.ChannelMessages(channelID, 100, "", "", ID)
	if err != nil {
		return nil, err
	}

	for _, m := range messages {
		if m.ID == ID {
			return m, nil
		}
	}

	return nil, fmt.Errorf("message not found")
}

func (DiscordUtils) DeleteMessages(session *discordgo.Session, channelID, userID string, amount int, delay time.Duration) (int, error) {
	var deleted int
	var lastID string

	// Fetch messages in chunks of 100 until desired amount is deleted
	for deleted < amount {
		messages, err := session.ChannelMessages(channelID, 100, lastID, "", "")
		if err != nil {
			return deleted, err
		}

		// No messages left to delete
		if len(messages) == 0 {
			break
		}

		for _, m := range messages {
			lastID = m.ID

			// Check if message is from the desired user
			if m.Author == nil || (userID != "" && m.Author.ID != userID) {
				continue
			}

			// Check if the message is deletable
			if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
				continue
			}

			// Delete message
			err := session.ChannelMessageDelete(channelID, m.ID)
			if err != nil {
				return deleted, err
			}

			deleted++

			time.Sleep(delay)

			// Check if the desired amount is deleted
			if deleted >= amount {
				break
			}
		}
	}

	return deleted, nil
}

func (DiscordUtils) SendFileFromURL(session *discordgo.Session, channelID, url, fileName string, headers ...map[string]string) error {
	var reqData http.RequestData
	if len(headers) > 0 {
		reqData.Headers = headers[0]
	}

	resp, err := http.DefaultClient().Get(url, reqData)
	if err != nil {
		return err
	}

	_, err = session.ChannelFileSend(
		channelID,
		fileName,
		resp.Body,
	)
	return err
}

func (DiscordUtils) UserChannelPermissions(session *discordgo.Session, guildID string, channelID string, userID string) (apermissions int64, err error) {
	_, err = session.State.Member(guildID, userID)
	if err != nil {
		_, _ = session.QueryMember(guildID, userID, true)
	}

	return session.State.UserChannelPermissions(userID, channelID)
}

func (DiscordUtils) ImpersonateMessage(s *discordgo.Session, user *discordgo.User, message *discordgo.Message, channelID string) (*discordgo.Message, error) {
	webhooks, err := s.ChannelWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	var webhook *discordgo.Webhook
	if len(webhooks) == 0 {
		webhook, err = s.WebhookCreate(channelID, "Impersonate", "")
		if err != nil {
			return nil, err
		}
	} else {
		webhook = webhooks[0]
	}

	return s.WebhookExecute(webhook.ID, webhook.Token, true, &discordgo.WebhookParams{
		Username:  user.Username,
		AvatarURL: user.AvatarURL(""),
		Content:   message.Content,
		Embeds:    message.Embeds,
	})
}

func (DiscordUtils) Impersonate(s *discordgo.Session, user *discordgo.User, content, channelID string) (*discordgo.Message, error) {
	return Discord.ImpersonateMessage(s, user, &discordgo.Message{Content: content}, channelID)
}

func (DiscordUtils) IsDM(session *discordgo.Session, message *discordgo.Message) bool {
	if message == nil {
		return false
	}

	stateChannel, err := session.State.Channel(message.ChannelID)
	if err == nil {
		return stateChannel.Type == discordgo.ChannelTypeDM
	}

	// Fallback check
	return message.GuildID == "" && // No server
		(message.Type == discordgo.MessageTypeDefault || message.Type == discordgo.MessageTypeReply) && // Normal message
		message.Flags&discordgo.MessageFlagsEphemeral == 0 // Not ephemeral
}

func (DiscordUtils) SaveAttachment(attachment *discordgo.MessageAttachment) (string, error) {
	filename := fs.SanitizeFilename(attachment.ID + "-" + attachment.Filename)
	filepath := path.Join(fs.GetMainFolder(), "Downloads", filename)

	var resp *nethttp.Response
	var err error

	resp, err = nethttp.Get(attachment.URL)
	if err != nil {
		var err error
		resp, err = nethttp.Get(attachment.ProxyURL)
		if err != nil {
			return "", fmt.Errorf("error downloading attachment: %w", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error downloading attachment: %s", resp.Status)
	}

	outFile, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("error creating attachment file: %w", err)
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error downloading attachment: %w", err)
	}

	return filepath, nil
}

var Discord = &DiscordUtils{}
