package service

import (
	"context"
	"fmt"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
)

func (s *Service) ListAll(ctx context.Context, page, pageSize int) ([]bookingmodel.Booking, int, error) {
	const op = "internal.booking.service.Service.ListAll"

	bookings, total, err := s.bookingRepository.ListAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: list all bookings: %w", op, err)
	}

	return bookings, total, nil
}
