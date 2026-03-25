package repository

import (
	"context"
	"fmt"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
)

func (r *Repository) Create(ctx context.Context, room roommodel.Room) error {
	const op = "internal.room.repository.Repository.Create"

	querier := postgres.QuerierFromContext(ctx, r.db)

	_, err := querier.Exec(ctx, `
		INSERT INTO rooms (id, name, description, capacity, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, room.ID, room.Name, room.Description, room.Capacity, room.CreatedAt)
	if err != nil {
		return fmt.Errorf("%s: exec insert room: %w", op, err)
	}

	return nil
}
