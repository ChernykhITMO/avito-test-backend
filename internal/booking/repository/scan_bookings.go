package repository

import (
	"database/sql"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
)

type bookingRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func scanBookings(rows bookingRows) ([]bookingmodel.Booking, error) {
	const op = "internal.booking.repository.scanBookings"

	bookings := make([]bookingmodel.Booking, 0)
	for rows.Next() {
		var booking bookingmodel.Booking
		var conferenceLink sql.NullString
		var cancelledAt sql.NullTime

		if err := rows.Scan(
			&booking.ID,
			&booking.SlotID,
			&booking.UserID,
			&booking.Status,
			&conferenceLink,
			&booking.CreatedAt,
			&cancelledAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan booking: %w", op, err)
		}

		if conferenceLink.Valid {
			value := conferenceLink.String
			booking.ConferenceLink = &value
		}

		if cancelledAt.Valid {
			value := cancelledAt.Time
			booking.CancelledAt = &value
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return bookings, nil
}
