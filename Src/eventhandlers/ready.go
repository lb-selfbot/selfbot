package eventhandlers

import (
	"fmt"
	"github.com/LightningDev1/discordgo"
)

func (e *Events) OnReady(_ *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Logged in as:", event.User.Username)
	fmt.Println(len(e.Router.Commands), "commands")
}
