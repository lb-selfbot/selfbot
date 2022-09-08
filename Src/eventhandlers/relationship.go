package eventhandlers

import (
	"time"

	"github.com/LightningDev1/LB-Selfbot/api"
	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnRelationshipAdd(_ *discordgo.Session, m *discordgo.RelationshipAdd) {
	// Incoming friend request
	if m.Type == 3 {
		e.Settings.friendRequests = append(e.Settings.friendRequests, &FriendRequest{Time: time.Now(), Relationship: m.Relationship})
		friendRequestsInLastMinute := 0
		for _, fr := range e.Settings.friendRequests {
			if time.Since(fr.Time) < time.Minute {
				friendRequestsInLastMinute++
			}
		}
		if friendRequestsInLastMinute > 10 {
			_ = e.Session.UserRelationshipsDelete(m.User.ID)
			_ = api.AllowFriendRequestsFrom([]string{})
			logging.UI.FriendRequestSpamBlocked()
			go func() {
				time.Sleep(5 * 60)
				_ = api.AllowFriendRequestsFrom([]string{"all"})
			}()
		}
	}
}
