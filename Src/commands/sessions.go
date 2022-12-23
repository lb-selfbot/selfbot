package commands

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/sessions"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SessionsCommands struct{}

func (s *SessionsCommands) RegisterCommands(router *dgc.Router, mainSession *discordgo.Session) {
	router.StartCategory("Sessions", "Discord session commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "viewsessions",
		Description: "View all your logged in sessions",
		Usage:       "[p]viewsessions",
		Handler:     s.viewsessions,
	})
}

func (s *SessionsCommands) viewsessions(ctx *dgc.Ctx) {
	sessionsList, err := sessions.GetSessions(ctx.Session)
	if err != nil {
		ctx.RespondText("Failed to get sessions: " + err.Error())
		return
	}

	emb := embed.NewEmbed().SetDescription("Sessions")
	for index, session := range sessionsList {
		emb.AddField(
			fmt.Sprintf("Session %d", index+1),
			fmt.Sprintf(
				"Session ID: %s\nOS: %s\nClient: %s\nLast Used: %s",
				session.IDHash,
				utils.String.Capitalize(session.ClientInfo.OS),
				utils.String.Capitalize(session.ClientInfo.Platform),
				sessions.ParseLastUsedTime(session.LastUsedTime),
			),
		)
	}
	emb.Send(ctx)
}
