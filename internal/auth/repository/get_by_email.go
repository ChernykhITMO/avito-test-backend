package repository

import (
	"context"
	"errors"
	"fmt"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetByEmail(ctx context.Context, email string) (*authmodel.UserWithPassword, error) {
	const op = "internal.auth.repository.Repository.GetByEmail"

	querier := postgres.QuerierFromContext(ctx, r.db)

	user := &authmodel.UserWithPassword{}
	err := querier.QueryRow(ctx, `
		SELECT id, email, role, created_at, COALESCE(password_hash, '')
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.Role,
		&user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: query user by email: %w", op, err)
	}

	return user, nil
}
