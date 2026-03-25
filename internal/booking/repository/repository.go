package repository

import "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"

type Repository struct {
	db postgres.DBTX
}

func New(db postgres.DBTX) *Repository {
	return &Repository{db: db}
}
