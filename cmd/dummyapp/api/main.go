package main

import (
	"flag"
	"github.com/go-easy-templ/internal/config"
	"github.com/go-easy-templ/internal/database"
	"github.com/go-easy-templ/internal/handler"
	"github.com/go-easy-templ/internal/logger"
	"github.com/go-easy-templ/internal/repository"
	"github.com/go-easy-templ/internal/server"
	"github.com/go-easy-templ/internal/service"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func main() {

	// Initialize dependencies
	var logLevel string

	flag.StringVar(&logLevel, "log-level", "INFO", "Log level")
	flag.Parse()

	logger := logger.NewSlogger(logLevel)

	config, err := config.InitConfig()
	if err != nil {
		logger.Error("Config fail to initialised", err)
		os.Exit(1)
	}

	// Establish database connection
	db, err := database.NewDatabase(config)
	if err != nil {
		logger.Error("database fail to establish", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	dummyRepo := repository.NewDummyRepository(db)
	repositories := repository.NewRepositories(dummyRepo)
	dummyService := service.NewDummy(config, logger, repositories)
	services := service.NewServices(dummyService)

	// Initialize handlers
	healthcheckHandler := handler.NewHealthcheckHandler(logger, config, services)
	handlers := handler.NewHandlers(healthcheckHandler)

	// Run server
	srv := server.NewServer(config, logger, handlers)

	if err := srv.Run(); err != nil {
		logger.Error("application unable to run", slog.Any("error", err))
		os.Exit(1)
	}
}
