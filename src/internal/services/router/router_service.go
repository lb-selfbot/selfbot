package router

import (
	"github.com/LightningDev1/discordgo"

	"github.com/LightningDev1/LB-Selfbot/internal/router"
	"github.com/LightningDev1/LB-Selfbot/internal/services/interfaces"
)

// RouterServiceImpl implements the RouterService interface
// It wraps the existing router implementation from internal/router.

// We alias the old router package as rt for delegation.

type RouterServiceImpl struct {
	rt router.Router
}

// NewRouterService creates a new RouterService instance.
// It takes an EmbedService for dependency and a RouterConfig for configuration.
func NewRouterService(embedSvc interfaces.EmbedService, cfg *router.RouterConfig) interfaces.RouterService {
	return &RouterServiceImpl{
		router.NewRouter(embedSvc, cfg),
	}
}

func (r *RouterServiceImpl) GetConfig() *router.RouterConfig {
	return r.rt.GetConfig()
}

func (r *RouterServiceImpl) RegisterHandler(session *discordgo.Session) {
	r.rt.RegisterHandler(session)
}

func (r *RouterServiceImpl) RegisterCategory(category *router.Category, commands []*router.Command) {
	r.rt.RegisterCategory(category, commands)
}

func (r *RouterServiceImpl) GetCategories() []*router.Category {
	return r.rt.GetCategories()
}

func (r *RouterServiceImpl) GetCategory(name string) *router.Category {
	return r.rt.GetCategory(name)
}

func (r *RouterServiceImpl) RegisterCmd(command *router.Command) {
	r.rt.RegisterCmd(command)
}

func (r *RouterServiceImpl) GetCmds() []*router.Command {
	return r.rt.GetCmds()
}

func (r *RouterServiceImpl) GetCmd(name string) *router.Command {
	return r.rt.GetCmd(name)
}

func (r *RouterServiceImpl) RegisterMiddleware(middleware router.Middleware) {
	r.rt.RegisterMiddleware(middleware)
}

func (r *RouterServiceImpl) InvokeRemoteCmd(session *discordgo.Session, channel *discordgo.Channel, author *discordgo.User, cmd *router.Command, args string) {
	r.rt.InvokeRemoteCmd(session, channel, author, cmd, args)
}
