package dumping

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

type DumpingCommands struct{ services *services.Services }

func NewDumpingCommands(s *services.Services) *DumpingCommands {
	return &DumpingCommands{services: s}
}
func (c *DumpingCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Dumping",
		Description: "Dumping commands",
	}, []*router.Command{
		{
			Name:        "dumpall",
			Description: "Dump all messages in a channel",
			Arguments:   "dumpall [channel]",
			Handler:     d.dumpall,
		},
		{
			Name:        "dumpimages",
			Description: "Dump all images in a channel",
			Arguments:   "dumpimages [channel]",
			Handler:     d.dumpimages,
		},
		{
			Name:        "dumpvideos",
			Description: "Dump all videos in a channel",
			Arguments:   "dumpvideos [channel]",
			Handler:     d.dumpvideos,
		},
		{
			Name:        "dumpaudio",
			Description: "Dump all audio files in a channel",
			Arguments:   "dumpaudio [channel]",
			Handler:     d.dumpaudio,
		},
		{
			Name:        "dumptext",
			Description: "Dump all text files in a channel",
			Arguments:   "dumptext [channel]",
			Handler:     d.dumptext,
		},
		{
			Name:        "dumpattachments",
			Description: "Dump all attachments in a channel",
			Arguments:   "dumpattachments [channel]",
			Handler:     d.dumpattachments,
		},
		{
			Name:        "dumpfiletype",
			Description: "Dump certain files in a channel",
			Arguments:   "dumpfiletype <extension> [channel]",
			Handler:     d.dumpfiletype,
		},
		{
			Name:        "dumpcontenttype",
			Description: "Dump certain files in a channel",
			Arguments:   "dumpcontenttype <content type> [channel]",
			Handler:     d.dumpcontenttype,
		},
		{
			Name:        "dumpemojis",
			Description: "Dump all emojis in a server",
			Arguments:   "dumpemojis [server id]",
			Handler:     d.dumpemojis,
		},
		{
			Name:        "dumpstickers",
			Description: "Dump all stickers in a server",
			Arguments:   "dumpstickers [server id]",
			Handler:     d.dumpstickers,
		},
		{
			Name:        "dumpavatars",
			Description: "Dump all avatars from members in a server",
			Arguments:   "dumpavatars [server id]",
			Handler:     d.dumpavatars,
		},
		{
			Name:        "dumpchannels",
			Description: "Dump all channels in a server",
			Arguments:   "dumpchannels [server id]",
			Handler:     d.dumpchannels,
		},
		{
			Name:        "dumpmembers",
			Description: "Dump all members in a server",
			Arguments:   "dumpmembers [server id]",
			Handler:     d.dumpmembers,
		},
	})
}

func (*DumpingCommands) dumpall(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "All")
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, utils.IO.SanitizeFilename(channel.ID+"_"+channel.Name+".txt"))
	_ = utils.IO.Write(filePath, logString)
	_ = utils.System.OpenURL(filePath)
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d messages from %s", len(messages), channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpimages(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Images", utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	images := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if !strings.HasPrefix(attachment.ContentType, "image/") {
				continue
			}

			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			images++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d images from %s", images, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpvideos(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Videos", utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	videos := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if !strings.HasPrefix(attachment.ContentType, "video/") {
				continue
			}

			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			videos++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d videos from %s", videos, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpaudio(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Audio", utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	audios := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if !strings.HasPrefix(attachment.ContentType, "audio/") {
				continue
			}

			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			audios++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d audios from %s", audios, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumptext(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Text", utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	texts := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if !strings.HasPrefix(attachment.ContentType, "text/") {
				continue
			}

			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			texts++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d texts from %s", texts, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpattachments(ctx *router.Ctx) {
	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(0))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Attachments", utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	attachments := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			attachments++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d attachments from %s", attachments, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpfiletype(ctx *router.Ctx) {
	extension := strings.ReplaceAll(utils.Arguments["<extension>"], ".", "")

	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", utils.IO.SanitizeFilename("FileType-"+extension), utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	attachments := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if !strings.HasSuffix(attachment.Filename, extension) {
				continue
			}

			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			attachments++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d attachments from %s", attachments, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpcontenttype(ctx *router.Ctx) {
	contentType := strings.ToLower(utils.Arguments["<content type>"])

	channel, success := utils.Discord.TryGetChannelFromMention(ctx, ctx.Arguments.Get(1))
	if !success {
		channel, _ = utils.Discord.GetChannel(ctx.Session, ctx.Message.ChannelID)
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
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", utils.IO.SanitizeFilename("ContentType-"+contentType), utils.IO.SanitizeFilename(channel.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	attachments := 0
	for _, message := range messages {
		for _, attachment := range message.Attachments {
			if !strings.Contains(attachment.ContentType, contentType) {
				continue
			}

			fileName := utils.IO.SanitizeFilename(attachment.ID + "_" + attachment.Filename)
			filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

			httpResponse := http.Get(utils.Misc.Or(attachment.URL, attachment.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}

			_ = utils.IO.WriteBytes(filePath, httpResponse.Body)

			attachments++
		}
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d attachments from %s", attachments, channel.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpemojis(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	emojis, err := ctx.Session.GuildEmojis(guild.ID)
	if err != nil {
		_ = ctx.RespondText("Error getting emojis: " + err.Error())
		return
	}
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Emojis", utils.IO.SanitizeFilename(guild.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	for _, emoji := range emojis {
		name := emoji.ID + "_" + emoji.Name

		if emoji.Animated {
			name += ".gif"
		} else {
			name += ".png"
		}

		filePath := filepath.Join(folder, utils.IO.SanitizeFilename(name))

		httpResponse := http.Get(emoji.URL())
		if httpResponse.Error != nil {
			continue
		}

		_ = utils.IO.WriteBytes(filePath, httpResponse.Body)
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d emojis from %s", len(emojis), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpstickers(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Stickers", utils.IO.SanitizeFilename(guild.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	for _, sticker := range guild.Stickers {
		var format string

		switch sticker.FormatType {
		case discordgo.StickerFormatTypePNG:
			format = ".png"
		case discordgo.StickerFormatTypeAPNG:
			format = ".apng"
		case discordgo.StickerFormatTypeLottie:
			format = ".json"
		case discordgo.StickerFormatTypeGIF:
			format = ".gif"
		}

		fileName := fmt.Sprintf("%s_%s%s", sticker.ID, sticker.Name, format)

		filePath := filepath.Join(folder, utils.IO.SanitizeFilename(fileName))

		httpResponse := http.Get(fmt.Sprintf("https://media.discordapp.net/stickers/%s%s", sticker.ID, format))
		if httpResponse.Error != nil {
			continue
		}

		_ = utils.IO.WriteBytes(filePath, httpResponse.Body)
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d stickers from %s", len(guild.Stickers), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpavatars(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	ctx.Session.FetchGuildMembers(discordgo.NewFetchGuildMembersParams(guild.ID))

	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Avatars", utils.IO.SanitizeFilename(guild.Name))
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = utils.System.OpenURL(folder)
	for _, member := range guild.Members {
		name := member.User.ID + "_" + member.User.String()

		if strings.HasPrefix(member.Avatar, "a_") {
			name += ".gif"
		} else {
			name += ".png"
		}

		filePath := filepath.Join(folder, utils.IO.SanitizeFilename(name))

		httpResponse := http.Get(member.AvatarURL(""))
		if httpResponse.Error != nil {
			continue
		}

		_ = utils.IO.WriteBytes(filePath, httpResponse.Body)
	}
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d avatars from %s", len(guild.Members), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpchannels(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}
	categories, categoryChannels, looseChannels := utils.Discord.MapCategories(guild.Channels)
	logString := ""

	for _, channel := range looseChannels {
		logString += fmt.Sprintf("%s\n", channel.Name)
	}

	for _, category := range categories {
		logString += fmt.Sprintf("\n\nCategory: %s\n\n", category.Name)
		for _, channel := range categoryChannels[category.ID] {
			logString += fmt.Sprintf("\t%s\n", channel.Name)
		}
	}

	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Channels")
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, utils.IO.SanitizeFilename(guild.ID+"_"+guild.Name+".txt"))
	_ = utils.IO.Write(filePath, logString)
	_ = utils.System.OpenURL(filePath)
	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d channels from %s", len(guild.Channels), guild.Name)).
		Send(ctx)
}

func (*DumpingCommands) dumpmembers(ctx *router.Ctx) {
	guild, err := utils.Discord.GetGuild(ctx.Session, utils.Misc.Or(utils.Arguments["[server id]"], ctx.Message.GuildID))
	if err != nil {
		_ = ctx.RespondText("Error getting server: " + err.Error())
		return
	}

	x := discordgo.NewFetchGuildMembersParams(guild.ID)
	x.ForceScraping = true
	ctx.Session.FetchGuildMembers(x)

	logString := fmt.Sprintf("Amount of dumped members: %d\n\n", len(guild.Members))
	for _, member := range guild.Members {
		logString += fmt.Sprintf("%s - %s\n", member.User.String(), member.User.ID)
	}
	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Members")
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, utils.IO.SanitizeFilename(guild.ID+"_"+guild.Name+".txt"))
	_ = utils.IO.Write(filePath, logString)
	_ = utils.System.OpenURL(filePath)

	runtime.GC()

	c.services.Embed.NewEmbed().
		SetDescription(fmt.Sprintf("Dumped %d members from %s", len(guild.Members), guild.Name)).
		Send(ctx)
}
