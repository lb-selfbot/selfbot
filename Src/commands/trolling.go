package commands

import (
	"fmt"
	"strings"

	"github.com/LightningDev/LBGo/constants"
	"github.com/LightningDev/LBGo/encodings"
	"github.com/LightningDev/LBGo/utils"
	"github.com/LightningDev1/dgc"
)

type TrollingCommands struct{}

func (t TrollingCommands) RegisterCommands(router *dgc.Router) {
	router.RegisterCmd(&dgc.Command{
		Name:        "purgehack",
		Description: "Purge a chat without permissions",
		Usage:       "[p]purgehack",
		Handler:     t.purgehack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ghostping",
		Description: "Ghostping a user",
		Usage:       "[p]ghostping <user>",
		Handler:     t.ghostping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "gprole",
		Description: "Ghostping a role",
		Usage:       "[p]gprole <role>",
		Handler:     t.gprole,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "token",
		Description: "Send the first part of someones token",
		Usage:       "[p]token <user>",
		Handler:     t.token,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "trollnitro",
		Description: "Send a fake nitro link with an image",
		Usage:       "[p]trollnitro <user>",
		Handler:     t.trollnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeninvite",
		Description: "Hide an invite behind text",
		Usage:       "[p]hiddeninvite <invite> <text>",
		Handler:     t.hiddeninvite,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddenping",
		Description: "Hide a ping behind text",
		Usage:       "[p]hiddenping <user> <text>",
		Handler:     t.hiddenping,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hiddeneveryone",
		Description: "Hide an everyone ping behind text",
		Usage:       "[p]hiddeneveryone <text>",
		Handler:     t.hiddeneveryone,
	})
	// forcenick
	// forcekick
	// noleave
	// copy
	// pinspam
	// smartcopy
	// replycopy
	// deletemessages
	// reactuser
	// unreact
	// forcedisconnect
}

func (t TrollingCommands) purgehack(ctx *dgc.Ctx) {
	ctx.RespondText("ﾠ" + strings.Repeat("\n", 1998) + "ﾠ")
}

func (t TrollingCommands) ghostping(ctx *dgc.Ctx) {
	// noop
}

func (t TrollingCommands) gprole(ctx *dgc.Ctx) {
	// noop
}

func (t TrollingCommands) token(ctx *dgc.Ctx) {
	user, success := utils.TryGetUserFromUserMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user.")
		return
	}
	ctx.RespondText(encodings.Base64(user.ID, false))
}

func (t TrollingCommands) trollnitro(ctx *dgc.Ctx) {
	ctx.RespondText(fmt.Sprintf("dıscord.gift/%s %s https://media.discordapp.net/attachments/759817334121824297/763209007170453544/xedy9ugkqo621.png", utils.RandomString(16), constants.InvisibleText))
}

func (t TrollingCommands) hiddeninvite(ctx *dgc.Ctx) {
	invite := ctx.Arguments.Get(0)
	text := utils.RemoveWords(ctx.Arguments.Raw(), 1)
	ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.InvisibleText, invite))
}

func (t TrollingCommands) hiddenping(ctx *dgc.Ctx) {
	user := ctx.Arguments.Get(0)
	text := utils.RemoveWords(ctx.Arguments.Raw(), 1)
	ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.InvisibleText, user))
}

func (t TrollingCommands) hiddeneveryone(ctx *dgc.Ctx) {
	text := ctx.Arguments.Raw()
	ctx.RespondText(fmt.Sprintf("%s %s @everyone", text, constants.InvisibleText))
}
