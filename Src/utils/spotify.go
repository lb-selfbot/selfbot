package utils

import "github.com/bitly/go-simplejson"

type SpotifyUtils struct{}

func (s *SpotifyUtils) GetBearerToken() string {
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