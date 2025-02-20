package admin

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services"
	"github.com/LightningDev1/dgc"
)

var (
	verificationLevels = map[int]string{
		0: "None",
		1: "Low",
		2: "Medium",
		3: "High",
		4: "Highest",
	}

	boostLevels = map[int]string{
		0: "None",
		1: "Level 1",
		2: "Level 2",
		3: "Level 3",
	}
)

type AdminCommands struct {
	services *services.Services
}

func NewAdminCommands(services *services.Services) *AdminCommands {
	return &AdminCommands{services: services}
}

func (m *AdminCommands) Register() {
	router := m.services.Discord.GetRouter()
	router.StartCategory("Administration", "Administration commands")

	router.RegisterCmd(&dgc.Command{
		Name:        "channelinfo",
		Description: "Get info about a channel",
		Usage:       "[p]channelinfo <#channel>",
		Handler:     m.channelinfo,
	})

	router.RegisterCmd(&dgc.Command{
		Name:        "textchannel",
		Description: "Create a text channel",
		Usage:       "[p]textchannel <name>",
		Flags:       []string{CmdFlagDangerous},
		Handler:     m.textchannel,
	})

	// ... Keep registering all other admin commands the same way ...
}

// Keep all original command handlers unchanged
func (m *AdminCommands) channelinfo(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

func (m *AdminCommands) textchannel(ctx *dgc.Ctx) {
	// ... Existing implementation ...
}

// ... Keep all other existing command implementations as is ...
