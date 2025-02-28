package services2

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services2/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services2/router"
)

// Services is the container for all service instances.
type Services struct {
	Config config.ConfigService

	Router router.RouterService
}

// InitServices initializes and wires all services.
func InitServices() *Services {
	cfg := config.NewConfigService()
	router := router.NewRouterService(cfg)

	return &Services{
		Config: cfg,
		Router: router,
	}
}
