package commands

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type ImageCommands struct{}

func (i ImageCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Image", "Image commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "tweet",
		Description: "Tweet a message",
		Usage:       "[p]tweet <username> <text>",
		Handler:     i.tweet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "meme",
		Description: "Send a random meme",
		Usage:       "[p]meme",
		Handler:     i.meme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ytcomment",
		Description: "Make a fake YouTube comment",
		Usage:       "[p]ytcomment <user> <text>",
		Handler:     i.ytcomment,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pokemon",
		Description: "Send a random Pokémon image",
		Usage:       "[p]pokemon",
		Handler:     i.pokemon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wasted",
		Description: "Add a wasted overlay on the users avatar",
		Usage:       "[p]wasted [user]",
		Handler:     i.wasted,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggered",
		Description: "Add a triggered overlay on the users avatar",
		Usage:       "[p]triggered [user]",
		Handler:     i.triggered,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "glass",
		Description: "Add a glass overlay on the users avatar",
		Usage:       "[p]glass [user]",
		Handler:     i.glass,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "greyscale",
		Description: "Add a greyscale overlay on the users avatar",
		Usage:       "[p]greyscale [user]",
		Handler:     i.greyscale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invert",
		Description: "Invert the users avatar",
		Usage:       "[p]invert [user]",
		Handler:     i.invert,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invertgreyscale",
		Description: "Invert and greyscale the users avatar",
		Usage:       "[p]invertgreyscale [user]",
		Handler:     i.invertgreyscale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "brightness",
		Description: "Make the users avatar brighter",
		Usage:       "[p]brightness [user]",
		Handler:     i.brightness,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "red",
		Description: "Make the users avatar red",
		Usage:       "[p]red [user]",
		Handler:     i.red,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "green",
		Description: "Make the users avatar green",
		Usage:       "[p]green [user]",
		Handler:     i.green,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blue",
		Description: "Make the users avatar blue",
		Usage:       "[p]blue [user]",
		Handler:     i.blue,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pixelate",
		Description: "Pixelate the users avatar",
		Usage:       "[p]pixelate [user]",
		Handler:     i.pixelate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "comrade",
		Description: "Add comrade overlay on users avatar",
		Usage:       "[p]comrade [user]",
		Handler:     i.comrade,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "missionpassed",
		Description: "Add Mission Passed overlay on users avatar",
		Usage:       "[p]missionpassed [user]",
		Handler:     i.missionpassed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jail",
		Description: "Add jail overlay on users avatar",
		Usage:       "[p]jail [user]",
		Handler:     i.jail,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nobitches",
		Description: "Make a no bitches meme",
		Usage:       "[p]nobitches <text>",
		Handler:     i.nobitches,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "threat",
		Description: "Make biggest threats meme",
		Usage:       "[p]threat [user]",
		Handler:     i.threat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "threshold",
		Description: "Add a threshold filter to a users avatar",
		Usage:       "[p]threshold <user>",
		Handler:     i.threshold,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blurple",
		Description: "Make the users avatar blurple",
		Usage:       "[p]blurple",
		Handler:     i.blurple,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blurple2",
		Description: "Make the users avatar new blurple",
		Usage:       "[p]blurple2",
		Handler:     i.blurple2,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "color",
		Description: "Make the users avatar a custom color",
		Usage:       "[p]color <user> <color>",
		Handler:     i.color,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sostupid",
		Description: "Make an it's so stupid meme",
		Usage:       "[p]sostupid <user> <text>",
		Handler:     i.sostupid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "simpcard",
		Description: "Make a simp card for a user",
		Usage:       "[p]simpcard <user>",
		Handler:     i.simpcard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "horny",
		Description: "Make a horny card for a user",
		Usage:       "[p]horny <user>",
		Handler:     i.horny,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lgbt",
		Description: "Add an lgbt border to the user",
		Usage:       "[p]lgbt <user>",
		Handler:     i.lgbt,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pansexual",
		Description: "Add a pansexual border to the user",
		Usage:       "[p]pansexual <user>",
		Handler:     i.pansexual,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nonbinary",
		Description: "Add a nonbinary border to the user",
		Usage:       "[p]nonbinary <user>",
		Handler:     i.nonbinary,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lesbian2",
		Description: "Add a lesbian border to the user",
		Usage:       "[p]lesbian2 <user>",
		Handler:     i.lesbian,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bixesual",
		Description: "Add a bixesual border to the user",
		Usage:       "[p]bixesual <user>",
		Handler:     i.bixesual,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trans",
		Description: "Add a trans border to the user",
		Usage:       "[p]trans <user>",
		Handler:     i.trans,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "oogway",
		Description: "Make a master oogway meme",
		Usage:       "[p]oogway <text>",
		Handler:     i.oogway,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "oogway2",
		Description: "Make a master oogway meme",
		Usage:       "[p]oogway2 <text>",
		Handler:     i.oogway2,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clyde",
		Description: "Make Clyde say something",
		Usage:       "[p]clyde <text>",
		Handler:     i.clyde,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gay",
		Description: "Add a rainbow overlay on the users avatar",
		Usage:       "[p]gay <user>",
		Handler:     i.gay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "distractedbf",
		Description: "Make a distracted boyfriend meme with the users",
		Usage:       "[p]distractedbf <user1> <user2>",
		Handler:     i.distractedbf,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "captcha",
		Description: "Make a captcha meme",
		Usage:       "[p]captcha <user> <text>",
		Handler:     i.captcha,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "whowouldwin",
		Description: "Make a who would win meme",
		Usage:       "[p]whowouldwin <user1> <user2>",
		Handler:     i.whowouldwin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jpeg",
		Description: "JPEG a users avatar",
		Usage:       "[p]jpeg <user>",
		Handler:     i.jpeg,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lolice",
		Description: "Make a lolice meme",
		Usage:       "[p]lolice <user>",
		Handler:     i.lolice,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kannagen",
		Description: "Make a kannagen meme",
		Usage:       "[p]kannagen <user1> <user2>",
		Handler:     i.kannagen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iphonex",
		Description: "Put users profile picture on an iPhone X",
		Usage:       "[p]iphonex <user1> <user2>",
		Handler:     i.iphonex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "awooify",
		Description: "Awooify the user",
		Usage:       "[p]awooify <user1> <user2>",
		Handler:     i.awooify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trap",
		Description: "Make a trap card",
		Usage:       "[p]trap <user>",
		Handler:     i.trap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trumptweet",
		Description: "Make a fake Trump tweet",
		Usage:       "[p]trumptweet <text>",
		Handler:     i.trumptweet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deepfry",
		Description: "Deepfry the users avatar",
		Usage:       "[p]deepfry <text>",
		Handler:     i.deepfry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blurpify",
		Description: "Blurpify the users avatar",
		Usage:       "[p]blurpify <text>",
		Handler:     i.blurpify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "phcomment",
		Description: "Make the user comment something on PornHub",
		Usage:       "[p]phcomment <user> <text>",
		Handler:     i.phcomment,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fact",
		Description: "Make a fact meme",
		Usage:       "[p]fact <user>",
		Handler:     i.fact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trash",
		Description: "Make a trash meme",
		Usage:       "[p]trash <user>",
		Handler:     i.trash,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stickbug",
		Description: "Make a stickbug meme",
		Usage:       "[p]stickbug <user>",
		Handler:     i.stickbug,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "magik",
		Description: "Make a users avatar magik",
		Usage:       "[p]magik <user>",
		Handler:     i.magik,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomgif",
		Description: "Send a random GIF from Giphy",
		Usage:       "[p]randomgif",
		Handler:     i.randomgif,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "searchgif",
		Description: "Search GIFs on Giphy",
		Usage:       "[p]searchgif <query>",
		Handler:     i.searchgif,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "communist",
		Description: "Make the user a communist",
		Usage:       "[p]communist <user>",
		Handler:     i.communist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snow",
		Description: "Add a snow overlay the the users avatar",
		Usage:       "[p]snow <user>",
		Handler:     i.snow,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hub",
		Description: "Make a PornHub logo with custom text",
		Usage:       "[p]hub <\"text1\"> <\"text2\">",
		Handler:     i.hub,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "floor",
		Description: "Make a floor meme",
		Usage:       "[p]floor <user> <text>",
		Handler:     i.floor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jokeoverhead",
		Description: "When someone doesn't understand a joke",
		Usage:       "[p]jokeoverhead <user>",
		Handler:     i.jokeoverhead,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "throwaway",
		Description: "Throw a user away",
		Usage:       "[p]throwaway <user>",
		Handler:     i.throwaway,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "salty",
		Description: "Make a user salty",
		Usage:       "[p]salty <user>",
		Handler:     i.salty,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "scroll",
		Description: "Make a scroll meme",
		Usage:       "[p]scroll <text>",
		Handler:     i.scroll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shame",
		Description: "Shame a user",
		Usage:       "[p]shame <user>",
		Handler:     i.shame,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ship",
		Description: "Ship users",
		Usage:       "[p]ship <user1> <user2>",
		Handler:     i.ship,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "what",
		Description: "Make a what meme",
		Usage:       "[p]what <user>",
		Handler:     i.what,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sadcat",
		Description: "Send an image of a sad cat",
		Usage:       "[p]sadcat",
		Handler:     i.sadcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bird",
		Description: "Send an image of a bird",
		Usage:       "[p]bird",
		Handler:     i.bird,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cat",
		Description: "Send an image of a cat",
		Usage:       "[p]cat",
		Handler:     i.cat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catgif",
		Description: "Send a GIF of a cat",
		Usage:       "[p]catgif",
		Handler:     i.catgif,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catsays",
		Description: "Make a cat say something",
		Usage:       "[p]catsays <text>",
		Handler:     i.catsays,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dog",
		Description: "Send an image of a dog",
		Usage:       "[p]dog",
		Handler:     i.dog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "panda",
		Description: "Send an image of a panda",
		Usage:       "[p]panda",
		Handler:     i.panda,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "redpanda",
		Description: "Send an image of a red panda",
		Usage:       "[p]redpanda",
		Handler:     i.redpanda,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "koala",
		Description: "Send an image of a koala",
		Usage:       "[p]koala",
		Handler:     i.koala,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "raccoon",
		Description: "Send an image of a raccoon",
		Usage:       "[p]raccoon",
		Handler:     i.raccoon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fox",
		Description: "Send an image of a fox",
		Usage:       "[p]fox",
		Handler:     i.fox,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shibe",
		Description: "Send an image of a shiba inu/shibe/doge",
		Usage:       "[p]shibe",
		Handler:     i.shibe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "whale",
		Description: "Send an image of a whale",
		Usage:       "[p]whale",
		Handler:     i.whale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kangaroo",
		Description: "Send an image of a kangaroo",
		Usage:       "[p]kangaroo",
		Handler:     i.kangaroo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bunny",
		Description: "Send an image of a bunny",
		Usage:       "[p]bunny",
		Handler:     i.bunny,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lion",
		Description: "Send an image of a lion",
		Usage:       "[p]lion",
		Handler:     i.lion,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "frog",
		Description: "Send an image of a frog",
		Usage:       "[p]frog",
		Handler:     i.frog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "duck",
		Description: "Send an image of a duck",
		Usage:       "[p]duck",
		Handler:     i.duck,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "penguin",
		Description: "Send an image of a penguin",
		Usage:       "[p]penguin",
		Handler:     i.penguin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dolphin",
		Description: "Send an image of a dolphin",
		Usage:       "[p]dolphin",
		Handler:     i.dolphin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bear",
		Description: "Send an image of a bear",
		Usage:       "[p]bear",
		Handler:     i.bear,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "axolotl",
		Description: "Send an image of an axolotl",
		Usage:       "[p]axolotl",
		Handler:     i.axolotl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "capybara",
		Description: "Send an image of a capybara",
		Usage:       "[p]capybara",
		Handler:     i.capybara,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "balloon",
		Description: "Make a balloon meme",
		Usage:       "[p]balloon <\"text1\"> <\"text2\">",
		Handler:     i.balloon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "abandon",
		Description: "Make an abandon meme",
		Usage:       "[p]abandon <text>",
		Handler:     i.abandon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aborted",
		Description: "Make an aborted meme",
		Usage:       "[p]aborted <user>",
		Handler:     i.aborted,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "affect",
		Description: "Make an affect meme",
		Usage:       "[p]affect <user>",
		Handler:     i.affect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "airpods",
		Description: "Add AirPods to the users avatar",
		Usage:       "[p]airpods <user>",
		Handler:     i.airpods,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "america",
		Description: "Patriotize the user",
		Usage:       "[p]america <user>",
		Handler:     i.america,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "armor",
		Description: "Make an armor meme",
		Usage:       "[p]armor <text>",
		Handler:     i.armor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bed",
		Description: "Make a bed meme",
		Usage:       "[p]bed <user1> <user2>",
		Handler:     i.bed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bongocat",
		Description: "Make a bongocat meme",
		Usage:       "[p]bongocat <user>",
		Handler:     i.bongocat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "boo",
		Description: "Make a boo meme",
		Usage:       "[p]boo <\"text1\"> <\"text2\">",
		Handler:     i.boo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bigbrain",
		Description: "Make a big brain meme",
		Usage:       "[p]bigbrain <\"item 1\"> <\"item 2\"> <\"item 3\"> <\"item 4\">",
		Handler:     i.bigbrain,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "brazzers",
		Description: "Make a brazzers meme",
		Usage:       "[p]brazzers <user>",
		Handler:     i.brazzers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "byemom",
		Description: "Make a byemom meme",
		Usage:       "[p]byemom <user> <text>",
		Handler:     i.byemom,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cancer",
		Description: "Make a Squidward cancer meme",
		Usage:       "[p]cancer <user>",
		Handler:     i.cancer,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "changemymind",
		Description: "Make a change my mind meme",
		Usage:       "[p]changemymind <text>",
		Handler:     i.changemymind,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cheating",
		Description: "Make a cheating meme",
		Usage:       "[p]cheating <\"text1\"> <\"text2\">",
		Handler:     i.cheating,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "citation",
		Description: "Make a citation meme",
		Usage:       "[p]citation <\"text1\"> <\"text2\"> <\"text3\">",
		Handler:     i.citation,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "communism",
		Description: "Make the user a communist",
		Usage:       "[p]communism <user>",
		Handler:     i.communism,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "confusedcat",
		Description: "Make a confused cat meme",
		Usage:       "[p]confusedcat <\"text1\"> <\"text2\">",
		Handler:     i.confusedcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "corporate",
		Description: "Make a corporate meme",
		Usage:       "[p]corporate <user>",
		Handler:     i.corporate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cry",
		Description: "Make a cry meme",
		Usage:       "[p]cry <text>",
		Handler:     i.cry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dab",
		Description: "Make a dab meme",
		Usage:       "[p]dab <user>",
		Handler:     i.dab,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dank",
		Description: "Make a dank meme",
		Usage:       "[p]dank <user>",
		Handler:     i.dank,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delete",
		Description: "Make a delete meme",
		Usage:       "[p]delete <user>",
		Handler:     i.delete,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "disability",
		Description: "Make a disability meme",
		Usage:       "[p]disability <user>",
		Handler:     i.disability,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "doglemon",
		Description: "Make a dog lemon meme",
		Usage:       "[p]doglemon <\"text1\"> <\"text2\">",
		Handler:     i.doglemon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "door",
		Description: "Make a door meme",
		Usage:       "[p]door <user>",
		Handler:     i.door,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "egg",
		Description: "Make a egg meme",
		Usage:       "[p]egg <user>",
		Handler:     i.egg,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "excuseme",
		Description: "Make an excuse me meme",
		Usage:       "[p]excuseme <text>",
		Handler:     i.excuseme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "expanddong",
		Description: "Make an expand dong meme",
		Usage:       "[p]expanddong <text>",
		Handler:     i.expanddong,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "expandingwwe",
		Description: "Make an expanding WWE meme",
		Usage:       "[p]expandingwwe <\"item 1\"> <\"item 2\"> <\"item 3\"> <\"item 4\"> <\"item 5\">",
		Handler:     i.expandingwwe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "recaptcha",
		Description: "Make a recaptcha",
		Usage:       "[p]recaptcha <text>",
		Handler:     i.recaptcha,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "qrcode",
		Description: "Make a QR Code",
		Usage:       "[p]qrcode <text>",
		Handler:     i.qrcode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "darkqrcode",
		Description: "Make a dark QR Code",
		Usage:       "[p]darkqrcode <text>",
		Handler:     i.darkqrcode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rule34",
		Description: "Random rule 34 image",
		Usage:       "[p]rule34",
		Handler:     i.rule34,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "achievement",
		Description: "Make a minecraft achievement",
		Usage:       "[p]achievement <text>",
		Handler:     i.achievement,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amiajoke",
		Description: "Make an am I a joke meme",
		Usage:       "[p]amiajoke <user>",
		Handler:     i.amiajoke,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bad",
		Description: "Make a bad meme",
		Usage:       "[p]bad <user>",
		Handler:     i.bad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calling",
		Description: "Make a calling meme",
		Usage:       "[p]calling <text>",
		Handler:     i.calling,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "didyoumean",
		Description: "Make a did you mean meme",
		Usage:       "[p]didyoumean <\"top\"> <\"bottom\">",
		Handler:     i.didyoumean,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "drake",
		Description: "Make a drake meme",
		Usage:       "[p]drake <\"top\"> <\"bottom\">",
		Handler:     i.drake,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ayano",
		Description: "Make an ayano meme",
		Usage:       "[p]ayano <\"top\"> <\"bottom\">",
		Handler:     i.ayano,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "facts",
		Description: "Make a facts meme",
		Usage:       "[p]facts <text>",
		Handler:     i.facts,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blur",
		Description: "Blur the user's avatar",
		Usage:       "[p]blur <user>",
		Handler:     i.blur,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "flipav",
		Description: "Flip the user's avatar",
		Usage:       "[p]flipav <user>",
		Handler:     i.flipav,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mirror",
		Description: "Mirror the user's avatar",
		Usage:       "[p]mirror <user>",
		Handler:     i.mirror,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wide",
		Description: "Make the user's avatar wide",
		Usage:       "[p]wide <user>",
		Handler:     i.wide,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sepia",
		Description: "Make the user's avatar sepia",
		Usage:       "[p]sepia <user>",
		Handler:     i.sepia,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adios",
		Description: "Make an adios meme",
		Usage:       "[p]adios <user>",
		Handler:     i.adios,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "batmanslap",
		Description: "Make a batmanslap meme",
		Usage:       "[p]batmanslap <user1> <user2> <\"text1\"> <\"text2\">",
		Handler:     i.batmanslap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "carreverse",
		Description: "Make a car reverse meme",
		Usage:       "[p]carreverse <text>",
		Handler:     i.carreverse,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "drip",
		Description: "Give the user some drip",
		Usage:       "[p]drip <user>",
		Handler:     i.drip,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "impostor",
		Description: "Make an Among Us impostor meme",
		Usage:       "[p]impostor <text>",
		Handler:     i.impostor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notimpostor",
		Description: "Make an Among Us not impostor meme",
		Usage:       "[p]notimpostor <text>",
		Handler:     i.notimpostor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "emergencymeeting",
		Description: "Make an Among Us emergency meeting meme",
		Usage:       "[p]emergencymeeting <text>",
		Handler:     i.emergencymeeting,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "firsttime",
		Description: "Make a first time meme",
		Usage:       "[p]firsttime <user>",
		Handler:     i.firsttime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "grave",
		Description: "Make a grave meme",
		Usage:       "[p]grave <user>",
		Handler:     i.grave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iamspeed",
		Description: "Make an i am speed meme",
		Usage:       "[p]iamspeed <user>",
		Handler:     i.iamspeed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "icanmilkyou",
		Description: "Make an i can milk you meme",
		Usage:       "[p]icanmilkyou <user1> <user2>",
		Handler:     i.icanmilkyou,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "heaven",
		Description: "Make a heaven meme",
		Usage:       "[p]heaven <user>",
		Handler:     i.heaven,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "npc",
		Description: "Make an npc meme",
		Usage:       "[p]npc <\"text1\"> <\"text2\">",
		Handler:     i.npc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stonks",
		Description: "Make a stonks meme",
		Usage:       "[p]stonks <user>",
		Handler:     i.stonks,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notstonks",
		Description: "Make a not stonks meme",
		Usage:       "[p]notstonks <user>",
		Handler:     i.notstonks,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fliptable",
		Description: "Make the user flip a table",
		Usage:       "[p]fliptable <user>",
		Handler:     i.fliptable,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "water",
		Description: "Make a water meme",
		Usage:       "[p]water <text>",
		Handler:     i.water,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wolverine",
		Description: "Make a wolverine meme",
		Usage:       "[p]wolverine <user>",
		Handler:     i.wolverine,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "womanyellingatcat",
		Description: "Make a woman yelling at cat meme",
		Usage:       "[p]womanyellingatcat <user1> <user2>",
		Handler:     i.womanyellingatcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pikachu",
		Description: "Send a random Pikachu image",
		Usage:       "[p]pikachu",
		Handler:     i.pikachu,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "speechbubble",
		Description: "Add a speech bubble to an image",
		Usage:       "[p]speechbubble [link/attachment] [force GIF: on/off]",
		Handler:     i.speechbubble,
	})
}

func (ImageCommands) tweet(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf(
		"https://nekobot.xyz/api/imagegen?type=tweet&text=%s&username=%s",
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
		url.QueryEscape(utils.Arguments["<username>"]),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) meme(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/meme")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	meme := jsonResult.Get("meme")
	embed.NewEmbed().
		SetTitle(meme.Get("title").MustString()).
		SetImage(meme.Get("imagelink").MustString()).
		AddField("URL", fmt.Sprintf("`%s`", meme.Get("permalink").MustString())).
		SetFooter(fmt.Sprintf("Subreddit: %s | Upvotes: %d | Posted by: u/%s", meme.Get("subreddit").MustString(), meme.Get("upvotes").MustInt(), meme.Get("author").MustString())).
		Send(ctx)
}

func (ImageCommands) ytcomment(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf(
		"https://some-random-api.com/canvas/youtube-comment?avatar=%s&comment=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
		url.QueryEscape(user.Username),
	))
}

func (ImageCommands) pokemon(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%d.png", utils.Random.Number(1, 898)))
}

func (ImageCommands) wasted(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/wasted?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) triggered(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/triggered?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) glass(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/glass?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) greyscale(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/greyscale?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) invert(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/invert?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) invertgreyscale(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/invertgreyscale?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) brightness(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/brightness?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) red(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/red?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) green(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/green?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) blue(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/blue?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) pixelate(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/pixelate?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) comrade(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/comrade?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) missionpassed(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/passed?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) jail(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/jail?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) nobitches(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/nobitches?no=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) threshold(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/threshold?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) blurple(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/blurple?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) blurple2(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/blurple2?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) color(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/color?avatar=%s&color=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(strings.ReplaceAll(utils.Arguments["<color>"], "#", "")),
	))
}

func (ImageCommands) sostupid(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/its-so-stupid?avatar=%s&dog=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
	))
}

func (ImageCommands) simpcard(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/simpcard?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) horny(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/horny?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) lgbt(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/lgbt?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) pansexual(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/pansexual?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) nonbinary(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/nonbinary?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) lesbian(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/lesbian?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) bixesual(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/bisexual?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) trans(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/transgender?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) oogway(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/oogway?quote=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) oogway2(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/oogway2?quote=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) threat(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=threats&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) clyde(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=clyde&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) distractedbf(ctx *dgc.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=ship&user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) captcha(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=captcha&url=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) whowouldwin(ctx *dgc.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=whowouldwin&user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) changemymind(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=changemymind&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) jpeg(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=jpeg&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) lolice(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=lolice&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) kannagen(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=kannagen&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) iphonex(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=iphonex&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) awooify(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=awooify&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) trap(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trap&image=%s&name=%s&author=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user.String()),
		url.QueryEscape(ctx.Event.Author.String()),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) trumptweet(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trumptweet&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) deepfry(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=deepfry&image=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) blurpify(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=blurpify&image=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) phcomment(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=phcomment&text=%s&username=%s&image=%s",
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
		url.QueryEscape(user.Username),
		url.QueryEscape(user.AvatarURL("")),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) fact(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=fact&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) trash(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trash&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) stickbug(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=stickbug&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) magik(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=magik&image=%s&intensity=5", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (ImageCommands) randomgif(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson("https://api.giphy.com/v1/gifs/random?api_key=ldQeNHnpL3WcCxJE1uO8HTk17ICn8i34&tag=&rating=R&limit=1&lang=en&offset=0")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.GetPath("data", "bitly_url").MustString())
}

func (ImageCommands) searchgif(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://api.giphy.com/v1/gifs/search?api_key=ldQeNHnpL3WcCxJE1uO8HTk17ICn8i34&q=%s&limit=1&offset=0&rating=R&lang=en", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	if len(jsonResult.Get("data").MustArray()) == 0 {
		_ = ctx.RespondText("No results found.")
		return
	}
	_ = ctx.RespondText(jsonResult.Get("data").GetIndex(0).Get("bitly_url").MustString())
}

func (ImageCommands) sadcat(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.alexflipnote.dev/sadcat", "file")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) bird(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.alexflipnote.dev/birb", "file")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) cat(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat", "cat.jpg")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) catgif(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat/gif", "cat.gif")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) catsays(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat/says/"+url.PathEscape(ctx.Arguments.Raw()), "cat.jpg")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) dog(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://dog.ceo/api/breeds/image/random", "message")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) panda(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/panda", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) redpanda(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/red_panda", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) koala(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/koala", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) raccoon(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/raccoon", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) fox(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://randomfox.ca/floof", "image")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) shibe(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson("https://shibe.online/api/shibes?count=1&urls=true&httpsUrls=true")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.GetIndex(0).MustString())
}

func (ImageCommands) whale(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/whale", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) kangaroo(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/kangaroo", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) bunny(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/bunny", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) lion(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/lion", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) frog(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/frog", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) duck(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/duck", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) penguin(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/penguin", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) dolphin(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/dolphin", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) bear(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/bear", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) axolotl(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/axolotl", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) capybara(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/capybara", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) bed(ctx *dgc.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/bed",
		"bed.png",
		map[string]string{
			"avatars": user1.AvatarURL("") + "," + user2.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) balloon(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/balloon",
		"balloon.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) abandon(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/abandon",
		"abandon.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) aborted(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/aborted",
		"aborted.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) affect(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/affect",
		"affect.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) airpods(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/airpods",
		"airpods.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) america(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/america",
		"america.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) armor(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/armor",
		"armor.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) bongocat(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/bongocat",
		"bongocat.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) boo(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/boo",
		"boo.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) bigbrain(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/brain",
		"bigbrain.png",
		map[string]string{
			"text": utils.Arguments["<\"item 1\">"] +
				"," + utils.Arguments["<\"item 2\">"] +
				"," + utils.Arguments["<\"item 3\">"] +
				"," + utils.Arguments["<\"item 4\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) brazzers(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/brazzers",
		"brazzers.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}

}

func (ImageCommands) byemom(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/byemom",
		"byemom.png",
		map[string]string{
			"avatars":   user.AvatarURL(""),
			"usernames": user.String(),
			"text":      utils.String.RemoveWords(ctx.Arguments.Raw(), 1),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) cancer(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/cancer",
		"cancer.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) cheating(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/cheating",
		"cheating.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) citation(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/citation",
		"citation.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"] + "," + utils.Arguments["<\"text3\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) communism(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/communism",
		"communism.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) confusedcat(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/confusedcat",
		"confusedcat.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) corporate(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/corporate",
		"corporate.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) cry(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/cry",
		"cry.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) dab(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/dab",
		"dab.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) dank(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/dank",
		"dank.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) delete(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/delete",
		"delete.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) disability(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/disability",
		"disability.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) doglemon(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/doglemon",
		"doglemon.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) door(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/door",
		"door.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) egg(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/egg",
		"egg.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) excuseme(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/excuseme",
		"excuseme.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) expanddong(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/expanddong",
		"expanddong.png",
		map[string]string{
			"text": ctx.Arguments.Raw(),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) expandingwwe(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.herokuapp.com/expandingwwe",
		"expandingwwe.png",
		map[string]string{
			"text": utils.Arguments["<\"item 1\">"] +
				"," + utils.Arguments["<\"item 2\">"] +
				"," + utils.Arguments["<\"item 3\">"] +
				"," + utils.Arguments["<\"item 4\">"] +
				"," + utils.Arguments["<\"item 5\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) qrcode(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		fmt.Sprintf("https://www.qrcoder.co.uk/api/v1/?type=png&text=%s&int=1&background=FFFFFF&foreground=000000&eo=000000&ei=000000&size=8&padding=2", url.QueryEscape(ctx.Arguments.Raw())),
		"qrcode.png",
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) darkqrcode(ctx *dgc.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		fmt.Sprintf("https://www.qrcoder.co.uk/api/v1/?type=png&text=%s&int=1&background=000&foreground=FFFFFF&eo=FFFFFF&ei=FFFFFF&size=8&padding=2", url.QueryEscape(ctx.Arguments.Raw())),
		"qrcode.png",
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (ImageCommands) rule34(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson("https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&json=1")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.GetIndex(rand.Intn(len(jsonResult.MustArray()))).Get("file_url").MustString())
}

func (ImageCommands) achievement(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/achievement?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) amiajoke(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/amiajoke?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) bad(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/bad?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) calling(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/calling?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) didyoumean(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/didyoumean?top=%s&bottom=%s",
		url.QueryEscape(utils.Arguments["<\"top\">"]),
		url.QueryEscape(utils.Arguments["<\"bottom\">"]),
	))
}

func (ImageCommands) drake(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/drake?top=%s&bottom=%s&ayano=false",
		url.QueryEscape(utils.Arguments["<\"top\">"]),
		url.QueryEscape(utils.Arguments["<\"bottom\">"]),
	))
}

func (ImageCommands) ayano(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/drake?top=%s&bottom=%s&ayano=true",
		url.QueryEscape(utils.Arguments["<\"top\">"]),
		url.QueryEscape(utils.Arguments["<\"bottom\">"]),
	))
}

func (ImageCommands) facts(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/facts?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) blur(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=blur", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) flipav(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=flip", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) mirror(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=mirror", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) wide(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=wide", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) sepia(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=sepia", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) recaptcha(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/captcha?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) gay(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=gay", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) communist(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=communist", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) snow(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=snow", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) hub(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/pornhub?text1=%s&text2=%s",
		url.QueryEscape(utils.Arguments["<\"text1\">"]),
		url.QueryEscape(utils.Arguments["<\"text2\">"]),
	))
}

func (ImageCommands) floor(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/floor?url=%s&text=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(utils.String.RemoveWords(ctx.Arguments.Raw(), 1)),
	))
}

func (ImageCommands) jokeoverhead(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/jokeoverhead?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) throwaway(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/trash?url=%s&trash_url=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Event.Author.AvatarURL("")),
	))
}

func (ImageCommands) salty(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/salty?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) scroll(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/scroll?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) shame(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/shame?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) ship(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/ship?url1=%s&url2=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user2.AvatarURL("")),
	))
}

func (ImageCommands) what(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/what?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) adios(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/adios?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) batmanslap(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf(
		"https://vacefron.nl/api/batmanslap?text1=%s&text2=%s&batman=%s&robin=%s",
		url.QueryEscape(utils.Arguments["<\"text1\">"]),
		url.QueryEscape(utils.Arguments["<\"text2\">"]),
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user2.AvatarURL("")),
	))
}

func (ImageCommands) carreverse(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/carreverse?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) drip(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/drip?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) impostor(ctx *dgc.Ctx) {
	crewmate := utils.RandomChoice([]string{"black", "blue", "brown", "cyan", "darkgreen", "lime", "orange", "pink", "purple", "red", "white", "yellow"})
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/ejected?name=%s&impostor=true&crewmate=%s", ctx.Arguments.Raw(), crewmate))
}

func (ImageCommands) notimpostor(ctx *dgc.Ctx) {
	crewmate := utils.RandomChoice([]string{"black", "blue", "brown", "cyan", "darkgreen", "lime", "orange", "pink", "purple", "red", "white", "yellow"})
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/ejected?name=%s&impostor=false&crewmate=%s", ctx.Arguments.Raw(), crewmate))
}

func (ImageCommands) emergencymeeting(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/emergencymeeting?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) firsttime(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/firsttime?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) grave(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/grave?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) iamspeed(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/iamspeed?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) icanmilkyou(ctx *dgc.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/icanmilkyou?user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
}

func (ImageCommands) heaven(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/heaven?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) npc(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/npc?text1=%s&text2=%s", url.QueryEscape(utils.Arguments["<\"text1\">"]), url.QueryEscape(utils.Arguments["<\"text2\">"])))
}

func (ImageCommands) stonks(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/stonks?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) notstonks(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/stonks?user=%s&notstonks=true", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) fliptable(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/tableflip?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) water(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/water?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (ImageCommands) wolverine(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/wolverine?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (ImageCommands) womanyellingatcat(ctx *dgc.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Event.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Event.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/womanyellingatcat?woman=%s&cat=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
}

func (ImageCommands) pikachu(ctx *dgc.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/pikachu", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (ImageCommands) speechbubble(ctx *dgc.Ctx) {
	link := utils.Arguments["[link/attachment]"]
	if link == "" {
		if len(ctx.Event.Message.Attachments) > 0 {
			link = ctx.Event.Message.Attachments[0].URL
		} else {
			_ = ctx.RespondText("Please provide a link or attachment.")
			return
		}
	}

	forceGIF := strings.ToLower(utils.Arguments["[force GIF: on/off]"])
	if forceGIF != "" && forceGIF != "on" && forceGIF != "off" {
		_ = ctx.RespondText("Invalid force GIF value. Use `on` or `off`.")
		return
	}

	format := "png"

	if forceGIF == "on" {
		format = "gif"
	}

	url := fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/speechbubble?url=%s&format=%s", url.QueryEscape(link), format)

	httpResponse := http.Get(url)
	if httpResponse.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + httpResponse.Error.Error())
		return
	}

	_ = ctx.RespondText("https://api.lightning-bot.com" + httpResponse.Headers["location"])
}
