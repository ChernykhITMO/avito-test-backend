package handler

import slotservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/service"

type Handler struct {
	service *slotservice.Service
}

func New(service *slotservice.Service) *Handler {
	return &Handler{service: service}
}
