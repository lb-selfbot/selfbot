package dumping

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/http"
	"github.com/lb-selfbot/selfbot/internal/services/router"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/system"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type DumpingCommands struct{ services *services.Services }

func NewDumpingCommands(s *services.Services) *DumpingCommands {
	return &DumpingCommands{services: s}
}

func (c *DumpingCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
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
			resp, err := http.DefaultClient().Get(utils.Misc.Or(att.URL, att.ProxyURL))
			if err != nil {
				continue
			}
			_ = os.WriteFile(filePath, resp.MustGetBody(), os.ModePerm)
			count++
		}
	}
	return count
}

func (c *DumpingCommands) dumpall(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	server, _ := ctx.Session.State.Guild(channel.GuildID)
	serverName := "Error"
	if server != nil {
		serverName = server.Name
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
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

	c.services.Embed.Respond(ctx, "Dumped %d messages from %s", len(messages), channel.Name)
	return nil
}

func (c *DumpingCommands) dumpimages(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", "Images", utils.IO.SanitizeFilename(channel.Name))
	images := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "image/")
	})

	c.services.Embed.Respond(ctx, "Dumped %d images from %s", images, channel.Name)
	return nil
}

func (c *DumpingCommands) dumpvideos(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", "Videos", utils.IO.SanitizeFilename(channel.Name))
	videos := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "video/")
	})

	c.services.Embed.Respond(ctx, "Dumped %d videos from %s", videos, channel.Name)
	return nil
}

func (c *DumpingCommands) dumpaudio(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", "Audio", utils.IO.SanitizeFilename(channel.Name))
	audios := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "audio/")
	})

	c.services.Embed.Respond(ctx, "Dumped %d audios from %s", audios, channel.Name)
	return nil
}

func (c *DumpingCommands) dumptext(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", "Text", utils.IO.SanitizeFilename(channel.Name))
	texts := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasPrefix(att.ContentType, "text/")
	})

	c.services.Embed.Respond(ctx, "Dumped %d texts from %s", texts, channel.Name)
	return nil
}

func (c *DumpingCommands) dumpattachments(ctx *router.Ctx) error {
	channel, err := ctx.Arguments.Get(0).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", "Attachments", utils.IO.SanitizeFilename(channel.Name))
	attachments := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return true
	})

	c.services.Embed.Respond(ctx, "Dumped %d attachments from %s", attachments, channel.Name)
	return nil
}

func (c *DumpingCommands) dumpfiletype(ctx *router.Ctx) error {
	extension := strings.ReplaceAll(ctx.Arguments.Get(0).Raw(), ".", "")

	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", utils.IO.SanitizeFilename("FileType-"+extension), utils.IO.SanitizeFilename(channel.Name))
	count := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.HasSuffix(att.Filename, extension)
	})

	c.services.Embed.Respond(ctx, "Dumped %d attachments from %s", count, channel.Name)
	return nil
}

func (c *DumpingCommands) dumpcontenttype(ctx *router.Ctx) error {
	contentType := strings.ToLower(ctx.Arguments.Get(0).Raw())

	channel, err := ctx.Arguments.Get(1).Channel(router.ArgumentFlagFallbackDefault)
	if err != nil {
		return fmt.Errorf("could not get channel: %v", err)
	}

	messages, err := utils.Discord.GetChannelMessages(ctx.Session, channel.ID, 1000)
	if err != nil {
		return fmt.Errorf("Error getting messages")
	}

	folder := c.createAndOpenFolder("Dumping", utils.IO.SanitizeFilename("ContentType-"+contentType), utils.IO.SanitizeFilename(channel.Name))
	count := c.processAttachments(messages, folder, func(att *discordgo.MessageAttachment) bool {
		return strings.Contains(att.ContentType, contentType)
	})

	c.services.Embed.Respond(ctx, "Dumped %d attachments from %s", count, channel.Name)
	return nil
}

func (c *DumpingCommands) dumpemojis(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("could not get server: %v", err)
	}

	emojis, err := ctx.Session.GuildEmojis(guild.ID)
	if err != nil {
		return fmt.Errorf("Error getting emojis: %v", err)
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
		resp, err := http.DefaultClient().Get(emoji.URL())
		if err != nil {
			continue
		}
		_ = os.WriteFile(filePath, resp.MustGetBody(), os.ModePerm)
	}

	c.services.Embed.Respond(ctx, "Dumped %d emojis from %s", len(emojis), guild.Name)
	return nil
}

func (c *DumpingCommands) dumpstickers(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("could not get server: %v", err)
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
		resp, err := http.DefaultClient().Get(fmt.Sprintf("https://media.discordapp.net/stickers/%s%s", sticker.ID, format))
		if err != nil {
			continue
		}
		_ = os.WriteFile(filePath, resp.MustGetBody(), os.ModePerm)
	}

	c.services.Embed.Respond(ctx, "Dumped %d stickers from %s", len(guild.Stickers), guild.Name)
	return nil
}

func (c *DumpingCommands) dumpavatars(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("could not get server: %v", err)
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
		resp, err := http.DefaultClient().Get(member.AvatarURL(""))
		if err != nil {
			continue
		}
		_ = os.WriteFile(filePath, resp.MustGetBody(), os.ModePerm)
	}

	c.services.Embed.Respond(ctx, "Dumped %d avatars from %s", len(guild.Members), guild.Name)
	return nil
}

func (c *DumpingCommands) dumpchannels(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("could not get server: %v", err)
	}

	categoryData := discordgo.MapCategories(guild.Channels)
	formattedData := strings.Join(fn.Map(categoryData, func(data *discordgo.CategoryData) string {
		channels := fn.Map(data.Channels, func(channel *discordgo.Channel) string {
			return fmt.Sprintf("  %s (ID: %s)", channel.Name, channel.ID)
		})
		if len(channels) == 0 {
			channels = append(channels, "No channels")
		}

		return fmt.Sprintf("Category: %s\n%s", data.GetCategoryName(), strings.Join(channels, "\n"))
	}), "\n\n")

	folder := filepath.Join(utils.IO.GetMainFolder(), "Dumping", "Channels")
	_ = os.MkdirAll(folder, os.ModePerm)
	filePath := filepath.Join(folder, utils.IO.SanitizeFilename(guild.ID+"_"+guild.Name+".txt"))
	_ = os.WriteFile(filePath, []byte(formattedData), os.ModePerm)
	_ = system.OpenURL(filePath)

	c.services.Embed.Respond(ctx, "Dumped %d channels from %s", len(guild.Channels), guild.Name)
	return nil
}

func (c *DumpingCommands) dumpmembers(ctx *router.Ctx) error {
	guild, err := ctx.Arguments.Get(0).Guild()
	if err != nil {
		return fmt.Errorf("could not get server: %v", err)
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

	c.services.Embed.Respond(ctx, "Dumped %d members from %s", len(guild.Members), guild.Name)
	return nil
}
