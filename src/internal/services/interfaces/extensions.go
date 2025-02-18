package interfaces

type ExtensionsService interface {
	Inject(uiSvc UIService, discordSvc DiscordService, embedSvc EmbedService)

	Register()
	GetLoadedExtensions() int
}
