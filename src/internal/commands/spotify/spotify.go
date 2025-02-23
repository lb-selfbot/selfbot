package spotify

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/internal/utils2"
	"github.com/LightningDev1/discordgo"
)

type SpotifyCommands struct{ services *services.Services }

func NewSpotifyCommands(s *services.Services) *SpotifyCommands {
	return &SpotifyCommands{services: s}
}

func (c *SpotifyCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Spotify",
		Description: "Control your Spotify music",
	}, []*router.Command{
		{
			Name:        "spause",
			Description: "Pause Spotify Music",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.spause,
		},
		{
			Name:        "sresume",
			Description: "Resume Spotify Music",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sresume,
		},
		{
			Name:        "sprevious",
			Description: "Play the previous song",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sprevious,
		},
		{
			Name:        "snext",
			Description: "Play the next song",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.snext,
		},
		{
			Name:        "splay",
			Description: "Play a song on Spotify",
			Arguments:   "<song name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.splay,
		},
		{
			Name:        "squeue",
			Description: "Queue a song to play after the current one",
			Arguments:   "<song name>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.squeue,
		},
		{
			Name:        "splaying",
			Description: "Current playing song info",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.splaying,
		},
		{
			Name:        "svolume",
			Description: "Change Spotify volume",
			Arguments:   "<volume>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.svolume,
		},
		{
			Name:        "sshuffle",
			Description: "Turn shuffle on or off",
			Arguments:   "<on/off>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sshuffle,
		},
		{
			Name:        "srepeat",
			Description: "Change the Spotify repeat mode",
			Arguments:   "<track/context/off>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.srepeat,
		},
		{
			Name:        "sseek",
			Description: "Seek to a specific time in the song",
			Arguments:   "<timestamp>",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sseek,
		},
		{
			Name:        "sdevices",
			Description: "Get your available Spotify devices",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sdevices,
		},
		{
			Name:        "sprofile",
			Description: "Get your Spotify profile info",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sprofile,
		},
		{
			Name:        "sbearertoken",
			Description: "Get your Spotify bearer token used for the API",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.sbearertoken,
		},
		{
			Name:        "slistenalong",
			Description: "Send a Spotify listen along invite",
			Arguments:   "[content]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     c.slistenalong,
		},
	})

}

func (c *SpotifyCommands) GetBearerToken(ctx *router.Ctx) (string, bool) {
	bearerToken := utils.Spotify.GetBearerToken(ctx.Session.Identify.Token)
	if bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
		return "", false
	}

	return bearerToken, true
}

func (c *SpotifyCommands) checkResponse(ctx *router.Ctx, httpResponse *http.HttpResponse) bool {
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return false
	}

	errorReason := httpResponse.Json().GetPath("error", "reason").MustString()
	errorReason = strings.ReplaceAll(errorReason, "NO_ACTIVE_DEVICE", "No active device found")

	if errorReason != "" {
		ctx.RespondText("A Spotify error occurred: `%s`", errorReason)
		return false
	}

	return true
}

func (c *SpotifyCommands) HandleRequest(ctx *router.Ctx, httpResponse *http.HttpResponse, onSuccess string) {
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	// Check the response if the body isn't empty
	if len(httpResponse.Body) != 0 && !c.checkResponse(ctx, httpResponse) {
		return
	}

	ctx.RespondText(onSuccess)
}

func (c *SpotifyCommands) spause(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Put("https://api.spotify.com/v1/me/player/pause", nil, map[string]string{"Authorization": bearerToken})

	c.HandleRequest(ctx, httpResponse, "Successfully paused music.")
}

func (c *SpotifyCommands) sresume(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Put("https://api.spotify.com/v1/me/player/play", nil, map[string]string{"Authorization": bearerToken})

	c.HandleRequest(ctx, httpResponse, "Successfully resumed music.")
}

func (c *SpotifyCommands) sprevious(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/previous", nil, map[string]string{"Authorization": bearerToken})

	c.HandleRequest(ctx, httpResponse, "Successfully playing previous song.")
}

func (c *SpotifyCommands) snext(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/next", nil, map[string]string{"Authorization": bearerToken})

	c.HandleRequest(ctx, httpResponse, "Successfully playing next song.")
}

func (c *SpotifyCommands) splay(ctx *router.Ctx) {
	songName := ctx.Arguments.Raw()

	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	headers := map[string]string{"Authorization": bearerToken}

	searchResponse := http.Get(
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)),
		headers,
	)
	if searchResponse.Error != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	search := searchResponse.Json().GetPath("tracks", "items")
	if len(search.MustArray()) == 0 {
		ctx.RespondText("No songs found.")
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

	c.HandleRequest(
		ctx,
		httpResponse,
		fmt.Sprintf("Successfully playing %s - %s.", strings.Join(songArtists, ", "), songName),
	)
}

func (c *SpotifyCommands) squeue(ctx *router.Ctx) {
	songName := ctx.Arguments.Raw()

	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	headers := map[string]string{"Authorization": bearerToken}

	searchResponse := http.Get(
		fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)),
		headers,
	)
	if searchResponse.Error != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	search := searchResponse.Json().GetPath("tracks", "items")
	if len(search.MustArray()) == 0 {
		ctx.RespondText("No songs found.")
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

	c.HandleRequest(
		ctx,
		httpResponse,
		fmt.Sprintf("Successfully queued %s - %s.", strings.Join(songArtists, ", "), songName),
	)
}

func (c *SpotifyCommands) splaying(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me/player/currently-playing", map[string]string{"Authorization": bearerToken})
	song, success := httpResponse.Json().CheckGet("item")
	if httpResponse.Error != nil || !success {
		ctx.RespondText("A Spotify error occurred: No song currently playing")
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

func (c *SpotifyCommands) svolume(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	volume, err := ctx.Arguments.Get(0).Int()
	volume = utils2.Clamp(volume, 0, 100)
	if err != nil {
		ctx.RespondText("Argument <volume> must be a number.")
		return
	}

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/volume?volume_percent=%d", volume),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	c.HandleRequest(ctx, httpResponse, fmt.Sprintf("Successfully set volume to: %d.", volume))
}

func (c *SpotifyCommands) sshuffle(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	toggle, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid value. Use `on` or `off`.")
		return
	}

	state := strconv.FormatBool(toggle)

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/shuffle?state=%s", state),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	c.HandleRequest(ctx, httpResponse, fmt.Sprintf("Successfully turned shuffle %s.", toggle))
}

func (c *SpotifyCommands) srepeat(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	state := ctx.Arguments.Get(0).Raw()
	if !slices.Contains([]string{"track", "context", "off"}, state) {
		ctx.RespondText("Invalid value. Use `track`, `context` or `off`.")
		return
	}

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/repeat?state=%s", state),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	c.HandleRequest(ctx, httpResponse, fmt.Sprintf("Successfully set Spotify repeat mode to %s.", state))
}

func (c *SpotifyCommands) sseek(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	timestamp, err := ctx.Arguments.Get(0).Duration()
	if err != nil {
		ctx.RespondText("Invalid value. Use a timestamp like `1m30s`.")
		return
	}

	httpResponse := http.Put(
		fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", timestamp.Milliseconds()),
		nil,
		map[string]string{"Authorization": bearerToken},
	)

	c.HandleRequest(ctx, httpResponse, fmt.Sprintf("Sucessfully seeked song to %s.", timestamp.String()))
}

func (c *SpotifyCommands) sdevices(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me/player/devices", map[string]string{"Authorization": bearerToken})
	if !c.checkResponse(ctx, httpResponse) {
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

func (c *SpotifyCommands) sprofile(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me", map[string]string{"Authorization": bearerToken})
	if !c.checkResponse(ctx, httpResponse) {
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

func (c *SpotifyCommands) sbearertoken(ctx *router.Ctx) {
	bearerToken, success := c.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me", map[string]string{"Authorization": bearerToken})
	if !c.checkResponse(ctx, httpResponse) {
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Spotify Account", httpResponse.Json().Get("display_name").MustString()).
		AddField("Bearer Token", bearerToken).
		Send(ctx)
}

func (c *SpotifyCommands) slistenalong(ctx *router.Ctx) {
	spotifySessions := utils.Filter(c.services.Sessions.GetGatewaySessions(), func(s *discordgo.GatewaySession) bool {
		return len(s.Activities) != 0 && s.Activities[0].Name == "Spotify"
	})

	if len(spotifySessions) == 0 {
		ctx.RespondText("No Spotify sessions found.")
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
		ctx.RespondText("An error has occurred while sending the invite: " + err.Error())
		return
	}
}
