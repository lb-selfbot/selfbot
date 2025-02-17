package interfaces

type ExtensionsService interface {
	Inject(uiSvc UIService, discordSvc DiscordService)

	Register()
	GetLoadedExtensions() int
}
