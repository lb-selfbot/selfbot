package sessions

import (
	"fmt"
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

	_ = refreshLoginSessions(session)

	session.AddHandler(onSessionsReplace)
}

func onSessionsReplace(session *discordgo.Session, replace *discordgo.SessionsReplace) {
	parseStatus(session, replace.Sessions, true)
	newSessions := filterSessions(replace.Sessions)

	// Refresh login sessions if a new session is added
	if len(newSessions) > len(GatewaySessions) {
		_ = refreshLoginSessions(session)
	}

	GatewaySessions = newSessions
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
	return session.State.UserSettings.Status.Status.Value
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
	fmt.Println("refreshLoginSessions")
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
