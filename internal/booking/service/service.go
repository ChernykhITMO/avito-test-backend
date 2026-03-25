package service

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	bookingRepository BookingRepository
	slotRepository    SlotRepository
	conferenceService ConferenceService
	clock             Clock
	transactor        Transactor
}

type noopConferenceService struct{}

func (noopConferenceService) CreateBookingLink(_ context.Context, _ uuid.UUID) (string, error) {
	return "", nil
}

func (noopConferenceService) CancelBookingLink(_ context.Context, _ string) error {
	return nil
}

func New(
	bookingRepository BookingRepository,
	slotRepository SlotRepository,
	conferenceService ConferenceService,
	clock Clock,
	transactor Transactor,
) *Service {
	if conferenceService == nil {
		conferenceService = noopConferenceService{}
	}

	return &Service{
		bookingRepository: bookingRepository,
		slotRepository:    slotRepository,
		conferenceService: conferenceService,
		clock:             clock,
		transactor:        transactor,
	}
}
