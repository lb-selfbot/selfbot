package services

import (
	"github.com/LightningDev1/LB-Selfbot/internal/services/config"
	"github.com/LightningDev1/LB-Selfbot/internal/services/router"
)

// Services is the container for all service instances.
type Services struct {
	Config config.Service

	Router router.RouterService
}

// InitServices initializes and wires all services.
func InitServices() *Services {
	cfg := config.NewService()
	router := router.NewRouterService(cfg)

	return &Services{
		Config: cfg,
		Router: router,
	}
}
