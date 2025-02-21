package spotify

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

func (s SpotifyCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Spotify", "Spotify control commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "spause",
		Description: "Pause Spotify Music",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.spause,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sresume",
		Description: "Resume Spotify Music",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sresume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sprevious",
		Description: "Play the previous song",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sprevious,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snext",
		Description: "Play the next song",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.snext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "splay",
		Description: "Play a song on Spotify",
		Arguments:   "<song name>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.splay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "squeue",
		Description: "Queue a song to play after the current one",
		Arguments:   "<song name>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.squeue,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "splaying",
		Description: "Current playing song info",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.splaying,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "svolume",
		Description: "Change Spotify volume",
		Arguments:   "<volume>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.svolume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sshuffle",
		Description: "Turn shuffle on or off",
		Arguments:   "<on/off>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sshuffle,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "srepeat",
		Description: "Change the Spotify repeat mode",
		Arguments:   "<track/context/off>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.srepeat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sseek",
		Description: "Seek to a specific time in the song",
		Arguments:   "<seconds>",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sseek,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sdevices",
		Description: "Get your available Spotify devices",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sdevices,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sprofile",
		Description: "Get your Spotify profile info",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sprofile,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sbearertoken",
		Description: "Get your Spotify bearer token used for the API",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.sbearertoken,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "slistenalong",
		Description: "Send a Spotify listen along invite",
		Arguments:   "[content]",
		Flags:       []router.CommandFlag{router.CmdFlagDangerous},
		Handler:     s.slistenalong,
	})
}

func (s SpotifyCommands) GetBearerToken(ctx *router.Ctx) (string, bool) {
	bearerToken := utils.Spotify.GetBearerToken(ctx.Session.Identify.Token)

	if bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
		return "", false
	}

	return bearerToken, true
}

func (s SpotifyCommands) checkResponse(ctx *router.Ctx, httpResponse *http.HttpResponse) bool {
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return false
	}

	errorReason := httpResponse.Json().GetPath("error", "reason").MustString()
	errorReason = strings.ReplaceAll(errorReason, "NO_ACTIVE_DEVICE", "No active device found")

	if errorReason != "" {
		_ = ctx.RespondText(fmt.Sprintf("A Spotify error occurred: `%s`", errorReason))
		return false
	}

	return true
}

func (s SpotifyCommands) HandleRequest(ctx *router.Ctx, httpResponse *http.HttpResponse, onSuccess string) {
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	// Check the response if the body isn't empty
	if len(httpResponse.Body) != 0 && !s.checkResponse(ctx, httpResponse) {
		return
	}

	_ = ctx.RespondText(onSuccess)
}

func (s SpotifyCommands) spause(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Put("https://api.spotify.com/v1/me/player/pause", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully paused music.")
}

func (s SpotifyCommands) sresume(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Put("https://api.spotify.com/v1/me/player/play", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully resumed music.")
}

func (s SpotifyCommands) sprevious(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/previous", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully playing previous song.")
}

func (s SpotifyCommands) snext(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/next", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully playing next song.")
}

func (s SpotifyCommands) splay(ctx *router.Ctx) {
	songName := ctx.Arguments.Raw()

	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	headers := map[string]string{"Authorization": bearerToken}

	searchResponse := http.Get(
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)),
		headers,
	)
	if searchResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	search := searchResponse.Json().GetPath("tracks", "items")
	if len(search.MustArray()) == 0 {
		_ = ctx.RespondText("No songs found.")
		return
	}

	song := search.GetIndex(0)
	songName = song.Get("name").MustString()
	songUrl := song.Get("external_urls").Get("spotify").MustString()
	songArtists := utils.Spotify.GetSongArtists(song)

	_ = http.Post(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", url.QueryEscape(songUrl)),
		nil,
		headers,
	)

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/next", nil, headers)

	s.HandleRequest(
		ctx,
		httpResponse,
		fmt.Sprintf("Successfully playing %s - %s.", strings.Join(songArtists, ", "), songName),
	)
}

func (s SpotifyCommands) squeue(ctx *router.Ctx) {
	songName := ctx.Arguments.Raw()

	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	headers := map[string]string{"Authorization": bearerToken}

	searchResponse := http.Get(
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)),
		headers,
	)
	if searchResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	search := searchResponse.Json().GetPath("tracks", "items")
	if len(search.MustArray()) == 0 {
		_ = ctx.RespondText("No songs found.")
		return
	}

	song := search.GetIndex(0)
	songName = song.Get("name").MustString()
	songUrl := song.Get("external_urls").Get("spotify").MustString()
	songArtists := utils.Spotify.GetSongArtists(song)

	httpResponse := http.Post(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", url.QueryEscape(songUrl)),
		nil,
		headers,
	)

	s.HandleRequest(
		ctx,
		httpResponse,
		fmt.Sprintf("Successfully queued %s - %s.", strings.Join(songArtists, ", "), songName),
	)
}

func (s SpotifyCommands) splaying(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me/player/currently-playing", map[string]string{"Authorization": bearerToken})
	song, success := httpResponse.Json().CheckGet("item")
	if httpResponse.Error != nil || !success {
		_ = ctx.RespondText("A Spotify error occurred: No song currently playing")
		return
	}

	songArtists := utils.Spotify.GetSongArtists(song)
	songArtistsString := strings.Join(songArtists, ", ")
	songReleaseDate := song.GetPath("album", "release_date").MustString("Unknown")
	songExplicit := utils.Misc.FormatBoolYesNo(song.Get("explicit").MustBool(false))
	duration := int(song.Get("duration_ms").MustInt(0)) / 1000
	songDuration := fmt.Sprintf("%d:%02d", duration/60, duration&60)
	songUrl := song.GetPath("external_urls", "spotify").MustString("Unknown")
	songImageUrl := song.GetPath("album", "images").GetIndex(0).Get("url").MustString("")

	c.services.Embed.NewEmbed().
		AddField("Title", song.Get("name").MustString()).
		AddField("Artist", songArtistsString).
		AddField("Release Date", songReleaseDate).
		AddField("Explicit", songExplicit).
		AddField("Duration", songDuration).
		AddField("URL", songUrl).
		SetImage(songImageUrl).
		Send(ctx)
}

func (s SpotifyCommands) svolume(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	volumeStr := utils.Arguments["<volume>"]
	volume, err := strconv.Atoi(volumeStr)
	if err != nil || volume < 0 || volume > 100 {
		_ = ctx.RespondText("Volume must be between 0 and 100.")
		return
	}

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/volume?volume_percent=%s", volumeStr),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	s.HandleRequest(ctx, httpResponse, fmt.Sprintf("Successfully set volume to: %s.", volumeStr))
}

func (s SpotifyCommands) sshuffle(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	state := strings.ReplaceAll(toggle, "on", "true")
	state = strings.ReplaceAll(state, "off", "false")

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/shuffle?state=%s", state),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	s.HandleRequest(ctx, httpResponse, fmt.Sprintf("Successfully turned shuffle %s.", toggle))
}

func (s SpotifyCommands) srepeat(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	state := strings.ToLower(utils.Arguments["<track/context/off>"])
	if state != "track" && state != "context" && state != "off" {
		_ = ctx.RespondText("Invalid value. Use `track`, `context` or `off`.")
		return
	}

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/repeat?state=%s", state),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	s.HandleRequest(ctx, httpResponse, fmt.Sprintf("Successfully set Spotify repeat mode to %s.", state))
}

func (s SpotifyCommands) sseek(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", seconds*1000),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	s.HandleRequest(ctx, httpResponse, fmt.Sprintf("Sucessfully seeked song to %d seconds.", seconds))
}

func (s SpotifyCommands) sdevices(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me/player/devices", map[string]string{"Authorization": bearerToken})
	if !s.checkResponse(ctx, httpResponse) {
		return
	}

	devices := httpResponse.Json().Get("devices")
	deviceLen := len(devices.MustArray())

	var deviceList []string

	for i := 0; i < deviceLen; i++ {
		device := devices.GetIndex(i)

		deviceData := fmt.Sprintf(
			"**%s**\nType: %s\nActive: %s\nVolume: %d%%\n",
			device.Get("name").MustString(),
			device.Get("type").MustString(),
			utils.Misc.FormatBoolYesNo(device.Get("is_active").MustBool(false)),
			device.Get("volume_percent").MustInt(0),
		)

		deviceList = append(deviceList, deviceData)
	}

	c.services.Embed.NewEmbed().
		AddField("Devices", strings.Join(deviceList, "\n")).
		Send(ctx)
}

func (s SpotifyCommands) sprofile(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me", map[string]string{"Authorization": bearerToken})
	if !s.checkResponse(ctx, httpResponse) {
		return
	}

	json := httpResponse.Json()

	c.services.Embed.NewEmbed().
		AddField("Username", json.Get("display_name").MustString()).
		AddField("ID", json.Get("id").MustString()).
		AddField("Country", json.Get("country").MustString()).
		AddField("Product", json.Get("product").MustString()).
		AddField("Followers", fmt.Sprint(json.GetPath("followers", "total").MustInt())).
		AddField("URL", json.GetPath("href").MustString()).
		Send(ctx)
}

func (s SpotifyCommands) sbearertoken(ctx *router.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me", map[string]string{"Authorization": bearerToken})
	if !s.checkResponse(ctx, httpResponse) {
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Spotify Account", httpResponse.Json().Get("display_name").MustString()).
		AddField("Bearer Token", bearerToken).
		Send(ctx)
}

func (s SpotifyCommands) slistenalong(ctx *router.Ctx) {
	spotifySessions := utils.Filter(sessions.GatewaySessions, func(s *discordgo.GatewaySession) bool {
		return len(s.Activities) != 0 && s.Activities[0].Name == "Spotify"
	})

	if len(spotifySessions) == 0 {
		_ = ctx.RespondText("No Spotify sessions found.")
		return
	}

	_, err := ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Content: ctx.Arguments.Raw(),
		Activity: &discordgo.MessageActivity{
			Type:      discordgo.MessageActivityTypeListen,
			SessionID: spotifySessions[0].SessionID,
			PartyID:   spotifySessions[0].Activities[0].Party.ID,
		},
		Nonce: discordgo.GenerateNonce(),
		Flags: 0,
	})
	if err != nil {
		_ = ctx.RespondText("An error has occurred while sending the invite: " + err.Error())
		return
	}
}

type SpotifyCommands struct{ services *services.Services }

func NewSpotifyCommands(s *services.Services) *SpotifyCommands {
	return &SpotifyCommands{services: s}
}
func (c *SpotifyCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

}
