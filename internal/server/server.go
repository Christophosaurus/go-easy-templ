package server

import (
	"context"
	"errors"
	"github.com/go-easy-templ/internal/config"
	"github.com/go-easy-templ/internal/handler"
	"github.com/go-easy-templ/internal/service"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	wg       sync.WaitGroup
	config   *config.Config
	logger   *slog.Logger
	handlers *handler.Handler
	services *service.Services
}

func NewServer(config *config.Config, logger *slog.Logger, handlers *handler.Handler) *Server {
	return &Server{
		config:   config,
		logger:   logger,
		handlers: handlers,
	}
}

func (srv *Server) Run() error {

	ctx := context.Background()

	httpServer := &http.Server{
		Addr:         net.JoinHostPort(srv.config.Server.Hostname, srv.config.Server.Port),
		Handler:      srv.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Graceful shutdown
	shutdownError := make(chan error)
	go srv.gracefulShutdown(httpServer, shutdownError)

	//srv.logger.Debug("starting server", "server_details", map[string]interface{}{
	//	"addr": httpServer.Addr,
	//})

	srv.logger.DebugContext(ctx, "starting server", slog.Group("server_details",
		slog.String("addr", httpServer.Addr),
	))

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait to receive the return value from Shutdown()
	if err := <-shutdownError; err != nil {
		return err
	}

	//srv.logger.Info("server stopped", "server_details", map[string]interface{}{
	//	"addr": httpServer.Addr,
	//})

	srv.logger.Debug("starting stopped", slog.Group("server_details",
		slog.String("addr", httpServer.Addr),
	))

	return nil
}

func (srv *Server) gracefulShutdown(httpServer *http.Server, shutdownError chan error) {

	// Listen for catchable stop signals with a buffer
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s := <-quit

	// Log the caught signal

	//srv.logger.Info("shutting down server", "signal_details", map[string]interface{}{
	//	"signal": s.String(),
	//})

	srv.logger.Info("starting stopped", slog.Group("signal_details",
		slog.String("signal", s.String()),
	))

	// Begin shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		shutdownError <- err
		srv.logger.Error("HTTP shutdown error", err)
	}

	// Log a message to say we're waiting for any background tasks

	//srv.logger.Info("completing background tasks", "bg_details", map[string]interface{}{
	//	"addr": httpServer.Addr,
	//})

	srv.logger.Info("completing background tasks", slog.Group("bg_details",
		slog.String("addr", httpServer.Addr),
	))

	// Wait for all goroutine to be done
	srv.wg.Wait()
	shutdownError <- nil
}
