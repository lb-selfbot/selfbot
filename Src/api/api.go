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

	"github.com/LightningDev1/LB-Selfbot/encryption"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/bitly/go-simplejson"
	"github.com/go-ping/ping"
)

func CreatePaste(text string) (string, error) {
	paste, err := http.PostJsonOld("https://paste.lightning-bot.com/documents", text)
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
	stats := pinger.Statistics()
	return stats.AvgRtt, nil
}

func WikipediaSearch(query string) (string, error) {
	search, err := http.GetJsonOld(fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=&titles=%s&format=json", url.QueryEscape(query)))
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

func SetHypesquad(session *discordgo.Session, hypesquad string) (err error) {
	hypesquadInt, err := strconv.Atoi(hypesquad)
	if err != nil {
		return err
	}

	data := map[string]int{
		"house_id": hypesquadInt,
	}

	_, err = session.Request("POST", "https://discord.com/api/v10/hypesquad/online", data)

	return err
}

func LeaveHypesquad(session *discordgo.Session) (err error) {
	_, err = session.Request("DELETE", "https://discord.com/api/v10/hypesquad/online", nil)

	return err
}

type RedditPost struct {
	Title    string
	ImageURL string
}

func GetImageFromSubreddit(subreddit string) (*RedditPost, bool) {
	sortBy := []string{"best", "new", "top", "hot"}[rand.Intn(4)]
	rawSearch, statusCode, err := http.GetJsonStatusCodeOld(fmt.Sprintf("https://www.reddit.com/r/%s/%s.json?limit=20", subreddit, sortBy), map[string]string{
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
	ip, err := http.GetOld("https://ipinfo.io/ip")
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
	ttsUrl := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), url.QueryEscape(language))
	tts, _, err := http.GetBytesOld(ttsUrl)
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

	translateUrl := fmt.Sprintf("https://translate.google.com/translate_a/single?client=gtx&sl=auto&tl=%s&dt=t&q=%s", url.QueryEscape(language), url.QueryEscape(text))
	jsonResult, err := http.GetJsonOld(translateUrl)
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
	data := map[string]interface{}{
		"channel_id": channelID,
		"guild_id":   guildID,
		"message_id": messageID,
		"reason":     reason,
	}

	_, err := session.Request("POST", "https://discord.com/api/v10/report", data)

	return err
}

func GoogleSearch(query string) (*simplejson.Json, error) {
	jsonResult, err := http.GetJsonOld("https://lightning-bot.com/api/v3/selfbot-utils/googlesearch?query=" + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}
	return jsonResult.Get("results"), nil
}

type CreditCard struct {
	Number     string
	CCV        string
	ExpiryDate string
}

func GenerateCreditCard(ccType string) CreditCard {
	jsonResult, err := http.GetJsonOld("https://lightning-bot.com/api/v3/selfbot-utils/ccgen?type=" + ccType)
	if err != nil {
		return CreditCard{}
	}
	return CreditCard{
		Number:     jsonResult.Get("number").MustString(),
		CCV:        jsonResult.Get("ccv").MustString(),
		ExpiryDate: jsonResult.Get("expiry_date").MustString(),
	}
}

func AllowFriendRequestsFrom(session *discordgo.Session, types []string) error {
	return ChangeSettings(
		session,
		map[string]interface{}{
			"friend_source_flags": map[string]bool{
				"all":            utils.SliceContains(types, "all"),
				"mutual_friends": utils.SliceContains(types, "mutual_friends"),
				"mutual_guilds":  utils.SliceContains(types, "mutual_guilds"),
			},
		},
		true,
	)
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
	encryptedStartParams, _, err := http.GetBytesOld("https://lightning-bot.com/api/v3/selfbot-utils/startparams")
	if err != nil {
		return nil
	}
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

func Auth(username, password string) (string, error) {
	data := map[string]string{
		"username":    username,
		"password":    password,
		"client_type": "ui",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return http.PostOld(
		"https://lightning-bot.com/api/v3/auth/auth",
		string(jsonData),
		map[string]string{
			"Content-Type": "application/json",
		},
	)
}
