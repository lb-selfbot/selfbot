package eventhandlers

import "github.com/LightningDev1/discordgo"

func (e *Events) OnGuildMemberUpdate(_ *discordgo.Session, memberUpdate *discordgo.GuildMemberUpdate) {
	if memberUpdate.BeforeUpdate == nil || memberUpdate.Member == nil || memberUpdate.User == nil {
		return
	}

	for _, event := range e.Settings.ForceNick {
		if event.GuildID != memberUpdate.GuildID {
			continue
		}

		if event.UserID != memberUpdate.User.ID {
			continue
		}

		if memberUpdate.BeforeUpdate.Nick != memberUpdate.Nick {
			e.Session.GuildMemberNickname(memberUpdate.GuildID, memberUpdate.User.ID, event.Argument)
		}
	}
}

func (e *Events) OnGuildMemberAdd(_ *discordgo.Session, memberAdd *discordgo.GuildMemberAdd) {
	for _, event := range e.Settings.ForceKick {
		if event.GuildID != memberAdd.GuildID {
			continue
		}

		if event.UserID != memberAdd.User.ID {
			continue
		}

		e.Session.GuildMemberNickname(memberAdd.GuildID, memberAdd.User.ID, event.Argument)
	}
}
