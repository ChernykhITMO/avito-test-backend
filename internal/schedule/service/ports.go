package service

import (
	"context"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule schedulemodel.Schedule) error
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*schedulemodel.Schedule, error)
	UpdateGeneratedUntil(ctx context.Context, roomID uuid.UUID, generatedUntil time.Time) error
}

type SlotGenerator interface {
	EnsureRange(ctx context.Context, roomID uuid.UUID, toDate time.Time) error
}

type Transactor = postgres.Transactor

type Clock interface {
	Now() time.Time
}
