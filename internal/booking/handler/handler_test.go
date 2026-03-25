package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	bookingservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/service"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

type bookingRepoHandlerStub struct {
	booking *bookingmodel.Booking
}

func (s *bookingRepoHandlerStub) Create(ctx context.Context, booking bookingmodel.Booking) error {
	s.booking = &booking
	return nil
}
func (s *bookingRepoHandlerStub) GetByID(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	return s.booking, nil
}
func (s *bookingRepoHandlerStub) GetByIDForUpdate(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	return s.booking, nil
}
func (s *bookingRepoHandlerStub) Cancel(ctx context.Context, bookingID uuid.UUID, cancelledAt time.Time) error {
	if s.booking != nil {
		s.booking.Status = bookingmodel.StatusCancelled
		s.booking.CancelledAt = &cancelledAt
	}
	return nil
}
func (s *bookingRepoHandlerStub) ListMyFuture(ctx context.Context, userID uuid.UUID, now time.Time) ([]bookingmodel.Booking, error) {
	if s.booking == nil {
		return nil, nil
	}
	return []bookingmodel.Booking{*s.booking}, nil
}
func (s *bookingRepoHandlerStub) ListAll(ctx context.Context, page, pageSize int) ([]bookingmodel.Booking, int, error) {
	if s.booking == nil {
		return nil, 0, nil
	}
	return []bookingmodel.Booking{*s.booking}, 1, nil
}

type bookingSlotRepoHandlerStub struct {
	slot *slotmodel.Slot
}

func (s bookingSlotRepoHandlerStub) GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error) {
	return s.slot, nil
}

type bookingConferenceStub struct{}

func (bookingConferenceStub) CreateBookingLink(ctx context.Context, bookingID uuid.UUID) (string, error) {
	_ = ctx
	return "https://conference.local/booking/" + bookingID.String(), nil
}

func (bookingConferenceStub) CancelBookingLink(ctx context.Context, link string) error {
	_ = ctx
	_ = link
	return nil
}

type bookingClock struct{}

func (bookingClock) Now() time.Time { return time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC) }

type bookingTx struct{}

func (bookingTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestCreateBookingHandler(t *testing.T) {
	slotID := uuid.New()
	repo := &bookingRepoHandlerStub{}
	service := bookingservice.New(
		repo,
		bookingSlotRepoHandlerStub{slot: &slotmodel.Slot{
			ID:      slotID,
			StartAt: time.Date(2026, 3, 21, 11, 0, 0, 0, time.UTC),
		}},
		bookingConferenceStub{},
		bookingClock{},
		bookingTx{},
	)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyUserID, authmodel.RoleUser)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(`{"slotId":"`+slotID.String()+`"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestCreateBookingHandlerRejectsEmptySlotID(t *testing.T) {
	repo := &bookingRepoHandlerStub{}
	service := bookingservice.New(
		repo,
		bookingSlotRepoHandlerStub{},
		bookingConferenceStub{},
		bookingClock{},
		bookingTx{},
	)
	handler := New(service)
	jwtManager := jwtplatform.New("secret", time.Hour)
	token, err := jwtManager.Issue(authmodel.DummyUserID, authmodel.RoleUser)
	if err != nil {
		t.Fatalf("Issue token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", strings.NewReader(`{"slotId":"   "}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	middleware.Authenticate(jwtManager)(handler.Create()).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
