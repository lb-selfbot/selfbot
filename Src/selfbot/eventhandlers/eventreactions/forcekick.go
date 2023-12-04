package eventreactions

import (
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type ForceKick struct {
	Guild *discordgo.Guild
	User  *discordgo.User
}

func NewForceKick(g *discordgo.Guild, u *discordgo.User) *ForceKick {
	return &ForceKick{Guild: g, User: u}
}

func NewForceKickFromData(session *discordgo.Session, data map[string]any) (*ForceKick, error) {
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

	return NewForceKick(guild, user), nil
}

func (f *ForceKick) Type() EventReactionType {
	return FORCE_KICK
}

func (f *ForceKick) Description() string {
	location := f.User.Username

	if f.Guild != nil {
		location = location + " - " + f.Guild.Name
	}

	return "Force Kick (" + location + ")"
}

func (f *ForceKick) Equals(other EventReaction) bool {
	otherForceKick, ok := other.(*ForceKick)
	if !ok {
		utils.Logging.Debug("ForceKick.Equals: Other is not a ForceKick")
		return false
	}

	return f.Guild == otherForceKick.Guild &&
		f.User == otherForceKick.User
}

func (f *ForceKick) Data() map[string]any {
	data := make(map[string]any)

	if f.Guild != nil {
		data["guild"] = f.Guild.ID
	}
	if f.User != nil {
		data["user"] = f.User.ID
	}

	return data
}

func (f *ForceKick) Handler(session *discordgo.Session, event any) {
	memberAdd, ok := event.(*discordgo.GuildMemberAdd)
	if !ok {
		utils.Logging.Debug("ForceKick.Handler: Event is not a guild member update event")
		return
	}

	if memberAdd.User.ID == session.State.User.ID {
		return
	}

	if f.Guild != nil && f.Guild.ID != memberAdd.GuildID {
		return
	} else if f.User != nil && f.User.ID != memberAdd.User.ID {
		return
	}

	_ = session.GuildMemberDeleteWithReason(memberAdd.GuildID, memberAdd.User.ID, "LightningBot - Force Kick")
}
