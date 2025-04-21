package utils

import (
	"fmt"
	"io"
	"math"
	nethttp "net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
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
			if !(m.Type == discordgo.MessageTypeDefault || m.Type == discordgo.MessageTypeReply) {
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

func (DiscordUtils) MuteUser(session *discordgo.Session, member *discordgo.Member) error {
	roles, err := session.GuildRoles(member.GuildID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if strings.ToLower(role.Name) != "muted" {
			continue
		}

		if slices.Contains(member.Roles, role.ID) {
			return fmt.Errorf("user is already muted")
		}

		newRoles := append(member.Roles, role.ID)

		_, err = session.GuildMemberEdit(member.GuildID, member.User.ID, &discordgo.GuildMemberParams{
			Roles: &newRoles,
		})
		return err
	}

	channels, err := session.GuildChannels(member.GuildID)
	if err != nil {
		return err
	}

	role, err := session.GuildRoleCreate(member.GuildID, &discordgo.RoleParams{
		Name: "Muted",
	})
	if err != nil {
		return err
	}

	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		err = session.ChannelPermissionSet(channel.ID, role.ID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			return err
		}
	}

	newRoles := append(member.Roles, role.ID)

	_, err = session.GuildMemberEdit(member.GuildID, member.User.ID, &discordgo.GuildMemberParams{
		Roles: &newRoles,
	})

	return err
}

func (DiscordUtils) UserChannelPermissions(session *discordgo.Session, guildID string, channelID string, userID string) (apermissions int64, err error) {
	userId, _ := strconv.Atoi(userID)

	_, err = session.State.Member(guildID, userID)
	if err != nil {
		params := discordgo.NewQueryGuildMembersParams(guildID)
		params.UserIDs = []int{userId}

		_, _ = session.QueryGuildMembers(params)
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

func (DiscordUtils) FormatStatusEmoji(status discordgo.Status) string {
	switch status {
	case discordgo.StatusOnline:
		return "`🟢` Online"
	case discordgo.StatusIdle:
		return "`🟡` Idle"
	case discordgo.StatusDoNotDisturb:
		return "`🔴` Do Not Disturb"
	case discordgo.StatusInvisible:
		return "`⚪` Invisible"
	case discordgo.StatusOffline:
		return "`⚪` Offline"
	default:
		return "`⚪` Unknown"
	}
}

func (DiscordUtils) FormatClientStatus(clientStatus discordgo.ClientStatus) string {
	return fmt.Sprintf(
		"Desktop: %s\nMobile: %s\nWeb: %s\nConsole: %s\n",
		Discord.FormatStatusEmoji(clientStatus.Desktop),
		Discord.FormatStatusEmoji(clientStatus.Mobile),
		Discord.FormatStatusEmoji(clientStatus.Web),
		Discord.FormatStatusEmoji(clientStatus.Embedded),
	)
}

func (DiscordUtils) FormatUserFlags(flags discordgo.UserFlags) string {
	var formatted []string

	if flags&discordgo.UserFlagDiscordEmployee != 0 {
		formatted = append(formatted, "Discord Employee")
	}
	if flags&discordgo.UserFlagDiscordPartner != 0 {
		formatted = append(formatted, "Discord Partner")
	}
	if flags&discordgo.UserFlagHypeSquadEvents != 0 {
		formatted = append(formatted, "HypeSquad Events")
	}
	if flags&discordgo.UserFlagBugHunterLevel1 != 0 {
		formatted = append(formatted, "Bug Hunter Level 1")
	}
	if flags&discordgo.UserFlagHouseBravery != 0 {
		formatted = append(formatted, "HypeSquad Bravery")
	}
	if flags&discordgo.UserFlagHouseBrilliance != 0 {
		formatted = append(formatted, "HypeSquad Brilliance")
	}
	if flags&discordgo.UserFlagHouseBalance != 0 {
		formatted = append(formatted, "HypeSquad Balance")
	}
	if flags&discordgo.UserFlagEarlySupporter != 0 {
		formatted = append(formatted, "Early Supporter")
	}
	if flags&discordgo.UserFlagTeamUser != 0 {
		formatted = append(formatted, "Team User")
	}
	if flags&discordgo.UserFlagSystem != 0 {
		formatted = append(formatted, "System")
	}
	if flags&discordgo.UserFlagBugHunterLevel2 != 0 {
		formatted = append(formatted, "Bug Hunter Level 2")
	}
	if flags&discordgo.UserFlagVerifiedBot != 0 {
		formatted = append(formatted, "Verified Bot")
	}
	if flags&discordgo.UserFlagVerifiedBotDeveloper != 0 {
		formatted = append(formatted, "Verified Bot Developer")
	}
	if flags&discordgo.UserFlagDiscordCertifiedModerator != 0 {
		formatted = append(formatted, "Discord Certified Moderator")
	}
	if flags&discordgo.UserFlagBotHTTPInteractions != 0 {
		formatted = append(formatted, "Bot HTTP Interactions")
	}
	if flags&discordgo.UserFlagActiveDeveloper != 0 {
		formatted = append(formatted, "Active Developer")
	}

	return strings.Join(formatted, ", ")
}

func (DiscordUtils) FormatPremiumType(premium discordgo.UserPremiumType) string {
	switch premium {
	case discordgo.UserPremiumTypeNone:
		return "None"
	case discordgo.UserPremiumTypeNitroClassic:
		return "Nitro Classic"
	case discordgo.UserPremiumTypeNitro:
		return "Nitro Boost"
	case discordgo.UserPremiumTypeNitroBasic:
		return "Nitro Basic"
	default:
		return "Unknown"
	}
}

func (DiscordUtils) FormatActivities(activities []*discordgo.Activity) string {
	var formatted []string

	for _, activity := range activities {
		var activityFormat string

		switch activity.Type {
		case discordgo.ActivityTypeGame:
			activityFormat = fmt.Sprintf("Playing **%s**", activity.Name)
		case discordgo.ActivityTypeStreaming:
			activityFormat = fmt.Sprintf("Streaming **%s**", activity.Name)
		case discordgo.ActivityTypeListening:
			activityFormat = fmt.Sprintf("Listening to **%s**", activity.Name)
		case discordgo.ActivityTypeWatching:
			activityFormat = fmt.Sprintf("Watching **%s**", activity.Name)
		case discordgo.ActivityTypeCustom:
			activityFormat = fmt.Sprintf("Custom Status: **%s**", activity.State)
		case discordgo.ActivityTypeCompeting:
			activityFormat = fmt.Sprintf("Competing in **%s**", activity.Name)
		default:
			activityFormat = fmt.Sprintf("Unknown Activity: **%s**", activity.Name)
		}

		if activity.Details != "" {
			activityFormat += fmt.Sprintf(" (%s)", activity.Details)
		}

		formatted = append(formatted, activityFormat)
	}

	return strings.Join(formatted, "\n")
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
