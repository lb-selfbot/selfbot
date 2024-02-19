package spotify

import (
	"errors"
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/bitly/go-simplejson"
)

type SpotifyAPI struct {
	connName        string
	connAccessToken string
	connRetrieved   time.Time

	currentlyPlaying          *simplejson.Json
	currentlyPlayingRetrieved time.Time
}

func (s *SpotifyAPI) formatAsBearerToken(token string) string {
	return "Bearer " + token
}

func (s *SpotifyAPI) GetBearerToken(token string) string {
	if s.connName != "" && time.Since(s.connRetrieved).Minutes() < 10 {
		return s.formatAsBearerToken(s.connAccessToken)
	}

	conns, err := http.GetJson(
		"https://discord.com/api/v10/users/@me/connections",
		map[string]string{"Authorization": token},
	)
	if err != nil {
		return ""
	}

	for i := 0; i < len(conns.MustArray()); i++ {
		conn := conns.GetIndex(i)
		if conn.Get("type").MustString() != "spotify" {
			continue
		}

		s.connName = conn.Get("name").MustString()
		s.connAccessToken = conn.Get("access_token").MustString()
		s.connRetrieved = time.Now()

		conn, err := http.GetJson(
			fmt.Sprintf("https://discord.com/api/v10/users/@me/connections/spotify/%s/access-token", s.connName),
			map[string]string{"Authorization": token},
		)
		if err != nil {
			return ""
		}

		s.connAccessToken = conn.Get("access_token").MustString()

		break
	}

	return s.formatAsBearerToken(s.connAccessToken)
}

func (s *SpotifyAPI) GetCurrentlyPlaying(token string) (*simplejson.Json, error) {
	if s.currentlyPlaying != nil && time.Since(s.currentlyPlayingRetrieved).Seconds() < 10 {
		return s.currentlyPlaying, nil
	}

	bearerToken := s.GetBearerToken(token)
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

func (s *SpotifyAPI) GetSongArtists(song *simplejson.Json) []string {
	var songArtists []string

	artistAmount := len(song.Get("artists").MustArray())

	for i := 0; i < artistAmount; i++ {
		artistName := song.Get("artists").GetIndex(i).Get("name").MustString()
		songArtists = append(songArtists, artistName)
	}

	return songArtists
}

var API = &SpotifyAPI{}
