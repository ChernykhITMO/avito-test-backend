package repository

import (
	"context"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
)

func (r *Repository) Create(ctx context.Context, booking bookingmodel.Booking) error {
	const op = "internal.booking.repository.Repository.Create"

	querier := postgres.QuerierFromContext(ctx, r.db)

	_, err := querier.Exec(ctx, `
		INSERT INTO bookings (id, slot_id, user_id, status, conference_link, created_at, cancelled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, booking.ID, booking.SlotID, booking.UserID, booking.Status, booking.ConferenceLink, booking.CreatedAt, booking.CancelledAt)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return bookingmodel.ErrSlotAlreadyBooked
		}

		return fmt.Errorf("%s: exec insert booking: %w", op, err)
	}

	return nil
}
