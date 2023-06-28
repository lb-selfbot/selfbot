package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SpotifyCommands struct{}

func (s SpotifyCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Spotify", "Spotify control commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "spause",
		Description: "Pause Spotify Music",
		Usage:       "[p]spause",
		Handler:     s.spause,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sresume",
		Description: "Resume Spotify Music",
		Usage:       "[p]sresume",
		Handler:     s.sresume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sprevious",
		Description: "Play the previous song",
		Usage:       "[p]sprevious",
		Handler:     s.sprevious,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snext",
		Description: "Play the next song",
		Usage:       "[p]snext",
		Handler:     s.snext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "splay",
		Description: "Play a song on Spotify",
		Usage:       "[p]splay <song name>",
		Handler:     s.splay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "squeue",
		Description: "Queue a song to play after the current one",
		Usage:       "[p]squeue <song name>",
		Handler:     s.squeue,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "splaying",
		Description: "Current playing song info",
		Usage:       "[p]splaying",
		Handler:     s.splaying,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "svolume",
		Description: "Change Spotify volume",
		Usage:       "[p]svolume <volume>",
		Handler:     s.svolume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "spresence",
		Description: "Set a fake Spotify presence (use simageid command for image id)",
		Usage:       "[p]spresence <on/off> [\"artist\"] [\"song\"] [\"album\"] [image id]",
		Handler:     s.spresence,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "simageid",
		Description: "Get an image id for the spresence command",
		Usage:       "[p]simageid <song name>",
		Handler:     s.simageid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sshuffle",
		Description: "Turn shuffle on or off",
		Usage:       "[p]sshuffle <on/off>",
		Handler:     s.sshuffle,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "srepeat",
		Description: "Change the Spotify repeat mode",
		Usage:       "[p]srepeat <track/context/off>",
		Handler:     s.srepeat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sseek",
		Description: "Seek to a specific time in the song",
		Usage:       "[p]sseek <seconds>",
		Handler:     s.sseek,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sbearertoken",
		Description: "Get your Spotify bearer token used for the API",
		Usage:       "[p]sbearertoken",
		Handler:     s.sbearertoken,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "slistenalong",
		Description: "Send a Spotify listen along invite",
		Usage:       "[p]slistenalong [content]",
		Handler:     s.slistenalong,
	})
}

func (s SpotifyCommands) GetBearerToken(ctx *dgc.Ctx) (string, bool) {
	bearerToken := utils.Spotify.GetBearerToken(ctx.Session)

	if bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
		return "", false
	}

	return bearerToken, true
}

func (s SpotifyCommands) checkRequest(ctx *dgc.Ctx, httpResponse *http.HttpResponse) bool {
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return false
	}

	errorReason := httpResponse.Json().GetPath("error", "reason").MustString()
	errorReason = strings.ReplaceAll(errorReason, "NO_ACTIVE_DEVICE", "No active device found")

	if errorReason != "" {
		_ = ctx.RespondText(fmt.Sprintf("A Spotify error occured: `%s`", errorReason))
		return false
	}

	return true
}

func (s SpotifyCommands) HandleRequest(ctx *dgc.Ctx, httpResponse *http.HttpResponse, onSuccess string) {
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	if httpResponse.Body != "" {
		if !s.checkRequest(ctx, httpResponse) {
			return
		}
	}
	_ = ctx.RespondText(onSuccess)
}

func (s SpotifyCommands) spause(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Put("https://api.spotify.com/v1/me/player/pause", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully paused music.")
}

func (s SpotifyCommands) sresume(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Put("https://api.spotify.com/v1/me/player/play", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully resumed music.")
}

func (s SpotifyCommands) sprevious(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/previous", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully playing previous song.")
}

func (s SpotifyCommands) snext(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Post("https://api.spotify.com/v1/me/player/next", nil, map[string]string{"Authorization": bearerToken})

	s.HandleRequest(ctx, httpResponse, "Successfully playing next song.")
}

func (s SpotifyCommands) splay(ctx *dgc.Ctx) {
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

func (s SpotifyCommands) squeue(ctx *dgc.Ctx) {
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
		fmt.Sprintf("Successfully queued %s - %s.", strings.Join(songArtists, ", "), songName),
	)
}

func (s SpotifyCommands) splaying(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me/player/currently-playing", map[string]string{"Authorization": bearerToken})
	song, success := httpResponse.Json().CheckGet("item")
	if httpResponse.Error != nil || !success {
		_ = ctx.RespondText("A Spotify error occured: No song currently playing")
		return
	}

	songArtists := utils.Spotify.GetSongArtists(song)
	songArtistsString := strings.Join(songArtists, ", ")
	songReleaseDate := song.GetPath("album", "release_date").MustString("Unknown")
	songExplicit := utils.Misc.FormatBoolYesNo(song.Get("explicit").MustBool(false))
	songDuration := fmt.Sprintf(
		"%d:%02d",
		int(song.Get("duration_ms").MustInt(0))/1000/60,
		int(song.Get("duration_ms").MustInt(0))/1000%60,
	)
	songUrl := song.GetPath("external_urls", "spotify").MustString("Unknown")
	songImageUrl := song.GetPath("album", "images").GetIndex(0).Get("url").MustString("")

	embed.NewEmbed().
		AddField("Title", song.Get("name").MustString()).
		AddField("Artist", songArtistsString).
		AddField("Release Date", songReleaseDate).
		AddField("Explicit", songExplicit).
		AddField("Duration", songDuration).
		AddField("URL", songUrl).
		SetImage(songImageUrl).
		Send(ctx)
}

func (s SpotifyCommands) svolume(ctx *dgc.Ctx) {
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

func (s SpotifyCommands) spresence(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"]) == "on"

	if !toggle {
		_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Game: &discordgo.Activity{
				Type: discordgo.ActivityType(4),
			},
			Activities: []*discordgo.Activity{},
		})
		return
	}

	imageID := utils.Json.Escape(utils.Arguments["[image id]"])
	if imageID == "" {
		imageID = "ab67616d0000b2735755e164993798e0c9ef7d7a"
	}

	_ = ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Type: discordgo.ActivityType(2),
				Name: "Spotify",
				Assets: &discordgo.Assets{
					LargeImageID: "spotify:" + imageID,
					LargeText:    utils.Arguments["[\"album\"]"],
				},
				Details: utils.Arguments["[\"song\"]"],
				State:   utils.Arguments["[\"artist\"]"],
				Timestamps: &discordgo.TimeStamps{
					StartTimestamp: time.Now().Unix() - 4140000,
					EndTimestamp:   time.Now().Unix() - 3880000,
				},
				Party: &discordgo.Party{
					ID: "spotify:" + ctx.Session.State.User.ID,
				},
				Flags: 48,
			},
		},
	})
}

func (s SpotifyCommands) simageid(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	headers := map[string]string{"Authorization": bearerToken}

	searchUrl := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(ctx.Arguments.Raw()))

	httpResponse := http.Get(searchUrl, headers)
	json := httpResponse.Json()

	song := json.GetPath("tracks", "items").GetIndex(0)

	imageUrl := song.Get("album").Get("images").GetIndex(0).Get("url").MustString()
	if httpResponse.Error != nil || imageUrl == "" {
		_ = ctx.RespondText("A Spotify error occured: Song not found.")
		return
	}

	splitUrl := strings.Split(imageUrl, "/")

	embed.NewEmbed().
		AddField("Song", song.Get("name").MustString()).
		AddField("Image ID", splitUrl[len(splitUrl)-1]).
		Send(ctx)
}

func (s SpotifyCommands) sshuffle(ctx *dgc.Ctx) {
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

func (s SpotifyCommands) srepeat(ctx *dgc.Ctx) {
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

func (s SpotifyCommands) sseek(ctx *dgc.Ctx) {
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

func (s SpotifyCommands) sbearertoken(ctx *dgc.Ctx) {
	bearerToken, success := s.GetBearerToken(ctx)
	if !success {
		return
	}

	httpResponse := http.Get("https://api.spotify.com/v1/me", map[string]string{"Authorization": bearerToken})
	if !s.checkRequest(ctx, httpResponse) {
		return
	}

	embed.NewEmbed().
		AddField("Spotify Account", httpResponse.Json().Get("display_name").MustString()).
		AddField("Bearer Token", bearerToken).
		Send(ctx)
}

func (s SpotifyCommands) slistenalong(ctx *dgc.Ctx) {
	for _, s := range sessions.GatewaySessions {
		if len(s.Activities) == 0 || s.Activities[0].Name != "Spotify" {
			continue
		}

		_, err := ctx.Session.ChannelMessageSendComplex(ctx.Event.ChannelID, &discordgo.MessageSend{
			Content: ctx.Arguments.Raw(),
			Activity: &discordgo.MessageActivity{
				Type:      discordgo.MessageActivityTypeListen,
				SessionID: s.SessionID,
				PartyID:   s.Activities[0].Party.ID,
			},
			Nonce: discordgo.GenerateNonce(),
			Flags: 0,
		})

		if err == nil {
			break
		}
	}
}
