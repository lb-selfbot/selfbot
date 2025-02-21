package backup

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/internal/api"
	"github.com/LightningDev1/LB-Selfbot/internal/http"
	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/services/backup"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/discordgo"
)

type BackupCommands struct{ services *services.Services }

func NewBackupCommands(s *services.Services) *BackupCommands {
	return &BackupCommands{services: s}
}
func (c *BackupCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Backup",
		Description: "Create and restore backups of your account",
	}, []*router.Command{
		{
			Name:        "brelationships",
			Description: "Backup or restore your friends and blocked users",
			Arguments:   "<backup/restore> [backup name]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:    "backup",
					Handler: c.backupRelationships,
				},
				{
					Name:    "restore",
					Handler: c.restoreRelationships,
				},
			},
		},
		{
			Name:        "bservers",
			Description: "Backup or restore your servers",
			Arguments:   "<backup/restore> [backup name]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:    "backup",
					Handler: c.backupServers,
				},
				{
					Name:    "restore",
					Handler: c.restoreServers,
				},
			},
		},
		{
			Name:        "bprofile",
			Description: "Backup or restore your profile",
			Arguments:   "<backup/restore> [backup name]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:    "backup",
					Handler: c.backupProfile,
				},
				{
					Name:    "restore",
					Handler: c.restoreProfile,
				},
			},
		},
		{
			Name:        "bsettings",
			Description: "Backup or restore your settings",
			Arguments:   "<backup/restore> [backup name]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:    "backup",
					Handler: c.backupSettings,
				},
				{
					Name:    "restore",
					Handler: c.restoreSettings,
				},
			},
		},
		{
			Name:        "bgifs",
			Description: "Backup or restore your saved GIFs",
			Arguments:   "<backup/restore> [backup name]",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			SubCommands: []*router.Command{
				{
					Name:    "backup",
					Handler: c.backupGifs,
				},
				{
					Name:    "restore",
					Handler: c.restoreGifs,
				},
			},
		},
	})
}

func (c *BackupCommands) backupRelationships(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

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

	err = c.services.Backup.SaveBackup("relationships", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription(fmt.Sprintf("Backed up %d friends and %d blocked users.", len(friends), len(blocked))).
		Send(ctx)
}

func (c *BackupCommands) restoreRelationships(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	var data backup.Relationships
	err := c.services.Backup.LoadBackup("relationships", ctx.Arguments.Get(0).Raw(), &data)
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

	c.services.Embed.NewEmbed().
		SetDescription("Below are the relationships from your backup. Please add or block the users manually (at your own risk).").
		AddField("Friends", strings.Join(utils.Map(data.Friends, userToString), "\n")).
		AddField("Blocked", strings.Join(utils.Map(data.Blocked, userToString), "\n")).
		Send(ctx)
}

func (c *BackupCommands) backupServers(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
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

		channel, err := c.services.Backup.FindInvitableChannel(guild)
		if err != nil {
			addServer(guild, "")
			continue
		}

		invitableGuilds = append(invitableGuilds, invitableGuild{Guild: guild, Channel: channel})
	}

	m, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, progressFunc())
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
		_, _ = ctx.Session.ChannelMessageEdit(ctx.Message.ChannelID, m.ID, progressFunc())
		time.Sleep(time.Second * 5)
	}

	err = c.services.Backup.SaveBackup("servers", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
		return
	}
}

func (c *BackupCommands) restoreServers(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	var data backup.Servers
	err := c.services.Backup.LoadBackup("servers", ctx.Arguments.Get(0).Raw(), &data)
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

	c.services.Embed.NewEmbed().
		SetDescription("Below are the servers from your backup. Please join the servers manually (at your own risk).").
		AddField("Servers", strings.Join(utils.Map(data.Servers, serverToString), "\n")).
		Send(ctx)
}

func (c *BackupCommands) backupProfile(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

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

	err = c.services.Backup.SaveBackup("profile", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription("Successfully backed up your profile!").
		Send(ctx)
}

func (c *BackupCommands) restoreProfile(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	var data backup.Profile
	err := c.services.Backup.LoadBackup("profile", ctx.Arguments.Get(0).Raw(), &data)
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

	c.services.Embed.NewEmbed().
		SetTitle("Profile restored").
		SetDescription("Successfully restored your profile!").
		Send(ctx)
}

func (c *BackupCommands) backupSettings(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	data := backup.Settings{Settings: ctx.Session.State.UserSettings}

	err := c.services.Backup.SaveBackup("settings", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription("Successfully backed up your settings!").
		Send(ctx)
}

func (c *BackupCommands) restoreSettings(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	var data backup.Settings
	err := c.services.Backup.LoadBackup("settings", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
		return
	}

	err = ctx.Session.UserSettingsUpdate(data.Settings, false)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while restoring the settings: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("Settings restored").
		SetDescription("Successfully restored your settings!").
		Send(ctx)
}

func (c *BackupCommands) backupGifs(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	settings, err := ctx.Session.UserFrecencySettings()
	if err != nil {
		_ = ctx.RespondText("An error has occurred fetching settings: " + err.Error())
		return
	}

	data := backup.GIFs{Settings: settings}

	err = c.services.Backup.SaveBackup("gifs", ctx.Arguments.Get(0).Raw(), data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while saving the backup: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription("Successfully backed up your GIFs!").
		Send(ctx)
}

func (c *BackupCommands) restoreGifs(ctx *router.Ctx) {
	if err := c.services.Backup.CheckEligibility(); err != nil {
		_ = ctx.RespondText("An error has occurred: " + err.Error())
		return
	}

	var data backup.GIFs
	err := c.services.Backup.LoadBackup("gifs", ctx.Arguments.Get(0).Raw(), &data)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while restoring the backup: " + err.Error())
		return
	}

	err = ctx.Session.FrecencySettingsUpdate(data.Settings, false)
	if err != nil {
		_ = ctx.RespondText("An error has occurred while restoring the GIFs: " + err.Error())
		return
	}

	c.services.Embed.NewEmbed().
		SetTitle("Settings restored").
		SetDescription("Successfully restored your GIFs!").
		Send(ctx)
}
