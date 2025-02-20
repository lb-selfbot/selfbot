package text

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/constants"
	"github.com/LightningDev1/LB-Selfbot/internal/encodings"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/common-nighthawk/go-figure"
)

type TextCommands struct {
	services *services.Services
}

func NewTextCommands(services *services.Services) *TextCommands {
	return &TextCommands{services: services}
}

func (c *TextCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Text",
		Description: "Text manipulation commands",
	}, []*router.Command{
		{
			Name:        "empty",
			Description: "Sends an empty message",
			Handler:     c.empty,
		},
		{
			Name:        "empty",
			Description: "Sends an empty message",
			Handler:     c.empty,
		},
		{
			Name:        "reverse",
			Description: "Reverses the text and sends it",
			Arguments:   "<text>",
			Handler:     c.reverse,
		},
		{
			Name:        "alphabet",
			Description: "Edits the message with every letter in the alphabet",
			Handler:     c.alphabet,
		},
		{
			Name:        "ascii",
			Description: "Sends the text as ASCII art",
			Arguments:   "<text>",
			Handler:     c.ascii,
		},
		{
			Name:        "bold",
			Description: "Sends the text in bold",
			Arguments:   "<text>",
			Handler:     c.bold,
		},
		{
			Name:        "strike",
			Description: "Add strikethrough effect to text",
			Arguments:   "<text>",
			Handler:     c.strike,
		},
		{
			Name:        "italic",
			Description: "Sends the text in italicized font",
			Arguments:   "<text>",
			Handler:     c.italic,
		},
		{
			Name:        "spoiler",
			Description: "Sends the text as a spoiler",
			Arguments:   "<text>",
			Handler:     c.spoiler,
		},
		{
			Name:        "lspoiler",
			Description: "Same as spoiler but spoilers ever letter seperately",
			Arguments:   "<text>",
			Handler:     c.lspoiler,
		},
		{
			Name:        "leet",
			Description: "Sends the text in a language only leet hackers can understand",
			Arguments:   "<text>",
			Handler:     c.leet,
		},
		{
			Name:        "devowel",
			Description: "Sends the text without the vowels",
			Arguments:   "<text>",
			Handler:     c.devowel,
		},
		{
			Name:        "edit",
			Description: "Makes an edit animation with the text",
			Arguments:   "<text>",
			Handler:     c.edit,
		},
		{
			Name:        "emojify",
			Description: "Sends the text as Discord emojis",
			Arguments:   "<text>",
			Handler:     c.emojify,
		},
		{
			Name:        "shrug",
			Description: "Sends a shrug",
			Arguments:   "<text>",
			Handler:     c.shrug,
		},
		{
			Name:        "lenny",
			Description: "Sends a lenny",
			Handler:     c.lenny,
		},
		{
			Name:        "tableflip",
			Description: "Sends a tableflip",
			Handler:     c.tableflip,
		},
		{
			Name:        "unflip",
			Description: "Sends an unflip",
			Handler:     c.unflip,
		},
		{
			Name:        "clap",
			Description: "Replaces spaces in the text into clap emojis",
			Arguments:   "<text>",
			Handler:     c.clap,
		},
		{
			Name:        "zalgo",
			Description: "Sends the text in an almost unreadable font",
			Arguments:   "<text>",
			Handler:     c.zalgo,
		},
		{
			Name:        "unzalgo",
			Description: "Makes zalgo text readable again",
			Arguments:   "<text>",
			Handler:     c.unzalgo,
		},
		{
			Name:        "upper",
			Description: "Sends the text in uppercase",
			Arguments:   "<text>",
			Handler:     c.upper,
		},
		{
			Name:        "lower",
			Description: "Sends the text in lowercase",
			Arguments:   "<text>",
			Handler:     c.lower,
		},
		{
			Name:        "fliptext",
			Description: "Sends the text upside down",
			Arguments:   "<text>",
			Handler:     c.fliptext,
		},
		{
			Name:        "unfliptext",
			Description: "Sends flipped text back to normal",
			Arguments:   "<text>",
			Handler:     c.unfliptext,
		},
		{
			Name:        "wave",
			Description: "Sends text in a wave",
			Arguments:   "<text>",
			Handler:     c.wave,
		},
		{
			Name:        "customwave",
			Description: "Sends text in a wave with custom parameters",
			Arguments:   "<width> <iterations> <text>",
			Handler:     c.customwave,
		},
		{
			Name:        "zwc",
			Description: "Sends text with zero width characters",
			Arguments:   "<text>",
			Handler:     c.zwc,
		},
		{
			Name:        "distance",
			Description: "Calculate the Levenshtein distance between two strings",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.distance,
		},
		{
			Name:        "codeblock",
			Description: "Send a codeblock with syntax highlighting",
			Arguments:   "<language> <code>",
			Handler:     c.codeblock,
		},

		// Animation commands
		{
			Name:        "virus",
			Description: "Send an animated virus message",
			Arguments:   "<virus name>",
			Handler:     c.virus,
		},
		{
			Name:        "100",
			Description: "Count to 100",
			Handler:     c._100,
		},
		{
			Name:        "bomb",
			Description: "Send an animated bomb message",
			Handler:     c.bomb,
		},
		{
			Name:        "fuckyou",
			Description: "Send an animated fuck you message",
			Handler:     c.fuckyou,
		},
		{
			Name:        "fap",
			Description: "Send a fap animation",
			Handler:     c.fap,
		},
		{
			Name:        "countup",
			Description: "Count up to a number",
			Arguments:   "<number>",
			Handler:     c.countup,
		},
		{
			Name:        "countdown",
			Description: "Count down to a number",
			Arguments:   "<number>",
			Handler:     c.countdown,
		},
		{
			Name:        "asciianimation",
			Description: "Send an animated ascii message",
			Arguments:   "<text>",
			Handler:     c.asciianimation,
		},
		{
			Name:        "customanimationhelp",
			Description: "View custom animation help",
			Handler:     c.customanimationhelp,
		},
		{
			Name:        "customanimations",
			Description: "View all custom animations",
			Handler:     c.customanimations,
		},
		{
			Name:        "customanimation",
			Description: "Invoke a custom animation",
			Arguments:   "<name>",
			Handler:     c.customanimation,
		},

		// Encoding commands
		{
			Name:        "encodings",
			Description: "Get a list of all available encodings",
			Handler:     c.encodings,
		},
		{
			Name:        "encode",
			Description: "Encode text using a specified encoding",
			Arguments:   "<encoding> <text>",
			Handler:     c.encode,
		},
		{
			Name:        "decode",
			Description: "Decode text using a specified encoding",
			Arguments:   "<encoding> <text>",
			Handler:     c.decode,
		},

		// Hashing commands
		{
			Name:        "hashes",
			Description: "Get a list of all available hashes",
			Handler:     c.hashes,
		},
		{
			Name:        "hash",
			Description: "Hash text using a specified hash",
			Arguments:   "<hash> <text>",
			Handler:     c.hash,
		},

		// Encryption commands
		{
			Name:        "encryptions",
			Description: "Get a list of all available encryptions",
			Handler:     c.encryptions,
		},
		{
			Name:        "encrypt",
			Description: "Encrypt text using a specified encryption",
			Arguments:   "<encryption> <\"key\"> <\"text\">",
			Handler:     c.encrypt,
		},
		{
			Name:        "decrypt",
			Description: "Decrypt text using a specified encryption",
			Arguments:   "<encryption> <\"key\"> <\"text\">",
			Handler:     c.decrypt,
		},
	})
}

func (c *TextCommands) empty(ctx *router.Ctx) {
	_ = ctx.RespondText("_ _")
}

func (c *TextCommands) reverse(ctx *router.Ctx) {
	_ = ctx.RespondText(utils.String.Reverse(ctx.Arguments.Raw()))
}

func (c *TextCommands) alphabet(ctx *router.Ctx) {
	message, _ := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "a")
	for i := 1; i < len(constants.ALPHABET); i++ {
		time.Sleep(1 * time.Second)
		char := string(constants.ALPHABET[i])
		_, _ = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, char)
	}
}

func (c *TextCommands) ascii(ctx *router.Ctx) {
	ascii := figure.NewFigure(ctx.Arguments.Raw(), "", false).String()
	_ = ctx.RespondText("```" + ascii + "```")
}

func (c *TextCommands) bold(ctx *router.Ctx) {
	_ = ctx.RespondText("**" + ctx.Arguments.Raw() + "**")
}

func (c *TextCommands) strike(ctx *router.Ctx) {
	_ = ctx.RespondText("~~" + ctx.Arguments.Raw() + "~~")
}

func (c *TextCommands) italic(ctx *router.Ctx) {
	_ = ctx.RespondText("*" + ctx.Arguments.Raw() + "*")
}

func (c *TextCommands) spoiler(ctx *router.Ctx) {
	_ = ctx.RespondText("||" + ctx.Arguments.Raw() + "||")
}

func (c *TextCommands) lspoiler(ctx *router.Ctx) {
	finalString := ""
	chars := []rune(ctx.Arguments.Raw())
	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		finalString = finalString + "||" + char + "||"
	}
	_ = ctx.RespondText(finalString)
}

func (c *TextCommands) leet(ctx *router.Ctx) {
	text := ctx.Arguments.Raw()
	leetChars := [][]string{
		{"a", "4"},
		{"A", "4"},
		{"e", "3"},
		{"E", "3"},
		{"i", "1"},
		{"I", "1"},
		{"o", "0"},
		{"O", "0"},
		{"u", "|_|"},
		{"U", "|_|"},
		{"t", "7"},
		{"T", "7"},
		{"g", "6"},
		{"G", "6"},
		{"s", "5"},
		{"S", "5"},
	}
	for _, element := range leetChars {
		text = strings.ReplaceAll(text, element[0], element[1])
	}
	_ = ctx.RespondText(text)
}

func (c *TextCommands) devowel(ctx *router.Ctx) {
	text := ctx.Arguments.Raw()
	vowels := []string{"a", "A", "e", "E", "i", "I", "o", "O", "u", "U"}
	for _, element := range vowels {
		text = strings.ReplaceAll(text, element, "")
	}
	_ = ctx.RespondText(text)
}

func (c *TextCommands) edit(ctx *router.Ctx) {
	fullText := ctx.Arguments.Raw()
	currentText := string(fullText[0])
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, currentText)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	for i := 1; i < len(fullText); i++ {
		time.Sleep(500 * time.Millisecond)
		currentText = currentText + string(fullText[i])
		_, _ = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, currentText)
	}
}

func (c *TextCommands) emojify(ctx *router.Ctx) {
	text := strings.ToLower(ctx.Arguments.Raw())
	emojiChars := [][]string{
		{"a", "🇦"},
		{"b", "🇧-"},
		{"c", "🇨-"},
		{"d", "🇩-"},
		{"e", "🇪-"},
		{"f", "🇫-"},
		{"g", "🇬-"},
		{"h", "🇭-"},
		{"i", "🇮-"},
		{"j", "🇯-"},
		{"k", "🇰-"},
		{"l", "🇱-"},
		{"m", "🇲-"},
		{"n", "🇳-"},
		{"o", "🇴-"},
		{"p", "🇵-"},
		{"q", "🇶-"},
		{"r", "🇷-"},
		{"s", "🇸-"},
		{"t", "🇹-"},
		{"u", "🇺-"},
		{"v", "🇻"},
		{"w", "🇼-"},
		{"x", "🇽-"},
		{"y", "🇾-"},
		{"z", "🇿-"},
		{"-", " "},
	}
	for _, element := range emojiChars {
		text = strings.ReplaceAll(text, element[0], element[1])
	}
	_ = ctx.RespondText(text)
}

func (c *TextCommands) shrug(ctx *router.Ctx) {
	_ = ctx.RespondText("¯\\_(ツ)_/¯")
}

func (c *TextCommands) lenny(ctx *router.Ctx) {
	_ = ctx.RespondText("( ͡° ͜ʖ ͡°)")
}

func (c *TextCommands) tableflip(ctx *router.Ctx) {
	_ = ctx.RespondText("(╯°□°）╯︵ ┻━┻")
}

func (c *TextCommands) unflip(ctx *router.Ctx) {
	_ = ctx.RespondText("┬─┬ ノ( ゜-゜ノ)")
}

func (c *TextCommands) clap(ctx *router.Ctx) {
	text := ctx.Arguments.Raw()
	clapChars := [][]string{
		{" ", "👏"},
		{"\n", "👏\n"},
	}
	for _, element := range clapChars {
		text = strings.ReplaceAll(text, element[0], element[1])
	}
	_ = ctx.RespondText(text)
}

func (c *TextCommands) zalgo(ctx *router.Ctx) {
	text := strings.ToLower(ctx.Arguments.Raw())
	zalgoChars := [][]string{
		{"a", "a̷̧̡̡̨͉͍̫͖͊̆̂"},
		{"b", "b̸̙͖̹͉͒̉͐̆̎͋͛̆"},
		{"c", "c̶̨̖̙̜̖̣͚̿͗̓͗̈́̏̒͝"},
		{"d", "d̷̩̣̟͈̟̙̀"},
		{"e", "e̶̢̟͛̿̿͊̆͝͝"},
		{"f", "f̵̝̫͋̓"},
		{"g", "g̷̨̩̫̫͎͉̙͔͙̎̊͒́͆͌͠"},
		{"h", "ḩ̴̻̪̼͐̄̀̐͐̄̑̚ͅ"},
		{"i", "ì̴̛̗̟̫̺̊͋̒̂̿̾"},
		{"j", "j̷̛̛̖̗̭̠̬͕̗̺̝̈́́͒̅̆́̐̃͜"},
		{"k", "ḱ̸̬͙̟͇̺͉̲̯̍̋̀̌̏́̈̆́ͅ"},
		{"l", "l̴͇̝̥̜͙̙̄̆̍ͅ"},
		{"m", "m̶̰͔̖͍͂"},
		{"n", "n̷͉͍̺̮̿̓̀"},
		{"o", "ọ̴̻̓̈͐͂͛̐͝͠"},
		{"p", "p̴͚̥͖̙̒̓̀"},
		{"q", "q̷͋̽͒͒̅̓̕"},
		{"r", "̘͎̖̖̦͚͈̓̾̑ṛ̸̢̡̧̜̟͇͖͛͂̆̽̐͌́"},
		{"s", "s̵̖̖̖̹̩͌̂̕͜"},
		{"t", "t̴̞͎͛͗͆̆͗̑̔"},
		{"u", "ǘ̴̞̃̏̓̄͠͝"},
		{"v", "v̸̥̻̩͕͖̂"},
		{"w", "w̵̢̛͇̪͓̄́̆̓̀͛̍͝͝"},
		{"x", "x̸̢̺̣̰͖͓̀̔̊"},
		{"y", "y̸̡̨̘̤͙͖͔̜̳͗̏̀͗͋͛͠"},
		{"z", "z̴͔̳̟̝̤̗̉"},
		{"-", " "},
	}
	for _, element := range zalgoChars {
		text = strings.ReplaceAll(text, element[0], element[1])
	}
	_ = ctx.RespondText(text)
}

func (c *TextCommands) unzalgo(ctx *router.Ctx) {
	text := utils.String.FilterAlphaNumeric(ctx.Arguments.Raw())
	_ = ctx.RespondText(text)
}

func (c *TextCommands) upper(ctx *router.Ctx) {
	text := strings.ToUpper(ctx.Arguments.Raw())
	_ = ctx.RespondText(text)
}

func (c *TextCommands) lower(ctx *router.Ctx) {
	text := strings.ToLower(ctx.Arguments.Raw())
	_ = ctx.RespondText(text)
}

func (c *TextCommands) fliptext(ctx *router.Ctx) {
	text := utils.String.Flip(ctx.Arguments.Raw(), false)
	_ = ctx.RespondText(text)
}

func (c *TextCommands) unfliptext(ctx *router.Ctx) {
	text := utils.String.Flip(ctx.Arguments.Raw(), true)
	_ = ctx.RespondText(text)
}

func (c *TextCommands) wave(ctx *router.Ctx) {
	sawtoothWave := utils.Misc.CreateSawtoothWave(8, 6)

	wave := ""

	for _, i := range sawtoothWave {
		wave += strings.Repeat(" ", i) + ctx.Arguments.Raw() + "\n"
	}

	_ = ctx.RespondText(wave)
}

func (c *TextCommands) customwave(ctx *router.Ctx) {
	width, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		_ = ctx.RespondText("Width must be a number!")
		return
	}

	if width <= 0 {
		_ = ctx.RespondText("Width must be a positive number!")
		return
	}

	iterations, err := ctx.Arguments.Get(1).Int()
	if err != nil {
		_ = ctx.RespondText("Iterations must be a number!")
	}

	if iterations <= 0 {
		_ = ctx.RespondText("Iterations must be a positive number!")
		return
	}

	sawtoothWave := utils.Misc.CreateSawtoothWave(width, iterations)

	text := ctx.Arguments.RawFrom(2)

	wave := ""

	for _, i := range sawtoothWave {
		wave += strings.Repeat(" ", i) + text + "\n"
	}

	_ = ctx.RespondText(wave)
}

func (c *TextCommands) zwc(ctx *router.Ctx) {
	text := ctx.Arguments.Raw()
	text = strings.Join(strings.Split(text, ""), "\u200B")
	_ = ctx.RespondText(text)
}

func (c *TextCommands) distance(ctx *router.Ctx) {
	text1, text2 := ctx.Arguments.Get(0).Raw(), ctx.Arguments.Get(1).Raw()

	c.services.Embed.NewEmbed().
		AddField("Text 1", text1).
		AddField("Text 2", text2).
		AddField("Distance", strconv.Itoa(utils.Misc.GetDistance(text1, text2))).
		Send(ctx)
}

func (c *TextCommands) codeblock(ctx *router.Ctx) {
	language := ctx.Arguments.Get(0).Raw()
	code := ctx.Arguments.RawFrom(1)

	_ = ctx.RespondText("```" + strings.ToLower(language) + "\n" + code + "\n```")
}

// Animation commands
var (
	virus = []string{
		"`[▓▓▓                    ] / %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓                ] - %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓           ] \\ %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓         ] | %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ] / %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓   ] - %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ] \\ %s.exe - Compiling`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] | %s.exe - Compiled`",
		"`[▓▓▓                    ] / %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓           ] - %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ] \\ %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ ] | %s.exe - Downloading`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] / %s.exe - Downloaded`",
		"`[▓▓▓▓▓▓▓                ] - %s.exe - Injecting`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓      ] \\ %s.exe - Injecting`",
		"`[▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓] | %s.exe - Injected`",
		"`Successfully Injected %s.exe`",
	}
	bomb = []string{
		":bomb: ---------------- :fire:",
		":bomb: --------------- :fire:",
		":bomb: -------------- :fire:",
		":bomb: ------------- :fire:",
		":bomb: ------------ :fire:",
		":bomb: ----------- :fire:",
		":bomb: ---------- :fire:",
		":bomb: --------- :fire:",
		":bomb: -------- :fire:",
		":bomb: ------- :fire:",
		":bomb: ------ :fire:",
		":bomb: ----- :fire:",
		":bomb: ---- :fire:",
		":bomb: --- :fire:",
		":bomb: -- :fire:",
		":bomb: - :fire:",
		":bomb:  :fire:",
		":boom:",
	}
	fuckYou = []string{
		"**F**",
		"**FU**",
		"**FUC**",
		"**FUCK**",
		"**FUCK Y**",
		"**FUCK YO**",
		"**FUCK YOU**",
		"**FUCK YOU!**",
	}
	fap = []string{
		"8:punch:=====Đ",
		"8=:punch:====Đ",
		"8==:punch:===Đ",
		"8===:punch:==Đ",
		"8====:punch:=Đ",
		"8=====:punch:Đ",
		"8====:punch:=Đ",
		"8===:punch:==Đ",
		"8==:punch:===Đ",
		"8=:punch:====Đ",
		"8:punch:=====Đ",
		"8=:punch:====Đ",
		"8==:punch:===Đ",
		"8===:punch:==Đ",
		"8====:punch:=Đ",
		"8=====:punch:Đ",
		"8=====Đ:sweat_drops:",
	}
)

func (c *TextCommands) virus(ctx *router.Ctx) {
	name := ctx.Arguments.Get(0).Raw()

	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(virus[0], name))
	if err != nil {
		return
	}

	for i := 1; i < len(virus); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, fmt.Sprintf(virus[i], name))
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) _100(ctx *router.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "1")
	if err != nil {
		return
	}

	for i := 2; i <= 100; i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, strconv.Itoa(i))
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) bomb(ctx *router.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, bomb[0])
	if err != nil {
		return
	}

	for i := 1; i < len(bomb); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, bomb[i])
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) fuckyou(ctx *router.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fuckYou[0])
	if err != nil {
		return
	}

	for i := 1; i < len(fuckYou); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, fuckYou[i])
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) fap(ctx *router.Ctx) {
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fap[0])
	if err != nil {
		return
	}

	for i := 1; i < len(fap); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, fap[i])
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) countup(ctx *router.Ctx) {
	number, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		_ = ctx.RespondText("Argument <number> must be a number.")
		return
	}

	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "1")
	if err != nil {
		return
	}

	for i := 2; i <= number; i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, strconv.Itoa(i))
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) countdown(ctx *router.Ctx) {
	number, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		_ = ctx.RespondText("Argument <number> must be a number.")
		return
	}

	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprint(number))
	if err != nil {
		return
	}

	for i := number; i > 0; i-- {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, fmt.Sprint(i-1))
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) asciianimation(ctx *router.Ctx) {
	text := ctx.Arguments.Get(0).Raw()
	firstCharAscii := figure.NewFigure(string(text[0]), "", false).String()
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("```%s```", firstCharAscii))
	if err != nil {
		return
	}
	currentText := string(text[0])
	for index, char := range text {
		if index == 0 {
			continue
		}
		time.Sleep(time.Second * 1)
		currentText = currentText + string(char)
		ascii := figure.NewFigure(currentText, "", false).String()
		_, err := ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, fmt.Sprintf("```%s```", ascii))
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) customanimationhelp(ctx *router.Ctx) {
	c.services.Embed.NewEmbed().
		AddField("Creating an animation", fmt.Sprintf(
			"To create an animation, create a text file in %s. Every line in the file will be a frame of the animation.\n",
			filepath.Join(utils.IO.GetMainFolder(), "Animations"),
		)).
		AddField("Using an animation", fmt.Sprintf(
			"To use an animation, use the command `%scustomanimation <filename without .txt>`.\n",
			ctx.Router.GetPrefixes()[0],
		)).
		Send(ctx)
}

func (c *TextCommands) customanimations(ctx *router.Ctx) {
	folder := filepath.Join(utils.IO.GetMainFolder(), "Animations")
	_ = os.MkdirAll(folder, os.ModePerm)

	customAnimString := ""
	_ = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		nameSplit := strings.Split(info.Name(), ".")
		if len(nameSplit) != 2 {
			return nil
		}
		if nameSplit[1] != "txt" {
			return nil
		}
		customAnimString += " - " + nameSplit[0] + "\n"
		return nil
	})

	c.services.Embed.NewEmbed().AddField("Custom Animations", customAnimString).Send(ctx)
}

func (c *TextCommands) customanimation(ctx *router.Ctx) {
	name := utils.Arguments["<name>"]
	folder := filepath.Join(utils.IO.GetMainFolder(), "Animations")
	_ = os.MkdirAll(folder, os.ModePerm)
	path := filepath.Join(folder, utils.IO.SanitizeFilename(name+".txt"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = ctx.RespondText("Animation does not exist.")
		return
	}
	rawAnimation, ok := utils.IO.Read(path)
	if !ok {
		_ = ctx.RespondText("Animation could not be read.")
		return
	}
	animationFrames := strings.Split(rawAnimation, "\n")
	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, animationFrames[0])
	if err != nil {
		return
	}
	for i := 1; i < len(animationFrames); i++ {
		time.Sleep(time.Second * 1)
		message, err = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, message.ID, animationFrames[i])
		if err != nil {
			return
		}
	}
}

func (c *TextCommands) encodings(ctx *router.Ctx) {
	var encodingNames []string
	for name := range encodings.Encodings {
		encodingNames = append(encodingNames, name)
	}

	c.services.Embed.NewEmbed().
		AddField("Available Encodings", strings.Join(encodingNames, ", ")).
		Send(ctx)
}

func (c *TextCommands) encode(ctx *router.Ctx) {
	encodingName := ctx.Arguments.Get(0).Raw()
	text := ctx.Arguments.RawFrom(1)

	encoding, ok := encodings.Encodings[encodingName]
	if !ok {
		_ = ctx.RespondText("Invalid encoding! Use `[p]encodings` to get a list of all available encodings.")
		return
	}

	encoded, err := encoding.Encode(text)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Text", text).
		AddField("Encoded", encoded).
		Send(ctx)
}

func (c *TextCommands) decode(ctx *router.Ctx) {
	encodingName := ctx.Arguments.Get(0).Raw()
	text := ctx.Arguments.RawFrom(1)

	encoding, ok := encodings.Encodings[encodingName]
	if !ok {
		_ = ctx.RespondText("Invalid encoding! Use `[p]encodings` to get a list of all available encodings.")
		return
	}

	if encoding.Decode == nil {
		encoding.Decode = encoding.Encode
	}

	decoded, err := encoding.Decode(text)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Text", text).
		AddField("Decoded", decoded).
		Send(ctx)
}

func (c *TextCommands) hashes(ctx *router.Ctx) {
	var hashNames []string
	for name := range encodings.Hashes {
		hashNames = append(hashNames, name)
	}

	c.services.Embed.NewEmbed().
		AddField("Available Hashes", strings.Join(hashNames, ", ")).
		Send(ctx)
}

func (c *TextCommands) hash(ctx *router.Ctx) {
	hashName := ctx.Arguments.Get(0).Raw()
	text := ctx.Arguments.RawFrom(1)

	hashFunc, ok := encodings.Hashes[hashName]
	if !ok {
		_ = ctx.RespondText("Invalid hash! Use `[p]hashes` to get a list of all available hashes.")
		return
	}

	hashed := hashFunc(text)

	c.services.Embed.NewEmbed().
		AddField("Text", text).
		AddField("Hashed", hashed).
		Send(ctx)
}

func (c *TextCommands) encryptions(ctx *router.Ctx) {
	var encryptionNames []string
	for name := range encodings.Encryptions {
		encryptionNames = append(encryptionNames, name)
	}

	c.services.Embed.NewEmbed().
		AddField("Available Encryptions", strings.Join(encryptionNames, ", ")).
		Send(ctx)
}

func (c *TextCommands) encrypt(ctx *router.Ctx) {
	encryptionName := ctx.Arguments.Get(0).Raw()
	key := utils.Arguments["<\"key\">"]
	text := utils.Arguments["<\"text\">"]

	encryption, ok := encodings.Encryptions[encryptionName]
	if !ok {
		_ = ctx.RespondText("Invalid encryption! Use `[p]encryptions` to get a list of all available encryptions.")
		return
	}

	encrypted, err := encryption.Encrypt(text, key)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Text", text).
		AddField("Encrypted", encrypted).
		Send(ctx)
}

func (c *TextCommands) decrypt(ctx *router.Ctx) {
	encryptionName := ctx.Arguments.Get(0).Raw()
	key := utils.Arguments["<\"key\">"]
	text := utils.Arguments["<\"text\">"]

	encryption, ok := encodings.Encryptions[encryptionName]
	if !ok {
		_ = ctx.RespondText("Invalid encryption! Use `[p]encryptions` to get a list of all available encryptions.")
		return
	}

	encrypted, err := encryption.Decrypt(text, key)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Text", text).
		AddField("Decrypted", encrypted).
		Send(ctx)
}
