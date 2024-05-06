package repository

import (
	"errors"
)

// Store all repo related

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repositories struct {
	DummyRepository DummyRepository
}

func NewRepositories(dummyRepository DummyRepository) Repositories {
	return Repositories{
		DummyRepository: dummyRepository,
	}
}
