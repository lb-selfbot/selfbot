package interfaces

type VariablesService interface {
	Inject(discordSvc DiscordService, extensionsSvc ExtensionsService, spotifySvc SpotifyService)

	Init()

	Parse(text string) string
}
