package commands

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type NSFWCommands struct{}

func (n NSFWCommands) RegisterCommands(router *dgc.Router) {
	nsfwCommands := map[string][]string{
		"thighs":  {"girlsinyogapants", "Thighs", "thighhighs", "ThickThighs", "UnderwearGW", "datgap", "leggingsgonewild", "pawg", "hipcleavage", "legs", "pantyhose"},
		"ass":     {"ass", "asstastic", "facedownassup", "assinthong", "bigasses"},
		"anal":    {"anal", "analgw", "painal"},
		"hentai":  {"hentai", "hentai_gif", "westernhentai", "hentai_irl"},
		"teen":    {"legalteens", "collegesluts", "adorableporn", "legalteensXXX", "gonewild18"},
		"milf":    {"milf", "gonewild30plus", "preggoporn", "realmoms"},
		"celeb":   {"celebnsfw", "WatchItForThePlot"},
		"thick":   {"gonewildcurvy", "curvy", "gonewildplus", "thick"},
		"pussy":   {"pussy", "rearpussy", "innie", "LabiaGW"},
		"boobs":   {"boobies", "TittyDrop", "boltedontits", "boobbounce"},
		"lesbian": {"lesbians", "StraightGirlsPlaying", "girlskissing"},
		"porn":    {"porn"},
		"gayporn": {"gaynsfw"},
		"feet":    {"feet"},
		"amateur": {"realgirls", "amateur", "homemadexxx"},
		"cam":     {"camwhores", "camsluts", "streamersgonewild"},
		"curvy":   {"CurvyCelebs", "curvygirls", "curvyxo"},
		"couples": {"gwcouples", "gonewildcouples", "gwcumsluts"},
		"bdsm":    {"bdsm", "bondage", "bdsmgw"},
		"blowjob": {"blowjobs", "lipsthatgrip", "deepthroat"},
		"nipples": {"pokies", "ghostnipples", "nipples"},
		"titfuck": {"titfuck", "clothedtitfuck"},
		"fit":     {"fitgirls", "fitnakedgirls"},
		"petite":  {"bustypetite", "dirtysmall", "petitegonewild"},
		"femboy":  {"femboys"},
	}
	for name, subReddits := range nsfwCommands {
		router.RegisterCmd(&dgc.Command{
			Name:        name,
			Description: fmt.Sprintf("Send a random %s picture", name),
			Usage:       "[p]" + name,
			Handler:     n.generateFunc(name, subReddits),
		})
	}
}

func (NSFWCommands) generateFunc(name string, subReddits []string) func(ctx *dgc.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.
	return func(ctx *dgc.Ctx) {
		image, success := api.GetImageFromSubreddit(utils.RandomChoice(subReddits).(string))
		if !success {
			ctx.RespondText("An error has occurred.")
			return
		}
		embed.NewEmbed().
			SetDescription(utils.Capitalize(name)).
			SetImage(image.ImageURL).
			Send(ctx)
	}
}
