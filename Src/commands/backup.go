package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/http"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type BackupCommands struct{}

func (b BackupCommands) RegisterCommands(router *dgc.Router) {
	router.StartCategory("Backup", "Create and restore backups of your account")

	router.RegisterCmd(&dgc.Command{
		Name:        "backuprelationships",
		Description: "Create a backup of your friends and blocked users",
		Usage:       "[p]backuprelationships [backup name]",
		Handler:     b.backuprelationships,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "backupservers",
		Description: "Create a backup of your servers",
		Usage:       "[p]backupservers [backup name]",
		Handler:     b.backupservers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "backupaccount",
		Description: "Create a backup of your account",
		Usage:       "[p]backupaccount [backup name]",
		Handler:     b.backupaccount,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "backupsettings",
		Description: "Create a backup of your settings",
		Usage:       "[p]backupsettings [backup name]",
		Handler:     b.backupsettings,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "restorerelationships",
		Description: "Restore your relationships",
		Usage:       "[p]restorerelationships [backup name]",
		Handler:     b.restorerelationships,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "restoreservers",
		Description: "Restore your servers",
		Usage:       "[p]restoreservers [backup name]",
		Handler:     b.restoreservers,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "restoreaccount",
		Description: "Restore your account information",
		Usage:       "[p]restoreaccount [backup name]",
		Handler:     b.restoreaccount,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "restoresettings",
		Description: "Restore your settings",
		Usage:       "[p]restoresettings [backup name]",
		Handler:     b.restoresettings,
	})
}

func (b BackupCommands) isPhoneVerified(session *discordgo.Session) bool {
	user, err := session.User("@me")
	if err != nil {
		user = session.State.User
	}

	return user.Phone != ""
}

func (b BackupCommands) checkPhoneVerified(ctx *dgc.Ctx) bool {
	if !b.isPhoneVerified(ctx.Session) {
		ctx.RespondText("You must have a verified phone number to use this command.")
		return false
	}

	return true
}

func (b BackupCommands) saveBackup(backup any, filename string) error {
	filename = utils.IO.SanitizeFilename(filename)

	enc, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return err
	}

	filePath := path.Join(utils.IO.GetMainFolder(), "Backups", filename)
	filePath = filepath.FromSlash(filePath)

	err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	err = utils.IO.WriteBytes(filePath, enc)
	if err != nil {
		return err
	}

	return utils.System.OpenURL(filePath)
}

func (b BackupCommands) restoreBackup(filename string) (any, error) {
	filename = utils.IO.SanitizeFilename(filename)

	filePath := path.Join(utils.IO.GetMainFolder(), "Backups", filename)
	filePath = filepath.FromSlash(filePath)

	backup, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var backupData any

	err = json.Unmarshal(backup, &backupData)
	if err != nil {
		return nil, err
	}

	return backupData, nil
}

func (b BackupCommands) getSuffix() string {
	backupName := utils.Arguments["[backup name]"]

	if backupName != "" {
		return "-" + backupName
	}

	return ""
}

func (b BackupCommands) relationshipToUser(rel *discordgo.Relationship) map[string]any {
	user := map[string]any{
		"id":       rel.User.ID,
		"username": rel.User.Username,
	}

	if rel.User.Discriminator == "0" {
		user["display_name"] = rel.User.DisplayName
	} else {
		user["discriminator"] = rel.User.Discriminator
	}

	return user
}

func (b BackupCommands) backuprelationships(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	relationships, err := ctx.Session.UserRelationships()
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	blocked := utils.Filter(relationships, func(rel *discordgo.Relationship) bool {
		return rel.Type == 2 && rel.User != nil
	})

	friends := utils.Filter(relationships, func(rel *discordgo.Relationship) bool {
		return rel.Type == 1 && rel.User != nil
	})

	backup := map[string]any{
		"blocked": utils.Map(blocked, b.relationshipToUser),
		"friends": utils.Map(friends, b.relationshipToUser),
	}

	err = b.saveBackup(backup, "relationships"+b.getSuffix()+".json")
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription(fmt.Sprintf("Backed up %d friends and %d blocked users", len(friends), len(blocked))).
		Send(ctx)
}

func (b BackupCommands) findInvitableChannel(session *discordgo.Session, guild *discordgo.Guild) *discordgo.Channel {
	userID := session.State.User.ID

	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		perms, err := utils.Discord.UserChannelPermissions(session, guild.ID, channel.ID, userID)
		if err != nil {
			continue
		}

		if perms&discordgo.PermissionCreateInstantInvite == 0 || perms&discordgo.PermissionViewChannel == 0 {
			continue
		}

		return channel
	}

	return nil
}

func (b BackupCommands) backupservers(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	guilds := ctx.Session.State.Guilds

	backup := []any{}

	for _, guild := range guilds {
		var inviteCode string

		if guild.VanityURLCode != "" {
			inviteCode = guild.VanityURLCode
		} else {
			channel := b.findInvitableChannel(ctx.Session, guild)
			if channel == nil {
				ctx.RespondText(fmt.Sprintf("Error backup up %s (%s): No invitable channel found", guild.Name, guild.ID))
				continue
			}

			maxAge := 2592000

			if utils.SliceContains(guild.Features, discordgo.GuildFeatureCommunity) {
				maxAge = 0
			}

			invite, err := ctx.Session.ChannelInviteCreate(channel.ID, discordgo.Invite{
				MaxAge:    maxAge,
				MaxUses:   0,
				Temporary: false,
				Unique:    false,
			})

			if err != nil {
				ctx.RespondText(fmt.Sprintf("Error backup up %s (%s): Couldn't create invite: %s", guild.Name, guild.ID, err.Error()))
				continue
			}

			inviteCode = invite.Code
		}

		backup = append(backup, map[string]any{
			"name":   guild.Name,
			"id":     guild.ID,
			"invite": inviteCode,
		})
	}

	err := b.saveBackup(backup, "servers"+b.getSuffix()+".json")
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription(fmt.Sprintf("Backed up %d servers", len(backup))).
		Send(ctx)
}

func (b BackupCommands) backupaccount(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	err := b.saveBackup(ctx.Session.State.User, "account"+b.getSuffix()+".json")
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	httpResponse := http.Get(ctx.Session.State.User.AvatarURL(""))
	if httpResponse.Error != nil {
		ctx.RespondText("Error creating backup: " + httpResponse.Error.Error())
		return
	}

	avatarFileName := utils.IO.SanitizeFilename("avatar" + b.getSuffix() + ".png")
	avatarFilePath := path.Join(utils.IO.GetMainFolder(), "Backups", avatarFileName)

	err = utils.IO.WriteBytes(avatarFilePath, httpResponse.BodyBytes)
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription("Your account has been successfully backed up.").
		Send(ctx)
}

func (b BackupCommands) backupsettings(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	settings, err := ctx.Session.UserSettings()
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	err = b.saveBackup(settings, "settings"+b.getSuffix()+".json")
	if err != nil {
		ctx.RespondText("Error creating backup: " + err.Error())
		return
	}

	embed.NewEmbed().
		SetTitle("Backup created").
		SetDescription("Your settings have been successfully backed up.").
		Send(ctx)
}

func (b BackupCommands) restorerelationships(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	rawBackup, err := b.restoreBackup("relationships" + b.getSuffix() + ".json")
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}

	backup := rawBackup.(map[string]any)

	for _, user := range backup["blocked"].([]any) {
		user := user.(map[string]any)

		err := ctx.Session.UserRelationshipCreate(user["id"].(string), 2)
		if err != nil {
			ctx.RespondText(fmt.Sprintf("Error blocking user %s (continuing restore): %s", user["username"], err.Error()))
		}

		time.Sleep(time.Second * 3)
	}

	for _, user := range backup["friends"].([]any) {
		user := user.(map[string]any)

		var discriminator int

		if discrim, ok := user["discriminator"]; ok {
			discrim_str := discrim.(string)

			if discrim_str == "" {
				discriminator = 0
			} else {
				discriminator, _ = strconv.Atoi(discrim_str)
			}
		}

		err := ctx.Session.UserSendFriendRequest(user["username"].(string), discriminator)
		if err != nil {
			ctx.RespondText(fmt.Sprintf("Error adding user %s as friend (continuing restore): %s", user["username"], err.Error()))
		}

		time.Sleep(time.Second * 3)
	}

	embed.NewEmbed().
		SetTitle("Backup restored").
		SetDescription(fmt.Sprintf("Restored %d friends and %d blocked users", len(backup["friends"].([]any)), len(backup["blocked"].([]any)))).
		Send(ctx)
}

func (b BackupCommands) restoreservers(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	rawBackup, err := b.restoreBackup("servers" + b.getSuffix() + ".json")
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}

	backup := rawBackup.([]any)

	failed := []string{}
	success := 0

	for _, guild := range backup {
		guild := guild.(map[string]any)

		inviteCode := guild["invite"].(string)
		inviteLink := "https://discord.gg/" + inviteCode

		if len(failed) > 5 {
			failed = append(failed, fmt.Sprintf("%s - %s", guild["name"], inviteLink))
			continue
		}

		invite, err := ctx.Session.InviteComplex(inviteCode, "", inviteLink, true, true)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s - %s", guild["name"], inviteLink))
			continue
		}

		_, err = ctx.Session.InviteAccept(invite.Code, invite.Guild.ID, invite.Channel.ID, int(invite.Channel.Type))
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s - %s", guild["name"], inviteLink))
			continue
		}

		success++

		time.Sleep(time.Second * 3)
	}

	embed.NewEmbed().
		SetTitle("Backup restored").
		SetDescription(fmt.Sprintf("Restored %d servers", success)).
		AddField("Failed", strings.Join(failed, "\n")).
		Send(ctx)
}

func (b BackupCommands) restoreaccount(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	rawBackup, err := b.restoreBackup("account" + b.getSuffix() + ".json")
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}

	backupMap := rawBackup.(map[string]any)

	avatarFileName := utils.IO.SanitizeFilename("avatar" + b.getSuffix() + ".png")
	avatarPath := path.Join(utils.IO.GetMainFolder(), "Backups", avatarFileName)
	avatar, err := os.ReadFile(avatarPath)
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}
	avatarBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(avatar)

	backup := map[string]any{
		"avatar":       avatarBase64,
		"banner_color": backupMap["banner_color"],
		"accent_color": backupMap["accent_color"],
		"bio":          backupMap["bio"],
		"pronouns":     backupMap["pronouns"],
	}

	if globalName, ok := backupMap["global_name"]; ok && globalName != "" {
		backup["global_name"] = globalName
	}

	err = api.ChangeSettings(ctx.Session, "", backup)
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}
}

func (b BackupCommands) restoresettings(ctx *dgc.Ctx) {
	if !b.checkPhoneVerified(ctx) {
		return
	}

	backup, err := b.restoreBackup("settings" + b.getSuffix() + ".json")
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}

	err = api.ChangeSettings(ctx.Session, "/settings", backup.(map[string]any))
	if err != nil {
		ctx.RespondText("Error restoring backup: " + err.Error())
		return
	}
}
