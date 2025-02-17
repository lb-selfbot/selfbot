package interfaces

import "github.com/bitly/go-simplejson"

type SpotifyService interface {
	Inject(discordSvc DiscordService)

	GetBearerToken() string
	GetCurrentlyPlaying() (*simplejson.Json, error)
	GetSongArtists(song *simplejson.Json) []string
}
