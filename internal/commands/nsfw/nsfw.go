package nsfw

import (
	"fmt"

	"github.com/lb-selfbot/selfbot/internal/api"
	"github.com/lb-selfbot/selfbot/internal/random"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type NSFWCommands struct{ services *services.Services }

func NewNSFWCommands(s *services.Services) *NSFWCommands {
	return &NSFWCommands{services: s}
}

func (c *NSFWCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

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

	commands := make([]*router.Command, 0, len(nsfwCommands))
	for name := range nsfwCommands {
		commands = append(commands, &router.Command{
			Name:        name,
			Description: fmt.Sprintf("Send a random %s picture", name),
			Handler:     c.generateFunc(name, nsfwCommands[name]),
		})
	}

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "NSFW",
		Description: "NSFW commands",
	}, commands)
}

func (c *NSFWCommands) generateFunc(name string, subReddits []string) func(ctx *router.Ctx) {
	// This makes sure that the function is not overwritten by the next command.
	// If you put the function in the RegisterCommands function, the element
	// variable will always be the last element in the command array.
	return func(ctx *router.Ctx) {
		image, success := api.GetImageFromSubreddit(random.Choice(subReddits))
		if !success {
			ctx.RespondText("An error has occurred.")
			return
		}

		c.services.Embed.NewEmbed().
			SetDescription(utils.String.Capitalize(name)).
			SetImage(image.ImageURL).
			Send(ctx)
	}
}
