package sessions

import (
	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var GatewaySessions []*discordgo.GatewaySession
var LoginSessions []*discordgo.LoginSession
var Status string

func RegisterSessionHandler(session *discordgo.Session) {
	parseStatus(session, session.State.Sessions)
	GatewaySessions = filterSessions(session.State.Sessions)

	refreshLoginSessions(session)

	session.AddHandler(onSessionsReplace)
}

func onSessionsReplace(session *discordgo.Session, replace *discordgo.SessionsReplace) {
	parseStatus(session, replace.Sessions)
	GatewaySessions = filterSessions(replace.Sessions)

	refreshLoginSessions(session)
}

func filterSessions(sessions []*discordgo.GatewaySession) []*discordgo.GatewaySession {
	return utils.Filter(sessions, func(session *discordgo.GatewaySession) bool {
		return session.SessionID != "all"
	})
}

func parseStatus(session *discordgo.Session, sessions []*discordgo.GatewaySession) {
	statuses := map[string]int{}

	for _, session := range sessions {
		if session.SessionID == "all" {
			Status = session.Status

			return
		}

		if _, ok := statuses[session.Status]; !ok {
			statuses[session.Status] = 0
		}

		statuses[session.Status]++
	}

	var maxStatus string
	var maxCount int

	for status, count := range statuses {
		if count > maxCount {
			maxCount = count
			maxStatus = status
		}
	}

	if maxStatus == "" {
		maxStatus = getFallbackStatus(session)
	}

	Status = maxStatus
}

func getFallbackStatus(session *discordgo.Session) string {
	status := "idle"

	settings, err := session.UserSettings()
	if err == nil {
		status = settings["status"].(string)
	}

	return status
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
