package utils

import (
	"fmt"

	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
)

type SpotifyUtils struct{}

func formatAsBearerToken(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

func (s *SpotifyUtils) GetBearerToken(session *discordgo.Session) string {
	response, err := session.Request("GET", "https://discord.com/api/v10/users/@me/connections", nil)
	if err != nil {
		return ""
	}

	connections, err := simplejson.NewJson(response)
	if err != nil {
		return ""
	}

	for i := 0; i < len(connections.MustArray()); i++ {
		connection := connections.GetIndex(i)
		connectionType := connection.Get("type").MustString()
		connectionName := connection.Get("name").MustString()
		connectionAccessToken := connection.Get("access_token").MustString()

		if connectionType != "spotify" {
			continue
		}

		response, err = session.Request(
			"GET",
			fmt.Sprintf("https://discord.com/api/v10/users/@me/connections/spotify/%s/access-token", connectionName),
			nil,
		)
		if err != nil {
			return formatAsBearerToken(connectionAccessToken)
		}

		json, err := simplejson.NewJson(response)
		if err != nil {
			return formatAsBearerToken(connectionAccessToken)
		}

		return formatAsBearerToken(json.Get("access_token").MustString())
	}

	return ""
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
