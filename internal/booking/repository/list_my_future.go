package repository

import (
	"context"
	"fmt"
	"time"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
)

func (r *Repository) ListMyFuture(ctx context.Context, userID uuid.UUID, now time.Time) ([]bookingmodel.Booking, error) {
	const op = "internal.booking.repository.Repository.ListMyFuture"

	querier := postgres.QuerierFromContext(ctx, r.db)

	rows, err := querier.Query(ctx, `
		SELECT b.id, b.slot_id, b.user_id, b.status, b.conference_link, b.created_at, b.cancelled_at
		FROM bookings b
		JOIN slots s ON s.id = b.slot_id
		WHERE b.user_id = $1
		  AND s.start_at >= $2
		ORDER BY s.start_at ASC, b.created_at ASC
	`, userID, now)
	if err != nil {
		return nil, fmt.Errorf("%s: query my future bookings: %w", op, err)
	}
	defer rows.Close()

	bookings, err := scanBookings(rows)
	if err != nil {
		return nil, fmt.Errorf("%s: scan bookings: %w", op, err)
	}

	return bookings, nil
}
