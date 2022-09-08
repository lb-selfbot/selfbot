package commands

import (
	"fmt"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DumpingCommands struct{}

func (d *DumpingCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Dumping", "Dumping commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpall",
		Description: "Dump all messages in a channel",
		Usage:       "[p]dumpall [channel]",
		Handler:     d.dumpall,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpimages",
		Description: "Dump all images in a channel",
		Usage:       "[p]dumpimages [channel]",
		Handler:     d.dumpimages,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpvideos",
		Description: "Dump all videos in a channel",
		Usage:       "[p]dumpvideos [channel]",
		Handler:     d.dumpvideos,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpaudio",
		Description: "Dump all audio files in a channel",
		Usage:       "[p]dumpaudio [channel]",
		Handler:     d.dumpaudio,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumptext",
		Description: "Dump all text files in a channel",
		Usage:       "[p]dumptext [channel]",
		Handler:     d.dumptext,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpemojis",
		Description: "Dump all emojis in a server",
		Usage:       "[p]dumpemojis [server id]",
		Handler:     d.dumpemojis,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpavatars",
		Description: "Dump all avatars from members in a server",
		Usage:       "[p]dumpavatars [server id]",
		Handler:     d.dumpavatars,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpchannels",
		Description: "Dump all channels in a server",
		Usage:       "[p]dumpchannels [server id]",
		Handler:     d.dumpchannels,
	})
	router.RegisterCmd(&dgc.Command{
		Name:        "dumpmembers",
		Description: "Dump all members in a server",
		Usage:       "[p]dumpmembers [server id]",
		Handler:     d.dumpmembers,
	})
}

func (*DumpingCommands) dumpall(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
	}
	if channel == nil {
		_ = ctx.RespondText("Error getting channel")
		return
	}
	server, _ := ctx.Session.State.Guild(channel.GuildID)
	serverName := "Error"
	if server != nil {
		serverName = server.Name
	}
	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	logString := fmt.Sprintf(
		"Channel Dump\n\nChannel: %s (ID: %s)\nServer: %s (ID: %s)\nType: Dump All\n\n",
		channel.Name,
		channel.ID,
		serverName,
		channel.GuildID,
	)
	// Iterate through messages in reverse order
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		logString += fmt.Sprintf(
			"[%s] Author: %s (ID: %s)\n%s\n\n",
			message.Timestamp.Local().Format("2006-01-02 15:04:05"),
			message.Author.String(),
			message.Author.ID,
			message.Content,
		)
	}
	folder := utils.File.GetMainFolder() + "/Dumping/All/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := folder + utils.File.SanitizeFilename(channel.Name+".txt")
	_ = utils.File.Write(filePath, logString)
	_ = utils.Process.OpenURL(filePath)
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d messages from %s", len(messages), channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpimages(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
	}
	if channel == nil {
		_ = ctx.RespondText("Error getting channel")
		return
	}
	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Images/" + utils.File.SanitizeFilename(channel.Name) + "/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.Process.OpenURL(folder)
	images := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if strings.HasPrefix(attachment.ContentType, "image/") {
				filePath := folder + utils.File.SanitizeFilename(attachment.Filename)
				fileBytes, _, _ := http.GetBytes(utils.Misc.Or(attachment.ProxyURL, attachment.URL))
				_ = utils.File.WriteBytes(filePath, fileBytes)
				images++
			}
		}
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d images from %s", images, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpvideos(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
	}
	if channel == nil {
		_ = ctx.RespondText("Error getting channel")
		return
	}
	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Videos/" + utils.File.SanitizeFilename(channel.Name) + "/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.Process.OpenURL(folder)
	videos := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if strings.HasPrefix(attachment.ContentType, "video/") {
				filePath := folder + utils.File.SanitizeFilename(attachment.Filename)
				fileBytes, _, _ := http.GetBytes(utils.Misc.Or(attachment.ProxyURL, attachment.URL))
				_ = utils.File.WriteBytes(filePath, fileBytes)
				videos++
			}
		}
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d videos from %s", videos, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpaudio(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
	}
	if channel == nil {
		_ = ctx.RespondText("Error getting channel")
		return
	}
	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Audio/" + utils.File.SanitizeFilename(channel.Name) + "/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.Process.OpenURL(folder)
	audios := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if strings.HasPrefix(attachment.ContentType, "audio/") {
				filePath := folder + utils.File.SanitizeFilename(attachment.Filename)
				fileBytes, _, _ := http.GetBytes(utils.Misc.Or(attachment.ProxyURL, attachment.URL))
				_ = utils.File.WriteBytes(filePath, fileBytes)
				audios++
			}
		}
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d audios from %s", audios, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumptext(ctx *dgc.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Event.ChannelID)
	}
	if channel == nil {
		_ = ctx.RespondText("Error getting channel")
		return
	}
	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		_ = ctx.RespondText("Error getting messages")
		return
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Text/" + utils.File.SanitizeFilename(channel.Name) + "/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.Process.OpenURL(folder)
	texts := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if strings.HasPrefix(attachment.ContentType, "text/") {
				filePath := folder + utils.File.SanitizeFilename(attachment.Filename)
				fileBytes, _, _ := http.GetBytes(utils.Misc.Or(attachment.ProxyURL, attachment.URL))
				_ = utils.File.WriteBytes(filePath, fileBytes)
				texts++
			}
		}
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d texts from %s", texts, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpemojis(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	emojis, err := ctx.Session.GuildEmojis(guild.ID)
	if err != nil {
		_ = ctx.RespondText("Error getting emojis: " + err.Error())
		return
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Emojis/" + utils.File.SanitizeFilename(guild.Name) + "/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.Process.OpenURL(folder)
	for _, emoji := range emojis {
		name := emoji.Name
		if emoji.Animated {
			name += ".gif"
		} else {
			name += ".png"
		}
		filePath := folder + utils.File.SanitizeFilename(name)
		fileBytes, _, _ := http.GetBytes(emoji.URL())
		_ = utils.File.WriteBytes(filePath, fileBytes)
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d emojis from %s", len(emojis), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpavatars(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Avatars/" + utils.File.SanitizeFilename(guild.Name) + "/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.Process.OpenURL(folder)
	for _, member := range guild.Members {
		name := member.User.String()
		if strings.HasPrefix(member.Avatar, "a_") {
			name += ".gif"
		} else {
			name += ".png"
		}
		filePath := folder + utils.File.SanitizeFilename(name)
		fileBytes, _, _ := http.GetBytes(member.AvatarURL(""))
		_ = utils.File.WriteBytes(filePath, fileBytes)
	}
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d avatars from %s", len(guild.Members), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpchannels(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	categories, categoryChannels := utils.Discord.MapCategories(guild.Channels)
	logString := ""
	for _, category := range categories {
		logString += fmt.Sprintf("\n\nCategory: %s\n\n", category.Name)
		for _, channel := range categoryChannels[category.ID] {
			logString += fmt.Sprintf("\t%s\n", channel.Name)
		}
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Channels/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := folder + utils.File.SanitizeFilename(guild.Name+".txt")
	_ = utils.File.Write(filePath, logString)
	_ = utils.Process.OpenURL(filePath)
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d channels from %s", len(guild.Channels), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpmembers(ctx *dgc.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Event.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	guildSubscription := discordgo.GuildSubscription{
		Session: ctx.Session,
		GuildID: guild.ID,
	}
	guildSubscription.Subscribe()
	time.Sleep(time.Second * 5)
	logString := fmt.Sprintf("Amount of dumped members: %d\n\n", len(guild.Members))
	for _, member := range guild.Members {
		logString += fmt.Sprintf("%s - %s\n", member.User.String(), member.User.ID)
	}
	folder := utils.File.GetMainFolder() + "/Dumping/Members/"
	folder = filepath.FromSlash(folder)
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := folder + utils.File.SanitizeFilename(guild.Name+".txt")
	_ = utils.File.Write(filePath, logString)
	_ = utils.Process.OpenURL(filePath)
	embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d members from %s", len(guild.Members), guild.Name)).
		Send(ctx)
}
