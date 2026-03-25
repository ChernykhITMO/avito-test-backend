package service

import (
	"context"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/google/uuid"
)

func (s *Service) Cancel(ctx context.Context, bookingID, userID uuid.UUID) (*bookingmodel.Booking, error) {
	const op = "internal.booking.service.Service.Cancel"

	var result *bookingmodel.Booking

	err := s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		booking, err := s.bookingRepository.GetByIDForUpdate(txCtx, bookingID)
		if err != nil {
			return fmt.Errorf("%s: get booking for update: %w", op, err)
		}
		if booking == nil {
			return ErrBookingNotFound
		}
		if booking.UserID != userID {
			return ErrForbiddenCancel
		}
		if booking.Status == bookingmodel.StatusCancelled {
			result = booking
			return nil
		}

		cancelledAt := s.clock.Now().UTC()
		if err := s.bookingRepository.Cancel(txCtx, bookingID, cancelledAt); err != nil {
			return fmt.Errorf("%s: cancel booking: %w", op, err)
		}

		booking.Status = bookingmodel.StatusCancelled
		booking.CancelledAt = &cancelledAt
		result = booking
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: transaction: %w", op, err)
	}

	return result, nil
}
