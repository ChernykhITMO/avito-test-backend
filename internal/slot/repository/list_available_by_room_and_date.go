package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

func (r *Repository) ListAvailableByRoomAndDate(
	ctx context.Context,
	roomID uuid.UUID,
	date time.Time,
) ([]slotmodel.Slot, error) {
	const op = "internal.slot.repository.Repository.ListAvailableByRoomAndDate"

	querier := postgres.QuerierFromContext(ctx, r.db)

	rows, err := querier.Query(ctx, `
		SELECT s.id, s.room_id, s.schedule_id, s.slot_date, s.start_at, s.end_at, s.created_at
		FROM slots s
		LEFT JOIN bookings b
			ON b.slot_id = s.id AND b.status = 'active'
		WHERE s.room_id = $1
		  AND s.slot_date = $2
		  AND b.id IS NULL
		ORDER BY s.start_at ASC
	`, roomID, date)
	if err != nil {
		return nil, fmt.Errorf("%s: query available slots: %w", op, err)
	}
	defer rows.Close()

	slots := make([]slotmodel.Slot, 0)
	for rows.Next() {
		var slot slotmodel.Slot
		if err := rows.Scan(
			&slot.ID,
			&slot.RoomID,
			&slot.ScheduleID,
			&slot.SlotDate,
			&slot.StartAt,
			&slot.EndAt,
			&slot.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: scan slot: %w", op, err)
		}

		slots = append(slots, slot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return slots, nil
}
