package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
)

func (r *Repository) List(ctx context.Context) ([]roommodel.Room, error) {
	const op = "internal.room.repository.Repository.List"

	querier := postgres.QuerierFromContext(ctx, r.db)

	rows, err := querier.Query(ctx, `
		SELECT id, name, description, capacity, created_at
		FROM rooms
		ORDER BY created_at DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: query rooms: %w", op, err)
	}
	defer rows.Close()

	rooms := make([]roommodel.Room, 0)

	for rows.Next() {
		var room roommodel.Room
		var description sql.NullString
		var capacity sql.NullInt32

		if err := rows.Scan(&room.ID, &room.Name, &description, &capacity, &room.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: scan room: %w", op, err)
		}

		if description.Valid {
			value := description.String
			room.Description = &value
		}

		if capacity.Valid {
			value := int(capacity.Int32)
			room.Capacity = &value
		}

		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return rooms, nil
}
