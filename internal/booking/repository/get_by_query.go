package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) getByQuery(ctx context.Context, query string, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	const op = "internal.booking.repository.Repository.getByQuery"

	querier := postgres.QuerierFromContext(ctx, r.db)

	booking := &bookingmodel.Booking{}
	var conferenceLink sql.NullString
	var cancelledAt sql.NullTime

	err := querier.QueryRow(ctx, query, bookingID).Scan(
		&booking.ID,
		&booking.SlotID,
		&booking.UserID,
		&booking.Status,
		&conferenceLink,
		&booking.CreatedAt,
		&cancelledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: query booking: %w", op, err)
	}

	if conferenceLink.Valid {
		value := conferenceLink.String
		booking.ConferenceLink = &value
	}

	if cancelledAt.Valid {
		value := cancelledAt.Time
		booking.CancelledAt = &value
	}

	return booking, nil
}
