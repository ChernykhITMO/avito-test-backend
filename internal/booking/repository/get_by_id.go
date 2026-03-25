package repository

import (
	"context"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/google/uuid"
)

func (r *Repository) GetByID(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	const op = "internal.booking.repository.Repository.GetByID"

	booking, err := r.getByQuery(ctx, `
		SELECT id, slot_id, user_id, status, conference_link, created_at, cancelled_at
		FROM bookings
		WHERE id = $1
	`, bookingID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return booking, nil
}
