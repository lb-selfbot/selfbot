package eventhandlers

import (
	"fmt"

	"github.com/LightningDev1/LB-Selfbot/logging"
	"github.com/LightningDev1/LB-Selfbot/utils"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnReady(ss *discordgo.Session, event *discordgo.Ready) {
	logging.Logging.OnReady(e.Router, event)
	fmt.Println(utils.System.GetIdleTime())
}
