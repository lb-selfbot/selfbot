package commands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SessionsCommands struct {
	Sessions []*discordgo.DiscordSession
}

func (s *SessionsCommands) RegisterCommands(router *dgc.Router, mainSession *discordgo.Session) {
	router.StartCategory("Sessions", "Discord session commands")
	router.RegisterCmd(&dgc.Command{
		Name:        "viewsessions",
		Description: "View all your logged in sessions",
		Usage:       "[p]viewsessions",
		Handler:     s.viewsessions,
	})

	// The sessions are currently only sent in the READY event.
	// The sessions API is not enabled yet outside of experiments,
	// So we have to connect to the gateway every time to refresh sessions

	// UPDATE 17/09/2022: The sessions API is now enabled for everyone,
	// but the API is not as reliable as the READY event.

	go func() {
		for {
			time.Sleep(time.Second * 60)
			cfg, err := config.GetConfig()
			if err != nil {
				continue
			}
			session, err := discordgo.New(cfg.Account.Discord.Token)
			if err != nil {
				continue
			}
			session.StateEnabled = false
			session.ShouldSubscribeGuilds = false
			session.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
				time.Sleep(time.Second * 2)
				_ = session.Close()

				tempSessions := make([]*discordgo.DiscordSession, 0)

				for index, dcSession := range e.Sessions {
					// Remove the "all" session and the current session
					if dcSession.SessionID != "all" &&
						dcSession.SessionID != session.State.SessionID &&
						dcSession.SessionID != mainSession.State.SessionID {
						tempSessions = append(tempSessions, e.Sessions[index])
					}
				}

				s.HandleSessions(s.Sessions, tempSessions)
				s.Sessions = tempSessions
			})
			if err = session.Open(); err != nil {
				continue
			}
		}
	}()
}

func (s *SessionsCommands) HandleSessions(oldSessions, newSessions []*discordgo.DiscordSession) {

	// oldSessions will be empty at the beginning, we don't want to
	// detect new sessions if they were already connected before us
	if len(oldSessions) == 0 {
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	var oldSessionIDs []string
	for _, oldSession := range oldSessions {
		oldSessionIDs = append(oldSessionIDs, oldSession.SessionID)
	}

	var newSessionIDs []string
	for _, newSession := range newSessions {
		newSessionIDs = append(newSessionIDs, newSession.SessionID)
	}

	for _, oldSession := range oldSessions {
		if !utils.SliceContains(newSessionIDs, oldSession.SessionID) &&
			cfg.Events.SessionDisconnected {
			logging.UI.SessionDisconnected(oldSession.SessionID, oldSession.ClientInfo.OS, oldSession.ClientInfo.Client)
		}
	}

	for _, newSession := range newSessions {
		if !utils.SliceContains(oldSessionIDs, newSession.SessionID) &&
			cfg.Events.SessionConnected {
			logging.UI.SessionConnected(newSession.SessionID, newSession.ClientInfo.OS, newSession.ClientInfo.Client)
		}
	}
}

func (s *SessionsCommands) viewsessions(ctx *dgc.Ctx) {
	emb := embed.NewEmbed().SetDescription("Sessions")
	for index, session := range s.Sessions {
		emb.AddField(
			fmt.Sprintf("Session %d", index+1),
			fmt.Sprintf("Session ID: %s\nOS: %s\nClient: %s", session.SessionID, utils.String.Capitalize(session.ClientInfo.OS), utils.String.Capitalize(session.ClientInfo.Client)),
		)
	}
	emb.Send(ctx)
}
