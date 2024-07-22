package commands

import (
	"fmt"
	"strings"

	"github.com/LightningDev1/LB-Selfbot/selfbot/backup"
	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
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

		exportRelationship := func(rel *discordgo.Relationship) map[string]any {
			return map[string]any{
				"id":       rel.User.ID,
				"username": rel.User.Username,
			}
		}

		data := map[string]any{
			"blocked": utils.Map(blocked, exportRelationship),
			"friends": utils.Map(friends, exportRelationship),
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

		friends := utils.Map(data.Friends, func(user backup.User) string {
			return fmt.Sprintf("%s - <@%s> - [click here](discord://-/users/%s)", user.Username, user.ID, user.ID)
		})

		blocked := utils.Map(data.Blocked, func(user backup.User) string {
			return fmt.Sprintf("%s - <@%s> - [click here](discord://-/users/%s)", user.Username, user.ID, user.ID)
		})

		embed.NewEmbed().
			SetDescription("Below are the relationships from your backup. Please add or block the users manually (at your own risk).").
			AddField("Friends", strings.Join(friends, "\n")).
			AddField("Blocked", strings.Join(blocked, "\n")).
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
	case "restore":
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
	case "restore":
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
	case "restore":
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
	case "restore":
	default:
		_ = ctx.RespondText("Invalid action. Please use `backup` or `restore`.")
		return
	}
}
