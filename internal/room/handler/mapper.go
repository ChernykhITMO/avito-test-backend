package handler

import (
	"time"

	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
)

func toRoomResponse(room roommodel.Room) roomResponse {
	return roomResponse{
		ID:          room.ID.String(),
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		CreatedAt:   room.CreatedAt.UTC().Format(time.RFC3339),
	}
}
