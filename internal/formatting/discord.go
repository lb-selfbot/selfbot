package formatting

import (
	"fmt"
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

// StatusEmoji formats the status of a user into a string with an emoji.
func StatusEmoji(status discordgo.Status) string {
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

// ClientStatus formats the client status of a user into a string.
func ClientStatus(clientStatus discordgo.ClientStatus) string {
	return fmt.Sprintf(
		"Desktop: %s\nMobile: %s\nWeb: %s\nConsole: %s\n",
		StatusEmoji(clientStatus.Desktop),
		StatusEmoji(clientStatus.Mobile),
		StatusEmoji(clientStatus.Web),
		StatusEmoji(clientStatus.Embedded),
	)
}

var userFlagNames = map[discordgo.UserFlags]string{
	discordgo.UserFlagDiscordEmployee:           "Discord Employee",
	discordgo.UserFlagDiscordPartner:            "Discord Partner",
	discordgo.UserFlagHypeSquadEvents:           "HypeSquad Events",
	discordgo.UserFlagBugHunterLevel1:           "Bug Hunter Level 1",
	discordgo.UserFlagHouseBravery:              "HypeSquad Bravery",
	discordgo.UserFlagHouseBrilliance:           "HypeSquad Brilliance",
	discordgo.UserFlagHouseBalance:              "HypeSquad Balance",
	discordgo.UserFlagEarlySupporter:            "Early Supporter",
	discordgo.UserFlagTeamUser:                  "Team User",
	discordgo.UserFlagSystem:                    "System",
	discordgo.UserFlagBugHunterLevel2:           "Bug Hunter Level 2",
	discordgo.UserFlagVerifiedBot:               "Verified Bot",
	discordgo.UserFlagVerifiedBotDeveloper:      "Verified Bot Developer",
	discordgo.UserFlagDiscordCertifiedModerator: "Discord Certified Moderator",
	discordgo.UserFlagBotHTTPInteractions:       "Bot HTTP Interactions",
	discordgo.UserFlagActiveDeveloper:           "Active Developer",
}

func UserFlags(flags discordgo.UserFlags) string {
	var formatted []string

	for flag, name := range userFlagNames {
		if flags&flag != 0 {
			formatted = append(formatted, name)
		}
	}

	return strings.Join(formatted, ", ")
}

func PremiumType(premium discordgo.UserPremiumType) string {
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

func Activities(activities []*discordgo.Activity) string {
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
