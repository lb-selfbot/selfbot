package game

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/go-jester"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/random"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type GameCommands struct{ services *services.Services }

func NewGameCommands(s *services.Services) *GameCommands {
	return &GameCommands{services: s}
}

func (c *GameCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Game",
		Description: "Info about games",
	}, []*router.Command{
		{
			Name:        "mcuuid",
			Description: "Get the minecraft uuid of a player",
			Arguments:   "<username>",
			Handler:     c.mcuuid,
		},
		{
			Name:        "mcnamehistory",
			Description: "Get the minecraft name history of a player",
			Arguments:   "<uuid>",
			Handler:     c.mcnamehistory,
		},
		{
			Name:        "fnmap",
			Description: "Get the current Fortnite map",
			Handler:     c.fnmap,
		},
		{
			Name:        "fnmappois",
			Description: "Get the current Fortnite map with POIs",
			Handler:     c.fnmappois,
		},
		{
			Name:        "fnrandomplaylist",
			Description: "Get a random Fortnite playlist",
			Handler:     c.fnrandomplaylist,
		},
		{
			Name:        "fnrandombanner",
			Description: "Get a random Fortnite banner",
			Handler:     c.fnrandombanner,
		},
		{
			Name:        "fnstats",
			Description: "Get a player's Fortnite stats",
			Arguments:   "<username>",
			Handler:     c.fnstats,
		},
		{
			Name:        "fncreatorcode",
			Description: "Get info about a Fortnite creator code",
			Arguments:   "<code>",
			Handler:     c.fncreatorcode,
		},
		{
			Name:        "fnrandomcosmetic",
			Description: "Get a random Fortnite cosmetic",
			Handler:     c.fnrandomcosmetic,
		},
		{
			Name:        "fnstatus",
			Description: "Get the current Fortnite status",
			Handler:     c.fnstatus,
		},
		{
			Name:        "alstatus",
			Description: "Get the Apex legends status",
			Handler:     c.alstatus,
		},
		{
			Name:        "alnews",
			Description: "Get the latest Apex legends news",
			Handler:     c.alnews,
		},
		{
			Name:        "alstats",
			Description: "Get a player's Apex legends stats",
			Arguments:   "<platform> <username>",
			Handler:     c.alstats,
		},
		{
			Name:        "gtastatus",
			Description: "Get the current GTA status",
			Handler:     c.gtastatus,
		},
		{
			Name:        "fivemplayers",
			Description: "Get the amount of players on a FiveM server",
			Arguments:   "<ip:port>",
			Handler:     c.fivemplayers,
		},
		{
			Name:        "fiveminfo",
			Description: "Get info about a FiveM server",
			Arguments:   "<ip:port>",
			Handler:     c.fiveminfo,
		},
		{
			Name:        "hltbsearch",
			Description: "Search for a game on HowLongToBeat",
			Arguments:   "<game>",
			Handler:     c.hltbsearch,
		},
		{
			Name:        "hltbinfo",
			Description: "Get info about a game on HowLongToBeat",
			Arguments:   "<game id>",
			Handler:     c.hltbinfo,
		},
	})
}

func (c *GameCommands) mcuuid(ctx *router.Ctx) {
	username := ctx.Arguments.Get(0).Raw()

	resp, err := http.DefaultClient().Post(
		"https://api.mojang.com/profiles/minecraft",
		[]string{username},
	)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	uuid := resp.Json().Get(0).Get("id").MustString()

	c.services.Embed.New(ctx).
		AddField("Username", username).
		AddField("UUID", utils.Misc.Or(uuid, "Not Found")).
		Send()
}

func (c *GameCommands) mcnamehistory(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://api.mojang.com/user/profiles/%s/names", ctx.Arguments.Get(0).Escaped()))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	var names []string
	for item := range resp.Json().Iterator() {
		name := item.Get("name").MustString()
		if name == "" {
			continue
		}

		names = append(names, name)
	}

	if len(names) == 0 {
		ctx.RespondText("No name history found (Invalid UUID).")
		return
	}
	c.services.Embed.New(ctx).
		AddField("UUID", ctx.Arguments.Get(0).Raw()).
		AddField("Name History", strings.Join(names, ", ")).
		Send()
}

func (c *GameCommands) fnmap(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://fortnite-api.com/v1/map")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Fortnite Map").
		SetImage(resp.Json().Get("data", "images", "blank").MustString()).
		Send()
}

func (c *GameCommands) fnmappois(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://fortnite-api.com/v1/map")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		SetDescription("Fortnite Map").
		SetImage(resp.Json().Get("data", "images", "pois").MustString()).
		Send()
}

func (c *GameCommands) fnrandomplaylist(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://fortnite-api.com/v1/playlists")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	data := resp.Json().Get("data")
	playlist := data.Get(random.Int(0, data.Len()-1))
	c.services.Embed.New(ctx).
		SetDescription("Random Fortnite Playlist").
		AddField("ID", playlist.Get("id").MustString("None")).
		AddField("Name", playlist.Get("name").MustString("None")).
		AddField("Description", playlist.Get("description").MustString("None")).
		AddField("Game Type", playlist.Get("gameType").MustString("None")).
		AddField("Max Players", fmt.Sprint(playlist.Get("maxPlayers").MustInt(-1))).
		AddField("Added At", playlist.Get("added").MustString("None")).
		Send()
}

func (c *GameCommands) fnrandombanner(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://fortnite-api.com/v1/banners")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	data := resp.Json().Get("data")
	banner := data.Get(random.Int(0, data.Len()-1))
	c.services.Embed.New(ctx).
		SetDescription("Random Fortnite Banner").
		SetImage(banner.Get("images", "icon").MustString()).
		AddField("ID", banner.Get("id").MustString("None")).
		AddField("Developer", banner.Get("devName").MustString("None")).
		AddField("Name", banner.Get("name").MustString("None")).
		AddField("Description", banner.Get("description").MustString("None")).
		AddField("Category", banner.Get("category").MustString("None")).
		Send()
}

func (c *GameCommands) fnstats(ctx *router.Ctx) {
	username := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get(
		fmt.Sprintf("https://fortnite-api.com/v1/stats/br/v2?name=%s", url.QueryEscape(username)),
		http.WithHeader("Authorization", "111a67c8-39d1-4d7c-a5bf-74842b0f112c"),
	)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	data := resp.Json()
	if data.Get("status").MustInt() != 200 {
		ctx.RespondText("Invalid username.")
		return
	}

	user := data.Get("data")
	c.services.Embed.New(ctx).
		SetDescription("Fortnite Stats - %s", user.Get("account", "name").MustString()).
		AddField("Name", user.Get("account", "name").MustString("None")).
		AddField("ID", user.Get("account", "id").MustString("None")).
		AddField("Battlepass Level", user.Get("battlePass", "level").MustString("None")).
		AddField("Battlepass Progress", user.Get("battlePass", "progress").MustString("None")).
		AddField("Image", user.Get("image").MustString("None")).
		AddField("Kills", fmt.Sprint(user.Get("stats", "all", "ltm", "kills").MustInt(-1))).
		AddField("Deaths", fmt.Sprint(user.Get("stats", "all", "ltm", "deaths").MustInt(-1))).
		AddField("Wins", fmt.Sprint(user.Get("stats", "all", "ltm", "wins").MustInt(-1))).
		Send()
}

func (c *GameCommands) fncreatorcode(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get(fmt.Sprintf("https://fortnite-api.com/v2/creatorcode/search?name=%s", ctx.Arguments.Get(0).Escaped()))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	data := resp.Json()
	if data.Get("status").MustInt() != 200 {
		ctx.RespondText("Invalid code.")
		return
	}
	data = data.Get("data")
	c.services.Embed.New(ctx).
		SetDescription("Fortnite Creator Code - %s", data.Get("code").MustString()).
		AddField("Account Name", data.Get("account", "name").MustString("None")).
		AddField("ID", data.Get("account", "id").MustString("None")).
		AddField("Status", data.Get("status").MustString("None")).
		AddField("Verified", utils.Misc.FormatBoolYesNo(data.Get("verified").MustBool())).
		Send()
}

func (c *GameCommands) fnrandomcosmetic(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://fortnite-api.com/v2/cosmetics/br")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	data := resp.Json().Get("data")
	cosmetic := data.Get(random.Int(0, data.Len()-1))
	image := cosmetic.Get("images", "featured").MustString(cosmetic.Get("images", "icon").MustString())
	c.services.Embed.New(ctx).
		SetDescription("Random Fortnite Cosmetic").
		AddField("ID", cosmetic.Get("id").MustString("None")).
		AddField("Name", cosmetic.Get("name").MustString("None")).
		AddField("Description", cosmetic.Get("description").MustString("None")).
		AddField("Type", cosmetic.Get("type", "displayValue").MustString("None")).
		AddField("Rarity", cosmetic.Get("rarity", "displayValue").MustString("None")).
		AddField("Added", cosmetic.Get("added").MustString("None")).
		SetImage(image).
		Send()
}

func (c *GameCommands) fnstatus(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://status.epicgames.com/api/v2/status.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.New(ctx).
		AddField("Fortnite Status", resp.Json().Get("status", "description").MustString("Error")).
		Send()
}

func (c *GameCommands) alstatus(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://api.mozambiquehe.re/servers?auth=IgUq5IRUFvvFxdXbrgmy")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		SetDescription("Apex Legends Status").
		AddField("Origin Login", data.Get("Origin_login", "EU-West", "Status").MustString("None")).
		AddField("EA Nova Fusion", data.Get("EA_novafusion", "EU-West", "Status").MustString("None")).
		AddField("EA Accounts", data.Get("EA_accounts", "EU-West", "Status").MustString("None")).
		AddField("Apex Legends PlayStation", data.Get("otherPlatforms", "Playstation-Network", "Status").MustString("None")).
		AddField("Apex Legends Xbox", data.Get("otherPlatforms", "Xbox-Live", "Status").MustString("None")).
		AddField("Apex Legends Cross Play", data.Get("ApexOauth_Crossplay", "EU-West", "Status").MustString("None")).
		AddField("Apex Legends API", data.Get("selfCoreTest", "Stats-API", "Status").MustString("None")).
		Send()
}

func (c *GameCommands) alnews(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://api.mozambiquehe.re/news?lang=en-us&auth=IgUq5IRUFvvFxdXbrgmy")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := c.services.Embed.New(ctx).SetDescription("Apex Legends News")
	for news := range resp.Json().Iterator() {
		emb.AddField(news.Get("title").MustString("None"), fmt.Sprintf("%s\n%s", news.Get("short_desc").MustString("None"), news.Get("link").MustString("None")))
	}
	emb.Fields = fn.Limit(emb.Fields, 4)
	emb.Send()
}

func (c *GameCommands) alstats(ctx *router.Ctx) {
	platform := ctx.Arguments.Get(0).Raw()
	if platform != "ps4" && platform != "xbox" && platform != "pc" {
		ctx.RespondText("Invalid platform.")
		return
	}
	platform = strings.ReplaceAll(platform, "xbox", "x1")
	resp, err := http.DefaultClient().Get(
		fmt.Sprintf(
			"https://api.mozambiquehe.re/bridge?version=4&platform=%s&player=%s&auth=IgUq5IRUFvvFxdXbrgmy",
			url.QueryEscape(strings.ToUpper(platform)),
			url.QueryEscape(ctx.Arguments.Get(1).Raw()),
		),
	)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	stats := resp.Json()
	if stats.Get("Error").MustString() != "" {
		ctx.RespondText("%s", stats.Get("Error").MustString())
		return
	}

	c.services.Embed.New(ctx).
		SetDescription("Apex Legends Stats - %s", stats.Get("global", "name").MustString()).
		AddField("Name", stats.Get("global", "name").MustString("None")).
		AddField("UID", fmt.Sprint(stats.Get("global", "uid").MustInt64())).
		AddField("Level", fmt.Sprint(stats.Get("global", "level").MustInt())).
		AddField("Rank", stats.Get("global", "rank", "rankName").MustString("None")).
		AddField("Battlepass Level", stats.Get("global", "battlepass", "level").MustString("None")).
		AddField("Selected Legend", stats.Get("legends", "selected", "LegendName").MustString("None")).
		AddField("Kills", fmt.Sprint(stats.Get("total", "kills", "value").MustInt())).
		AddField("Winning Kills", fmt.Sprint(stats.Get("total", "winning_kills", "value").MustInt())).
		AddField("Executions", fmt.Sprint(stats.Get("total", "executions", "value").MustInt())).
		AddField("Headshots", fmt.Sprint(stats.Get("total", "headshots", "value").MustInt())).
		AddField("Bamboozles", fmt.Sprint(stats.Get("total", "bamboozles", "value").MustInt())).
		AddField("Wins as full squad", fmt.Sprint(stats.Get("total", "wins_as_full_squad", "value").MustInt())).
		AddField("Revives", fmt.Sprint(stats.Get("total", "revives", "value").MustInt())).
		AddField("Kill/Death", stats.Get("total", "kd", "value").MustString("None")).
		Send()
}

func (c *GameCommands) gtastatus(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Get("https://support.rockstargames.com/services/status.json")
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	var gtaStatus *jester.Data
	for status := range resp.Json().Get("statuses").Iterator() {
		if status.Get("name").MustString() == "Grand Theft Auto Online" {
			gtaStatus = status
			break
		}
	}

	emb := c.services.Embed.New(ctx).
		SetDescription("GTA Online Status").
		AddField("Overall Status", gtaStatus.Get("status_tag").MustString())
	for platform := range gtaStatus.Get("services_platforms").Iterator() {
		emb.AddField(platform.Get("name").MustString(), platform.Get("status").MustString())
	}
	emb.Send()
}

func (c *GameCommands) fivemplayers(ctx *router.Ctx) {
	host := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get("http://" + host + "/dynamic.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	json := resp.Json()

	c.services.Embed.New(ctx).
		AddField("Host", host).
		AddField("Players", fmt.Sprint(json.Get("clients").MustInt())).
		AddField("Max Players", json.Get("sv_maxclients").MustString()).
		Send()
}

func (c *GameCommands) fiveminfo(ctx *router.Ctx) {
	host := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get("http://" + host + "/dynamic.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	json := resp.Json()

	c.services.Embed.New(ctx).
		AddField("Host", host).
		AddField("Server Name", json.Get("hostname").MustString()).
		AddField("Game Type", json.Get("gametype").MustString()).
		AddField("Map", json.Get("mapname").MustString()).
		AddField("Players", fmt.Sprint(json.Get("clients").MustInt())).
		AddField("Max Players", json.Get("sv_maxclients").MustString()).
		Send()
}

var HLTBHeaders = http.WithHeaders(map[string]string{
	"Authority":          "howlongtobeat.com",
	"Method":             "POST",
	"Path":               "/api/search",
	"Scheme":             "https",
	"Accept":             "application/json, text/plain, */*",
	"Accept-Encoding":    "gzip, deflate, br",
	"Accept-Language":    "en-US,en;q=0.9",
	"Origin":             "https://howlongtobeat.com",
	"Referer":            "https://howlongtobeat.com/",
	"Sec-CH-UA":          "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"119\", \"Google Chrome\";v=\"119\"",
	"Sec-CH-UA-Mobile":   "?0",
	"Sec-CH-UA-Platform": "Windows",
	"Sec-Fetch-Dest":     "empty",
	"Sec-Fetch-Mode":     "cors",
	"Sec-Fetch-Site":     "same-origin",
	"User-Agent":         discordgo.UserAgentWeb,
})

func (c *GameCommands) hltbsearch(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Post(
		"https://howlongtobeat.com/api/search",
		map[string]any{
			"searchType":  "games",
			"searchTerms": []string{ctx.Arguments.Raw()},
			"searchPage":  1,
			"size":        20,
			"searchOptions": map[string]any{
				"games": map[string]any{
					"userId":        0,
					"platform":      "",
					"sortCategory":  "popular",
					"rangeCategory": "main",
					"rangeTime":     map[string]any{"min": nil, "max": nil},
					"gameplay":      map[string]any{"perspective": "", "flow": "", "genre": ""},
					"rangeYear":     map[string]any{"min": "", "max": ""},
					"modifier":      "",
				},
				"users":      map[string]any{"sortCategory": "postcount"},
				"lists":      map[string]any{"sortCategory": "follows"},
				"filter":     "",
				"sort":       0,
				"randomizer": 0,
			},
		},
		HLTBHeaders,
	)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	if resp.Json().Get("count").MustInt() == 0 {
		ctx.RespondText("No results found.")
		return
	}

	var games []string

	for game := range resp.Json().Get("data").Iterator() {
		games = append(games, fmt.Sprintf(
			"%s (%d)",
			game.Get("game_name").MustString(),
			game.Get("game_id").MustInt(),
		))
	}

	c.services.Embed.New(ctx).
		AddField("Games", strings.Join(games, "\n")).
		Send()
}

func secondsToTime(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds - (hours * 3600)) / 60
	return fmt.Sprintf("%d Hours, %d Minutes", hours, minutes)
}

func (c *GameCommands) hltbinfo(ctx *router.Ctx) {
	resp, err := http.DefaultClient().Post(
		"https://howlongtobeat.com/game?id="+ctx.Arguments.Get(0).Raw(),
		nil,
		HLTBHeaders,
	)
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	dataMatches := regexp.MustCompile("<script id=\"__NEXT_DATA__\" type=\"application/json\">(.*?)</script>").FindSubmatch(resp.MustGetBody())
	if len(dataMatches) < 2 {
		ctx.RespondText("An error has occurred: No data matches")
		return
	}

	data, err := jester.NewJson(dataMatches[1])
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	games := data.Get("props", "pageProps", "game", "data", "game")
	if games.Len() == 0 {
		ctx.RespondText("An error has occurred: Game not found")
		return
	}

	game := games.Get(0)

	c.services.Embed.New(ctx).
		AddField("Name", game.Get("game_name").MustString("None")).
		AddField("Aliases", game.Get("game_alias").MustString("None")).
		AddField("Main Story", secondsToTime(game.Get("comp_main").MustInt())).
		AddField("Main + Extras", secondsToTime(game.Get("comp_plus").MustInt())).
		AddField("Completionist", secondsToTime(game.Get("comp_100").MustInt())).
		AddField("All Styles", secondsToTime(game.Get("comp_all").MustInt())).
		AddField("Rating", fmt.Sprintf("%d/100", game.Get("review_score").MustInt())).
		AddField("Currently Playing", fmt.Sprintf("%d Players", game.Get("count_playing").MustInt())).
		AddField("Developer", game.Get("profile_dev").MustString("None")).
		AddField("Platforms", game.Get("profile_platform").MustString("None")).
		Send()
}
