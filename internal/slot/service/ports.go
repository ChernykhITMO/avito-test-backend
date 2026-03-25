package service

import (
	"context"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

type ScheduleRepository interface {
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error)
	UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error
}

type SlotRepository interface {
	ListAvailableByRoomAndDate(ctx context.Context, roomID uuid.UUID, date time.Time) ([]slotmodel.Slot, error)
	CreateBatch(ctx context.Context, slots []slotmodel.Slot) error
	GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error)
}

type RoomRepository interface {
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type Transactor = postgres.Transactor

type Clock interface {
	Now() time.Time
}
