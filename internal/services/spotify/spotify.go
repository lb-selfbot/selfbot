package spotify

import (
	"errors"
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/http"
)

// SpotifyService handles Spotify-related operations.
type SpotifyService interface {
	// GetBearerToken retrieves the Spotify bearer token.
	GetBearerToken(session *discordgo.Session) string

	// GetCurrentlyPlaying retrieves the currently playing track.
	GetCurrentlyPlaying(session *discordgo.Session) (*simplejson.Json, error)

	// GetSongArtists retrieves the artists of a song.
	GetSongArtists(song *simplejson.Json) []string
}

type spotifyConnection struct {
	name        string
	accessToken string
	retrievedAt time.Time
}

func (s spotifyConnection) isValid() bool {
	return s.name != "" && time.Since(s.retrievedAt).Minutes() < 10
}

func (s spotifyConnection) getBearerToken() string {
	return "Bearer " + s.accessToken
}

type spotifyService struct {
	conn spotifyConnection

	currentlyPlaying            *simplejson.Json
	currentlyPlayingRetrievedAt time.Time
}

// NewSpotifyService creates a new instance of SpotifyService.
func NewSpotifyService() SpotifyService {
	return &spotifyService{}
}

func (s *spotifyService) GetBearerToken(session *discordgo.Session) string {
	if s.conn.isValid() {
		return s.conn.getBearerToken()
	}

	conns, err := session.UserConnections()
	if err != nil {
		return ""
	}

	for _, conn := range conns {
		if conn.Type != "spotify" {
			continue
		}

		s.conn = spotifyConnection{
			name:        conn.Name,
			accessToken: conn.AccessToken,
			retrievedAt: time.Now(),
		}

		data, err := session.Request("GET", fmt.Sprintf("/users/@me/connections/spotify/%s/access-token", s.conn.name), nil)
		if err != nil {
			return ""
		}

		jsonData, err := simplejson.NewJson(data)
		if err != nil {
			return ""
		}

		s.conn.accessToken = jsonData.Get("access_token").MustString()
		break
	}

	return s.conn.getBearerToken()
}

func (s *spotifyService) GetCurrentlyPlaying(session *discordgo.Session) (*simplejson.Json, error) {
	if s.currentlyPlaying != nil && time.Since(s.currentlyPlayingRetrievedAt).Seconds() < 10 {
		return s.currentlyPlaying, nil
	}

	bearerToken := s.GetBearerToken(session)
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

	s.currentlyPlayingRetrievedAt = time.Now()

	return s.currentlyPlaying, nil
}

func (s *spotifyService) GetSongArtists(song *simplejson.Json) []string {
	var songArtists []string

	artistAmount := len(song.Get("artists").MustArray())

	for i := range artistAmount {
		artistName := song.Get("artists").GetIndex(i).Get("name").MustString()
		songArtists = append(songArtists, artistName)
	}

	return songArtists
}
