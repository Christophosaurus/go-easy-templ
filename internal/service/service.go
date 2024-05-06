package service

// Store all service related

type Services struct {
	DummyService DummyService
}

func NewServices(dummyService DummyService) *Services {
	return &Services{
		DummyService: dummyService,
	}
}
