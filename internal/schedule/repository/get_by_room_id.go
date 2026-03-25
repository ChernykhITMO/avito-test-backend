package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error) {
	const op = "internal.schedule.repository.Repository.GetByRoomID"

	querier := postgres.QuerierFromContext(ctx, r.db)

	schedule := &schedulemodel.Schedule{}
	var daysOfWeek []int16
	err := querier.QueryRow(ctx, `
		SELECT id, room_id, days_of_week, start_minute, end_minute, generated_until, created_at
		FROM schedules
		WHERE room_id = $1
	`, roomID).Scan(
		&schedule.ID,
		&schedule.RoomID,
		&daysOfWeek,
		&schedule.StartMinute,
		&schedule.EndMinute,
		&schedule.GeneratedUntil,
		&schedule.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s: query schedule by room: %w", op, err)
	}

	schedule.DaysOfWeek = fromSmallIntDays(daysOfWeek)

	return schedule, nil
}
