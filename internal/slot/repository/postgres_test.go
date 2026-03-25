package repository

import (
	"context"
	"testing"
	"time"

	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestSlotRepositoryQueries(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer mock.Close()

	repo := New(mock)
	roomID := uuid.New()
	scheduleID := uuid.New()
	slotID := uuid.New()
	date := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	startAt := date.Add(9 * time.Hour)
	endAt := startAt.Add(30 * time.Minute)
	createdAt := time.Now().UTC()

	rows := pgxmock.NewRows([]string{"id", "room_id", "schedule_id", "slot_date", "start_at", "end_at", "created_at"}).
		AddRow(slotID, roomID, scheduleID, date, startAt, endAt, createdAt)
	mock.ExpectQuery("SELECT s.id, s.room_id, s.schedule_id, s.slot_date, s.start_at, s.end_at, s.created_at").
		WithArgs(roomID, date).
		WillReturnRows(rows)

	slots, err := repo.ListAvailableByRoomAndDate(context.Background(), roomID, date)
	if err != nil {
		t.Fatalf("ListAvailableByRoomAndDate: %v", err)
	}
	if len(slots) != 1 {
		t.Fatalf("expected 1 slot, got %d", len(slots))
	}

	batch := []slotmodel.Slot{{
		ID:         slotID,
		RoomID:     roomID,
		ScheduleID: scheduleID,
		SlotDate:   date,
		StartAt:    startAt,
		EndAt:      endAt,
		CreatedAt:  createdAt,
	}}
	mock.ExpectExec("INSERT INTO slots").
		WithArgs(
			slotID,
			roomID,
			scheduleID,
			date,
			startAt,
			endAt,
			createdAt,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err := repo.CreateBatch(context.Background(), batch); err != nil {
		t.Fatalf("CreateBatch: %v", err)
	}

	getRows := pgxmock.NewRows([]string{"id", "room_id", "schedule_id", "slot_date", "start_at", "end_at", "created_at"}).
		AddRow(slotID, roomID, scheduleID, date, startAt, endAt, createdAt)
	mock.ExpectQuery("SELECT id, room_id, schedule_id, slot_date, start_at, end_at, created_at FROM slots").
		WithArgs(slotID).
		WillReturnRows(getRows)

	got, err := repo.GetByID(context.Background(), slotID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got == nil || got.ID != slotID {
		t.Fatalf("unexpected slot: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
