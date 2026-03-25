package service

import (
	"context"
	"errors"
	"testing"
	"time"

	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	"github.com/google/uuid"
)

type roomRepoStub struct {
	exists bool
	err    error
}

func (s roomRepoStub) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	return s.exists, s.err
}

type scheduleRepoStub struct {
	created *schedulemodel.Schedule
}

func (s *scheduleRepoStub) Create(ctx context.Context, schedule schedulemodel.Schedule) error {
	s.created = &schedule
	return nil
}

func (s *scheduleRepoStub) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error) {
	return nil, nil
}

func (s *scheduleRepoStub) UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error {
	return nil
}

type slotGeneratorStub struct {
	roomID uuid.UUID
	toDate time.Time
}

func (s *slotGeneratorStub) EnsureRange(ctx context.Context, roomID uuid.UUID, toDate time.Time) error {
	s.roomID = roomID
	s.toDate = toDate
	return nil
}

type noopTransactor struct{}

func (noopTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

func TestCreateScheduleRoomNotFound(t *testing.T) {
	service := New(roomRepoStub{exists: false}, &scheduleRepoStub{}, &slotGeneratorStub{}, noopTransactor{}, fixedClock{now: time.Now().UTC()}, 30)

	_, err := service.Create(context.Background(), CreateInput{
		RoomID:     uuid.New(),
		DaysOfWeek: []int{1, 2},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("expected room not found, got %v", err)
	}
}

func TestCreateScheduleGeneratesInitialRange(t *testing.T) {
	now := time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC)
	repo := &scheduleRepoStub{}
	generator := &slotGeneratorStub{}
	roomID := uuid.New()
	service := New(roomRepoStub{exists: true}, repo, generator, noopTransactor{}, fixedClock{now: now}, 30)

	schedule, err := service.Create(context.Background(), CreateInput{
		RoomID:     roomID,
		DaysOfWeek: []int{1, 3, 5},
		StartTime:  "09:00",
		EndTime:    "10:00",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if schedule.RoomID != roomID {
		t.Fatalf("unexpected room id: %s", schedule.RoomID)
	}
	if repo.created == nil {
		t.Fatal("expected schedule to be persisted")
	}
	if generator.roomID != roomID {
		t.Fatalf("unexpected generator room id: %s", generator.roomID)
	}
}
