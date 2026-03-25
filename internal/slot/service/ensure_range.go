package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *Service) EnsureRange(ctx context.Context, roomID uuid.UUID, toDate time.Time) error {
	const op = "internal.slot.service.Service.EnsureRange"

	if err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		schedule, err := s.scheduleRepository.GetByRoomID(txCtx, roomID)
		if err != nil {
			return fmt.Errorf("%s: get schedule by room: %w", op, err)
		}
		if schedule == nil {
			return nil
		}

		targetDate := normalizeSlotDate(toDate)
		startDate := normalizeSlotDate(s.clock.Now())
		if schedule.GeneratedUntil.After(startDate.AddDate(0, 0, -1)) {
			startDate = schedule.GeneratedUntil.AddDate(0, 0, 1)
		}

		if startDate.After(targetDate) {
			return nil
		}

		slots := generateSlots(*schedule, startDate, targetDate)
		if len(slots) > 0 {
			if err := s.slotRepository.CreateBatch(txCtx, slots); err != nil {
				return fmt.Errorf("%s: create slots batch: %w", op, err)
			}
		}

		if err := s.scheduleRepository.UpdateGeneratedUntil(txCtx, roomID, targetDate); err != nil {
			return fmt.Errorf("%s: update generated until: %w", op, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("%s: transaction: %w", op, err)
	}

	return nil
}
