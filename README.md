# LightningBot

## Completion

- Text Commands:                 100%
- Encoding Commands:             100%
- Codeblock Commands:            100%
- Crypto Commands:               100%
- Fun Commands:                  100%
- Trolling Commands:              60%
- Event Reaction Commands:         0%
- Image Commands:                 10%
- Network Commands:              100%
- Moderation Commands:           100%
- Util Commands:                   0%
- Spotify Commands:              100%
- Game Commands:                 100%
- Activity Commands:             100%
- Animation Commands:             75%
- Config Commands:               100%
- RPC Commands:                  100%
- NSFW Commands:                 100%
- Reddit Commands:               100%
- Share Commands:                  0%
- Custom Commands:                 0%
- Help Command:                    0%
- Event Listeners:                 5%
- IPC:                            75%
- encoding/json to simplejson:   100%
- Better argument handling:      100%
- Remove usage of ctx.Arguments: 100%

## Non-completed

- Trolling Commands:              60%
- Event Reaction Commands:         0%
- Image Commands:                 10%
- Util Commands:                   0%
- Animation Commands:             75%
- Share Commands:                  0%
- Custom Commands:                 0%
- Help Command:                    0%
- Event Listeners:                 5%
- IPC:                            75%


# Nekobot Avatar

func (ImageCommands) threat(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		user = ctx.Event.Author
	}
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=threats&url=%s", url.QueryEscape(user.AvatarURL(""))))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}

# Nekobot Text

func (ImageCommands) tweet(ctx *dgc.Ctx) {
	jsonResult, err := http.GetJson(fmt.Sprintf("https://nekobot.xyz/api/imagegen?type=tweet&text=%s", url.QueryEscape(ctx.Arguments.Raw())))
	if err != nil {
		ctx.RespondText("An error has occurred: " + err.Error())
		return
	}
	ctx.RespondText(jsonResult.Get("message").MustString())
}