package handler

import bookingservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/service"

type Handler struct {
	service *bookingservice.Service
}

func New(service *bookingservice.Service) *Handler {
	return &Handler{service: service}
}
