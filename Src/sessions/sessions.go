package sessions

import (
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var GatewaySessions []*discordgo.GatewaySession
var LoginSessions []*discordgo.LoginSession

func RegisterSessionHandler(session *discordgo.Session) {
	GatewaySessions = session.State.Sessions
	refreshLoginSessions(session)

	session.AddHandler(onSessionsReplace)
}

func onSessionsReplace(session *discordgo.Session, replace *discordgo.SessionsReplace) {
	GatewaySessions = replace.Sessions

	refreshLoginSessions(session)
}

func refreshLoginSessions(session *discordgo.Session) error {
	sessions, err := session.LoginSessions()
	if err != nil {
		return err
	}

	if len(LoginSessions) != 0 {
		newSessions := utils.Filter(sessions, func(session *discordgo.LoginSession) bool {
			for _, loginSession := range LoginSessions {
				if loginSession.IDHash == session.IDHash {
					return false
				}
			}
			return true
		})

		removedSessions := utils.Filter(LoginSessions, func(session *discordgo.LoginSession) bool {
			for _, loginSession := range sessions {
				if loginSession.IDHash == session.IDHash {
					return false
				}
			}
			return true
		})

		handleChangedSessions(newSessions, removedSessions)
	}

	LoginSessions = sessions
	return nil
}

func handleChangedSessions(newSessions, removedSessions []*discordgo.LoginSession) {
	for _, session := range newSessions {
		logging.Logging.SessionConnected(session.IDHash, session.ClientInfo.OS, session.ClientInfo.Platform)
	}

	for _, session := range removedSessions {
		logging.Logging.SessionDisconnected(session.IDHash, session.ClientInfo.OS, session.ClientInfo.Platform)
	}
}
