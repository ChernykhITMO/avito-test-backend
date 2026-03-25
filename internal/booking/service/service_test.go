package service

import (
	"context"
	"errors"
	"testing"
	"time"

	bookingmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/model"
	slotmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/model"
	"github.com/google/uuid"
)

type bookingRepoStub struct {
	createErr        error
	created          *bookingmodel.Booking
	getForUpdateResp *bookingmodel.Booking
	cancelled        bool
}

func (s *bookingRepoStub) Create(ctx context.Context, booking bookingmodel.Booking) error {
	s.created = &booking
	return s.createErr
}

func (s *bookingRepoStub) GetByID(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	return nil, nil
}

func (s *bookingRepoStub) GetByIDForUpdate(ctx context.Context, bookingID uuid.UUID) (*bookingmodel.Booking, error) {
	return s.getForUpdateResp, nil
}

func (s *bookingRepoStub) Cancel(ctx context.Context, bookingID uuid.UUID, cancelledAt time.Time) error {
	s.cancelled = true
	return nil
}

func (s *bookingRepoStub) ListMyFuture(ctx context.Context, userID uuid.UUID, now time.Time) ([]bookingmodel.Booking, error) {
	return []bookingmodel.Booking{}, nil
}

func (s *bookingRepoStub) ListAll(ctx context.Context, page, pageSize int) ([]bookingmodel.Booking, int, error) {
	return []bookingmodel.Booking{}, 0, nil
}

type bookingSlotRepoStub struct {
	slot *slotmodel.Slot
}

func (s bookingSlotRepoStub) GetByID(ctx context.Context, slotID uuid.UUID) (*slotmodel.Slot, error) {
	return s.slot, nil
}

type conferenceStub struct {
	link         string
	createErr    error
	cancelErr    error
	cancelCalled bool
	cancelLink   string
}

func (s *conferenceStub) CreateBookingLink(ctx context.Context, bookingID uuid.UUID) (string, error) {
	_ = ctx
	if s.createErr != nil {
		return "", s.createErr
	}
	if s.link != "" {
		return s.link, nil
	}
	return "https://conference.local/booking/" + bookingID.String(), nil
}

func (s *conferenceStub) CancelBookingLink(ctx context.Context, link string) error {
	_ = ctx
	s.cancelCalled = true
	s.cancelLink = link
	return s.cancelErr
}

type noopTransactor struct{}

func (noopTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

func TestCreateBookingPastSlot(t *testing.T) {
	service := New(
		&bookingRepoStub{},
		bookingSlotRepoStub{
			slot: &slotmodel.Slot{
				ID:      uuid.New(),
				StartAt: time.Date(2026, 3, 20, 9, 0, 0, 0, time.UTC),
			},
		},
		&conferenceStub{},
		fixedClock{now: time.Date(2026, 3, 21, 9, 0, 0, 0, time.UTC)},
		noopTransactor{},
	)

	_, err := service.Create(context.Background(), CreateInput{
		SlotID: uuid.New(),
		UserID: uuid.New(),
	})
	if !errors.Is(err, ErrSlotInPast) {
		t.Fatalf("expected slot in past, got %v", err)
	}
}

func TestCancelBookingForbidden(t *testing.T) {
	ownerID := uuid.New()
	service := New(
		&bookingRepoStub{
			getForUpdateResp: &bookingmodel.Booking{
				ID:     uuid.New(),
				UserID: ownerID,
				Status: bookingmodel.StatusActive,
			},
		},
		bookingSlotRepoStub{},
		&conferenceStub{},
		fixedClock{now: time.Now().UTC()},
		noopTransactor{},
	)

	_, err := service.Cancel(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, ErrForbiddenCancel) {
		t.Fatalf("expected forbidden cancel, got %v", err)
	}
}

func TestCreateBookingConferenceUnavailableDoesNotFailBooking(t *testing.T) {
	repo := &bookingRepoStub{}
	service := New(
		repo,
		bookingSlotRepoStub{
			slot: &slotmodel.Slot{
				ID:      uuid.New(),
				StartAt: time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC),
			},
		},
		&conferenceStub{createErr: errors.New("conference unavailable")},
		fixedClock{now: time.Date(2026, 3, 21, 9, 0, 0, 0, time.UTC)},
		noopTransactor{},
	)

	created, err := service.Create(context.Background(), CreateInput{
		SlotID:               uuid.New(),
		UserID:               uuid.New(),
		CreateConferenceLink: true,
	})
	if err != nil {
		t.Fatalf("expected booking creation without conference link, got err %v", err)
	}
	if created == nil {
		t.Fatal("expected booking to be returned")
	}
	if created.ConferenceLink != nil {
		t.Fatal("expected conference link to be nil")
	}
	if repo.created == nil {
		t.Fatal("expected booking to be written")
	}
}

func TestCreateBookingCompensatesConferenceLinkOnInsertError(t *testing.T) {
	repo := &bookingRepoStub{createErr: bookingmodel.ErrSlotAlreadyBooked}
	conference := &conferenceStub{link: "https://conference.local/booking/test"}
	service := New(
		repo,
		bookingSlotRepoStub{
			slot: &slotmodel.Slot{
				ID:      uuid.New(),
				StartAt: time.Date(2026, 3, 21, 10, 0, 0, 0, time.UTC),
			},
		},
		conference,
		fixedClock{now: time.Date(2026, 3, 21, 9, 0, 0, 0, time.UTC)},
		noopTransactor{},
	)

	_, err := service.Create(context.Background(), CreateInput{
		SlotID:               uuid.New(),
		UserID:               uuid.New(),
		CreateConferenceLink: true,
	})
	if !errors.Is(err, ErrSlotAlreadyBooked) {
		t.Fatalf("expected ErrSlotAlreadyBooked, got %v", err)
	}
	if !conference.cancelCalled {
		t.Fatal("expected compensation cancel call")
	}
	if conference.cancelLink == "" {
		t.Fatal("expected cancel link to be passed")
	}
}
