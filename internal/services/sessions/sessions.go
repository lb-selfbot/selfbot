package sessions

import (
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services/eventbroker"
	"github.com/lb-selfbot/selfbot/internal/services/ui/featurelogs"
)

// Service handles the fetching of Discord sessions.
type Service interface {
	// GetGatewaySessions returns the list of gateway sessions.
	GetGatewaySessions() []*discordgo.GatewaySession

	// GetLoginSessions returns the list of login sessions.
	GetLoginSessions() []*discordgo.LoginSession

	// GetStatus returns the current status.
	GetStatus() string

	// SetStatus sets the current status.
	SetStatus(status string)
}

type sessionsService struct {
	eventBroker eventbroker.Service

	gatewaySessions  []*discordgo.GatewaySession
	loginSessions    []*discordgo.LoginSession
	lastStatusChange time.Time
	status           string
}

// NewService creates a new instance of SessionsService.
func NewService(eventBroker eventbroker.Service) Service {
	svc := &sessionsService{eventBroker: eventBroker}
	svc.subscribe()
	return svc
}

func (s *sessionsService) subscribe() {
	s.eventBroker.Subscribe(eventbroker.EventSessionReady, func(payload any) {
		session, ok := payload.(*discordgo.Session)
		if !ok {
			return
		}

		s.register(session)
	})
}

func (s *sessionsService) register(session *discordgo.Session) {
	s.parseStatus(session, session.State.Sessions, false)
	s.gatewaySessions = s.filterSessions(session.State.Sessions)

	if err := s.refreshLoginSessions(session); err != nil {
		s.eventBroker.Publish(eventbroker.EventLogError, "Failed to refresh login sessions: "+err.Error())
	}

	session.AddHandler(s.onSessionsReplace)
}

func (s *sessionsService) GetGatewaySessions() []*discordgo.GatewaySession {
	return s.gatewaySessions
}

func (s *sessionsService) GetLoginSessions() []*discordgo.LoginSession {
	return s.loginSessions
}

func (s *sessionsService) GetStatus() string {
	return s.status
}

func (s *sessionsService) SetStatus(status string) {
	s.status = status
	s.lastStatusChange = time.Now()
}

func (s *sessionsService) onSessionsReplace(session *discordgo.Session, e *discordgo.SessionsReplace) {
	s.parseStatus(session, e.Sessions, true)
	newSessions := s.filterSessions(e.Sessions)

	// Refresh login sessions if a new session is added
	if len(newSessions) > len(s.gatewaySessions) {
		if err := s.refreshLoginSessions(session); err != nil {
			s.eventBroker.Publish(eventbroker.EventLogError, "Failed to refresh login sessions: "+err.Error())
		}
	}

	s.gatewaySessions = newSessions
}

func (s *sessionsService) filterSessions(sessions []*discordgo.GatewaySession) []*discordgo.GatewaySession {
	return fn.Filter(sessions, func(session *discordgo.GatewaySession) bool {
		return session.SessionID != "all"
	})
}

func (s *sessionsService) parseStatus(session *discordgo.Session, sessions []*discordgo.GatewaySession, isReplace bool) {
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

func (s *sessionsService) getMostCommonStatus(statuses map[string]int) string {
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

func (s *sessionsService) setStatus(session *discordgo.Session, status string, isReplace bool) {
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

func (s *sessionsService) getFallbackStatus(session *discordgo.Session) string {
	return session.State.UserSettings.Status.Status.Value
}

func (s *sessionsService) refreshLoginSessions(session *discordgo.Session) error {
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

func (s *sessionsService) findNewSessions(currentSessions []*discordgo.LoginSession) []*discordgo.LoginSession {
	return fn.Filter(currentSessions, func(session *discordgo.LoginSession) bool {
		for _, loginSession := range s.loginSessions {
			if loginSession.IDHash == session.IDHash {
				return false
			}
		}
		return true
	})
}

func (s *sessionsService) findRemovedSessions(currentSessions []*discordgo.LoginSession) []*discordgo.LoginSession {
	return fn.Filter(s.loginSessions, func(session *discordgo.LoginSession) bool {
		for _, loginSession := range currentSessions {
			if loginSession.IDHash == session.IDHash {
				return false
			}
		}
		return true
	})
}

func (s *sessionsService) handleChangedSessions(newSessions, removedSessions []*discordgo.LoginSession) {
	for _, session := range newSessions {
		s.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.SessionConnected{
			ID:       session.IDHash,
			OS:       session.ClientInfo.OS,
			Client:   session.ClientInfo.Platform,
			Location: session.ClientInfo.Location,
		})
	}

	for _, session := range removedSessions {
		s.eventBroker.Publish(eventbroker.EventLogFeature, &featurelogs.SessionDisconnected{
			ID:       session.IDHash,
			OS:       session.ClientInfo.OS,
			Client:   session.ClientInfo.Platform,
			Location: session.ClientInfo.Location,
		})
	}
}
