package handler

import (
	"errors"
	"net/http"
	"strings"

	authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

// DummyLogin godoc
// @Summary Get test JWT by role
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dummyLoginRequest true "Role payload"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dummyLogin [post]
func (h *Handler) DummyLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req dummyLoginRequest
		if err := httpcommon.DecodeJSON(r, &req); err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if strings.TrimSpace(req.Role) == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		token, err := h.service.DummyLogin(ctx, req.Role)
		if err != nil {
			if errors.Is(err, authservice.ErrInvalidRole) {
				httpcommon.WriteInvalidRequest(w)
				return
			}

			httpcommon.WriteInternalError(w)
			return
		}

		httpcommon.WriteJSON(w, http.StatusOK, tokenResponse{Token: token})
	})
}
