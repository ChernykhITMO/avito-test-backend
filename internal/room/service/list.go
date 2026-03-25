package service

import (
	"context"
	"fmt"

	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
)

func (s *Service) List(ctx context.Context) ([]roommodel.Room, error) {
	const op = "internal.room.service.Service.List"

	rooms, err := s.roomRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: list rooms: %w", op, err)
	}

	return rooms, nil
}
