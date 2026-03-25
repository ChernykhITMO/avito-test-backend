package service

import (
	"context"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/google/uuid"
)

func (s *Service) ListMyFuture(ctx context.Context, userID uuid.UUID) ([]bookingmodel.Booking, error) {
	const op = "internal.booking.service.Service.ListMyFuture"

	bookings, err := s.bookingRepository.ListMyFuture(ctx, userID, s.clock.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("%s: list my future bookings: %w", op, err)
	}

	return bookings, nil
}
