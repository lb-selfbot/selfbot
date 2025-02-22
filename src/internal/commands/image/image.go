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
	"github.com/LightningDev1/dgc"
)

func (c *ImageCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Image", "Image commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "tweet",
		Description: "Tweet a message",
		Arguments:   "<username> <text>",
		Handler:     c.tweet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "meme",
		Description: "Send a random meme",
		Handler:     c.meme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ytcomment",
		Description: "Make a fake YouTube comment",
		Arguments:   "<user> <text>",
		Handler:     c.ytcomment,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pokemon",
		Description: "Send a random Pokémon image",
		Handler:     c.pokemon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wasted",
		Description: "Add a wasted overlay on the users avatar",
		Arguments:   "[user]",
		Handler:     c.wasted,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggered",
		Description: "Add a triggered overlay on the users avatar",
		Arguments:   "[user]",
		Handler:     c.triggered,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "glass",
		Description: "Add a glass overlay on the users avatar",
		Arguments:   "[user]",
		Handler:     c.glass,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "greyscale",
		Description: "Add a greyscale overlay on the users avatar",
		Arguments:   "[user]",
		Handler:     c.greyscale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invert",
		Description: "Invert the users avatar",
		Arguments:   "[user]",
		Handler:     c.invert,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "invertgreyscale",
		Description: "Invert and greyscale the users avatar",
		Arguments:   "[user]",
		Handler:     c.invertgreyscale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "brightness",
		Description: "Make the users avatar brighter",
		Arguments:   "[user]",
		Handler:     c.brightness,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "red",
		Description: "Make the users avatar red",
		Arguments:   "[user]",
		Handler:     c.red,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "green",
		Description: "Make the users avatar green",
		Arguments:   "[user]",
		Handler:     c.green,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blue",
		Description: "Make the users avatar blue",
		Arguments:   "[user]",
		Handler:     c.blue,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pixelate",
		Description: "Pixelate the users avatar",
		Arguments:   "[user]",
		Handler:     c.pixelate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "comrade",
		Description: "Add comrade overlay on users avatar",
		Arguments:   "[user]",
		Handler:     c.comrade,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "missionpassed",
		Description: "Add Mission Passed overlay on users avatar",
		Arguments:   "[user]",
		Handler:     c.missionpassed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jail",
		Description: "Add jail overlay on users avatar",
		Arguments:   "[user]",
		Handler:     c.jail,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nobitches",
		Description: "Make a no bitches meme",
		Arguments:   "<text>",
		Handler:     c.nobitches,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "threat",
		Description: "Make biggest threats meme",
		Arguments:   "[user]",
		Handler:     c.threat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "threshold",
		Description: "Add a threshold filter to a users avatar",
		Arguments:   "<user>",
		Handler:     c.threshold,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blurple",
		Description: "Make the users avatar blurple",
		Handler:     c.blurple,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blurple2",
		Description: "Make the users avatar new blurple",
		Handler:     c.blurple2,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "color",
		Description: "Make the users avatar a custom color",
		Arguments:   "<user> <color>",
		Handler:     c.color,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sostupid",
		Description: "Make an it's so stupid meme",
		Arguments:   "<user> <text>",
		Handler:     c.sostupid,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "simpcard",
		Description: "Make a simp card for a user",
		Arguments:   "<user>",
		Handler:     c.simpcard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "horny",
		Description: "Make a horny card for a user",
		Arguments:   "<user>",
		Handler:     c.horny,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lgbt",
		Description: "Add an lgbt border to the user",
		Arguments:   "<user>",
		Handler:     c.lgbt,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pansexual",
		Description: "Add a pansexual border to the user",
		Arguments:   "<user>",
		Handler:     c.pansexual,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nonbinary",
		Description: "Add a nonbinary border to the user",
		Arguments:   "<user>",
		Handler:     c.nonbinary,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lesbian2",
		Description: "Add a lesbian border to the user",
		Arguments:   "<user>",
		Handler:     c.lesbian,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bisexual",
		Description: "Add a bisexual border to the user",
		Arguments:   "<user>",
		Handler:     c.bisexual,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trans",
		Description: "Add a trans border to the user",
		Arguments:   "<user>",
		Handler:     c.trans,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "oogway",
		Description: "Make a master oogway meme",
		Arguments:   "<text>",
		Handler:     c.oogway,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "oogway2",
		Description: "Make a master oogway meme",
		Arguments:   "<text>",
		Handler:     c.oogway2,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "clyde",
		Description: "Make Clyde say something",
		Arguments:   "<text>",
		Handler:     c.clyde,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gay",
		Description: "Add a rainbow overlay on the users avatar",
		Arguments:   "<user>",
		Handler:     c.gay,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "distractedbf",
		Description: "Make a distracted boyfriend meme with the users",
		Arguments:   "<user1> <user2>",
		Handler:     c.distractedbf,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "captcha",
		Description: "Make a captcha meme",
		Arguments:   "<user> <text>",
		Handler:     c.captcha,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "whowouldwin",
		Description: "Make a who would win meme",
		Arguments:   "<user1> <user2>",
		Handler:     c.whowouldwin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jpeg",
		Description: "JPEG a users avatar",
		Arguments:   "<user>",
		Handler:     c.jpeg,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lolice",
		Description: "Make a lolice meme",
		Arguments:   "<user>",
		Handler:     c.lolice,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kannagen",
		Description: "Make a kannagen meme",
		Arguments:   "<user1> <user2>",
		Handler:     c.kannagen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iphonex",
		Description: "Put users profile picture on an iPhone X",
		Arguments:   "<user1> <user2>",
		Handler:     c.iphonex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "awooify",
		Description: "Awooify the user",
		Arguments:   "<user1> <user2>",
		Handler:     c.awooify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trap",
		Description: "Make a trap card",
		Arguments:   "<user>",
		Handler:     c.trap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trumptweet",
		Description: "Make a fake Trump tweet",
		Arguments:   "<text>",
		Handler:     c.trumptweet,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deepfry",
		Description: "Deepfry the users avatar",
		Arguments:   "<text>",
		Handler:     c.deepfry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blurpify",
		Description: "Blurpify the users avatar",
		Arguments:   "<text>",
		Handler:     c.blurpify,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "phcomment",
		Description: "Make the user comment something on PornHub",
		Arguments:   "<user> <text>",
		Handler:     c.phcomment,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fact",
		Description: "Make a fact meme",
		Arguments:   "<user>",
		Handler:     c.fact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trash",
		Description: "Make a trash meme",
		Arguments:   "<user>",
		Handler:     c.trash,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stickbug",
		Description: "Make a stickbug meme",
		Arguments:   "<user>",
		Handler:     c.stickbug,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "magik",
		Description: "Make a users avatar magik",
		Arguments:   "<user>",
		Handler:     c.magik,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "randomgif",
		Description: "Send a random GIF from Giphy",
		Handler:     c.randomgif,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "searchgif",
		Description: "Search GIFs on Giphy",
		Arguments:   "<query>",
		Handler:     c.searchgif,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "communist",
		Description: "Make the user a communist",
		Arguments:   "<user>",
		Handler:     c.communist,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snow",
		Description: "Add a snow overlay the the users avatar",
		Arguments:   "<user>",
		Handler:     c.snow,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hub",
		Description: "Make a PornHub logo with custom text",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.hub,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "floor",
		Description: "Make a floor meme",
		Arguments:   "<user> <text>",
		Handler:     c.floor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jokeoverhead",
		Description: "When someone doesn't understand a joke",
		Arguments:   "<user>",
		Handler:     c.jokeoverhead,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "throwaway",
		Description: "Throw a user away",
		Arguments:   "<user>",
		Handler:     c.throwaway,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "salty",
		Description: "Make a user salty",
		Arguments:   "<user>",
		Handler:     c.salty,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "scroll",
		Description: "Make a scroll meme",
		Arguments:   "<text>",
		Handler:     c.scroll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shame",
		Description: "Shame a user",
		Arguments:   "<user>",
		Handler:     c.shame,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ship",
		Description: "Ship users",
		Arguments:   "<user1> <user2>",
		Handler:     c.ship,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "what",
		Description: "Make a what meme",
		Arguments:   "<user>",
		Handler:     c.what,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sadcat",
		Description: "Send an image of a sad cat",
		Handler:     c.sadcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bird",
		Description: "Send an image of a bird",
		Handler:     c.bird,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cat",
		Description: "Send an image of a cat",
		Handler:     c.cat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catgif",
		Description: "Send a GIF of a cat",
		Handler:     c.catgif,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "catsays",
		Description: "Make a cat say something",
		Arguments:   "<text>",
		Handler:     c.catsays,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dog",
		Description: "Send an image of a dog",
		Handler:     c.dog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "panda",
		Description: "Send an image of a panda",
		Handler:     c.panda,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "redpanda",
		Description: "Send an image of a red panda",
		Handler:     c.redpanda,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "koala",
		Description: "Send an image of a koala",
		Handler:     c.koala,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "raccoon",
		Description: "Send an image of a raccoon",
		Handler:     c.raccoon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fox",
		Description: "Send an image of a fox",
		Handler:     c.fox,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shibe",
		Description: "Send an image of a shiba inu/shibe/doge",
		Handler:     c.shibe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "whale",
		Description: "Send an image of a whale",
		Handler:     c.whale,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "kangaroo",
		Description: "Send an image of a kangaroo",
		Handler:     c.kangaroo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bunny",
		Description: "Send an image of a bunny",
		Handler:     c.bunny,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "lion",
		Description: "Send an image of a lion",
		Handler:     c.lion,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "frog",
		Description: "Send an image of a frog",
		Handler:     c.frog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "duck",
		Description: "Send an image of a duck",
		Handler:     c.duck,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "penguin",
		Description: "Send an image of a penguin",
		Handler:     c.penguin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dolphin",
		Description: "Send an image of a dolphin",
		Handler:     c.dolphin,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bear",
		Description: "Send an image of a bear",
		Handler:     c.bear,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "axolotl",
		Description: "Send an image of an axolotl",
		Handler:     c.axolotl,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "capybara",
		Description: "Send an image of a capybara",
		Handler:     c.capybara,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "balloon",
		Description: "Make a balloon meme",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.balloon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "abandon",
		Description: "Make an abandon meme",
		Arguments:   "<text>",
		Handler:     c.abandon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "aborted",
		Description: "Make an aborted meme",
		Arguments:   "<user>",
		Handler:     c.aborted,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "affect",
		Description: "Make an affect meme",
		Arguments:   "<user>",
		Handler:     c.affect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "airpods",
		Description: "Add AirPods to the users avatar",
		Arguments:   "<user>",
		Handler:     c.airpods,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "america",
		Description: "Patriotize the user",
		Arguments:   "<user>",
		Handler:     c.america,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "armor",
		Description: "Make an armor meme",
		Arguments:   "<text>",
		Handler:     c.armor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bed",
		Description: "Make a bed meme",
		Arguments:   "<user1> <user2>",
		Handler:     c.bed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bongocat",
		Description: "Make a bongocat meme",
		Arguments:   "<user>",
		Handler:     c.bongocat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "boo",
		Description: "Make a boo meme",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.boo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bigbrain",
		Description: "Make a big brain meme",
		Arguments:   "<\"item 1\"> <\"item 2\"> <\"item 3\"> <\"item 4\">",
		Handler:     c.bigbrain,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "brazzers",
		Description: "Make a brazzers meme",
		Arguments:   "<user>",
		Handler:     c.brazzers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "byemom",
		Description: "Make a byemom meme",
		Arguments:   "<user> <text>",
		Handler:     c.byemom,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cancer",
		Description: "Make a Squidward cancer meme",
		Arguments:   "<user>",
		Handler:     c.cancer,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "changemymind",
		Description: "Make a change my mind meme",
		Arguments:   "<text>",
		Handler:     c.changemymind,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cheating",
		Description: "Make a cheating meme",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.cheating,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "citation",
		Description: "Make a citation meme",
		Arguments:   "<\"text1\"> <\"text2\"> <\"text3\">",
		Handler:     c.citation,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "communism",
		Description: "Make the user a communist",
		Arguments:   "<user>",
		Handler:     c.communism,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "confusedcat",
		Description: "Make a confused cat meme",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.confusedcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "corporate",
		Description: "Make a corporate meme",
		Arguments:   "<user>",
		Handler:     c.corporate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cry",
		Description: "Make a cry meme",
		Arguments:   "<text>",
		Handler:     c.cry,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dab",
		Description: "Make a dab meme",
		Arguments:   "<user>",
		Handler:     c.dab,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dank",
		Description: "Make a dank meme",
		Arguments:   "<user>",
		Handler:     c.dank,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "delete",
		Description: "Make a delete meme",
		Arguments:   "<user>",
		Handler:     c.delete,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "disability",
		Description: "Make a disability meme",
		Arguments:   "<user>",
		Handler:     c.disability,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "doglemon",
		Description: "Make a dog lemon meme",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.doglemon,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "door",
		Description: "Make a door meme",
		Arguments:   "<user>",
		Handler:     c.door,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "egg",
		Description: "Make a egg meme",
		Arguments:   "<user>",
		Handler:     c.egg,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "excuseme",
		Description: "Make an excuse me meme",
		Arguments:   "<text>",
		Handler:     c.excuseme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "expanddong",
		Description: "Make an expand dong meme",
		Arguments:   "<text>",
		Handler:     c.expanddong,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "expandingwwe",
		Description: "Make an expanding WWE meme",
		Arguments:   "<\"item 1\"> <\"item 2\"> <\"item 3\"> <\"item 4\"> <\"item 5\">",
		Handler:     c.expandingwwe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "recaptcha",
		Description: "Make a recaptcha",
		Arguments:   "<text>",
		Handler:     c.recaptcha,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "qrcode",
		Description: "Make a QR Code",
		Arguments:   "<text>",
		Handler:     c.qrcode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "darkqrcode",
		Description: "Make a dark QR Code",
		Arguments:   "<text>",
		Handler:     c.darkqrcode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "rule34",
		Description: "Random rule 34 image",
		Handler:     c.rule34,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "achievement",
		Description: "Make a minecraft achievement",
		Arguments:   "<text>",
		Handler:     c.achievement,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amiajoke",
		Description: "Make an am I a joke meme",
		Arguments:   "<user>",
		Handler:     c.amiajoke,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "bad",
		Description: "Make a bad meme",
		Arguments:   "<user>",
		Handler:     c.bad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calling",
		Description: "Make a calling meme",
		Arguments:   "<text>",
		Handler:     c.calling,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "didyoumean",
		Description: "Make a did you mean meme",
		Arguments:   "<\"top\"> <\"bottom\">",
		Handler:     c.didyoumean,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "drake",
		Description: "Make a drake meme",
		Arguments:   "<\"top\"> <\"bottom\">",
		Handler:     c.drake,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ayano",
		Description: "Make an ayano meme",
		Arguments:   "<\"top\"> <\"bottom\">",
		Handler:     c.ayano,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "facts",
		Description: "Make a facts meme",
		Arguments:   "<text>",
		Handler:     c.facts,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "blur",
		Description: "Blur the user's avatar",
		Arguments:   "<user>",
		Handler:     c.blur,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "flipav",
		Description: "Flip the user's avatar",
		Arguments:   "<user>",
		Handler:     c.flipav,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mirror",
		Description: "Mirror the user's avatar",
		Arguments:   "<user>",
		Handler:     c.mirror,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wide",
		Description: "Make the user's avatar wide",
		Arguments:   "<user>",
		Handler:     c.wide,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sepia",
		Description: "Make the user's avatar sepia",
		Arguments:   "<user>",
		Handler:     c.sepia,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adios",
		Description: "Make an adios meme",
		Arguments:   "<user>",
		Handler:     c.adios,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "batmanslap",
		Description: "Make a batmanslap meme",
		Arguments:   "<user1> <user2> <\"text1\"> <\"text2\">",
		Handler:     c.batmanslap,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "carreverse",
		Description: "Make a car reverse meme",
		Arguments:   "<text>",
		Handler:     c.carreverse,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "drip",
		Description: "Give the user some drip",
		Arguments:   "<user>",
		Handler:     c.drip,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "impostor",
		Description: "Make an Among Us impostor meme",
		Arguments:   "<text>",
		Handler:     c.impostor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notimpostor",
		Description: "Make an Among Us not impostor meme",
		Arguments:   "<text>",
		Handler:     c.notimpostor,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "emergencymeeting",
		Description: "Make an Among Us emergency meeting meme",
		Arguments:   "<text>",
		Handler:     c.emergencymeeting,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "firsttime",
		Description: "Make a first time meme",
		Arguments:   "<user>",
		Handler:     c.firsttime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "grave",
		Description: "Make a grave meme",
		Arguments:   "<user>",
		Handler:     c.grave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "iamspeed",
		Description: "Make an i am speed meme",
		Arguments:   "<user>",
		Handler:     c.iamspeed,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "icanmilkyou",
		Description: "Make an i can milk you meme",
		Arguments:   "<user1> <user2>",
		Handler:     c.icanmilkyou,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "heaven",
		Description: "Make a heaven meme",
		Arguments:   "<user>",
		Handler:     c.heaven,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "npc",
		Description: "Make an npc meme",
		Arguments:   "<\"text1\"> <\"text2\">",
		Handler:     c.npc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stonks",
		Description: "Make a stonks meme",
		Arguments:   "<user>",
		Handler:     c.stonks,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notstonks",
		Description: "Make a not stonks meme",
		Arguments:   "<user>",
		Handler:     c.notstonks,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fliptable",
		Description: "Make the user flip a table",
		Arguments:   "<user>",
		Handler:     c.fliptable,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "water",
		Description: "Make a water meme",
		Arguments:   "<text>",
		Handler:     c.water,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "wolverine",
		Description: "Make a wolverine meme",
		Arguments:   "<user>",
		Handler:     c.wolverine,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "womanyellingatcat",
		Description: "Make a woman yelling at cat meme",
		Arguments:   "<user1> <user2>",
		Handler:     c.womanyellingatcat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pikachu",
		Description: "Send a random Pikachu image",
		Handler:     c.pikachu,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "speechbubble",
		Description: "Add a speech bubble to an image",
		Arguments:   "[force GIF: on/off] [link/attachment]",
		Handler:     c.speechbubble,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "quotereply",
		Description: "Quote the message you replied to",
		Arguments:   "[background url]",
		Handler:     c.quotereply,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "quotemessage",
		Description: "Quote a message",
		Arguments:   "<message id> [channel id]",
		Handler:     c.quotemessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "quoteuser",
		Description: "Quote a user",
		Arguments:   "<@user> <\"text\"> [background url]",
		Handler:     c.quoteuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "quotetext",
		Description: "Quote using text",
		Arguments:   "<\"author\"> <\"text\"> [background url]",
		Handler:     c.quotetext,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "discordquotereply",
		Description: "Quote the Discord message you replied to",
		Handler:     c.discordquotereply,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discordquotemessage",
		Description: "Quote a Discord message",
		Arguments:   "<message id> [channel id]",
		Handler:     c.discordquotemessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discordquoteuser",
		Description: "Quote a user",
		Arguments:   "<@user> <\"text\">",
		Handler:     c.discordquoteuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discordquotetext",
		Description: "Quote a Discord message",
		Arguments:   "<\"author\"> <\"text\"> <avatar url>",
		Handler:     c.discordquotetext,
	})
}

func (c *ImageCommands) tweet(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf(
		"https://nekobot.xyz/api/imagegen?type=tweet&text=%s&username=%s",
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
		url.QueryEscape(utils.Arguments["<username>"]),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) meme(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://api.lightning-bot.com/api/v3/misc/meme")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
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
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf(
		"https://some-random-api.com/canvas/youtube-comment?avatar=%s&comment=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
		url.QueryEscape(user.Username),
	))
}

func (c *ImageCommands) pokemon(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/%d.png", random.Int(1, 898)))
}

func (c *ImageCommands) wasted(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/wasted?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) triggered(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/triggered?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) glass(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/glass?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) greyscale(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/greyscale?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) invert(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/invert?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) invertgreyscale(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/invertgreyscale?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) brightness(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/brightness?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) red(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/red?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) green(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/green?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) blue(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/blue?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) pixelate(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/pixelate?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) comrade(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/comrade?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) missionpassed(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/passed?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) jail(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/jail?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) nobitches(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/nobitches?no=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) threshold(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/threshold?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) blurple(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/blurple?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) blurple2(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/blurple2?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) color(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/color?avatar=%s&color=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(strings.ReplaceAll(utils.Arguments["<color>"], "#", "")),
	))
}

func (c *ImageCommands) sostupid(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/its-so-stupid?avatar=%s&dog=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
}

func (c *ImageCommands) simpcard(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/simpcard?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) horny(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/horny?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) lgbt(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/lgbt?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) pansexual(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/pansexual?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) nonbinary(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/nonbinary?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) lesbian(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/lesbian?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) bisexual(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/bisexual?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) trans(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/transgender?avatar=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) oogway(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/oogway?quote=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) oogway2(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://some-random-api.com/canvas/oogway2?quote=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) threat(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=threats&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) clyde(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=clyde&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) distractedbf(ctx *router.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=ship&user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) captcha(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=captcha&url=%s&username=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) whowouldwin(ctx *router.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=whowouldwin&user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) changemymind(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=changemymind&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) jpeg(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=jpeg&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) lolice(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=lolice&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) kannagen(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=kannagen&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) iphonex(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=iphonex&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) awooify(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=awooify&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) trap(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trap&image=%s&name=%s&author=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user.String()),
		url.QueryEscape(ctx.Message.Author.String()),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) trumptweet(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trumptweet&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) deepfry(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=deepfry&image=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) blurpify(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=blurpify&image=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) phcomment(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=phcomment&text=%s&username=%s&image=%s",
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
		url.QueryEscape(user.Username),
		url.QueryEscape(user.AvatarURL("")),
	))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) fact(ctx *router.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=fact&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) trash(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=trash&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) stickbug(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=stickbug&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) magik(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=magik&image=%s&intensity=5", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.Get("message").MustString())
}

func (c *ImageCommands) randomgif(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://api.giphy.com/v1/gifs/random?api_key=ldQeNHnpL3WcCxJE1uO8HTk17ICn8i34&tag=&rating=R&limit=1&lang=en&offset=0")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.GetPath("data", "bitly_url").MustString())
}

func (c *ImageCommands) searchgif(ctx *router.Ctx) {
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

func (c *ImageCommands) sadcat(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.alexflipnote.dev/sadcat", "file")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bird(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.alexflipnote.dev/birb", "file")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) cat(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat", "cat.jpg")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) catgif(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat/gif", "cat.gif")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) catsays(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(ctx, "https://cataas.com/cat/says/"+url.PathEscape(ctx.Arguments.Raw()), "cat.jpg")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) dog(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://dog.ceo/api/breeds/image/random", "message")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) panda(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/panda", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) redpanda(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/red_panda", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) koala(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/koala", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) raccoon(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/raccoon", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) fox(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://randomfox.ca/floof", "image")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) shibe(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://shibe.online/api/shibes?count=1&urls=true&httpsUrls=true")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.GetIndex(0).MustString())
}

func (c *ImageCommands) whale(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/whale", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) kangaroo(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/kangaroo", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bunny(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/bunny", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) lion(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/lion", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) frog(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/frog", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) duck(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/duck", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) penguin(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/penguin", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) dolphin(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/dolphin", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bear(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/bear", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) axolotl(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/axolotl", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) capybara(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://api.animality.xyz/img/capybara", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) bed(ctx *router.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/bed",
		"bed.png",
		map[string]string{
			"avatars": user1.AvatarURL("") + "," + user2.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) balloon(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/balloon",
		"balloon.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
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
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) aborted(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/aborted",
		"aborted.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) affect(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/affect",
		"affect.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) airpods(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/airpods",
		"airpods.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) america(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/america",
		"america.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
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
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) bongocat(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/bongocat",
		"bongocat.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) boo(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/boo",
		"boo.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) bigbrain(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/brain",
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

func (c *ImageCommands) brazzers(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/brazzers",
		"brazzers.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}

}

func (c *ImageCommands) byemom(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
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
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) cancer(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/cancer",
		"cancer.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) cheating(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/cheating",
		"cheating.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) citation(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/citation",
		"citation.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"] + "," + utils.Arguments["<\"text3\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) communism(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/communism",
		"communism.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) confusedcat(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/confusedcat",
		"confusedcat.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) corporate(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/corporate",
		"corporate.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
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
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) dab(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/dab",
		"dab.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) dank(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/dank",
		"dank.gif",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) delete(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/delete",
		"delete.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) disability(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/disability",
		"disability.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) doglemon(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/doglemon",
		"doglemon.png",
		map[string]string{
			"text": utils.Arguments["<\"text1\">"] + "," + utils.Arguments["<\"text2\">"],
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) door(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/door",
		"door.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) egg(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/egg",
		"egg.png",
		map[string]string{
			"avatars": user.AvatarURL(""),
		},
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
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
		_ = ctx.RespondText("An error has occurred: " + err.Error())
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
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) expandingwwe(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		"https://imgen.lightning-bot.com/expandingwwe",
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

func (c *ImageCommands) qrcode(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		fmt.Sprintf("https://www.qrcoder.co.uk/api/v1/?type=png&text=%s&int=1&background=FFFFFF&foreground=000000&eo=000000&ei=000000&size=8&padding=2", url.QueryEscape(ctx.Arguments.Raw())),
		"qrcode.png",
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) darkqrcode(ctx *router.Ctx) {
	err := utils.Discord.SendFileFromURL(
		ctx,
		fmt.Sprintf("https://www.qrcoder.co.uk/api/v1/?type=png&text=%s&int=1&background=000&foreground=FFFFFF&eo=FFFFFF&ei=FFFFFF&size=8&padding=2", url.QueryEscape(ctx.Arguments.Raw())),
		"qrcode.png",
	)
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
	}
}

func (c *ImageCommands) rule34(ctx *router.Ctx) {
	jsonResult, err := http.GetJson("https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&json=1")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(jsonResult.GetIndex(rand.Intn(len(jsonResult.MustArray()))).Get("file_url").MustString())
}

func (c *ImageCommands) achievement(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/achievement?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) amiajoke(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/amiajoke?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) bad(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/bad?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) calling(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/calling?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) didyoumean(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/didyoumean?top=%s&bottom=%s",
		url.QueryEscape(utils.Arguments["<\"top\">"]),
		url.QueryEscape(utils.Arguments["<\"bottom\">"]),
	))
}

func (c *ImageCommands) drake(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/drake?top=%s&bottom=%s&ayano=false",
		url.QueryEscape(utils.Arguments["<\"top\">"]),
		url.QueryEscape(utils.Arguments["<\"bottom\">"]),
	))
}

func (c *ImageCommands) ayano(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/drake?top=%s&bottom=%s&ayano=true",
		url.QueryEscape(utils.Arguments["<\"top\">"]),
		url.QueryEscape(utils.Arguments["<\"bottom\">"]),
	))
}

func (c *ImageCommands) facts(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/facts?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) blur(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=blur", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) flipav(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=flip", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) mirror(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=mirror", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) wide(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=wide", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) sepia(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=sepia", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) recaptcha(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/captcha?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) gay(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=gay", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) communist(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=communist", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) snow(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/filter?url=%s&filter=snow", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) hub(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/pornhub?text1=%s&text2=%s",
		url.QueryEscape(utils.Arguments["<\"text1\">"]),
		url.QueryEscape(utils.Arguments["<\"text2\">"]),
	))
}

func (c *ImageCommands) floor(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/floor?url=%s&text=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
}

func (c *ImageCommands) jokeoverhead(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/jokeoverhead?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) throwaway(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/trash?url=%s&trash_url=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Message.Author.AvatarURL("")),
	))
}

func (c *ImageCommands) salty(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/salty?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) scroll(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/scroll?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) shame(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/shame?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) ship(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/ship?url1=%s&url2=%s",
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user2.AvatarURL("")),
	))
}

func (c *ImageCommands) what(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://api.lightning-bot.com/api/v3/imagegen/what?url=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) adios(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/adios?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) batmanslap(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf(
		"https://vacefron.nl/api/batmanslap?text1=%s&text2=%s&batman=%s&robin=%s",
		url.QueryEscape(utils.Arguments["<\"text1\">"]),
		url.QueryEscape(utils.Arguments["<\"text2\">"]),
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(user2.AvatarURL("")),
	))
}

func (c *ImageCommands) carreverse(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/carreverse?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) drip(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/drip?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) impostor(ctx *router.Ctx) {
	crewmate := random.Choice([]string{"black", "blue", "brown", "cyan", "darkgreen", "lime", "orange", "pink", "purple", "red", "white", "yellow"})
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/ejected?name=%s&impostor=true&crewmate=%s", ctx.Arguments.Raw(), crewmate))
}

func (c *ImageCommands) notimpostor(ctx *router.Ctx) {
	crewmate := random.Choice([]string{"black", "blue", "brown", "cyan", "darkgreen", "lime", "orange", "pink", "purple", "red", "white", "yellow"})
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/ejected?name=%s&impostor=false&crewmate=%s", ctx.Arguments.Raw(), crewmate))
}

func (c *ImageCommands) emergencymeeting(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/emergencymeeting?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) firsttime(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/firsttime?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) grave(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/grave?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) iamspeed(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/iamspeed?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) icanmilkyou(ctx *router.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/icanmilkyou?user1=%s&user2=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
}

func (c *ImageCommands) heaven(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/heaven?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) npc(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/npc?text1=%s&text2=%s", url.QueryEscape(utils.Arguments["<\"text1\">"]), url.QueryEscape(utils.Arguments["<\"text2\">"])))
}

func (c *ImageCommands) stonks(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/stonks?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) notstonks(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/stonks?user=%s&notstonks=true", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) fliptable(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/tableflip?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) water(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/water?text=%s", url.QueryEscape(ctx.Arguments.Raw())))
}

func (c *ImageCommands) wolverine(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/wolverine?user=%s", url.QueryEscape(user.AvatarURL(""))))
}

func (c *ImageCommands) womanyellingatcat(ctx *router.Ctx) {
	user1, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user1 = ctx.Message.Author
	}
	user2, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(1))
	if !success {
		user2 = ctx.Message.Author
	}
	_ = ctx.RespondText(fmt.Sprintf("https://vacefron.nl/api/womanyellingatcat?woman=%s&cat=%s", url.QueryEscape(user1.AvatarURL("")), url.QueryEscape(user2.AvatarURL(""))))
}

func (c *ImageCommands) pikachu(ctx *router.Ctx) {
	imageUrl, err := utils.Misc.GetImageUrlFromJson("https://some-random-api.com/img/pikachu", "link")
	if err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	_ = ctx.RespondText(imageUrl)
}

func (c *ImageCommands) speechbubble(ctx *router.Ctx) {
	link := utils.Arguments["[link/attachment]"]

	attachments := ctx.Message.Attachments
	if len(attachments) > 0 {
		link = attachments[0].URL
	}

	if link == "" {
		_ = ctx.RespondText("Please provide a link or attachment.")
		return
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

	_ = ctx.RespondText("https://api.lightning-bot.com" + httpResponse.Header.Get("Location"))
}

func (c *ImageCommands) quotereply(ctx *router.Ctx) {
	if ctx.Message.ReferencedMessage == nil {
		c.services.Embed.NewEmbed().
			SetDescription("You must reply to a message to use quotereply.").
			Send(ctx)
		return
	}

	m := ctx.Message.ReferencedMessage

	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(ctx.Message.ReferencedMessage.ContentWithMentionsReplaced()),
		url.QueryEscape(utils.Arguments["[background url]"]),
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

	m, err := utils.Discord.GetChannelMessage(ctx.Session, channel.ID, utils.Arguments["<message id>"])
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find message: " + err.Error()).
			Send(ctx)
		return
	}

	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(m.ContentWithMentionsReplaced()),
		url.QueryEscape(utils.Arguments["[background url]"]),
	))
}

func (c *ImageCommands) quoteuser(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find user.").
			Send(ctx)
		return
	}

	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Misc.Or(user.DisplayName, user.Username)),
		url.QueryEscape(utils.Arguments["<\"text\">"]),
		url.QueryEscape(utils.Arguments["[background url]"]),
	))
}

func (c *ImageCommands) quotetext(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/quote?author=%s&text=%s&background=%s",
		url.QueryEscape(utils.Arguments["<\"author\">"]),
		url.QueryEscape(utils.Arguments["<\"text\">"]),
		url.QueryEscape(utils.Arguments["[background url]"]),
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

	m, err := utils.Discord.GetChannelMessage(ctx.Session, channel.ID, utils.Arguments["<message id>"])
	if err != nil {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find message: " + err.Error()).
			Send(ctx)
		return
	}

	authorColor := c.getRoleColor(ctx, m.Author.ID)

	_ = ctx.RespondText(fmt.Sprintf(
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

	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&author_color=%s&text=%s&timestamp=%s",
		url.QueryEscape(utils.Misc.Or(m.Author.DisplayName, m.Author.Username)),
		url.QueryEscape(m.Author.AvatarURL("")),
		url.QueryEscape(authorColor),
		url.QueryEscape(m.ContentWithMentionsReplaced()),
		url.QueryEscape(m.Timestamp.Format("01/02/2006 3:04 PM")),
	))
}

func (c *ImageCommands) discordquoteuser(ctx *router.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		c.services.Embed.NewEmbed().
			SetDescription("Could not find user.").
			Send(ctx)
		return
	}

	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&text=%s",
		url.QueryEscape(utils.Misc.Or(user.DisplayName, user.Username)),
		url.QueryEscape(user.AvatarURL("")),
		url.QueryEscape(ctx.Arguments.RawFrom(1)),
	))
}

func (c *ImageCommands) discordquotetext(ctx *router.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf(
		"https://api.lightning-bot.com/api/v3/imagegen/discordquote?author=%s&author_url=%s&text=%s",
		url.QueryEscape(utils.Arguments["<\"author\">"]),
		url.QueryEscape(utils.Arguments["<avatar url>"]),
		url.QueryEscape(utils.Arguments["<\"text\">"]),
	))
}

type ImageCommands struct{ services *services.Services }

func NewImageCommands(s *services.Services) *ImageCommands {
	return &ImageCommands{services: s}
}
func (c *ImageCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

}
