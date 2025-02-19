package interfaces

type VariablesService interface {
	Inject(discordSvc DiscordService, eventReactionsSvc EventReactionsService, extensionsSvc ExtensionsService, spotifySvc SpotifyService)

	Init()

	Parse(text string) string
}
