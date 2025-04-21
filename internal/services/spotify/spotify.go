package spotify

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	fhttp "github.com/bogdanfinn/fhttp"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/http"
)

// Public error variables
var (
	// ErrNotConnected is returned when Spotify is not connected to the Discord account
	ErrNotConnected = errors.New("spotify is not connected to your discord account")

	// ErrNoCurrentlyPlaying is returned when no song is currently playing
	ErrNoCurrentlyPlaying = errors.New("no song is currently playing on spotify")

	// ErrInvalidVolume is returned when the volume is not between 0 and 100
	ErrInvalidVolume = errors.New("volume must be between 0 and 100")

	// ErrInvalidRepeatMode is returned when the repeat mode is not 'off', 'track', or 'context'
	ErrInvalidRepeatMode = errors.New("repeat mode must be 'off', 'track', or 'context'")

	// ErrEmptyQuery is returned when a search query is empty
	ErrEmptyQuery = errors.New("query cannot be empty")

	// ErrNoSongsFound is returned when no songs match the search criteria
	ErrNoSongsFound = errors.New("no songs found")

	// ErrNegativePosition is returned when the position is negative
	ErrNegativePosition = errors.New("position cannot be negative")

	// ErrNoActiveDevice is returned when there is no active device
	ErrNoActiveDevice = errors.New("no active device found. you must be listening to Spotify on a device")
)

// Service handles Spotify-related operations.
type Service interface {
	// GetBearerToken retrieves the Spotify bearer token.
	GetBearerToken(session *discordgo.Session) (string, error)

	// GetProfile retrieves the Spotify profile.
	GetProfile(session *discordgo.Session) (*jester.Data, error)

	// GetDevices retrieves the available devices.
	GetDevices(session *discordgo.Session) (*jester.Data, error)

	// GetCurrentlyPlaying retrieves the currently playing track.
	GetCurrentlyPlaying(session *discordgo.Session) (*jester.Data, error)

	// GetSongArtists retrieves the artists of a song.
	GetSongArtists(song *jester.Data) []string

	// Pause pauses the currently playing track.
	Pause(session *discordgo.Session) error

	// Resume resumes the currently playing track.
	Resume(session *discordgo.Session) error

	// Previous skips to the previous track.
	Previous(session *discordgo.Session) error

	// Next skips to the next track.
	Next(session *discordgo.Session) error

	// SetVolume sets the volume of the currently playing track.
	SetVolume(session *discordgo.Session, volume int) error

	// SetShuffle sets the shuffle mode.
	SetShuffle(session *discordgo.Session, shuffle bool) error

	// SetRepeat sets the repeat mode.
	SetRepeat(session *discordgo.Session, mode string) error

	// Seek seeks to a specific position in the track.
	Seek(session *discordgo.Session, position time.Duration) error

	// Search searches for a track.
	Search(session *discordgo.Session, query string) (*jester.Data, error)

	// Play plays a track.
	Play(session *discordgo.Session, query string) (*jester.Data, error)

	// Queue queues a track.
	Queue(session *discordgo.Session, query string) (*jester.Data, error)
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
	if !s.isValid() {
		return ""
	}

	return "Bearer " + s.accessToken
}

type spotifyService struct {
	conn spotifyConnection

	currentlyPlaying            *jester.Data
	currentlyPlayingRetrievedAt time.Time
}

// NewService creates a new instance of SpotifyService.
func NewService() Service {
	return &spotifyService{}
}

func (s *spotifyService) GetBearerToken(session *discordgo.Session) (string, error) {
	if s.conn.isValid() {
		return s.conn.getBearerToken(), nil
	}

	conns, err := session.UserConnections()
	if err != nil {
		return "", err
	}

	for _, conn := range conns {
		if conn.Type != "spotify" {
			continue
		}

		resp, err := session.Request(
			"GET",
			fmt.Sprintf("https://discord.com/api/v9/users/@me/connections/spotify/%s/access-token", conn.Name),
			nil,
		)
		if err != nil {
			return "", err
		}

		data, err := jester.NewJson(resp)
		if err != nil {
			return "", err
		}

		s.conn = spotifyConnection{
			name:        conn.Name,
			accessToken: data.Get("access_token").MustString(),
			retrievedAt: time.Now(),
		}

		return s.conn.getBearerToken(), nil
	}

	return "", ErrNotConnected
}

func (s *spotifyService) request(session *discordgo.Session, method, endpoint string) (*jester.Data, error) {
	bearerToken, err := s.GetBearerToken(session)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient().Request(
		method,
		"https://api.spotify.com/v1"+endpoint,
		nil,
		http.WithHeader("Authorization", bearerToken),
	)
	if err != nil {
		return nil, err
	}

	data := resp.Json()
	if !resp.IsSuccess() {
		reason := data.Get("error", "reason").MustString()
		if reason == "NO_ACTIVE_DEVICE" {
			return nil, ErrNoActiveDevice
		}

		return nil, fmt.Errorf("failed to get data from Spotify API: %s %v", endpoint, data.Interface())
	}

	return data, nil
}

func (s *spotifyService) GetProfile(session *discordgo.Session) (*jester.Data, error) {
	return s.request(session, fhttp.MethodGet, "/me")
}

func (s *spotifyService) GetDevices(session *discordgo.Session) (*jester.Data, error) {
	data, err := s.request(nil, fhttp.MethodGet, "/me/player/devices")
	if err != nil {
		return nil, err
	}

	return data.Get("devices"), nil
}

func (s *spotifyService) GetCurrentlyPlaying(session *discordgo.Session) (*jester.Data, error) {
	if s.currentlyPlaying != nil && time.Since(s.currentlyPlayingRetrievedAt).Seconds() < 10 {
		return s.currentlyPlaying, nil
	}

	data, err := s.request(session, fhttp.MethodGet, "/me/player/currently-playing")
	if err != nil {
		return nil, err
	}

	s.currentlyPlaying = data.Get("item")
	if s.currentlyPlaying.Interface() == nil {
		return nil, ErrNoCurrentlyPlaying
	}

	s.currentlyPlayingRetrievedAt = time.Now()

	return s.currentlyPlaying, nil
}

func (s *spotifyService) GetSongArtists(song *jester.Data) []string {
	artists := make([]string, 0, song.Get("artists").Len())

	for artist := range song.Get("artists").Iterator() {
		artists = append(artists, artist.Get("name").MustString())
	}

	return artists
}

func (s *spotifyService) Pause(session *discordgo.Session) error {
	_, err := s.request(session, fhttp.MethodPut, "/me/player/pause")
	return err
}

func (s *spotifyService) Resume(session *discordgo.Session) error {
	_, err := s.request(session, fhttp.MethodPut, "/me/player/play")
	return err
}

func (s *spotifyService) Previous(session *discordgo.Session) error {
	_, err := s.request(session, fhttp.MethodPost, "/me/player/previous")
	return err
}

func (s *spotifyService) Next(session *discordgo.Session) error {
	_, err := s.request(session, fhttp.MethodPost, "/me/player/next")
	return err
}

func (s *spotifyService) SetVolume(session *discordgo.Session, volume int) error {
	if volume < 0 || volume > 100 {
		return ErrInvalidVolume
	}

	_, err := s.request(session, fhttp.MethodPut, fmt.Sprintf("/me/player/volume?volume_percent=%d", volume))
	return err
}

func (s *spotifyService) SetShuffle(session *discordgo.Session, shuffle bool) error {
	_, err := s.request(session, fhttp.MethodPut, fmt.Sprintf("/me/player/shuffle?state=%t", shuffle))
	return err
}

func (s *spotifyService) SetRepeat(session *discordgo.Session, mode string) error {
	if mode != "off" && mode != "track" && mode != "context" {
		return ErrInvalidRepeatMode
	}

	_, err := s.request(session, fhttp.MethodPut, fmt.Sprintf("/me/player/repeat?state=%s", url.QueryEscape(mode)))
	return err
}

func (s *spotifyService) Seek(session *discordgo.Session, position time.Duration) error {
	if position < 0 {
		return ErrNegativePosition
	}

	_, err := s.request(session, fhttp.MethodPut, fmt.Sprintf("/me/player/seek?position_ms=%d", int(position.Milliseconds())))
	return err
}

func (s *spotifyService) Search(session *discordgo.Session, query string) (*jester.Data, error) {
	if query == "" {
		return nil, ErrEmptyQuery
	}

	data, err := s.request(session, fhttp.MethodGet, fmt.Sprintf("/search?q=%s&type=track", url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}

	return data.Get("tracks", "items"), nil
}

func (s *spotifyService) Play(session *discordgo.Session, query string) (*jester.Data, error) {
	track, err := s.Queue(session, query)
	if err != nil {
		return nil, err
	}

	// Wait for the queue to be processed
	time.Sleep(100 * time.Millisecond)

	if err := s.Next(session); err != nil {
		return nil, err
	}

	return track, nil
}

func (s *spotifyService) Queue(session *discordgo.Session, query string) (*jester.Data, error) {
	tracks, err := s.Search(session, query)
	if err != nil {
		return nil, err
	}

	if tracks.Len() == 0 {
		return nil, ErrNoSongsFound
	}

	track := tracks.Get(0)
	trackURI := track.Get("uri").MustString()

	_, err = s.request(session, fhttp.MethodPost, fmt.Sprintf("/me/player/queue?uri=%s", url.QueryEscape(trackURI)))
	if err != nil {
		return nil, err
	}

	return track, nil
}
