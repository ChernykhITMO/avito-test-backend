package repository

import (
	"context"
	"fmt"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
)

func (r *Repository) Create(ctx context.Context, params authmodel.CreateUserParams) error {
	const op = "internal.auth.repository.Repository.Create"

	querier := postgres.QuerierFromContext(ctx, r.db)

	_, err := querier.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, params.User.ID, params.User.Email, params.PasswordHash, params.User.Role, params.User.CreatedAt)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return authmodel.ErrEmailAlreadyExists
		}

		return fmt.Errorf("%s: exec insert user: %w", op, err)
	}

	return nil
}
