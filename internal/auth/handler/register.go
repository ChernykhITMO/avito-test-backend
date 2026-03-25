package handler

import (
	"errors"
	"net/http"
	"strings"

	authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

// Register godoc
// @Summary Register user by email/password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Register payload"
// @Success 201 {object} registerResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func (h *Handler) Register() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req registerRequest
		if err := httpcommon.DecodeJSON(r, &req); err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.Role) == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		user, err := h.service.Register(ctx, authservice.RegisterInput{
			Email:    req.Email,
			Password: req.Password,
			Role:     req.Role,
		})
		if err != nil {
			switch {
			case errors.Is(err, authservice.ErrInvalidRole), errors.Is(err, authservice.ErrInvalidCredentials):
				httpcommon.WriteInvalidRequest(w)
			case errors.Is(err, authservice.ErrEmailAlreadyExists):
				httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.CodeInvalidRequest, "email already exists")
			default:
				httpcommon.WriteInternalError(w)
			}
			return
		}

		httpcommon.WriteJSON(w, http.StatusCreated, registerResponse{User: toUserResponse(*user)})
	})
}
