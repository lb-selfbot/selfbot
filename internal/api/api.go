package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/random"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bitly/go-simplejson"
	"github.com/lb-selfbot/selfbot/internal/constants"
	"github.com/lb-selfbot/selfbot/internal/http"
)

func CreatePaste(text string) (string, error) {
	resp, err := http.PostJson("https://paste.lightning-bot.com/documents", []byte(text))
	if err != nil {
		return "", err
	}

	return "https://paste.lightning-bot.com/%s" + resp.Get("key").MustString(), nil
}

func CreateMeme(templateID, text0, text1 string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "An error has occurred: " + err.Error(), err
	}

	resp, err := http.PostJson("https://api.lightning-bot.com/api/v3/selfbot/memegen", map[string]any{
		"username":    cfg.Account.LightningBot.Username,
		"password":    cfg.Account.LightningBot.Password,
		"template_id": templateID,
		"text0":       text0,
		"text1":       text1,
	})
	if err != nil {
		return "An error has occurred: " + err.Error(), err
	}

	imageURL := resp.Get("url").MustString()
	if imageURL == "" {
		return "An error has occurred: " + resp.Get("error").MustString(), fmt.Errorf("empty response")
	}

	return imageURL, nil
}

func Ping(host string) (time.Duration, error) {
	beforeRequest := time.Now()

	httpResponse := http.CustomMethodWithData(host, "HEAD", nil)
	if httpResponse.Error != nil {
		return time.Duration(0), httpResponse.Error
	}

	return time.Since(beforeRequest), nil
}

func WikipediaSearch(query string) (string, error) {
	searchUrl := fmt.Sprintf(
		"https://en.wikipedia.org/w/api.php?action=query&prop=extracts&explaintext=&titles=%s&format=json",
		url.QueryEscape(query),
	)

	resp, err := http.GetJson(searchUrl)
	if err != nil {
		return "No results found", err
	}

	pages := resp.Get("query").Get("pages")

	for pageId := range pages.MustMap() {
		page := pages.Get(pageId)

		if _, missing := page.CheckGet("missing"); missing {
			continue
		}

		return page.Get("extract").MustString(), nil
	}

	return "No results found", fmt.Errorf("no results found")
}

func SetHypesquad(session *discordgo.Session, hypesquad int) (err error) {
	_, err = session.Request(
		"POST",
		"https://discord.com/api/v9/hypesquad/online",
		map[string]int{"house_id": hypesquad},
	)

	return err
}

func LeaveHypesquad(session *discordgo.Session) (err error) {
	_, err = session.Request("DELETE", "https://discord.com/api/v9/hypesquad/online", nil)

	return err
}

var RedditHeaders = map[string]string{
	"Authority":                 "www.reddit.com",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"Accept-Language":           "en-US,en;q=0.9",
	"Cache-Control":             "max-age=0",
	"Sec-CH-UA":                 "\"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"",
	"Sec-CH-UA-Mobile":          "?0",
	"Sec-CH-UA-Platform":        "\"Windows\"",
	"Sec-Fetch-Dest":            "empty",
	"Sec-Fetch-Mode":            "cors",
	"Sec-Fetch-Site":            "same-origin",
	"Upgrade-Insecure-Requests": "1",
	"User-Agent":                discordgo.UserAgentWeb,
}

type RedditPost struct {
	Title    string
	ImageURL string
}

func GetImageFromSubreddit(subreddit string) (*RedditPost, bool) {
	url := fmt.Sprintf("https://api.pullpush.io/reddit/submission/search?html_decode=True&subreddit=%s&size=50", subreddit)

	resp, err := http.GetJson(url, RedditHeaders)
	if err != nil {
		return nil, false
	}

	data := resp.Get("data")
	dataLen := len(data.MustArray())

	if dataLen == 0 {
		return nil, false
	}

	// Iterating over the posts randomly
	idxSlice := make([]int, dataLen)
	for i := 0; i < dataLen; i++ {
		idxSlice[i] = i
	}

	idxSlice = random.Shuffle(idxSlice)

	for i := 0; i < dataLen; i++ {
		post := data.GetIndex(idxSlice[i])

		if post.Get("post_hint").MustString() != "image" {
			continue
		}

		url := post.Get("url").MustString()
		if regexp.MustCompile(`(\.jpg|\.png|\.gif|\.jpeg)`).MatchString(url) {
			return &RedditPost{
				Title:    post.Get("title").MustString(),
				ImageURL: url,
			}, true
		}
	}

	return nil, false
}

func GetPublicIP() string {
	resp := http.Get("https://api.ipify.org")
	if resp.Error != nil {
		return "Error"
	}

	return string(resp.Body)
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

	resp, err := http.GetJson(translateUrl)
	if err != nil {
		return &TranslateResult{Error: err}
	}

	result := resp.GetIndex(0).GetIndex(0)
	return &TranslateResult{
		Original:            result.GetIndex(1).MustString(),
		Result:              result.GetIndex(0).MustString(),
		SourceLanguage:      resp.GetIndex(2).MustString(),
		DestinationLanguage: language,
		Confidence:          resp.GetIndex(6).MustFloat64(),
	}
}

func GoogleSearch(query string) (*simplejson.Json, error) {
	url := fmt.Sprintf("https://api.lightning-bot.com/api/v3/misc/search?query=%s", url.QueryEscape(query))

	resp, err := http.GetJson(url)
	if err != nil {
		return nil, err
	}

	return resp.Get("results"), nil
}

type CreditCard struct {
	Number        string
	CVV           string
	ExpiryDate    string
	Visualisation string
}

func GenerateCreditCard(ccType, name string) CreditCard {
	url := fmt.Sprintf("https://api.lightning-bot.com/api/v3/misc/ccgen?type=%s&name=%s", ccType, name)

	resp, err := http.GetJson(url)
	if err != nil {
		return CreditCard{}
	}

	return CreditCard{
		Number:        resp.Get("number").MustString(),
		CVV:           resp.Get("cvv").MustString(),
		ExpiryDate:    resp.Get("expiry_date").MustString(),
		Visualisation: resp.Get("visualisation").MustString(),
	}
}

func OCR(image string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	data := map[string]string{
		"username": cfg.Account.LightningBot.Username,
		"password": cfg.Account.LightningBot.Password,
		"image":    image,
	}

	resp, err := http.PostJson("https://api.lightning-bot.com/api/v3/selfbot/ocr", data)
	if err != nil {
		return "", err
	}

	return resp.Get("text").String()
}

func Auth(username, password string) *http.HttpResponse {
	return http.Post(
		"https://api.lightning-bot.com/api/v3/auth/auth",
		map[string]any{
			"username":    username,
			"password":    password,
			"client_type": "ui",
		},
	)
}

func GenImageKey(imageUrl string) string {
	resp, err := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/gen-image-key",
		map[string]any{"image_url": imageUrl},
	)
	if err != nil {
		return "Error generating image key: " + err.Error()
	}

	return resp.Get("image_key").MustString()
}

func Obfuscate(lang, code string) string {
	resp, err := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/obfuscate",
		map[string]any{"lang": lang, "code": code},
	)

	code = resp.Get("code").MustString()
	if code == "" {
		return "An error has occurred: " + err.Error()
	}

	return code
}

type URLEmbed struct {
	Title        string
	Description  string
	URL          string
	Color        string
	Image        string
	Large        bool
	Video        string
	ProviderName string
	ProviderURL  string
	AuthorName   string
	AuthorURL    string
}

func GenerateURLEmbed(embed URLEmbed) (string, error) {
	const embedURL = "https://embed.lightning-bot.com/api/embed"
	values := url.Values{
		"title":         {embed.Title},
		"description":   {embed.Description},
		"url":           {embed.URL},
		"color":         {embed.Color},
		"image":         {embed.Image},
		"large":         {fmt.Sprintf("%t", embed.Large)},
		"video":         {embed.Video},
		"provider_name": {embed.ProviderName},
		"provider_url":  {embed.ProviderURL},
		"author_name":   {embed.AuthorName},
		"author_url":    {embed.AuthorURL},
	}
	httpResponse := http.Post(
		embedURL+"?"+values.Encode(),
		nil,
	)
	if httpResponse.Error != nil {
		return "", httpResponse.Error
	}

	return embedURL + "?key=" + string(httpResponse.Body), nil
}

type Pokemon struct {
	ID            string `json:"id"`
	Region        string `json:"region"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	NameJapanese  string `json:"name.ja"`
	NameJapaneseT string `json:"name.ja_t"`
	NameEnglish   string `json:"name.en"`
	NameGerman    string `json:"name.de"`
	NameFrench    string `json:"name.fr"`
	Type          string `json:"type.0"`
	Event         string `json:"event"`
	Height        string `json:"height"`
	Weight        string `json:"weight"`
	HP            string `json:"base.hp"`
	Attack        string `json:"base.atk"`
	Defense       string `json:"base.def"`
}

func PokeTwo(image string, cfg *config.Config) (*Pokemon, float64, error) {
	httpResponse := http.Post(
		"https://lightningdev1-poketwo.hf.space/run/predict",
		map[string]any{"data": []any{image}},
	)
	if httpResponse.Error != nil {
		return nil, 0, httpResponse.Error
	}

	var resp struct {
		Data []struct {
			Pokemon    *Pokemon `json:"pokemon"`
			Confidence float64  `json:"confidence"`
		}
	}

	err := json.Unmarshal(httpResponse.Body, &resp)
	if err != nil {
		return nil, 0, err
	}

	if len(resp.Data) == 0 {
		return nil, 0, fmt.Errorf("no data")
	}

	return resp.Data[0].Pokemon, resp.Data[0].Confidence, nil
}

func PokeTwoVerify(userID string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	resp, err := http.PostJson(
		"https://api.lightning-bot.com/api/v3/selfbot/poketwo-verify",
		map[string]any{
			"username": cfg.Account.LightningBot.Username,
			"password": cfg.Account.LightningBot.Password,
			"user_id":  userID,
		},
	)
	if err != nil {
		return err
	}

	taskID := resp.Get("task_id").MustString()
	if taskID == "" {
		return fmt.Errorf("no task ID")
	}

	var solution string

	const maxTries = 20

	for i := 0; i < maxTries; i++ {
		time.Sleep(time.Duration(2000+500*i) * time.Millisecond)

		resp, err = http.PostJson(
			"https://api.lightning-bot.com/api/v3/selfbot/poketwo-status",
			map[string]any{
				"username": cfg.Account.LightningBot.Username,
				"password": cfg.Account.LightningBot.Password,
				"task_id":  taskID,
			},
		)
		if err != nil {
			return err
		}

		if resp.Get("status").MustString() == "ready" {
			solution = resp.Get("solution").MustString()
			break
		}
	}

	if solution == "" {
		return fmt.Errorf("failed to solve captcha")
	}

	params := url.Values{
		"token": {solution},
		"uid":   {userID},
	}

	httpResponse := http.Post(
		"https://verify.poketwo.net/api/recaptcha?"+params.Encode(),
		nil,
		map[string]string{
			"Accept":                      "*/*",
			"Accept-Language":             "en-US,en;q=0.9",
			"Content-Length":              "0",
			"Origin":                      "https://verify.poketwo.net",
			"Referer":                     "https://verify.poketwo.net/captcha/" + userID,
			"Sec-Ch-Ua":                   "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"119\"",
			"Sec-Ch-Ua-Arch":              "x86",
			"Sec-Ch-Ua-Bitness":           "64",
			"Sec-Ch-Ua-Full-Version":      "119.0.0.0",
			"Sec-Ch-Ua-Full-Version-List": "\"Not_A Brand\";v=\"8.0.0.0\", \"Chromium\";v=\"119.0.0.0\"",
			"Sec-Ch-Ua-Mobile":            "?0",
			"Sec-Ch-Ua-Model":             "",
			"Sec-Ch-Ua-Platform":          "Windows",
			"Sec-Ch-Ua-Platform-Version":  "10.0.0",
			"Sec-Fetch-Dest":              "empty",
			"Sec-Fetch-Mode":              "cors",
			"Sec-Fetch-Site":              "same-origin",
			"User-Agent":                  discordgo.UserAgentWeb,
		},
	)
	fmt.Println(httpResponse.StatusCode, string(httpResponse.Body))
	if httpResponse.Error != nil {
		return httpResponse.Error
	}

	if httpResponse.StatusCode != 204 {
		return fmt.Errorf("error verifying captcha: %d %s", httpResponse.StatusCode, httpResponse.Body)
	}

	return nil
}

func GetTokenData(session *discordgo.Session, token string) ([]byte, error) {
	headers := map[string]string{
		"Authorization": token,
	}

	return session.RequestWithBucketID("GET", discordgo.EndpointUser("@me"), nil, discordgo.EndpointUser(""), headers)
}

func ChangeSettings(session *discordgo.Session, path string, settings map[string]any) error {
	_, err := session.Request("PATCH", "https://discord.com/api/v9/users/@me"+path, settings)

	return err
}

func AllowFriendRequests(session *discordgo.Session, allow bool) error {
	if allow {
		session.State.UserSettings.Privacy.FriendSourceFlags = wrapperspb.UInt32(14)
	} else {
		session.State.UserSettings.Privacy.FriendSourceFlags = wrapperspb.UInt32(0)
	}

	return session.UserSettingsUpdate(session.State.UserSettings, true)
}

// AI

func GPTCustom(messages []any, models ...string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "An error has occurred: " + err.Error(), err
	}

	model := "gpt-3.5-turbo"
	if len(models) > 0 {
		model = models[0]
	}

	resp, err := http.PostJson("https://api.lightning-bot.com/api/v3/selfbot/gpt", map[string]any{
		"username": cfg.Account.LightningBot.Username,
		"password": cfg.Account.LightningBot.Password,
		"model":    model,
		"messages": messages,
	})
	if err != nil {
		return "An error has occurred: " + err.Error(), err
	}

	content := resp.Get("response").MustString()

	if content == "" {
		data := fmt.Sprint(resp.Interface())
		return "An error has occurred: " + data, fmt.Errorf("empty response: %s", data)
	}

	return content, nil
}

func GPT(instructions, message string, model ...string) (string, error) {
	messages := []any{}

	if instructions != "" {
		messages = append(messages, map[string]any{"role": "system", "content": instructions})
	}

	if message != "" {
		messages = append(messages, map[string]any{"role": "user", "content": message})
	}

	return GPTCustom(messages, model...)
}

var DefaultStyles = map[string]string{
	"anime":        "anime atmospheric, atmospheric anime, anime character; full body art, digital anime art, beautiful anime art style, anime picture, anime arts, beautiful anime style, digital advanced anime art, anime painting, anime artwork, beautiful anime art, detailed digital anime art, anime epic artwork",
	"space":        "in cosmic atmosphere, humanitys cosmic future,  space art concept, space landscape, scene in space, cosmic space, beautiful space star planet, background in space, realistic, cinematic, breathtaking view",
	"comic":        "comic book, john romita senior, inspired by Alton Tobey, by Alan Davis, arachnophobia, by Alton Tobey, as a panel of a marvel comic, marvel comic",
	"dystopian":    "cifi world, cybernetic civilizations, peter gric and dan mumford,   brutalist dark futuristic, dystopian brutalist atmosphere, dark dystopian world, cinematic 8k, end of the world, doomsday",
	"surrealism":   "salvador dali painting, highly detailed surrealist art, surrealist conceptual art,  masterpiece surrealism, surreal architecture, surrealistic digital artwork, whimsical surrealism, bizarre art",
	"graffiti":     "graffiti background, colorful graffiti, graffiti art style, colorful mural, ravi supa, symbolic mural, juxtapoz, pablo picasso, street art",
	"medieval":     "movie still from game of thrones, powerful fantasy epic, middle ages, lush green landscape, olden times, roman empire, 1400 ce, highly detailed background, cinematic lighting, 8k render, high quality, bright colours",
	"fantasy":      "fantasy matte painting,  fantasy landscape, ( ( thomas kinkade ) ), whimsical, dreamy, alice in wonderland, daydreaming, epic scene, high exposure, highly detailed, tim white, michael whelan",
	"abstract":     "vibrant, editorial, abstract elements, colorful, color splatter, realism, Inspired by the style of Ryan Hewett, dynamic realism, soft lighting and intricate details",
	"painting":     "atmospheric dreamscape painting, by Mac Conner, vibrant gouache painting scenery, vibrant painting, vivid painting, a beautiful painting, dream scenery art, instagram art, psychedelic painting, lofi art, bright art",
	"icon":         "single vector graphics icon, ios icon, smooth shape, vector",
	"render":       "isometric, polaroid octane render, 3 d render 1 5 0 mm lens, keyshot product render, rendered, keyshot product render pinterest, 3 d product render, 3 d cgi render, 3d cgi render, ultra wide angle isometric view",
	"cyberpunk":    "synthwave image, (neotokyo), dreamy colorful cyberpunk colors, cyberpunk blade runner art, retrofuturism, cyberpunk, beautiful cyberpunk style, cgsociety 9",
	"landscape":    "landscape 4k, beautiful landscapes, nature wallpaper, 8k photography",
	"architecture": "modern architecture design, luxury architecture, bright, very beautiful, trending on unsplash, breath taking",
	"retro":        "retro futuristic illustration, featured on illustrationx, art deco illustration, beautiful retro art, stylized digital illustration, highly detailed vector art, ( ( mads berg ) ), automotive design art, epic smooth illustration, by mads berg, stylized illustration, ash thorp khyzyl saleem, clean vector art",
	"vibrant":      "Psychedelic, water colors spots, vibrant color scheme, highly detailed, romanticism style, cinematic, artstation, greg rutkowski",
	"futuristic":   "futuristic, elegant atmosphere, glowing lights, highly detailed, digital painting, artstation, concept art, smooth sharp focus, illustration, mars ravelo, gereg rutkowski",
}

var DefaultModels = map[string]string{
	"3 Guofeng3 V3.4":                         "3Guofeng3_v34.safetensors [50f420de]",
	"Absolute Reality V1.6":                   "absolutereality_V16.safetensors [37db0fc3]",
	"Absolute Reality V1.8.1":                 "absolutereality_v181.safetensors [3d9d4d2b]",
	"Am I Real V4.1":                          "amIReal_V41.safetensors [0a8a2e61]",
	"Analog V1":                               "analog-diffusion-1.0.ckpt [9ca13f02]",
	"Anything V3":                             "anythingv3_0-pruned.ckpt [2700c435]",
	"Anything V4.5":                           "anything-v4.5-pruned.ckpt [65745d25]",
	"Anything V5":                             "anythingV5_PrtRE.safetensors [893e49b9]",
	"AbyssOrangeMix V3":                       "AOM3A3_orangemixs.safetensors [9600da17]",
	"Blazing Drive V10g":                      "blazing_drive_v10g.safetensors [ca1c1eab]",
	"CetusMix Version35":                      "cetusMix_Version35.safetensors [de2f2560]",
	"Children's Stories V1 3D":                "childrensStories_v13D.safetensors [9dfaabcb]",
	"Children's Stories V1 SemiReal":          "childrensStories_v1SemiReal.safetensors [a1c56dbb]",
	"Children's Stories V1 Toon-Anime":        "childrensStories_v1ToonAnime.safetensors [2ec7b88b]",
	"Counterfeit V3.0":                        "Counterfeit_v30.safetensors [9e2a8f19]",
	"CuteYukimix MidChapter3":                 "cuteyukimixAdorable_midchapter3.safetensors [04bdffe6]",
	"CyberRealistic V3.3":                     "cyberrealistic_v33.safetensors [82b0d085]",
	"Dalcefo V4":                              "dalcefo_v4.safetensors [425952fe]",
	"Deliberate V2":                           "deliberate_v2.safetensors [10ec4b29]",
	"Deliberate V3":                           "deliberate_v3.safetensors [afd9d2d4]",
	"Dreamlike Anime V1":                      "dreamlike-anime-1.0.safetensors [4520e090]",
	"Dreamlike Diffusion V1":                  "dreamlike-diffusion-1.0.safetensors [5c9fd6e0]",
	"Dreamlike Photoreal V2":                  "dreamlike-photoreal-2.0.safetensors [fdcf65e7]",
	"Dreamshaper 6 baked vae":                 "dreamshaper_6BakedVae.safetensors [114c8abb]",
	"Dreamshaper 7":                           "dreamshaper_7.safetensors [5cf5ae06]",
	"Dreamshaper 8":                           "dreamshaper_8.safetensors [9d40847d]",
	"Edge of Realism EOR V2.0":                "edgeOfRealism_eorV20.safetensors [3ed5de15]",
	"Eimis Anime Diffusion V1.0":              "EimisAnimeDiffusion_V1.ckpt [4f828a15]",
	"Elldreth's Vivid":                        "elldreths-vivid-mix.safetensors [342d9d26]",
	"EpiCRealism Natural Sin RC1":             "epicrealism_naturalSinRC1VAE.safetensors [90a4c676]",
	"I Cant Believe Its Not Photography Seco": "ICantBelieveItsNotPhotography_seco.safetensors [4e7a3dfd]",
	"Juggernaut Aftermath":                    "juggernaut_aftermath.safetensors [5e20c455]",
	"Lofi V4":                                 "lofi_v4.safetensors [ccc204d6]",
	"Lyriel V1.6":                             "lyriel_v16.safetensors [68fceea2]",
	"MajicMix Realistic V4":                   "majicmixRealistic_v4.safetensors [29d0de58]",
	"MechaMix V1.0":                           "mechamix_v10.safetensors [ee685731]",
	"MeinaMix Meina V9":                       "meinamix_meinaV9.safetensors [2ec66ab0]",
	"MeinaMix Meina V11":                      "meinamix_meinaV11.safetensors [b56ce717]",
	"Neverending Dream V1.22":                 "neverendingDream_v122.safetensors [f964ceeb]",
	"Openjourney V4":                          "openjourney_V4.ckpt [ca2f377f]",
	"Pastel-Mix":                              "pastelMixStylizedAnime_pruned_fp16.safetensors [793a26e8]",
	"Portrait+ V1":                            "portraitplus_V1.0.safetensors [1400e684]",
	"Protogen x3.4":                           "protogenx34.safetensors [5896f8d5]",
	"Realistic Vision V1.4":                   "Realistic_Vision_V1.4-pruned-fp16.safetensors [8d21810b]",
	"Realistic Vision V2.0":                   "Realistic_Vision_V2.0.safetensors [79587710]",
	"Realistic Vision V4.0":                   "Realistic_Vision_V4.0.safetensors [29a7afaa]",
	"Realistic Vision V5.0":                   "Realistic_Vision_V5.0.safetensors [614d1063]",
	"Redshift Diffusion V1.0":                 "redshift_diffusion-V10.safetensors [1400e684]",
	"ReV Animated V1.2.2":                     "revAnimated_v122.safetensors [3f4fefd9]",
	"RunDiffusion FX 2.5D V1.0":               "rundiffusionFX25D_v10.safetensors [cd12b0ee]",
	"RunDiffusion FX Photorealistic V1.0":     "rundiffusionFX_v10.safetensors [cd4e694d]",
	"SD V1.4":                                 "sdv1_4.ckpt [7460a6fa]",
	"SD V1.5":                                 "v1-5-pruned-emaonly.safetensors [d7049739]",
	"Shonin's Beautiful People V1.0":          "shoninsBeautiful_v10.safetensors [25d8c546]",
	"TheAlly's Mix II":                        "theallys-mix-ii-churned.safetensors [5d9225a4]",
	"Timeless V1":                             "timeless-1.0.ckpt [7c4971d4]",
	"ToonYou Beta 6":                          "toonyou_beta6.safetensors [980f6b15]",
}

func GetImageModels() map[string]string {
	httpResponse := http.Get("https://app.prodia.com")
	if httpResponse.Error != nil {
		return DefaultModels
	}

	script := regexp.MustCompile(`<script defer="defer" src="\/js\/app.(.*?).js"><\/script>`).FindStringSubmatch(string(httpResponse.Body))
	if len(script) < 2 {
		return DefaultModels
	}

	scriptUrl := fmt.Sprintf("https://app.prodia.com/js/app.%s.js", script[1])

	httpResponse = http.Get(scriptUrl)
	if httpResponse.Error != nil {
		return DefaultModels
	}

	modelsList := regexp.MustCompile(`VUE_APP_AI_MODELS:'(.*?)',`).FindSubmatch(httpResponse.Body)
	if len(modelsList) < 2 {
		return DefaultModels
	}

	models := bytes.ReplaceAll(modelsList[1], []byte(`\'`), []byte(`'`))

	_ = json.Unmarshal(models, &DefaultModels)

	return DefaultModels
}

func CreateImageJob(prompt, negPrompt, model string) (string, error) {
	seed := rand.Intn(1000000)

	params := url.Values{
		"new":             {"true"},
		"prompt":          {prompt},
		"model":           {model},
		"negative_prompt": {negPrompt},
		"steps":           {"30"},
		"cfg":             {"9.5"},
		"seed":            {strconv.Itoa(seed)},
		"sampler":         {"DPM++ 2M Karras"},
		"aspect_ratio":    {"square"},
	}

	headers := map[string]string{
		"authority":          "api.prodia.com",
		"accept":             "*/*",
		"accept-language":    "en-US,en;q=0.6",
		"dnt":                "1",
		"origin":             "https://app.prodia.com",
		"referer":            "https://app.prodia.com/",
		"sec-ch-ua":          "\"Chromium\";v=\"119\", \"Not-A.Brand\";v=\"24\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"sec-gpc":            "1",
		"user-agent":         discordgo.UserAgentWeb,
	}

	resp, err := http.GetJson("https://api.prodia.com/generate?"+params.Encode(), headers)
	if err != nil {
		return "", err
	}

	job := resp.Get("job").MustString()
	if job == "" {
		return "", fmt.Errorf("could not create job: %v", resp.Interface())
	}

	return job, nil
}

func GenerateImage(prompt, negPrompt, model string) (*bytes.Reader, error) {
	if strings.TrimSpace(negPrompt) == "" {
		negPrompt = "porn, nsfw, verybadimagenegative_v1.3, ng_deepnegative_v1_75t, ugly face, caricature, body horror, mutant, facebook, youtube, food"
	}

	if strings.TrimSpace(model) == "" {
		model = "Realistic_Vision_V5.0.safetensors [614d1063]"
	}

	job, err := CreateImageJob(prompt, negPrompt, model)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"authority": "api.prodia.com",
		"accept":    "*/*",
	}

	start := time.Now()
	maxDuration := 30 * time.Second

	for {
		time.Sleep(500 * time.Millisecond)
		resp, err := http.GetJson("https://api.prodia.com/job/"+job, headers)
		if err != nil {
			return nil, err
		}

		status := resp.Get("status").MustString()

		if status == "succeeded" {
			break
		}

		if time.Since(start) > maxDuration {
			return nil, fmt.Errorf("timed out")
		}
	}

	imageUrl := fmt.Sprintf("https://images.prodia.xyz/%s.png?download=1", job)

	httpResponse := http.Get(imageUrl)
	if httpResponse.Error != nil {
		return nil, httpResponse.Error
	}

	return bytes.NewReader(httpResponse.Body), nil
}

// ReviewDB

var reviewDBToken string

func GetReviewDBToken(session *discordgo.Session) (string, error) {
	if reviewDBToken != "" {
		return reviewDBToken, nil
	}

	authorizeUrl := "https://discord.com/api/v9/oauth2/authorize?client_id=915703782174752809&response_type=code&redirect_uri=https%3A%2F%2Fmanti.vendicated.dev%2Fapi%2Freviewdb%2Fauth&scope=identify"

	rawData, err := session.Request("POST", authorizeUrl, map[string]any{
		"authorize":    true,
		"permissions":  "0",
		"user_install": false,
	})
	if err != nil {
		return "", err
	}

	var data map[string]any
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return "", err
	}

	location, ok := data["location"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response: no location: %s", string(rawData))
	}

	resp, err := http.GetJson(location + "&clientMod=vencord")
	if err != nil {
		return "", err
	}

	token := resp.Get("token").MustString()
	if token == "" {
		return "", fmt.Errorf("invalid response: no token")
	}

	return token, nil
}

type Review struct {
	ID        int64
	Username  string
	Comment   string
	Timestamp time.Time
}

func GetReviewDBReviews(userID string) ([]Review, error) {
	resp, err := http.GetJson("https://manti.vendicated.dev/api/reviewdb/users/" + userID + "/reviews?flags=0&offset=0")
	if err != nil {
		return nil, err
	}

	if !resp.Get("success").MustBool() {
		return nil, fmt.Errorf("invalid response: %s", resp.Get("message").MustString())
	}

	rawReviews := resp.Get("reviews")
	reviews := []Review{}

	for i := 0; i < len(rawReviews.MustArray()); i++ {
		review := rawReviews.GetIndex(i)

		reviews = append(reviews, Review{
			ID:        review.Get("id").MustInt64(),
			Username:  review.GetPath("sender", "username").MustString(),
			Comment:   review.Get("comment").MustString(),
			Timestamp: time.Unix(review.Get("timestamp").MustInt64(), 0),
		})
	}

	return reviews, nil
}

func AddReviewDBReview(session *discordgo.Session, userID, review string) error {
	token, err := GetReviewDBToken(session)
	if err != nil {
		return err
	}

	return http.PutJson("https://manti.vendicated.dev/api/reviewdb/users/"+userID+"/reviews", map[string]any{
		"token":   strings.TrimSpace(token),
		"comment": review,
		"userid":  userID,
	}).Error
}

type weatherData struct {
	weather *simplejson.Json
	time    time.Time
}

var weatherCache = map[string]weatherData{}

func GetWeather(location string) (*simplejson.Json, error) {
	if w, ok := weatherCache[location]; ok && time.Since(w.time) < time.Minute*5 {
		return w.weather, nil
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?appid=533d00f29739edba54c5fb2b1a8262b3&q=%s",
		url.QueryEscape(location),
	)

	resp, err := http.GetJson(url)
	if err != nil {
		return nil, err
	}

	if resp.Get("cod").MustInt() != 200 {
		return nil, fmt.Errorf("invalid location")
	}

	weatherCache[location] = weatherData{
		weather: resp,
		time:    time.Now(),
	}

	return resp, nil
}
