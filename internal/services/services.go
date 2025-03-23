package services

import (
	"github.com/lb-selfbot/selfbot/internal/services/config"
	"github.com/lb-selfbot/selfbot/internal/services/router"
)

// Services is the container for all service instances.
type Services struct {
	Config config.Service

	Router router.Service
}

// InitServices initializes and wires all services.
func InitServices() *Services {
	cfg := config.NewService()
	router := router.NewService(cfg)

	return &Services{
		Config: cfg,
		Router: router,
	}
}
