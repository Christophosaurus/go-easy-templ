package handler

import (
	"github.com/go-easy-templ/internal/config"
	"github.com/go-easy-templ/internal/helper"
	"github.com/go-easy-templ/internal/service"
	"github.com/go-easy-templ/internal/types"
	"log/slog"
	"net/http"
	"time"
)

type HealthcheckHandler struct {
	logger   *slog.Logger
	config   *config.Config
	services *service.Services
}

func NewHealthcheckHandler(logger *slog.Logger, config *config.Config, services *service.Services) *HealthcheckHandler {
	return &HealthcheckHandler{
		logger:   logger,
		config:   config,
		services: services,
	}
}

func (h *HealthcheckHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()
	ctx := r.Context()

	response := types.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"version":     config.Version,
			"environment": h.config.Server.Env,
		},
	}

	err := h.services.DummyService.DoSomething(ctx)
	if err != nil {
		h.logger.Error("Something went wrong:", err)
	}

	if err := helper.Encode(w, http.StatusOK, nil, response); err != nil {
		h.logger.ErrorContext(ctx, "Unable to encode", slog.Any("error", err))
	}

	h.logger.DebugContext(ctx, "outgoing response", slog.Group("response_details",
		slog.Int("status_code", http.StatusOK),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.Any("response_body", response),
		slog.Duration("time_taken", time.Duration(time.Since(startTime).Milliseconds())),
	))

}
