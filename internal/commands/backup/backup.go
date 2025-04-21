package backup

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/lb-selfbot/discordgo"
	"github.com/lb-selfbot/selfbot/internal/fn"
	"github.com/lb-selfbot/selfbot/internal/http"

	"github.com/lb-selfbot/selfbot/internal/services"
	"github.com/lb-selfbot/selfbot/internal/services/backup"
	"github.com/lb-selfbot/selfbot/internal/services/router"
	"github.com/lb-selfbot/selfbot/internal/utils"
)

type BackupCommands struct{ services *services.Services }

func NewBackupCommands(s *services.Services) *BackupCommands {
	return &BackupCommands{services: s}
}

func (c *BackupCommands) Register() {
	c.services.Router.RegisterCategory(&router.Category{
		Name:        "Backup",
		Description: "Create and restore backups of your account",
	}, []*router.Command{
		{
			Name:  "brelationships",
			Flags: []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "backup",
					Description: "Back up your friends and blocked users",
					Arguments:   "[backup name]",
					Handler:     c.backupRelationships,
				},
				{
					Name:        "restore",
					Description: "Restore your friends and blocked users",
					Arguments:   "[backup name]",
					Handler:     c.restoreRelationships,
				},
			},
		},
		{
			Name:  "bservers",
			Flags: []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "backup",
					Description: "Back up your servers",
					Arguments:   "[backup name]",
					Handler:     c.backupServers,
				},
				{
					Name:        "restore",
					Description: "Restore your servers",
					Arguments:   "[backup name]",
					Handler:     c.restoreServers,
				},
			},
		},
		{
			Name:  "bprofile",
			Flags: []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "backup",
					Description: "Back up your profile",
					Arguments:   "[backup name]",
					Handler:     c.backupProfile,
				},
				{
					Name:        "restore",
					Description: "Restore your profile",
					Arguments:   "[backup name]",
					Handler:     c.restoreProfile,
				},
			},
		},
		{
			Name:  "bsettings",
			Flags: []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "backup",
					Description: "Back up your settings",
					Arguments:   "[backup name]",
					Handler:     c.backupSettings,
				},
				{
					Name:        "restore",
					Description: "Restore your settings",
					Arguments:   "[backup name]",
					Handler:     c.restoreSettings,
				},
			},
		},
		{
			Name:  "bgifs",
			Flags: []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:        "backup",
					Description: "Back up your saved GIFs",
					Arguments:   "[backup name]",
					Handler:     c.backupGifs,
				},
				{
					Name:        "restore",
					Description: "Restore your saved GIFs",
					Arguments:   "[backup name]",
					Handler:     c.restoreGifs,
				},
			},
		},
	})
}

func (c *BackupCommands) backupRelationships(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	rels, err := ctx.Session.UserRelationships()
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	if len(rels) == 0 {
		return fmt.Errorf("No relationships found.")
	}

	blocked := fn.Filter(rels, func(rel *discordgo.Relationship) bool {
		return rel.Type == 2 && rel.User != nil
	})

	friends := fn.Filter(rels, func(rel *discordgo.Relationship) bool {
		return rel.Type == 1 && rel.User != nil
	})

	exportRelationship := func(rel *discordgo.Relationship) backup.User {
		return backup.User{
			ID:       rel.User.ID,
			Username: rel.User.Username,
		}
	}

	data := backup.Relationships{
		Blocked: fn.Map(blocked, exportRelationship),
		Friends: fn.Map(friends, exportRelationship),
	}

	err = c.services.Backup.SaveBackup("relationships", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		return fmt.Errorf("An error has occurred while saving the backup: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Backup created").
		SetDescription("Backed up %d friends and %d blocked users.", len(friends), len(blocked)).
		Send()
	return nil
}

func (c *BackupCommands) restoreRelationships(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	var data backup.Relationships
	err := c.services.Backup.LoadBackup("relationships", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the backup: %v", err)
	}

	if len(data.Friends) == 0 && len(data.Blocked) == 0 {
		return fmt.Errorf("No relationships found in the backup.")
	}

	var currentFriends, currentBlocked []string

	rels, err := ctx.Session.UserRelationships()
	if err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	for _, rel := range rels {
		if rel.Type == 1 {
			currentFriends = append(currentFriends, rel.User.ID)
		} else if rel.Type == 2 {
			currentBlocked = append(currentBlocked, rel.User.ID)
		}
	}

	data.Friends = fn.Filter(data.Friends, func(user backup.User) bool {
		return !slices.Contains(currentFriends, user.ID)
	})

	data.Blocked = fn.Filter(data.Blocked, func(user backup.User) bool {
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

	c.services.Embed.New(ctx).
		SetDescription("Below are the relationships from your backup. Please add or block the users manually (at your own risk).").
		AddField("Friends", strings.Join(fn.Map(data.Friends, userToString), "\n")).
		AddField("Blocked", strings.Join(fn.Map(data.Blocked, userToString), "\n")).
		Send()
	return nil
}

func (c *BackupCommands) backupServers(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

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
		if guild.VanityURLCode != "" {
			addServer(guild, guild.VanityURLCode)
			continue
		}

		channel, err := c.services.Backup.FindInvitableChannel(ctx.Session, guild)
		if err != nil {
			addServer(guild, "")
			continue
		}

		invitableGuilds = append(invitableGuilds, invitableGuild{Guild: guild, Channel: channel})
	}

	m, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, progressFunc())
	if err != nil {
		return err
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
			ctx.RespondText("Too many errors occurred while backing up servers: %v", err)
			break
		}

		time.Sleep(time.Second * 5)
		_, _ = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, m.ID, progressFunc())
		time.Sleep(time.Second * 5)
	}

	err = c.services.Backup.SaveBackup("servers", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		return fmt.Errorf("An error has occurred while saving the backup: %v", err)
	}
	return nil
}

func (c *BackupCommands) restoreServers(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	var data backup.Servers
	err := c.services.Backup.LoadBackup("servers", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the backup: %v", err)
	}

	if len(data.Servers) == 0 {
		return fmt.Errorf("No servers found in the backup.")
	}

	currentServerIDs := fn.Map(ctx.Session.State.Guilds, func(guild *discordgo.Guild) string {
		return guild.ID
	})

	data.Servers = fn.Filter(data.Servers, func(server backup.Server) bool {
		return !slices.Contains(currentServerIDs, server.ID) && server.Invite != ""
	})

	serverToString := func(server backup.Server) string {
		return fmt.Sprintf("%s - %s - [click here (Discord App URL)](discord://-/invite/%s) - [click here (Website URL)](https://discord.gg/%s)", server.Name, server.ID, server.Invite, server.Invite)
	}

	c.services.Embed.New(ctx).
		SetDescription("Below are the servers from your backup. Please join the servers manually (at your own risk).").
		AddField("Servers", strings.Join(fn.Map(data.Servers, serverToString), "\n")).
		Send()
	return nil
}

func (c *BackupCommands) backupProfile(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	resp, err := http.DefaultClient().Get(ctx.Session.State.User.AvatarURL(""))
	if err != nil {
		return fmt.Errorf("An error has occurred while fetching your avatar: %v", err)
	}

	var avatar string
	if body := resp.MustGetBody(); body != nil {
		avatar = utils.Misc.ImageToBase64(body)
	}

	profile, err := ctx.Session.UserProfile(ctx.Session.State.User.ID)
	if err != nil {
		return fmt.Errorf("An error has occurred while fetching your profile: %v", err)
	}

	data := backup.Profile{
		User:    ctx.Session.State.User,
		Profile: profile.UserProfile,
		Avatar:  avatar,
	}

	err = c.services.Backup.SaveBackup("profile", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		return fmt.Errorf("An error has occurred while saving the backup: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Backup created").
		SetDescription("Successfully backed up your profile!").
		Send()
	return nil
}

func (c *BackupCommands) restoreProfile(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	var data backup.Profile
	err := c.services.Backup.LoadBackup("profile", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the backup: %v", err)
	}

	baseSettings := map[string]any{"avatar": data.Avatar}

	err = c.services.Discord.UpdateSettings("", baseSettings)
	if err != nil {
		return fmt.Errorf("Error restoring profile [1]: %v", err)
	}

	profile := map[string]any{}

	validKeys := []string{"accent_color", "bio", "pronouns", "theme_colors"}
	for _, key := range validKeys {
		if val, ok := data.Profile[key]; ok {
			profile[key] = val
		}
	}

	err = c.services.Discord.UpdateSettings("/profile", profile)
	if err != nil {
		return fmt.Errorf("Error restoring profile [2]: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Profile restored").
		SetDescription("Successfully restored your profile!").
		Send()
	return nil
}

func (c *BackupCommands) backupSettings(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	data := backup.Settings{Settings: ctx.Session.State.UserSettings}

	err := c.services.Backup.SaveBackup("settings", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		return fmt.Errorf("An error has occurred while saving the backup: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Backup created").
		SetDescription("Successfully backed up your settings!").
		Send()
	return nil
}

func (c *BackupCommands) restoreSettings(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	var data backup.Settings
	err := c.services.Backup.LoadBackup("settings", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the backup: %v", err)
	}

	err = ctx.Session.UserSettingsUpdate(data.Settings, false)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the settings: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Settings restored").
		SetDescription("Successfully restored your settings!").
		Send()
	return nil
}

func (c *BackupCommands) backupGifs(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	settings, err := ctx.Session.UserFrecencySettings()
	if err != nil {
		return fmt.Errorf("An error has occurred fetching settings: %v", err)
	}

	data := backup.GIFs{Settings: settings}

	err = c.services.Backup.SaveBackup("gifs", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		return fmt.Errorf("An error has occurred while saving the backup: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Backup created").
		SetDescription("Successfully backed up your GIFs!").
		Send()
	return nil
}

func (c *BackupCommands) restoreGifs(ctx *router.Ctx) error {
	if err := c.services.Backup.CheckEligibility(ctx.Session); err != nil {
		return fmt.Errorf("An error has occurred: %v", err)
	}

	var data backup.GIFs
	err := c.services.Backup.LoadBackup("gifs", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the backup: %v", err)
	}

	err = ctx.Session.FrecencySettingsUpdate(data.Settings, false)
	if err != nil {
		return fmt.Errorf("An error has occurred while restoring the GIFs: %v", err)
	}

	c.services.Embed.New(ctx).
		SetTitle("Settings restored").
		SetDescription("Successfully restored your GIFs!").
		Send()
	return nil
}
