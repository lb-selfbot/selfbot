package interfaces

type RPCService interface {
	Inject(uiSvc UIService, discordSvc DiscordService, sesionsSvc SessionsService)

	Start()
}
