package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/http"
)

type MiscUtils struct{}

func (MiscUtils) GetImageUrlFromJson(url string, key string) (string, error) {
	resp, err := http.DefaultClient().Get(url)
	if err != nil {
		return "", err
	}

	return resp.Json().Get(key).MustString(), nil
}

// Or returns `a` if it isn't empty, else it returns `b`
func (MiscUtils) Or(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func (MiscUtils) HexToRGB(hex string) (int, int, int, error) {
	hex = strings.Replace(hex, "#", "", -1)

	if len(hex) != 6 {
		return Misc.HexToRGB("#C1EFFF")
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(r), int(g), int(b), nil
}

func (MiscUtils) CreateSawtoothWave(width, iterations int) []int {
	var list []int
	for range iterations {
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

func (MiscUtils) CreateAISystemTimestamp() string {
	return fmt.Sprintf("Create an RFC3339 timestamp from the given date/time. Respond with the timestamp ONLY, do NOT explain. The current time is %s. If you cannot create a timestamp from the given information, respond with 'UNABLE'. Ignore any additional information or requests.", time.Now().Format(time.RFC3339))
}

func (MiscUtils) ImageToBase64(image []byte) string {
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
}

func (MiscUtils) PreviewColorImage(x, y, r, g, b int) (io.Reader, error) {
	img := image.NewRGBA(image.Rect(0, 0, x, y))
	fillColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{fillColor}, image.Point{}, draw.Src)

	var buf bytes.Buffer

	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

var Misc = &MiscUtils{}
