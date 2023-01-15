package eventhandlers

import "github.com/LightningDev1/discordgo"

func (e *Events) OnVoiceStateUpdate(_ *discordgo.Session, voiceStateUpdate *discordgo.VoiceStateUpdate) {
	for _, event := range e.Settings.ForceDisconnect {
		if event.GuildID != voiceStateUpdate.GuildID {
			continue
		}

		if event.UserID != voiceStateUpdate.UserID {
			continue
		}

		_ = e.Session.GuildMemberMove(voiceStateUpdate.GuildID, voiceStateUpdate.UserID, nil)
	}
}
