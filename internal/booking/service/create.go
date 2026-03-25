package service

import (
	"context"
	"errors"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/google/uuid"
)

func (s *Service) Create(ctx context.Context, input CreateInput) (*bookingmodel.Booking, error) {
	const op = "internal.booking.service.Service.Create"

	var booking *bookingmodel.Booking

	err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		slot, err := s.slotRepository.GetByID(txCtx, input.SlotID)
		if err != nil {
			return fmt.Errorf("%s: get slot by id: %w", op, err)
		}
		if slot == nil {
			return ErrSlotNotFound
		}

		now := s.clock.Now().UTC()
		if slot.StartAt.Before(now) {
			return ErrSlotInPast
		}

		booking = &bookingmodel.Booking{
			ID:        uuid.New(),
			SlotID:    input.SlotID,
			UserID:    input.UserID,
			Status:    bookingmodel.StatusActive,
			CreatedAt: now,
		}

		if input.CreateConferenceLink {
			link, err := s.conferenceService.CreateBookingLink(txCtx, booking.ID)
			if err == nil {
				booking.ConferenceLink = &link
			}
		}

		if err := s.bookingRepository.Create(txCtx, *booking); err != nil {
			if booking.ConferenceLink != nil {
				_ = s.conferenceService.CancelBookingLink(txCtx, *booking.ConferenceLink)
			}
			return fmt.Errorf("%s: create booking: %w", op, err)
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, bookingmodel.ErrSlotAlreadyBooked) {
			return nil, ErrSlotAlreadyBooked
		}

		return nil, fmt.Errorf("%s: transaction: %w", op, err)
	}

	return booking, nil
}
