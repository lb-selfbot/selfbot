package commands

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/error_handling"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/json"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/kkdai/youtube/v2"
	googlesearch "github.com/rocketlaunchr/google-search"
)

type NetworkCommands struct{}

func (n NetworkCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "ping",
		Description: "Check the latency of the bot",
		Usage:       "[p]ping",
		Handler:     n.ping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pingweb",
		Description: "Check the latency to a website",
		Usage:       "[p]pingweb <host>",
		Handler:     n.pingweb,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "googlesearch",
		Description: "Search someting on Google",
		Usage:       "[p]googlesearch <query>",
		Handler:     n.googlesearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "youtubesearch",
		Description: "Search someting on Youtube",
		Usage:       "[p]youtubesearch <query>",
		Handler:     n.youtubesearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "urbandictionary",
		Description: "Search someting on the Urban Dictionary",
		Usage:       "[p]urbandictionary <word>",
		Handler:     n.urbandictionary,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wikipedia",
		Description: "Search someting on Wikipedia",
		Usage:       "[p]wikipedia <query>",
		Handler:     n.wikipedia,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "weather",
		Description: "Check the weather in a city",
		Usage:       "[p]weather <city>",
		Handler:     n.weather,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hypesquad",
		Description: "Change your Hypesquad House",
		Usage:       "[p]hypesquad <bravery/brilliance/balance>",
		Handler:     n.hypesquad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tinyurl",
		Description: "Shorten a URL with TinyURL",
		Usage:       "[p]tinyurl <url>",
		Handler:     n.tinyurl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ipinfo",
		Description: "Get info about an IP Address",
		Usage:       "[p]ipinfo <ip>",
		Handler:     n.ipinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "youtubedl",
		Description: "Download a YouTube video to the desktop",
		Usage:       "[p]youtubedl <query>",
		Handler:     n.youtubedl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lvbypass",
		Description: "Bypass a Linkvertise link",
		Usage:       "[p]lvbypass <url>",
		Handler:     n.lvbypass,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "portscan",
		Description: "Scan a host for open ports",
		Usage:       "[p]portscan <host>",
		Handler:     n.portscan,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "theme",
		Description: "Change your Discord theme",
		Usage:       "[p]theme <light/dark>",
		Handler:     n.theme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hastebin",
		Description: "Upload text to Hastebin",
		Usage:       "[p]hastebin <text>",
		Handler:     n.hastebin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "checktoken",
		Description: "Check if a Discord token is valid",
		Usage:       "[p]checktoken <token>",
		Handler:     n.checktoken,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delwebhook",
		Description: "Delete a webook",
		Usage:       "[p]delwebhook <webhook url>",
		Handler:     n.delwebhook,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "compact",
		Description: "Toggle Discord compact mode",
		Usage:       "[p]compact <on/off>",
		Handler:     n.compact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "devmode",
		Description: "Toggle Discord developer mode",
		Usage:       "[p]devmode <on/off>",
		Handler:     n.devmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "resolve",
		Description: "Get a domain's IP address",
		Usage:       "[p]resolve <domain>",
		Handler:     n.resolve,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reversedns",
		Description: "Find all domains hosted on an IP address",
		Usage:       "[p]reversedns <ip>",
		Handler:     n.reversedns,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reversemac",
		Description: "Find the vender of a MAC address",
		Usage:       "[p]reversemac <mac>",
		Handler:     n.reversemac,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amazon",
		Description: "Search products on Amazon",
		Usage:       "[p]amazon <query>",
		Handler:     n.amazon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "steam",
		Description: "Search games on Steam",
		Usage:       "[p]steam <query>",
		Handler:     n.steam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "github",
		Description: "Search repos on GitHub",
		Usage:       "[p]github <query>",
		Handler:     n.github,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stackoverflow",
		Description: "Search answers on Stack Overflow",
		Usage:       "[p]stackoverflow <query>",
		Handler:     n.stackoverflow,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "covid",
		Description: "See COVID-19 statistics for a country",
		Usage:       "[p]covid <country>",
		Handler:     n.covid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discordstatus",
		Description: "Check the Discord status",
		Usage:       "[p]discordstatus",
		Handler:     n.discordstatus,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "open",
		Description: "Open a URL in the browser",
		Usage:       "[p]open <url>",
		Handler:     n.open,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "searchcountry",
		Description: "Get info about a country",
		Usage:       "[p]searchcountry <country>",
		Handler:     n.searchcountry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomcountry",
		Description: "Get info about a random country",
		Usage:       "[p]randomcountry <country>",
		Handler:     n.randomcountry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "headers",
		Description: "Get URL response headers",
		Usage:       "[p]headers <url>",
		Handler:     n.headers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "interface",
		Description: "See your public/local ip and MAC address",
		Usage:       "[p]interface",
		Handler:     n._interface,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "revimage",
		Description: "Reverse image search",
		Usage:       "[p]revimage <image url>",
		Handler:     n.revimage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "paste",
		Description: "Create a paste",
		Usage:       "[p]pastebin <text>",
		Handler:     n.pastebin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pastefile",
		Description: "Create a paste from a file",
		Usage:       "[p]pastefile <filepath>",
		Handler:     n.pastefile,
	})
}

func (NetworkCommands) ping(ctx *dgc.Ctx) {
	emb := embed.NewEmbed()
	time, err := api.Ping("lightning-bot.com")
	if err != nil {
		emb.AddField("Ping to LightningBot", "Error")
	} else {
		emb.AddField("Ping to LightningBot", fmt.Sprintf("%dms", time.Milliseconds()))
	}
	time, err = api.Ping("google.com")
	if err != nil {
		emb.AddField("Ping to Google", "Error")
	} else {
		emb.AddField("Ping to Google", fmt.Sprintf("%dms", time.Milliseconds()))
	}
	emb.
		AddField("Ping to Discord", fmt.Sprintf("%dms", ctx.Session.HeartbeatLatency().Milliseconds())).
		Send(ctx)
}

func (NetworkCommands) pingweb(ctx *dgc.Ctx) {
	host := utils.Arguments["<host>"]
	time, err := api.Ping(host)
	emb := embed.NewEmbed()
	if err != nil {
		emb.AddField("Ping to "+host, "Error")
	} else {
		emb.AddField("Ping to "+host, fmt.Sprintf("%dms", time.Milliseconds()))
	}
	emb.Send(ctx)
}

func (NetworkCommands) googlesearch(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	results, err := googlesearch.Search(context.TODO(), query)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := embed.NewEmbed()
	for index, result := range results {
		emb.AddField(result.Title, result.URL)
		if index == 5 {
			break
		}
	}
	emb.Send(ctx)
}

func (NetworkCommands) youtubesearch(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	results, err := googlesearch.Search(context.TODO(), query+" site:youtube.com")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := embed.NewEmbed()
	for index, result := range results {
		if strings.Contains(result.URL, "playlist") {
			index -= 1
			continue
		}
		emb.AddField(result.Title, result.URL)
		if index == 5 {
			break
		}
	}
	emb.Send(ctx)
}

func (NetworkCommands) urbandictionary(ctx *dgc.Ctx) {
	word := ctx.Arguments.Raw()
	results, err := http.GetJson(fmt.Sprintf("https://api.urbandictionary.com/v0/define?term=%s", url.QueryEscape(word)))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	list := results.Get("list")
	if len(list.MustArray()) == 0 {
		ctx.RespondText("No results found.")
		return
	}
	result := list.GetIndex(0)
	embed.NewEmbed().
		AddField("Word", word).
		AddField("Definition", result.Get("definition").MustString()).
		AddField("Permalink", result.Get("permalink").MustString()).
		AddField("Thumbs Up", fmt.Sprint(result.Get("thumbs_up").MustFloat64())).
		AddField("Thumbs Down", fmt.Sprint(result.Get("thumbs_down").MustFloat64())).
		AddField("Author", result.Get("author").MustString()).
		Send(ctx)
}

func (NetworkCommands) wikipedia(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	article, err := api.WikipediaSearch(query)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Query", query).
		AddField("Wikipedia Summary", article[:1500]).
		Send(ctx)
}

func (NetworkCommands) weather(ctx *dgc.Ctx) {
	city := ctx.Arguments.Raw()
	weather, err := http.GetJson(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?appid=533d00f29739edba54c5fb2b1a8262b3&q=%s", city))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if fmt.Sprint(weather.Get("cod").MustInt()) == "404" {
		ctx.RespondText("No results found.")
		return
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Weather in %s, %s", weather.Get("name").MustString(), weather.GetPath("sys", "country").MustString())).
		AddField("Description", utils.String.Capitalize(weather.Get("weather").GetIndex(0).Get("description").MustString())).
		AddField("Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "temp").MustFloat64()-273.15))).
		AddField("Feels Like", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "feels_like").MustFloat64()-273.15))).
		AddField("Min. Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "temp_min").MustFloat64()-273.15))).
		AddField("Max. Temperature", fmt.Sprintf("%.2f °C", math.Round(weather.GetPath("main", "temp_max").MustFloat64()-273.15))).
		AddField("Pressure", fmt.Sprintf("%.2f hPa", weather.GetPath("main", "pressure").MustFloat64())).
		AddField("Humidity", fmt.Sprintf("%.2f %%", weather.GetPath("main", "humidity").MustFloat64())).
		AddField("Visibility", fmt.Sprintf("%.2f m", weather.Get("visibility").MustFloat64())).
		AddField("Wind Speed", fmt.Sprintf("%.2f Km/h", math.Round(weather.GetPath("wind", "speed").MustFloat64()*3.6))).
		AddField("Wind Gust", fmt.Sprintf("%.2f Km/h", math.Round(weather.GetPath("wind", "gust").MustFloat64()*3.6))).
		AddField("Wind Direction", fmt.Sprintf("%.2f°", weather.GetPath("wind", "deg").MustFloat64())).
		Send(ctx)
}

func (NetworkCommands) hypesquad(ctx *dgc.Ctx) {
	house := utils.Arguments["<bravery/brilliance/balance>"]
	house = strings.Replace(house, "bravery", "1", -1)
	house = strings.Replace(house, "brilliance", "2", -1)
	house = strings.Replace(house, "balance", "3", -1)
	if house != "1" && house != "2" && house != "3" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<bravery/brilliance/balance>")
		return
	}
	api.SetHypesquad(house)
}

func (NetworkCommands) tinyurl(ctx *dgc.Ctx) {
	url := ctx.Arguments.Raw()
	tinyurl, err := http.Get(fmt.Sprintf("http://tinyurl.com/api-create.php?url=%s", url))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("URL", url).
		AddField("TinyURL", tinyurl).
		Send(ctx)
}

func (NetworkCommands) ipinfo(ctx *dgc.Ctx) {
	ip := utils.Arguments["<ip>"]
	ipinfo, err := http.GetJson(fmt.Sprintf("http://ipinfo.io/%s/json", ip))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if ipinfo.Get("hostname").MustInt(0) == 404 {
		ctx.RespondText("No results found.")
		return
	}
	embed.NewEmbed().
		AddField("IP", ip).
		AddField("Hostname", ipinfo.Get("hostname").MustString()).
		AddField("City", ipinfo.Get("city").MustString()).
		AddField("Region", ipinfo.Get("region").MustString()).
		AddField("Country", ipinfo.Get("country").MustString()).
		AddField("Location", ipinfo.Get("loc").MustString()).
		Send(ctx)
}

func (NetworkCommands) youtubedl(ctx *dgc.Ctx) {
	results, err := googlesearch.Search(context.TODO(), fmt.Sprintf("%s site:youtube.com", ctx.Arguments.Raw()))
	if err != nil {
		fmt.Println("An error has occurred.")
		return
	}
	for _, result := range results {
		if strings.Contains(result.URL, "/watch?v=") {
			client := youtube.Client{}
			video, err := client.GetVideo(result.URL)
			if err != nil {
				ctx.RespondText("An error has occurred: " + err.Error())
				return
			}
			ctx.RespondText(fmt.Sprintf("Downloading %s...", video.Title))
			formats := video.Formats.WithAudioChannels()
			stream, _, err := client.GetStream(video, &formats[0])
			if err != nil {
				ctx.RespondText("An error has occurred: " + err.Error())
				return
			}
			filepath := fmt.Sprintf("%s\\%s", utils.GetWindowsFolder(0), utils.File.SanitizeFilename(fmt.Sprintf("%s.mp4", video.Title)))
			file, err := os.Create(filepath)
			if err != nil {
				ctx.RespondText("An error has occurred: " + err.Error())
				return
			}
			defer file.Close()
			_, err = io.Copy(file, stream)
			if err != nil {
				ctx.RespondText("An error has occurred: " + err.Error())
				return
			}
			exec.Command("explorer", filepath).Run()
			ctx.RespondText("Successfully downloaded " + video.Title + ".")
			break
		}
	}
}

func (NetworkCommands) lvbypass(ctx *dgc.Ctx) {
	result, err := http.PostJson("https://api.bypass.vip/", fmt.Sprintf("{\"url\": \"%s\"}", ctx.Arguments.Raw()), map[string]string{"Content-Type": "application/json"})
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(fmt.Sprintf("Bypassed Link: %s", result.Get("destination").MustString()))
}

func (NetworkCommands) portscan(ctx *dgc.Ctx) {
	host := utils.Arguments["<host>"]
	ports := map[string]int{
		"FTP":            21,
		"SSH":            22,
		"Telnet":         23,
		"SMTP":           25,
		"DNS":            53,
		"HTTP":           80,
		"POP3":           110,
		"RPC":            135,
		"IMAP":           143,
		"HTTPS":          443,
		"SMB":            445,
		"MS-SQL":         1433,
		"MySQL":          3306,
		"Remote Desktop": 3389,
		"VNC":            5900,
		"HTTP Alt":       8080,
	}
	emb := embed.NewEmbed()
	for name, port := range ports {
		if utils.Misc.IsPortOpen(host, port) {
			emb.AddField(fmt.Sprintf("%s (%d)", name, port), "Open")
		}
	}
	emb.Send(ctx)
}

func (NetworkCommands) theme(ctx *dgc.Ctx) {
	theme := strings.ToLower(utils.Arguments["<light/dark>"])
	if theme != "light" && theme != "dark" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<light/dark> must be either \"light\" or \"dark\"")
		return
	}
	api.ChangeSettings(map[string]interface{}{"theme": theme}, true)
}

func (NetworkCommands) hastebin(ctx *dgc.Ctx) {
	text := ctx.Arguments.Raw()
	result, err := http.PostJson("https://www.toptal.com/developers/hastebin/documents", text)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(fmt.Sprintf("https://hastebin.com/%s", result.Get("key").MustString()))
}

func (NetworkCommands) checktoken(ctx *dgc.Ctx) {
	token := utils.Arguments["<token>"]
	response, statusCode, err := http.GetJsonStatusCode("https://discord.com/api/v10/users/@me", http.GetDiscordHeaders(token))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if statusCode != 200 {
		embed.NewEmbed().
			AddField("Status", "Invalid Token").
			Send(ctx)
		return
	}
	embed.NewEmbed().
		AddField("Status", "Valid Token").
		AddField("Username", fmt.Sprintf("%s#%s", response.Get("username").MustString(), response.Get("discriminator").MustString())).
		AddField("ID", response.Get("id").MustString()).
		AddField("Public Flags", fmt.Sprintf("%.f", response.Get("public_flags").MustFloat64())).
		AddField("Flags", fmt.Sprintf("%.f", response.Get("flags").MustFloat64())).
		AddField("Purchased Flags", fmt.Sprintf("%.f", response.Get("purchased_flags").MustFloat64())).
		AddField("Banner", response.Get("banner").MustString()).
		AddField("Banner Color", response.Get("banner_color").MustString()).
		AddField("Bio", response.Get("bio").MustString()).
		AddField("Locale", response.Get("locale").MustString()).
		AddField("NSFW Allowed", utils.Misc.FormatBool(response.Get("nsfw_allowed").MustBool())).
		AddField("2FA Enabled", utils.Misc.FormatBool(response.Get("mfa_enabled").MustBool())).
		AddField("Email", response.Get("email").MustString()).
		AddField("Verified", utils.Misc.FormatBool(response.Get("verified").MustBool())).
		AddField("Phone", response.Get("phone").MustString()).
		Send(ctx)
}

func (NetworkCommands) delwebhook(ctx *dgc.Ctx) {
	http.Delete(utils.Arguments["<webhook url>"])
	embed.NewEmbed().
		SetDescription("Successfully deleted webhook.").
		Send(ctx)
}

func (NetworkCommands) compact(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<on/off> must be either \"on\" or \"off\"")
		return
	}
	api.ChangeSettings(map[string]interface{}{"message_display_compact": toggle == "on"}, true)
}

func (NetworkCommands) devmode(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<on/off> must be either \"on\" or \"off\"")
		return
	}
	api.ChangeSettings(map[string]interface{}{"developer_mode": toggle == "on"}, true)
}

func (NetworkCommands) resolve(ctx *dgc.Ctx) {
	domain := utils.Arguments["<domain>"]
	response, err := http.GetJson(fmt.Sprintf("https://json.geoiplookup.io/%s", domain))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if response.Get("ip").MustString() == "" {
		embed.NewEmbed().SetDescription("Invalid domain.").Send(ctx)
		return
	}
	embed.NewEmbed().
		AddField("Domain", domain).
		AddField("IP", response.Get("ip").MustString()).
		Send(ctx)
}

func (NetworkCommands) reversedns(ctx *dgc.Ctx) {
	ip := utils.Arguments["<ip>"]
	response, err := http.GetJson(fmt.Sprintf("https://api.domaintools.com/v1/%s/host-domains/", ip))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if response.Get("error").MustArray() != nil {
		embed.NewEmbed().SetDescription("Invalid IP Address.").Send(ctx)
		return
	}
	domains := response.GetPath("response", "ip_addresses", "domain_names").MustStringArray()
	embed.NewEmbed().
		AddField("IP Address", ip).
		AddField("Domains", strings.Join(domains, ", ")).
		Send(ctx)
}

func (NetworkCommands) reversemac(ctx *dgc.Ctx) {
	mac := utils.Arguments["<mac>"]
	response, err := http.Get(fmt.Sprintf("https://api.macvendors.com/%s", mac))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	if strings.Contains(response, "Not Found") {
		embed.NewEmbed().SetDescription("Invalid MAC Address.").Send(ctx)
		return
	}
	embed.NewEmbed().
		AddField("MAC Address", mac).
		AddField("Vendor", response).
		Send(ctx)
}

func (NetworkCommands) amazon(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Amazon Search", fmt.Sprintf("https://www.amazon.com/s?k=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (NetworkCommands) steam(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Amazon Search", fmt.Sprintf("https://store.steampowered.com/search/?term=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (NetworkCommands) github(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Amazon Search", fmt.Sprintf("https://github.com/search?q=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (NetworkCommands) stackoverflow(ctx *dgc.Ctx) {
	query := ctx.Arguments.Raw()
	embed.NewEmbed().
		AddField("Amazon Search", fmt.Sprintf("https://stackoverflow.com/search?q=%s", url.QueryEscape(query))).
		Send(ctx)
}

func (NetworkCommands) covid(ctx *dgc.Ctx) {
	country := ctx.Arguments.Raw()
	covid, err := http.GetJson(fmt.Sprintf("https://coronavirus-19-api.herokuapp.com/countries/%s", country))
	if err != nil {
		ctx.RespondText("Country not found.")
		return
	}
	embed.NewEmbed().
		SetAuthor(fmt.Sprintf("COVID-19 - %s", country)).
		AddField("Total Cases", fmt.Sprint(covid.Get("cases").MustInt())).
		AddField("Cases Today", fmt.Sprint(covid.Get("todayCases").MustInt())).
		AddField("Total Deaths", fmt.Sprint(covid.Get("deaths").MustInt())).
		AddField("Deaths Today", fmt.Sprint(covid.Get("todayDeaths").MustInt())).
		AddField("Recovered", fmt.Sprint(covid.Get("recovered").MustInt())).
		AddField("Active", fmt.Sprint(covid.Get("active").MustInt())).
		AddField("Critical", fmt.Sprint(covid.Get("critical").MustInt())).
		AddField("Cases/Million", fmt.Sprint(covid.Get("casesPerOneMillion").MustInt())).
		AddField("Deaths/Million", fmt.Sprint(covid.Get("deathsPerOneMillion").MustInt())).
		AddField("Total Tests", fmt.Sprint(covid.Get("totalTests").MustInt())).
		AddField("Tests/Million", fmt.Sprint(covid.Get("testsPerOneMillion").MustInt())).
		Send(ctx)
}

func (NetworkCommands) discordstatus(ctx *dgc.Ctx) {
	status, err := http.GetJson("https://srhpyqt94yxb.statuspage.io/api/v2/components.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emb := embed.NewEmbed()
	components := status.Get("components")
	for i := 0; i < len(components.MustArray()); i++ {
		component := components.GetIndex(i)
		emb.AddField(component.Get("name").MustString(), component.Get("status").MustString())
	}
	emb.Send(ctx)
}

func (NetworkCommands) open(ctx *dgc.Ctx) {
	utils.Process.OpenURL(ctx.Arguments.Raw())
}

func (NetworkCommands) searchcountry(ctx *dgc.Ctx) {
	response, err := http.Get(fmt.Sprintf("https://restcountries.com/v3.1/name/%s?fullText=true", ctx.Arguments.Raw()))
	countryList, err2 := json.SimpleUnmarshal(response)
	if err != nil || err2 != nil {
		ctx.RespondText("Country not found.")
		return
	}
	if status, ok := countryList.CheckGet("status"); ok {
		if statusInt, err := status.Int(); err == nil && statusInt == 404 {
			ctx.RespondText("Country not found.")
			return
		}
	}
	country := countryList.GetIndex(0)
	rawCurrencies := country.Get("currencies").Interface().(map[string]interface{})
	currencies := []string{}
	for shortName, ifa := range rawCurrencies {
		currencies = append(currencies, fmt.Sprintf("%s (%s)", shortName, ifa.(map[string]interface{})["name"].(string)))
	}
	rawLanguages := country.Get("languages").Interface().(map[string]interface{})
	languages := []string{}
	for _, language := range rawLanguages {
		languages = append(languages, language.(string))
	}
	embed.NewEmbed().
		SetAuthor(fmt.Sprintf("%s - %s", country.GetPath("name", "common").MustString(), country.Get("cca2").MustString())).
		AddField("Name", country.GetPath("name", "common").MustString()).
		AddField("Top-Level Domain", strings.Join(country.Get("tld").MustStringArray(), ", ")).
		AddField("Country Code", country.Get("cca2").MustString()).
		AddField("Capital", strings.Join(country.Get("capital").MustStringArray(), ", ")).
		AddField("Region", country.Get("region").MustString()).
		AddField("Subregion", country.Get("subregion").MustString()).
		AddField("Population", fmt.Sprint(country.Get("population").MustInt())).
		AddField("Latitude/Longitude", fmt.Sprintf("%f, %f", country.Get("latlng").GetIndex(0).MustFloat64(), country.Get("latlng").GetIndex(1).MustFloat64())).
		AddField("Timezone", strings.Join(country.Get("timezones").MustStringArray(), ", ")).
		AddField("Currency", strings.Join(currencies, ", ")).
		AddField("Language", strings.Join(languages, ", ")).
		Send(ctx)
}

func (NetworkCommands) randomcountry(ctx *dgc.Ctx) {
	response, err := http.Get("https://restcountries.com/v3.1/all")
	countryList, err2 := json.SimpleUnmarshal(response)
	if err != nil || err2 != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	country := countryList.GetIndex(utils.Random.Number(0, len(countryList.MustArray())))
	rawCurrencies := country.Get("currencies").Interface().(map[string]interface{})
	currencies := []string{}
	for shortName, ifa := range rawCurrencies {
		currencies = append(currencies, fmt.Sprintf("%s (%s)", shortName, ifa.(map[string]interface{})["name"].(string)))
	}
	rawLanguages := country.Get("languages").Interface().(map[string]interface{})
	languages := []string{}
	for _, language := range rawLanguages {
		languages = append(languages, language.(string))
	}
	embed.NewEmbed().
		SetAuthor(fmt.Sprintf("%s - %s", country.GetPath("name", "common").MustString(), country.Get("cca2").MustString())).
		AddField("Name", country.GetPath("name", "common").MustString()).
		AddField("Top-Level Domain", strings.Join(country.Get("tld").MustStringArray(), ", ")).
		AddField("Country Code", country.Get("cca2").MustString()).
		AddField("Capital", strings.Join(country.Get("capital").MustStringArray(), ", ")).
		AddField("Region", country.Get("region").MustString()).
		AddField("Subregion", country.Get("subregion").MustString()).
		AddField("Population", fmt.Sprint(country.Get("population").MustInt())).
		AddField("Latitude/Longitude", fmt.Sprintf("%f, %f", country.Get("latlng").GetIndex(0).MustFloat64(), country.Get("latlng").GetIndex(1).MustFloat64())).
		AddField("Timezone", strings.Join(country.Get("timezones").MustStringArray(), ", ")).
		AddField("Currency", strings.Join(currencies, ", ")).
		AddField("Language", strings.Join(languages, ", ")).
		Send(ctx)
}

func (NetworkCommands) headers(ctx *dgc.Ctx) {
	url := utils.Arguments["<url>"]
	headers, err := http.GetHeadersString(url)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("URL", url).
		AddField("Headers", headers).
		Send(ctx)
}

func (NetworkCommands) _interface(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Public IP", api.GetPublicIP()).
		AddField("Local IP", api.GetLocalIP()).
		AddField("MAC Address", api.GetMacAddress()).
		Send(ctx)
}

func (NetworkCommands) revimage(ctx *dgc.Ctx) {
	imageUrl := utils.Arguments["<image url>"]
	embed.NewEmbed().
		AddField("Reverse Image Search URL", fmt.Sprintf("https://www.google.com/searchbyimage?image_url=%s", imageUrl)).
		Send(ctx)
}

func (NetworkCommands) pastebin(ctx *dgc.Ctx) {
	paste, err := api.CreatePaste(ctx.Arguments.Raw())
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Paste URL", paste).
		Send(ctx)
}

func (NetworkCommands) pastefile(ctx *dgc.Ctx) {
	filePath := utils.Arguments["<filepath>"]
	if content, success := utils.File.Read(filePath); success {
		paste, err := api.CreatePaste(content)
		if err != nil {
			ctx.RespondText("An error has occurred.")
			return
		}
		embed.NewEmbed().
			AddField("Paste URL", paste).
			Send(ctx)
	} else {
		ctx.RespondText("That file does not exist.")
	}
}
