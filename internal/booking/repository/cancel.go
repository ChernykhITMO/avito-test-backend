package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
)

func (r *Repository) Cancel(ctx context.Context, bookingID uuid.UUID, cancelledAt time.Time) error {
	const op = "internal.booking.repository.Repository.Cancel"

	querier := postgres.QuerierFromContext(ctx, r.db)

	_, err := querier.Exec(ctx, `
		UPDATE bookings
		SET status = 'cancelled', cancelled_at = $2
		WHERE id = $1
	`, bookingID, cancelledAt)
	if err != nil {
		return fmt.Errorf("%s: exec update booking cancel: %w", op, err)
	}

	return nil
}
