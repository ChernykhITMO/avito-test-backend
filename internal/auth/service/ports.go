package service

import (
	"context"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, params authmodel.CreateUserParams) error
	GetByEmail(ctx context.Context, email string) (*authmodel.UserWithPassword, error)
}

type TokenManager interface {
	Issue(userID uuid.UUID, role authmodel.Role) (string, error)
}

type PasswordManager interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type Transactor = postgres.Transactor
