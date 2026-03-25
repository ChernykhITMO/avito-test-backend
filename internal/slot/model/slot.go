package model

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID         uuid.UUID
	RoomID     uuid.UUID
	ScheduleID uuid.UUID
	SlotDate   time.Time
	StartAt    time.Time
	EndAt      time.Time
	CreatedAt  time.Time
}
