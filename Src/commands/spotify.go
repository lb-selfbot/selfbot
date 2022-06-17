package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev/LB-Selfbot/config"
	"github.com/LightningDev/LB-Selfbot/embed"
	"github.com/LightningDev/LB-Selfbot/http"
	"github.com/LightningDev/LB-Selfbot/json"
	"github.com/LightningDev/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SpotifyCommands struct{}

func (t SpotifyCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "spause",
		Description: "Pause Spotify Music",
		Usage:       "[p]spause",
		Handler:     t.spause,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sresume",
		Description: "Resume Spotify Music",
		Usage:       "[p]sresume",
		Handler:     t.sresume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sprevious",
		Description: "Play the previous song",
		Usage:       "[p]sprevious",
		Handler:     t.sprevious,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snext",
		Description: "Play the next song",
		Usage:       "[p]snext",
		Handler:     t.snext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "splay",
		Description: "Play a song on Spotify",
		Usage:       "[p]splay <song name>",
		Handler:     t.splay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "splaying",
		Description: "Current playing song info",
		Usage:       "[p]splaying",
		Handler:     t.splaying,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "svolume",
		Description: "Change Spotify volume",
		Usage:       "[p]svolume <volume>",
		Handler:     t.svolume,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "spresence",
		Description: "Set a fake Spotify presence (use simageid command for image id)",
		Usage:       "[p]spresence <on/off> [\"artist\"] [\"song\"] [\"album\"] [image id]",
		Handler:     t.spresence,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "simageid",
		Description: "Get an image id for the spresence command",
		Usage:       "[p]simageid <song name>",
		Handler:     t.simageid,
	})
}

func (t SpotifyCommands) GetBearerToken() string {
	cfg, err := config.GetConfig()
	if err != nil {
		return ""
	}
	discordHeaders := http.GetDiscordHeaders(cfg.Account.Discord.Token)
	connections, err := http.GetSimpleJson("https://discord.com/api/v8/users/@me/connections", discordHeaders)
	if err != nil {
		return ""
	}
	for _, connectionInterface := range connections.MustArray() {
		connection := connectionInterface.(map[string]interface{})
		if connType := connection["type"]; connType != nil && connType.(string) == "spotify" {
			bearerToken := fmt.Sprintf("Bearer %s", connection["access_token"].(string))
			res, err := http.GetSimpleJson("https://api.spotify.com/v1/search?q=Despacito&type=track", map[string]string{"Authorization": bearerToken})
			if err != nil {
				return ""
			}
			if _, success := res.GetPath("tracks", "items").GetIndex(0).Get("external_urls").CheckGet("spotify"); success {
				return bearerToken
			}
			res, err = http.GetSimpleJson(fmt.Sprintf("https://discord.com/api/v8/users/@me/connections/spotify/%s/access-token", connection["name"].(string)), discordHeaders)
			if err != nil {
				return ""
			}
			return "Bearer " + res.Get("access_token").MustString()
		}
	}
	return ""
}

func (t SpotifyCommands) HandleRequest(ctx *dgc.Ctx, response string, err error, onSuccess string) {
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if response != "" {
		responseJson, err := json.SimpleUnmarshal(response)
		if err != nil {
			ctx.RespondText("An error has occurred.")
			return
		}
		errorReason := responseJson.GetPath("error", "reason").MustString()
		if errorReason == "NO_ACTIVE_DEVICE" {
			ctx.RespondText("A Spotify error occured: `No active device found`")
			return
		}
	}
	ctx.RespondText(onSuccess)
}

func (t SpotifyCommands) spause(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Put("https://api.spotify.com/v1/me/player/pause", map[string]string{"Authorization": bearerToken})
		t.HandleRequest(ctx, response, err, "Successfully paused music.")
	}
}

func (t SpotifyCommands) sresume(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Put("https://api.spotify.com/v1/me/player/play", map[string]string{"Authorization": bearerToken})
		t.HandleRequest(ctx, response, err, "Successfully resumed music.")
	}
}

func (t SpotifyCommands) sprevious(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Post("https://api.spotify.com/v1/me/player/previous", "", map[string]string{"Authorization": bearerToken})
		t.HandleRequest(ctx, response, err, "Successfully playing previous song.")
	}
}

func (t SpotifyCommands) snext(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.Post("https://api.spotify.com/v1/me/player/next", "", map[string]string{"Authorization": bearerToken})
		t.HandleRequest(ctx, response, err, "Successfully playing next song.")
	}
}

func (t SpotifyCommands) splay(ctx *dgc.Ctx) {
	songName := ctx.Arguments.Raw()
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		headers := map[string]string{"Authorization": bearerToken}
		searchResponse, err := http.GetSimpleJson(fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(songName)), headers)
		if err != nil {
			ctx.RespondText("An error has occurred.")
			return
		}
		search := searchResponse.GetPath("tracks", "items").MustArray()
		if len(search) == 0 {
			ctx.RespondText("No songs found.")
			return
		}
		song := search[0].(map[string]interface{})
		songUrl := song["external_urls"].(map[string]interface{})["spotify"].(string)
		songName := song["name"].(string)
		songArtists := []string{}
		for _, artist := range song["artists"].([]interface{}) {
			songArtists = append(songArtists, artist.(map[string]interface{})["name"].(string))
		}
		http.Post(fmt.Sprintf("https://api.spotify.com/v1/me/player/queue?uri=%s", url.QueryEscape(songUrl)), "", headers)
		response, err := http.Post("https://api.spotify.com/v1/me/player/next", "", headers)
		t.HandleRequest(ctx, response, err, fmt.Sprintf("Successfully playing %s - %s.", strings.Join(songArtists, ", "), songName))
	}
}

func (t SpotifyCommands) splaying(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		response, err := http.GetSimpleJson("https://api.spotify.com/v1/me/player/currently-playing", map[string]string{"Authorization": bearerToken})
		song, success := response.CheckGet("item")
		if err != nil || !success {
			ctx.RespondText("A Spotify error occured: No song currently playing")
			return
		}
		songArtists := []string{}
		for _, artist := range song.Get("artists").MustArray() {
			songArtists = append(songArtists, artist.(map[string]interface{})["name"].(string))
		}
		embed.NewEmbed().
			AddField("Title", song.Get("name").MustString()).
			AddField("Artist", strings.Join(songArtists, ", ")).
			AddField("Release Date", song.GetPath("album", "release_date").MustString("Unknown")).
			AddField("Explicit", utils.FormatBool(song.Get("explicit").MustBool(false))).
			AddField("Duration", fmt.Sprintf("%d:%02d", int(song.Get("duration_ms").MustInt(0))/1000/60, int(song.Get("duration_ms").MustInt(0))/1000%60)).
			AddField("URL", song.GetPath("external_urls", "spotify").MustString("Unknown")).
			SetImage(song.GetPath("album", "images").GetIndex(0).Get("url").MustString("")).
			Send(ctx)
	}
}

func (t SpotifyCommands) svolume(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		volumeStr := utils.Arguments["<volume>"]
		volume, err := strconv.Atoi(volumeStr)
		if err != nil || volume < 0 || volume > 100 {
			ctx.RespondText("Volume must be between 0 and 100.")
			return
		}
		res, err := http.Put("https://api.spotify.com/v1/me/player/volume?volume_percent="+volumeStr, map[string]string{"Authorization": bearerToken})
		t.HandleRequest(ctx, res, err, fmt.Sprintf("Successfully set volume to: %s.", volumeStr))
	}
}

func (t SpotifyCommands) spresence(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"]) == "on"
	if !toggle {
		ctx.Session.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{},
		})
		return
	}
	imageID := utils.JsonEscape(utils.Arguments["[\"image_id\"]"])
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
		utils.JsonEscape(utils.Arguments["[\"artist\"]"]),
		utils.JsonEscape(utils.Arguments["[\"song\"]"]),
		utils.JsonEscape(utils.Arguments["[\"album\"]"]),
		imageID,
	)
	ctx.Session.SendWsData(data)
}

func (t SpotifyCommands) simageid(ctx *dgc.Ctx) {
	if bearerToken := t.GetBearerToken(); bearerToken == "" {
		ctx.RespondText("Could not find a Spotify connection. Please connect your Spotify account to your Discord.")
	} else {
		headers := map[string]string{"Authorization": bearerToken}
		searchResponse, err := http.GetSimpleJson(fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=track", url.QueryEscape(utils.Arguments["<song name>"])), headers)
		if err != nil {
			ctx.RespondText("A Spotify error occured: No song found")
			return
		}
		url := searchResponse.GetPath("tracks", "items").GetIndex(0).GetPath("album", "images").GetIndex(0).Get("url").MustString()
		if url == "" {
			ctx.RespondText("A Spotify error occured: Song not found.")
			return
		}
		splitUrl := strings.Split(url, "/")
		embed.NewEmbed().
			AddField("Song", searchResponse.GetPath("tracks", "items").GetIndex(0).Get("name").MustString()).
			AddField("Image ID", splitUrl[len(splitUrl)-1]).
			Send(ctx)
	}
}
