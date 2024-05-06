package handler

// Store all handlers related

type Handler struct {
	HealthcheckHandler *HealthcheckHandler
}

func NewHandlers(healthcheckHandler *HealthcheckHandler) *Handler {
	return &Handler{
		HealthcheckHandler: healthcheckHandler,
	}
}
