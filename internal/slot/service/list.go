package service

import (
	"context"
	"fmt"
	"time"

	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

func (s *Service) ListAvailable(ctx context.Context, roomID uuid.UUID, date time.Time) ([]slotmodel.Slot, error) {
	const op = "internal.slot.service.Service.ListAvailable"

	exists, err := s.roomRepository.Exists(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: check room exists: %w", op, err)
	}
	if !exists {
		return nil, ErrRoomNotFound
	}

	schedule, err := s.scheduleRepository.GetByRoomID(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s: get schedule by room: %w", op, err)
	}
	if schedule == nil {
		return []slotmodel.Slot{}, nil
	}

	requestedDate := normalizeSlotDate(date)
	if schedule.GeneratedUntil.Before(requestedDate) {
		if err := s.EnsureRange(ctx, roomID, requestedDate); err != nil {
			return nil, fmt.Errorf("%s: ensure range: %w", op, err)
		}
	}

	slots, err := s.slotRepository.ListAvailableByRoomAndDate(ctx, roomID, requestedDate)
	if err != nil {
		return nil, fmt.Errorf("%s: list available slots: %w", op, err)
	}

	return slots, nil
}
