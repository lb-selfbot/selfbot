package commands

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SessionsCommands struct{}

func (s *SessionsCommands) RegisterCommands(router *dgc.Router, mainSession *discordgo.Session) {
	router.StartCategory("Sessions", "Discord session commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "viewgatewaysessions",
		Description: "View all your gateway sessions (active devices)",
		Usage:       "[p]viewgatewaysessions",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.viewgatewaysessions,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "viewloginsessions",
		Description: "View all your login sessions (tokens)",
		Usage:       "[p]viewloginsessions",
		Flags:       []string{"DANGEROUS"},
		Handler:     s.viewloginsessions,
	})
}

func (s *SessionsCommands) viewgatewaysessions(ctx *dgc.Ctx) {
	emb := embed.NewEmbed().SetDescription("Sessions")
	for index, session := range sessions.GatewaySessions {
		var activity string
		if len(session.Activities) > 0 {
			activity = session.Activities[0].Name
		} else {
			activity = "None"
		}

		emb.AddField(
			fmt.Sprintf("Session %d", index+1),
			fmt.Sprintf(
				"Session ID: %s\nOS: %s\nClient: %s\nActivity: %s",
				session.SessionID,
				utils.String.Capitalize(session.ClientInfo.OS),
				utils.String.Capitalize(session.ClientInfo.Client),
				activity,
			),
		)
	}
	emb.Send(ctx)
}

func (s *SessionsCommands) viewloginsessions(ctx *dgc.Ctx) {
	emb := embed.NewEmbed().SetDescription("Sessions")
	for index, session := range sessions.LoginSessions {
		emb.AddField(
			fmt.Sprintf("Session %d", index+1),
			fmt.Sprintf(
				"Session ID: %s\nOS: %s\nPlatform: %s\nLast Used: %s",
				session.IDHash,
				utils.String.Capitalize(session.ClientInfo.OS),
				utils.String.Capitalize(session.ClientInfo.Platform),
				session.LastUsedTime.Format("2006-01-02 15:04:05"),
			),
		)
	}
	emb.Send(ctx)
}
