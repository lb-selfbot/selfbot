package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/eventhandlers"
	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/constants"
	"github.com/LightningDev1/LB-Selfbot/encodings"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
)

type TrollingCommands struct{}

func (t TrollingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Trolling", "Trolling commands")
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
		Usage:       "[p]trollnitro",
		Handler:     t.trollnitro,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "fakenitro",
		Description: "Generate fake nitro codes",
		Usage:       "[p]fakenitro <amount>",
		Handler:     t.fakenitro,
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
	router.RegisterCmd(&dgc.Command{
		Name:        "forcenick",
		Description: "Force a users nickname",
		Usage:       "[p]forcenick <user> <nick>",
		Handler:     t.forcenick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopforcenick",
		Description: "Stop forcing a users nickname",
		Usage:       "[p]stopforcenick <user>",
		Handler:     t.stopforcenick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "forcekick",
		Description: "Force kick a user",
		Usage:       "[p]forcekick <user>",
		Handler:     t.forcekick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopforcekick",
		Description: "Stop force kicking a user",
		Usage:       "[p]stopforcekick <user>",
		Handler:     t.stopforcekick,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "noleave",
		Description: "Stop user from leaving a group chat",
		Usage:       "[p]noleave <user>",
		Handler:     t.noleave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopnoleave",
		Description: "Stop no leave",
		Usage:       "[p]stopnoleave <user>",
		Handler:     t.stopnoleave,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "voicemove",
		Description: "Move to random voice channels",
		Usage:       "[p]voicemove <amount>",
		Handler:     t.voicemove,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "copy",
		Description: "Copy a users messages",
		Usage:       "[p]copy <user>",
		Handler:     t.copy,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopcopy",
		Description: "Stop copying a user",
		Usage:       "[p]stopcopy <user>",
		Handler:     t.stopcopy,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "mock",
		Description: "Mock a users messages",
		Usage:       "[p]mock <user>",
		Handler:     t.mock,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopmock",
		Description: "Stop mocking a user",
		Usage:       "[p]stopmock <user>",
		Handler:     t.stopmock,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "replycopy",
		Description: "Copy a users messages with reply",
		Usage:       "[p]replycopy <user>",
		Handler:     t.replycopy,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopreplycopy",
		Description: "Stop reply copying a user",
		Usage:       "[p]stopreplycopy <user>",
		Handler:     t.stopreplycopy,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "pinspam",
		Description: "Pin every message a user sends",
		Usage:       "[p]pinspam <user>",
		Handler:     t.pinspam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stoppinspam",
		Description: "Stop pin spam",
		Usage:       "[p]stoppinspam <user>",
		Handler:     t.stoppinspam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "unpinspam",
		Description: "Unpin messages pinned by pinspam",
		Usage:       "[p]unpinspam",
		Handler:     t.unpinspam,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "deletemessages",
		Description: "Delete every message a user sends",
		Usage:       "[p]deletemessages <user>",
		Handler:     t.deletemessages,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopdeletemessages",
		Description: "Stop deleting messages",
		Usage:       "[p]stopdeletemessages <user>",
		Handler:     t.stopdeletemessages,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "reactuser",
		Description: "React to every message a user sends",
		Usage:       "[p]reactuser <user> <emoji>",
		Handler:     t.reactuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopreactuser",
		Description: "Stop react user",
		Usage:       "[p]stopreactuser <user>",
		Handler:     t.stopreactuser,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "removereactions",
		Description: "Delete every reaction a user makes",
		Usage:       "[p]removereactions <user>",
		Handler:     t.removereactions,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopremovereactions",
		Description: "Stop removing reactions",
		Usage:       "[p]stopremovereactions <user>",
		Handler:     t.stopremovereactions,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "forcedisconnect",
		Description: "Force disconnect user from voice channels",
		Usage:       "[p]forcedisconnect <user>",
		Handler:     t.forcedisconnect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "stopforcedisconnect",
		Description: "Stop force disconnecting",
		Usage:       "[p]stopforcedisconnect <user>",
		Handler:     t.stopforcedisconnect,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "triggertyping",
		Description: "Fake \"User is typing\" message",
		Usage:       "[p]triggertyping <seconds> [channel]",
		Handler:     t.triggertyping,
	})
}

func (TrollingCommands) purgehack(ctx *dgc.Ctx) {
	_ = ctx.RespondText("ﾠ" + strings.Repeat("\n", 1998) + "ﾠ")
}

func (TrollingCommands) ghostping(*dgc.Ctx) {
	// noop
}

func (TrollingCommands) gprole(*dgc.Ctx) {
	// noop
}

func (TrollingCommands) token(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	_ = ctx.RespondText(encodings.Base64(user.ID, false))
}

func (TrollingCommands) trollnitro(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("dıscord.gift/%s %s https://media.discordapp.net/attachments/759817334121824297/763209007170453544/xedy9ugkqo621.png", utils.Random.String(16), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) fakenitro(ctx *dgc.Ctx) {
	amount, err := ctx.Arguments.Get(0).AsInt()
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}
	for i := 0; i < amount; i++ {
		_ = ctx.RespondText("https://discord.gift/" + utils.Random.String(16))
	}
}

func (TrollingCommands) hiddeninvite(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.INVISIBLE_TEXT, utils.Arguments["<invite>"]))
}

func (TrollingCommands) hiddenping(ctx *dgc.Ctx) {
	text := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)
	_ = ctx.RespondText(fmt.Sprintf("%s %s %s", text, constants.INVISIBLE_TEXT, utils.Arguments["<user>"]))
}

func (TrollingCommands) hiddeneveryone(ctx *dgc.Ctx) {
	_ = ctx.RespondText(fmt.Sprintf("%s %s @everyone", ctx.Arguments.Raw(), constants.INVISIBLE_TEXT))
}

func (TrollingCommands) forcenick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}
	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.ForceNick {
		if entry.GuildID == ctx.Event.GuildID && entry.UserID == member.User.ID {
			_ = ctx.RespondText("Users nickname is already being forced")
			return
		}
	}

	nick := utils.String.RemoveWords(ctx.Arguments.Raw(), 1)

	ctx.Session.GuildMemberNickname(ctx.Event.GuildID, member.User.ID, nick)

	eventhandlers.EventHandlers.Settings.ForceNick = append(eventhandlers.EventHandlers.Settings.ForceNick, eventhandlers.EventEntry{
		GuildID:  ctx.Event.GuildID,
		UserID:   member.User.ID,
		Argument: nick,
	})
}

func (TrollingCommands) stopforcenick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.ForceNick {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == member.User.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.ForceNick = newEntries
}

func (TrollingCommands) forcekick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	for _, entry := range eventhandlers.EventHandlers.Settings.ForceKick {
		if entry.GuildID == ctx.Event.GuildID && entry.UserID == member.User.ID {
			_ = ctx.RespondText("User is already being force kicked")
			return
		}
	}

	eventhandlers.EventHandlers.Settings.ForceKick = append(eventhandlers.EventHandlers.Settings.ForceKick, eventhandlers.EventEntry{
		GuildID: ctx.Event.GuildID,
		UserID:  member.User.ID,
	})

	_ = ctx.Session.GuildMemberDelete(ctx.Event.GuildID, member.User.ID)
}

func (TrollingCommands) stopforcekick(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.ForceKick {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == member.User.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.ForceKick = newEntries
}

func (TrollingCommands) noleave(ctx *dgc.Ctx) {
	channel, err := ctx.Session.State.Channel(ctx.Event.ChannelID)
	if err != nil {
		channel, err = ctx.Session.Channel(ctx.Event.ChannelID)
		if err != nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	if channel.Type != discordgo.ChannelTypeGroupDM {
		_ = ctx.RespondText("This command can only be used in a group DM")
		return
	}
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	isFriend := false
	for _, rel := range ctx.Session.State.Relationships {
		if rel.Type == 1 && rel.User.ID == user.ID {
			isFriend = true
			break
		}
	}
	if !isFriend {
		_ = ctx.RespondText("User must be a friend.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.NoLeave {
		if entry.GuildID == ctx.Event.GuildID && entry.UserID == user.ID {
			_ = ctx.RespondText("User is already being force kicked")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.NoLeave = append(eventhandlers.EventHandlers.Settings.NoLeave, eventhandlers.EventEntry{
		ChannelID: ctx.Event.ChannelID,
		UserID:    user.ID,
	})
}

func (TrollingCommands) stopnoleave(ctx *dgc.Ctx) {
	if !utils.Discord.CheckIfInServer(ctx) {
		return
	}

	member, success := utils.Discord.TryGetMemberFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.NoLeave {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == member.User.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.NoLeave = newEntries
}

func (TrollingCommands) voicemove(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, ctx.Event.GuildID)
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	amount, err := strconv.Atoi(utils.Arguments["<amount>"])
	if err != nil {
		_ = ctx.RespondText("Argument <amount> must be a number.")
		return
	}
	var voiceChannels []*discordgo.Channel
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			voiceChannels = append(voiceChannels, channel)
		}
	}
	if len(voiceChannels) == 0 {
		_ = ctx.RespondText("No voice channels found")
		return
	}
	var voiceConnection *discordgo.VoiceConnection
	for i := 0; i < amount; i++ {
		channel := utils.RandomChoice(voiceChannels)
		if voiceConnection != nil {
			if voiceConnection.ChannelID == channel.ID {
				i--
				continue
			}
		}
		voice, err := ctx.Session.ChannelVoiceJoin(ctx.Event.GuildID, channel.ID, false, false)
		if err == nil {
			time.Sleep(time.Second * 2)
			_ = voice.Disconnect()
			voiceConnection = voice
		}
	}
}

func getEventEntry(ctx *dgc.Ctx, userID, argument string) (e eventhandlers.EventEntry) {
	e.UserID = userID
	e.Argument = argument
	if ctx.Event.GuildID != "" {
		e.GuildID = ctx.Event.GuildID
	} else {
		e.ChannelID = ctx.Event.ChannelID
	}
	return
}

func (TrollingCommands) copy(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.Copy {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already copied")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.Copy = append(eventhandlers.EventHandlers.Settings.Copy, getEventEntry(ctx, user.ID, ""))
}

func (TrollingCommands) stopcopy(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.Copy {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.Copy = newEntries
}

func (TrollingCommands) mock(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.Mock {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already mocked")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.Mock = append(eventhandlers.EventHandlers.Settings.Mock, getEventEntry(ctx, user.ID, ""))
}

func (TrollingCommands) stopmock(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.Mock {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.Mock = newEntries
}

func (TrollingCommands) replycopy(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.ReplyCopy {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already reply copied")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.ReplyCopy = append(eventhandlers.EventHandlers.Settings.ReplyCopy, getEventEntry(ctx, user.ID, ""))
}

func (TrollingCommands) stopreplycopy(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.ReplyCopy {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.ReplyCopy = newEntries
}

func (TrollingCommands) pinspam(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.PinSpam {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already pin spammed")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.PinSpam = append(eventhandlers.EventHandlers.Settings.PinSpam, getEventEntry(ctx, user.ID, ""))
}

func (TrollingCommands) stoppinspam(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.PinSpam {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.PinSpam = newEntries
}

func (TrollingCommands) unpinspam(ctx *dgc.Ctx) {
	for _, message := range eventhandlers.EventHandlers.Settings.PinSpamMessages {
		_ = ctx.Session.ChannelMessageUnpin(message.ChannelID, message.ID)
	}
}

func (TrollingCommands) deletemessages(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.DeleteMessages {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already being deleted")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.DeleteMessages = append(eventhandlers.EventHandlers.Settings.DeleteMessages, getEventEntry(ctx, user.ID, ""))
}

func (TrollingCommands) stopdeletemessages(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.DeleteMessages {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.DeleteMessages = newEntries
}

func (TrollingCommands) reactuser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	for _, entry := range eventhandlers.EventHandlers.Settings.ReactUser {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already being reacted")
			return
		}
	}

	emoji := utils.Arguments["<emoji>"]
	emoji = strings.ReplaceAll(emoji, "<", "")
	emoji = strings.ReplaceAll(emoji, ">", "")

	eventhandlers.EventHandlers.Settings.ReactUser = append(eventhandlers.EventHandlers.Settings.ReactUser, getEventEntry(ctx, user.ID, emoji))
}

func (TrollingCommands) stopreactuser(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.ReactUser {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.ReactUser = newEntries
}

func (TrollingCommands) removereactions(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.RemoveReactions {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("Users reactions are already being removed")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.RemoveReactions = append(eventhandlers.EventHandlers.Settings.RemoveReactions, getEventEntry(ctx, user.ID, utils.Arguments["<emoji>"]))
}

func (TrollingCommands) stopremovereactions(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.RemoveReactions {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.RemoveReactions = newEntries
}

func (TrollingCommands) forcedisconnect(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}
	for _, entry := range eventhandlers.EventHandlers.Settings.ForceDisconnect {
		if entry.UserID == user.ID && (entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID) {
			_ = ctx.RespondText("User is already being force disconnected")
			return
		}
	}
	eventhandlers.EventHandlers.Settings.ForceDisconnect = append(eventhandlers.EventHandlers.Settings.ForceDisconnect, getEventEntry(ctx, user.ID, utils.Arguments["<emoji>"]))
}

func (TrollingCommands) stopforcedisconnect(ctx *dgc.Ctx) {
	user, success := utils.Discord.TryGetUserFromMention(ctx.Session, ctx.Arguments.Get(0))
	if !success {
		_ = ctx.RespondText("Could not find user.")
		return
	}

	newEntries := make([]eventhandlers.EventEntry, 0)

	for _, entry := range eventhandlers.EventHandlers.Settings.ForceDisconnect {
		isGuildOrChannel := entry.GuildID == ctx.Event.GuildID || entry.ChannelID == ctx.Event.ChannelID
		isMember := isGuildOrChannel && entry.UserID == user.ID

		if !isMember {
			newEntries = append(newEntries, entry)
		}
	}

	eventhandlers.EventHandlers.Settings.ForceDisconnect = newEntries
}

func (TrollingCommands) triggertyping(ctx *dgc.Ctx) {
	seconds, err := strconv.Atoi(utils.Arguments["<seconds>"])
	if err != nil {
		_ = ctx.RespondText("Argument <seconds> must be a number.")
		return
	}
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
		if channel == nil {
			_ = ctx.RespondText("Could not find channel")
			return
		}
	}
	utils.Discord.SleepWithTypingFix(ctx.Session, channel.ID, time.Duration(seconds)*time.Second)
}
