package sessions

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var GatewaySessions []*discordgo.GatewaySession
var LoginSessions []*discordgo.LoginSession
var LastStatusChange time.Time
var Status string

func RegisterSessionHandler(session *discordgo.Session) {
	parseStatus(session, session.State.Sessions, false)
	GatewaySessions = filterSessions(session.State.Sessions)

	refreshLoginSessions(session)

	session.AddHandler(onSessionsReplace)
}

func onSessionsReplace(session *discordgo.Session, replace *discordgo.SessionsReplace) {
	parseStatus(session, replace.Sessions, true)
	GatewaySessions = filterSessions(replace.Sessions)

	refreshLoginSessions(session)
}

func filterSessions(sessions []*discordgo.GatewaySession) []*discordgo.GatewaySession {
	return utils.Filter(sessions, func(session *discordgo.GatewaySession) bool {
		return session.SessionID != "all"
	})
}

func parseStatus(session *discordgo.Session, sessions []*discordgo.GatewaySession, isReplace bool) {
	statuses := map[string]int{}

	for _, gatewaySession := range sessions {
		if gatewaySession.SessionID == session.State.SessionID {
			continue
		}

		if gatewaySession.SessionID == "all" {
			setStatus(session, gatewaySession.Status, isReplace)

			return
		}

		if _, ok := statuses[gatewaySession.Status]; !ok {
			statuses[gatewaySession.Status] = 0
		}

		statuses[gatewaySession.Status]++
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

	setStatus(session, maxStatus, isReplace)
}

func getFallbackStatus(session *discordgo.Session) string {
	status := "idle"

	settings, err := session.UserSettings()
	if err == nil {
		status = settings["status"].(string)
	}

	return status
}

func setStatus(session *discordgo.Session, status string, isReplace bool) {
	if status == Status {
		return
	}

	if time.Since(LastStatusChange) < time.Second*5 {
		return
	}

	Status = status

	if isReplace {
		_ = session.UpdateStatusComplex(discordgo.UpdateStatusData{
			AFK:    true,
			Status: status,
		})
	}
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
