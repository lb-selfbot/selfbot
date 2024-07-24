package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/encodings"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/common-nighthawk/go-figure"
)

type TextCommands struct{}

func (t TextCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Text", "Text commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "empty",
		Description: "Sends an empty message",
		Usage:       "[p]empty",
		Handler:     t.empty,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reverse",
		Description: "Reverses the text and sends it",
		Usage:       "[p]reverse <text>",
		Handler:     t.reverse,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "alphabet",
		Description: "Edits the message with every letter in the alphabet",
		Usage:       "[p]alphabet",
		Handler:     t.alphabet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ascii",
		Description: "Sends the text as ASCII art",
		Usage:       "[p]ascii <text>",
		Handler:     t.ascii,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bold",
		Description: "Sends the text in bold",
		Usage:       "[p]bold <text>",
		Handler:     t.bold,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "strike",
		Description: "Add strikethrough effect to text",
		Usage:       "[p]strike <text>",
		Handler:     t.strike,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "italic",
		Description: "Sends the text in italicized font",
		Usage:       "[p]italic <text>",
		Handler:     t.italic,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "spoiler",
		Description: "Sends the text as a spoiler",
		Usage:       "[p]spoiler <text>",
		Handler:     t.spoiler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lspoiler",
		Description: "Same as spoiler but spoilers ever letter seperately",
		Usage:       "[p]lspoiler <text>",
		Handler:     t.lspoiler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leet",
		Description: "Sends the text in a language only leet hackers can understand",
		Usage:       "[p]leet <text>",
		Handler:     t.leet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "devowel",
		Description: "Sends the text without the vowels",
		Usage:       "[p]devowel <text>",
		Handler:     t.devowel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "edit",
		Description: "Makes an edit animation with the text",
		Usage:       "[p]edit <text>",
		Handler:     t.edit,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "emojify",
		Description: "Sends the text as Discord emojis",
		Usage:       "[p]emojify <text>",
		Handler:     t.emojify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shrug",
		Description: "Sends a shrug",
		Usage:       "[p]shrug <text>",
		Handler:     t.shrug,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lenny",
		Description: "Sends a lenny",
		Usage:       "[p]lenny",
		Handler:     t.lenny,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tableflip",
		Description: "Sends a tableflip",
		Usage:       "[p]tableflip",
		Handler:     t.tableflip,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unflip",
		Description: "Sends an unflip",
		Usage:       "[p]unflip",
		Handler:     t.unflip,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clap",
		Description: "Replaces spaces in the text into clap emojis",
		Usage:       "[p]clap <text>",
		Handler:     t.clap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "zalgo",
		Description: "Sends the text in an almost unreadable font",
		Usage:       "[p]zalgo <text>",
		Handler:     t.zalgo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unzalgo",
		Description: "Makes zalgo text readable again",
		Usage:       "[p]unzalgo <text>",
		Handler:     t.unzalgo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "upper",
		Description: "Sends the text in uppercase",
		Usage:       "[p]upper <text>",
		Handler:     t.upper,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lower",
		Description: "Sends the text in lowercase",
		Usage:       "[p]lower <text>",
		Handler:     t.lower,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fliptext",
		Description: "Sends the text upside down",
		Usage:       "[p]fliptext <text>",
		Handler:     t.fliptext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unfliptext",
		Description: "Sends flipped text back to normal",
		Usage:       "[p]unfliptext <text>",
		Handler:     t.unfliptext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wave",
		Description: "Sends text in a wave",
		Usage:       "[p]wave <text>",
		Handler:     t.wave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "customwave",
		Description: "Sends text in a wave with custom parameters",
		Usage:       "[p]customwave <width> <iterations> <text>",
		Handler:     t.customwave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "zwc",
		Description: "Sends text with zero width characters",
		Usage:       "[p]zwc <text>",
		Handler:     t.zwc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "distance",
		Description: "Calculate the Levenshtein distance between two strings",
		Usage:       "[p]distance <\"text1\"> <\"text2\">",
		Handler:     t.distance,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "codeblock",
		Description: "Send a codeblock with syntax highlighting",
		Usage:       "[p]codeblock <language> <code>",
		Handler:     t.codeblock,
	})

	// Encoding commands

	router.RegisterCmd(&dgc.Command{
		Name:        "encodings",
		Description: "Get a list of all available encodings",
		Usage:       "[p]encodings",
		Handler:     t.encodings,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "encode",
		Description: "Encode text using a specified encoding",
		Usage:       "[p]encode <encoding> <text>",
		Handler:     t.encode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "decode",
		Description: "Decode text using a specified encoding",
		Usage:       "[p]decode <encoding> <text>",
		Handler:     t.decode,
	})

	// Hashing commands

	router.RegisterCmd(&dgc.Command{
		Name:        "hashes",
		Description: "Get a list of all available hashes",
		Usage:       "[p]hashes",
		Handler:     t.hashes,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hash",
		Description: "Hash text using a specified hash",
		Usage:       "[p]hash <hash> <text>",
		Handler:     t.hash,
	})

	// Encryption commands

	router.RegisterCmd(&dgc.Command{
		Name:        "encryptions",
		Description: "Get a list of all available encryptions",
		Usage:       "[p]encryptions",
		Handler:     t.encryptions,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "encrypt",
		Description: "Encrypt text using a specified encryption",
		Usage:       "[p]encrypt <encryption> <\"key\"> <\"text\">",
		Handler:     t.encrypt,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "decrypt",
		Description: "Decrypt text using a specified encryption",
		Usage:       "[p]decrypt <encryption> <\"key\"> <\"text\">",
		Handler:     t.decrypt,
	})

}

func (TextCommands) empty(ctx *dgc.Ctx) {
	_ = ctx.RespondText("_ _")
}

func (TextCommands) reverse(ctx *dgc.Ctx) {
	_ = ctx.RespondText(utils.String.Reverse(ctx.Arguments.Raw()))
}

func (TextCommands) alphabet(ctx *dgc.Ctx) {
	message, _ := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, "a")
	for i := 1; i < len(constants.ALPHABET); i++ {
		time.Sleep(1 * time.Second)
		char := string(constants.ALPHABET[i])
		_, _ = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, char)
	}
}

func (TextCommands) ascii(ctx *dgc.Ctx) {
	ascii := figure.NewFigure(ctx.Arguments.Raw(), "", false).String()
	_ = ctx.RespondText("```" + ascii + "```")
}

func (TextCommands) bold(ctx *dgc.Ctx) {
	_ = ctx.RespondText("**" + ctx.Arguments.Raw() + "**")
}

func (TextCommands) strike(ctx *dgc.Ctx) {
	_ = ctx.RespondText("~~" + ctx.Arguments.Raw() + "~~")
}

func (TextCommands) italic(ctx *dgc.Ctx) {
	_ = ctx.RespondText("*" + ctx.Arguments.Raw() + "*")
}

func (TextCommands) spoiler(ctx *dgc.Ctx) {
	_ = ctx.RespondText("||" + ctx.Arguments.Raw() + "||")
}

func (TextCommands) lspoiler(ctx *dgc.Ctx) {
	finalString := ""
	chars := []rune(ctx.Arguments.Raw())
	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		finalString = finalString + "||" + char + "||"
	}
	_ = ctx.RespondText(finalString)
}

func (TextCommands) leet(ctx *dgc.Ctx) {
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

func (TextCommands) devowel(ctx *dgc.Ctx) {
	text := ctx.Arguments.Raw()
	vowels := []string{"a", "A", "e", "E", "i", "I", "o", "O", "u", "U"}
	for _, element := range vowels {
		text = strings.ReplaceAll(text, element, "")
	}
	_ = ctx.RespondText(text)
}

func (TextCommands) edit(ctx *dgc.Ctx) {
	fullText := ctx.Arguments.Raw()
	currentText := string(fullText[0])
	message, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, currentText)
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

func (TextCommands) emojify(ctx *dgc.Ctx) {
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

func (TextCommands) shrug(ctx *dgc.Ctx) {
	_ = ctx.RespondText("¯\\_(ツ)_/¯")
}

func (TextCommands) lenny(ctx *dgc.Ctx) {
	_ = ctx.RespondText("( ͡° ͜ʖ ͡°)")
}

func (TextCommands) tableflip(ctx *dgc.Ctx) {
	_ = ctx.RespondText("(╯°□°）╯︵ ┻━┻")
}

func (TextCommands) unflip(ctx *dgc.Ctx) {
	_ = ctx.RespondText("┬─┬ ノ( ゜-゜ノ)")
}

func (TextCommands) clap(ctx *dgc.Ctx) {
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

func (TextCommands) zalgo(ctx *dgc.Ctx) {
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

func (TextCommands) unzalgo(ctx *dgc.Ctx) {
	text := utils.String.FilterAlphaNumeric(ctx.Arguments.Raw())
	_ = ctx.RespondText(text)
}

func (TextCommands) upper(ctx *dgc.Ctx) {
	text := strings.ToUpper(ctx.Arguments.Raw())
	_ = ctx.RespondText(text)
}

func (TextCommands) lower(ctx *dgc.Ctx) {
	text := strings.ToLower(ctx.Arguments.Raw())
	_ = ctx.RespondText(text)
}

func (TextCommands) fliptext(ctx *dgc.Ctx) {
	text := utils.String.Flip(ctx.Arguments.Raw(), false)
	_ = ctx.RespondText(text)
}

func (TextCommands) unfliptext(ctx *dgc.Ctx) {
	text := utils.String.Flip(ctx.Arguments.Raw(), true)
	_ = ctx.RespondText(text)
}

func (TextCommands) wave(ctx *dgc.Ctx) {
	sawtoothWave := utils.Misc.CreateSawtoothWave(8, 6)

	wave := ""

	for _, i := range sawtoothWave {
		wave += strings.Repeat(" ", i) + ctx.Arguments.Raw() + "\n"
	}

	_ = ctx.RespondText(wave)
}

func (TextCommands) customwave(ctx *dgc.Ctx) {
	width, err := ctx.Arguments.Get(0).AsInt()
	if err != nil {
		_ = ctx.RespondText("Width must be a number!")
		return
	}

	if width <= 0 {
		_ = ctx.RespondText("Width must be a positive number!")
		return
	}

	iterations, err := ctx.Arguments.Get(1).AsInt()
	if err != nil {
		_ = ctx.RespondText("Iterations must be a number!")
	}

	if iterations <= 0 {
		_ = ctx.RespondText("Iterations must be a positive number!")
		return
	}

	sawtoothWave := utils.Misc.CreateSawtoothWave(width, iterations)

	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 2)

	wave := ""

	for _, i := range sawtoothWave {
		wave += strings.Repeat(" ", i) + text + "\n"
	}

	_ = ctx.RespondText(wave)
}

func (TextCommands) zwc(ctx *dgc.Ctx) {
	text := ctx.Arguments.Raw()
	text = strings.Join(strings.Split(text, ""), "\u200B")
	_ = ctx.RespondText(text)
}

func (TextCommands) distance(ctx *dgc.Ctx) {
	distance := utils.Misc.GetDistance(
		utils.Arguments["<\"text1\">"],
		utils.Arguments["<\"text2\">"],
	)

	embed.NewEmbed().
		AddField("Text 1", utils.Arguments["<\"text1\">"]).
		AddField("Text 2", utils.Arguments["<\"text2\">"]).
		AddField("Distance", fmt.Sprint(distance)).
		Send(ctx)
}

func (TextCommands) codeblock(ctx *dgc.Ctx) {
	language := utils.Arguments["<language>"]
	code := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	_ = ctx.RespondText("```" + strings.ToLower(language) + "\n" + code + "\n```")
}

func (TextCommands) encodings(ctx *dgc.Ctx) {
	var encodingNames []string
	for name := range encodings.Encodings {
		encodingNames = append(encodingNames, name)
	}

	embed.NewEmbed().
		AddField("Available Encodings", strings.Join(encodingNames, ", ")).
		Send(ctx)
}

func (TextCommands) encode(ctx *dgc.Ctx) {
	encodingName := ctx.Arguments.Get(0).Raw()
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

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

	embed.NewEmbed().
		AddField("Text", text).
		AddField("Encoded", encoded).
		Send(ctx)
}

func (TextCommands) decode(ctx *dgc.Ctx) {
	encodingName := ctx.Arguments.Get(0).Raw()
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

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

	embed.NewEmbed().
		AddField("Text", text).
		AddField("Decoded", decoded).
		Send(ctx)
}

func (TextCommands) hashes(ctx *dgc.Ctx) {
	var hashNames []string
	for name := range encodings.Hashes {
		hashNames = append(hashNames, name)
	}

	embed.NewEmbed().
		AddField("Available Hashes", strings.Join(hashNames, ", ")).
		Send(ctx)
}

func (TextCommands) hash(ctx *dgc.Ctx) {
	hashName := ctx.Arguments.Get(0).Raw()
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	hashFunc, ok := encodings.Hashes[hashName]
	if !ok {
		_ = ctx.RespondText("Invalid hash! Use `[p]hashes` to get a list of all available hashes.")
		return
	}

	hashed := hashFunc(text)

	embed.NewEmbed().
		AddField("Text", text).
		AddField("Hashed", hashed).
		Send(ctx)
}

func (TextCommands) encryptions(ctx *dgc.Ctx) {
	var encryptionNames []string
	for name := range encodings.Encryptions {
		encryptionNames = append(encryptionNames, name)
	}

	embed.NewEmbed().
		AddField("Available Encryptions", strings.Join(encryptionNames, ", ")).
		Send(ctx)
}

func (TextCommands) encrypt(ctx *dgc.Ctx) {
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

	embed.NewEmbed().
		AddField("Text", text).
		AddField("Encrypted", encrypted).
		Send(ctx)
}

func (TextCommands) decrypt(ctx *dgc.Ctx) {
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

	embed.NewEmbed().
		AddField("Text", text).
		AddField("Decrypted", encrypted).
		Send(ctx)
}
