package image

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/random"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/internal/utils2"
)

type ImageCommands struct{ services *services.Services }

func NewImageCommands(s *services.Services) *ImageCommands {
	return &ImageCommands{services: s}
}
func (c *ImageCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Image",
		Description: "Image commands",
	}, []*router.Command{
		{
			Name:        "tweet",
			Description: "Tweet a message",
			Arguments:   "<username> <text>",
			Handler:     c.tweet,
		},
		{
			Name:        "meme",
			Description: "Send a random meme",
			Handler:     c.meme,
		},
		{
			Name:        "ytcomment",
			Description: "Make a fake YouTube comment",
			Arguments:   "<user> <text>",
			Handler:     c.ytcomment,
		},
		{
			Name:        "pokemon",
			Description: "Send a random Pokémon image",
			Handler:     c.pokemon,
		},
		{
			Name:        "wasted",
			Description: "Add a wasted overlay on the users avatar",
			Arguments:   "[user]",
			Handler:     c.wasted,
		},
		{
			Name:        "triggered",
			Description: "Add a triggered overlay on the users avatar",
			Arguments:   "[user]",
			Handler:     c.triggered,
		},
		{
			Name:        "glass",
			Description: "Add a glass overlay on the users avatar",
			Arguments:   "[user]",
			Handler:     c.glass,
		},
		{
			Name:        "greyscale",
			Description: "Add a greyscale overlay on the users avatar",
			Arguments:   "[user]",
			Handler:     c.greyscale,
		},
		{
			Name:        "invert",
			Description: "Invert the users avatar",
			Arguments:   "[user]",
			Handler:     c.invert,
		},
		{
			Name:        "invertgreyscale",
			Description: "Invert and greyscale the users avatar",
			Arguments:   "[user]",
			Handler:     c.invertgreyscale,
		},
		{
			Name:        "brightness",
			Description: "Make the users avatar brighter",
			Arguments:   "[user]",
			Handler:     c.brightness,
		},
		{
			Name:        "red",
			Description: "Make the users avatar red",
			Arguments:   "[user]",
			Handler:     c.red,
		},
		{
			Name:        "green",
			Description: "Make the users avatar green",
			Arguments:   "[user]",
			Handler:     c.green,
		},
		{
			Name:        "blue",
			Description: "Make the users avatar blue",
			Arguments:   "[user]",
			Handler:     c.blue,
		},
		{
			Name:        "pixelate",
			Description: "Pixelate the users avatar",
			Arguments:   "[user]",
			Handler:     c.pixelate,
		},
		{
			Name:        "comrade",
			Description: "Add comrade overlay on users avatar",
			Arguments:   "[user]",
			Handler:     c.comrade,
		},
		{
			Name:        "missionpassed",
			Description: "Add Mission Passed overlay on users avatar",
			Arguments:   "[user]",
			Handler:     c.missionpassed,
		},
		{
			Name:        "jail",
			Description: "Add jail overlay on users avatar",
			Arguments:   "[user]",
			Handler:     c.jail,
		},
		{
			Name:        "nobitches",
			Description: "Make a no bitches meme",
			Arguments:   "<text>",
			Handler:     c.nobitches,
		},
		{
			Name:        "threat",
			Description: "Make biggest threats meme",
			Arguments:   "[user]",
			Handler:     c.threat,
		},
		{
			Name:        "threshold",
			Description: "Add a threshold filter to a users avatar",
			Arguments:   "<user>",
			Handler:     c.threshold,
		},
		{
			Name:        "blurple",
			Description: "Make the users avatar blurple",
			Handler:     c.blurple,
		},
		{
			Name:        "blurple2",
			Description: "Make the users avatar new blurple",
			Handler:     c.blurple2,
		},
		{
			Name:        "color",
			Description: "Make the users avatar a custom color",
			Arguments:   "<user> <color>",
			Handler:     c.color,
		},
		{
			Name:        "sostupid",
			Description: "Make an it's so stupid meme",
			Arguments:   "<user> <text>",
			Handler:     c.sostupid,
		},
		{
			Name:        "simpcard",
			Description: "Make a simp card for a user",
			Arguments:   "<user>",
			Handler:     c.simpcard,
		},
		{
			Name:        "horny",
			Description: "Make a horny card for a user",
			Arguments:   "<user>",
			Handler:     c.horny,
		},
		{
			Name:        "lgbt",
			Description: "Add an lgbt border to the user",
			Arguments:   "<user>",
			Handler:     c.lgbt,
		},
		{
			Name:        "pansexual",
			Description: "Add a pansexual border to the user",
			Arguments:   "<user>",
			Handler:     c.pansexual,
		},
		{
			Name:        "nonbinary",
			Description: "Add a nonbinary border to the user",
			Arguments:   "<user>",
			Handler:     c.nonbinary,
		},
		{
			Name:        "lesbian2",
			Description: "Add a lesbian border to the user",
			Arguments:   "<user>",
			Handler:     c.lesbian,
		},
		{
			Name:        "bisexual",
			Description: "Add a bisexual border to the user",
			Arguments:   "<user>",
			Handler:     c.bisexual,
		},
		{
			Name:        "trans",
			Description: "Add a trans border to the user",
			Arguments:   "<user>",
			Handler:     c.trans,
		},
		{
			Name:        "oogway",
			Description: "Make a master oogway meme",
			Arguments:   "<text>",
			Handler:     c.oogway,
		},
		{
			Name:        "oogway2",
			Description: "Make a master oogway meme",
			Arguments:   "<text>",
			Handler:     c.oogway2,
		},
		{
			Name:        "clyde",
			Description: "Make Clyde say something",
			Arguments:   "<text>",
			Handler:     c.clyde,
		},
		{
			Name:        "gay",
			Description: "Add a rainbow overlay on the users avatar",
			Arguments:   "<user>",
			Handler:     c.gay,
		},
		{
			Name:        "distractedbf",
			Description: "Make a distracted boyfriend meme with the users",
			Arguments:   "<user1> <user2>",
			Handler:     c.distractedbf,
		},
		{
			Name:        "captcha",
			Description: "Make a captcha meme",
			Arguments:   "<user> <text>",
			Handler:     c.captcha,
		},
		{
			Name:        "whowouldwin",
			Description: "Make a who would win meme",
			Arguments:   "<user1> <user2>",
			Handler:     c.whowouldwin,
		},
		{
			Name:        "jpeg",
			Description: "JPEG a users avatar",
			Arguments:   "<user>",
			Handler:     c.jpeg,
		},
		{
			Name:        "lolice",
			Description: "Make a lolice meme",
			Arguments:   "<user>",
			Handler:     c.lolice,
		},
		{
			Name:        "kannagen",
			Description: "Make a kannagen meme",
			Arguments:   "<user1> <user2>",
			Handler:     c.kannagen,
		},
		{
			Name:        "iphonex",
			Description: "Put users profile picture on an iPhone X",
			Arguments:   "<user1> <user2>",
			Handler:     c.iphonex,
		},
		{
			Name:        "awooify",
			Description: "Awooify the user",
			Arguments:   "<user1> <user2>",
			Handler:     c.awooify,
		},
		{
			Name:        "trap",
			Description: "Make a trap card",
			Arguments:   "<user>",
			Handler:     c.trap,
		},
		{
			Name:        "trumptweet",
			Description: "Make a fake Trump tweet",
			Arguments:   "<text>",
			Handler:     c.trumptweet,
		},
		{
			Name:        "deepfry",
			Description: "Deepfry the users avatar",
			Arguments:   "<text>",
			Handler:     c.deepfry,
		},
		{
			Name:        "blurpify",
			Description: "Blurpify the users avatar",
			Arguments:   "<text>",
			Handler:     c.blurpify,
		},
		{
			Name:        "phcomment",
			Description: "Make the user comment something on PornHub",
			Arguments:   "<user> <text>",
			Handler:     c.phcomment,
		},
		{
			Name:        "fact",
			Description: "Make a fact meme",
			Arguments:   "<user>",
			Handler:     c.fact,
		},
		{
			Name:        "trash",
			Description: "Make a trash meme",
			Arguments:   "<user>",
			Handler:     c.trash,
		},
		{
			Name:        "stickbug",
			Description: "Make a stickbug meme",
			Arguments:   "<user>",
			Handler:     c.stickbug,
		},
		{
			Name:        "magik",
			Description: "Make a users avatar magik",
			Arguments:   "<user>",
			Handler:     c.magik,
		},
		{
			Name:        "randomgif",
			Description: "Send a random GIF from Giphy",
			Handler:     c.randomgif,
		},
		{
			Name:        "searchgif",
			Description: "Search GIFs on Giphy",
			Arguments:   "<query>",
			Handler:     c.searchgif,
		},
		{
			Name:        "communist",
			Description: "Make the user a communist",
			Arguments:   "<user>",
			Handler:     c.communist,
		},
		{
			Name:        "snow",
			Description: "Add a snow overlay the the users avatar",
			Arguments:   "<user>",
			Handler:     c.snow,
		},
		{
			Name:        "hub",
			Description: "Make a PornHub logo with custom text",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.hub,
		},
		{
			Name:        "floor",
			Description: "Make a floor meme",
			Arguments:   "<user> <text>",
			Handler:     c.floor,
		},
		{
			Name:        "jokeoverhead",
			Description: "When someone doesn't understand a joke",
			Arguments:   "<user>",
			Handler:     c.jokeoverhead,
		},
		{
			Name:        "throwaway",
			Description: "Throw a user away",
			Arguments:   "<user>",
			Handler:     c.throwaway,
		},
		{
			Name:        "salty",
			Description: "Make a user salty",
			Arguments:   "<user>",
			Handler:     c.salty,
		},
		{
			Name:        "scroll",
			Description: "Make a scroll meme",
			Arguments:   "<text>",
			Handler:     c.scroll,
		},
		{
			Name:        "shame",
			Description: "Shame a user",
			Arguments:   "<user>",
			Handler:     c.shame,
		},
		{
			Name:        "ship",
			Description: "Ship users",
			Arguments:   "<user1> <user2>",
			Handler:     c.ship,
		},
		{
			Name:        "what",
			Description: "Make a what meme",
			Arguments:   "<user>",
			Handler:     c.what,
		},
		{
			Name:        "sadcat",
			Description: "Send an image of a sad cat",
			Handler:     c.sadcat,
		},
		{
			Name:        "bird",
			Description: "Send an image of a bird",
			Handler:     c.bird,
		},
		{
			Name:        "cat",
			Description: "Send an image of a cat",
			Handler:     c.cat,
		},
		{
			Name:        "catgif",
			Description: "Send a GIF of a cat",
			Handler:     c.catgif,
		},
		{
			Name:        "catsays",
			Description: "Make a cat say something",
			Arguments:   "<text>",
			Handler:     c.catsays,
		},
		{
			Name:        "dog",
			Description: "Send an image of a dog",
			Handler:     c.dog,
		},
		{
			Name:        "panda",
			Description: "Send an image of a panda",
			Handler:     c.panda,
		},
		{
			Name:        "redpanda",
			Description: "Send an image of a red panda",
			Handler:     c.redpanda,
		},
		{
			Name:        "koala",
			Description: "Send an image of a koala",
			Handler:     c.koala,
		},
		{
			Name:        "raccoon",
			Description: "Send an image of a raccoon",
			Handler:     c.raccoon,
		},
		{
			Name:        "fox",
			Description: "Send an image of a fox",
			Handler:     c.fox,
		},
		{
			Name:        "shibe",
			Description: "Send an image of a shiba inu/shibe/doge",
			Handler:     c.shibe,
		},
		{
			Name:        "whale",
			Description: "Send an image of a whale",
			Handler:     c.whale,
		},
		{
			Name:        "kangaroo",
			Description: "Send an image of a kangaroo",
			Handler:     c.kangaroo,
		},
		{
			Name:        "bunny",
			Description: "Send an image of a bunny",
			Handler:     c.bunny,
		},
		{
			Name:        "lion",
			Description: "Send an image of a lion",
			Handler:     c.lion,
		},
		{
			Name:        "frog",
			Description: "Send an image of a frog",
			Handler:     c.frog,
		},
		{
			Name:        "duck",
			Description: "Send an image of a duck",
			Handler:     c.duck,
		},
		{
			Name:        "penguin",
			Description: "Send an image of a penguin",
			Handler:     c.penguin,
		},
		{
			Name:        "dolphin",
			Description: "Send an image of a dolphin",
			Handler:     c.dolphin,
		},
		{
			Name:        "bear",
			Description: "Send an image of a bear",
			Handler:     c.bear,
		},
		{
			Name:        "axolotl",
			Description: "Send an image of an axolotl",
			Handler:     c.axolotl,
		},
		{
			Name:        "capybara",
			Description: "Send an image of a capybara",
			Handler:     c.capybara,
		},
		{
			Name:        "balloon",
			Description: "Make a balloon meme",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.balloon,
		},
		{
			Name:        "abandon",
			Description: "Make an abandon meme",
			Arguments:   "<text>",
			Handler:     c.abandon,
		},
		{
			Name:        "aborted",
			Description: "Make an aborted meme",
			Arguments:   "<user>",
			Handler:     c.aborted,
		},
		{
			Name:        "affect",
			Description: "Make an affect meme",
			Arguments:   "<user>",
			Handler:     c.affect,
		},
		{
			Name:        "airpods",
			Description: "Add AirPods to the users avatar",
			Arguments:   "<user>",
			Handler:     c.airpods,
		},
		{
			Name:        "america",
			Description: "Patriotize the user",
			Arguments:   "<user>",
			Handler:     c.america,
		},
		{
			Name:        "armor",
			Description: "Make an armor meme",
			Arguments:   "<text>",
			Handler:     c.armor,
		},
		{
			Name:        "bed",
			Description: "Make a bed meme",
			Arguments:   "<user1> <user2>",
			Handler:     c.bed,
		},
		{
			Name:        "bongocat",
			Description: "Make a bongocat meme",
			Arguments:   "<user>",
			Handler:     c.bongocat,
		},
		{
			Name:        "boo",
			Description: "Make a boo meme",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.boo,
		},
		{
			Name:        "bigbrain",
			Description: "Make a big brain meme",
			Arguments:   "<\"item 1\"> <\"item 2\"> <\"item 3\"> <\"item 4\">",
			Handler:     c.bigbrain,
		},
		{
			Name:        "brazzers",
			Description: "Make a brazzers meme",
			Arguments:   "<user>",
			Handler:     c.brazzers,
		},
		{
			Name:        "byemom",
			Description: "Make a byemom meme",
			Arguments:   "<user> <text>",
			Handler:     c.byemom,
		},
		{
			Name:        "cancer",
			Description: "Make a Squidward cancer meme",
			Arguments:   "<user>",
			Handler:     c.cancer,
		},
		{
			Name:        "changemymind",
			Description: "Make a change my mind meme",
			Arguments:   "<text>",
			Handler:     c.changemymind,
		},
		{
			Name:        "cheating",
			Description: "Make a cheating meme",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.cheating,
		},
		{
			Name:        "citation",
			Description: "Make a citation meme",
			Arguments:   "<\"text1\"> <\"text2\"> <\"text3\">",
			Handler:     c.citation,
		},
		{
			Name:        "communism",
			Description: "Make the user a communist",
			Arguments:   "<user>",
			Handler:     c.communism,
		},
		{
			Name:        "confusedcat",
			Description: "Make a confused cat meme",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.confusedcat,
		},
		{
			Name:        "corporate",
			Description: "Make a corporate meme",
			Arguments:   "<user>",
			Handler:     c.corporate,
		},
		{
			Name:        "cry",
			Description: "Make a cry meme",
			Arguments:   "<text>",
			Handler:     c.cry,
		},
		{
			Name:        "dab",
			Description: "Make a dab meme",
			Arguments:   "<user>",
			Handler:     c.dab,
		},
		{
			Name:        "dank",
			Description: "Make a dank meme",
			Arguments:   "<user>",
			Handler:     c.dank,
		},
		{
			Name:        "delete",
			Description: "Make a delete meme",
			Arguments:   "<user>",
			Handler:     c.delete,
		},
		{
			Name:        "disability",
			Description: "Make a disability meme",
			Arguments:   "<user>",
			Handler:     c.disability,
		},
		{
			Name:        "doglemon",
			Description: "Make a dog lemon meme",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.doglemon,
		},
		{
			Name:        "door",
			Description: "Make a door meme",
			Arguments:   "<user>",
			Handler:     c.door,
		},
		{
			Name:        "egg",
			Description: "Make a egg meme",
			Arguments:   "<user>",
			Handler:     c.egg,
		},
		{
			Name:        "excuseme",
			Description: "Make an excuse me meme",
			Arguments:   "<text>",
			Handler:     c.excuseme,
		},
		{
			Name:        "expanddong",
			Description: "Make an expand dong meme",
			Arguments:   "<text>",
			Handler:     c.expanddong,
		},
		{
			Name:        "expandingwwe",
			Description: "Make an expanding WWE meme",
			Arguments:   "<\"item 1\"> <\"item 2\"> <\"item 3\"> <\"item 4\"> <\"item 5\">",
			Handler:     c.expandingwwe,
		},
		{
			Name:        "recaptcha",
			Description: "Make a recaptcha",
			Arguments:   "<text>",
			Handler:     c.recaptcha,
		},
		{
			Name:        "qrcode",
			Description: "Make a QR Code",
			Arguments:   "<text>",
			Handler:     c.qrcode,
		},
		{
			Name:        "darkqrcode",
			Description: "Make a dark QR Code",
			Arguments:   "<text>",
			Handler:     c.darkqrcode,
		},
		{
			Name:        "rule34",
			Description: "Random rule 34 image",
			Handler:     c.rule34,
		},
		{
			Name:        "achievement",
			Description: "Make a minecraft achievement",
			Arguments:   "<text>",
			Handler:     c.achievement,
		},
		{
			Name:        "amiajoke",
			Description: "Make an am I a joke meme",
			Arguments:   "<user>",
			Handler:     c.amiajoke,
		},
		{
			Name:        "bad",
			Description: "Make a bad meme",
			Arguments:   "<user>",
			Handler:     c.bad,
		},
		{
			Name:        "calling",
			Description: "Make a calling meme",
			Arguments:   "<text>",
			Handler:     c.calling,
		},
		{
			Name:        "didyoumean",
			Description: "Make a did you mean meme",
			Arguments:   "<\"top\"> <\"bottom\">",
			Handler:     c.didyoumean,
		},
		{
			Name:        "drake",
			Description: "Make a drake meme",
			Arguments:   "<\"top\"> <\"bottom\">",
			Handler:     c.drake,
		},
		{
			Name:        "ayano",
			Description: "Make an ayano meme",
			Arguments:   "<\"top\"> <\"bottom\">",
			Handler:     c.ayano,
		},
		{
			Name:        "facts",
			Description: "Make a facts meme",
			Arguments:   "<text>",
			Handler:     c.facts,
		},
		{
			Name:        "blur",
			Description: "Blur the user's avatar",
			Arguments:   "<user>",
			Handler:     c.blur,
		},
		{
			Name:        "flipav",
			Description: "Flip the user's avatar",
			Arguments:   "<user>",
			Handler:     c.flipav,
		},
		{
			Name:        "mirror",
			Description: "Mirror the user's avatar",
			Arguments:   "<user>",
			Handler:     c.mirror,
		},
		{
			Name:        "wide",
			Description: "Make the user's avatar wide",
			Arguments:   "<user>",
			Handler:     c.wide,
		},
		{
			Name:        "sepia",
			Description: "Make the user's avatar sepia",
			Arguments:   "<user>",
			Handler:     c.sepia,
		},
		{
			Name:        "adios",
			Description: "Make an adios meme",
			Arguments:   "<user>",
			Handler:     c.adios,
		},
		{
			Name:        "batmanslap",
			Description: "Make a batmanslap meme",
			Arguments:   "<user1> <user2> <\"text1\"> <\"text2\">",
			Handler:     c.batmanslap,
		},
		{
			Name:        "carreverse",
			Description: "Make a car reverse meme",
			Arguments:   "<text>",
			Handler:     c.carreverse,
		},
		{
			Name:        "drip",
			Description: "Give the user some drip",
			Arguments:   "<user>",
			Handler:     c.drip,
		},
		{
			Name:        "impostor",
			Description: "Make an Among Us impostor meme",
			Arguments:   "<text>",
			Handler:     c.impostor,
		},
		{
			Name:        "notimpostor",
			Description: "Make an Among Us not impostor meme",
			Arguments:   "<text>",
			Handler:     c.notimpostor,
		},
		{
			Name:        "emergencymeeting",
			Description: "Make an Among Us emergency meeting meme",
			Arguments:   "<text>",
			Handler:     c.emergencymeeting,
		},
		{
			Name:        "firsttime",
			Description: "Make a first time meme",
			Arguments:   "<user>",
			Handler:     c.firsttime,
		},
		{
			Name:        "grave",
			Description: "Make a grave meme",
			Arguments:   "<user>",
			Handler:     c.grave,
		},
		{
			Name:        "iamspeed",
			Description: "Make an i am speed meme",
			Arguments:   "<user>",
			Handler:     c.iamspeed,
		},
		{
			Name:        "icanmilkyou",
			Description: "Make an i can milk you meme",
			Arguments:   "<user1> <user2>",
			Handler:     c.icanmilkyou,
		},
		{
			Name:        "heaven",
			Description: "Make a heaven meme",
			Arguments:   "<user>",
			Handler:     c.heaven,
		},
		{
			Name:        "npc",
			Description: "Make an npc meme",
			Arguments:   "<\"text1\"> <\"text2\">",
			Handler:     c.npc,
		},
		{
			Name:        "stonks",
			Description: "Make a stonks meme",
			Arguments:   "<user>",
			Handler:     c.stonks,
		},
		{
			Name:        "notstonks",
			Description: "Make a not stonks meme",
			Arguments:   "<user>",
			Handler:     c.notstonks,
		},
		{
			Name:        "fliptable",
			Description: "Make the user flip a table",
			Arguments:   "<user>",
			Handler:     c.fliptable,
		},
		{
			Name:        "water",
			Description: "Make a water meme",
			Arguments:   "<text>",
			Handler:     c.water,
		},
		{
			Name:        "wolverine",
			Description: "Make a wolverine meme",
			Arguments:   "<user>",
			Handler:     c.wolverine,
		},
		{
			Name:        "womanyellingatcat",
			Description: "Make a woman yelling at cat meme",
			Arguments:   "<user1> <user2>",
			Handler:     c.womanyellingatcat,
		},
		{
			Name:        "pikachu",
			Description: "Send a random Pikachu image",
			Handler:     c.pikachu,
		},
		{
			Name:        "speechbubble",
			Description: "Add a speech bubble to an image",
			Arguments:   "[force GIF: on/off] [link/attachment]",
			Handler:     c.speechbubble,
		},

		{
			Name:        "quotereply",
			Description: "Quote the message you replied to",
			Arguments:   "[background url]",
			Handler:     c.quotereply,
		},
		{
			Name:        "quotemessage",
			Description: "Quote a message",
			Arguments:   "<message id> [channel id] [background url]",
			Handler:     c.quotemessage,
		},
		{
			Name:        "quoteuser",
			Description: "Quote a user",
			Arguments:   "<@user> <\"text\"> [background url]",
			Handler:     c.quoteuser,
		},
		{
			Name:        "quotetext",
			Description: "Quote using text",
			Arguments:   "<\"author\"> <\"text\"> [background url]",
			Handler:     c.quotetext,
		},

		{
			Name:        "discordquotereply",
			Description: "Quote the Discord message you replied to",
			Handler:     c.discordquotereply,
		},
		{
			Name:        "discordquotemessage",
			Description: "Quote a Discord message",
			Arguments:   "<message id> [channel id]",
			Handler:     c.discordquotemessage,
		},
		{
			Name:        "discordquoteuser",
			Description: "Quote a user",
			Arguments:   "<@user> <\"text\">",
			Handler:     c.discordquoteuser,
		},
		{
			Name:        "discordquotetext",
			Description: "Quote a Discord message",
			Arguments:   "<\"author\"> <\"text\"> <avatar url>",
			Handler:     c.discordquotetext,
		},
	})

}

func (c *ImageCommands) tweet(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf(
		"https://nekobot.xyz/api/imagegen?type=tweet&text=%s&username=%s",
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
	))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) meme(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/meme")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	meme := jsonResult.Get("meme")
	c.services.Embed.NewEmbed().
		SetTitle(meme.Get("title").MustString()).
		SetImage(meme.Get("imagelink").MustString()).
		AddField("URL", fmt.Sprintf("`%s`", meme.Get("permalink").MustString())).
		SetFooter(fmt.Sprintf("Subreddit: %s | Upvotes: %d | Posted by: u/%s", meme.Get("subreddit").MustString(), meme.Get("upvotes").MustInt(), meme.Get("author").MustString())).
		Send(ctx)
}

func (c *ImageCommands) ytcomment(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf(
		"https://some-random-api.com/canvas/youtube-comment?avatar=%s&comment=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
		url.QueryEscape(user.Username),
	))
}

func (c *ImageCommands) pokemon(ctx *router.Ctx) {
	ctx.RespondText("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%d.png", random.Int(1, 898))
}

func (c *ImageCommands) wasted(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/wasted?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) triggered(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/triggered?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) glass(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/glass?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) greyscale(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/greyscale?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) invert(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/invert?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) invertgreyscale(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/invertgreyscale?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) brightness(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/brightness?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) red(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/red?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) green(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/green?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) blue(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/blue?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) pixelate(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/pixelate?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) comrade(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/comrade?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) missionpassed(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/passed?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) jail(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/jail?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) nobitches(ctx *router.Ctx) {
	ctx.RespondText("https://some-random-api.com/canvas/nobitches?no=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) threshold(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/threshold?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) blurple(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/blurple?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) blurple2(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/blurple2?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) color(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/color?avatar=%s&color=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(strings.ReplaceAll(ctx.Arguments.Get(1).Raw(), "#", "")),
	))
}

func (c *ImageCommands) sostupid(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/its-so-stupid?avatar=%s&dog=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
}

func (c *ImageCommands) simpcard(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/simpcard?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) horny(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/horny?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) lgbt(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/lgbt?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) pansexual(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/pansexual?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) nonbinary(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/nonbinary?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) lesbian(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/lesbian?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) bisexual(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/bisexual?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) trans(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://some-random-api.com/canvas/transgender?avatar=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) oogway(ctx *router.Ctx) {
	ctx.RespondText("https://some-random-api.com/canvas/oogway?quote=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) oogway2(ctx *router.Ctx) {
	ctx.RespondText("https://some-random-api.com/canvas/oogway2?quote=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) threat(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=threats&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) clyde(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=clyde&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) distractedbf(ctx *router.Ctx) {
	user1, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=ship&user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) captcha(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=captcha&url=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) whowouldwin(ctx *router.Ctx) {
	user1, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=whowouldwin&user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) changemymind(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=changemymind&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) jpeg(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=jpeg&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) lolice(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=lolice&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) kannagen(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=kannagen&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) iphonex(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=iphonex&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) awooify(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=awooify&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) trap(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trap&image=%s&name=%s&author=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user.String()),
		url.QueryEscape(ctx.Message.Author.String()),
	))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) trumptweet(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trumptweet&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) deepfry(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=deepfry&image=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) blurpify(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=blurpify&image=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) phcomment(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=phcomment&text=%s&username=%s&image=%s",
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
		url.QueryEscape(user.Username),
		url.QueryEscape(user.AvatarURL("")),
	))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) fact(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=fact&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) trash(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trash&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) stickbug(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=stickbug&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) magik(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=magik&image=%s&intensity=5", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) randomgif(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://api.giphy.com/v1/gifs/random?api_key=ldQeNHnpL3WcCxJE1uO8HTk17ICn8i34&tag=&rating=R&limit=1&lang=en&offset=0")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.GetPath("data", "bitly_url").MustString())
}

func (c *ImageCommands) searchgif(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://api.giphy.com/v1/gifs/search?api_key=ldQeNHnpL3WcCxJE1uO8HTk17ICn8i34&q=%s&limit=1&offset=0&rating=R&lang=en", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	if len(jsonResult.Get("data").MustArray()) == 0 {
		ctx.RespondText("No results found.")
		return
	}
	ctx.RespondText(jsonResult.Get("data").GetIndex(0).Get("bitly_url").MustString())
}

func (c *ImageCommands) sadcat(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.alexflipnote.dev/sadcat", "file")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bird(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.alexflipnote.dev/birb", "file")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) cat(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat", "cat.jpg")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) catgif(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat/gif", "cat.gif")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) catsays(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat/says/"+url.PathEscape(ctx.Arguments.Raw()), "cat.jpg")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) dog(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://dog.ceo/api/breeds/image/random", "message")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) panda(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/panda", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) redpanda(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/red_panda", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) koala(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/koala", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) raccoon(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/raccoon", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) fox(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://randomfox.ca/floof", "image")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) shibe(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://shibe.online/api/shibes?count=1&urls=true&httpsUrls=true")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.GetIndex(0).MustString())
}

func (c *ImageCommands) whale(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/whale", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) kangaroo(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/kangaroo", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bunny(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/bunny", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) lion(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/lion", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) frog(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/frog", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) duck(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/duck", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) penguin(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/penguin", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) dolphin(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/dolphin", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bear(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/bear", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) axolotl(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/axolotl", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) capybara(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/capybara", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bed(ctx *router.Ctx) {
	user1, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/bed",
		"bed.png",
		map[string]string{
			"avatars": user1.AvatarURL("") + "," + user2.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) balloon(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/balloon",
		"balloon.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() + "," + ctx.Arguments.Get(1).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) abandon(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/abandon",
		"abandon.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) aborted(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/aborted",
		"aborted.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) affect(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/affect",
		"affect.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) airpods(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/airpods",
		"airpods.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) america(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/america",
		"america.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) armor(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/armor",
		"armor.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) bongocat(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/bongocat",
		"bongocat.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) boo(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/boo",
		"boo.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() + "," + ctx.Arguments.Get(1).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) bigbrain(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/brain",
		"bigbrain.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() +
				"," + ctx.Arguments.Get(1).Raw() +
				"," + ctx.Arguments.Get(2).Raw() +
				"," + ctx.Arguments.Get(3).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) brazzers(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/brazzers",
		"brazzers.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}

}

func (c *ImageCommands) byemom(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/byemom",
		"byemom.png",
		map[string]string{
			"avatars":   user.AvatarURL(""),
			"usernames": user.String(),
			"text":      ctx.Arguments.RawFrom(1),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) cancer(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/cancer",
		"cancer.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) cheating(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/cheating",
		"cheating.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() + "," + ctx.Arguments.Get(1).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) citation(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/citation",
		"citation.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() + "," + ctx.Arguments.Get(1).Raw() + "," + ctx.Arguments.Get(2).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) communism(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/communism",
		"communism.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) confusedcat(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/confusedcat",
		"confusedcat.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() + "," + ctx.Arguments.Get(1).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) corporate(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/corporate",
		"corporate.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) cry(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/cry",
		"cry.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) dab(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/dab",
		"dab.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) dank(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/dank",
		"dank.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) delete(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/delete",
		"delete.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) disability(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/disability",
		"disability.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) doglemon(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/doglemon",
		"doglemon.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() + "," + ctx.Arguments.Get(1).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) door(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/door",
		"door.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) egg(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	err = utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/egg",
		"egg.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) excuseme(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/excuseme",
		"excuseme.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) expanddong(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/expanddong",
		"expanddong.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) expandingwwe(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/expandingwwe",
		"expandingwwe.png",
		map[string]string{
			"text": ctx.Arguments.Get(0).Raw() +
				"," + ctx.Arguments.Get(1).Raw() +
				"," + ctx.Arguments.Get(2).Raw() +
				"," + ctx.Arguments.Get(3).Raw() +
				"," + ctx.Arguments.Get(4).Raw(),
		},
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) qrcode(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		fmt.Sprintf("https://www.qrcoder.co.uk/api/v1/?type=png&text=%s&int=1&background=FFFFFF&foreground=000000&eo=000000&ei=000000&size=8&padding=2", url.QueryEscape(ctx.Arguments.Raw())),
		"qrcode.png",
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) darkqrcode(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		fmt.Sprintf("https://www.qrcoder.co.uk/api/v1/?type=png&text=%s&int=1&background=000&foreground=FFFFFF&eo=FFFFFF&ei=FFFFFF&size=8&padding=2", url.QueryEscape(ctx.Arguments.Raw())),
		"qrcode.png",
	)
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) rule34(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&json=1")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.GetIndex(rand.Intn(len(jsonResult.MustArray()))).Get("file_url").MustString())
}

func (c *ImageCommands) achievement(ctx *router.Ctx) {
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/achievement?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) amiajoke(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/amiajoke?url=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) bad(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/bad?url=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) calling(ctx *router.Ctx) {
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/calling?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) didyoumean(ctx *router.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/didyoumean?top=%s&bottom=%s",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
	))
}

func (c *ImageCommands) drake(ctx *router.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/drake?top=%s&bottom=%s&ayano=false",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
	))
}

func (c *ImageCommands) ayano(ctx *router.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/drake?top=%s&bottom=%s&ayano=true",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
	))
}

func (c *ImageCommands) facts(ctx *router.Ctx) {
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/facts?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) blur(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=blur", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) flipav(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=flip", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) mirror(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=mirror", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) wide(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=wide", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) sepia(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=sepia", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) recaptcha(ctx *router.Ctx) {
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/captcha?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) gay(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=gay", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) communist(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=communist", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) snow(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=snow", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) hub(ctx *router.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/pornhub?text1=%s&text2=%s",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
	))
}

func (c *ImageCommands) floor(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/floor?url=%s&text=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
}

func (c *ImageCommands) jokeoverhead(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/jokeoverhead?url=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) throwaway(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/trash?url=%s&trash_url=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Message.Author.AvatarURL("")),
	))
}

func (c *ImageCommands) salty(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/salty?url=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) scroll(ctx *router.Ctx) {
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/scroll?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) shame(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/shame?url=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) ship(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/ship?url1=%s&url2=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user2.AvatarURL("")),
	))
}

func (c *ImageCommands) what(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://api.lightning-bot.com/api/v3/imagegen/what?url=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) adios(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/adios?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) batmanslap(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText(fmt.Sprintf(
		"https://vacefron.nl/api/batmanslap?text1=%s&text2=%s&batman=%s&robin=%s",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user2.AvatarURL("")),
	))
}

func (c *ImageCommands) carreverse(ctx *router.Ctx) {
	ctx.RespondText("https://vacefron.nl/api/carreverse?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) drip(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/drip?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) impostor(ctx *router.Ctx) {
	crewmate := random.Choice([]string{"black", "blue", "brown", "cyan", "darkgreen", "lime", "orange", "pink", "purple", "red", "white", "yellow"})
	ctx.RespondText("https://vacefron.nl/api/ejected?name=%s&impostor=true&crewmate=%s", ctx.Arguments.Raw(), crewmate)
}

func (c *ImageCommands) notimpostor(ctx *router.Ctx) {
	crewmate := random.Choice([]string{"black", "blue", "brown", "cyan", "darkgreen", "lime", "orange", "pink", "purple", "red", "white", "yellow"})
	ctx.RespondText("https://vacefron.nl/api/ejected?name=%s&impostor=false&crewmate=%s", ctx.Arguments.Raw(), crewmate)
}

func (c *ImageCommands) emergencymeeting(ctx *router.Ctx) {
	ctx.RespondText("https://vacefron.nl/api/emergencymeeting?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) firsttime(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/firsttime?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) grave(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/grave?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) iamspeed(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/iamspeed?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) icanmilkyou(ctx *router.Ctx) {
	user1, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/icanmilkyou?user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL("")))
}

func (c *ImageCommands) heaven(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/heaven?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) npc(ctx *router.Ctx) {
	ctx.RespondText("https://vacefron.nl/api/npc?text1=%s&text2=%s", url.QueryEscape(ctx.Arguments.Get(0).Raw()), url.QueryEscape(ctx.Arguments.Get(1).Raw()))
}

func (c *ImageCommands) stonks(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/stonks?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) notstonks(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/stonks?user=%s&notstonks=true", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) fliptable(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/tableflip?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) water(ctx *router.Ctx) {
	ctx.RespondText("https://vacefron.nl/api/water?text=%s", url.QueryEscape(ctx.Arguments.Raw()))
}

func (c *ImageCommands) wolverine(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/wolverine?user=%s", url.QueryEscape(user.AvatarURL("")))
}

func (c *ImageCommands) womanyellingatcat(ctx *router.Ctx) {
	user1, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	user2, err := ctx.Arguments.Get(1).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Cannot find user")
		return
	}
	ctx.RespondText("https://vacefron.nl/api/womanyellingatcat?woman=%s&cat=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL("")))
}

func (c *ImageCommands) pikachu(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/pikachu", "link")
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(imageUrl)
}

func (c *ImageCommands) speechbubble(ctx *router.Ctx) {
	link := ctx.Arguments.Get(1).Raw()

	attachments := ctx.Message.Attachments
	if len(attachments) > 0 {
		link = attachments[0].URL
	}

	if link == "" {
		ctx.RespondText("Please provide a link or attachment.")
		return
	}

	forceGIF, err := ctx.Arguments.Get(0).Bool()
	if err != nil {
		ctx.RespondText("Invalid force GIF value. Use `on` or `off`.")
		return
	}

	format := utils2.Ternary(forceGIF, "gif", "png")
	url := fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/speechbubble?url=%s&format=%s", url.QueryEscape(link), format)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	ctx.RespondText("https://api.lightning-bot.com" + httpResponse.Header.Get("Location"))
}

func (c *ImageCommands) quotereply(ctx *router.Ctx) {
	if ctx.Message.ReferencedMessage == nil {
		c.services.Embed.NewEmbed().
			SetDescription("You must reply to a message to use quotereply.").
			Send(ctx)
		return
	}

	m := ctx.Message.ReferencedMessage

	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(ctx.Message.ReferencedMessage.ContentWithMentionsReplaced()),
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
	))
}

func (c *ImageCommands) quotemessage(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find channel.").
			Send(ctx)
		return
	}

	m, err := utils.Discord.GetChannelMessage(ctx.Session, channel.ID, ctx.Arguments.Get(0).Raw())
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find message: " + err.Error()).
			Send(ctx)
		return
	}

	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(m.ContentWithMentionsReplaced()),
		url.QueryEscape(ctx.Arguments.Get(2).Raw()),
	))
}

func (c *ImageCommands) quoteuser(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find user.").
			Send(ctx)
		return
	}

	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Misc.Or(user.DisplayName, user.Username)),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
		url.QueryEscape(ctx.Arguments.Get(2).Raw()),
	))
}

func (c *ImageCommands) quotetext(ctx *router.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
		url.QueryEscape(ctx.Arguments.Get(2).Raw()),
	))
}

func (c *ImageCommands) getRoleColor(ctx *router.Ctx, userID string) string {
	if ctx.Message.GuildID == "" {
		return ""
	}

	member, ok := utils.Discord.GetMember(ctx.Session, ctx.Message.GuildID, userID)
	if !ok {
		return ""
	}

	var authorColor string

	for _, roleID := range member.Roles {
		role, err := ctx.Session.State.Role(ctx.Message.GuildID, roleID)
		if err != nil {
			continue
		}

		if role.Color != 0 {
			authorColor = fmt.Sprintf("#%x", role.Color)
		}
	}

	return authorColor
}

func (c *ImageCommands) discordquotemessage(ctx *router.Ctx) {
	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find channel.").
			Send(ctx)
		return
	}

	m, err := utils.Discord.GetChannelMessage(ctx.Session, channel.ID, ctx.Arguments.Get(0).Raw())
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find message: " + err.Error()).
			Send(ctx)
		return
	}

	authorColor := c.getRoleColor(ctx, m.Author.ID)

	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&author_color=%s&text=%s&timestamp=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(m.Author.AvatarURL("")),
		url.QueryEscape(authorColor),
		url.QueryEscape(m.ContentWithMentionsReplaced()),
		url.QueryEscape(m.Timestamp.Format("01/02/2006 3:04 PM")),
	))
}

func (c *ImageCommands) discordquotereply(ctx *router.Ctx) {
	if ctx.Message.ReferencedMessage == nil {
		c.services.Embed.NewEmbed().
			SetDescription("You must reply to a message to use quotereply.").
			Send(ctx)
		return
	}

	m := ctx.Message.ReferencedMessage

	authorColor := c.getRoleColor(ctx, m.Author.ID)

	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&author_color=%s&text=%s&timestamp=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(m.Author.AvatarURL("")),
		url.QueryEscape(authorColor),
		url.QueryEscape(m.ContentWithMentionsReplaced()),
		url.QueryEscape(m.Timestamp.Format("01/02/2006 3:04 PM")),
	))
}

func (c *ImageCommands) discordquoteuser(ctx *router.Ctx) {
	user, err := ctx.Arguments.Get(0).User(router.ArgumentFlagFallbackDefault)
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find user.").
			Send(ctx)
		return
	}

	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&text=%s",
		url.QueryEscape(utils.Misc.Or(user.DisplayName, user.Username)),
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
}

func (c *ImageCommands) discordquotetext(ctx *router.Ctx) {
	ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&text=%s",
		url.QueryEscape(ctx.Arguments.Get(0).Raw()),
		url.QueryEscape(ctx.Arguments.Get(2).Raw()),
		url.QueryEscape(ctx.Arguments.Get(1).Raw()),
	))
}
