package service

import (
	"context"

	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Create(ctx context.Context, room roommodel.Room) error
	List(ctx context.Context) ([]roommodel.Room, error)
	Exists(ctx context.Context, roomID uuid.UUID) (bool, error)
}
