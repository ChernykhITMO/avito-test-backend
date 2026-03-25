package handler

import (
	"errors"
	"net/http"
	"strings"

	authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

// Login godoc
// @Summary Login by email/password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login payload"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func (h *Handler) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req loginRequest
		if err := httpcommon.DecodeJSON(r, &req); err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		token, err := h.service.Login(ctx, authservice.LoginInput{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			switch {
			case errors.Is(err, authservice.ErrInvalidCredentials):
				httpcommon.WriteUnauthorized(w, "invalid credentials")
			default:
				httpcommon.WriteInternalError(w)
			}
			return
		}

		httpcommon.WriteJSON(w, http.StatusOK, tokenResponse{Token: token})
	})
}
