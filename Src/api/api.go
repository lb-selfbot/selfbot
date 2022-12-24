package api

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/bitly/go-simplejson"
	"github.com/go-ping/ping"
)

func CreatePaste(text string) (string, error) {
	httpResult := http.Post("https://paste.lightning-bot.com/documents", []byte(text))
	if httpResult.Error != nil {
		return "", httpResult.Error
	}

	key := httpResult.Json().Get("key").MustString()

	return fmt.Sprintf("https://paste.lightning-bot.com/%s", key), nil
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
	stats := pinger.Statistics()
	return stats.AvgRtt, nil
}

func WikipediaSearch(query string) (string, error) {
	searchUrl := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=&titles=%s&format=json",
		url.QueryEscape(query),
	)

	httpResult := http.Get(searchUrl)
	if httpResult.Error != nil {
		return "No results found", httpResult.Error
	}

	pages := httpResult.Json().Get("query").Get("pages")

	for pageId := range pages.MustMap() {
		page := pages.Get(pageId)

		_, missing := page.CheckGet("missing")
		if missing {
			continue
		}

		return page.Get("extract").MustString(), nil
	}

	return "No results found", nil
}

func SetHypesquad(session *discordgo.Session, hypesquad string) (err error) {
	hypesquadInt, err := strconv.Atoi(hypesquad)
	if err != nil {
		return err
	}

	_, err = session.Request(
		"POST",
		"https://discord.com/api/v10/hypesquad/online",
		map[string]int{
			"house_id": hypesquadInt,
		},
	)

	return err
}

func LeaveHypesquad(session *discordgo.Session) (err error) {
	_, err = session.Request("DELETE", "https://discord.com/api/v10/hypesquad/online", nil)

	return err
}

var RedditHeaders = map[string]string{
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"Accept-Encoding":           "gzip, deflate, br",
	"Accept-Language":           "en-US,en;q=0.9",
	"Sec-Ch-Ua":                 "\"Opera GX\";v=\"93\", \"Not\\A)Brand\";v=\"99\", \"Chromium\";v=\"93\"",
	"Sec-Ch-Ua-Mobile":          "?0",
	"Sec-Fetch-Dest":            "document",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Site":            "none",
	"Sec-Fetch-User":            "?1",
	"Upgrade-Insecure-Requests": "1",
}

type RedditPost struct {
	Title    string
	ImageURL string
}

func GetImageFromSubreddit(subreddit string) (*RedditPost, bool) {
	sortBy := []string{"best", "new", "top", "hot"}[rand.Intn(4)]

	httpResponse := http.Get(
		fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?limit=20", subreddit, sortBy),
		RedditHeaders,
	)

	if httpResponse.Error != nil || httpResponse.StatusCode != 200 {
		return nil, false
	}

	return ExtractRedditURL(httpResponse.Json(), 0)
}

func ExtractRedditURL(jsonData *simplejson.Json, tries int) (*RedditPost, bool) {
	if tries > 25 {
		return nil, false
	}

	posts := jsonData.GetPath("data", "children")
	postsAmount := len(posts.MustArray())

	post := posts.GetIndex(rand.Intn(postsAmount - 1)).Get("data")

	if regexp.MustCompile(`(\.jpg|\.png|\.gif|\.jpeg)`).MatchString(post.Get("url_overridden_by_dest").MustString()) {
		return &RedditPost{
			Title:    post.Get("title").MustString(),
			ImageURL: strings.ReplaceAll(post.Get("url_overridden_by_dest").MustString(), " .gifv", ".gif"),
		}, true
	}

	if post.Get("is_video").MustBool() || post.Get("media").MustMap() == nil {
		return ExtractRedditURL(jsonData, tries+1)
	}

	postUrl := post.GetPath("media", "oembed", "thumbnail_url").MustString()
	if strings.Contains(postUrl, "gfycat") {
		return ExtractRedditURL(jsonData, tries+1)
	}

	return &RedditPost{
		Title:    post.Get("title").MustString(),
		ImageURL: strings.ReplaceAll(postUrl, " .gifv", ".gif"),
	}, true
}

func GetPublicIP() string {
	ip := http.Get("https://api.ipify.org").Body

	return utils.Misc.Or(ip, "Error")
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
			if !strings.Contains(addr.String(), currentIP) {
				continue
			}

			netInterface, err := net.InterfaceByName(interf.Name)
			if err != nil {
				continue
			}
			return strings.ToUpper(netInterface.HardwareAddr.String())
		}
	}

	return "Error"
}

func TTS(text, language string) *http.HttpResponse {
	ttsUrl := fmt.Sprintf(
		"https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s",
		url.QueryEscape(text),
		url.QueryEscape(language),
	)

	return http.Get(ttsUrl)
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

	if _, ok := constants.LANGUAGES[language]; !ok {
		if _, ok := constants.LANGUAGE_SPECIAL_CASES[language]; ok {
			language = constants.LANGUAGE_SPECIAL_CASES[language]
		} else if _, ok := constants.LANGUAGE_CODES[language]; ok {
			language = constants.LANGUAGE_CODES[language]
		} else {
			return &TranslateResult{
				Error: fmt.Errorf("invalid destination language"),
			}
		}
	}

	translateUrl := fmt.Sprintf(
		"https://translate.google.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=t&q=%s",
		url.QueryEscape(language),
		url.QueryEscape(text),
	)

	httpResponse := http.Get(translateUrl)
	if httpResponse.Error != nil {
		return &TranslateResult{
			Error: httpResponse.Error,
		}
	}

	json := httpResponse.Json()

	translateData := json.GetIndex(0).GetIndex(0)
	return &TranslateResult{
		Original:            translateData.GetIndex(1).MustString(),
		Result:              translateData.GetIndex(0).MustString(),
		SourceLanguage:      json.GetIndex(2).MustString(),
		DestinationLanguage: language,
		Confidence:          json.GetIndex(6).MustFloat64(),
	}
}

func ChangeSettings(session *discordgo.Session, settings map[string]any, settingsPath bool) error {
	settingsUrl := ""
	if settingsPath {
		settingsUrl = "https://discord.com/api/v10/users/@me/settings"
	} else {
		settingsUrl = "https://discord.com/api/v10/users/@me"
	}

	_, err := session.Request("PATCH", settingsUrl, settings)

	return err
}

func ReportMessage(session *discordgo.Session, guildID, channelID, messageID string, reason int) error {
	_, err := session.Request(
		"POST",
		"https://discord.com/api/v10/report",
		map[string]any{
			"channel_id": channelID,
			"guild_id":   guildID,
			"message_id": messageID,
			"reason":     reason,
		},
	)

	return err
}

func GoogleSearch(query string) (*simplejson.Json, error) {
	url := fmt.Sprintf("https://api.lightning-bot.com/api/v3/misc/search?query=%s", url.QueryEscape(query))

	httpResult := http.Get(url)
	if httpResult.Error != nil {
		return nil, httpResult.Error
	}

	return httpResult.Json().Get("results"), nil
}

type CreditCard struct {
	Number     string
	CCV        string
	ExpiryDate string
}

func GenerateCreditCard(ccType string) CreditCard {
	url := fmt.Sprintf("https://api.lightning-bot.com/api/v3/misc/ccgen?type=%s", ccType)

	httpResult := http.Get(url)
	if httpResult.Error != nil {
		return CreditCard{}
	}

	json := httpResult.Json()

	return CreditCard{
		Number:     json.Get("number").MustString(),
		CCV:        json.Get("ccv").MustString(),
		ExpiryDate: json.Get("expiry_date").MustString(),
	}
}

func AllowFriendRequestsFrom(session *discordgo.Session, types []string) error {
	newSettings := map[string]any{
		"friend_source_flags": map[string]any{
			"all":            utils.SliceContains(types, "all"),
			"mutual_friends": utils.SliceContains(types, "mutual_friends"),
			"mutual_guilds":  utils.SliceContains(types, "mutual_guilds"),
		},
	}

	return ChangeSettings(session, newSettings, true)
}

func GetStartParams() *simplejson.Json {
	// XOR encryption
	encrypted := []int{2, 15, 13, 48, 11, 25, 22, 16, 31, 51, 53, 13, 44, 20, 32, 7, 34, 53, 48, 37, 13, 5, 7, 62, 59, 30, 36, 20, 25, 7, 36, 0, 59, 3, 61, 56, 22, 8, 48, 34, 10, 54, 7, 8, 29, 58, 15, 52, 11, 49, 39, 62, 49, 39, 51, 17, 49, 21, 43, 19, 15, 61, 63, 11, 39, 58, 12, 2, 37, 22, 42, 11, 31, 57, 60, 41, 18, 45, 42, 33, 53, 50, 8, 48, 57, 2, 59, 63, 62, 50, 25, 43, 46, 42, 34, 17, 4, 35, 2, 24, 15}
	keyInt := []int{73, 104, 122, 122, 67, 99, 111, 100, 67, 111, 78, 104, 114, 117, 74, 115, 113, 119, 114, 120, 79, 106, 69, 111, 65, 106, 114, 78, 84, 65, 110, 103, 76, 99, 71, 86, 73, 114, 77, 122, 87, 102, 66, 103, 99, 99, 116, 66, 86}
	decrypted := ""
	for i := len(encrypted) - 1; i > -1; i-- {
		decrypted += string(rune(encrypted[i] ^ keyInt[(len(encrypted)-1-i)%len(keyInt)]))
	}
	keyEncryption := &encryption.AESCBC{Key: decrypted}
	encryptedKey := []byte{78, 70, 57, 114, 89, 55, 75, 67, 112, 105, 57, 68, 102, 104, 71, 80, 121, 53, 107, 86, 53, 78, 76, 74, 104, 85, 70, 78, 74, 51, 77, 121, 102, 53, 88, 121, 56, 43, 48, 89, 106, 77, 111, 61, 58, 87, 97, 53, 113, 115, 76, 51, 99, 116, 107, 118, 110, 79, 120, 48, 102, 56, 86, 57, 70, 117, 103, 61, 61, 58, 102, 104, 105, 66, 105, 113, 48, 116, 79, 119, 80, 54, 87, 82, 101, 71, 113, 118, 102, 117, 66, 103, 43, 83, 102, 88, 111, 53, 106, 109, 68, 86, 90, 55, 99, 82, 71, 87, 57, 121, 109, 72, 120, 120, 113, 55, 78, 56, 98, 85, 121, 116, 74, 86, 117, 70, 104, 100, 49, 122, 43, 74, 43, 67, 69, 78, 68, 66, 82, 54, 55, 57, 57, 105, 99, 103, 74, 51, 90, 87, 103, 73, 57, 52, 66, 119, 61, 61}
	key, _ := keyEncryption.Decrypt(encryptedKey)
	Encryption := &encryption.AESCBC{Key: string(key)}

	httpResponse := http.Get("https://lightning-bot.com/api/v3/selfbot/startparams")
	if httpResponse.Error != nil {
		return nil
	}

	encryptedStartParams := httpResponse.BodyBytes

	decryptedStartParams, err := Encryption.Decrypt(encryptedStartParams)
	if err != nil {
		return nil
	}

	jsonStartParams, err := simplejson.NewJson(decryptedStartParams)
	if err != nil {
		return nil
	}

	return jsonStartParams
}

func Auth(username, password string) *http.HttpResponse {
	return http.PostJson(
		"https://api.lightning-bot.com/api/v3/auth/auth",
		map[string]any{
			"username":    username,
			"password":    password,
			"client_type": "ui",
		},
	)
}

func GetTokenData(session *discordgo.Session, token string) ([]byte, error) {
	return session.RequestWithBucketID("GET", discordgo.EndpointUser("@me"), nil, discordgo.EndpointUser(""))
}
