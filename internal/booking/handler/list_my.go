package handler

import (
	"net/http"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

// ListMyFuture godoc
// @Summary List current user future bookings
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Success 200 {object} listBookingsResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bookings/my [get]
func (h *Handler) ListMyFuture() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasRole(ctx, authmodel.RoleUser) {
			httpcommon.WriteForbidden(w, "my bookings require user role")
			return
		}

		userID, ok := userIDFromContext(ctx)
		if !ok {
			httpcommon.WriteUnauthorized(w, "authenticated user id is missing from token")
			return
		}

		bookings, err := h.service.ListMyFuture(ctx, userID)
		if err != nil {
			httpcommon.WriteInternalError(w)
			return
		}

		response := make([]bookingResponse, 0, len(bookings))
		for _, booking := range bookings {
			response = append(response, toBookingResponse(booking))
		}

		httpcommon.WriteJSON(w, http.StatusOK, listBookingsResponse{Bookings: response})
	})
}
