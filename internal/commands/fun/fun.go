package fun

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	owoifygo "github.com/deadshot465/owoify-go"
	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/formatting"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/random"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type FunCommands struct{ services *services.Services }

func NewFunCommands(s *services.Services) *FunCommands {
	return &FunCommands{services: s}
}

func (c *FunCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Fun",
		Description: "Fun commands",
	}, []*router.Command{
		{
			Name:        "eightball",
			Description: "Asks the 8-Ball a question",
			Arguments:   "<question>",
			Handler:     c.eightball,
		},
		{
			Name:        "howgay",
			Description: "See how gay a person is",
			Arguments:   "<user>",
			Handler:     c.howgay,
		},
		{
			Name:        "howsimp",
			Description: "See how simp a person is",
			Arguments:   "<user>",
			Handler:     c.howsimp,
		},
		{
			Name:        "penis",
			Description: "See how big a person's penis is",
			Arguments:   "<user>",
			Handler:     c.penis,
		},
		{
			Name:        "ego",
			Description: "See how big a person's ego is",
			Arguments:   "<user>",
			Handler:     c.ego,
		},
		{
			Name:        "iqtest",
			Description: "Test a person's IQ",
			Arguments:   "<user>",
			Handler:     c.iqtest,
		},
		{
			Name:        "hack",
			Description: "Fake hack a person",
			Arguments:   "<user>",
			Handler:     c.hack,
		},
		{
			Name:        "insult",
			Description: "Insult a person",
			Handler:     c.insult,
		},
		{
			Name:        "dadjoke",
			Description: "Send a random dad joke",
			Handler:     c.dadjoke,
		},
		{
			Name:        "yomomma",
			Description: "Send a random yo momma joke",
			Handler:     c.yomomma,
		},
		{
			Name:        "wouldyourather",
			Description: "Send a random would you rather question",
			Handler:     c.wouldyourather,
		},
		{
			Name:        "topic",
			Description: "Send a conversation topic",
			Handler:     c.topic,
		},
		{
			Name:        "slot",
			Description: "Play the slot machine",
			Handler:     c.slot,
		},
		{
			Name:        "advice",
			Description: "Send random advice",
			Handler:     c.advice,
		},
		{
			Name:        "catfact",
			Description: "Send a random cat fact",
			Handler:     c.catfact,
		},
		{
			Name:        "dogfact",
			Description: "Send a random dog fact",
			Handler:     c.dogfact,
		},
		{
			Name:        "pandafact",
			Description: "Send a random panda fact",
			Handler:     c.pandafact,
		},
		{
			Name:        "foxfact",
			Description: "Send a random fox fact",
			Handler:     c.foxfact,
		},
		{
			Name:        "koalafact",
			Description: "Send a random koala fact",
			Handler:     c.koalafact,
		},
		{
			Name:        "birdfact",
			Description: "Send a random bird fact",
			Handler:     c.birdfact,
		},
		{
			Name:        "redpandafact",
			Description: "Send a random red panda fact",
			Handler:     c.redpandafact,
		},
		{
			Name:        "whalefact",
			Description: "Send a random whale fact",
			Handler:     c.whalefact,
		},
		{
			Name:        "kangaroofact",
			Description: "Send a random whale fact",
			Handler:     c.kangaroofact,
		},
		{
			Name:        "bunnyfact",
			Description: "Send a random bunny fact",
			Handler:     c.bunnyfact,
		},
		{
			Name:        "lionfact",
			Description: "Send a random lion fact",
			Handler:     c.lionfact,
		},
		{
			Name:        "frogfact",
			Description: "Send a random frog fact",
			Handler:     c.frogfact,
		},
		{
			Name:        "duckfact",
			Description: "Send a random duck fact",
			Handler:     c.duckfact,
		},
		{
			Name:        "penguinfact",
			Description: "Send a random penguin fact",
			Handler:     c.penguinfact,
		},
		{
			Name:        "dolphinfact",
			Description: "Send a random dolphin fact",
			Handler:     c.dolphinfact,
		},
		{
			Name:        "bearfact",
			Description: "Send a random bear fact",
			Handler:     c.bearfact,
		},
		{
			Name:        "axolotlfact",
			Description: "Send a random axolotl fact",
			Handler:     c.axolotlfact,
		},
		{
			Name:        "capybarafact",
			Description: "Send a random capybara fact",
			Handler:     c.capybarafact,
		},
		{
			Name:        "trivia",
			Description: "Send a random trivia question",
			Handler:     c.trivia,
		},
		{
			Name:        "dice",
			Description: "Roll a dice",
			Handler:     c.dice,
		},
		{
			Name:        "randomemoji",
			Description: "Send a random server emoji",
			Arguments:   "[server]",
			Handler:     c.randomemoji,
		},
		{
			Name:        "randomchannel",
			Description: "Send a random channel from the server",
			Arguments:   "[server]",
			Handler:     c.randomchannel,
		},
		{
			Name:        "choose",
			Description: "Choose a random item from the list of items",
			Arguments:   "<option1> <option2> ...",
			Handler:     c.choose,
		},
		{
			Name:        "poll",
			Description: "Create a poll",
			Arguments:   "<question>",
			Handler:     c.poll,
		},
		{
			Name:        "numberfact",
			Description: "Send a random fact about a number",
			Arguments:   "<number>",
			Handler:     c.numberfact,
		},
		{
			Name:        "joke",
			Description: "Send a random joke",
			Handler:     c.joke,
		},
		{
			Name:        "fml",
			Description: "Send a random fuck my life quote",
			Handler:     c.fml,
		},
		{
			Name:        "tronald",
			Description: "Send a random Donald Trump quote",
			Handler:     c.tronald,
		},
		{
			Name:        "kanyewest",
			Description: "Send a random Kanye West quote",
			Handler:     c.kanyewest,
		},
		{
			Name:        "rps",
			Description: "Play rock paper scissors",
			Arguments:   "<rock/paper/scissors>",
			Handler:     c.rps,
		},
		{
			Name:        "covidtest",
			Description: "Test someone for Covid-19",
			Arguments:   "<user>",
			Handler:     c.covidtest,
		},
		{
			Name:        "hug",
			Description: "Send someone",
			Handler:     c.hug,
		},
		{
			Name:        "kiss",
			Description: "Kiss someone",
			Handler:     c.kiss,
		},
		{
			Name:        "slap",
			Description: "Slap someone",
			Handler:     c.slap,
		},
		{
			Name:        "cuddle",
			Description: "Cuddle someone",
			Handler:     c.cuddle,
		},
		{
			Name:        "pat",
			Description: "Pat someone",
			Handler:     c.pat,
		},
		{
			Name:        "feed",
			Description: "Feed someone",
			Handler:     c.feed,
		},
		{
			Name:        "wink",
			Description: "Wink to someone",
			Handler:     c.wink,
		},
		{
			Name:        "chucknorris",
			Description: "Send a random Chuck Norris joke",
			Handler:     c.chucknorris,
		},
		{
			Name:        "iss",
			Description: "Send the current ISS coordinates",
			Handler:     c.iss,
		},
		{
			Name:        "apod",
			Description: "Send NASA's Astronomy Picture Of The Day",
			Handler:     c.apod,
		},
		{
			Name:        "httpcat",
			Description: "Send a HTTP Cat image",
			Arguments:   "<statuscode>",
			Handler:     c.httpcat,
		},
		{
			Name:        "riddle",
			Description: "Send a random riddle",
			Handler:     c.riddle,
		},
		{
			Name:        "quote",
			Description: "Send a random quote",
			Handler:     c.quote,
		},
		{
			Name:        "evilinsult",
			Description: "Send a random evil insult",
			Handler:     c.evilinsult,
		},
		{
			Name:        "owoify",
			Description: "Owoify text",
			Arguments:   "<text>",
			Handler:     c.owoify,
		},
		{
			Name:        "uwuify",
			Description: "Uwuify text",
			Arguments:   "<text>",
			Handler:     c.uwuify,
		},
		{
			Name:        "uvuify",
			Description: "Uvuify text",
			Arguments:   "<text>",
			Handler:     c.uvuify,
		},
		{
			Name:        "rickrollurl",
			Description: "Send the rick roll url",
			Handler:     c.rickrollurl,
		},
		{
			Name:        "rickroll",
			Description: "Send the rick roll lyrics",
			Handler:     c.rickroll,
		},
		{
			Name:        "coinflip",
			Description: "Flip a coin",
			Handler:     c.coinflip,
		},
		{
			Name:        "notfunny",
			Description: "Send a not funny copypasta",
			Handler:     c.notfunny,
		},
		{
			Name:        "randomquestion",
			Description: "Send a random question",
			Handler:     c.randomquestion,
		},
		{
			Name:        "lyrics",
			Description: "Get the lyrics of a song",
			Arguments:   "<song name>",
			Handler:     c.lyrics,
		},
		{
			Name:        "breakingbadquote",
			Description: "Send a random Breaking Badquote",
			Handler:     c.breakingbadquote,
		},
		{
			Name:        "gameofthronesquote",
			Description: "Send a random Game of Thrones quote",
			Handler:     c.gameofthronesquote,
		},
		{
			Name:        "luciferquote",
			Description: "Send a random Lucifer quote",
			Handler:     c.luciferquote,
		},
		{
			Name:        "strangerthingsquote",
			Description: "Send a random Stranger Things quote",
			Handler:     c.strangerthingsquote,
		},
		{
			Name:        "xkcd",
			Description: "Send the current XKCD comic",
			Handler:     c.xkcd,
		},
		{
			Name:        "boredactivity",
			Description: "Send a random activity to do when bored",
			Handler:     c.boredactivity,
		},
		{
			Name:        "ermahgerd",
			Description: "Translate text to ERMAHGERD",
			Arguments:   "<text>",
			Handler:     c.ermahgerd,
		},
		{
			Name:        "iseven",
			Description: "Check if a number is even",
			Arguments:   "<number>",
			Handler:     c.iseven,
		},
		{
			Name:        "agify",
			Description: "Determine the age of a first name",
			Arguments:   "<name>",
			Handler:     c.agify,
		},
		{
			Name:        "randomword",
			Description: "Get a random word",
			Handler:     c.randomword,
		},
		{
			Name:        "randomweirdword",
			Description: "Get a random weird word",
			Handler:     c.randomweirdword,
		},
		{
			Name:        "startupidea",
			Description: "Get a random startup idea",
			Handler:     c.startupidea,
		},
		{
			Name:        "fakeperson",
			Description: "Generate an image of a fake person",
			Handler:     c.fakeperson,
		},
		{
			Name:        "commitmessage",
			Description: "Generate a random commit message",
			Handler:     c.commitmessage,
		},
		{
			Name:        "showerthought",
			Description: "Send a random shower thought",
			Handler:     c.showerthought,
		},
		{
			Name:        "truth",
			Description: "Play Truth or Dare",
			Handler:     c.truth,
		},
		{
			Name:        "dare",
			Description: "Play Truth or Dare",
			Handler:     c.dare,
		},
		{
			Name:        "neverhaveiever",
			Description: "Play Never have I ever...",
			Handler:     c.neverhaveiever,
		},
		{
			Name:        "paranoia",
			Description: "Send a random paranoia question",
			Handler:     c.paranoia,
		},
		{
			Name:        "randomcolor",
			Description: "Generate a random color",
			Handler:     c.randomcolor,
		},
		{
			Name:        "pickupline",
			Description: "Generate a pickup line",
			Handler:     c.pickupline,
		},
		{
			Name:        "elementinfo",
			Description: "Get info about a periodic table element",
			Arguments:   "<element>",
			Handler:     c.elementinfo,
		},
		{
			Name:        "randomelement",
			Description: "Get info about a random periodic table element",
			Handler:     c.randomelement,
		},
	})
}

func (c *FunCommands) eightball(ctx *router.Ctx) error {
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
	c.services.Embed.New(ctx).
		AddField("Question", question).
		AddField("Answer", answer).
		Send()
	return nil
}

func (c *FunCommands) howgay(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("Gay Percentage", random.Percentage()).
		Send()
	return nil
}

func (c *FunCommands) howsimp(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("Simp Percentage", random.Percentage()).
		Send()
	return nil
}

func (c *FunCommands) penis(ctx *router.Ctx) error {
	penisLength := random.Int(1, 20)
	penis := "8"
	for range penisLength {
		penis += "="
	}
	penis += "Đ"
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("Penis", penis).
		AddField("Penis Length", fmt.Sprintf("%d cm", penisLength)).
		Send()
	return nil
}

func (c *FunCommands) ego(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("Ego", random.Percentage()).
		Send()
	return nil
}

func (c *FunCommands) iqtest(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("IQ", fmt.Sprint(random.Int(1, 200))).
		Send()
	return nil
}

func (c *FunCommands) hack(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}

	resp, err := http.DefaultClient().Get("https://randomuser.me/api/")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}

	data := resp.Json().Get("results", 0)
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("Real Name", data.Get("name", "first").MustString()+" "+data.Get("name", "last").MustString()).
		AddField("Email", strings.ReplaceAll(data.Get("email").MustString(), "example", "gmail")).
		AddField("Phone", data.Get("cell").MustString()).
		AddField("Address", strconv.Itoa(data.Get("location", "street", "number").MustInt())+" "+data.Get("location", "street", "name").MustString()+", "+data.Get("location", "city").MustString()+", "+data.Get("location", "country").MustString()).
		AddField(data.Get("id", "name").MustString(), data.Get("id", "value").MustString()).
		Send()
	return nil
}

func (c *FunCommands) insult(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://insult.mattbas.org/api/insult")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Insult", string(resp.MustGetBody()))
	return nil
}

func (c *FunCommands) dadjoke(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://icanhazdadjoke.com/", http.WithHeader("Accept", "application/json"))
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Dad Joke", resp.Json().Get("joke").MustString())
	return nil
}

func (c *FunCommands) yomomma(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.yomomma.info/")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Yo Momma Joke", resp.Json().Get("joke").MustString())
	return nil
}

func (c *FunCommands) wouldyourather(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/wyr")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		SetDescription("Would you rather...").
		AddField("A", data.Get("question", "a").MustString()).
		AddField("B", data.Get("question", "b").MustString()).
		Send()
	return nil
}

func (c *FunCommands) topic(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/topic")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Topic", resp.Json().Get("topic").MustString())
	return nil
}

func (c *FunCommands) slot(ctx *router.Ctx) error {
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

	cfg := c.services.Config.GetConfig()

	getEmbed := func(final bool) string {
		text, winning := getSlot()
		if final {
			text += "\n\n" + winning
		}
		return c.services.Embed.New(ctx).
			SetDescription("%s", text).
			ToString(cfg.Selfbot.CommandMode, 2000)[0]
	}

	message, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, getEmbed(false))
	if err != nil {
		return err
	}

	for range 5 {
		time.Sleep(500 * time.Millisecond)

		_, err = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, getEmbed(false))
		if err != nil {
			return err
		}
	}

	_, err = ctx.Session.ChannelMessageEdit(message.ChannelID, message.ID, getEmbed(true))
	return err
}

func (c *FunCommands) advice(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.adviceslip.com/advice")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Advice", resp.Json().Get("slip", "advice").MustString())
	return nil
}

func (c *FunCommands) catfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://some-random-api.com/facts/cat")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Cat Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) dogfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://some-random-api.com/facts/dog")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Dog Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) pandafact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://some-random-api.com/facts/panda")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Panda Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) foxfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://some-random-api.com/facts/fox")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Fox Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) koalafact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://some-random-api.com/facts/koala")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Koala Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) birdfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/bird")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Bird Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) redpandafact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/redpanda")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Red Panda Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) whalefact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/whale")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Whale Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) kangaroofact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/kangaroo")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Kangaroo Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) bunnyfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/bunny")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Bunny Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) lionfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/lion")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Lion Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) frogfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/frog")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Frog Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) duckfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/duck")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Duck Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) penguinfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/penguin")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Penguin Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) dolphinfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/dolphin")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Dolphin Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) bearfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/bear")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Bear Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) axolotlfact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/axolotl")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Axolotl Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) capybarafact(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.animality.xyz/fact/capybara")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Capybara Fact", resp.Json().Get("fact").MustString())
	return nil
}

func (c *FunCommands) trivia(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://opentdb.com/api.php?amount=1")
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("An error has occurred.")
	}

	data := resp.Json()

	if data.Get("response_code").MustInt64() != 0 {
		return fmt.Errorf("An error has occurred.")
	}

	data = data.Get("results").Get(0)

	answers := data.Get("incorrect_answers").MustStringSlice()
	answers = append(answers, data.Get("correct_answer").MustString())
	answers = random.Shuffle(answers)

	c.services.Embed.New(ctx).
		AddField("Question", html.UnescapeString(data.Get("question").MustString())).
		AddField("Category", data.Get("category").MustString()).
		AddField("Difficulty", utils.String.Capitalize(data.Get("difficulty").MustString())).
		AddField("Answers", html.UnescapeString(strings.Join(answers, ", "))).
		AddField("Correct Answer", fmt.Sprintf("||%s||", html.UnescapeString(data.Get("correct_answer").MustString()))).
		Send()
	return nil
}

func (c *FunCommands) dice(ctx *router.Ctx) error {
	c.services.Embed.Respond(ctx, "You rolled: %d", random.Int(1, 6))
	return nil
}

func (c *FunCommands) randomemoji(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find server!")
	}

	emojis, err := ctx.Session.GuildEmojis(guild.ID)
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	emoji := random.Choice(emojis)
	return ctx.RespondText("%s", emoji.MessageFormat())
}

func (c *FunCommands) randomchannel(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find server!")
	}

	channels, err := ctx.Session.GuildChannels(guild.ID)
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	channels = fn.Filter(channels, func(channel *discordgo.Channel) bool {
		return channel.Type == discordgo.ChannelTypeGuildText
	})

	channel := random.Choice(channels)
	return ctx.RespondText("%s", channel.Mention())
}

func (c *FunCommands) choose(ctx *router.Ctx) error {
	choice := random.Choice(strings.Split(ctx.Arguments.Raw(), " "))
	c.services.Embed.Respond(ctx, "I choose: %s", choice)
	return nil
}

func (c *FunCommands) poll(ctx *router.Ctx) error {
	messages := c.services.Embed.New(ctx).
		SetDescription("Answer to this poll by reacting with :thumbsup: or :thumbsdown:").
		AddField("Question", ctx.Arguments.Raw()).
		Send()

	if len(messages) > 0 {
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, messages[0].ID, "👍")
		_ = ctx.Session.MessageReactionAdd(ctx.Message.ChannelID, messages[0].ID, "👎")
	}

	return nil
}

func (c *FunCommands) numberfact(ctx *router.Ctx) error {
	number := ctx.Arguments.Get(0).Raw()
	resp, err := http.DefaultClient().Get(fmt.Sprintf("http://numbersapi.com/%s", number))
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.New(ctx).
		AddField("Number", number).
		AddField("Fact", string(resp.MustGetBody())).
		Send()
	return nil
}

func (c *FunCommands) joke(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://v2.jokeapi.dev/joke/Any")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	emb := c.services.Embed.New(ctx)
	if data.Get("type").MustString() == "twopart" {
		emb.AddField("Setup", data.Get("setup").MustString()).
			AddField("Delivery", data.Get("delivery").MustString())
	} else {
		emb.AddField("Joke", data.Get("joke").MustString())
	}
	emb.AddField("Category", data.Get("category").MustString()).Send()
	return nil
}

func (c *FunCommands) fml(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/fml")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "FML", resp.Json().Get("fml").MustString())
	return nil
}

func (c *FunCommands) tronald(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.tronalddump.io/random/quote", http.WithHeader("Accept", "application/hal+json"))
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Tronald Dump Quote", resp.Json().Get("value").MustString())
	return nil
}

func (c *FunCommands) kanyewest(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.kanye.rest")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Kanye West Quote", resp.Json().Get("quote").MustString())
	return nil
}

func (c *FunCommands) rps(ctx *router.Ctx) error {
	choice := ctx.Arguments.Get(0).Raw()
	if choice != "rock" && choice != "paper" && choice != "scissors" {
		return fmt.Errorf("You can only choose rock, paper or scissors.")
	}
	c.services.Embed.New(ctx).
		AddField("Your choice", choice).
		AddField("Bot's choice", random.Choice([]string{"rock", "paper", "scissors"})).
		Send()
	return nil
}

func (c *FunCommands) covidtest(ctx *router.Ctx) error {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("Couldn't find user!")
	}
	c.services.Embed.New(ctx).
		AddField("User", user.String()).
		AddField("Result", random.Choice([]string{"Positive", "Negative"})).
		Send()
	return nil
}

func (c *FunCommands) hug(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/hug", "url")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) kiss(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/kiss", "url")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) slap(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/slap", "url")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) cuddle(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/cuddle", "url")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) pat(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/pat", "url")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) feed(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://nekos.life/api/v2/img/feed", "url")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) wink(ctx *router.Ctx) error {
	url, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/animu/wink", "link")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	return ctx.RespondText("%s", url)
}

func (c *FunCommands) chucknorris(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Chuck Norris Joke", resp.Json().Get("value").MustString())
	return nil
}

func (c *FunCommands) iss(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("http://api.open-notify.org/iss-now.json")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	latitude := data.Get("iss_position", "latitude").MustString()
	longitude := data.Get("iss_position", "longitude").MustString()
	c.services.Embed.New(ctx).
		SetTitle("International Space Station Info").
		AddField("Latitude", latitude).
		AddField("Longitude", longitude).
		AddField("Google Maps URL", fmt.Sprintf("https://maps.google.com/?q=%s,%s", latitude, longitude)).
		Send()
	return nil
}

func (c *FunCommands) apod(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.nasa.gov/planetary/apod?api_key=rDenYNQFftu3YADWdUmF9s8iw0cyZckV7vBcibGe")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	copyright := data.Get("copyright").MustString()
	if copyright == "" {
		copyright = "NASA"
	}
	explanation, err := api.LBCreatePaste(data.Get("explanation").MustString())
	if err != nil {
		explanation = "An error has occurred."
	}
	c.services.Embed.New(ctx).
		SetImage(data.Get("url").MustString()).
		SetTitle("Astronomy Picture of the Day - ©️%s", copyright).
		AddField("Title", data.Get("title").MustString()).
		AddField("Explanation", explanation).
		AddField("Date", data.Get("date").MustString()).
		Send()
	return nil
}

func (c *FunCommands) httpcat(ctx *router.Ctx) error {
	statuscode := ctx.Arguments.Get(0).Raw()
	c.services.Embed.New(ctx).
		SetImage(fmt.Sprintf("https://http.cat/%s.jpg", statuscode)).
		Send()
	return nil
}

func (c *FunCommands) riddle(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/riddle")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Riddle", data.Get("question").MustString()).
		AddField("Answer", "||"+data.Get("answer").MustString()+"||").
		Send()
	return nil
}

func (c *FunCommands) quote(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.quotable.io/random")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Quote", data.Get("content").MustString()).
		AddField("Author", data.Get("author").MustString()).
		Send()
	return nil
}

func (c *FunCommands) evilinsult(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://evilinsult.com/generate_insult.php?lang=en")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Evil Insult", string(resp.MustGetBody()))
	return nil
}

func (c *FunCommands) owoify(ctx *router.Ctx) error {
	return ctx.RespondText("%s", owoifygo.Owoify(ctx.Arguments.Raw(), "owo"))
}

func (c *FunCommands) uwuify(ctx *router.Ctx) error {
	return ctx.RespondText("%s", owoifygo.Owoify(ctx.Arguments.Raw(), "uwu"))
}

func (c *FunCommands) uvuify(ctx *router.Ctx) error {
	return ctx.RespondText("%s", owoifygo.Owoify(ctx.Arguments.Raw(), "uvu"))
}

func (c *FunCommands) rickrollurl(ctx *router.Ctx) error {
	return ctx.RespondText("https://www.youtube.com/watch?v=dQw4w9WgXcQ")
}

func (c *FunCommands) rickroll(ctx *router.Ctx) error {
	return ctx.RespondText("**We're no strangers to love\nYou know the rules and so do I\nA full commitment's what I'm thinking of\nYou wouldn't get this from any other guy\nI just wanna tell you how I'm feeling\nGotta make you understand\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nWe've known each other for so long\nYour heart's been aching but you're too shy to say it\nInside we both know what's been going on\nWe know the game and we're gonna play it\nAnd if you ask me how I'm feeling\nDon't tell me you're too blind to see\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give, never gonna give\n(Give you up)\n(Ooh) Never gonna give, never gonna give\n(Give you up)\nWe've known each other for so long\nYour heart's been aching but you're too shy to say it\nInside we both know what's been going on\nWe know the game and we're gonna play it\nI just wanna tell you how I'm feeling\nGotta make you understand\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry\nNever gonna say goodbye\nNever gonna tell a lie and hurt you\nNever gonna give you up\nNever gonna let you down\nNever gonna run around and desert you\nNever gonna make you cry**")
}

func (c *FunCommands) coinflip(ctx *router.Ctx) error {
	choice := random.Choice([]string{"Heads!", "Tails!"})
	images := map[string]string{
		"Heads!": "https://www.ssaurel.com/blog/wp-content/uploads/2017/01/heads.png",
		"Tails!": "https://www.ssaurel.com/blog/wp-content/uploads/2017/01/tails.png",
	}
	c.services.Embed.New(ctx).
		SetDescription("%s", choice).
		SetImage(images[choice]).
		Send()
	return nil
}

func (c *FunCommands) notfunny(ctx *router.Ctx) error {
	ctx.RespondText("**Not funny I didn't laugh. Your joke is so bad I would have preferred the joke went over my head and you gave up re-telling me the joke. To be honest this is a horrid attempt at trying to get a laugh out of me. Not a chuckle, not a hehe, not even a subtle burst of air out of my esophagus. Science says before you laugh your brain preps your face muscles but I didn't even feel the slightest twitch. 0/10 this joke is so bad I cannot believe anyone legally allowed you to be creative at all. The amount of brain power you must have put into that joke has the potential to power every house on Earth. Get a personality and learn how to make jokes, read a book. I'm not saying this to be funny I genuinely mean it on how this is just bottom barrel embarrassment at comedy. You've single-handedly killed humor and every comedic act on the planet. I'm so disappointed that society has failed as a whole in being able to teach you how to be funny. Honestly if I put in all my power and time to try and make your joke funny it would require Einstein himself to build a device to strap me into so I can be connected to the energy of a billion stars to do it, and even then all that joke would get from people is a subtle scuff. You're lucky I still have the slightest of empathy for you after telling that joke otherwise I would have committed every war crime in the book just to prevent you from attempting any humor ever again. We should put that joke in text books so future generations can be wary of becoming such an absolute comedic failure. Im disappointed, hurt, and outright offended that my precious time has been wasted in my brain understanding that joke. In the time that took I was planning on helping kids who have been orphaned, but because of that you've waisted my time explaining the obscene integrity of your terrible attempt at comedy. Now those kids are suffering without meals and there's nobody to blame but you.**")
	ctx.RespondText("**I hope you're happy with what you have done and I truly hope you can move on and learn from this piss poor attempt**")
	return nil
}

func (c *FunCommands) randomquestion(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://nekos.life/api/v2/why")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Question", resp.Json().Get("why").MustString())
	return nil
}

func (c *FunCommands) lyrics(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://some-random-api.com/lyrics?title=" + ctx.Arguments.Escaped())
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	if data.Get("error").MustString() != "" {
		return fmt.Errorf("Could not find lyrics.")
	}
	c.services.Embed.New(ctx).
		AddField("Artist", data.Get("author").MustString()).
		AddField("Title", data.Get("title").MustString()).
		AddField("Lyrics", data.Get("lyrics").MustString()).
		Send()
	return nil
}

func (c *FunCommands) breakingbadquote(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/breaking-bad-quote")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Quote", data.Get("quote").MustString()).
		AddField("Author", data.Get("author").MustString()).
		Send()
	return nil
}

func (c *FunCommands) gameofthronesquote(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/game-of-thrones-quote")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Quote", data.Get("quote").MustString()).
		AddField("Author", data.Get("author").MustString()).
		Send()
	return nil
}

func (c *FunCommands) strangerthingsquote(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/stranger-things-quote")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Quote", data.Get("quote").MustString()).
		AddField("Author", data.Get("author").MustString()).
		Send()
	return nil
}

func (c *FunCommands) luciferquote(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.lightning-bot.com/api/v3/misc/lucifer-quote")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Quote", data.Get("quote").MustString()).
		AddField("Author", data.Get("author").MustString()).
		Send()
	return nil
}

func (c *FunCommands) xkcd(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://xkcd.com/info.0.json")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()

	date := data.Get("year").MustString() + "-" + data.Get("month").MustString() + "-" + data.Get("day").MustString()

	c.services.Embed.New(ctx).
		AddField("Title", data.Get("title").MustString()).
		AddField("Date", date).
		SetImage(data.Get("img").MustString()).
		Send()
	return nil
}

func (c *FunCommands) boredactivity(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://www.boredapi.com/api/activity")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Activity", data.Get("activity").MustString()).
		AddField("Type", utils.String.Capitalize(data.Get("type").MustString())).
		SetImage(data.Get("img").MustString()).
		Send()
	return nil
}

func (c *FunCommands) ermahgerd(ctx *router.Ctx) error {
	return ctx.RespondText("%s", utils.String.Ermahgerd(ctx.Arguments.Raw()))
}

func (c *FunCommands) iseven(ctx *router.Ctx) error {
	number, err := ctx.Arguments.Get(0).Int()
	if err != nil {
		return fmt.Errorf("Argument <number> must be a number.")
	}

	c.services.Embed.New(ctx).
		AddField("Number", ctx.Arguments.Get(0).Raw()).
		AddField("Is Even", formatting.BoolYesNo(number%2 == 0)).
		Send()
	return nil
}

func (c *FunCommands) agify(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.agify.io?name=" + ctx.Arguments.Escaped())
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Name", data.Get("name").MustString()).
		AddField("Age", fmt.Sprint(data.Get("age").MustInt())).
		AddField("Count", fmt.Sprint(data.Get("count").MustInt())).
		Send()
	return nil
}

func (c *FunCommands) randomword(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://randomwordgenerator.com/json/words_ws.json")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}

	words := resp.Json().Get("data")
	word := words.Get(random.Int(0, words.Len())).Get("word")

	c.services.Embed.RespondField(ctx, "Random Word", word.Get("value").MustString())
	return nil
}

func (c *FunCommands) randomweirdword(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://randomwordgenerator.com/json/unusual-words.json")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}

	words := resp.Json().Get("data")
	word := words.Get(random.Int(0, words.Len())).Get("word")

	c.services.Embed.New(ctx).
		AddField("Random Weird Word", word.Get("value").MustString()).
		AddField("Definition", word.Get("definition").MustString()).
		Send()
	return nil
}

func (c *FunCommands) startupidea(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://itsthisforthat.com/api.php?text")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Startup Idea", string(resp.MustGetBody()))
	return nil
}

func (c *FunCommands) fakeperson(ctx *router.Ctx) error {
	c.services.Embed.New(ctx).
		SetDescription("Fake Person").
		SetImage("https://thispersondoesnotexist.com?" + fmt.Sprint(random.Int(1000, 9999))).
		Send()
	return nil
}

func (c *FunCommands) commitmessage(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://whatthecommit.com/index.txt")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Commit Message", string(resp.MustGetBody()))
	return nil
}

func (c *FunCommands) showerthought(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://showerthoughts.nickf.me/api/thoughts/random?limit=1")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Shower Thought", resp.Json().Get("title").MustString())
	return nil
}

func (c *FunCommands) truth(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.truthordarebot.xyz/v1/truth")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Truth Question", resp.Json().Get("question").MustString())
	return nil
}

func (c *FunCommands) dare(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.truthordarebot.xyz/v1/dare")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Dare Question", resp.Json().Get("question").MustString())
	return nil
}

func (c *FunCommands) neverhaveiever(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.truthordarebot.xyz/v1/nhie")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Never have I ever...", resp.Json().Get("question").MustString())
	return nil
}

func (c *FunCommands) paranoia(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.truthordarebot.xyz/v1/paranoia")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Paranoia Question", resp.Json().Get("question").MustString())
	return nil
}

func (c *FunCommands) randomcolor(ctx *router.Ctx) error {
	r, g, b := random.Int(0, 255), random.Int(0, 255), random.Int(0, 255)
	hex := fmt.Sprintf("%02x%02x%02x", r, g, b)
	name := "Unknown"

	resp, err := http.DefaultClient().Get("https://www.thecolorapi.com/id?hex=" + hex)
	if err == nil {
		name = resp.Json().Get("name", "value").MustString()
	}

	c.services.Embed.New(ctx).
		AddField("Color", name).
		AddField("RGB", fmt.Sprintf("%d, %d, %d", r, g, b)).
		AddField("Hex", "#"+hex).
		Send()

	preview, err := utils.Misc.PreviewColorImage(500, 100, r, g, b)
	if err != nil {
		return fmt.Errorf("failed to create preview: %v", err)
	}

	_, err = ctx.Session.ChannelFileSend(ctx.Message.ChannelID, "color.png", preview)
	if err != nil {
		return fmt.Errorf("failed to send preview: %v", err)
	}

	return nil
}

func (c *FunCommands) pickupline(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.popcat.xyz/pickuplines")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	c.services.Embed.RespondField(ctx, "Pickup Line", resp.Json().Get("pickupline").MustString())
	return nil
}

func (c *FunCommands) elementinfo(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.popcat.xyz/periodic-table?element=" + ctx.Arguments.Escaped())
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Element Name", data.Get("name").MustString()).
		AddField("Element Symbol", data.Get("symbol").MustString()).
		AddField("Atomic Number", fmt.Sprint(data.Get("atomic_number").MustInt())).
		AddField("Atomic Mass", fmt.Sprint(data.Get("atomic_mass").MustFloat64())).
		AddField("Period", fmt.Sprint(data.Get("period").MustInt())).
		AddField("Phase", data.Get("phase").MustString()).
		AddField("Discovered By", data.Get("discovered_by").MustString()).
		AddField("Summary", data.Get("summary").MustString()).
		SetImage(data.Get("image").MustString()).
		Send()
	return nil
}

func (c *FunCommands) randomelement(ctx *router.Ctx) error {
	resp, err := http.DefaultClient().Get("https://api.popcat.xyz/periodic-table/random")
	if err != nil {
		return fmt.Errorf("An error has occurred.")
	}
	data := resp.Json()
	c.services.Embed.New(ctx).
		AddField("Element Name", data.Get("name").MustString()).
		AddField("Element Symbol", data.Get("symbol").MustString()).
		AddField("Atomic Number", fmt.Sprint(data.Get("atomic_number").MustInt())).
		AddField("Atomic Mass", fmt.Sprint(data.Get("atomic_mass").MustFloat64())).
		AddField("Period", fmt.Sprint(data.Get("period").MustInt())).
		AddField("Phase", data.Get("phase").MustString()).
		AddField("Discovered By", data.Get("discovered_by").MustString()).
		AddField("Summary", data.Get("summary").MustString()).
		SetImage(data.Get("image").MustString()).
		Send()
	return nil
}
