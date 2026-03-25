package service

import (
	"context"
	"errors"
	"fmt"

	schedulemodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/model"
	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, input CreateInput) (*schedulemodel.Schedule, error) {
	const op = "internal.schedule.service.Service.Create"

	exists, err := s.roomRepository.Exists(ctx, input.RoomID)
	if err != nil {
		return nil, fmt.Errorf("%s: check room exists: %w", op, err)
	}
	if !exists {
		return nil, ErrRoomNotFound
	}

	daysOfWeek, err := schedulemodel.NormalizeDaysOfWeek(input.DaysOfWeek)
	if err != nil {
		if errors.Is(err, schedulemodel.ErrInvalidSchedule) {
			return nil, ErrInvalidSchedule
		}

		return nil, fmt.Errorf("%s: normalize days of week: %w", op, err)
	}

	startMinute, err := schedulemodel.ParseClock(input.StartTime)
	if err != nil {
		if errors.Is(err, schedulemodel.ErrInvalidSchedule) {
			return nil, ErrInvalidSchedule
		}

		return nil, fmt.Errorf("%s: parse start time: %w", op, err)
	}

	endMinute, err := schedulemodel.ParseClock(input.EndTime)
	if err != nil {
		if errors.Is(err, schedulemodel.ErrInvalidSchedule) {
			return nil, ErrInvalidSchedule
		}

		return nil, fmt.Errorf("%s: parse end time: %w", op, err)
	}

	if startMinute >= endMinute || (endMinute-startMinute)%30 != 0 {
		return nil, ErrInvalidSchedule
	}

	now := s.clock.Now().UTC()
	today := normalizeScheduleDate(now)
	schedule := schedulemodel.Schedule{
		ID:             uuid.New(),
		RoomID:         input.RoomID,
		DaysOfWeek:     daysOfWeek,
		StartMinute:    startMinute,
		EndMinute:      endMinute,
		GeneratedUntil: today.AddDate(0, 0, -1),
		CreatedAt:      now,
	}

	horizonDate := today.AddDate(0, 0, s.windowDays-1)

	if err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.scheduleRepository.Create(txCtx, schedule); err != nil {
			return fmt.Errorf("%s: create schedule: %w", op, err)
		}

		if err := s.slotGenerator.EnsureRange(txCtx, input.RoomID, horizonDate); err != nil {
			return fmt.Errorf("%s: ensure slot range: %w", op, err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, schedulemodel.ErrScheduleAlreadyExists) {
			return nil, ErrScheduleAlreadyExists
		}

		return nil, fmt.Errorf("%s: transaction: %w", op, err)
	}

	schedule.GeneratedUntil = horizonDate

	return &schedule, nil
}
