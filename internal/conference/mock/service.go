package mock

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrUnavailable = errors.New("conference service unavailable")
)

type Service struct {
	baseURL string
	mode    string
}

func New(baseURL, mode string) *Service {
	return &Service{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/") + "/",
		mode:    strings.TrimSpace(mode),
	}
}

func (s *Service) CreateBookingLink(ctx context.Context, bookingID uuid.UUID) (string, error) {
	_ = ctx

	if s.mode == "create_unavailable" {
		return "", ErrUnavailable
	}

	return s.baseURL + bookingID.String(), nil
}

func (s *Service) CancelBookingLink(ctx context.Context, link string) error {
	_ = ctx
	_ = link

	if s.mode == "cancel_unavailable" {
		return ErrUnavailable
	}

	return nil
}
