package commands

import (
	"fmt"
	"math/rand"
	"net/url"

	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type MemeGenCommands struct{}

func (m *MemeGenCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Memegen", "Generate memes using imgflip")
	router.RegisterCmd(&dgc.Command{
		Name:        "randommeme",
		Description: "Send a random meme with custom text",
		Usage:       "[p]randommeme <\"top\"> <\"bottom\">",
		Handler:     m.randommeme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "custommeme",
		Description: "Send a meme based on a template",
		Usage:       "[p]custommeme <template id> <\"top\"> <\"bottom\">",
		Handler:     m.custommeme,
	})

	memes := [][]string{
		{"drakememe", "181913649", "Drake"},
		{"twobuttonsmeme", "87743020", "Two Buttons"},
		{"draw25cardsmeme", "217743513", "Draw 25 Cards"},
		{"turnrightmeme", "124822590", "Turn Right"},
		{"dogememe", "247375501", "Doge"},
		{"changemymindmeme", "129242436", "Change My Mind"},
		{"askingsupportmeme", "222403160", "Asking For Support"},
		{"skeletonmeme", "4087833", "Skeleton"},
		{"grumeme", "131940431", "Gru"},
		{"epichandshakememe", "135256802", "Epic Handshake"},
		{"disastergirlmeme", "97984", "Disaster Girl"},
		{"alwayshasbeenmeme", "252600902", "Always Has Been"},
		{"sadescobar", "80707627", "Sad Escobar"},
		{"safememe", "161865971", "Marked Safe"},
		{"everywherememe", "91538330", "Everywhere"},
		{"megamindmeme", "370867422", "Megamind"},
		{"winniememe", "178591752", "Winnie The Pooh"},
		{"spongebobmeme", "102156234", "Spongebob"},
		{"monkeypuppetmeme", "148909805", "Monkey Puppet"},
		{"tradeoffermeme", "309868304", "Trade Offer"},
		{"samepicturememe", "180190441", "Same Picture"},
		{"nutbuttonmeme", "119139145", "Nut Button"},
		{"haroldmeme", "27813981", "Harold"},
		{"notsimplymeme", "61579", "One Does Not Simply"},
		{"bikefallmeme", "79132341", "Bike Fall"},
		{"ancientaliensmeme", "101470", "Ancient Aliens"},
		{"nobitchesmeme", "371619279", "No Bitches"},
		{"thisisfinememe", "55311130", "This Is Fine"},
		{"wheremonkeymeme", "316466202", "Where Monkey"},
		{"anymorememe", "124055727", "Y'all got anymore"},
		{"pigeonmeme", "100777631", "Pigeon"},
		{"trophymeme", "3218037", "Trophy"},
		{"oprahmeme", "28251713", "Oprah"},
		{"scrolloftruthmeme", "123999232", "Scroll Of Truth"},
		{"sleepingshaqmeme", "99683372", "Sleeping Shaq"},
		{"thinkingmeme", "89370399", "Thinking"},
		{"pikachumeme", "155067746", "Surprised Pikachu"},
		{"laughingleomeme", "259237855", "Laughing Leo"},
		{"hannibalmeme", "135678846", "Hannibal"},
		{"flextapememe", "166969924", "Flex Tape"},
		{"therockmeme", "21735", "The Rock"},
		{"successkidmeme", "61544", "Success Kid"},
		{"evilkermitmeme", "84341851", "Evil Kermit"},
		{"drowningmeme", "252758727", "Drowning"},
		{"notthesamememe", "342785297", "Not The Same"},
		{"theydontknowmeme", "284929871", "They don't know"},
		{"trainbuscrashmeme", "247113703", "Train & Bus Crash"},
		{"gravememe", "221578498", "Grave"},
		{"unsettledtommeme", "175540452", "Unsettled Tom"},
		{"allhomieshatememe", "216523697", "All Homies Hate"},
		{"dissappointedmeme", "50421420", "Disappointed"},
		{"gettingpaidmeme", "177682295", "Getting Paid"},
		{"girlhidingmeme", "224514655", "Girl Hiding"},
		{"skepticalkidmeme", "101288", "Skeptical Kid"},
		{"grandmameme", "61556", "Grandma"},
		{"spongebobmeme2", "196652226", "Spongebob"},
		{"goosebumpsmeme", "101956210", "Goosebumps"},
		{"neverlandmeme", "6235864", "Neverland"},
		{"spidermanmeme", "110133729", "Spiderman"},
		{"hardpillmeme", "132769734", "Hard Pill"},
		{"futuramameme", "61520", "Futurama"},
		{"cheersmeme", "5496396", "Cheers"},
		{"squidwardmeme", "67452763", "Squidward"},
		{"besticandomeme", "77045868", "Best I Can Do"},
		{"dominomeme", "162372564", "Domino"},
		{"spongebobmeme3", "139971723", "Spongebob"},
		{"trumpmeme", "91545132", "Trump"},
		{"lookatmememe", "29617627", "Look At Me"},
		{"dogememe2", "8072285", "Doge"},
		{"dragonmeme", "187102311", "Dragon"},
		{"bugsbunnymeme", "249257686", "Bugs Bunny"},
		{"protectmeme", "171305372", "Protect"},
		{"powermeme", "177984372", "Feelings of Power"},
		{"undertakermeme", "234202281", "Undertaker"},
		{"greatpleasurememe", "313581692", "Great Pleasure"},
		{"isthisbutterflymeme", "142009471", "Is This Butterfly"},
		{"friendshipendedmeme", "137501417", "Friendship Ended"},
		{"simpsonsmeme", "72525473", "Simpsons"},
		{"dadmeme", "58045129", "Dad"},
		{"averageenjoyer", "316989420", "Average Enjoyer"},
		{"elmomeme", "140667789", "Elmo"},
		{"brainsleepmeme", "189601074", "Brain Before Sleep"},
		{"teacherscopymeme", "318009221", "Teacher's Copy"},
		{"theofficememe", "153695877", "The Office"},
		{"afraidtoaskmeme", "24557067", "Afraid to Ask"},
	}

	for _, meme := range memes {
		router.RegisterCmd(&dgc.Command{
			Name:        meme[0],
			Description: fmt.Sprintf("Send a %s meme with custom text", meme[2]),
			Usage:       fmt.Sprintf("[p]%s <\"top\"> <\"bottom\">", meme[0]),
			Handler:     m.generateFunc(meme[1]),
		})
	}
}

func (MemeGenCommands) randommeme(ctx *dgc.Ctx) {
	template, err := http.GetJson("https://api.imgflip.com/get_memes")
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	memes := template.Get("data").Get("memes")
	memesLen := len(memes.MustArray())

	// Get a random meme from the list
	meme := memes.GetIndex(rand.Intn(memesLen))
	templateId := meme.Get("id").MustString()

	uri := "https://api.imgflip.com/caption_image?template_id=%s&username=LightningBotSB2&password=LightningBotSB2&text0=%s&text1=%s"
	uri = fmt.Sprintf(uri, templateId, url.QueryEscape(utils.Arguments["<\"top\">"]), url.QueryEscape(utils.Arguments["<\"bottom\">"]))

	image, err := http.GetJson(uri)
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	embed.NewEmbed().
		AddField("Used Meme", meme.Get("name").MustString()).
		SetImage(image.GetPath("data", "url").MustString()).
		Send(ctx)
}

func (MemeGenCommands) custommeme(ctx *dgc.Ctx) {
	uri := "https://api.imgflip.com/caption_image?template_id=%s&username=LightningBotSB2&password=LightningBotSB2&text0=%s&text1=%s"
	uri = fmt.Sprintf(uri, url.QueryEscape(utils.Arguments["<template id>"]), url.QueryEscape(utils.Arguments["<\"top\">"]), url.QueryEscape(utils.Arguments["<\"bottom\">"]))

	image, err := http.GetJson(uri)
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}

	_ = ctx.RespondText(image.GetPath("data", "url").MustString())
}

func (MemeGenCommands) generateFunc(templateID string) func(ctx *dgc.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.

	return func(ctx *dgc.Ctx) {
		uri := "https://api.imgflip.com/caption_image?template_id=%s&username=LightningBotSB2&password=LightningBotSB2&text0=%s&text1=%s"
		uri = fmt.Sprintf(uri, templateID, url.QueryEscape(utils.Arguments["<\"top\">"]), url.QueryEscape(utils.Arguments["<\"bottom\">"]))

		image, err := http.GetJson(uri)
		if err != nil {
			_ = ctx.RespondText("An error has occurred.")
			return
		}

		_ = ctx.RespondText(image.GetPath("data", "url").MustString())
	}
}
