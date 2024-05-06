package service

import (
	"context"
	"fmt"
	"github.com/go-easy-templ/internal/config"
	"github.com/go-easy-templ/internal/repository"
	"log/slog"
)

type Dummy struct {
	config       *config.Config
	logger       *slog.Logger
	repositories repository.Repositories
}

type DummyService interface {
	DoSomething(ctx context.Context) error
}

func NewDummy(config *config.Config, logger *slog.Logger, repositories repository.Repositories) DummyService {
	return &Dummy{
		config:       config,
		logger:       logger,
		repositories: repositories,
	}
}

func (s *Dummy) DoSomething(ctx context.Context) error {
	s.logger.InfoContext(ctx, "Doing something important")
	get, err := s.repositories.DummyRepository.Get(ctx)
	if err != nil {
		return err
	}
	s.logger.DebugContext(ctx, fmt.Sprintf("ID: %s", get.ID))
	return nil
}
