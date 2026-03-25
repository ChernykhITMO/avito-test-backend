package handler

import scheduleservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/service"

type Handler struct {
	service *scheduleservice.Service
}

func New(service *scheduleservice.Service) *Handler {
	return &Handler{service: service}
}
