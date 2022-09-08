package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
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
)

var LowercaseAlphabet = "abcdefghijklmnopqrstuvwxyz"

var Arguments = make(map[string]string)

var InviteRegex = regexp.MustCompile(`(?:https?://)?discord(?:(?:app)?\.com/invite|\.gg)/?[a-zA-Z\d]+/?`)

type DiscordUtils struct{}

func (DiscordUtils) TryGetUserFromMention(session *discordgo.Session, mention *dgc.Argument) (*discordgo.User, bool) {
	userID := Misc.Or(mention.AsUserMentionID(), mention.Raw())
	user, err := session.State.GetUser(userID)
	if err != nil {
		user, err = session.User(userID)
	}
	if err != nil {
		return nil, false
	}
	return user, true
}

func (DiscordUtils) TryGetMemberFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Member, bool) {
	userID := Misc.Or(mention.AsUserMentionID(), mention.Raw())
	member, err := ctx.Session.State.Member(ctx.Event.GuildID, userID)
	if err != nil {
		member, err = ctx.Session.GuildMember(ctx.Event.GuildID, Misc.Or(mention.AsUserMentionID(), mention.Raw()))
	}
	if err != nil {
		return nil, false
	}
	return member, true
}

func (d DiscordUtils) TryGetUsernameFromMention(session *discordgo.Session, mention *dgc.Argument) (string, bool) {
	user, success := d.TryGetUserFromMention(session, mention)
	if !success {
		return "", false
	}
	return user.String(), true
}

func (DiscordUtils) GetChannel(session *discordgo.Session, channelID string) (*discordgo.Channel, error) {
	channel, err := session.State.Channel(channelID)
	if err != nil {
		channel, err = session.Channel(channelID)
		if err != nil {
			return nil, err
		}
	}
	return channel, nil
}

func (d DiscordUtils) TryGetChannelFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Channel, bool) {
	channelID := Misc.Or(mention.AsChannelMentionID(), mention.Raw())
	channel, err := d.GetChannel(ctx.Session, channelID)
	if err != nil {
		return nil, false
	}
	return channel, true
}

func (DiscordUtils) TryGetRoleFromMention(ctx *dgc.Ctx, mention *dgc.Argument) (*discordgo.Role, bool) {
	roleID := Misc.Or(mention.AsRoleMentionID(), mention.Raw())
	role, err := ctx.Session.State.Role(ctx.Event.GuildID, roleID)
	if err != nil {
		role, err = ctx.Session.State.Role(ctx.Event.GuildID, roleID)
		if err != nil {
			return nil, false
		}
	}
	return role, true
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
	commandArguments := regexp.MustCompile(`(<(.*?)>|\[(.*?)])`).FindAllString(rawArguments, -1)
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

func (DiscordUtils) GetChannelMessages(session *discordgo.Session, channelID string, limit int) ([]*discordgo.Message, error) {
	// The maximum number of messages that can be retrieved at once is 100.
	if limit <= 100 {
		return session.ChannelMessages(channelID, limit, "", "", "")
	}

	amountOfRequests := int(math.Ceil(float64(limit) / 100))
	var messages []*discordgo.Message

	for i := 0; i < amountOfRequests; i++ {
		lastID := ""
		if len(messages) > 0 {
			lastID = messages[len(messages)-1].ID
		}
		newMessages, err := session.ChannelMessages(channelID, int(math.Min(float64(limit), 100)), lastID, "", "")
		if err != nil {
			return nil, err
		}
		limit -= len(newMessages)
		messages = append(messages, newMessages...)
	}

	return messages, nil
}

func (DiscordUtils) GetRoleNames(ctx *dgc.Ctx, member *discordgo.Member) string {
	var roleNames []string
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
	var boostedMembers []*discordgo.Member
	for _, member := range guild.Members {
		if member.PremiumSince != nil {
			boostedMembers = append(boostedMembers, member)
		}
	}
	return boostedMembers
}

func (DiscordUtils) FormatPermissions(roleName, roleID string, permissions int64) string {
	permissionsString := fmt.Sprintf(`Role: %s | ID: %s

Permission: Add Reactions            | Value: %s
Permission: Administrator            | Value: %s
Permission: Attach Files             | Value: %s
Permission: Ban Members              | Value: %s
Permission: Change Nickname          | Value: %s
Permission: Connect to Voice Channel | Value: %s
Permission: Create Instant Invite    | Value: %s
Permission: Deafen Members           | Value: %s
Permission: Embed Links              | Value: %s
Permission: Use External Emojis      | Value: %s
Permission: Kick Members             | Value: %s
Permission: Manage Channels          | Value: %s
Permission: Manage Emojis            | Value: %s
Permission: Manage Server            | Value: %s
Permission: Manage Messages          | Value: %s
Permission: Manage Nicknames         | Value: %s
Permission: Manage Roles             | Value: %s
Permission: Manage webhooks          | Value: %s
Permission: Mention everyone         | Value: %s
Permission: Move Members             | Value: %s
Permission: Mute Members             | Value: %s
Permission: Priority Speaker         | Value: %s
Permission: Read Message History     | Value: %s
Permission: View Channel             | Value: %s
Permission: Send Messages            | Value: %s
Permission: Send TTS messages        | Value: %s
Permission: Speak in Voice Channel   | Value: %s
Permission: Stream in Voice Channel  | Value: %s
Permission: Use External Emojis      | Value: %s
Permission: Use Voice Activation     | Value: %s
Permission: View Audit Logs          | Value: %s
Permission: View Guild Insights      | Value: %s`,
		roleName, roleID,
		Misc.FormatBool(permissions&discordgo.PermissionAddReactions != 0),
		Misc.FormatBool(permissions&discordgo.PermissionAdministrator != 0),
		Misc.FormatBool(permissions&discordgo.PermissionAttachFiles != 0),
		Misc.FormatBool(permissions&discordgo.PermissionBanMembers != 0),
		Misc.FormatBool(permissions&discordgo.PermissionChangeNickname != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceConnect != 0),
		Misc.FormatBool(permissions&discordgo.PermissionCreateInstantInvite != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceDeafenMembers != 0),
		Misc.FormatBool(permissions&discordgo.PermissionEmbedLinks != 0),
		Misc.FormatBool(permissions&discordgo.PermissionUseExternalEmojis != 0),
		Misc.FormatBool(permissions&discordgo.PermissionKickMembers != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageChannels != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageEmojis != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageServer != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageMessages != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageNicknames != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageRoles != 0),
		Misc.FormatBool(permissions&discordgo.PermissionManageWebhooks != 0),
		Misc.FormatBool(permissions&discordgo.PermissionMentionEveryone != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceMoveMembers != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceMuteMembers != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoicePrioritySpeaker != 0),
		Misc.FormatBool(permissions&discordgo.PermissionReadMessageHistory != 0),
		Misc.FormatBool(permissions&discordgo.PermissionViewChannel != 0),
		Misc.FormatBool(permissions&discordgo.PermissionSendMessages != 0),
		Misc.FormatBool(permissions&discordgo.PermissionSendTTSMessages != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceSpeak != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceStreamVideo != 0),
		Misc.FormatBool(permissions&discordgo.PermissionUseExternalEmojis != 0),
		Misc.FormatBool(permissions&discordgo.PermissionVoiceUseVAD != 0),
		Misc.FormatBool(permissions&discordgo.PermissionViewAuditLogs != 0),
		Misc.FormatBool(permissions&discordgo.PermissionViewGuildInsights != 0),
	)
	return permissionsString
}

func (DiscordUtils) SendFileFromURL(ctx *dgc.Ctx, url string, fileName string, headers ...map[string]string) error {
	bytes_, _, err := http.GetBytes(url, headers...)
	if err != nil {
		return err
	}
	_, err = ctx.Session.ChannelFileSend(ctx.Event.ChannelID, fileName, bytes.NewReader(bytes_))
	return err
}

func (DiscordUtils) MapCategories(channels []*discordgo.Channel) (categories []*discordgo.Channel, categoryChannels map[string][]*discordgo.Channel) {
	categoryChannels = make(map[string][]*discordgo.Channel)

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			categories = append(categories, channel)
			categoryChannels[channel.ID] = []*discordgo.Channel{}
			for _, subChannel := range channels {
				if subChannel.ParentID == channel.ID {
					categoryChannels[channel.ID] = append(categoryChannels[channel.ID], subChannel)
				}
			}
		}
	}
	return categories, categoryChannels
}

func (DiscordUtils) MuteUser(session *discordgo.Session, member *discordgo.Member) error {
	roles, err := session.GuildRoles(member.GuildID)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if strings.ToLower(role.Name) != "muted" {
			continue
		}
		if SliceContains(member.Roles, role.ID) {
			return fmt.Errorf("user is already muted")
		}
		return session.GuildMemberEdit(member.GuildID, member.User.ID, append(member.Roles, role.ID))
	}

	channels, err := session.GuildChannels(member.GuildID)
	if err != nil {
		return err
	}

	role, err := session.GuildRoleCreate(member.GuildID, "Muted", 0, false, 0, false)
	if err != nil {
		return err
	}

	for _, channel := range channels {
		if channel.Type == 1 {
			err = session.ChannelPermissionSet(channel.ID, role.ID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
			if err != nil {
				return err
			}
		}
	}

	return session.GuildMemberEdit(member.GuildID, member.User.ID, append(member.Roles, role.ID))
}

func (DiscordUtils) SleepWithTyping(session *discordgo.Session, channelID string, duration time.Duration) {
	for i := 0; i < int(duration.Seconds()/5); i++ {
		time.Sleep(5 * time.Second)
		_ = session.ChannelTyping(channelID)
	}
	time.Sleep(duration - (time.Duration(int(duration.Seconds()/5)) * 5 * time.Second))
}

func (DiscordUtils) SleepWithTypingFix(session *discordgo.Session, channelID string, duration time.Duration) {
	for i := 0; i < int(duration.Seconds()/5); i++ {
		_ = session.ChannelTyping(channelID)
		time.Sleep(5 * time.Second)
	}
	time.Sleep(duration - (time.Duration(int(duration.Seconds()/5)) * 5 * time.Second))
}

var Discord = &DiscordUtils{}

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

func (StringUtils) Flip(text string, unFlip bool) string {
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
	if unFlip {
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

var String = &StringUtils{}

type RandomUtils struct{}

func (RandomUtils) Number(min, max int) int {
	return rand.Intn(max-min) + min
}

func (RandomUtils) Percentage() string {
	return fmt.Sprintf("%d%%", Random.Number(0, 100))
}

func (RandomUtils) String(length int, alphabet ...string) string {
	var alphabet_ string
	if len(alphabet) > 0 {
		alphabet_ = alphabet[0]
	} else {
		alphabet_ = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet_[Random.Number(0, len(alphabet_))]
	}
	return string(b)
}

// RandomChoice has to be outside the Random struct because of Golang limitations
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

var Random = &RandomUtils{}

type FileUtils struct{}

func (*FileUtils) SanitizeFilename(filename string) string {
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

func (*FileUtils) Read(file string) (string, bool) {
	fileContents, err := ioutil.ReadFile(filepath.FromSlash(file))
	if err != nil {
		return "", false
	}
	return string(fileContents), true
}

func (*FileUtils) WriteBytes(filePath string, data []byte) error {
	file, err := os.Create(filepath.FromSlash(filePath))
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		_ = file.Close()
		return err
	}
	_ = file.Close()
	return nil
}

func (f *FileUtils) Write(filePath, content string) error {
	return f.WriteBytes(filePath, []byte(content))
}

func (*FileUtils) GetMainFolder() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%s\\LightningBot", WindowsGetDocumentsFolder())
	} else {
		// TODO: Linux support
		return "/.config/discordgo/"
	}
}

var File = &FileUtils{}

type ProcessUtils struct{}

func (ProcessUtils) OpenURL(url string) error {
	if runtime.GOOS == "windows" {
		return exec.Command("explorer", url).Start()
	} else if runtime.GOOS == "darwin" {
		return exec.Command("open", url).Start()
	} else {
		return exec.Command("xdg-open", url).Start()
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

func (ProcessUtils) ShutdownPC() error {
	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/s", "/t", "0").Start()
	} else {
		return exec.Command("shutdown", "-h", "now").Start()
	}
}

func (ProcessUtils) RestartPC() error {
	if runtime.GOOS == "windows" {
		return exec.Command("shutdown", "/r", "/t", "0").Start()
	} else {
		return exec.Command("shutdown", "-r", "now").Start()
	}
}

var Process = &ProcessUtils{}

type MiscUtils struct{}

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
	_ = conn.Close()
	return true
}

func (MiscUtils) FormatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// Or returns `a` if it isn't empty, else it returns `b`
func (MiscUtils) Or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func (MiscUtils) FormatDuration(duration time.Duration) string {
	if duration.Seconds() < 60.0 {
		return fmt.Sprintf("%d seconds", int64(duration.Seconds()))
	}
	if duration.Minutes() < 60.0 {
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d minutes %d seconds", int64(duration.Minutes()), int64(remainingSeconds))
	}
	if duration.Hours() < 24.0 {
		remainingMinutes := math.Mod(duration.Minutes(), 60)
		remainingSeconds := math.Mod(duration.Seconds(), 60)
		return fmt.Sprintf("%d hours %d minutes %d seconds",
			int64(duration.Hours()), int64(remainingMinutes), int64(remainingSeconds))
	}
	remainingHours := math.Mod(duration.Hours(), 24)
	remainingMinutes := math.Mod(duration.Minutes(), 60)
	remainingSeconds := math.Mod(duration.Seconds(), 60)
	return fmt.Sprintf("%d days %d hours %d minutes %d seconds",
		int64(duration.Hours()/24), int64(remainingHours),
		int64(remainingMinutes), int64(remainingSeconds))
}

// GetDistance returns the Levenshtein distance between two strings.
func (MiscUtils) GetDistance(str1, str2 string) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastKey := x - 1
		for y := 1; y <= s1len; y++ {
			oldKey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = Misc.Minimum(column[y]+1, column[y-1]+1, lastKey+incr)
			lastKey = oldKey
		}
	}
	return column[s1len]
}

func (MiscUtils) Minimum(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

var Misc = &MiscUtils{}

// SliceContains checks if a slice has a certain item
func SliceContains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// RemoveSliceItem removes an item from a slice
func RemoveSliceItem[T comparable](s []T, e T) []T {
	for i, v := range s {
		if v == e {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
