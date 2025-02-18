package interfaces

import "github.com/LightningDev1/LB-Selfbot/internal/embed"

type EmbedService interface {
	Inject(variablesSvc VariablesService)

	NewEmbed() *embed.Embed
}
