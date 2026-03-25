package repository

import (
	"context"
	"fmt"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
)

func (r *Repository) Create(ctx context.Context, schedule schedulemodel.Schedule) error {
	const op = "internal.schedule.repository.Repository.Create"

	querier := postgres.QuerierFromContext(ctx, r.db)

	_, err := querier.Exec(ctx, `
		INSERT INTO schedules (id, room_id, days_of_week, start_minute, end_minute, generated_until, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, schedule.ID, schedule.RoomID, toSmallIntDays(schedule.DaysOfWeek), schedule.StartMinute, schedule.EndMinute, schedule.GeneratedUntil, schedule.CreatedAt)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return schedulemodel.ErrScheduleAlreadyExists
		}

		return fmt.Errorf("%s: exec insert schedule: %w", op, err)
	}

	return nil
}
