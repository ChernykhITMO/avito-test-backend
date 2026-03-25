package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error) {
	const op = "internal.slot.repository.Repository.GetByID"

	querier := postgres.QuerierFromContext(ctx, r.db)

	slot := &slotmodel.Slot{}
	err := querier.QueryRow(ctx, `
		SELECT id, room_id, schedule_id, slot_date, start_at, end_at, created_at
		FROM slots
		WHERE id = $1
	`, slotID).Scan(
		&slot.ID,
		&slot.RoomID,
		&slot.ScheduleID,
		&slot.SlotDate,
		&slot.StartAt,
		&slot.EndAt,
		&slot.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: query slot by id: %w", op, err)
	}

	return slot, nil
}
