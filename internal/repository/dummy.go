package repository

import (
	"context"
	"database/sql"
)

type Dummy struct {
	ID string
}

type DummyModel struct {
	db *sql.DB
}

type DummyRepository interface {
	Get(ctx context.Context) (*Dummy, error)
}

func NewDummyRepository(db *sql.DB) DummyRepository {
	return &DummyModel{
		db: db,
	}
}

func (m DummyModel) Get(ctx context.Context) (*Dummy, error) {
	return &Dummy{
		ID: "123",
	}, nil
}
