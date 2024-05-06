package server

import (
	"net/http"
)

func (srv *Server) routes() http.Handler {

	router := http.NewServeMux()

	router.HandleFunc("GET /api/v1/healthcheck", srv.handlers.HealthcheckHandler.Healthcheck)

	// Initialise middlewares
	middlewares := srv.addMiddlewares(
		srv.recovery(),
		srv.correlationID(),
		srv.requestInterceptor(),
	)

	return middlewares(router)
}
