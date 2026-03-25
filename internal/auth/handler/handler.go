package handler

import authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"

type Handler struct {
	service *authservice.Service
}

func New(service *authservice.Service) *Handler {
	return &Handler{service: service}
}
