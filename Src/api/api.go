package api

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/bitly/go-simplejson"
	"github.com/go-ping/ping"
)

func CreatePaste(text string) (string, error) {
	paste, err := http.PostJson("https://paste.lightning-bot.com/documents", text)
	if err != nil {
		return "", err
	}
	return "https://paste.lightning-bot.com/" + paste.Get("key").MustString(), nil
}

func Ping(host string) (time.Duration, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return time.Duration(0), err
	}
	pinger.Count = 3
	pinger.SetPrivileged(true)
	err = pinger.Run()
	if err != nil {
		return time.Duration(0), err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	return stats.AvgRtt, nil
}

func WikipediaSearch(query string) (string, error) {
	search, err := http.GetJson(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=&titles=%s&format=json", url.QueryEscape(query)))
	if err != nil {
		return "", err
	}
	pages := search.GetPath("query", "pages").MustMap()
	for _, page := range pages {
		page := page.(map[string]interface{})
		if page["missing"] != nil {
			return "No results found", nil
		}
		return page["extract"].(string), nil
	}
	return "No results found", nil
}

func SetHypesquad(hypesquad string) {
	hypesquadInt, _ := strconv.Atoi(hypesquad)
	jsonData := fmt.Sprintf(`{"hypesquad": %d}`, hypesquadInt)
	cfg, err := config.GetConfig()
	if err != nil {
		return
	}
	http.Post("https://discord.com/api/v10/hypesquad/online", jsonData, http.GetDiscordHeaders(cfg.Account.Discord.Token))
}

type RedditPost struct {
	Title    string
	ImageURL string
}

func GetImageFromSubreddit(subreddit string) (*RedditPost, bool) {
	sortBy := []string{"best", "new", "top", "hot"}[utils.RandomNumber(0, 3)]
	rawSearch, statusCode, err := http.GetJsonStatusCode(fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?limit=20", subreddit, sortBy), map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36 Edg/96.0.1054.62",
	})
	if err != nil {
		return nil, false
	}
	if statusCode != 200 {
		return nil, false
	}
	return ExtractRedditURL(rawSearch, 0)
}

func ExtractRedditURL(jsonData *simplejson.Json, tries int) (*RedditPost, bool) {
	if tries > 25 {
		return nil, false
	}
	posts := jsonData.GetPath("data", "children")
	amount := len(posts.MustArray())
	post := posts.GetIndex(utils.RandomNumber(0, amount-1)).Get("data")
	if regexp.MustCompile(`(\.jpg|\.png|\.gif|\.jpeg)`).MatchString(post.Get("url_overridden_by_dest").MustString()) {
		return &RedditPost{
			Title:    post.Get("title").MustString(),
			ImageURL: strings.ReplaceAll(post.Get("url_overridden_by_dest").MustString(), " .gifv", ".gif"),
		}, true
	}
	if post.Get("is_video").MustBool() || post.Get("media").MustMap() == nil {
		return ExtractRedditURL(jsonData, tries+1)
	}
	url := post.GetPath("media", "oembed", "thumbnail_url").MustString()
	if strings.Contains(url, "gfycat") {
		return ExtractRedditURL(jsonData, tries+1)
	}
	return &RedditPost{
		Title:    post.Get("title").MustString(),
		ImageURL: strings.ReplaceAll(url, " .gifv", ".gif"),
	}, true
}

func GetPublicIP() string {
	ip, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		return "Error"
	}
	return ip
}

func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "Error"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetMacAddress() string {
	currentIP := GetLocalIP()
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		addrs, err := interf.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if strings.Contains(addr.String(), currentIP) {
				netInterface, err := net.InterfaceByName(interf.Name)
				if err != nil {
					continue
				}
				return strings.ToUpper(netInterface.HardwareAddr.String())
			}
		}
	}
	return "Error"
}
