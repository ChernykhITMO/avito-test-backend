package service

import (
	"context"
	"time"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

type BookingRepository interface {
	Create(ctx context.Context, booking bookingmodel.Booking) error
	GetByID(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error)
	GetByIDForUpdate(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error)
	Cancel(ctx context.Context, bookingID uuid.UUID, cancelledAt time.Time) error
	ListMyFuture(ctx context.Context, userID uuid.UUID, now time.Time) ([]bookingmodel.Booking, error)
	ListAll(ctx context.Context, page, pageSize int) ([]bookingmodel.Booking, int, error)
}

type SlotRepository interface {
	GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error)
}

type Clock interface {
	Now() time.Time
}

type ConferenceService interface {
	CreateBookingLink(ctx context.Context, bookingID uuid.UUID) (string, error)
	CancelBookingLink(ctx context.Context, link string) error
}

type Transactor = postgres.Transactor
