package sessions

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

type SessionsServiceImpl struct {
	uiSvc      interfaces.UIService
	discordSvc interfaces.DiscordService

	gatewaySessions  []*discordgo.GatewaySession
	loginSessions    []*discordgo.LoginSession
	lastStatusChange time.Time
	status           string
}

func (s *SessionsServiceImpl) Inject(uiSvc interfaces.UIService, discordSvc interfaces.DiscordService) {
	s.uiSvc = uiSvc
	s.discordSvc = discordSvc
}

func (s *SessionsServiceImpl) Register() {
	session := s.discordSvc.GetSession()

	s.parseStatus(session, session.State.Sessions, false)
	s.gatewaySessions = s.filterSessions(session.State.Sessions)

	if err := s.refreshLoginSessions(session); err != nil {
		s.uiSvc.LogError("Failed to refresh login sessions:", err)
	}

	session.AddHandler(s.onSessionsReplace)
}

func (s *SessionsServiceImpl) GetGatewaySessions() []*discordgo.GatewaySession {
	return s.gatewaySessions
}

func (s *SessionsServiceImpl) GetLoginSessions() []*discordgo.LoginSession {
	return s.loginSessions
}

func (s *SessionsServiceImpl) SetLastStatusChange(time time.Time) {
	s.lastStatusChange = time
}

func (s *SessionsServiceImpl) GetStatus() string {
	return s.status
}

func (s *SessionsServiceImpl) onSessionsReplace(session *discordgo.Session, e *discordgo.SessionsReplace) {
	s.parseStatus(session, e.Sessions, true)
	newSessions := s.filterSessions(e.Sessions)

	// Refresh login sessions if a new session is added
	if len(newSessions) > len(s.gatewaySessions) {
		if err := s.refreshLoginSessions(session); err != nil {
			s.uiSvc.LogError("Failed to refresh login sessions:", err)
		}
	}

	s.gatewaySessions = newSessions
}

func (s *SessionsServiceImpl) filterSessions(sessions []*discordgo.GatewaySession) []*discordgo.GatewaySession {
	return utils.Filter(sessions, func(session *discordgo.GatewaySession) bool {
		return session.SessionID != "all"
	})
}

func (s *SessionsServiceImpl) parseStatus(session *discordgo.Session, sessions []*discordgo.GatewaySession, isReplace bool) {
	statuses := make(map[string]int)

	for _, gatewaySession := range sessions {
		if gatewaySession.SessionID == session.State.SessionID {
			continue
		}

		if gatewaySession.SessionID == "all" {
			s.setStatus(session, gatewaySession.Status, isReplace)
			return
		}

		if _, ok := statuses[gatewaySession.Status]; !ok {
			statuses[gatewaySession.Status] = 0
		}

		statuses[gatewaySession.Status]++
	}

	maxStatus := s.getMostCommonStatus(statuses)
	if maxStatus == "" {
		maxStatus = s.getFallbackStatus(session)
	}

	s.setStatus(session, maxStatus, isReplace)
}

func (s *SessionsServiceImpl) getMostCommonStatus(statuses map[string]int) string {
	var maxStatus string
	var maxCount int

	for status, count := range statuses {
		if count > maxCount {
			maxCount = count
			maxStatus = status
		}
	}
	return maxStatus
}

func (s *SessionsServiceImpl) setStatus(session *discordgo.Session, status string, isReplace bool) {
	if s.status == status {
		return
	}

	s.status = status
	s.lastStatusChange = time.Now()

	if isReplace {
		_ = session.UpdateStatusComplex(discordgo.UpdateStatusData{
			AFK:    true,
			Status: status,
		})
	}
}

func (s *SessionsServiceImpl) getFallbackStatus(session *discordgo.Session) string {
	return session.State.UserSettings.Status.Status.Value
}

func (s *SessionsServiceImpl) refreshLoginSessions(session *discordgo.Session) error {
	sessions, err := session.LoginSessions()
	if err != nil {
		return err
	}

	if len(s.loginSessions) > 0 {
		newSessions := s.findNewSessions(sessions)
		removedSessions := s.findRemovedSessions(sessions)
		s.handleChangedSessions(newSessions, removedSessions)
	}

	s.loginSessions = sessions
	return nil
}

func (s *SessionsServiceImpl) findNewSessions(currentSessions []*discordgo.LoginSession) []*discordgo.LoginSession {
	return utils.Filter(currentSessions, func(session *discordgo.LoginSession) bool {
		for _, loginSession := range s.loginSessions {
			if loginSession.IDHash == session.IDHash {
				return false
			}
		}
		return true
	})
}

func (s *SessionsServiceImpl) findRemovedSessions(currentSessions []*discordgo.LoginSession) []*discordgo.LoginSession {
	return utils.Filter(s.loginSessions, func(session *discordgo.LoginSession) bool {
		for _, loginSession := range currentSessions {
			if loginSession.IDHash == session.IDHash {
				return false
			}
		}
		return true
	})
}

func (s *SessionsServiceImpl) handleChangedSessions(newSessions, removedSessions []*discordgo.LoginSession) {
	for _, session := range newSessions {
		s.uiSvc.LogSessionConnected(session.IDHash, session.ClientInfo.OS, session.ClientInfo.Platform)
	}

	for _, session := range removedSessions {
		s.uiSvc.LogSessionDisconnected(session.IDHash, session.ClientInfo.OS, session.ClientInfo.Platform)
	}
}
