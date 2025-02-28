package interfaces

import (
	"github.com/LightningDev1/discordgo"
)



type RouterService interface {
	GetConfig() *router.RouterConfig
	RegisterHandler()
	RegisterCategory(category *router.Category, commands []*router.Command)
	GetCategories() []*router.Category
	GetCategory(name string) *router.Category
	RegisterCmd(command *router.Command)
	GetCmds() []*router.Command
	GetCmd(name string) *router.Command
	RegisterMiddleware(middleware router.Middleware)
	InvokeRemoteCmd(session *discordgo.Session, channel *discordgo.Channel, author *discordgo.User, cmd *router.Command, args string)
}
