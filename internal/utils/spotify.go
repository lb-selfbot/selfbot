package utils

import (
	"fmt"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/lb-selfbot/selfbot/internal/http"
)

type SpotifyUtils struct {
	connName        string
	connAccessToken string
	connRetrieved   time.Time
}

func formatAsBearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func (s *SpotifyUtils) GetBearerToken(token string) string {
	if s.connName == "" || time.Since(s.connRetrieved).Minutes() > 10 {
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

			connID := conn.Get("id").MustString()

			conn, err := http.GetJson(
				fmt.Sprintf("https://discord.com/api/v10/users/@me/connections/spotify/%s/access-token", connID),
				map[string]string{"Authorization": token},
			)
			if err != nil {
				return ""
			}

			s.connAccessToken = conn.Get("access_token").MustString()

			break
		}
	}

	return formatAsBearerToken(s.connAccessToken)
}

func (s SpotifyUtils) GetSongArtists(song *simplejson.Json) []string {
	var songArtists []string

	artistAmount := len(song.Get("artists").MustArray())

	for i := 0; i < artistAmount; i++ {
		artistName := song.Get("artists").GetIndex(i).Get("name").MustString()
		songArtists = append(songArtists, artistName)
	}

	return songArtists
}

var Spotify = SpotifyUtils{}
