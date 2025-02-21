package fun

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"github.com/brianvoe/gofakeit/v6"
	owoifygo "github.com/deadshot465/owoify-go"
)

func (f FunCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Fun", "Fun commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "eightball",
		Description: "Asks the 8-Ball a question",
		Arguments:   "<question>",
		Handler:     f.eightball,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "howgay",
		Description: "See how gay a person is",
		Arguments:   "<user>",
		Handler:     f.howgay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "howsimp",
		Description: "See how simp a person is",
		Arguments:   "<user>",
		Handler:     f.howsimp,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "penis",
		Description: "See how big a person's penis is",
		Arguments:   "<user>",
		Handler:     f.penis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ego",
		Description: "See how big a person's ego is",
		Arguments:   "<user>",
		Handler:     f.ego,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iqtest",
		Description: "Test a person's IQ",
		Arguments:   "<user>",
		Handler:     f.iqtest,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hack",
		Description: "Fake hack a person",
		Arguments:   "<user>",
		Handler:     f.hack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "insult",
		Description: "Insult a person",
		Usage:       "[p]insult",
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
		Name:        "dolphinfact",
		Description: "Send a random dolphin fact",
		Usage:       "[p]dolphinfact",
		Handler:     f.dolphinfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bearfact",
		Description: "Send a random bear fact",
		Usage:       "[p]bearfact",
		Handler:     f.bearfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "axolotlfact",
		Description: "Send a random axolotl fact",
		Usage:       "[p]axolotlfact",
		Handler:     f.axolotlfact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "capybarafact",
		Description: "Send a random capybara fact",
		Usage:       "[p]capybarafact",
		Handler:     f.capybarafact,
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
		Arguments:   "<option1> <option2> ...",
		Handler:     f.choose,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "poll",
		Description: "Create a poll",
		Arguments:   "<question>",
		Handler:     f.poll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "numberfact",
		Description: "Send a random fact about a number",
		Arguments:   "<number>",
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
		Arguments:   "<rock/paper/scissors>",
		Handler:     f.rps,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "covidtest",
		Description: "Test someone for Covid-19",
		Arguments:   "<user>",
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
		Arguments:   "<statuscode>",
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
		Arguments:   "<text>",
		Handler:     f.owoify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uwuify",
		Description: "Uwuify text",
		Arguments:   "<text>",
		Handler:     f.uwuify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uvuify",
		Description: "Uvuify text",
		Arguments:   "<text>",
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
	router.RegisterCmd(&dgc.Command{
		Name:        "lyrics",
		Description: "Get the lyrics of a song",
		Arguments:   "<song name>",
		Handler:     f.lyrics,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "breakingbadquote",
		Description: "Send a random Breaking Badquote",
		Usage:       "[p]breakingbadquote",
		Handler:     f.breakingbadquote,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gameofthronesquote",
		Description: "Send a random Game of Thrones quote",
		Usage:       "[p]gameofthronesquote",
		Handler:     f.gameofthronesquote,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "luciferquote",
		Description: "Send a random Lucifer quote",
		Usage:       "[p]luciferquote",
		Handler:     f.luciferquote,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "strangerthingsquote",
		Description: "Send a random Stranger Things quote",
		Usage:       "[p]strangerthingsquote",
		Handler:     f.strangerthingsquote,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "xkcd",
		Description: "Send the current XKCD comic",
		Usage:       "[p]xkcd",
		Handler:     f.xkcd,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "boredactivity",
		Description: "Send a random activity to do when bored",
		Usage:       "[p]boredactivity",
		Handler:     f.boredactivity,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ermahgerd",
		Description: "Translate text to ERMAHGERD",
		Arguments:   "<text>",
		Handler:     f.ermahgerd,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iseven",
		Description: "Check if a number is even",
		Arguments:   "<number>",
		Handler:     f.iseven,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "agify",
		Description: "Determine the age of a first name",
		Arguments:   "<name>",
		Handler:     f.agify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomword",
		Description: "Get a random word",
		Usage:       "[p]randomword",
		Handler:     f.randomword,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomweirdword",
		Description: "Get a random weird word",
		Usage:       "[p]randomweirdword",
		Handler:     f.randomweirdword,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "startupidea",
		Description: "Get a random startup idea",
		Usage:       "[p]startupidea",
		Handler:     f.startupidea,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakeperson",
		Description: "Generate an image of a fake person",
		Usage:       "[p]fakeperson",
		Handler:     f.fakeperson,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "commitmessage",
		Description: "Generate a random commit message",
		Usage:       "[p]commitmessage",
		Handler:     f.commitmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "showerthought",
		Description: "Send a random shower thought",
		Usage:       "[p]showerthought",
		Handler:     f.showerthought,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "truth",
		Description: "Play Truth or Dare",
		Usage:       "[p]truth",
		Handler:     f.truth,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dare",
		Description: "Play Truth or Dare",
		Usage:       "[p]dare",
		Handler:     f.dare,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "neverhaveiever",
		Description: "Play Never have I ever...",
		Usage:       "[p]neverhaveiever",
		Handler:     f.neverhaveiever,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "paranoia",
		Description: "Send a random paranoia question",
		Usage:       "[p]paranoia",
		Handler:     f.paranoia,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomcolor",
		Description: "Generate a random color",
		Usage:       "[p]randomcolor",
		Handler:     f.randomcolor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pickupline",
		Description: "Generate a pickup line",
		Usage:       "[p]pickupline",
		Handler:     f.pickupline,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "elementinfo",
		Description: "Get info about a periodic table element",
		Arguments:   "<element>",
		Handler:     f.elementinfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomelement",
		Description: "Get info about a random periodic table element",
		Usage:       "[p]randomelement",
		Handler:     f.randomelement,
	})
}

func (FunCommands) eightball(ctx *router.Ctx) {
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
	answer := random.Choice(possibleAnswers)
	c.services.Embed.NewEmbed().
		AddField("Question", question).
		AddField("Answer", answer).
		Send(ctx)
}

func (FunCommands) howgay(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
		AddField("Gay Percentage", random.Percentage()).
		Send(ctx)
}

func (FunCommands) howsimp(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
		AddField("Simp Percentage", random.Percentage()).
		Send(ctx)
}

func (FunCommands) penis(ctx *router.Ctx) {
	penisLength := random.Int(1, 20)
	penis := "8"
	for i := 0; i < penisLength; i++ {
		penis += "="
	}
	penis += "Đ"
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
		AddField("Penis", penis).
		AddField("Penis Length", fmt.Sprintf("%d cm", penisLength)).
		Send(ctx)
}

func (FunCommands) ego(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
		AddField("Ego", random.Percentage()).
		Send(ctx)
}

func (FunCommands) iqtest(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
		AddField("IQ", fmt.Sprint(random.Int(1, 200))).
		Send(ctx)
}

func (FunCommands) hack(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	cc := gofakeit.CreditCard()
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
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

func (FunCommands) insult(ctx *router.Ctx) {
	httpResponse := http.Get("https://insult.mattbas.org/api/insult")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Insult", string(httpResponse.Body)).
		Send(ctx)
}

func (FunCommands) dadjoke(ctx *router.Ctx) {
	joke, err := http.GetJson("https://icanhazdadjoke.com/", map[string]string{"Accept": "application/json"})
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Dad Joke", joke.Get("joke").MustString()).
		Send(ctx)
}

func (FunCommands) yomomma(ctx *router.Ctx) {
	joke, err := http.GetJson("https://api.yomomma.info/")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Yo Momma Joke", joke.Get("joke").MustString()).
		Send(ctx)
}

func (FunCommands) wouldyourather(ctx *router.Ctx) {
	wyrJson, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/wyr")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		SetDescription("Would you rather...").
		AddField("A", wyrJson.GetPath("question", "a").MustString()).
		AddField("B", wyrJson.GetPath("question", "b").MustString()).
		Send(ctx)
}

func (FunCommands) topic(ctx *router.Ctx) {
	topicJson, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/topic")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Topic", topicJson.Get("topic").MustString()).
		Send(ctx)
}

func (FunCommands) slot(ctx *router.Ctx) {
	slots := strings.Split("🍎🍊🍐🍋🍉🍇🍓🍒", "")
	getSlot := func() (string, string) {
		left := random.Choice(slots)
		mid := random.Choice(slots)
		right := random.Choice(slots)
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
		return c.services.Embed.NewEmbed().
			SetDescription(text).
			ToString("load")[0]
	}

	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, getEmbed(false))
	if err != nil {
		return
	}

	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)

		_, err = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, getEmbed(false))
		if err != nil {
			return
		}
	}

	_, _ = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, getEmbed(true))
}

func (FunCommands) advice(ctx *router.Ctx) {
	advice, err := http.GetJson("https://api.adviceslip.com/advice")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Advice", advice.GetPath("slip", "advice").MustString()).
		Send(ctx)
}

func (FunCommands) catfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://some-random-api.com/facts/cat")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Cat Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) dogfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://some-random-api.com/facts/dog")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Dog Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) pandafact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://some-random-api.com/facts/panda")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Panda Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) foxfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://some-random-api.com/facts/fox")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Fox Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) koalafact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://some-random-api.com/facts/koala")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Koala Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) birdfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/bird")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Bird Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) redpandafact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/redpanda")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Red Panda Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) whalefact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/whale")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Whale Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) kangaroofact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/kangaroo")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Kangaroo Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) bunnyfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/bunny")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Bunny Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) lionfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/lion")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Lion Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) frogfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/frog")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Frog Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) duckfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/duck")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Duck Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) penguinfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/penguin")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Penguin Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) dolphinfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/dolphin")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Dolphin Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) bearfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/bear")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Bear Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) axolotlfact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/axolotl")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Axolotl Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) capybarafact(ctx *router.Ctx) {
	fact, err := http.GetJson("https://api.animality.xyz/fact/capybara")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Capybara Fact", fact.Get("fact").MustString()).
		Send(ctx)
}

func (FunCommands) trivia(ctx *router.Ctx) {
	httpResponse := http.Get("https://opentdb.com/api.php?amount=1")
	if httpResponse.StatusCode != 200 {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	data := httpResponse.Json()

	if data.Get("response_code").MustInt64() != 0 {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	data = data.Get("results").GetIndex(0)

	answers := data.Get("incorrect_answers").MustStringArray()
	answers = append(answers, data.Get("correct_answer").MustString())
	answers = random.Shuffle(answers)

	c.services.Embed.NewEmbed().
		AddField("Question", html.UnescapeString(data.Get("question").MustString())).
		AddField("Category", data.Get("category").MustString()).
		AddField("Difficulty", utils.String.Capitalize(data.Get("difficulty").MustString())).
		AddField("Answers", html.UnescapeString(strings.Join(answers, ", "))).
		AddField("Correct Answer", fmt.Sprintf("||%s||", html.UnescapeString(data.Get("correct_answer").MustString()))).
		Send(ctx)
}

func (FunCommands) dice(ctx *router.Ctx) {
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("You rolled: %d", random.Int(1, 6))).
		Send(ctx)
}

func (FunCommands) randomemoji(ctx *router.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	emojis, err := ctx.Session.GuildEmojis(ctx.Message.GuildID)
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	emoji := random.Choice(emojis)
	_ = ctx.RespondText(emoji.MessageFormat())
}

func (FunCommands) randomchannel(ctx *router.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	channels, err := ctx.Session.GuildChannels(ctx.Message.GuildID)
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	for index, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			// Remove the channel from the list if it is not a text channel
			channels = append(channels[index:], channels[index+1:]...)
		}
	}
	channel := random.Choice(channels)
	_ = ctx.RespondText(channel.Mention())
}

func (FunCommands) choose(ctx *router.Ctx) {
	choice := random.Choice(strings.Split(ctx.Arguments.Raw(), " "))
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("I choose: %s", choice)).
		Send(ctx)
}

func (FunCommands) poll(ctx *router.Ctx) {
	messages := c.services.Embed.NewEmbed().
		SetDescription("Answer to this poll by reacting with :thumbsup: or :thumbsdown:").
		AddField("Question", ctx.Arguments.Raw()).
		Send(ctx)

	if len(messages) > 0 {
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, messages[0].ID, "👍")
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, messages[0].ID, "👎")
	}
}

func (FunCommands) numberfact(ctx *router.Ctx) {
	number := utils.Arguments["<number>"]
	httpResponse := http.Get(fmt.Sprintf("http://numbersapi.com/%s", number))
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Number", number).
		AddField("Fact", string(httpResponse.Body)).
		Send(ctx)
}

func (FunCommands) joke(ctx *router.Ctx) {
	joke, err := http.GetJson("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	emb := c.services.Embed.NewEmbed()
	if joke.Get("type").MustString() == "twopart" {
		emb.AddField("Setup", joke.Get("setup").MustString()).
			AddField("Delivery", joke.Get("delivery").MustString())
	} else {
		emb.AddField("Joke", joke.Get("joke").MustString())
	}
	emb.AddField("Category", joke.Get("category").MustString()).Send(ctx)
}

func (FunCommands) fml(ctx *router.Ctx) {
	fmlJson, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/fml")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("FML", fmlJson.Get("fml").MustString()).
		Send(ctx)
}

func (FunCommands) tronald(ctx *router.Ctx) {
	tronald, err := http.GetJson("https://api.tronalddump.io/random/quote", map[string]string{"Accept": "application/hal+json"})
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Tronald Dump Quote", tronald.Get("value").MustString()).
		Send(ctx)
}

func (FunCommands) kanyewest(ctx *router.Ctx) {
	kanyewest, err := http.GetJson("https://api.kanye.rest")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Kanye West Quote", kanyewest.Get("quote").MustString()).
		Send(ctx)
}

func (FunCommands) rps(ctx *router.Ctx) {
	choice := utils.Arguments["<rock/paper/scissors>"]
	if choice != "rock" && choice != "paper" && choice != "scissors" {
		_ = ctx.RespondText("You can only choose rock, paper or scissors.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Your choice", choice).
		AddField("Bot's choice", random.Choice([]string{"rock", "paper", "scissors"})).
		Send(ctx)
}

func (FunCommands) covidtest(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Couldn't find user!")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("User", user.String()).
		AddField("Result", random.Choice([]string{"Positive", "Negative"})).
		Send(ctx)
}

func (FunCommands) hug(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/hug", "url")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) kiss(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/kiss", "url")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) slap(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/slap", "url")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) cuddle(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/cuddle", "url")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) pat(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/pat", "url")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) feed(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/feed", "url")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) wink(ctx *router.Ctx) {
	url, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/animu/wink", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	_ = ctx.RespondText(url)
}

func (FunCommands) chucknorris(ctx *router.Ctx) {
	joke, err := http.GetJson("https://api.chucknorris.io/jokes/random")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Chuck Norris Joke", joke.Get("value").MustString()).
		Send(ctx)
}

func (FunCommands) iss(ctx *router.Ctx) {
	iss, err := http.GetJson("http://api.open-notify.org/iss-now.json")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	latitude := iss.GetPath("iss_position", "latitude").MustString()
	longitude := iss.GetPath("iss_position", "longitude").MustString()
	c.services.Embed.NewEmbed().
		SetAuthor("International Space Station Info").
		AddField("Latitude", latitude).
		AddField("Longitude", longitude).
		AddField("Google Maps URL", fmt.Sprintf("https://maps.google.com/?q=%s,%s", latitude, longitude)).
		Send(ctx)
}

func (FunCommands) apod(ctx *router.Ctx) {
	apod, err := http.GetJson("https://api.nasa.gov/planetary/apod?api_key=rDenYNQFftu3YADWdUmF9s8iw0cyZckV7vBcibGe")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
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
	c.services.Embed.NewEmbed().
		SetImage(apod.Get("url").MustString()).
		SetAuthor(fmt.Sprintf("Astronomy Picture of the Day - ©️%s", copyright)).
		AddField("Title", apod.Get("title").MustString()).
		AddField("Explanation", explanation).
		AddField("Date", apod.Get("date").MustString()).
		Send(ctx)
}

func (FunCommands) httpcat(ctx *router.Ctx) {
	statuscode := utils.Arguments["<statuscode>"]
	c.services.Embed.NewEmbed().
		SetImage(fmt.Sprintf("https://http.cat/%s.jpg", statuscode)).
		Send(ctx)
}

func (FunCommands) riddle(ctx *router.Ctx) {
	riddleJson, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/riddle")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Riddle", riddleJson.Get("question").MustString()).
		AddField("Answer", "||"+riddleJson.Get("answer").MustString()+"||").
		Send(ctx)
}

func (FunCommands) quote(ctx *router.Ctx) {
	quote, err := http.GetJson("https://api.quotable.io/random")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Quote", quote.Get("content").MustString()).
		AddField("Author", quote.Get("author").MustString()).
		Send(ctx)
}

func (FunCommands) evilinsult(ctx *router.Ctx) {
	httpResponse := http.Get("https://evilinsult.com/generate_insult.php?lang=en")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Evil Insult", string(httpResponse.Body)).
		Send(ctx)
}

func (FunCommands) owoify(ctx *router.Ctx) {
	_ = ctx.RespondText(owoifygo.Owoify(ctx.Arguments.Raw(), "owo"))
}

func (FunCommands) uwuify(ctx *router.Ctx) {
	_ = ctx.RespondText(owoifygo.Owoify(ctx.Arguments.Raw(), "uwu"))
}

func (FunCommands) uvuify(ctx *router.Ctx) {
	_ = ctx.RespondText(owoifygo.Owoify(ctx.Arguments.Raw(), "uvu"))
}

func (FunCommands) rickrollurl(ctx *router.Ctx) {
	_ = ctx.RespondText("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
}

func (FunCommands) rickroll(ctx *router.Ctx) {
	_ = ctx.RespondText("**We're no strangers to love\nYou know the rules and so do I\nA full commitment's what I'm thinking of\nYou wouldn't get this from any other guy\nI just wanna tell you how I'm feeling\nGotta make you understand\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nWe've known each other for so long\nYour heart's been aching but you're too shy to say it\nInside we both know what's been going on\nWe know the game and we're gonna play it\nAnd if you ask me how I'm feeling\nDon't tell me you're too blind to see\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give, never gonna give\n(Give you up)\n(Ooh) Never gonna give, never gonna give\n(Give you up)\nWe've known each other for so long\nYour heart's been aching but you're too shy to say it\nInside we both know what's been going on\nWe know the game and we're gonna play it\nI just wanna tell you how I'm feeling\nGotta make you understand\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry**")
}

func (FunCommands) coinflip(ctx *router.Ctx) {
	choice := random.Choice([]string{"Heads!", "Tails!"})
	images := map[string]string{
		"Heads!": "https://www.ssaurel.com/blog/wp-content/uploads/2017/01/heads.png",
		"Tails!": "https://www.ssaurel.com/blog/wp-content/uploads/2017/01/tails.png",
	}
	c.services.Embed.NewEmbed().
		SetDescription(choice).
		SetImage(images[choice]).
		Send(ctx)
}

func (FunCommands) notfunny(ctx *router.Ctx) {
	_ = ctx.RespondText("**Not funny I didn't laugh. Your joke is so bad I would have preferred the joke went over my head and you gave up re-telling me the joke. To be honest this is a horrid attempt at trying to get a laugh out of me. Not a chuckle, not a hehe, not even a subtle burst of air out of my esophagus. Science says before you laugh your brain preps your face muscles but I didn't even feel the slightest twitch. 0/10 this joke is so bad I cannot believe anyone legally allowed you to be creative at all. The amount of brain power you must have put into that joke has the potential to power every house on Earth. Get a personality and learn how to make jokes, read a book. I'm not saying this to be funny I genuinely mean it on how this is just bottom barrel embarrassment at comedy. You've single-handedly killed humor and every comedic act on the planet. I'm so disappointed that society has failed as a whole in being able to teach you how to be funny. Honestly if I put in all my power and time to try and make your joke funny it would require Einstein himself to build a device to strap me into so I can be connected to the energy of a billion stars to do it, and even then all that joke would get from people is a subtle scuff. You're lucky I still have the slightest of empathy for you after telling that joke otherwise I would have committed every war crime in the book just to prevent you from attempting any humor ever again. We should put that joke in text books so future generations can be wary of becoming such an absolute comedic failure. Im disappointed, hurt, and outright offended that my precious time has been wasted in my brain understanding that joke. In the time that took I was planning on helping kids who have been orphaned, but because of that you've waisted my time explaining the obscene integrity of your terrible attempt at comedy. Now those kids are suffering without meals and there's nobody to blame but you.**")
	_ = ctx.RespondText("**I hope you're happy with what you have done and I truly hope you can move on and learn from this piss poor attempt**")
}

func (FunCommands) randomquestion(ctx *router.Ctx) {
	question, err := http.GetJson("https://nekos.life/api/v2/why")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Question", question.Get("why").MustString()).
		Send(ctx)
}

func (FunCommands) lyrics(ctx *router.Ctx) {
	lyrics, err := http.GetJson("https://some-random-api.com/lyrics?title=" + url.QueryEscape(ctx.Arguments.Raw()))
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	if lyrics.Get("error").MustString() != "" {
		_ = ctx.RespondText("Could not find lyrics.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Artist", lyrics.Get("author").MustString()).
		AddField("Title", lyrics.Get("title").MustString()).
		AddField("Lyrics", lyrics.Get("lyrics").MustString()).
		Send(ctx)
}

func (FunCommands) breakingbadquote(ctx *router.Ctx) {
	quote, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/breaking-bad-quote")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Quote", quote.Get("quote").MustString()).
		AddField("Author", quote.Get("author").MustString()).
		Send(ctx)
}

func (FunCommands) gameofthronesquote(ctx *router.Ctx) {
	quote, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/game-of-thrones-quote")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Quote", quote.Get("quote").MustString()).
		AddField("Author", quote.Get("author").MustString()).
		Send(ctx)
}

func (FunCommands) strangerthingsquote(ctx *router.Ctx) {
	quote, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/stranger-things-quote")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Quote", quote.Get("quote").MustString()).
		AddField("Author", quote.Get("author").MustString()).
		Send(ctx)
}

func (FunCommands) luciferquote(ctx *router.Ctx) {
	quote, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/lucifer-quote")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	c.services.Embed.NewEmbed().
		AddField("Quote", quote.Get("quote").MustString()).
		AddField("Author", quote.Get("author").MustString()).
		Send(ctx)
}

func (FunCommands) xkcd(ctx *router.Ctx) {
	json, err := http.GetJson("https://xkcd.com/info.0.json")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	date := json.Get("day").MustString() + "/" + json.Get("month").MustString() + "/" + json.Get("year").MustString()

	c.services.Embed.NewEmbed().
		AddField("Title", json.Get("title").MustString()).
		AddField("Date", date).
		SetImage(json.Get("img").MustString()).
		Send(ctx)
}

func (FunCommands) boredactivity(ctx *router.Ctx) {
	json, err := http.GetJson("https://www.boredapi.com/api/activity")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	activityType := utils.String.Capitalize(json.Get("type").MustString())

	c.services.Embed.NewEmbed().
		AddField("Activity", json.Get("activity").MustString()).
		AddField("Type", activityType).
		SetImage(json.Get("img").MustString()).
		Send(ctx)
}

func (FunCommands) ermahgerd(ctx *router.Ctx) {
	ctx.RespondText(utils.String.Ermahgerd(ctx.Arguments.Raw()))
}

func (FunCommands) iseven(ctx *router.Ctx) {
	number, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		_ = ctx.RespondText("Argument <number> must be a number.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Number", ctx.Arguments.Get(0).Raw()).
		AddField("Is Even", utils.Misc.FormatBoolYesNo(number%2 == 0)).
		Send(ctx)
}

func (FunCommands) agify(ctx *router.Ctx) {
	json, err := http.GetJson("https://api.agify.io?name=" + url.QueryEscape(ctx.Arguments.Raw()))
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Name", json.Get("name").MustString()).
		AddField("Age", fmt.Sprint(json.Get("age").MustInt())).
		AddField("Count", fmt.Sprint(json.Get("count").MustInt())).
		Send(ctx)
}

func (FunCommands) randomword(ctx *router.Ctx) {
	json, err := http.GetJson("https://randomwordgenerator.com/json/words_ws.json")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	words := json.Get("data")
	word := words.GetIndex(random.Int(0, len(words.MustArray()))).Get("word")

	c.services.Embed.NewEmbed().
		AddField("Random Word", word.Get("value").MustString()).
		Send(ctx)
}

func (FunCommands) randomweirdword(ctx *router.Ctx) {
	json, err := http.GetJson("https://randomwordgenerator.com/json/unusual-words.json")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	words := json.Get("data")
	word := words.GetIndex(random.Int(0, len(words.MustArray()))).Get("word")

	c.services.Embed.NewEmbed().
		AddField("Random Weird Word", word.Get("value").MustString()).
		AddField("Definition", word.Get("definition").MustString()).
		Send(ctx)
}

func (FunCommands) startupidea(ctx *router.Ctx) {
	httpResponse := http.Get("https://itsthisforthat.com/api.php?text")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Startup Idea", string(httpResponse.Body)).
		Send(ctx)
}

func (FunCommands) fakeperson(ctx *router.Ctx) {
	c.services.Embed.NewEmbed().
		SetDescription("Fake Person").
		SetImage("https://thispersondoesnotexist.com?" + fmt.Sprint(random.Int(1000, 9999))).
		Send(ctx)
}

func (FunCommands) commitmessage(ctx *router.Ctx) {
	httpResponse := http.Get("https://whatthecommit.com/index.txt")
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Commit Message", string(httpResponse.Body)).
		Send(ctx)
}

func (FunCommands) showerthought(ctx *router.Ctx) {
	data, err := http.GetJson("https://showerthoughts.nickf.me/api/thoughts/random?limit=1")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Shower Thought", data.Get("title").MustString()).
		Send(ctx)
}

func (FunCommands) truth(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.truthordarebot.xyz/v1/truth")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Truth Question", data.Get("question").MustString()).
		Send(ctx)
}

func (FunCommands) dare(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.truthordarebot.xyz/v1/dare")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Dare Question", data.Get("question").MustString()).
		Send(ctx)
}

func (FunCommands) neverhaveiever(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.truthordarebot.xyz/v1/nhie")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Never have I ever...", data.Get("question").MustString()).
		Send(ctx)
}

func (FunCommands) paranoia(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.truthordarebot.xyz/v1/paranoia")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Paranoia Question", data.Get("question").MustString()).
		Send(ctx)
}

func (FunCommands) randomcolor(ctx *router.Ctx) {
	r, g, b := random.Int(0, 255), random.Int(0, 255), random.Int(0, 255)
	hex := fmt.Sprintf("%02x%02x%02x", r, g, b)
	name := "Unknown"

	data, err := http.GetJson("https://www.thecolorapi.com/id?hex=" + hex)
	if err == nil {
		name = data.GetPath("name", "value").MustString()
	}

	c.services.Embed.NewEmbed().
		AddField("Color", name).
		AddField("RGB", fmt.Sprintf("%d, %d, %d", r, g, b)).
		AddField("Hex", "#"+hex).
		Send(ctx)

	preview, err := utils.Misc.PreviewColorImage(500, 100, r, g, b)
	if err == nil {
		_, _ = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "color.png", preview)
	}
}

func (FunCommands) pickupline(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.popcat.xyz/pickuplines")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Pickup Line", data.Get("pickupline").MustString()).
		Send(ctx)
}

func (FunCommands) elementinfo(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.popcat.xyz/periodic-table?element=" + url.QueryEscape(ctx.Arguments.Raw()))
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Element Name", data.Get("name").MustString()).
		AddField("Element Symbol", data.Get("symbol").MustString()).
		AddField("Atomic Number", fmt.Sprint(data.Get("atomic_number").MustInt())).
		AddField("Atomic Mass", fmt.Sprint(data.Get("atomic_mass").MustFloat64())).
		AddField("Period", fmt.Sprint(data.Get("period").MustInt())).
		AddField("Phase", data.Get("phase").MustString()).
		AddField("Discovered By", data.Get("discovered_by").MustString()).
		AddField("Summary", data.Get("summary").MustString()).
		SetImage(data.Get("image").MustString()).
		Send(ctx)
}

func (FunCommands) randomelement(ctx *router.Ctx) {
	data, err := http.GetJson("https://api.popcat.xyz/periodic-table/random")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	c.services.Embed.NewEmbed().
		AddField("Element Name", data.Get("name").MustString()).
		AddField("Element Symbol", data.Get("symbol").MustString()).
		AddField("Atomic Number", fmt.Sprint(data.Get("atomic_number").MustInt())).
		AddField("Atomic Mass", fmt.Sprint(data.Get("atomic_mass").MustFloat64())).
		AddField("Period", fmt.Sprint(data.Get("period").MustInt())).
		AddField("Phase", data.Get("phase").MustString()).
		AddField("Discovered By", data.Get("discovered_by").MustString()).
		AddField("Summary", data.Get("summary").MustString()).
		SetImage(data.Get("image").MustString()).
		Send(ctx)
}

type FunCommands struct{ services *services.Services }

func NewFunCommands(s *services.Services) *FunCommands {
	return &FunCommands{services: s}
}
func (c *FunCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

}
