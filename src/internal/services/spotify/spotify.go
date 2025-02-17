package spotify

import (
	"errors"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/bitly/go-simplejson"
)

type SpotifyServiceImpl struct {
	discordSvc interfaces.DiscordService

	connName        string
	connAccessToken string
	connRetrieved   time.Time

	currentlyPlaying          *simplejson.Json
	currentlyPlayingRetrieved time.Time
}

func (s *SpotifyServiceImpl) Inject(discordSvc interfaces.DiscordService) {
	s.discordSvc = discordSvc
}

func (s *SpotifyServiceImpl) GetBearerToken() string {
	if s.connName != "" && time.Since(s.connRetrieved).Minutes() < 10 {
		return s.formatAsBearerToken(s.connAccessToken)
	}

	session := s.discordSvc.GetSession()

	conns, err := session.UserConnections()
	if err != nil {
		return ""
	}

	for _, conn := range conns {
		if conn.Type != "spotify" {
			continue
		}

		s.connName = conn.Name
		s.connAccessToken = conn.AccessToken
		s.connRetrieved = time.Now()

		data, err := session.Request("GET", fmt.Sprintf("/users/@me/connections/spotify/%s/access-token", s.connName), nil)
		if err != nil {
			return ""
		}

		jsonData, err := simplejson.NewJson(data)
		if err != nil {
			return ""
		}

		s.connAccessToken = jsonData.Get("access_token").MustString()

		break
	}

	return s.formatAsBearerToken(s.connAccessToken)
}

func (s *SpotifyServiceImpl) formatAsBearerToken(token string) string {
	return "Bearer " + token
}

func (s *SpotifyServiceImpl) GetCurrentlyPlaying() (*simplejson.Json, error) {
	if s.currentlyPlaying != nil && time.Since(s.currentlyPlayingRetrieved).Seconds() < 10 {
		return s.currentlyPlaying, nil
	}

	bearerToken := s.GetBearerToken()
	if bearerToken == "" {
		return nil, errors.New("spotify is not connected to your discord account")
	}

	resp, err := http.GetJson(
		"https://api.spotify.com/v1/me/player/currently-playing",
		map[string]string{"Authorization": bearerToken},
	)
	if err != nil {
		return nil, err
	}

	var ok bool
	s.currentlyPlaying, ok = resp.CheckGet("item")
	if !ok {
		return nil, errors.New("no song currently playing on spotify")
	}

	s.currentlyPlayingRetrieved = time.Now()

	return s.currentlyPlaying, nil
}

func (s *SpotifyServiceImpl) GetSongArtists(song *simplejson.Json) []string {
	var songArtists []string

	artistAmount := len(song.Get("artists").MustArray())

	for i := range artistAmount {
		artistName := song.Get("artists").GetIndex(i).Get("name").MustString()
		songArtists = append(songArtists, artistName)
	}

	return songArtists
}
