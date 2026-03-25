package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
)

func (r *Repository) UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error {
	const op = "internal.schedule.repository.Repository.UpdateGeneratedUntil"

	querier := postgres.QuerierFromContext(ctx, r.db)

	_, err := querier.Exec(ctx, `
		UPDATE schedules
		SET generated_until = GREATEST(generated_until, $2::date)
		WHERE room_id = $1
	`, roomID, generatedUntil)
	if err != nil {
		return fmt.Errorf("%s: exec update generated_until: %w", op, err)
	}

	return nil
}
