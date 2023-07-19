package utils

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/http"
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
	userId := Misc.Or(mention.AsUserMentionID(), mention.Raw())
	member, err := ctx.Session.State.Member(ctx.Event.GuildID, userId)
	if err != nil {
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return nil, false
		}

		params := discordgo.NewQueryGuildMembersParams(ctx.Event.GuildID)
		params.UserIDs = []int{userIdInt}

		members, err := ctx.Session.QueryGuildMembers(params)
		if err != nil {
			return nil, false
		}

		if len(members) > 0 {
			member = members[0]
		}
	}

	if member == nil {
		return nil, false
	}

	return member, true
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
	for index, arg := range commandArguments {
		raw := ctx.Arguments.Get(index).Raw()
		isRequired := strings.Contains(arg, "<") && strings.Contains(arg, ">")

		if strings.Contains(arg, "\"") && strings.Count(raw, "\"") >= 2 {
			arguments[arg] = regexp.MustCompile("\"(.*?)\"").FindString(content)
			content = strings.Replace(content, arguments[arg], "", -1)
		} else {
			arguments[arg] = raw
			content = String.RemoveWords(content, 1)
		}

		if arguments[arg] == "" && isRequired {
			_ = ctx.RespondText("Missing required argument: " + arg)
			return map[string]string{}, false
		}
	}
	return arguments, true
}

func (DiscordUtils) GetChannelMessages(session *discordgo.Session, channelID string, limit int) ([]*discordgo.Message, error) {
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
	guild, err = session.State.Guild(guildID)
	if err != nil {
		guild, err = session.Guild(guildID)
		if err != nil {
			return nil, err
		}
	}
	return guild, nil
}

func (DiscordUtils) GetMember(session *discordgo.Session, guildID, userID string) (member *discordgo.Member, err error) {
	member, err = session.State.Member(guildID, userID)
	if err != nil {
		member, err = session.GuildMember(guildID, userID)
		if err != nil {
			return nil, err
		}
	}
	return member, nil
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
		bytes.NewReader(httpResponse.BodyBytes),
	)
	return err
}

func (DiscordUtils) MapCategories(channels []*discordgo.Channel) (categories []*discordgo.Channel, categoryChannels map[string][]*discordgo.Channel) {
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
	return categories, categoryChannels
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

		if SliceContains(member.Roles, role.ID) {
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

	return SliceContains(guildsList, guild.ID) || SliceContains(guildsList, guild.Name)
}

func (DiscordUtils) IsUserInList(session *discordgo.Session, userID string, usersList []string) bool {
	return SliceContains(usersList, userID)
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

var Discord = &DiscordUtils{}
