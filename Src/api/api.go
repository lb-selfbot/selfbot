package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/http"
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
	sortBy := []string{"best", "new", "top", "hot"}[rand.Intn(4)]
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
	post := posts.GetIndex(rand.Intn(amount - 1)).Get("data")
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

func TTS(text, language string) ([]byte, error) {
	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), url.QueryEscape(language))
	tts, _, err := http.GetBytes(url)
	return tts, err
}

type TranslateResult struct {
	Original            string
	Result              string
	SourceLanguage      string
	DestinationLanguage string
	Confidence          float64
	Error               error
}

func Translate(text, language string) *TranslateResult {
	language = strings.ToLower(language)
	language = strings.Split(language, "_")[0]

	if _, ok := constants.Languages[language]; !ok {
		if _, ok := constants.LanguageSpecialCases[language]; ok {
			language = constants.LanguageSpecialCases[language]
		} else if _, ok := constants.LanguageCodes[language]; ok {
			language = constants.LanguageCodes[language]
		} else {
			return &TranslateResult{
				Error: fmt.Errorf("invalid destination language"),
			}
		}
	}

	url := fmt.Sprintf("https://translate.google.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=t&q=%s", url.QueryEscape(language), url.QueryEscape(text))
	jsonResult, err := http.GetJson(url)
	if err != nil {
		return &TranslateResult{
			Error: err,
		}
	}

	translateData := jsonResult.GetIndex(0).GetIndex(0)
	return &TranslateResult{
		Original:            translateData.GetIndex(1).MustString(),
		Result:              translateData.GetIndex(0).MustString(),
		SourceLanguage:      jsonResult.GetIndex(2).MustString(),
		DestinationLanguage: language,
		Confidence:          jsonResult.GetIndex(6).MustFloat64(),
	}
}

func ChangeSettings(settings map[string]interface{}, settingsPath bool) (*simplejson.Json, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}
	url := ""
	if settingsPath {
		url = "https://discord.com/api/v10/users/@me/settings"
	} else {
		url = "https://discord.com/api/v10/users/@me"
	}
	return http.PatchJson(url, string(jsonBytes), http.GetDiscordHeaders(cfg.Account.Discord.Token))
}

func ReportMessage(guildID, channelID, messageID string, reason int) (*simplejson.Json, error) {
	payload := map[string]interface{}{
		"channel_id": channelID,
		"guild_id":   guildID,
		"message_id": messageID,
		"reason":     reason,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return http.PostJson("https://discord.com/api/v10/report", string(jsonBytes), http.GetDiscordHeaders(cfg.Account.Discord.Token))
}
