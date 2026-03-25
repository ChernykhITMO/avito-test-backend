package repository

import (
	"context"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
)

func (r *Repository) ListAll(ctx context.Context, page, pageSize int) ([]bookingmodel.Booking, int, error) {
	const op = "internal.booking.repository.Repository.ListAll"

	querier := postgres.QuerierFromContext(ctx, r.db)

	var total int
	if err := querier.QueryRow(ctx, `SELECT COUNT(*) FROM bookings`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%s: count bookings: %w", op, err)
	}

	offset := (page - 1) * pageSize
	rows, err := querier.Query(ctx, `
		SELECT id, slot_id, user_id, status, conference_link, created_at, cancelled_at
		FROM bookings
		ORDER BY created_at DESC, id DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: query bookings page: %w", op, err)
	}
	defer rows.Close()

	bookings, err := scanBookings(rows)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: scan bookings: %w", op, err)
	}

	return bookings, total, nil
}
