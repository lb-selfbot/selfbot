package eventhandlers

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnRelationshipAdd(_ *discordgo.Session, m *discordgo.RelationshipAdd) {
	if m.Type != 3 {
		return
	}

	e.Settings.friendRequests = append(
		e.Settings.friendRequests,
		&FriendRequest{
			Time:         time.Now(),
			Relationship: m.Relationship,
		},
	)

	friendRequestsInLastMinute := utils.Filter(e.Settings.friendRequests, func(fr *FriendRequest) bool {
		return time.Since(fr.Time) < time.Minute
	})

	if len(friendRequestsInLastMinute) < 10 {
		return
	}

	_ = e.Session.UserRelationshipDelete(m.User.ID)

	_ = api.AllowFriendRequestsFrom(e.Session, []string{})

	logging.Logging.FriendRequestSpamBlocked()

	go func() {
		time.Sleep(5 * 60)
		_ = api.AllowFriendRequestsFrom(e.Session, []string{"all"})
	}()
}
