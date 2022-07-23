package utils

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/text/unicode/norm"

	"github.com/LightningDev1/LB-Selfbot/error_handling"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/PuerkitoBio/goquery"
)

var LowercaseAlphabet string = "abcdefghijklmnopqrstuvwxyz"

var Arguments map[string]string = make(map[string]string)

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

func GetMainFolder() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\LightningBot", WindowsGetDocumentsFolder())
	} else {
		// TODO: Linux support
		return "/.config/discordgo/"
	}
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomChoice[T comparable](choices []T) T {
	return choices[rand.Intn(len(choices))]
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

func TryGetUserFromMention(session *discordgo.Session, mention *dgc.Argument) (*discordgo.User, bool) {
	userID := mention.AsUserMentionID()
	if userID == "" {
		userID = mention.Raw()
	}
	user, err := session.User(userID)
	if err != nil {
		return nil, false
	}
	return user, true
}

func TryGetMemberFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Member, bool) {
	userID := mention.AsUserMentionID()
	if userID == "" {
		userID = mention.Raw()
	}
	member, err := ctx.Session.GuildMember(ctx.Event.GuildID, userID)
	if err != nil {
		return nil, false
	}
	return member, true
}

func TryGetUsernameFromMention(session *discordgo.Session, mention *dgc.Argument) (string, bool) {
	user, success := TryGetUserFromMention(session, mention)
	if !success {
		return "", false
	}
	return fmt.Sprintf("%s#%s", user.Username, user.Discriminator), true
}

func TryGetChannelFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Channel, bool) {
	if mention.AsChannelMentionID() == "" {
		return nil, false
	}
	channel, err := ctx.Session.Channel(mention.AsChannelMentionID())
	if err != nil {
		return nil, false
	}
	return channel, true
}

func SnowflakeToTime(snowflake string) time.Time {
	// parse snowflake to int64
	snowflakeInt, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}
	}
	// convert int64 to time
	return time.Unix(((snowflakeInt>>22)+1420070400000)/1000, 0)
}

func GetElementById(html string, id string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return doc.Find(id).Text()
}

func Capitalize(text string) string {
	return strings.ToUpper(text[:1]) + text[1:]
}

func Shuffle[T comparable](slice []T) []T {
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
	return json.Get(jsonQuery).MustString(), nil
}

func SliceContains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
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

func Restart() error {
	self, err := os.Executable()
	if err != nil {
		return err
	}
	env := os.Environ()
	if runtime.GOOS == "windows" {
		cmd := exec.Command(self, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = env
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}
	return syscall.Exec(self, os.Args, env)
}

func GetChannelMessages(ctx *dgc.Ctx, limit int) ([]*discordgo.Message, error) {
	// The maximum number of messages that can be retrieved at once is 100.
	if limit <= 100 {
		return ctx.Session.ChannelMessages(ctx.Event.ChannelID, limit, "", "", ctx.Event.Message.ID)
	}

	amountOfRequests := int(math.Ceil(float64(limit) / 100))
	messages := []*discordgo.Message{ctx.Event.Message}
	// Subtract 1 from the limit because we already have the message that executed the command
	limit -= 1

	for i := 0; i < amountOfRequests; i++ {
		newMessages, err := ctx.Session.ChannelMessages(ctx.Event.ChannelID, int(math.Min(float64(limit), 100)), messages[len(messages)-1].ID, "", "")
		if err != nil {
			return nil, err
		}
		limit -= 100
		messages = append(messages, newMessages...)
	}

	return messages, nil
}

func GetRoleNames(ctx *dgc.Ctx, member *discordgo.Member) string {
	roleNames := []string{}
	for _, role := range member.Roles {
		role, err := ctx.Session.State.Role(member.GuildID, role)
		if err == nil {
			roleNames = append(roleNames, role.Name)
		}
	}
	return strings.Join(roleNames, ", ")
}
