package handler

import roomservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/service"

type Handler struct {
	service *roomservice.Service
}

func New(service *roomservice.Service) *Handler {
	return &Handler{service: service}
}
