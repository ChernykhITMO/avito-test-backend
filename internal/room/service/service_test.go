package service

import (
	"context"
	"errors"
	"testing"

	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
	"github.com/google/uuid"
)

type roomRepoServiceStub struct {
	rooms []roommodel.Room
}

func (s *roomRepoServiceStub) Create(ctx context.Context, room roommodel.Room) error {
	s.rooms = append(s.rooms, room)
	return nil
}

func (s *roomRepoServiceStub) List(ctx context.Context) ([]roommodel.Room, error) {
	return s.rooms, nil
}

func (s *roomRepoServiceStub) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	return true, nil
}

func TestCreateAndListRooms(t *testing.T) {
	repo := &roomRepoServiceStub{}
	service := New(repo)

	room, err := service.Create(context.Background(), CreateInput{Name: "Alpha"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if room.Name != "Alpha" {
		t.Fatalf("unexpected room: %#v", room)
	}

	rooms, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(rooms) != 1 {
		t.Fatalf("expected 1 room, got %d", len(rooms))
	}
}

func TestCreateRoomRejectsInvalidInput(t *testing.T) {
	repo := &roomRepoServiceStub{}
	service := New(repo)

	_, err := service.Create(context.Background(), CreateInput{Name: "   "})
	if !errors.Is(err, ErrInvalidRoom) {
		t.Fatalf("expected invalid room, got %v", err)
	}

	capacity := 0
	_, err = service.Create(context.Background(), CreateInput{Name: "Alpha", Capacity: &capacity})
	if !errors.Is(err, ErrInvalidRoom) {
		t.Fatalf("expected invalid room, got %v", err)
	}

	if len(repo.rooms) != 0 {
		t.Fatalf("expected no rooms to be created, got %d", len(repo.rooms))
	}
}
