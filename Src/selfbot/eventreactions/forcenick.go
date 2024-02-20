package eventreactions

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type ForceNick struct {
	Guild *discordgo.Guild
	User  *discordgo.User

	Nickname string
}

func NewForceNick(g *discordgo.Guild, u *discordgo.User, nick string) *ForceNick {
	return &ForceNick{Guild: g, User: u, Nickname: nick}
}

func NewForceNickFromData(session *discordgo.Session, data map[string]any) (*ForceNick, error) {
	var guild *discordgo.Guild
	var user *discordgo.User
	var err error

	guildID, ok := data["guild"].(string)
	if ok {
		guild, err = utils.Discord.GetGuild(session, guildID)
		if err != nil {
			guild = &discordgo.Guild{ID: guildID}
		}
	}

	userID, ok := data["user"].(string)
	if ok {
		user, ok = utils.Discord.GetUser(session, userID)
		if !ok {
			user = &discordgo.User{ID: userID}
		}
	}

	nickname, ok := data["nickname"].(string)
	if !ok {
		utils.Logging.Error("ForceNick.NewForceNickFromData: Invalid nickname")
		return nil, fmt.Errorf("invalid nickname")
	}

	return NewForceNick(guild, user, nickname), nil
}

func (f *ForceNick) Type() EventReactionType {
	return FORCE_NICK
}

func (f *ForceNick) Description() string {
	location := f.User.Username

	if f.Guild != nil {
		location = location + " - " + f.Guild.Name
	}

	return "Force Nickname (" + location + ")"
}

func (f *ForceNick) Equals(other EventReaction) bool {
	otherForceNick, ok := other.(*ForceNick)
	if !ok {
		utils.Logging.Debug("ForceNick.Equals: Other is not a ForceNick")
		return false
	}

	return f.Guild == otherForceNick.Guild &&
		f.User == otherForceNick.User &&
		f.Nickname == otherForceNick.Nickname
}

func (f *ForceNick) Data() map[string]any {
	data := make(map[string]any)

	if f.Guild != nil {
		data["guild"] = f.Guild.ID
	}
	if f.User != nil {
		data["user"] = f.User.ID
	}
	data["nickname"] = f.Nickname

	return data
}

func (f *ForceNick) Handler(session *discordgo.Session, event any) {
	memberUpdate, ok := event.(*discordgo.GuildMemberUpdate)
	if !ok {
		utils.Logging.Debug("ForceNick.Handler: Event is not a guild member update event")
		return
	}

	if f.Guild != nil && f.Guild.ID != memberUpdate.GuildID {
		return
	} else if f.User != nil && f.User.ID != memberUpdate.User.ID {
		return
	}

	if memberUpdate.BeforeUpdate.Nick != memberUpdate.Nick {
		_ = session.GuildMemberNickname(memberUpdate.GuildID, memberUpdate.User.ID, f.Nickname)
	}
}
