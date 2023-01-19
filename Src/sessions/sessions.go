package sessions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

var Sessions []*UserSessions

type ClientInfo struct {
	OS       string `json:"os"`
	Platform string `json:"platform"`
	Location string `json:"location"`
}

type UserSessions struct {
	IDHash     string     `json:"id_hash"`
	ClientInfo ClientInfo `json:"client_info"`

	LastUsedTime time.Time
}

// Parse the time string into a time.Time object when json is unmarshalled
func (u *UserSessions) UnmarshalJSON(data []byte) error {
	type Alias UserSessions
	aux := &struct {
		ApproxLastUsedTime string `json:"approx_last_used_time"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339, aux.ApproxLastUsedTime)
	if err != nil {
		return err
	}

	u.LastUsedTime = parsedTime

	return nil
}

type SessionData struct {
	UserSessions []*UserSessions `json:"user_sessions"`
}

func GetSessions(session *discordgo.Session) ([]*UserSessions, error) {
	response, err := session.Request("GET", "https://discord.com/api/v10/auth/sessions", nil)
	if err != nil {
		return nil, err
	}

	sessionData := &SessionData{}

	err = json.Unmarshal(response, sessionData)
	if err != nil {
		return nil, err
	}

	userSessions := sessionData.UserSessions

	// Sort the sessions by last used time
	for i := 0; i < len(userSessions); i++ {
		for j := 0; j < len(userSessions); j++ {
			if userSessions[i].LastUsedTime.Before(userSessions[j].LastUsedTime) {
				userSessions[i], userSessions[j] = userSessions[j], userSessions[i]
			}
		}
	}

	return userSessions, nil
}

func CheckSessions(session *discordgo.Session) {
	for {
		newSessions, err := GetSessions(session)
		if err != nil {
			time.Sleep(time.Second * 10)
			continue
		}

		HandleSessions(newSessions)

		Sessions = newSessions

		time.Sleep(time.Second * 10)
	}
}

func StartCheckingSessions(session *discordgo.Session) {
	go CheckSessions(session)
}

func HandleSessions(newSessions []*UserSessions) {
	// oldSessions will be empty at the beginning, we don't want to
	// detect new sessions if they were already connected before us
	if len(Sessions) == 0 {
		return
	}

	oldSessionIDHashes := getHashes(Sessions)

	newSessionIDHashes := getHashes(newSessions)

	for _, oldSession := range Sessions {
		if utils.SliceContains(newSessionIDHashes, oldSession.IDHash) {
			continue
		}

		logging.Logging.SessionDisconnected(
			oldSession.IDHash,
			oldSession.ClientInfo.OS,
			oldSession.ClientInfo.Platform,
		)
	}

	for _, newSession := range newSessions {
		if utils.SliceContains(oldSessionIDHashes, newSession.IDHash) {
			continue
		}

		logging.Logging.SessionConnected(
			newSession.IDHash,
			newSession.ClientInfo.OS,
			newSession.ClientInfo.Platform,
		)
	}
}

func ParseLastUsedTime(t time.Time) string {
	// format as x minutes/hours/days/months ago

	diff := time.Since(t)

	// if the difference is less than a minute, return "Just now"
	if diff.Minutes() < 1 {
		return "Less than a minute ago"
	}

	// if the difference is less than an hour, return "x minutes ago"
	if diff.Hours() < 1 {
		minutesInt := int(diff.Minutes())

		if minutesInt == 1 {
			return "1 minute ago"
		}

		return fmt.Sprintf("%d minutes ago", minutesInt)
	}

	// if the difference is less than a day, return "x hours ago"
	if diff.Hours() < 24 {
		hoursInt := int(diff.Hours())

		if hoursInt == 1 {
			return "1 hour ago"
		}

		return fmt.Sprintf("%d hours ago", hoursInt)
	}

	// if the difference is less than a month, return "x days ago"
	if diff.Hours() < 24*30 {
		daysInt := int(diff.Hours() / 24)

		if daysInt == 1 {
			return "1 day ago"
		}

		return fmt.Sprintf("%d days ago", daysInt)
	}

	// if the difference is less than a year, return "x months ago"
	if diff.Hours() < 24*30*12 {
		monthsInt := int(diff.Hours() / 24 / 30)

		if monthsInt == 1 {
			return "1 month ago"
		}

		return fmt.Sprintf("%d months ago", monthsInt)
	}

	// if the difference is more than a year, return the date
	return t.Format("2006-01-02 15:04:05")
}

func getHashes(sessions []*UserSessions) []string {
	var hashes []string

	for _, session := range sessions {
		hashes = append(hashes, session.IDHash)
	}

	return hashes
}
