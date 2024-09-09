package commands

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/selfbot/backup"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type BackupCommands struct{}

func (b BackupCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Backup", "Create and restore backups of your account")

	router.RegisterCmd(&dgc.Command{
		Name:        "brelationships",
		Description: "Backup or restore your friends and blocked users",
		Usage:       "[p]brelationships <backup/restore> [backup name]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     b.brelationships,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "bservers",
		Description: "Backup or restore your servers",
		Usage:       "[p]bservers <backup/restore> [backup name]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     b.bservers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "bprofile",
		Description: "Backup or restore your profile",
		Usage:       "[p]bprofile <backup/restore> [backup name]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     b.bprofile,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "bsettings",
		Description: "Backup or restore your settings",
		Usage:       "[p]bsettings <backup/restore> [backup name]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     b.bsettings,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "bgifs",
		Description: "Backup or restore your saved GIF's",
		Usage:       "[p]bgifs <backup/restore> [backup name]",
		Flags:       []string{CmdFlagDangerous},
		Handler:     b.bgifs,
	})
}

func (b BackupCommands) brelationships(ctx *dgc.Ctx) {
	if err := backup.CheckEligible(ctx.Session); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	switch strings.ToLower(utils.Arguments["<backup/restore>"]) {
	case "backup":
		rels, err := ctx.Session.UserRelationships()
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}

		if len(rels) == 0 {
			_ = ctx.RespondText("No relationships found.")
			return
		}

		blocked := utils.Filter(rels, func(rel *discordgo.Relationship) bool {
			return rel.Type == 2 && rel.User != nil
		})

		friends := utils.Filter(rels, func(rel *discordgo.Relationship) bool {
			return rel.Type == 1 && rel.User != nil
		})

		exportRelationship := func(rel *discordgo.Relationship) backup.User {
			return backup.User{
				ID:       rel.User.ID,
				Username: rel.User.Username,
			}
		}

		data := backup.Relationships{
			Blocked: utils.Map(blocked, exportRelationship),
			Friends: utils.Map(friends, exportRelationship),
		}

		err = backup.SaveBackup("relationships", utils.Arguments["[backup name]"], data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Backup created").
			SetDescription(fmt.Sprintf("Backed up %d friends and %d blocked users.", len(friends), len(blocked))).
			Send(ctx)
	case "restore":
		var data backup.Relationships
		err := backup.RestoreBackup("relationships", utils.Arguments["[backup name]"], &data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
			return
		}

		if len(data.Friends) == 0 && len(data.Blocked) == 0 {
			_ = ctx.RespondText("No relationships found in the backup.")
			return
		}

		var currentFriends, currentBlocked []string

		rels, err := ctx.Session.UserRelationships()
		if err != nil {
			_ = ctx.RespondText("An error has occurred: " + err.Error())
			return
		}

		for _, rel := range rels {
			if rel.Type == 1 {
				currentFriends = append(currentFriends, rel.User.ID)
			} else if rel.Type == 2 {
				currentBlocked = append(currentBlocked, rel.User.ID)
			}
		}

		data.Friends = utils.Filter(data.Friends, func(user backup.User) bool {
			return !slices.Contains(currentFriends, user.ID)
		})

		data.Blocked = utils.Filter(data.Blocked, func(user backup.User) bool {
			return !slices.Contains(currentBlocked, user.ID)
		})

		userToString := func(user backup.User) string {
			return fmt.Sprintf(
				"%s - <@%s> - [click here (Discord App URL)](discord://-/users/%s) - [click here (Website URL)](https://discord.com/users/%s)",
				user.Username,
				user.ID,
				user.ID,
				user.ID,
			)
		}

		embed.NewEmbed().
			SetDescription("Below are the relationships from your backup. Please add or block the users manually (at your own risk).").
			AddField("Friends", strings.Join(utils.Map(data.Friends, userToString), "\n")).
			AddField("Blocked", strings.Join(utils.Map(data.Blocked, userToString), "\n")).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid action. Please use `backup` or `restore`.")
		return
	}
}

func (b BackupCommands) bservers(ctx *dgc.Ctx) {
	if err := backup.CheckEligible(ctx.Session); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	switch strings.ToLower(utils.Arguments["<backup/restore>"]) {
	case "backup":
		data := backup.Servers{Servers: []backup.Server{}}

		addServer := func(guild *discordgo.Guild, code string) {
			data.Servers = append(data.Servers, backup.Server{
				ID:     guild.ID,
				Name:   guild.Name,
				Invite: code,
			})
		}

		errors := 0

		progressFunc := func() string {
			done := len(data.Servers)
			total := len(ctx.Session.State.Guilds)
			percentage := (float64(done) / float64(total)) * 100

			return fmt.Sprintf("Backed up %d/%d servers (%.1f%%) (%d errors)", done, total, percentage, errors)
		}

		type invitableGuild struct {
			Guild   *discordgo.Guild
			Channel *discordgo.Channel
		}

		invitableGuilds := []invitableGuild{}

		for _, guild := range ctx.Session.State.Guilds {
			// Find servers with vanity URLs
			if guild.VanityURLCode != "" {
				addServer(guild, guild.VanityURLCode)
				continue
			}

			// Find servers with invitable channels
			channel := backup.FindInvitableChannel(ctx.Session, guild)
			if channel != nil {
				invitableGuilds = append(
					invitableGuilds,
					invitableGuild{Guild: guild, Channel: channel},
				)
				continue
			}

			// Servers not backed up
			addServer(guild, "")
		}

		m, err := ctx.Session.ChannelMessageSend(ctx.Event.ChannelID, progressFunc())
		if err != nil {
			return
		}

		time.Sleep(time.Second * 5)

		for _, guild := range invitableGuilds {
			invite, err := ctx.Session.ChannelInviteCreate(guild.Channel.ID, 0, 0, false)
			if err != nil {
				addServer(guild.Guild, "")
				errors++
			} else {
				addServer(guild.Guild, invite.Code)
			}

			if errors > 2 {
				_ = ctx.RespondText("Too many errors occurred while backing up servers: " + err.Error())
				break
			}

			time.Sleep(time.Second * 5)

			_, _ = ctx.Session.ChannelMessageEdit(ctx.Event.ChannelID, m.ID, progressFunc())

			time.Sleep(time.Second * 5)
		}

		err = backup.SaveBackup("servers", utils.Arguments["[backup name]"], data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
			return
		}
	case "restore":
		var data backup.Servers
		err := backup.RestoreBackup("servers", utils.Arguments["[backup name]"], &data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
			return
		}

		if len(data.Servers) == 0 {
			_ = ctx.RespondText("No servers found in the backup.")
			return
		}

		currentServerIDs := utils.Map(ctx.Session.State.Guilds, func(guild *discordgo.Guild) string {
			return guild.ID
		})

		data.Servers = utils.Filter(data.Servers, func(server backup.Server) bool {
			return !slices.Contains(currentServerIDs, server.ID) && server.Invite != ""
		})

		serverToString := func(server backup.Server) string {
			return fmt.Sprintf("%s - %s - [click here (Discord App URL)](discord://-/invite/%s) - [click here (Website URL)](https://discord.gg/%s)", server.Name, server.ID, server.Invite, server.Invite)
		}

		embed.NewEmbed().
			SetDescription("Below are the servers from your backup. Please join the servers manually (at your own risk).").
			AddField("Servers", strings.Join(utils.Map(data.Servers, serverToString), "\n")).
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid action. Please use `backup` or `restore`.")
		return
	}
}

func (b BackupCommands) bprofile(ctx *dgc.Ctx) {
	if err := backup.CheckEligible(ctx.Session); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	switch strings.ToLower(utils.Arguments["<backup/restore>"]) {
	case "backup":
		avatarResponse := http.Get(ctx.Session.State.User.AvatarURL(""))
		if avatarResponse.Error != nil {
			_ = ctx.RespondText("An error has occurred while fetching your avatar: " + avatarResponse.Error.Error())
			return
		}

		var avatar string
		if avatarResponse.Body != nil {
			avatar = "data:image/png;base64," + base64.StdEncoding.EncodeToString(avatarResponse.Body)
		}

		profile, err := ctx.Session.UserProfile(ctx.Session.State.User.ID)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while fetching your profile: " + err.Error())
			return
		}

		data := backup.Profile{
			User:    ctx.Session.State.User,
			Profile: profile.UserProfile,
			Avatar:  avatar,
		}

		err = backup.SaveBackup("profile", utils.Arguments["[backup name]"], data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Backup created").
			SetDescription("Successfully backed up your profile!").
			Send(ctx)
	case "restore":
		var data backup.Profile
		err := backup.RestoreBackup("profile", utils.Arguments["[backup name]"], &data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
			return
		}

		baseSettings := map[string]any{"avatar": data.Avatar}

		err = api.ChangeSettings(ctx.Session, "", baseSettings)
		if err != nil {
			_ = ctx.RespondText("Error restoring profile [1]: " + err.Error())
			return
		}

		profile := map[string]any{}

		validKeys := []string{"accent_color", "bio", "pronouns", "theme_colors"}
		for _, key := range validKeys {
			if val, ok := data.Profile[key]; ok {
				profile[key] = val
			}
		}

		err = api.ChangeSettings(ctx.Session, "/profile", profile)
		if err != nil {
			_ = ctx.RespondText("Error restoring profile [2]: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Profile restored").
			SetDescription("Successfully restored your profile!").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid action. Please use `backup` or `restore`.")
		return
	}
}

func (b BackupCommands) bsettings(ctx *dgc.Ctx) {
	if err := backup.CheckEligible(ctx.Session); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	switch strings.ToLower(utils.Arguments["<backup/restore>"]) {
	case "backup":
		data := backup.Settings{Settings: ctx.Session.State.UserSettings}

		err := backup.SaveBackup("settings", utils.Arguments["[backup name]"], data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Backup created").
			SetDescription("Successfully backed up your settings!").
			Send(ctx)
	case "restore":
		var data backup.Settings
		err := backup.RestoreBackup("settings", utils.Arguments["[backup name]"], &data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
			return
		}

		err = ctx.Session.UserSettingsUpdate(data.Settings, false)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the settings: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Settings restored").
			SetDescription("Successfully restored your settings!").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid action. Please use `backup` or `restore`.")
		return
	}
}

func (b BackupCommands) bgifs(ctx *dgc.Ctx) {
	if err := backup.CheckEligible(ctx.Session); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	switch strings.ToLower(utils.Arguments["<backup/restore>"]) {
	case "backup":
		settings, err := ctx.Session.UserFrecencySettings()
		if err != nil {
			_ = ctx.RespondText("An error has occurred fetching settings: " + err.Error())
			return
		}

		data := backup.GIFs{Settings: settings}

		err = backup.SaveBackup("gifs", utils.Arguments["[backup name]"], data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Backup created").
			SetDescription("Successfully backed up your GIFs!").
			Send(ctx)
	case "restore":
		var data backup.GIFs
		err := backup.RestoreBackup("gifs", utils.Arguments["[backup name]"], &data)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
			return
		}

		err = ctx.Session.FrecencySettingsUpdate(data.Settings, false)
		if err != nil {
			_ = ctx.RespondText("An error has occurred while restoring the GIFs: " + err.Error())
			return
		}

		embed.NewEmbed().
			SetTitle("Settings restored").
			SetDescription("Successfully restored your GIFs!").
			Send(ctx)
	default:
		_ = ctx.RespondText("Invalid action. Please use `backup` or `restore`.")
		return
	}
}
