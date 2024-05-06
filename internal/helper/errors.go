package helper

import (
	"github.com/go-easy-templ/internal/types"
	"log/slog"
	"net/http"
)

func logError(logger *slog.Logger, r *http.Request, err error) {

	logger.Error("server error", err, slog.Group("request_details",
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String())),
	)
}

func ErrorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, status int, message string) {
	env := types.Envelope{"error": message}

	err := Encode(w, status, nil, env)
	if err != nil {
		logError(logger, r, err)
	}
}

func ServerErrorResponse(logger *slog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	logError(logger, r, err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(logger, w, r, http.StatusInternalServerError, message)
}
