package formatting

import (
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
)

// TimestampFlag is a flag for formatting timestamps.
type TimestampFlag string

const (
	// TimestampFlagShortTime is a flag for short time format.
	TimestampFlagShortTime TimestampFlag = "t"

	// TimestampFlagLongTime is a flag for long time format.
	TimestampFlagLongTime TimestampFlag = "T"

	// TimestampFlagShortDate is a flag for short date format.
	TimestampFlagShortDate TimestampFlag = "d"

	// TimestampFlagLongDate is a flag for long date format.
	TimestampFlagLongDate TimestampFlag = "D"

	// TimestampFlagShortDateTime is a flag for short date and time format.
	TimestampFlagShortDateTime TimestampFlag = "f"

	// TimestampFlagLongDateTime is a flag for long date and time format.
	TimestampFlagLongDateTime TimestampFlag = "F"

	// TimestampFlagRelative is a flag for relative time format.
	TimestampFlagRelative TimestampFlag = "R"
)

// Timestamp formats a time.Time value into a Discord timestamp string.
// Only the first flag is used if multiple are provided.
// Format: <t:unix> without flags; <t:unix:flag> with a flag.
func Timestamp(t time.Time, flags ...TimestampFlag) string {
	var builder strings.Builder

	builder.WriteString("<t:")
	builder.WriteString(strconv.FormatInt(t.Unix(), 10))
	if len(flags) == 1 {
		builder.WriteString(":")
		builder.WriteString(string(flags[0]))
	}
	builder.WriteString(">")

	return builder.String()
}

var permissionNames = map[int]string{
	discordgo.PermissionAddReactions:         "Add Reactions",
	discordgo.PermissionAdministrator:        "**Administrator**",
	discordgo.PermissionAttachFiles:          "Attach Files",
	discordgo.PermissionBanMembers:           "Ban Members",
	discordgo.PermissionChangeNickname:       "Change Nickname",
	discordgo.PermissionVoiceConnect:         "Connect to Voice Channel",
	discordgo.PermissionCreateInstantInvite:  "Create Instant Invite",
	discordgo.PermissionVoiceDeafenMembers:   "Deafen Members",
	discordgo.PermissionEmbedLinks:           "Embed Links",
	discordgo.PermissionUseExternalEmojis:    "Use External Emojis",
	discordgo.PermissionKickMembers:          "Kick Members",
	discordgo.PermissionManageChannels:       "Manage Channels",
	discordgo.PermissionManageEmojis:         "Manage Emojis",
	discordgo.PermissionManageServer:         "Manage Server",
	discordgo.PermissionManageMessages:       "Manage Messages",
	discordgo.PermissionManageNicknames:      "Manage Nicknames",
	discordgo.PermissionManageRoles:          "Manage Roles",
	discordgo.PermissionManageWebhooks:       "Manage Webhooks",
	discordgo.PermissionMentionEveryone:      "Mention Everyone",
	discordgo.PermissionVoiceMoveMembers:     "Move Members",
	discordgo.PermissionVoiceMuteMembers:     "Mute Members",
	discordgo.PermissionVoicePrioritySpeaker: "Priority Speaker",
	discordgo.PermissionReadMessageHistory:   "Read Message History",
	discordgo.PermissionViewChannel:          "View Channel",
	discordgo.PermissionSendMessages:         "Send Messages",
	discordgo.PermissionSendTTSMessages:      "Send TTS messages",
	discordgo.PermissionVoiceSpeak:           "Speak in Voice Channel",
	discordgo.PermissionVoiceStreamVideo:     "Stream in Voice Channel",
	discordgo.PermissionVoiceUseVAD:          "Use Voice Activation",
	discordgo.PermissionViewAuditLogs:        "View Audit Logs",
	discordgo.PermissionViewGuildInsights:    "View Guild Insights",
}

// Permissions formats the permissions of a role into a string.
func Permissions(permissions int64) string {
	var builder strings.Builder

	for permission, name := range permissionNames {
		if permissions&int64(permission) != 0 {
			builder.WriteString("- ")
			builder.WriteString(name)
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
