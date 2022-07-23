package commands

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/json"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v6"
	owoify_go "github.com/deadshot465/owoify-go"
)

type FunCommands struct{}

func (f FunCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "eightball",
		Description: "Asks the 8-Ball a question",
		Usage:       "[p]eightball <question>",
		Handler:     f.eightball,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "howgay",
		Description: "See how gay a person is",
		Usage:       "[p]howgay <user>",
		Handler:     f.howgay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "howsimp",
		Description: "See how simp a person is",
		Usage:       "[p]howsimp <user>",
		Handler:     f.howsimp,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "penis",
		Description: "See how big a person's penis is",
		Usage:       "[p]penis <user>",
		Handler:     f.penis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ego",
		Description: "See how big a person's ego is",
		Usage:       "[p]ego <user>",
		Handler:     f.ego,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hack",
		Description: "Fake hack a person",
		Usage:       "[p]hack <user>",
		Handler:     f.hack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "insult",
		Description: "Insult a person",
		Usage:       "[p]insult <user>",
		Handler:     f.insult,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dadjoke",
		Description: "Send a random dad joke",
		Usage:       "[p]dadjoke",
		Handler:     f.dadjoke,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "yomomma",
		Description: "Send a random yo momma joke",
		Usage:       "[p]yomomma",
		Handler:     f.yomomma,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wouldyourather",
		Description: "Send a random would you rather question",
		Usage:       "[p]wouldyourather",
		Handler:     f.wouldyourather,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "topic",
		Description: "Send a conversation topic",
		Usage:       "[p]topic",
		Handler:     f.topic,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "slot",
		Description: "Play the slot machine",
		Usage:       "[p]slot",
		Handler:     f.slot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tts",
		Description: "Send a Text-To-Speech message",
		Usage:       "[p]tts <text>",
		Handler:     f.tts,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "advice",
		Description: "Send random advice",
		Usage:       "[p]advice",
		Handler:     f.advice,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catfact",
		Description: "Send a random cat fact",
		Usage:       "[p]catfact",
		Handler:     f.catfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dogfact",
		Description: "Send a random dog fact",
		Usage:       "[p]dogfact",
		Handler:     f.dogfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pandafact",
		Description: "Send a random panda fact",
		Usage:       "[p]pandafact",
		Handler:     f.pandafact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "foxfact",
		Description: "Send a random fox fact",
		Usage:       "[p]foxfact",
		Handler:     f.foxfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "koalafact",
		Description: "Send a random koala fact",
		Usage:       "[p]koalafact",
		Handler:     f.koalafact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "birdfact",
		Description: "Send a random bird fact",
		Usage:       "[p]birdfact",
		Handler:     f.birdfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "redpandafact",
		Description: "Send a random red panda fact",
		Usage:       "[p]redpandafact",
		Handler:     f.redpandafact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "whalefact",
		Description: "Send a random whale fact",
		Usage:       "[p]whalefact",
		Handler:     f.whalefact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kangaroofact",
		Description: "Send a random whale fact",
		Usage:       "[p]kangaroofact",
		Handler:     f.kangaroofact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bunnyfact",
		Description: "Send a random bunny fact",
		Usage:       "[p]bunnyfact",
		Handler:     f.bunnyfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lionfact",
		Description: "Send a random lion fact",
		Usage:       "[p]lionfact",
		Handler:     f.lionfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "frogfact",
		Description: "Send a random frog fact",
		Usage:       "[p]frogfact",
		Handler:     f.frogfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "duckfact",
		Description: "Send a random duck fact",
		Usage:       "[p]duckfact",
		Handler:     f.duckfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "penguinfact",
		Description: "Send a random penguin fact",
		Usage:       "[p]penguinfact",
		Handler:     f.penguinfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trivia",
		Description: "Send a random trivia question",
		Usage:       "[p]trivia",
		Handler:     f.trivia,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dice",
		Description: "Roll a dice",
		Usage:       "[p]dice",
		Handler:     f.dice,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomemoji",
		Description: "Send a random emoji from ther current server [server only]",
		Usage:       "[p]randomemoji",
		Handler:     f.randomemoji,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomchannel",
		Description: "Send a random channel from ther current server [server only]",
		Usage:       "[p]randomchannel",
		Handler:     f.randomchannel,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "choose",
		Description: "Choose a random item from the list of items",
		Usage:       "[p]choose <option1> <option2> ...",
		Handler:     f.choose,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "poll",
		Description: "Make a poll",
		Usage:       "[p]poll <question>",
		Handler:     f.poll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "numberfact",
		Description: "Send a random fact about a number",
		Usage:       "[p]numberfact <number>",
		Handler:     f.numberfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "joke",
		Description: "Send a random joke",
		Usage:       "[p]joke",
		Handler:     f.joke,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fml",
		Description: "Send a random fuck my life quote",
		Usage:       "[p]fml",
		Handler:     f.fml,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tronald",
		Description: "Send a random Donald Trump quote",
		Usage:       "[p]tronald",
		Handler:     f.tronald,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kanyewest",
		Description: "Send a random Kanye West quote",
		Usage:       "[p]kanyewest",
		Handler:     f.kanyewest,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rps",
		Description: "Play rock paper scissors",
		Usage:       "[p]rps <rock/paper/scissors>",
		Handler:     f.rps,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "covidtest",
		Description: "Test someone for Covid-19",
		Usage:       "[p]covidtest <user>",
		Handler:     f.covidtest,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hug",
		Description: "Send someone",
		Usage:       "[p]hug",
		Handler:     f.hug,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kiss",
		Description: "Kiss someone",
		Usage:       "[p]kiss",
		Handler:     f.kiss,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "slap",
		Description: "Slap someone",
		Usage:       "[p]slap",
		Handler:     f.slap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cuddle",
		Description: "Cuddle someone",
		Usage:       "[p]cuddle",
		Handler:     f.cuddle,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pat",
		Description: "Pat someone",
		Usage:       "[p]pat",
		Handler:     f.pat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "feed",
		Description: "Feed someone",
		Usage:       "[p]feed",
		Handler:     f.feed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wink",
		Description: "Wink to someone",
		Usage:       "[p]wink",
		Handler:     f.wink,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "chucknorris",
		Description: "Send a random Chuck Norris joke",
		Usage:       "[p]chucknorris",
		Handler:     f.chucknorris,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iss",
		Description: "Send the current ISS coordinates",
		Usage:       "[p]iss",
		Handler:     f.iss,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "apod",
		Description: "Send NASA's Astronomy Picture Of The Day",
		Usage:       "[p]apod",
		Handler:     f.apod,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "httpcat",
		Description: "Send a HTTP Cat image",
		Usage:       "[p]httpcat <statuscode>",
		Handler:     f.httpcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "riddle",
		Description: "Send a random riddle",
		Usage:       "[p]riddle",
		Handler:     f.riddle,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "quote",
		Description: "Send a random quote",
		Usage:       "[p]quote",
		Handler:     f.quote,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "evilinsult",
		Description: "Send a random evil insult",
		Usage:       "[p]evilinsult",
		Handler:     f.evilinsult,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "owoify",
		Description: "Owoify text",
		Usage:       "[p]owoify <text>",
		Handler:     f.owoify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uwuify",
		Description: "Uwuify text",
		Usage:       "[p]uwuify <text>",
		Handler:     f.uwuify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uvuify",
		Description: "Uvuify text",
		Usage:       "[p]uvuify <text>",
		Handler:     f.uvuify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rickrollurl",
		Description: "Send the rick roll url",
		Usage:       "[p]rickrollurl",
		Handler:     f.rickrollurl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rickroll",
		Description: "Send the rick roll lyrics",
		Usage:       "[p]rickroll",
		Handler:     f.rickroll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "impersonate",
		Description: "Impersonate a user using a webhook",
		Usage:       "[p]impersonate <user> <text>",
		Handler:     f.impersonate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "coinflip",
		Description: "Flip a coin",
		Usage:       "[p]coinflip",
		Handler:     f.coinflip,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notfunny",
		Description: "Send a not funny copypasta",
		Usage:       "[p]notfunny",
		Handler:     f.notfunny,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomquestion",
		Description: "Send a random question",
		Usage:       "[p]randomquestion",
		Handler:     f.randomquestion,
	})
}

func (FunCommands) eightball(ctx *dgc.Ctx) {
	question := ctx.Arguments.Raw()
	possibleAnswers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it, yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	answer := utils.RandomChoice(possibleAnswers)
	embed.NewEmbed().
		AddField("Question", question).
		AddField("Answer", answer).
		Send(ctx)
}

func (FunCommands) howgay(ctx *dgc.Ctx) {
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	embed.NewEmbed().
		AddField("User", user).
		AddField("Gay Percentage", utils.RandomPercentage()).
		Send(ctx)
}

func (FunCommands) howsimp(ctx *dgc.Ctx) {
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	embed.NewEmbed().
		AddField("User", user).
		AddField("Simp Percentage", utils.RandomPercentage()).
		Send(ctx)
}

func (FunCommands) penis(ctx *dgc.Ctx) {
	penisLength := utils.RandomNumber(1, 20)
	penis := "8"
	for i := 0; i < penisLength; i++ {
		penis += "="
	}
	penis += "D"
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	embed.NewEmbed().
		AddField("User", user).
		AddField("Penis", penis).
		AddField("Penis Length", fmt.Sprintf("%d cm", penisLength)).
		Send(ctx)
}

func (FunCommands) ego(ctx *dgc.Ctx) {
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	embed.NewEmbed().
		AddField("User", user).
		AddField("Ego", utils.RandomPercentage()).
		Send(ctx)
}

func (FunCommands) hack(ctx *dgc.Ctx) {
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	cc := gofakeit.CreditCard()
	embed.NewEmbed().
		AddField("User", user).
		AddField("Real Name", gofakeit.Name()).
		AddField("Email", gofakeit.Email()).
		AddField("Phone", gofakeit.PhoneFormatted()).
		AddField("Credit Card Number", cc.Number).
		AddField("CCV", cc.Cvv).
		AddField("Address", gofakeit.Address().Address).
		AddField("Job", gofakeit.JobTitle()).
		AddField("Favorite Minecraft Biome", gofakeit.MinecraftBiome()).
		AddField("Idol", gofakeit.CelebrityActor()).
		Send(ctx)
}

func (FunCommands) insult(ctx *dgc.Ctx) {
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	insult, err := http.GetJson(fmt.Sprintf("https://lightning-bot.com/api/v2/selfbot/roast?mention=%s", user))
	if err != nil {
		ctx.RespondText("An error has occurred.")
	}
	embed.NewEmbed().
		AddField("Insult", insult.Get("roast").MustString()).
		Send(ctx)
}

func (FunCommands) dadjoke(ctx *dgc.Ctx) {
	joke, err := http.GetJson("https://icanhazdadjoke.com/", map[string]string{"Accept": "application/json"})
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Dad Joke", joke.Get("joke").MustString()).
		Send(ctx)
}

func (FunCommands) yomomma(ctx *dgc.Ctx) {
	joke, err := http.GetJson("https://api.yomomma.info/")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Yo Momma Joke", joke.Get("joke").MustString()).
		Send(ctx)
}

func (FunCommands) wouldyourather(ctx *dgc.Ctx) {
	html, err := http.Get("https://www.conversationstarters.com/wyrqlist.php")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		SetDescription("Would you rather...").
		AddField("A", utils.GetElementById(html, "#qa")).
		AddField("B", utils.GetElementById(html, "#qb")).
		Send(ctx)
}

func (FunCommands) topic(ctx *dgc.Ctx) {
	html, err := http.Get("https://www.conversationstarters.com/generator.php")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		SetDescription("Would you rather...").
		AddField("Topic", utils.GetElementById(html, "#random")).
		Send(ctx)
}

func (FunCommands) slot(ctx *dgc.Ctx) {
	slots := strings.Split("🍎🍊🍐🍋🍉🍇🍓🍒", "")
	getSlot := func() (string, string) {
		left := utils.RandomChoice(slots)
		mid := utils.RandomChoice(slots)
		right := utils.RandomChoice(slots)
		text := "No matches, you lost!"
		if left == mid && mid == right {
			text = "3 in a row, you win!"
		} else if left == mid || mid == right || left == right {
			text = "2 in a row, you win!"
		}
		return fmt.Sprintf("%s | %s | %s", left, mid, right), text
	}
	getEmbed := func(final bool) string {
		text, winning := getSlot()
		if final {
			text += "\n\n" + winning
		}
		return embed.NewEmbed().
			SetDescription(text).
			ToString()[0]
	}
	message, _ := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, getEmbed(false))
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, getEmbed(false))
	}
	ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, getEmbed(true))
}

func (FunCommands) tts(ctx *dgc.Ctx) {
	ttsFile, err := api.TTS(ctx.Arguments.Raw(), "en")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "tts.mp3", bytes.NewReader(ttsFile))
}

func (FunCommands) advice(ctx *dgc.Ctx) {
	advice, err := http.GetJson("https://api.adviceslip.com/advice")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Advice", advice.GetPath("slip", "advice").MustString()).
		Send(ctx)
}

func (FunCommands) catfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://some-random-api.ml/facts/cat")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Cat Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) dogfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://some-random-api.ml/facts/dog")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Dog Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) pandafact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://some-random-api.ml/facts/panda")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Panda Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) foxfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://some-random-api.ml/facts/fox")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Fox Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) koalafact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://some-random-api.ml/facts/koala")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Koala Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) birdfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/bird")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Bird Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) redpandafact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/redpanda")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Red Panda Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) whalefact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/whale")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Whale Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) kangaroofact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/kangaroo")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Kangaroo Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) bunnyfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/bunny")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Bunny Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) lionfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/lion")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Lion Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) frogfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/frog")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Frog Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) duckfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/duck")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Duck Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) penguinfact(ctx *dgc.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/penguin")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Penguin Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) trivia(ctx *dgc.Ctx) {
	rawJson, err1 := http.Get("https://opentdb.com/api.php?amount=1")
	rawJson = strings.Replace(rawJson, "&quot;", "\"", -1)
	rawJson = strings.Replace(rawJson, "&#039;", "'", -1)
	trivia, err2 := json.SimpleUnmarshal(rawJson)
	if err1 != nil || err2 != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	trivia = trivia.Get("results").GetIndex(0)
	answers := trivia.Get("incorrect_answers").MustStringArray()
	answers = append(answers, trivia.Get("correct_answer").MustString())
	answers = utils.Shuffle(answers)
	embed.NewEmbed().
		AddField("Question", trivia.Get("question").MustString()).
		AddField("Category", trivia.Get("category").MustString()).
		AddField("Difficulty", utils.Capitalize(trivia.Get("difficulty").MustString())).
		AddField("Answers", strings.Join(answers, ", ")).
		AddField("Correct Answer", fmt.Sprintf("||%s||", trivia.Get("correct_answer").MustString())).
		Send(ctx)
}

func (FunCommands) dice(ctx *dgc.Ctx) {
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("You rolled: %d", utils.RandomNumber(1, 6))).
		Send(ctx)
}

func (FunCommands) randomemoji(ctx *dgc.Ctx) {
	if !utils.CheckIfInServer(ctx) {
		return
	}
	emojis, err := ctx.Session.GuildEmojis(ctx.Event.GuildID)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	emoji := utils.RandomChoice(emojis)
	ctx.RespondText(emoji.MessageFormat())
}

func (FunCommands) randomchannel(ctx *dgc.Ctx) {
	if !utils.CheckIfInServer(ctx) {
		return
	}
	channels, err := ctx.Session.GuildChannels(ctx.Event.GuildID)
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	for index, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			// Remove the channel from the list if it is not a text channel
			channels = append(channels[index:], channels[index+1:]...)
		}
	}
	channel := utils.RandomChoice(channels)
	ctx.RespondText(channel.Mention())
}

func (FunCommands) choose(ctx *dgc.Ctx) {
	choice := utils.RandomChoice(strings.Split(ctx.Arguments.Raw(), " "))
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("I choose: %s", choice)).
		Send(ctx)
}

func (FunCommands) poll(ctx *dgc.Ctx) {
	message := embed.NewEmbed().
		SetDescription("Answer to this poll by reacting with :thumbsup: or :thumbsdown:").
		AddField("Question", ctx.Arguments.Raw()).
		Send(ctx)
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, message.ID, "👍")
	ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, message.ID, "👎")
}

func (FunCommands) numberfact(ctx *dgc.Ctx) {
	number := utils.Arguments["<number>"]
	fact, err := http.Get(fmt.Sprintf("http://numbersapi.com/%s", number))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Number", number).
		AddField("Fact", fact).
		Send(ctx)
}

func (FunCommands) joke(ctx *dgc.Ctx) {
	joke, err := http.GetJson("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed := embed.NewEmbed()
	if joke.Get("type").MustString() == "twopart" {
		embed.AddField("Setup", joke.Get("setup").MustString()).
			AddField("Delivery", joke.Get("delivery").MustString())
	} else {
		embed.AddField("Joke", joke.Get("joke").MustString())
	}
	embed.AddField("Category", joke.Get("category").MustString()).Send(ctx)
}

func (FunCommands) fml(ctx *dgc.Ctx) {
	html, err := http.Get("https://www.fmylife.com/random")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("FML", doc.Find("a.block text-blue-500 my-4").Text()).
		Send(ctx)
}

func (FunCommands) tronald(ctx *dgc.Ctx) {
	tronald, err := http.GetJson("https://api.tronalddump.io/random/quote", map[string]string{"Accept": "application/hal+json"})
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Tronald Dump Quote", tronald.Get("value").MustString()).
		Send(ctx)
}

func (FunCommands) kanyewest(ctx *dgc.Ctx) {
	kanyewest, err := http.GetJson("https://api.kanye.rest")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Kanye West Quote", kanyewest.Get("quote").MustString()).
		Send(ctx)
}

func (FunCommands) rps(ctx *dgc.Ctx) {
	choice := utils.Arguments["<rock/paper/scissors>"]
	if choice != "rock" && choice != "paper" && choice != "scissors" {
		ctx.RespondText("You can only choose rock, paper or scissors.")
		return
	}
	embed.NewEmbed().
		AddField("Your choice", choice).
		AddField("Bot's choice", utils.RandomChoice([]string{"rock", "paper", "scissors"})).
		Send(ctx)
}

func (FunCommands) covidtest(ctx *dgc.Ctx) {
	user, _ := utils.TryGetUsernameFromMention(ctx.Session, ctx.Arguments.Get(0))
	embed.NewEmbed().
		AddField("User", user).
		AddField("Result", utils.RandomChoice([]string{"Positive", "Negative"})).
		Send(ctx)
}

func (FunCommands) hug(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://nekos.life/api/v2/img/hug", "url")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) kiss(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://nekos.life/api/v2/img/kiss", "url")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) slap(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://nekos.life/api/v2/img/slap", "url")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) cuddle(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://nekos.life/api/v2/img/cuddle", "url")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) pat(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://nekos.life/api/v2/img/pat", "url")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) feed(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://nekos.life/api/v2/img/feed", "url")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) wink(ctx *dgc.Ctx) {
	url, err := utils.GetImageUrlFromJson("https://some-random-api.ml/animu/wink", "link")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	ctx.RespondText(url)
}

func (FunCommands) chucknorris(ctx *dgc.Ctx) {
	joke, err := http.GetJson("https://api.chucknorris.io/jokes/random")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Chuck Norris Joke", joke.Get("value").MustString()).
		Send(ctx)
}

func (FunCommands) iss(ctx *dgc.Ctx) {
	iss, err := http.GetJson("http://api.open-notify.org/iss-now.json")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	latitude := iss.GetPath("iss_position", "latitude").MustString()
	longitude := iss.GetPath("iss_position", "longitude").MustString()
	embed.NewEmbed().
		SetAuthor("International Space Station Info").
		AddField("Latitude", latitude).
		AddField("Longitude", longitude).
		AddField("Google Maps URL", fmt.Sprintf("https://maps.google.com/?q=%s,%s", latitude, longitude)).
		Send(ctx)
}

func (FunCommands) apod(ctx *dgc.Ctx) {
	apod, err := http.GetJson("https://api.nasa.gov/planetary/apod?api_key=rDenYNQFftu3YADWdUmF9s8iw0cyZckV7vBcibGe")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	copyright := apod.Get("copyright").MustString()
	if copyright == "" {
		copyright = "NASA"
	}
	explanation, err := api.CreatePaste(apod.Get("explanation").MustString())
	if err != nil {
		explanation = "An error has occurred."
	}
	embed.NewEmbed().
		SetImage(apod.Get("url").MustString()).
		SetAuthor(fmt.Sprintf("Astronomy Picture of the Day - ©️%s", copyright)).
		AddField("Title", apod.Get("title").MustString()).
		AddField("Explanation", explanation).
		AddField("Date", apod.Get("date").MustString()).
		Send(ctx)
}

func (FunCommands) httpcat(ctx *dgc.Ctx) {
	statuscode := utils.Arguments["<statuscode>"]
	embed.NewEmbed().
		SetImage(fmt.Sprintf("https://http.cat/%s.jpg", statuscode)).
		Send(ctx)
}

func (FunCommands) riddle(ctx *dgc.Ctx) {
	html, err := http.Get("https://generatorfun.com/code/model/generatorcontent.php?recordtable=generator&recordkey=69&gen=Y&itemnumber=10&randomoption=undefined&genimage=Yes&geneditor=No&nsfw=undefined&keyword=undefined&searchfilter=&searchfilterexclude=&tone=Normal&prefix=None&randomai=No")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	riddles := doc.Find("#gencont").Find("p").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	riddle := strings.Replace(utils.RandomChoice(riddles), "SHOW", "", -1)
	embed.NewEmbed().
		AddField("Riddle", riddle).
		Send(ctx)
}

func (FunCommands) quote(ctx *dgc.Ctx) {
	quote, err := http.GetJson("https://api.quotable.io/random")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Quote", quote.Get("content").MustString()).
		AddField("Author", quote.Get("author").MustString()).
		Send(ctx)
}

func (FunCommands) evilinsult(ctx *dgc.Ctx) {
	insult, err := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Evil Insult", insult).
		Send(ctx)
}

func (FunCommands) owoify(ctx *dgc.Ctx) {
	ctx.RespondText(owoify_go.Owoify(ctx.Arguments.Raw(), "owo"))
}

func (FunCommands) uwuify(ctx *dgc.Ctx) {
	ctx.RespondText(owoify_go.Owoify(ctx.Arguments.Raw(), "uwu"))
}

func (FunCommands) uvuify(ctx *dgc.Ctx) {
	ctx.RespondText(owoify_go.Owoify(ctx.Arguments.Raw(), "uvu"))
}

func (FunCommands) rickrollurl(ctx *dgc.Ctx) {
	ctx.RespondText("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
}

func (FunCommands) rickroll(ctx *dgc.Ctx) {
	ctx.RespondText("**We're no strangers to love\nYou know the rules and so do I\nA full commitment's what I'm thinking of\nYou wouldn't get this from any other guy\nI just wanna tell you how I'm feeling\nGotta make you understand\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nWe've known each other for so long\nYour heart's been aching but you're too shy to say it\nInside we both know what's been going on\nWe know the game and we're gonna play it\nAnd if you ask me how I'm feeling\nDon't tell me you're too blind to see\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give, never gonna give\n(Give you up)\n(Ooh) Never gonna give, never gonna give\n(Give you up)\nWe've known each other for so long\nYour heart's been aching but you're too shy to say it\nInside we both know what's been going on\nWe know the game and we're gonna play it\nI just wanna tell you how I'm feeling\nGotta make you understand\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry**")
}

func (FunCommands) impersonate(ctx *dgc.Ctx) {
	user, success := utils.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user.")
		return
	}
	text := utils.RemoveWords(ctx.Arguments.Raw(), 1)
	webhooks, err := ctx.Session.ChannelWebhooks(ctx.Event.ChannelID)
	if err != nil {
		ctx.RespondText("You need the manage_webhooks permission to use this command.")
		return
	}
	var webhook *discordgo.Webhook
	if len(webhooks) == 0 {
		webhook, err = ctx.Session.WebhookCreate(ctx.Event.ChannelID, "Impersonate", "")
		if err != nil {
			ctx.RespondText("An error has occurred creating a webhook.")
			return
		}
	} else {
		webhook = webhooks[0]
	}
	ctx.Session.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Username:  user.Username,
		AvatarURL: user.AvatarURL(""),
		Content:   text,
	})
}

func (FunCommands) coinflip(ctx *dgc.Ctx) {
	choice := utils.RandomChoice([]string{"Heads!", "Tails!"})
	images := map[string]string{
		"Heads!": "https://www.ssaurel.com/blog/wp-content/uploads/2017/01/heads.png",
		"Tails!": "https://www.ssaurel.com/blog/wp-content/uploads/2017/01/tails.png",
	}
	embed.NewEmbed().
		SetDescription(choice).
		SetImage(images[choice]).
		Send(ctx)
}

func (FunCommands) notfunny(ctx *dgc.Ctx) {
	ctx.RespondText("**Not funny I didn't laugh. Your joke is so bad I would have preferred the joke went over my head and you gave up re-telling me the joke. To be honest this is a horrid attempt at trying to get a laugh out of me. Not a chuckle, not a hehe, not even a subtle burst of air out of my esophagus. Science says before you laugh your brain preps your face muscles but I didn't even feel the slightest twitch. 0/10 this joke is so bad I cannot believe anyone legally allowed you to be creative at all. The amount of brain power you must have put into that joke has the potential to power every house on Earth. Get a personality and learn how to make jokes, read a book. I'm not saying this to be funny I genuinely mean it on how this is just bottom barrel embarrassment at comedy. You've single handedly killed humor and every comedic act on the planet. I'm so disappointed that society has failed as a whole in being able to teach you how to be funny. Honestly if I put in all my power and time to try and make your joke funny it would require Einstein himself to build a device to strap me into so I can be connected to the energy of a billion stars to do it, and even then all that joke would get from people is a subtle scuff. You're lucky I still have the slightest of empathy for you after telling that joke otherwise I would have committed every war crime in the book just to prevent you from attempting any humor ever again. We should put that joke in text books so future generations can be wary of becoming such an absolute comedic failure. Im disappointed, hurt, and outright offended that my precious time has been wasted in my brain understanding that joke. In the time that took I was planning on helping kids who have been orphaned, but because of that you've waisted my time explaining the obscene integrity of your terrible attempt at comedy. Now those kids are suffering without meals and there's nobody to blame but you.**")
	ctx.RespondText("**I hope you're happy with what you have done and I truly hope you can move on and learn from this piss poor attempt**")
}

func (FunCommands) randomquestion(ctx *dgc.Ctx) {
	question, err := http.GetJson("https://nekos.life/api/v2/why")
	if err != nil {
		ctx.RespondText("An error has occurred.")
		return
	}
	embed.NewEmbed().
		AddField("Question", question.Get("why").MustString()).
		Send(ctx)
}
