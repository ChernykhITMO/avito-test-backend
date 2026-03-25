package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
)

func (r *Repository) ListRoomIDsNeedingGeneration(
	ctx context.Context,
	generatedUntilBefore time.Time,
	limit int,
) ([]uuid.UUID, error) {
	const op = "internal.schedule.repository.Repository.ListRoomIDsNeedingGeneration"

	querier := postgres.QuerierFromContext(ctx, r.db)

	rows, err := querier.Query(ctx, `
		SELECT room_id
		FROM schedules
		WHERE generated_until < $1
		ORDER BY generated_until ASC, room_id ASC
		LIMIT $2
	`, generatedUntilBefore, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: query schedules for refill: %w", op, err)
	}
	defer rows.Close()

	roomIDs := make([]uuid.UUID, 0, limit)
	for rows.Next() {
		var roomID uuid.UUID
		if err := rows.Scan(&roomID); err != nil {
			return nil, fmt.Errorf("%s: scan room id: %w", op, err)
		}

		roomIDs = append(roomIDs, roomID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return roomIDs, nil
}
