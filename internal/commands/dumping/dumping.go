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
	"github.com/LightningDev1/LB-Selfbot/internal/system"
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
			Arguments:   "[channel]",
			Handler:     c.dumpall,
		},
		{
			Name:        "dumpimages",
			Description: "Dump all images in a channel",
			Arguments:   "[channel]",
			Handler:     c.dumpimages,
		},
		{
			Name:        "dumpvideos",
			Description: "Dump all videos in a channel",
			Arguments:   "[channel]",
			Handler:     c.dumpvideos,
		},
		{
			Name:        "dumpaudio",
			Description: "Dump all audio files in a channel",
			Arguments:   "[channel]",
			Handler:     c.dumpaudio,
		},
		{
			Name:        "dumptext",
			Description: "Dump all text files in a channel",
			Arguments:   "[channel]",
			Handler:     c.dumptext,
		},
		{
			Name:        "dumpattachments",
			Description: "Dump all attachments in a channel",
			Arguments:   "[channel]",
			Handler:     c.dumpattachments,
		},
		{
			Name:        "dumpfiletype",
			Description: "Dump certain files in a channel",
			Arguments:   "<extension> [channel]",
			Handler:     c.dumpfiletype,
		},
		{
			Name:        "dumpcontenttype",
			Description: "Dump certain files in a channel",
			Arguments:   "<content type> [channel]",
			Handler:     c.dumpcontenttype,
		},
		{
			Name:        "dumpemojis",
			Description: "Dump all emojis in a server",
			Arguments:   "[server id]",
			Handler:     c.dumpemojis,
		},
		{
			Name:        "dumpstickers",
			Description: "Dump all stickers in a server",
			Arguments:   "[server id]",
			Handler:     c.dumpstickers,
		},
		{
			Name:        "dumpavatars",
			Description: "Dump all avatars from members in a server",
			Arguments:   "[server id]",
			Handler:     c.dumpavatars,
		},
		{
			Name:        "dumpchannels",
			Description: "Dump all channels in a server",
			Arguments:   "[server id]",
			Handler:     c.dumpchannels,
		},
		{
			Name:        "dumpmembers",
			Description: "Dump all members in a server",
			Arguments:   "[server id]",
			Handler:     c.dumpmembers,
		},
	})
}

// Helper: create folder relative to main folder, make dirs and open URL
func (c *DumpingCommands) createAndOpenFolder(pathElements ...string) string {
	parts := append([]string{utils.IO.GetMainFolder()}, pathElements...)
	folder := filepath.Join(parts...)
	_ = os.MkdirAll(folder, os.ModePerm)
	_ = system.OpenURL(folder)
	return folder
}

// Helper: Process attachments in messages with a filter function
func (c *DumpingCommands) processAttachments(messages []*discordgo.Message, folder string, filter func(att *discordgo.MessageAttachment) bool) int {
	count := 0
	for _, message := range messages {
		for _, att := range message.Attachments {
			if !filter(att) {
				continue
			}
			fileName := utils.IO.SanitizeFilename(att.ID + "_" + att.Filename)
			filePath := filepath.Join(folder, fileName)
			httpResponse := http.Get(utils.Misc.Or(att.URL, att.ProxyURL))
			if httpResponse.Error != nil {
				continue
			}
			_ = os.WriteFile(filePath, httpResponse.Body, os.ModePerm)
			count++
		}
	}
	return count
}

// Helper: Retrieve channel from context arguments
func (c *DumpingCommands) getChannel(ctx *router.Ctx, index int) (*discordgo.Channel, bool) {
	channel, err := ctx.Arguments.Get(index).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		ctx.RespondText("Error getting channel")
		return nil, false
	}
	return channel, true
}

// Helper: Retrieve guild from context arguments
func (c *DumpingCommands) getGuild(ctx *router.Ctx, index int) (*discordgo.Guild, bool) {
	guild, err := ctx.Arguments.Get(index).Guild()
	if err != nil {
		ctx.RespondText("Error getting server: %v", err)
		return nil, false
	}
	return guild, true
}

func (c *DumpingCommands) dumpall(ctx *router.Ctx) {
	channel, ok := c.getChannel(ctx, 0)
	if !ok {
		return
	}

	server, _ := ctx.Session.State.Guild(channel.GuildID)
	serverName := "Error"
	if server != nil {
		serverName = server.Name
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	logString := fmt.Sprintf("Channel Dump\n\nChannel: %s (ID: %s)\nServer: %s (ID: %s)\nType: Dump All\n\n",
		channel.Name,
		channel.ID,
		serverName,
		channel.GuildID,
	)
	// Iterate in reverse order
	for i := len(messages) - 1; i >= 0; i-- {
		message := messages[i]
		logString += fmt.Sprintf("[%s] Author: %s (ID: %s)\n%s\n\n",
			message.Timestamp.Local().Format("2006-01-02 15:04:05"),
			message.Author.String(),
			message.Author.ID,
			message.Content,
		)
	}

	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "All")
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, utils.IO.SanitizeFilename(channel.ID+"_"+channel.Name+".txt"))
	_ = os.WriteFile(filePath, []byte(logString), os.ModePerm)
	_ = system.OpenURL(filePath)

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d messages from %s", len(messages), channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpimages(ctx *router.Ctx) {
	channel, ok := c.getChannel(ctx, 0)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Images", utils.IO.SanitizeFilename(channel.Name))
	images := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "image/")
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d images from %s", images, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpvideos(ctx *router.Ctx) {
	channel, ok := c.getChannel(ctx, 0)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Videos", utils.IO.SanitizeFilename(channel.Name))
	videos := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "video/")
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d videos from %s", videos, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpaudio(ctx *router.Ctx) {
	channel, ok := c.getChannel(ctx, 0)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Audio", utils.IO.SanitizeFilename(channel.Name))
	audios := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "audio/")
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d audios from %s", audios, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumptext(ctx *router.Ctx) {
	channel, ok := c.getChannel(ctx, 0)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Text", utils.IO.SanitizeFilename(channel.Name))
	texts := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "text/")
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d texts from %s", texts, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpattachments(ctx *router.Ctx) {
	channel, ok := c.getChannel(ctx, 0)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Attachments", utils.IO.SanitizeFilename(channel.Name))
	attachments := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return true
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d attachments from %s", attachments, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpfiletype(ctx *router.Ctx) {
	extension := strings.ReplaceAll(ctx.Arguments.Get(0).Raw(), ".", "")
	channel, ok := c.getChannel(ctx, 1)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", utils.IO.SanitizeFilename("FileType-"+extension), utils.IO.SanitizeFilename(channel.Name))
	count := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasSuffix(att.Filename, extension)
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d attachments from %s", count, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpcontenttype(ctx *router.Ctx) {
	contentType := strings.ToLower(ctx.Arguments.Get(0).Raw())
	channel, ok := c.getChannel(ctx, 1)
	if !ok {
		return
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		ctx.RespondText("Error getting messages")
		return
	}

	folder := c.createAndOpenFolder("Dumping", utils.IO.SanitizeFilename("ContentType-"+contentType), utils.IO.SanitizeFilename(channel.Name))
	count := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.Contains(att.ContentType, contentType)
	})

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d attachments from %s", count, channel.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpemojis(ctx *router.Ctx) {
	guild, ok := c.getGuild(ctx, 0)
	if !ok {
		return
	}

	emojis, err := ctx.Session.GuildEmojis(guild.ID)
	if err != nil {
		ctx.RespondText("Error getting emojis: %v", err)
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Emojis", utils.IO.SanitizeFilename(guild.Name))
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
		_ = os.WriteFile(filePath, httpResponse.Body, os.ModePerm)
	}

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d emojis from %s", len(emojis), guild.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpstickers(ctx *router.Ctx) {
	guild, ok := c.getGuild(ctx, 0)
	if !ok {
		return
	}

	folder := c.createAndOpenFolder("Dumping", "Stickers", utils.IO.SanitizeFilename(guild.Name))
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
		_ = os.WriteFile(filePath, httpResponse.Body, os.ModePerm)
	}

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d stickers from %s", len(guild.Stickers), guild.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpavatars(ctx *router.Ctx) {
	guild, ok := c.getGuild(ctx, 0)
	if !ok {
		return
	}

	ctx.Session.FetchGuildMembers(discordgo.NewFetchGuildMembersParams(guild.ID))
	folder := c.createAndOpenFolder("Dumping", "Avatars", utils.IO.SanitizeFilename(guild.Name))
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
		_ = os.WriteFile(filePath, httpResponse.Body, os.ModePerm)
	}

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d avatars from %s", len(guild.Members), guild.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpchannels(ctx *router.Ctx) {
	guild, ok := c.getGuild(ctx, 0)
	if !ok {
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
	_ = os.WriteFile(filePath, []byte(logString), os.ModePerm)
	_ = system.OpenURL(filePath)

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d channels from %s", len(guild.Channels), guild.Name).
		Send(ctx)
}

func (c *DumpingCommands) dumpmembers(ctx *router.Ctx) {
	guild, ok := c.getGuild(ctx, 0)
	if !ok {
		return
	}

	x := discordgo.NewFetchGuildMembersParams(guild.ID)
	x.ForceScraping = true
	ctx.Session.FetchGuildMembers(x)

	logString := fmt.Sprintf("Amount of dumped members: %d\n\n", len(guild.Members))
	for _, member := range guild.Members {
		logString += fmt.Sprintf("%s - %s\n", member.User.String(), member.User.ID)
	}

	folder := c.createAndOpenFolder("Dumping", "Members")
	filePath := filepath.Join(folder, utils.IO.SanitizeFilename(guild.ID+"_"+guild.Name+".txt"))
	_ = os.WriteFile(filePath, []byte(logString), os.ModePerm)
	_ = system.OpenURL(filePath)

	runtime.GC()

	c.services.Embed.NewEmbed().
		SetDescription("Dumped %d members from %s", len(guild.Members), guild.Name).
		Send(ctx)
}
