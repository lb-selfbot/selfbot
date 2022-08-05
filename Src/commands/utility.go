package commands

import (
	"bytes"
	"fmt"
	"github.com/LightningDev1/LB-Selfbot/error_handling"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type UtilityCommands struct{}

func (u UtilityCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Utility", "Utility commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "copyserver",
		Description: "Make an exact copy of a server",
		Usage:       "[p]copyserver [server id]",
		Handler:     u.copyserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leaveservers",
		Description: "Leave all of your servers",
		Usage:       "[p]leaveservers",
		Handler:     u.leaveservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leavegroups",
		Description: "Leave all of your groups",
		Usage:       "[p]leavegroups",
		Handler:     u.leavegroups,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "ack",
		Description: "Mark a server as read",
		Usage:       "[p]ack [server id]",
		Handler:     u.ack,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tts",
		Description: "Send a Text-To-Speech message",
		Usage:       "[p]tts <language> <text>",
		Handler:     u.tts,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "translate",
		Description: "Translate text to a language",
		Usage:       "[p]translate <language> <text>",
		Handler:     u.translate,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "detectlang",
		Description: "Detect a texts language",
		Usage:       "[p]detectlang <text>",
		Handler:     u.detectlang,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "english",
		Description: "Translate text to English",
		Usage:       "[p]english <text>",
		Handler:     u.english,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "adminservers",
		Description: "See all servers you have administator in",
		Usage:       "[p]adminservers",
		Handler:     u.adminservers,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "admingroups",
		Description: "See all groupds you have owner in",
		Usage:       "[p]admingroups",
		Handler:     u.admingroups,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "optionpoll",
		Description: "Create a poll with options",
		Usage:       "[p]optionpoll <\"question\"> <\"option1\"> <\"option2\">",
		Handler:     u.optionpoll,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leaveserver",
		Description: "Leave a server",
		Usage:       "[p]leaveserver [server id]",
		Handler:     u.leaveserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deleteserver",
		Description: "Delete a server",
		Usage:       "[p]deleteserver [server id]",
		Handler:     u.deleteserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "systeminfo",
		Description: "See info about your system",
		Usage:       "[p]systeminfo",
		Handler:     u.systeminfo,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpemojis",
		Description: "Download all the server emojis",
		Usage:       "[p]dumpemojis [server id]",
		Handler:     u.dumpemojis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "addemojis",
		Description: "Add all images in the folder to the server",
		Usage:       "[p]addemojis <\"folder\"> [server id]",
		Handler:     u.addemojis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "findserver",
		Description: "Find a server by ID",
		Usage:       "[p]findserver <server id>",
		Handler:     u.findserver,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "howold",
		Description: "See how old a user is",
		Usage:       "[p]howold <user>",
		Handler:     u.howold,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "passgen",
		Description: "Generate a password",
		Usage:       "[p]passgen <character amount>",
		Handler:     u.passgen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "tokengen",
		Description: "Generate a Discord token",
		Usage:       "[p]tokengen",
		Handler:     u.tokengen,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "nickscan",
		Description: "Search for servers with a nickname",
		Usage:       "[p]nickscan",
		Handler:     u.nickscan,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "emoji",
		Description: "Send bigger emoji (emoji url)",
		Usage:       "[p]emoji <emoji>",
		Handler:     u.emoji,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "chatdump",
		Description: "Dump all messages in the channel",
		Usage:       "[p]chatdump <limit>",
		Handler:     u.chatdump,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "shudownpc",
		Description: "Shut your PC down",
		Usage:       "[p]shudownpc",
		Handler:     u.shudownpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartpc",
		Description: "Restart your PC",
		Usage:       "[p]restartpc",
		Handler:     u.restartpc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "snipe",
		Description: "See the last deleted message in this channel",
		Usage:       "[p]snipe",
		Handler:     u.snipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "editsnipe",
		Description: "See the last edited message in this channel",
		Usage:       "[p]editsnipe",
		Handler:     u.editsnipe,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "repeat",
		Description: "Repeat the last used command",
		Usage:       "[p]repeat",
		Handler:     u.repeat,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "history",
		Description: "See used command history",
		Usage:       "[p]history",
		Handler:     u.history,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "website",
		Description: "Open the LightningBot website",
		Usage:       "[p]website",
		Handler:     u.website,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dashboard",
		Description: "Opens dashboard in browser",
		Usage:       "[p]dashboard",
		Handler:     u.dashboard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "restartbot",
		Description: "Restart LightningBot",
		Usage:       "[p]restartbot",
		Handler:     u.restartbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopbot",
		Description: "Stop LightningBot",
		Usage:       "[p]stopbot",
		Handler:     u.stopbot,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "changelog",
		Description: "See the latest changelog",
		Usage:       "[p]changelog",
		Handler:     u.changelog,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "uptime",
		Description: "See how long the bot has been running for",
		Usage:       "[p]uptime",
		Handler:     u.uptime,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "cls",
		Description: "Clear UI console",
		Usage:       "[p]cls",
		Handler:     u.cls,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mastercard",
		Description: "Generate a mastercard credit card",
		Usage:       "[p]mastercard",
		Handler:     u.mastercard,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa16",
		Description: "Generate a visa credit card (16 characters)",
		Usage:       "[p]visa16",
		Handler:     u.visa16,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "visa13",
		Description: "Generate a visa credit card (13 characters)",
		Usage:       "[p]visa13",
		Handler:     u.visa13,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "amex",
		Description: "Generate an amex credit card",
		Usage:       "[p]amex",
		Handler:     u.amex,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "discover",
		Description: "Generate a discover credit card",
		Usage:       "[p]discover",
		Handler:     u.discover,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "diners",
		Description: "Generate a diners credit card",
		Usage:       "[p]diners",
		Handler:     u.diners,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "jcb",
		Description: "Generate a JCB credit card",
		Usage:       "[p]jcb",
		Handler:     u.jcb,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "revavatarsearch",
		Description: "Reverse image search a users avatar",
		Usage:       "[p]revavatarsearch <user>",
		Handler:     u.revavatarsearch,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "calc",
		Description: "Open the calculator",
		Usage:       "[p]calc",
		Handler:     u.calc,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "notepad",
		Description: "Open notepad",
		Usage:       "[p]notepad",
		Handler:     u.notepad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "firstmessage",
		Description: "See the first message in the channel",
		Usage:       "[p]firstmessage [channel]",
		Handler:     u.firstmessage,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "hypesquad",
		Description: "Change your Hypesquad House",
		Usage:       "[p]hypesquad <bravery/brilliance/balance>",
		Handler:     u.hypesquad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "leavehypesquad",
		Description: "Leave Hypesquad",
		Usage:       "[p]leavehypesquad",
		Handler:     u.leavehypesquad,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "theme",
		Description: "Change your Discord theme",
		Usage:       "[p]theme <light/dark>",
		Handler:     u.theme,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "compact",
		Description: "Toggle Discord compact mode",
		Usage:       "[p]compact <on/off>",
		Handler:     u.compact,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "devmode",
		Description: "Toggle Discord developer mode",
		Usage:       "[p]devmode <on/off>",
		Handler:     u.devmode,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "sessionid",
		Description: "See LightningBot's Discord session ID",
		Usage:       "[p]sessionid",
		Handler:     u.sessionid,
	})
}

func getChannelType(channelType discordgo.ChannelType) discordgo.ChannelType {
	if channelType == discordgo.ChannelTypeGuildVoice || channelType == discordgo.ChannelTypeGuildStageVoice {
		return discordgo.ChannelTypeGuildVoice
	}
	return discordgo.ChannelTypeGuildText
}

func (UtilityCommands) copyserver(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	icon, _, _ := http.GetBytes(guild.IconURL())
	newGuild, err := ctx.Session.GuildCreate(guild.Name+" Copy", icon)
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}
	_, _ = ctx.Session.GuildEdit(newGuild.ID, discordgo.GuildParams{
		VerificationLevel:           &guild.VerificationLevel,
		DefaultMessageNotifications: int(guild.DefaultMessageNotifications),
		AfkChannelID:                guild.AfkChannelID,
		AfkTimeout:                  guild.AfkTimeout,
	})

	for _, role := range guild.Roles {
		if role.ID == guild.ID || role.Managed {
			continue
		}
		_, err := ctx.Session.GuildRoleCreate(newGuild.ID, role.Name, role.Color, role.Hoist, role.Permissions, role.Mentionable)
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}
	}

	newChannels, err := ctx.Session.GuildChannels(newGuild.ID)
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}

	for _, channel := range newChannels {
		_, err := ctx.Session.ChannelDelete(channel.ID)
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}
	}

	channels, err := ctx.Session.GuildChannels(guild.ID)
	if err != nil {
		_ = ctx.RespondText("Error copying server: " + err.Error())
		return
	}

	categories, categoryChannels := utils.Discord.MapCategories(channels)

	for _, category := range categories {
		newCategory, err := ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
			Name:                 category.Name,
			Type:                 discordgo.ChannelTypeGuildCategory,
			Position:             category.Position,
			PermissionOverwrites: category.PermissionOverwrites,
		})
		if err != nil {
			_ = ctx.RespondText("Error copying server: " + err.Error())
			return
		}
		for _, channel := range categoryChannels[category.ID] {
			_, err = ctx.Session.GuildChannelCreateComplex(newGuild.ID, discordgo.GuildChannelCreateData{
				Name:                 channel.Name,
				Type:                 getChannelType(channel.Type),
				Topic:                channel.Topic,
				Bitrate:              int(math.Min(float64(channel.Bitrate), 96000)),
				UserLimit:            channel.UserLimit,
				RateLimitPerUser:     channel.RateLimitPerUser,
				Position:             channel.Position,
				PermissionOverwrites: channel.PermissionOverwrites,
				ParentID:             newCategory.ID,
				NSFW:                 channel.NSFW,
			})
			if err != nil {
				_ = ctx.RespondText("Error copying server: " + err.Error())
				return
			}
		}
	}
	embed.NewEmbed().
		SetDescription("Server copied successfully!").
		Send(ctx)
}

func (UtilityCommands) leaveservers(ctx *dgc.Ctx) { // Admin
	//     for Server in bot.guilds:
	//         await Server.leave()
}

func (UtilityCommands) leavegroups(ctx *dgc.Ctx) { // Admin
	//     for Channel in bot.private_channels:
	//         if isinstance(Channel, discord.GroupChannel):
	//             await Channel.leave()
}

func (UtilityCommands) ack(ctx *dgc.Ctx) { // Admin
	//     Server = Server or ctx.guild
	//     await Server.ack()
}

func (UtilityCommands) tts(ctx *dgc.Ctx) {
	ttsFile, err := api.TTS(utils.String.RemoveWords(ctx.Arguments.Raw(), 1), utils.Arguments["<language>"])
	if err != nil {
		_ = ctx.RespondText("An error has occurred.")
		return
	}
	if strings.Contains(string(ttsFile), "Bad Request") {
		_ = ctx.RespondText("Invalid language. Please use an ISO 639-1 language code. See https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes for more information.")
		return
	}
	_, _ = ctx.Session.ChannelFileSend(ctx.Event.ChannelID, "tts.mp3", bytes.NewReader(ttsFile))
}

func (UtilityCommands) translate(ctx *dgc.Ctx) {
	translateResult := api.Translate(utils.String.RemoveWords(ctx.Arguments.Raw(), 1), utils.Arguments["<language>"])
	if translateResult.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + translateResult.Error.Error())
		return
	}
	embed.NewEmbed().
		AddField("Source Language", utils.String.Capitalize(constants.Languages[translateResult.SourceLanguage])).
		AddField("Destination Language", utils.String.Capitalize(constants.Languages[translateResult.DestinationLanguage])).
		AddField("Original Text", translateResult.Original).
		AddField("Translated Text", translateResult.Result).
		Send(ctx)
}

func (UtilityCommands) detectlang(ctx *dgc.Ctx) {
	translateResult := api.Translate(ctx.Arguments.Raw(), "en")
	if translateResult.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + translateResult.Error.Error())
		return
	}
	embed.NewEmbed().
		AddField("Text", translateResult.Original).
		AddField("Detected Language", utils.String.Capitalize(constants.Languages[translateResult.SourceLanguage])).
		AddField("Confidence", fmt.Sprintf("%.2f", translateResult.Confidence)).
		Send(ctx)
}

func (UtilityCommands) english(ctx *dgc.Ctx) {
	translateResult := api.Translate(ctx.Arguments.Raw(), "en")
	if translateResult.Error != nil {
		_ = ctx.RespondText("An error has occurred: " + translateResult.Error.Error())
		return
	}
	embed.NewEmbed().
		AddField("Source Language", utils.String.Capitalize(constants.Languages[translateResult.SourceLanguage])).
		AddField("Original Text", translateResult.Original).
		AddField("Translated Text", translateResult.Result).
		Send(ctx)
}

func (UtilityCommands) adminservers(ctx *dgc.Ctx) { // Admin
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Admin Servers', value=', '.join(Server.name for Server in bot.guilds if Server.me.guild_permissions.administrator)))
}

func (UtilityCommands) admingroups(ctx *dgc.Ctx) { // Admin
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name='Admin Groups', value=', '.join(str(Group) for Group in bot.private_channels if isinstance(Group, discord.GroupChannel) and Group.owner == bot.user)))
}

func (UtilityCommands) optionpoll(ctx *dgc.Ctx) {
	message := embed.NewEmbed().
		SetDescription("Poll: "+utils.Arguments["<\"question\">"]).
		AddField("👍", utils.Arguments["<\"option1\">"]).
		AddField("👎", utils.Arguments["<\"option2\">"]).
		Send(ctx)
	if message != nil {
		_ = ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, message.ID, "👍")
		_ = ctx.Session.MessageReactionAdd(ctx.Event.ChannelID, message.ID, "👎")
	}
}

func (UtilityCommands) leaveserver(ctx *dgc.Ctx) { // Admin
	//     Server = Server or ctx.guild
	//     await Server.leave()
}

func (UtilityCommands) deleteserver(ctx *dgc.Ctx) { // Admin
	//     Server = Server or ctx.guild
	//     await Server.delete()
}

func (UtilityCommands) systeminfo(ctx *dgc.Ctx) {
	//     CPUInfo = cpuinfo.get_cpu_info(); Embed = bot.DefaultEmbed().add_field(name="CPU Info", value=f"CPU Name: {CPUInfo['brand_raw']}\nBits: {CPUInfo['bits']}\nCPU Threads: {CPUInfo['count']}\nCPU Arch: {CPUInfo['arch']}/{CPUInfo['arch_string_raw']}\nCPU Base Frequency: {CPUInfo['hz_advertised_friendly']}\nCPU Usage: {round(psutil.cpu_percent())}%", inline=False).add_field(name="RAM Info", value=f"Total RAM: {round(psutil.virtual_memory().total / 1000000000, 2)}GB\nUsed RAM: {round(psutil.virtual_memory().percent)}%", inline=False)
	//     for GPU in GPUtil.getGPUs():
	//         Embed.add_field(name=f"GPU {GPU.id} Info", value=f"GPU Name: {GPU.name}\nGPU Usage: {GPU.load * 100}%\nFree VRAM: {GPU.memoryFree}MB\nUsed VRAM: {GPU.memoryUsed}MB\nTotal VRAM: {GPU.memoryTotal}MB\nGPU Temp: {GPU.temperature} °C\nGPU UUID: {GPU.uuid}", inline=False)
	//     await bot.Helper.SendMessage(ctx, embed=Embed)
}

func (UtilityCommands) dumpemojis(ctx *dgc.Ctx) {
	//     Server = Server or ctx.guild; Emojis = await Server.fetch_emojis()
	//     try: os.makedirs(f'{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Emojis\\{bot.Helper.SanitizeFilename(Server.name)}', exist_ok=True)
	//     except FileExistsError: pass
	//     for Emoji in Emojis:
	//         with open(f"{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Emojis\\{bot.Helper.SanitizeFilename(Server.name)}\\{bot.Helper.SanitizeFilename(Emoji.name)}.{str(Emoji.url).split('.')[-1]}", 'wb+') as f: shutil.copyfileobj(requests.get(Emoji.url, stream=True).raw, f)
}

func (UtilityCommands) addemojis(ctx *dgc.Ctx) {
	//     if not pathlib.Path(Folder).exists(): return await ctx.send("Invalid folder!")
	//     for Image in os.listdir(Folder):
	//         if Image[-4::1] in [".png", ".jpg", ".gif"]:
	//             with open(f"{Folder}\\{Image}", "rb") as File:
	//                 try: await ctx.guild.create_custom_emoji(name=Image[:-4], image=File.read())
	//                 except: return await ctx.send("You are missing permissions to use this command!")
}

func (UtilityCommands) findserver(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Server Name", value=Server.name).add_field(name="Created At", value=Server.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')))
}

func (UtilityCommands) howold(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="User", value=str(User)).add_field(name="Created At", value=User.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')))
}

func (UtilityCommands) passgen(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Password", value=''.join(random.choices(string.ascii_letters + string.digits + string.punctuation, k=Characters))))
}

func (UtilityCommands) tokengen(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Token", value=requests.get('https://some-random-api.ml/bottoken').json()['token']))
}

func (UtilityCommands) nickscan(ctx *dgc.Ctx) {
	//     Nicks = '\n'.join(f"{Server} | {Server.me.nick}" for Server in bot.guilds if Server.me.nick != None)
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Nicknames", value=Nicks if Nicks else "No nicknames found!"))
}

func (UtilityCommands) emoji(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Server Name", value=Server.name).add_field(name="Created At", value=Server.created_at.__format__('%A, %d. %B %Y @ %H:%M:%S')))
}

func (UtilityCommands) chatdump(ctx *dgc.Ctx) {
	//     try: os.mkdir(f'{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Dumps\\{bot.Helper.SanitizeFilename(ctx.guild.name)}')
	//     except FileExistsError: pass
	//     Messages = [message async for message in ctx.channel.history(limit=Limit)]; Messages.sort(key=lambda e: e.created_at)
	//     with open(f"{bot.Helper.GetDocumentsFolder()}\\LightningBot\\Dumps\\{bot.Helper.SanitizeFilename(ctx.guild.name)}\\{bot.Helper.SanitizeFilename(ctx.channel.name)}.txt", 'a+') as f:
	//         for M in Messages:
	//             f.write(f"[{M.created_at.strftime('%H:%M %p')}] Author: {M.author} | Server: {M.guild} | Attachments: {(', '.join(A for A in M.attachments) if M.attachments else 'None')} | Jump URL: {M.jump_url} | Content: {M.content}\n")
}

func (UtilityCommands) shudownpc(ctx *dgc.Ctx) {
	//     os.system('shutdown -s -t 0')
}

func (UtilityCommands) restartpc(ctx *dgc.Ctx) {
	//     os.system('shutdown -t 0 -r -f')
}

func (UtilityCommands) snipe(ctx *dgc.Ctx) {
	//     if DeletedMsg := bot.EventListenerClass.LastDeleted.get(ctx.channel.id, None):
	//         await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_author(name=DeletedMsg.author, icon_url=DeletedMsg.author.avatar_url).add_field(name="Content", value=DeletedMsg.content))
	//     else: return await ctx.send("No deleted message found in this channel!")
}

func (UtilityCommands) editsnipe(ctx *dgc.Ctx) {
	//     if EditedMsg := bot.EventListenerClass.LastEdited.get(ctx.channel.id, None):
	//         await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().set_author(name=EditedMsg[0].author, icon_url=EditedMsg[0].author.avatar_url).add_field(name="Before", value=EditedMsg[0].content).add_field(name="After", value=EditedMsg[1].content, inline=False))
	//     else: return await ctx.send("No edited message found in this channel!")
}

func (UtilityCommands) repeat(ctx *dgc.Ctx) {
	//     if not bot.EventListenerClass.LastUsedCommand: return await ctx.send("You haven't used any commands yet!")
	//     await ctx.send(bot.EventListenerClass.LastUsedCommand)
}

func (UtilityCommands) history(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="User", value=str(User)).add_field(name="Is On Mobile", value=User.is_on_mobile()))
}

func (UtilityCommands) website(ctx *dgc.Ctx) {
	_ = utils.Process.OpenURL("https://lightning-bot.com/")
}

func (UtilityCommands) dashboard(ctx *dgc.Ctx) {
	_ = utils.Process.OpenURL("https://lightning-bot.com/dashboard")
}

func (UtilityCommands) restartbot(ctx *dgc.Ctx) {
	_ = utils.Process.Restart()
}

func (UtilityCommands) stopbot(ctx *dgc.Ctx) {
	os.Exit(0)
}

func (UtilityCommands) changelog(ctx *dgc.Ctx) {
	//     Release = requests.get("https://lightning-bot.com/api/v2/releases/getlatest?branch=main").json()
	//     Changelog = Release["changelog"].replace("&nbsp;", " ").replace("<p>", "").replace("</p>", "\n")[:1021]
	//     if len(Changelog) == 1021: Changelog += "..."
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Version", value=Release["version"]).add_field(name="Changelog", value=Changelog, inline=False))
}

func (UtilityCommands) uptime(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Uptime", value=display_time(int(time.time() - bot.StartTime))))
}

func (UtilityCommands) cls(ctx *dgc.Ctx) {
	//     bot.LightningBot.CLS()
	//     bot.LightningBot.PrintBanner()
}

func (UtilityCommands) mastercard(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.mastercardPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='MasterCard',icon_url='https://upload.wikimedia.org/wikipedia/commons/thumb/2/2a/Mastercard-logo.svg/1920px-Mastercard-logo.svg.png'))
}

func (UtilityCommands) visa16(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.visaPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Visa',icon_url='https://th.bing.com/th/id/R161dd5c504711b4644b769827cbd2aaf?rik=2BzQFCfHGrAJMQ&riu=http%3a%2f%2flogok.org%2fwp-content%2fuploads%2f2014%2f03%2fvisa-logo.png&ehk=DlfrlLCoE7o7DPrrnIiLPwVyCLovXEfIKvLBTD%2fVcFs%3d&risl=&pid=ImgRaw'))
}

func (UtilityCommands) visa13(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.visaPrefixList, 13, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Visa',icon_url='https://th.bing.com/th/id/R161dd5c504711b4644b769827cbd2aaf?rik=2BzQFCfHGrAJMQ&riu=http%3a%2f%2flogok.org%2fwp-content%2fuploads%2f2014%2f03%2fvisa-logo.png&ehk=DlfrlLCoE7o7DPrrnIiLPwVyCLovXEfIKvLBTD%2fVcFs%3d&risl=&pid=ImgRaw'))
}

func (UtilityCommands) amex(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.amexPrefixList, 15, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='American Express',icon_url='https://images.fastcompany.net/image/upload/w_596,c_limit,q_auto:best,f_auto,fl_lossy/wp-cms/uploads/sites/4/2018/04/4-you-might-not-notice-amex-new-brand.jpg'))
}

func (UtilityCommands) discover(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.discoverPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Discover',icon_url='https://th.bing.com/th/id/R875eb1db718b104f8b263a8241a79923?rik=GPJOPmauU249aw&riu=http%3a%2f%2fwww.hustlermoneyblog.com%2fwp-content%2fuploads%2f2016%2f10%2fDiscover-Logo.jpg&ehk=VURPeIHuOFuFNNRpdRW%2buoDTrsG46EwHoJ5QkjnopKs%3d&risl=&pid=ImgRaw'))
}

func (UtilityCommands) diners(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.discoverPrefixList, 14, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='Diners',icon_url='https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQedJZA8K4inF8fS2ufGUvk4GQTRW_ympFDJs4ksvVH4ANAFfCG_9gNN3P6wPkkL6Rm8gU&usqp=CAU'))
}

func (UtilityCommands) jcb(ctx *dgc.Ctx) {
	//     CC = bot.CreditCardGenerator.credit_card_number(bot.CreditCardGenerator.discoverPrefixList, 16, 1)[0]
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="Number", value=CC[0]).add_field(name="CCV", value=CC[1]).add_field(name="Expiry Date", value=CC[2])
	//         .set_author(name='JCB',icon_url='https://upload.wikimedia.org/wikipedia/commons/thumb/4/40/JCB_logo.svg/1200px-JCB_logo.svg.png'))
}

func (UtilityCommands) revavatarsearch(ctx *dgc.Ctx) {
	//     await bot.Helper.SendMessage(ctx, embed=bot.DefaultEmbed().add_field(name="URL", value=f"https://images.google.com/searchbyimage?image_url={User.avatar_url}"))
}

func (UtilityCommands) calc(ctx *dgc.Ctx) {
	if runtime.GOOS != "windows" {
		_ = ctx.RespondText("This command is only available on Windows")
	}
	_ = exec.Command("C:\\Windows\\System32\\calc.exe").Start()
}

func (UtilityCommands) notepad(ctx *dgc.Ctx) {
	if runtime.GOOS != "windows" {
		_ = ctx.RespondText("This command is only available on Windows")
	}
	_ = exec.Command("C:\\Windows\\System32\\notepad.exe").Start()
}

func (UtilityCommands) firstmessage(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = ctx.Session.State.Channel(ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	messages, err := ctx.Session.ChannelMessages(channel.ID, 1, "", "1", "")
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	if len(messages) != 1 {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	embed.NewEmbed().
		SetDescription("First Message").
		AddField("Author", messages[0].Author.String()).
		AddField("Message", messages[0].Content).
		AddField("Jump URL", fmt.Sprintf("https://discord.com/channels/%s/%s/%s", channel.GuildID, channel.ID, messages[0].ID)).
		Send(ctx)
}

func (UtilityCommands) hypesquad(ctx *dgc.Ctx) {
	house := utils.Arguments["<bravery/brilliance/balance>"]
	house = strings.ReplaceAll(house, "bravery", "1")
	house = strings.ReplaceAll(house, "brilliance", "2")
	house = strings.ReplaceAll(house, "balance", "3")
	fmt.Println(house)
	if house != "1" && house != "2" && house != "3" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<bravery/brilliance/balance>")
		return
	}
	_ = api.SetHypesquad(house)
}

func (UtilityCommands) leavehypesquad(ctx *dgc.Ctx) {
	_ = api.LeaveHypesquad()
}

func (UtilityCommands) theme(ctx *dgc.Ctx) {
	theme := strings.ToLower(utils.Arguments["<light/dark>"])
	if theme != "light" && theme != "dark" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<light/dark> must be either \"light\" or \"dark\"")
		return
	}
	_, _ = api.ChangeSettings(map[string]interface{}{"theme": theme}, true)
}

func (UtilityCommands) compact(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<on/off> must be either \"on\" or \"off\"")
		return
	}
	_, _ = api.ChangeSettings(map[string]interface{}{"message_display_compact": toggle == "on"}, true)
}

func (UtilityCommands) devmode(ctx *dgc.Ctx) {
	toggle := strings.ToLower(utils.Arguments["<on/off>"])
	if toggle != "on" && toggle != "off" {
		error_handling.OnError(ctx, error_handling.InvalidArgument, "<on/off> must be either \"on\" or \"off\"")
		return
	}
	_, _ = api.ChangeSettings(map[string]interface{}{"developer_mode": toggle == "on"}, true)
}

func (UtilityCommands) sessionid(ctx *dgc.Ctx) {
	embed.NewEmbed().
		AddField("Session ID", ctx.Session.State.SessionID).
		Send(ctx)
}
