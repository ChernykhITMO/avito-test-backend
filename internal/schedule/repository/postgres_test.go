package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestScheduleRepositoryCreateConflictAndRead(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer mock.Close()

	repo := New(mock)
	schedule := schedulemodel.Schedule{
		ID:             uuid.New(),
		RoomID:         uuid.New(),
		DaysOfWeek:     []int{1, 3},
		StartMinute:    540,
		EndMinute:      600,
		GeneratedUntil: time.Date(2026, 3, 21, 0, 0, 0, 0, time.UTC),
		CreatedAt:      time.Now().UTC(),
	}

	mock.ExpectExec("INSERT INTO schedules").
		WithArgs(
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
		).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	err = repo.Create(context.Background(), schedule)
	if !errors.Is(err, schedulemodel.ErrScheduleAlreadyExists) {
		t.Fatalf("expected schedule exists error, got %v", err)
	}

	rows := pgxmock.NewRows([]string{"id", "room_id", "days_of_week", "start_minute", "end_minute", "generated_until", "created_at"}).
		AddRow(schedule.ID, schedule.RoomID, []int16{1, 3}, schedule.StartMinute, schedule.EndMinute, schedule.GeneratedUntil, schedule.CreatedAt)
	mock.ExpectQuery("SELECT id, room_id, days_of_week, start_minute, end_minute, generated_until, created_at").
		WithArgs(schedule.RoomID).
		WillReturnRows(rows)

	got, err := repo.GetByRoomID(context.Background(), schedule.RoomID)
	if err != nil {
		t.Fatalf("GetByRoomID: %v", err)
	}
	if got == nil || got.ID != schedule.ID {
		t.Fatalf("unexpected schedule: %#v", got)
	}
	if len(got.DaysOfWeek) != 2 || got.DaysOfWeek[0] != 1 || got.DaysOfWeek[1] != 3 {
		t.Fatalf("unexpected days of week: %#v", got.DaysOfWeek)
	}

	mock.ExpectExec("UPDATE schedules SET generated_until").
		WithArgs(schedule.RoomID, schedule.GeneratedUntil).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err := repo.UpdateGeneratedUntil(context.Background(), schedule.RoomID, schedule.GeneratedUntil); err != nil {
		t.Fatalf("UpdateGeneratedUntil: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
