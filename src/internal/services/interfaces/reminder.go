package interfaces

type ReminderService interface {
	Inject(uiSvc UIService)

	Start()
}
