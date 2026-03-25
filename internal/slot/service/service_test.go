package service

import (
	"context"
	"errors"
	"testing"
	"time"

	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

type slotRoomRepoStub struct {
	exists bool
}

func (s slotRoomRepoStub) Exists(ctx context.Context, roomID uuid.UUID) (bool, error) {
	return s.exists, nil
}

type slotScheduleRepoStub struct {
	schedule *schedulemodel.Schedule
	updated  bool
}

func (s *slotScheduleRepoStub) GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error) {
	return s.schedule, nil
}

func (s *slotScheduleRepoStub) UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error {
	s.updated = true
	return nil
}

type slotRepoStub struct {
	slots       []slotmodel.Slot
	batch       []slotmodel.Slot
	listErr     error
	getByIDSlot *slotmodel.Slot
}

func (s *slotRepoStub) ListAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]slotmodel.Slot, error) {
	return s.slots, s.listErr
}

func (s *slotRepoStub) CreateBatch(ctx context.Context, slots []slotmodel.Slot) error {
	s.batch = slots
	return nil
}

func (s *slotRepoStub) GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error) {
	return s.getByIDSlot, nil
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

func TestListAvailableRoomNotFound(t *testing.T) {
	service := New(slotRoomRepoStub{exists: false}, &slotScheduleRepoStub{}, &slotRepoStub{}, noopTransactor{}, fixedClock{now: time.Now().UTC()})

	_, err := service.ListAvailable(context.Background(), uuid.New(), time.Now().UTC())
	if !errors.Is(err, ErrRoomNotFound) {
		t.Fatalf("expected room not found, got %v", err)
	}
}

func TestEnsureRangeGeneratesSlots(t *testing.T) {
	roomID := uuid.New()
	schedule := &schedulemodel.Schedule{
		ID:             uuid.New(),
		RoomID:         roomID,
		DaysOfWeek:     []int{1},
		StartMinute:    540,
		EndMinute:      600,
		GeneratedUntil: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
	}
	scheduleRepo := &slotScheduleRepoStub{schedule: schedule}
	slotRepo := &slotRepoStub{}
	service := New(slotRoomRepoStub{exists: true}, scheduleRepo, slotRepo, noopTransactor{}, fixedClock{now: time.Date(2026, 3, 21, 12, 0, 0, 0, time.UTC)})

	if err := service.EnsureRange(context.Background(), roomID, time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("EnsureRange returned error: %v", err)
	}
	if len(slotRepo.batch) == 0 {
		t.Fatal("expected slots to be generated")
	}
	if !scheduleRepo.updated {
		t.Fatal("expected generated_until to be updated")
	}
}
