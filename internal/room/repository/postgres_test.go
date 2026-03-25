package repository

import (
	"context"
	"testing"
	"time"

	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
	"github.com/google/uuid"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestRoomRepositoryCRUDQueries(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("NewPool: %v", err)
	}
	defer mock.Close()

	repo := New(mock)
	roomID := uuid.New()
	now := time.Now().UTC()
	description := "Room A"
	capacity := 10

	mock.ExpectExec("INSERT INTO rooms").
		WithArgs(roomID, "A", &description, &capacity, now).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err := repo.Create(context.Background(), roommodel.Room{
		ID:          roomID,
		Name:        "A",
		Description: &description,
		Capacity:    &capacity,
		CreatedAt:   now,
	}); err != nil {
		t.Fatalf("Create: %v", err)
	}

	rows := pgxmock.NewRows([]string{"id", "name", "description", "capacity", "created_at"}).
		AddRow(roomID, "A", description, int32(capacity), now)
	mock.ExpectQuery("SELECT id, name, description, capacity, created_at FROM rooms").
		WillReturnRows(rows)

	list, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 room, got %d", len(list))
	}

	existsRows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(roomID).
		WillReturnRows(existsRows)

	exists, err := repo.Exists(context.Background(), roomID)
	if err != nil {
		t.Fatalf("Exists: %v", err)
	}
	if !exists {
		t.Fatal("expected room to exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
