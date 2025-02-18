package interfaces

type O11yService interface {
	Inject(uiSvc UIService)

	Init() error
}
