package game

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/lb-selfbot/discordgo"
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

	resp, err := http.PostJson(
		"https://api.mojang.com/profiles/minecraft",
		[]string{username},
	)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	uuid := resp.GetIndex(0).Get("id").MustString()

	ctx.NewEmbed().
		AddField("Username", username).
		AddField("UUID", utils.Misc.Or(uuid, "Not Found")).
		Send()
}

func (c *GameCommands) mcnamehistory(ctx *router.Ctx) {
	uuid := ctx.Arguments.Get(0).Raw()
	res, err := http.GetJson(fmt.Sprintf("https://api.mojang.com/user/profiles/%s/names", url.QueryEscape(uuid)))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	nameLen := len(res.MustArray())
	var names []string
	for i := 0; i < nameLen; i++ {
		names = append(names, res.GetIndex(i).Get("name").MustString())
	}
	if len(names) == 0 {
		names = append(names, "Invalid UUID.")
	}
	ctx.NewEmbed().
		AddField("UUID", uuid).
		AddField("Name History", strings.Join(names, ", ")).
		Send()
}

func (c *GameCommands) fnmap(ctx *router.Ctx) {
	res, err := http.GetJson("https://fortnite-api.com/v1/map")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.NewEmbed().
		SetDescription("Fortnite Map").
		SetImage(res.GetPath("data", "images", "blank").MustString()).
		Send()
}

func (c *GameCommands) fnmappois(ctx *router.Ctx) {
	res, err := http.GetJson("https://fortnite-api.com/v1/map")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.NewEmbed().
		SetDescription("Fortnite Map").
		SetImage(res.GetPath("data", "images", "pois").MustString()).
		Send()
}

func (c *GameCommands) fnrandomplaylist(ctx *router.Ctx) {
	playlists, err := http.GetJson("https://fortnite-api.com/v1/playlists")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	playlist := playlists.Get("data").GetIndex(random.Int(0, len(playlists.Get("data").MustArray())-1))
	ctx.NewEmbed().
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
	banners, err := http.GetJson("https://fortnite-api.com/v1/banners")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	banner := banners.Get("data").GetIndex(random.Int(0, len(banners.Get("data").MustArray())-1))
	ctx.NewEmbed().
		SetDescription("Random Fortnite Banner").
		SetImage(banner.GetPath("images", "icon").MustString()).
		AddField("ID", banner.Get("id").MustString("None")).
		AddField("Developer", banner.Get("devName").MustString("None")).
		AddField("Name", banner.Get("name").MustString("None")).
		AddField("Description", banner.Get("description").MustString("None")).
		AddField("Category", banner.Get("category").MustString("None")).
		Send()
}

func (c *GameCommands) fnstats(ctx *router.Ctx) {
	username := ctx.Arguments.Get(0).Raw()
	user, err := http.GetJson(
		fmt.Sprintf("https://fortnite-api.com/v1/stats/br/v2?name=%s", url.QueryEscape(username)),
		map[string]string{
			"Authorization": "111a67c8-39d1-4d7c-a5bf-74842b0f112c",
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if user.Get("status").MustInt(-1) != 200 {
		ctx.RespondText("Invalid username.")
		return
	}
	user = user.GetPath("data")
	ctx.NewEmbed().
		SetDescription("Fortnite Stats - %s", user.GetPath("account", "name").MustString()).
		AddField("Name", user.GetPath("account", "name").MustString("None")).
		AddField("ID", user.GetPath("account", "id").MustString("None")).
		AddField("Battlepass Level", user.GetPath("battlePass", "level").MustString("None")).
		AddField("Battlepass Progress", user.GetPath("battlePass", "progress").MustString("None")).
		AddField("Image", user.GetPath("image").MustString("None")).
		AddField("Kills", fmt.Sprint(user.GetPath("stats", "all", "ltm", "kills").MustInt(-1))).
		AddField("Deaths", fmt.Sprint(user.GetPath("stats", "all", "ltm", "deaths").MustInt(-1))).
		AddField("Wins", fmt.Sprint(user.GetPath("stats", "all", "ltm", "wins").MustInt(-1))).
		Send()
}

func (c *GameCommands) fncreatorcode(ctx *router.Ctx) {
	code := ctx.Arguments.Get(0).Raw()
	res, err := http.GetJson(fmt.Sprintf("https://fortnite-api.com/v2/creatorcode/search?name=%s", url.QueryEscape(code)))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if res.Get("status").MustInt(-1) != 200 {
		ctx.RespondText("Invalid code.")
		return
	}
	res = res.GetPath("data")
	ctx.NewEmbed().
		SetDescription("Fortnite Creator Code - %s", res.GetPath("code").MustString()).
		AddField("Account Name", res.GetPath("account", "name").MustString("None")).
		AddField("ID", res.GetPath("account", "id").MustString("None")).
		AddField("Status", res.GetPath("status").MustString("None")).
		AddField("Verified", utils.Misc.FormatBoolYesNo(res.GetPath("verified").MustBool())).
		Send()
}

func (c *GameCommands) fnrandomcosmetic(ctx *router.Ctx) {
	cosmetics, err := http.GetJson("https://fortnite-api.com/v2/cosmetics/br")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	cosmetic := cosmetics.Get("data").GetIndex(random.Int(0, len(cosmetics.Get("data").MustArray())-1))
	image := cosmetic.GetPath("images", "featured").MustString()
	if image == "" {
		image = cosmetic.GetPath("images", "icon").MustString()
	}
	ctx.NewEmbed().
		SetDescription("Random Fortnite Cosmetic").
		AddField("ID", cosmetic.Get("id").MustString("None")).
		AddField("Name", cosmetic.Get("name").MustString("None")).
		AddField("Description", cosmetic.Get("description").MustString("None")).
		AddField("Type", cosmetic.GetPath("type", "displayValue").MustString("None")).
		AddField("Rarity", cosmetic.GetPath("rarity", "displayValue").MustString("None")).
		AddField("Added", cosmetic.GetPath("added").MustString("None")).
		SetImage(image).
		Send()
}

func (c *GameCommands) fnstatus(ctx *router.Ctx) {
	httpResponse := http.Get("https://status.epicgames.com/api/v2/status.json")
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	json := httpResponse.Json()

	ctx.NewEmbed().
		AddField("Fortnite Status", json.GetPath("status", "description").MustString("Error")).
		Send()
}

func (c *GameCommands) alstatus(ctx *router.Ctx) {
	status, err := http.GetJson("https://api.mozambiquehe.re/servers?auth=IgUq5IRUFvvFxdXbrgmy")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.NewEmbed().
		SetDescription("Apex Legends Status").
		AddField("Origin Login", status.GetPath("Origin_login", "EU-West", "Status").MustString("None")).
		AddField("EA Nova Fusion", status.GetPath("EA_novafusion", "EU-West", "Status").MustString("None")).
		AddField("EA Accounts", status.GetPath("EA_accounts", "EU-West", "Status").MustString("None")).
		AddField("Apex Legends PlayStation", status.GetPath("otherPlatforms", "Playstation-Network", "Status").MustString("None")).
		AddField("Apex Legends Xbox", status.GetPath("otherPlatforms", "Xbox-Live", "Status").MustString("None")).
		AddField("Apex Legends Cross Play", status.GetPath("ApexOauth_Crossplay", "EU-West", "Status").MustString("None")).
		AddField("Apex Legends API", status.GetPath("selfCoreTest", "Stats-API", "Status").MustString("None")).
		Send()
}

func (c *GameCommands) alnews(ctx *router.Ctx) {
	newsList, err := http.GetJson("https://api.mozambiquehe.re/news?lang=en-us&auth=IgUq5IRUFvvFxdXbrgmy")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := ctx.NewEmbed().SetDescription("Apex Legends News")
	for i := 0; i < len(newsList.MustArray()); i++ {
		news := newsList.GetIndex(i)
		emb.AddField(news.GetPath("title").MustString("None"), fmt.Sprintf("%s\n%s", news.GetPath("short_desc").MustString("None"), news.GetPath("link").MustString("None")))
		if i > 4 {
			break
		}
	}
	emb.Send()
}

func (c *GameCommands) alstats(ctx *router.Ctx) {
	platform := ctx.Arguments.Get(0).Raw()
	if platform != "ps4" && platform != "xbox" && platform != "pc" {
		ctx.RespondText("Invalid platform.")
		return
	}
	platform = strings.ReplaceAll(platform, "xbox", "x1")
	stats, err := http.GetJson(
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
	if stats.Get("Error").MustString() != "" {
		ctx.RespondText(stats.Get("Error").MustString())
		return
	}
	ctx.NewEmbed().
		SetDescription("Apex Legends Stats - %s", stats.GetPath("global", "name").MustString()).
		AddField("Name", stats.GetPath("global", "name").MustString("None")).
		AddField("UID", fmt.Sprint(stats.GetPath("global", "uid").MustInt64())).
		AddField("Level", fmt.Sprint(stats.GetPath("global", "level").MustInt())).
		AddField("Rank", stats.GetPath("global", "rank", "rankName").MustString("None")).
		AddField("Battlepass Level", stats.GetPath("global", "battlepass", "level").MustString("None")).
		AddField("Selected Legend", stats.GetPath("legends", "selected", "LegendName").MustString("None")).
		AddField("Kills", fmt.Sprint(stats.GetPath("total", "kills", "value").MustInt())).
		AddField("Winning Kills", fmt.Sprint(stats.GetPath("total", "winning_kills", "value").MustInt())).
		AddField("Executions", fmt.Sprint(stats.GetPath("total", "executions", "value").MustInt())).
		AddField("Headshots", fmt.Sprint(stats.GetPath("total", "headshots", "value").MustInt())).
		AddField("Bamboozles", fmt.Sprint(stats.GetPath("total", "bamboozles", "value").MustInt())).
		AddField("Wins as full squad", fmt.Sprint(stats.GetPath("total", "wins_as_full_squad", "value").MustInt())).
		AddField("Revives", fmt.Sprint(stats.GetPath("total", "revives", "value").MustInt())).
		AddField("Kill/Death", stats.GetPath("total", "kd", "value").MustString("None")).
		Send()
}

func (c *GameCommands) gtastatus(ctx *router.Ctx) {
	status, err := http.GetJson("https://support.rockstargames.com/services/status.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	status = status.Get("statuses")
	for i := 0; i < len(status.MustArray()); i++ {
		if status.GetIndex(i).Get("name").MustString() == "Grand Theft Auto Online" {
			status = status.GetIndex(i)
			break
		}
	}
	emb := ctx.NewEmbed().
		SetDescription("GTA Online Status").
		AddField("Overall Status", status.Get("status_tag").MustString())
	for i := 0; i < len(status.Get("services_platforms").MustArray()); i++ {
		platform := status.Get("services_platforms").GetIndex(i)
		emb.AddField(platform.Get("name").MustString(), platform.Get("service_status").Get("status").MustString())
	}
	emb.Send()
}

func (c *GameCommands) fivemplayers(ctx *router.Ctx) {
	host := ctx.Arguments.Get(0).Raw()
	httpResponse := http.Get("http://" + host + "/dynamic.json")
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	json := httpResponse.Json()

	ctx.NewEmbed().
		AddField("Host", host).
		AddField("Players", fmt.Sprint(json.Get("clients").MustInt())).
		AddField("Max Players", json.Get("sv_maxclients").MustString()).
		Send()
}

func (c *GameCommands) fiveminfo(ctx *router.Ctx) {
	host := ctx.Arguments.Get(0).Raw()
	httpResponse := http.Get("http://" + host + "/dynamic.json")
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred.")
		return
	}

	json := httpResponse.Json()

	ctx.NewEmbed().
		AddField("Host", host).
		AddField("Server Name", json.Get("hostname").MustString()).
		AddField("Game Type", json.Get("gametype").MustString()).
		AddField("Map", json.Get("mapname").MustString()).
		AddField("Players", fmt.Sprint(json.Get("clients").MustInt())).
		AddField("Max Players", json.Get("sv_maxclients").MustString()).
		Send()
}

var HLTBHeaders = map[string]string{
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
}

func (c *GameCommands) hltbsearch(ctx *router.Ctx) {
	resp, err := http.PostJson(
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

	if resp.Get("count").MustInt() == 0 {
		ctx.RespondText("No results found.")
		return
	}

	var games []string

	for i := 0; i < len(resp.Get("data").MustArray()); i++ {
		game := resp.Get("data").GetIndex(i)
		games = append(games, fmt.Sprintf(
			"%s (%d)",
			game.Get("game_name").MustString(),
			game.Get("game_id").MustInt(),
		))
	}

	ctx.NewEmbed().
		AddField("Games", strings.Join(games, "\n")).
		Send()
}

func secondsToTime(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds - (hours * 3600)) / 60
	return fmt.Sprintf("%d Hours, %d Minutes", hours, minutes)
}

func (c *GameCommands) hltbinfo(ctx *router.Ctx) {
	httpResponse := http.Post(
		"https://howlongtobeat.com/game?id="+ctx.Arguments.Get(0).Raw(),
		nil,
		HLTBHeaders,
	)
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred: %v", httpResponse.Error)
		return
	}

	dataMatches := regexp.MustCompile("<script id=\"__NEXT_DATA__\" type=\"application/json\">(.*?)</script>").FindStringSubmatch(string(httpResponse.Body))
	if len(dataMatches) < 2 {
		ctx.RespondText("An error has occurred: No data matches")
		return
	}

	json, err := simplejson.NewJson([]byte(dataMatches[1]))
	if err != nil {
		ctx.RespondText("An error has occurred: %v", err)
		return
	}

	games := json.GetPath("props", "pageProps", "game", "data", "game")
	if len(games.MustArray()) == 0 {
		ctx.RespondText("An error has occurred: Game not found")
		return
	}

	game := games.GetIndex(0)

	ctx.NewEmbed().
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
