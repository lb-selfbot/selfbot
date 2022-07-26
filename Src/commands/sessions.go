package commands

import (
	"fmt"
	"time"

	"github.com/LightningDev1/LB-Selfbot/config"
	"github.com/LightningDev1/LB-Selfbot/embed"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/dgc"
	"github.com/LightningDev1/discordgo"
)

type SessionsCommands struct {
	Sessions []*discordgo.DiscordSession
}

func (s *SessionsCommands) RegisterCommands(router *dgc.Router, mainSession *discordgo.Session) {
	router.RegisterCmd(&dgc.Command{
		Name:        "viewsessions",
		Description: "View all your logged in sessions",
		Usage:       "[p]viewsessions",
		Handler:     s.viewsessions,
	})

	// The sessions are currently only sent in the READY event.
	// The sessions API is not enabled yet outside of experiments,
	// So we have to connect to the gateway every time to refresh sessions

	go func() {
		for {
			time.Sleep(time.Second * 60)
			cfg, err := config.GetConfig()
			if err != nil {
				fmt.Println("Error getting config:", err)
			}
			session, err := discordgo.New(cfg.Account.Discord.Token)
			if err != nil {
				fmt.Println("Error creating session:", err)
			}
			session.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
				time.Sleep(time.Second * 2)
				session.Close()

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
				fmt.Println("Error opening session:", err)
			}
		}
	}()
}

func (s *SessionsCommands) HandleSessions(oldSessions, newSessions []*discordgo.DiscordSession) {
	if len(oldSessions) != len(newSessions) {
		fmt.Println("Sessions:", len(oldSessions), "->", len(newSessions))
	}

	// oldSessions will be empty at the beginning, we don't want to
	// detect new sessions if they were already connected before us
	if len(oldSessions) == 0 {
		return
	}

	oldSessionIDs := []string{}
	for _, oldSession := range oldSessions {
		oldSessionIDs = append(oldSessionIDs, oldSession.SessionID)
	}

	newSessionIDs := []string{}
	for _, newSession := range newSessions {
		newSessionIDs = append(newSessionIDs, newSession.SessionID)
	}

	for _, oldSession := range oldSessions {
		if !utils.SliceContains(newSessionIDs, oldSession.SessionID) {
			fmt.Println("Session disconnected:", oldSession.SessionID, " OS: ", oldSession.ClientInfo.OS, " Browser: ", oldSession.ClientInfo.Client)
		}
	}

	for _, newSession := range newSessions {
		if !utils.SliceContains(oldSessionIDs, newSession.SessionID) {
			fmt.Println("Session connected:", newSession.SessionID, " OS: ", newSession.ClientInfo.OS, " Browser: ", newSession.ClientInfo.Client)
		}
	}
}

func (s *SessionsCommands) viewsessions(ctx *dgc.Ctx) {
	embed := embed.NewEmbed().SetDescription("Sessions")
	for index, session := range s.Sessions {
		embed.AddField(
			fmt.Sprintf("Session %d", index+1),
			fmt.Sprintf("Session ID: %s\nOS: %s\nClient: %s", session.SessionID, utils.String.Capitalize(session.ClientInfo.OS), utils.String.Capitalize(session.ClientInfo.Client)),
		)
	}
	embed.Send(ctx)
}
