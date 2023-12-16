package utils

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/discordgo"
)

type MiscUtils struct{}

func (MiscUtils) GetImageUrlFromJson(url string, jsonQuery string) (string, error) {
	httpResult := http.Get(url)
	if httpResult.Error != nil {
		return "", httpResult.Error
	}

	return httpResult.Json().Get(jsonQuery).MustString(), nil
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

func (MiscUtils) FormatBoolYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func (MiscUtils) FormatBoolOnOff(b bool) string {
	if b {
		return "On"
	}
	return "Off"
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

func (MiscUtils) HexToRGB(hex string) (int, int, int) {
	hex = strings.Replace(hex, "#", "", -1)

	if len(hex) != 6 {
		return Misc.HexToRGB("#C1EFFF")
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 64)
	g, _ := strconv.ParseInt(hex[2:4], 16, 64)
	b, _ := strconv.ParseInt(hex[4:6], 16, 64)
	return int(r), int(g), int(b)
}

func (MiscUtils) CreateSawtoothWave(width, iterations int) []int {
	var list []int
	for i := 0; i < iterations; i++ {
		for j := 0; j <= width; j++ {
			if j == width {
				for k := width - 1; k >= 0; k-- {
					list = append(list, k)
				}
			} else {
				list = append(list, j)
			}
		}
	}
	return list
}

func (MiscUtils) CreateAISystemPing(system string, session *discordgo.Session, message *discordgo.Message) string {
	if system == "" {
		system = fmt.Sprintf("You are a Discord user named %s. A user named %s just mentioned/pinged you. You can mention/ping them back by saying `%s` The following message is what they sent. Respond to them in the style of their message or like a normal Discord user. Use as little punctuation and uppercase letters as possible.", session.State.User.Username, message.Author.Username, message.Author.Mention())
	}
	return system
}

func (MiscUtils) CreateAISystemPingMessages(system string, session *discordgo.Session, messages []*discordgo.Message) []any {
	system = Misc.CreateAISystemPing(system, session, messages[len(messages)-1])

	gptMessages := []any{
		map[string]any{
			"role":    "system",
			"content": system,
		},
	}

	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]

		role := "user"
		if message.Author.ID == session.State.User.ID {
			role = "assistant"
		}

		gptMessages = append(gptMessages, map[string]any{
			"role":    role,
			"content": message.Content,
		})
	}

	return gptMessages
}

func (MiscUtils) CreateAISystemAFK(system string, session *discordgo.Session, message *discordgo.Message) string {
	return fmt.Sprintf("You are a Discord user named %s. A user named %s just messaged you. You can mention/ping them by saying `%s` The following message is what they sent. Respond to them in the style of their message or like a normal Discord user. Use as little punctuation and uppercase letters as possible. Extra instructions: %s", session.State.User.Username, message.Author.Username, message.Author.Mention(), system)
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

// Map applies a function to every element of a slice
func Map[T, U any](s []T, f func(T) U) []U {
	r := make([]U, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// Filter filters a slice based on a function
func Filter[T any](s []T, f func(T) bool) []T {
	r := make([]T, len(s))
	counter := 0
	for i := 0; i < len(s); i++ {
		if f(s[i]) {
			r[counter] = s[i]
			counter++
		}
	}
	return r[:counter]
}
