package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/json"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/bitly/go-simplejson"
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
}

func (s SpotifyCommands) GetBearerToken() string {
	cfg, err := config.GetConfig()
	if err != nil {
		return ""
	}
	discordHeaders := http.GetDiscordHeaders(cfg.Account.Discord.Token)
	connections, err := http.GetJson("https://discord.com/api/v10/users/@me/connections", discordHeaders)
	if err != nil {
		return ""
	}
	for _, connectionInterface := range connections.MustArray() {
		connection := connectionInterface.(map[string]interface{})
		if connType := connection["type"]; connType != nil && connType.(string) == "spotify" {
			bearerToken := fmt.Sprintf("Bearer %s", connection["access_token"].(string))
			res, err := http.GetJson("https://api.spotify.com/v1/search?q=Despacito&type=track", map[string]string{"Authorization": bearerToken})
			if err != nil {
				return ""
			}
			if _, success := res.GetPath("tracks", "items").GetIndex(0).Get("external_urls").CheckGet("spotify"); success {
				return bearerToken
			}
			res, err = http.GetJson(fmt.Sprintf("https://discord.com/api/v10/users/@me/connections/spotify/%s/access-token", connection["name"].(string)), discordHeaders)
			if err != nil {
				return ""
			}
			return "Bearer " + res.Get("access_token").MustString()
		}
	}
	return ""
}

func (s SpotifyCommands) checkRequest(ctx *dgc.Ctx, response *simplejson.Json, err error) bool {
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return false
	}
	errorReason := response.GetPath("error", "reason").MustString()
	if errorReason == "NO_ACTIVE_DEVICE" {
		_ = ctx.RespondText("A Spotify error occured: `No active device found`")
		return false
	}
	return true
}

func (s SpotifyCommands) HandleRequest(ctx *dgc.Ctx, response string, err error, onSuccess string) {
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	if response != "" {
		responseJson, err := json.SimpleUnmarshal(response)
		if !s.checkRequest(ctx, responseJson, err) {
			return
		}
	}
	_ = ctx.RespondText(onSuccess)
}

func (s SpotifyCommands) spause(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Put("https://api.spotify.com/v1/me/player/pause", map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, response, err, "Successfully paused music.")
	}
}

func (s SpotifyCommands) sresume(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Put("https://api.spotify.com/v1/me/player/play", map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, response, err, "Successfully resumed music.")
	}
}

func (s SpotifyCommands) sprevious(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Post("https://api.spotify.com/v1/me/player/previous", "", map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, response, err, "Successfully playing previous song.")
	}
}

func (s SpotifyCommands) snext(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Post("https://api.spotify.com/v1/me/player/next", "", map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, response, err, "Successfully playing next song.")
	}
}

func (s SpotifyCommands) splay(ctx *dgc.Ctx) {
	songName := ctx.Arguments.Raw()
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		headers := map[string]string{"Authorization": bearerToken}
		searchResponse, err := http.GetJson(fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)), headers)
		if err != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}
		search := searchResponse.GetPath("tracks", "items").MustArray()
		if len(search) == 0 {
			_ = ctx.RespondText("No songs found.")
			return
		}
		song := search[0].(map[string]interface{})
		songUrl := song["external_urls"].(map[string]interface{})["spotify"].(string)
		songName := song["name"].(string)
		var songArtists []string
		for _, artist := range song["artists"].([]interface{}) {
			songArtists = append(songArtists, artist.(map[string]interface{})["name"].(string))
		}
		_, _ = http.Post(fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", url.QueryEscape(songUrl)), "", headers)
		response, err := http.Post("https://api.spotify.com/v1/me/player/next", "", headers)
		s.HandleRequest(ctx, response, err, fmt.Sprintf("Successfully playing %s - %s.", strings.Join(songArtists, ", "), songName))
	}
}

func (s SpotifyCommands) squeue(ctx *dgc.Ctx) {
	songName := ctx.Arguments.Raw()
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		headers := map[string]string{"Authorization": bearerToken}
		searchResponse, err := http.GetJson(fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)), headers)
		if err != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}
		search := searchResponse.GetPath("tracks", "items").MustArray()
		if len(search) == 0 {
			_ = ctx.RespondText("No songs found.")
			return
		}
		song := search[0].(map[string]interface{})
		songUrl := song["external_urls"].(map[string]interface{})["spotify"].(string)
		songName := song["name"].(string)
		var songArtists []string
		for _, artist := range song["artists"].([]interface{}) {
			songArtists = append(songArtists, artist.(map[string]interface{})["name"].(string))
		}
		response, err := http.Post(fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", url.QueryEscape(songUrl)), "", headers)
		s.HandleRequest(ctx, response, err, fmt.Sprintf("Successfully queued %s - %s.", strings.Join(songArtists, ", "), songName))
	}
}

func (s SpotifyCommands) splaying(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.GetJson("https://api.spotify.com/v1/me/player/currently-playing", map[string]string{"Authorization": bearerToken})
		song, success := response.CheckGet("item")
		if err != nil || !success {
			_ = ctx.RespondText("A Spotify error occured: No song currently playing")
			return
		}
		var songArtists []string
		for _, artist := range song.Get("artists").MustArray() {
			songArtists = append(songArtists, artist.(map[string]interface{})["name"].(string))
		}
		embed.NewEmbed().
			AddField("Title", song.Get("name").MustString()).
			AddField("Artist", strings.Join(songArtists, ", ")).
			AddField("Release Date", song.GetPath("album", "release_date").MustString("Unknown")).
			AddField("Explicit", utils.Misc.FormatBool(song.Get("explicit").MustBool(false))).
			AddField("Duration", fmt.Sprintf("%d:%02d", int(song.Get("duration_ms").MustInt(0))/1000/60, int(song.Get("duration_ms").MustInt(0))/1000%60)).
			AddField("URL", song.GetPath("external_urls", "spotify").MustString("Unknown")).
			SetImage(song.GetPath("album", "images").GetIndex(0).Get("url").MustString("")).
			Send(ctx)
	}
}

func (s SpotifyCommands) svolume(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		volumeStr := utils.Arguments["<volume>"]
		volume, err := strconv.Atoi(volumeStr)
		if err != nil || volume < 0 || volume > 100 {
			_ = ctx.RespondText("Volume must be between 0 and 100.")
			return
		}
		res, err := http.Put("https://api.spotify.com/v1/me/player/volume?volume_percent="+volumeStr, map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, res, err, fmt.Sprintf("Successfully set volume to: %s.", volumeStr))
	}
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
	imageID := json.Escape(utils.Arguments["[\"image_id\"]"])
	if imageID == "" {
		imageID = "ab67616d0000b2735755e164993798e0c9ef7d7a"
	}
	data := fmt.Sprintf(`{
		"op": 3,
		"d": {
			"status": "",
			"since": 0,
			"afk": false,
			"activities": [
				{
					"type": 2,
					"timestamps": {
						"start": %d,
						"end": %d
					},
					"sync_id": "0ydLMgK66rwPhJusUeJpou",
					"state": "%s",
					"session_id": "8175592baa82669363f6ce10ea908c0e",
					"party": {
						"id": "spotify:733933799859028009"
					},
					"name": "Spotify",
					"id": "spotify:1",
					"flags": 48,
					"details": "%s",
					"created_at": 0,
					"assets": {
						"large_text": "%s",
						"large_image": "spotify:%s"
					}
				}
			]
		}
	}`,
		time.Now().Unix()-4140000,
		time.Now().Unix()-3880000,
		json.Escape(utils.Arguments["[\"artist\"]"]),
		json.Escape(utils.Arguments["[\"song\"]"]),
		json.Escape(utils.Arguments["[\"album\"]"]),
		imageID,
	)
	_ = ctx.Session.SendWsData(data)
}

func (s SpotifyCommands) simageid(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		headers := map[string]string{"Authorization": bearerToken}
		searchResponse, err := http.GetJson(fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1))), headers)
		if err != nil {
			_ = ctx.RespondText("A Spotify error occured: No song found")
			return
		}
		imageUrl := searchResponse.GetPath("tracks", "items").GetIndex(0).GetPath("album", "images").GetIndex(0).Get("url").MustString()
		if imageUrl == "" {
			_ = ctx.RespondText("A Spotify error occured: Song not found.")
			return
		}
		splitUrl := strings.Split(imageUrl, "/")
		embed.NewEmbed().
			AddField("Song", searchResponse.GetPath("tracks", "items").GetIndex(0).Get("name").MustString()).
			AddField("Image ID", splitUrl[len(splitUrl)-1]).
			Send(ctx)
	}
}

func (s SpotifyCommands) sshuffle(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		toggle := strings.ToLower(utils.Arguments["<on/off>"])
		if toggle != "on" && toggle != "off" {
			_ = ctx.RespondText("Invalid value. Use `on` or `off`.")
			return
		}
		state := ""
		if toggle == "on" {
			state = "true"
		} else {
			state = "false"
		}
		res, err := http.Put("https://api.spotify.com/v1/me/player/shuffle?state="+state, map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, res, err, fmt.Sprintf("Successfully turned shuffle %s.", toggle))
	}
}

func (s SpotifyCommands) srepeat(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		state := strings.ToLower(utils.Arguments["<track/context/off>"])
		if state != "track" && state != "context" && state != "off" {
			_ = ctx.RespondText("Invalid value. Use `track`, `context` or `off`.")
			return
		}
		res, err := http.Put("https://api.spotify.com/v1/me/player/repeat?state="+state, map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, res, err, fmt.Sprintf("Sucessfully set Spotify repeat mode to %s.", state))
	}
}

func (s SpotifyCommands) sseek(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
		if err != nil {
			_ = ctx.RespondText("Argument <seconds> must be a number.")
			return
		}
		res, err := http.Put(fmt.Sprintf("https://api.spotify.com/v1/me/player/seek?position_ms=%d", seconds*1000), map[string]string{"Authorization": bearerToken})
		s.HandleRequest(ctx, res, err, fmt.Sprintf("Sucessfully seeked song to %d seconds.", seconds))
	}
}

func (s SpotifyCommands) sbearertoken(ctx *dgc.Ctx) {
	if bearerToken := s.GetBearerToken(); bearerToken == "" {
		_ = ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		jsonData, err := http.GetJson("https://api.spotify.com/v1/me", map[string]string{"Authorization": bearerToken})
		if !s.checkRequest(ctx, jsonData, err) {
			return
		}

		embed.NewEmbed().
			AddField("Spotify Account", jsonData.Get("display_name").MustString()).
			AddField("Bearer Token", bearerToken).
			Send(ctx)
	}
}
