package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/url"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/text/unicode/norm"

	"github.com/LightningDev/LBGo/error_handling"
	"github.com/LightningDev/LBGo/http"
	"github.com/LightningDev/LBGo/ipc"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/PuerkitoBio/goquery"
)

var LowercaseAlphabet string = "abcdefghijklmnopqrstuvwxyz"

var Arguments map[string]string = make(map[string]string)

func SendNotification(text string) {
	jsonBytes, _ := json.Marshal(ipc.JsonEvent{
		Op: ipc.Ipc.OpCodes.NOTIFICATION,
		Data: ipc.JsonEventData{
			Text: text,
		},
	})
	ipc.Ipc.Write(string(jsonBytes))
}

func ReverseString(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func FilterAlphaNumeric(str string) string {
	s := []byte(str)
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

func FlipText(text string, unflip bool) string {
	flipChars := [][]string{
		{"a", "ɐ"},
		{"b", "q"},
		{"c", "ɔ"},
		{"d", "p"},
		{"e", "ǝ"},
		{"f", "ɟ"},
		{"g", "ƃ"},
		{"h", "ɥ"},
		{"i", "ᴉ"},
		{"j", "ɾ"},
		{"k", "ʞ"},
		{"m", "ɯ"},
		{"n", "u"},
		{"r", "ɹ"},
		{"t", "ʇ"},
		{"v", "ʌ"},
		{"w", "ʍ"},
		{"y", "ʎ"},
		{"A", "∀"},
		{"C", "Ɔ"},
		{"E", "Ǝ"},
		{"F", "Ⅎ"},
		{"G", "פ"},
		{"H", "H"},
		{"I", "I"},
		{"J", "ſ"},
		{"L", "˥"},
		{"M", "W"},
		{"N", "N"},
		{"P", "Ԁ"},
		{"T", "┴"},
		{"U", "∩"},
		{"V", "Λ"},
		{"Y", "⅄"},
		{"1", "Ɩ"},
		{"2", "ᄅ"},
		{"3", "Ɛ"},
		{"4", "ㄣ"},
		{"5", "ϛ"},
		{"6", "9"},
		{"7", "ㄥ"},
		{"8", "8"},
		{"9", "6"},
		{"0", "0"},
		{".", "˙"},
		{",", "\x27"},
		{"\x27", ","},
		{"\x22", ",,"},
		{"`", ","},
		{"?", "¿"},
		{"!", "¡"},
		{"[", "]"},
		{"]", "["},
		{"(", ")"},
		{")", "("},
		{"{", "}"},
		{"}", "{"},
		{"<", ">"},
		{">", "<"},
		{"&", "⅋"},
		{"_", "‾"},
		{"∴", "∵"},
		{"⁅", "⁆"},
	}
	if unflip {
		for _, element := range flipChars {
			text = strings.ReplaceAll(text, element[1], element[0])
		}
	} else {
		for _, element := range flipChars {
			text = strings.ReplaceAll(text, element[0], element[1])
		}
	}
	return ReverseString(text)
}

func RemoveWords(text string, amount int) string {
	words := strings.Split(text, " ")
	if amount > len(words) {
		return ""
	}
	return strings.Join(words[amount:], " ")
}

func EncodingCommand(ctx *dgc.Ctx, encodeFunc func(string, bool) string) {
	encodeOrDecode := ctx.Arguments.Get(0).Raw()
	text := RemoveWords(ctx.Arguments.Raw(), 1)
	if encodeOrDecode == "encode" {
		ctx.RespondText(encodeFunc(text, false))
	} else if encodeOrDecode == "decode" {
		ctx.RespondText(encodeFunc(text, true))
	} else {
		ctx.RespondText("Invalid option for <encode/decode: " + encodeOrDecode)
	}
}

func HashCommand(ctx *dgc.Ctx, hashFunc func(string) string) {
	text := ctx.Arguments.Raw()
	ctx.RespondText(hashFunc(text))
}

func UnmarshalJson(jsonString string) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonMap)
	return jsonMap, err
}

func GetWindowsFolder(folder int) string {
	user32DLL := windows.NewLazyDLL("shell32.dll")
	getFolderPath := user32DLL.NewProc("SHGetFolderPathW")
	buffer := [windows.MAX_PATH]uint16{}
	nilPointer := uintptr(unsafe.Pointer(nil))
	getFolderPath.Call(nilPointer, uintptr(folder), nilPointer, 0, uintptr(unsafe.Pointer(&buffer[0])))
	return windows.UTF16ToString(buffer[:])
}

func WindowsGetDocumentsFolder() string {
	return GetWindowsFolder(5)
}

func GetMainFolder() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\LightningBot", WindowsGetDocumentsFolder())
	} else {
		// TODO: Linux support
		return "/.config/discordgo/"
	}
}

func SeedRandom() {
	rand.Seed(time.Now().UnixNano())
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomChoice(choices interface{}) interface{} {
	reflectType := reflect.TypeOf(choices)
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	if reflectType.Kind() != reflect.Slice {
		return "An error has occurred"
	}
	val := reflect.ValueOf(choices)
	return val.Index(rand.Intn(val.Len())).Interface()
}

func RandomPercentage() string {
	return fmt.Sprintf("%d%%", RandomNumber(0, 100))
}

func RandomString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(RandomNumber(65, 90))
	}
	return string(bytes)
}

func TryGetUserFromUserMention(session *discordgo.Session, mention *dgc.Argument) (*discordgo.User, bool) {
	if mention.AsUserMentionID() == "" {
		return nil, false
	}
	user, err := session.User(mention.AsUserMentionID())
	if err != nil {
		return nil, false
	}
	return user, true
}

func TryGetNameFromUserMention(session *discordgo.Session, mention *dgc.Argument) (string, bool) {
	if mention.AsUserMentionID() == "" {
		return mention.Raw(), false
	}
	user, err := session.User(mention.AsUserMentionID())
	if err != nil {
		return "", false
	}
	return fmt.Sprintf("%s#%s", user.Username, user.Discriminator), true
}

func GetElementById(html string, id string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return doc.Find(id).Text()
}

func TTS(text string, language string) ([]byte, error) {
	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), url.QueryEscape(language))
	tts, _, err := http.GetBytes(url)
	if err != nil {
		return nil, err
	}
	return tts, nil
}

func Capitalize(text string) string {
	return strings.ToUpper(text[:1]) + text[1:]
}

func Shuffle(slice []string) []string {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func CheckIfInServer(ctx *dgc.Ctx) bool {
	if ctx.Event.GuildID == "" {
		error_handling.OnError(ctx, error_handling.ServerOnlyCommand, ctx.Command.Name)
	}
	return ctx.Event.GuildID != ""
}

func GetImageUrlFromJson(url string, jsonQuery string) (string, error) {
	json, err := http.GetJson(url)
	if err != nil {
		return "", err
	}
	return json[jsonQuery].(string), nil
}

func SliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SanitizeFilename(filename string) string {
	blacklist := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", "\x00"}
	reserved := []string{
		"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5",
		"COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5",
		"LPT6", "LPT7", "LPT8", "LPT9",
	}
	for _, c := range blacklist {
		filename = strings.Replace(filename, c, "", -1)
	}
	filename = norm.NFKD.String(filename)
	filename = strings.TrimSuffix(filename, ". ")
	filename = strings.TrimSpace(filename)
	isDots := true
	for _, c := range filename {
		// Remove all characters in filename below code point 32
		if c < ' ' {
			filename = strings.Replace(filename, string(c), "", -1)
		}
		if isDots && c != '.' {
			isDots = false
		}
	}
	if isDots {
		filename = "__" + filename
	}
	if SliceContains(reserved, filename) {
		filename = "__" + filename
	}
	if len(filename) > 255 {
		filename = filename[:255]
	}
	return filename
}

func IsPortOpen(hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func FormatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func OpenURL(url string) {
	if runtime.GOOS == "windows" {
		exec.Command("explorer", url).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("open", url).Start()
	} else {
		exec.Command("xdg-open", url).Start()
	}
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

func GetFileContents(file string) (string, bool) {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return "", false
	}
	return string(fileContents), true
}

func GetArguments(ctx *dgc.Ctx) (map[string]string, bool) {
	rawArguments := strings.Replace(ctx.Command.Usage, fmt.Sprintf("[p]%s", ctx.Command.Name), "", -1)
	commandArguments := regexp.MustCompile(`(<(.*?)>|\[(.*?)\])`).FindAllString(rawArguments, -1)
	arguments := map[string]string{}
	content := ctx.Event.Message.Content
	for index, arg := range commandArguments {
		raw := ctx.Arguments.Get(index).Raw()
		isRequired := strings.Contains(arg, "<") && strings.Contains(arg, ">")

		if strings.Contains(arg, "\"") && strings.Count(raw, "\"") >= 2 {
			arguments[arg] = regexp.MustCompile("\"(.*?)\"").FindString(content)
			content = strings.Replace(content, arguments[arg], "", -1)
		} else {
			arguments[arg] = raw
			content = RemoveWords(content, 1)
		}

		if arguments[arg] == "" && isRequired {
			error_handling.OnError(ctx, error_handling.MissingRequiredArgument, arg)
			return map[string]string{}, false
		}
	}
	return arguments, true
}

func JsonEscape(i string) string {
    b, err := json.Marshal(i)
    if err != nil {
        panic(err)
    }
    // Trim the beginning and trailing " character
    return string(b[1:len(b)-1])
}
