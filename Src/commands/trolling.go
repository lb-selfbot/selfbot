package commands

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/encodings"
	"github.com/LightningDev1/LB-Selfbot/utils"
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

func (TrollingCommands) purgehack(ctx *dgc.Ctx) {
	ctx.RespondText("ﾠ" + strings.Repeat("\n", 1998) + "ﾠ")
}

func (TrollingCommands) ghostping(ctx *dgc.Ctx) {
	// noop
}

func (TrollingCommands) gprole(ctx *dgc.Ctx) {
	// noop
}

func (TrollingCommands) token(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		ctx.RespondText("Could not find user.")
		return
	}
	ctx.RespondText(encodings.Base64(user.ID, false))
}

func (TrollingCommands) trollnitro(ctx *dgc.Ctx) {
	ctx.RespondText(fmt.Sprintf("dıscord.gift/%s %s https://media.discordapp.net/attachments/759817334121824297/763209007170453544/xedy9ugkqo621.png", utils.Random.String(16), constants.InvisibleText))
}

func (TrollingCommands) hiddeninvite(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.InvisibleText, utils.Arguments["<invite>"]))
}

func (TrollingCommands) hiddenping(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.InvisibleText, utils.Arguments["<user>"]))
}

func (TrollingCommands) hiddeneveryone(ctx *dgc.Ctx) {
	ctx.RespondText(fmt.Sprintf("%s %s @everyone", ctx.Arguments.Raw(), constants.InvisibleText))
}
