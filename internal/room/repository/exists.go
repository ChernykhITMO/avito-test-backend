package repository

import (
	"context"
	"fmt"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
)

func (r *Repository) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	const op = "internal.room.repository.Repository.Exists"

	querier := postgres.QuerierFromContext(ctx, r.db)

	var exists bool
	err := querier.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)`, roomID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query exists: %w", op, err)
	}

	return exists, nil
}
