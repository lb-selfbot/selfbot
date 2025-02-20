package sessions

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/LB-Selfbot/internal/utils"
	"github.com/LightningDev1/LB-Selfbot/selfbot/sessions"
)

type SessionsCommands struct{ services *services.Services }

func NewSessionsCommands(s *services.Services) *SessionsCommands {
	return &SessionsCommands{services: s}
}
func (c *SessionsCommands) Register() {
	cmdRouter := c.services.Discord.GetRouter()

	cmdRouter.RegisterCategory(&router.Category{
		Name:        "Sessions",
		Description: "Discord session commands",
	}, []*router.Command{
		{
			Name:        "viewgatewaysessions",
			Description: "View all your gateway sessions (active devices)",
			Arguments:   "viewgatewaysessions",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     s.viewgatewaysessions,
		},
		{
			Name:        "viewloginsessions",
			Description: "View all your login sessions (tokens)",
			Arguments:   "viewloginsessions",
			Flags:       []router.CommandFlag{router.CmdFlagDangerous},
			Handler:     s.viewloginsessions,
		},
	})
}

func (s *SessionsCommands) viewgatewaysessions(ctx *router.Ctx) {
	emb := c.services.Embed.NewEmbed().SetDescription("Sessions")
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

func (s *SessionsCommands) viewloginsessions(ctx *router.Ctx) {
	emb := c.services.Embed.NewEmbed().SetDescription("Sessions")
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
