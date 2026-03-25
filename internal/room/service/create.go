package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	roommodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/model"
	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, input CreateInput) (*roommodel.Room, error) {
	const op = "internal.room.service.Service.Create"

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidRoom
	}
	if input.Capacity != nil && *input.Capacity <= 0 {
		return nil, ErrInvalidRoom
	}

	room := roommodel.Room{
		ID:          uuid.New(),
		Name:        name,
		Description: input.Description,
		Capacity:    input.Capacity,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.roomRepository.Create(ctx, room); err != nil {
		return nil, fmt.Errorf("%s: create room: %w", op, err)
	}

	return &room, nil
}
