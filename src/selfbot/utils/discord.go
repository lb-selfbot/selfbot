package utils

import (
	"bytes"
	"fmt"
	"io"
	"math"
	nethttp "net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

var Arguments = make(map[string]string)

var Users = sync.Map{}
var nilUser = &discordgo.User{}

type DiscordUtils struct{}

func (DiscordUtils) TryGetUserFromMention(session *discordgo.Session, mention *dgc.Argument) (*discordgo.User, bool) {
	userID := Misc.Or(mention.AsUserMentionID(), mention.Raw())

	if user, ok := Users.Load(userID); ok {
		if user == nilUser {
			return nil, false
		}

		return user.(*discordgo.User), true
	}

	user, err := session.State.GetUser(userID)
	if err != nil {
		var err error

		user, err = session.User(userID)
		if err != nil {
			Users.Store(userID, nilUser)
			return nil, false
		}
	}

	Users.Store(userID, user)

	return user, true
}

func (DiscordUtils) TryGetMemberFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Member, bool) {
	return Discord.GetMember(ctx.Session, ctx.Event.GuildID, Misc.Or(mention.AsUserMentionID(), mention.Raw()))
}

func (d DiscordUtils) TryGetUsernameFromMention(session *discordgo.Session, mention *dgc.Argument) (string, bool) {
	user, success := d.TryGetUserFromMention(session, mention)
	if !success {
		return "", false
	}
	return user.String(), true
}

func (DiscordUtils) GetChannel(session *discordgo.Session, channelID string) (*discordgo.Channel, error) {
	channel, err := session.State.Channel(channelID)
	if err != nil {
		channel, err = session.Channel(channelID)
		if err != nil {
			return nil, err
		}
	}
	return channel, nil
}

func (d DiscordUtils) TryGetChannelFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Channel, bool) {
	channelID := Misc.Or(mention.AsChannelMentionID(), mention.Raw())
	channel, err := d.GetChannel(ctx.Session, channelID)
	if err != nil {
		return nil, false
	}
	return channel, true
}

func (d DiscordUtils) TryGetChannelFromMentionWithDefault(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Channel, bool) {
	channelID := Misc.Or(mention.AsChannelMentionID(), mention.Raw())
	if channelID == "" {
		channelID = ctx.Event.ChannelID
	}

	channel, err := d.GetChannel(ctx.Session, channelID)
	if err != nil {
		return nil, false
	}
	return channel, true
}

func (DiscordUtils) TryGetRoleFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Role, bool) {
	roleID := Misc.Or(mention.AsRoleMentionID(), mention.Raw())
	role, err := ctx.Session.State.Role(ctx.Event.GuildID, roleID)
	if err != nil {
		role, err = ctx.Session.State.Role(ctx.Event.GuildID, roleID)
		if err != nil {
			return nil, false
		}
	}
	return role, true
}

func (DiscordUtils) SnowflakeToTime(snowflake string) time.Time {
	// parse snowflake to int64
	snowflakeInt, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}
	}
	// convert int64 to time
	return time.Unix(((snowflakeInt>>22)+1420070400000)/1000, 0)
}

func (DiscordUtils) GetTimestamp(t *time.Time, defaultString string) string {
	if t == nil {
		return defaultString
	}
	return fmt.Sprintf("<t:%d>", t.Unix())
}

func (DiscordUtils) ParseInvite(invite string) string {
	inviteRegexMatch := string(constants.INVITE_REGEX.Find([]byte(invite)))
	if inviteRegexMatch != "" {
		splitInvite := strings.Split(inviteRegexMatch, "/")
		return splitInvite[len(splitInvite)-1]
	}
	return invite
}

func (DiscordUtils) CheckIfInServer(ctx *dgc.Ctx) bool {
	if ctx.Event.GuildID == "" {
		_ = ctx.RespondText(fmt.Sprintf("Error: The `%s` command is only available in a server.", ctx.Command.Name))
	}
	return ctx.Event.GuildID != ""
}

func (DiscordUtils) GetArguments(ctx *dgc.Ctx) (map[string]string, bool) {
	rawArguments := strings.Replace(ctx.Command.Usage, fmt.Sprintf("[p]%s", ctx.Command.Name), "", -1)
	commandArguments := regexp.MustCompile(`(<(.*?)>|\[(.*?)])`).FindAllString(rawArguments, -1)
	arguments := map[string]string{}
	content := ctx.Event.Message.Content
	for index, argName := range commandArguments {
		raw := ctx.Arguments.Get(index).Raw()
		isRequired := strings.Contains(argName, "<") && strings.Contains(argName, ">")

		if strings.Contains(argName, "\"") && strings.Count(raw, "\"") >= 2 {
			arguments[argName] = regexp.MustCompile("\"(.*?)\"").FindString(content)
			content = strings.Replace(content, arguments[argName], "", -1)
		} else {
			arguments[argName] = raw
			content = String.RemoveWords(content, 1)
		}

		if arguments[argName] == "" && isRequired {
			_ = ctx.RespondText("Missing required argument: " + argName)
			return map[string]string{}, false
		}
	}
	return arguments, true
}

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

	for i := 0; i < amountOfRequests; i++ {
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

func (DiscordUtils) GetRoleNames(ctx *dgc.Ctx, member *discordgo.Member) string {
	var roleNames []string
	for _, role := range member.Roles {
		role, err := ctx.Session.State.Role(member.GuildID, role)
		if err == nil {
			roleNames = append(roleNames, role.Name)
		}
	}
	return strings.Join(roleNames, ", ")
}

func (DiscordUtils) GetGuild(session *discordgo.Session, guildID string) (guild *discordgo.Guild, err error) {
	if guildID == "" {
		return nil, fmt.Errorf("guildID is empty")
	}

	guild, err = session.State.Guild(guildID)
	if err != nil {
		guild, err = session.Guild(guildID)
		if err != nil {
			return nil, err
		}
	}
	return guild, nil
}

func (DiscordUtils) GetMember(session *discordgo.Session, guildID, userID string, reload ...bool) (*discordgo.Member, bool) {
	member, err := session.State.Member(guildID, userID)
	if err != nil || len(reload) > 0 && reload[0] {
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			return nil, false
		}

		params := discordgo.NewQueryGuildMembersParams(guildID)
		params.UserIDs = []int{userIDInt}

		members, err := session.QueryGuildMembers(params)
		if err != nil {
			return nil, false
		}

		if len(members) == 0 {
			// User is not in this server, find mutual servers
			profile, err := session.UserProfile(userID)
			if err != nil {
				return nil, false
			}

			if len(profile.MutualGuilds) == 0 {
				// No mutuals, can't get member
				return nil, false
			}

			return Discord.GetMember(session, profile.MutualGuilds[0].ID, userID)
		}

		member = members[0]
	}

	if member == nil {
		return nil, false
	}

	return member, true
}

func (DiscordUtils) GetBoosters(guild *discordgo.Guild) []*discordgo.Member {
	var boostedMembers []*discordgo.Member
	for _, member := range guild.Members {
		if member.PremiumSince != nil {
			boostedMembers = append(boostedMembers, member)
		}
	}
	return boostedMembers
}

func (DiscordUtils) FormatPermissions(roleName, roleID string, permissions int64) string {
	permissionsString := fmt.Sprintf(`Role: %s | ID: %s

Permission: Add Reactions            | Value: %s
Permission: Administrator            | Value: %s
Permission: Attach Files             | Value: %s
Permission: Ban Members              | Value: %s
Permission: Change Nickname          | Value: %s
Permission: Connect to Voice Channel | Value: %s
Permission: Create Instant Invite    | Value: %s
Permission: Deafen Members           | Value: %s
Permission: Embed Links              | Value: %s
Permission: Use External Emojis      | Value: %s
Permission: Kick Members             | Value: %s
Permission: Manage Channels          | Value: %s
Permission: Manage Emojis            | Value: %s
Permission: Manage Server            | Value: %s
Permission: Manage Messages          | Value: %s
Permission: Manage Nicknames         | Value: %s
Permission: Manage Roles             | Value: %s
Permission: Manage webhooks          | Value: %s
Permission: Mention everyone         | Value: %s
Permission: Move Members             | Value: %s
Permission: Mute Members             | Value: %s
Permission: Priority Speaker         | Value: %s
Permission: Read Message History     | Value: %s
Permission: View Channel             | Value: %s
Permission: Send Messages            | Value: %s
Permission: Send TTS messages        | Value: %s
Permission: Speak in Voice Channel   | Value: %s
Permission: Stream in Voice Channel  | Value: %s
Permission: Use External Emojis      | Value: %s
Permission: Use Voice Activation     | Value: %s
Permission: View Audit Logs          | Value: %s
Permission: View Guild Insights      | Value: %s`,
		roleName, roleID,
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionAddReactions != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionAdministrator != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionAttachFiles != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionBanMembers != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionChangeNickname != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceConnect != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionCreateInstantInvite != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceDeafenMembers != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionEmbedLinks != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionUseExternalEmojis != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionKickMembers != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageChannels != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageEmojis != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageServer != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageMessages != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageNicknames != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageRoles != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionManageWebhooks != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionMentionEveryone != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceMoveMembers != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceMuteMembers != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoicePrioritySpeaker != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionReadMessageHistory != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionViewChannel != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionSendMessages != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionSendTTSMessages != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceSpeak != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceStreamVideo != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionUseExternalEmojis != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionVoiceUseVAD != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionViewAuditLogs != 0),
		Misc.FormatBoolYesNo(permissions&discordgo.PermissionViewGuildInsights != 0),
	)
	return permissionsString
}

func (DiscordUtils) SendFileFromURL(ctx *dgc.Ctx, url string, fileName string, headers ...map[string]string) error {
	httpResponse := http.Get(url, headers...)
	if httpResponse.Error != nil {
		return httpResponse.Error
	}

	_, err := ctx.Session.ChannelFileSend(
		ctx.Event.ChannelID,
		fileName,
		bytes.NewReader(httpResponse.Body),
	)
	return err
}

func (DiscordUtils) MapCategories(channels []*discordgo.Channel) (categories []*discordgo.Channel, categoryChannels map[string][]*discordgo.Channel, looseChannels []*discordgo.Channel) {
	categoryChannels = make(map[string][]*discordgo.Channel)

	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		categories = append(categories, channel)
		categoryChannels[channel.ID] = []*discordgo.Channel{}
		for _, subChannel := range channels {
			if subChannel.ParentID == channel.ID {
				categoryChannels[channel.ID] = append(categoryChannels[channel.ID], subChannel)
			}
		}
	}

	looseChannels = Filter(channels, func(channel *discordgo.Channel) bool {
		return channel.Type != discordgo.ChannelTypeGuildCategory && channel.ParentID == ""
	})

	return categories, categoryChannels, looseChannels
}

func (DiscordUtils) MapPermissionOverwrites(overwrites []*discordgo.PermissionOverwrite, roleIDs map[string]string) []*discordgo.PermissionOverwrite {
	var newOverwrites []*discordgo.PermissionOverwrite

	for _, overwrite := range overwrites {
		id := overwrite.ID

		if newID, ok := roleIDs[overwrite.ID]; ok {
			id = newID
		}

		newOverwrites = append(newOverwrites, &discordgo.PermissionOverwrite{
			ID:    id,
			Type:  overwrite.Type,
			Deny:  overwrite.Deny,
			Allow: overwrite.Allow,
		})
	}

	return newOverwrites
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
		if channel.Type == 1 {
			err = session.ChannelPermissionSet(channel.ID, role.ID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
			if err != nil {
				return err
			}
		}
	}

	newRoles := append(member.Roles, role.ID)

	_, err = session.GuildMemberEdit(member.GuildID, member.User.ID, &discordgo.GuildMemberParams{
		Roles: &newRoles,
	})

	return err
}

func (DiscordUtils) SleepWithTyping(session *discordgo.Session, channelID string, duration time.Duration) {
	for i := 0; i < int(duration.Seconds()/5); i++ {
		time.Sleep(5 * time.Second)
		_ = session.ChannelTyping(channelID)
	}
	time.Sleep(duration - (time.Duration(int(duration.Seconds()/5)) * 5 * time.Second))
}

func (DiscordUtils) SleepWithTypingFix(session *discordgo.Session, channelID string, duration time.Duration) {
	for i := 0; i < int(duration.Seconds()/5); i++ {
		_ = session.ChannelTyping(channelID)
		time.Sleep(5 * time.Second)
	}
	time.Sleep(duration - (time.Duration(int(duration.Seconds()/5)) * 5 * time.Second))
}

func (DiscordUtils) GetGuildName(session *discordgo.Session, guildID string) string {
	guild, _ := Discord.GetGuild(session, guildID)
	guildName := "No Server"
	if guild != nil {
		guildName = guild.Name
	}
	return guildName
}

func (DiscordUtils) GetChannelName(session *discordgo.Session, channelID string) string {
	channel, _ := Discord.GetChannel(session, channelID)
	channelName := "No Channel"
	if channel != nil {
		channelName = channel.Name
	}
	return channelName
}

func (DiscordUtils) GetUser(session *discordgo.Session, userID string) (*discordgo.User, bool) {
	arg := &dgc.Argument{}
	return Discord.TryGetUserFromMention(session, arg.SetRaw(userID))
}

func (DiscordUtils) GetUserName(session *discordgo.Session, userID string) string {
	arg := &dgc.Argument{}
	userName, _ := Discord.TryGetUsernameFromMention(session, arg.SetRaw(userID))
	return Misc.Or(userName, "Unknown")
}

func (DiscordUtils) IsGuildInList(session *discordgo.Session, guildID string, guildsList []string) bool {
	guild, err := Discord.GetGuild(session, guildID)
	if err != nil {
		return false
	}

	return slices.Contains(guildsList, guild.ID) || slices.Contains(guildsList, guild.Name)
}

func (DiscordUtils) IsUserInList(session *discordgo.Session, userID string, usersList []string) bool {
	return slices.Contains(usersList, userID)
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

func (DiscordUtils) GetMessageLink(m *discordgo.Message) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", Misc.Or(m.GuildID, "@me"), m.ChannelID, m.ID)
}

func (DiscordUtils) SaveAttachment(attachment *discordgo.MessageAttachment) (string, error) {
	filename := IO.SanitizeFilename(attachment.ID + "-" + attachment.Filename)
	filepath := path.Join(IO.GetMainFolder(), "Downloads", filename)

	outFile, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("error creating attachment file: %w", err)
	}

	defer outFile.Close()

	var resp *nethttp.Response
	resp, err = nethttp.Get(attachment.URL)
	if err != nil {
		var err error
		resp, err = nethttp.Get(attachment.ProxyURL)
		if err != nil {
			return "", fmt.Errorf("error downloading attachment: %w", err)
		}
	}

	defer resp.Body.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error downloading attachment: %w", err)
	}

	return filepath, nil
}

func (DiscordUtils) IsUserMentioned(message *discordgo.Message, userID string) bool {
	for _, user := range message.Mentions {
		if user.ID == userID {
			return true
		}
	}

	return false
}

func (DiscordUtils) WebhookURL(id, token string) string {
	return fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", id, token)
}

var Discord = &DiscordUtils{}
