package embed

import (
	"github.com/LightningDev1/LB-Selfbot/internal/embed"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
	"github.com/LightningDev1/LB-Selfbot/selfbot/config"
	"github.com/LightningDev1/discordgo"
)

type EmbedServiceImpl struct {
	variablesSvc interfaces.VariablesService
}

func (e *EmbedServiceImpl) Inject(variablesSvc interfaces.VariablesService) {
	e.variablesSvc = variablesSvc
}

func (e *EmbedServiceImpl) NewEmbed() *embed.Embed {
	embed := &embed.Embed{MessageEmbed: &discordgo.MessageEmbed{}}

	cfg, err := config.GetConfig()
	if err != nil {
		return embed
	}

	embed.SetTitle(e.variablesSvc.Parse(cfg.Selfbot.Title))
	embed.SetFooter(e.variablesSvc.Parse(cfg.Selfbot.Footer))

	return embed
}
