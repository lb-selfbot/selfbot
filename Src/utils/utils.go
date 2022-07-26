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

var InviteRegex *regexp.Regexp = regexp.MustCompile(`(?:https?://)?discord(?:(?:app)?\.com/invite|\.gg)/?[a-zA-Z0-9]+/?`)

type DiscordUtils struct{}

func (DiscordUtils) TryGetUserFromMention(session *discordgo.Session, mention *dgc.Argument) (*discordgo.User, bool) {
	user, err := session.User(Misc.Or(mention.AsUserMentionID(), mention.Raw()))
	if err != nil {
		return nil, false
	}
	return user, true
}

func (DiscordUtils) TryGetMemberFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Member, bool) {
	member, err := ctx.Session.GuildMember(ctx.Event.GuildID, Misc.Or(mention.AsUserMentionID(), mention.Raw()))
	if err != nil {
		return nil, false
	}
	return member, true
}

func (d *DiscordUtils) TryGetUsernameFromMention(session *discordgo.Session, mention *dgc.Argument) (string, bool) {
	user, success := d.TryGetUserFromMention(session, mention)
	if !success {
		return "", false
	}
	return fmt.Sprintf("%s#%s", user.Username, user.Discriminator), true
}

func (DiscordUtils) TryGetChannelFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Channel, bool) {
	channel, err := ctx.Session.Channel(Misc.Or(mention.AsChannelMentionID(), mention.Raw()))
	if err != nil {
		return nil, false
	}
	return channel, true
}

func (DiscordUtils) TryGetRoleFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Channel, bool) {
	channel, err := ctx.Session.Channel(Misc.Or(mention.AsRoleMentionID(), mention.Raw()))
	if err != nil {
		return nil, false
	}
	return channel, true
}

func (DiscordUtils) SnowflakeToTime(snowflake string) time.Time {
	// parse snowflake to int64
	snowflakeInt, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}
	}
	// convert int64 to time
	return time.Unix(((snowflakeInt>>22)+1420070400000)/1000, 0)
}

func (DiscordUtils) GetTimestamp(t *time.Time, defaultString string) string {
	if t == nil {
		return defaultString
	}
	return fmt.Sprintf("<t:%d>", t.Unix())
}

func (DiscordUtils) ParseInvite(invite string) string {
	inviteRegexMatch := string(InviteRegex.Find([]byte(invite)))
	if inviteRegexMatch != "" {
		splitInvite := strings.Split(inviteRegexMatch, "/")
		return splitInvite[len(splitInvite)-1]
	}
	return invite
}

func (DiscordUtils) CheckIfInServer(ctx *dgc.Ctx) bool {
	if ctx.Event.GuildID == "" {
		error_handling.OnError(ctx, error_handling.ServerOnlyCommand, ctx.Command.Name)
	}
	return ctx.Event.GuildID != ""
}

func (DiscordUtils) GetArguments(ctx *dgc.Ctx) (map[string]string, bool) {
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
			content = String.RemoveWords(content, 1)
		}

		if arguments[arg] == "" && isRequired {
			error_handling.OnError(ctx, error_handling.MissingRequiredArgument, arg)
			return map[string]string{}, false
		}
	}
	return arguments, true
}

func (DiscordUtils) GetChannelMessages(ctx *dgc.Ctx, limit int) ([]*discordgo.Message, error) {
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

func (DiscordUtils) GetRoleNames(ctx *dgc.Ctx, member *discordgo.Member) string {
	roleNames := []string{}
	for _, role := range member.Roles {
		role, err := ctx.Session.State.Role(member.GuildID, role)
		if err == nil {
			roleNames = append(roleNames, role.Name)
		}
	}
	return strings.Join(roleNames, ", ")
}

func (DiscordUtils) GetGuild(session *discordgo.Session, guildID string) (guild *discordgo.Guild, err error) {
	guild, err = session.State.Guild(guildID)
	if err != nil {
		guild, err = session.Guild(guildID)
		if err != nil {
			return nil, err
		}
	}
	return guild, nil
}

func (DiscordUtils) GetMember(session *discordgo.Session, guildID, userID string) (member *discordgo.Member, err error) {
	member, err = session.State.Member(guildID, userID)
	if err != nil {
		member, err = session.GuildMember(guildID, userID)
		if err != nil {
			return nil, err
		}
	}
	return member, nil
}

func (DiscordUtils) GetBoosters(guild *discordgo.Guild) []*discordgo.Member {
	boostedMembers := []*discordgo.Member{}
	for _, member := range guild.Members {
		if member.PremiumSince != nil {
			boostedMembers = append(boostedMembers, member)
		}
	}
	return boostedMembers
}

var Discord *DiscordUtils = &DiscordUtils{}

type StringUtils struct{}

func (StringUtils) Reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}

func (StringUtils) FilterAlphaNumeric(str string) string {
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

func (StringUtils) Flip(text string, unflip bool) string {
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
	return String.Reverse(text)
}

func (StringUtils) RemoveWords(text string, amount int) string {
	words := strings.Split(text, " ")
	if amount > len(words) {
		return ""
	}
	return strings.Join(words[amount:], " ")
}

func (StringUtils) Capitalize(text string) string {
	return strings.ToUpper(text[:1]) + text[1:]
}

var String *StringUtils = &StringUtils{}

type RandomUtils struct{}

func (RandomUtils) Number(min, max int) int {
	return rand.Intn(max-min) + min
}

func (RandomUtils) Percentage() string {
	return fmt.Sprintf("%d%%", Random.Number(0, 100))
}

func (RandomUtils) String(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(Random.Number(65, 90))
	}
	return string(bytes)
}

// Golang prevents you from using a generic function on a
// normal struct so these have to be outside RandomUtils
func RandomChoice[T comparable](choices []T) T {
	return choices[rand.Intn(len(choices))]
}

func RandomShuffle[T comparable](slice []T) []T {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

var Random *RandomUtils = &RandomUtils{}

type FileUtils struct{}

func (FileUtils) SanitizeFilename(filename string) string {
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

func (FileUtils) Read(file string) (string, bool) {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return "", false
	}
	return string(fileContents), true
}

func (FileUtils) GetMainFolder() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\LightningBot", WindowsGetDocumentsFolder())
	} else {
		// TODO: Linux support
		return "/.config/discordgo/"
	}
}

var File *FileUtils = &FileUtils{}

type ProcessUtils struct{}

func (ProcessUtils) OpenURL(url string) {
	if runtime.GOOS == "windows" {
		exec.Command("explorer", url).Start()
	} else if runtime.GOOS == "darwin" {
		exec.Command("open", url).Start()
	} else {
		exec.Command("xdg-open", url).Start()
	}
}

func (ProcessUtils) Restart() error {
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

var Process *ProcessUtils = &ProcessUtils{}

type MiscUtils struct{}

func (MiscUtils) GetElementById(html string, id string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	return doc.Find(id).Text()
}

func (MiscUtils) GetImageUrlFromJson(url string, jsonQuery string) (string, error) {
	json, err := http.GetJson(url)
	if err != nil {
		return "", err
	}
	return json.Get(jsonQuery).MustString(), nil
}

func (MiscUtils) IsPortOpen(hostname string, port int) bool {
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (MiscUtils) FormatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// If a is not empty, returns a. If a is empty, returns b.
func (MiscUtils) Or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

var Misc *MiscUtils = &MiscUtils{}

// Golang generics limitation
func SliceContains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
