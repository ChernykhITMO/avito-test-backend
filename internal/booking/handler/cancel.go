package handler

import (
	"errors"
	"net/http"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	bookingservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/service"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	"github.com/google/uuid"
)

// Cancel godoc
// @Summary Cancel own booking
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param bookingId path string true "Booking ID"
// @Success 200 {object} cancelBookingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bookings/{bookingId}/cancel [post]
func (h *Handler) Cancel() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasRole(ctx, authmodel.RoleUser) {
			httpcommon.WriteForbidden(w, "booking cancellation requires user role")
			return
		}

		userID, ok := userIDFromContext(ctx)
		if !ok {
			httpcommon.WriteUnauthorized(w, "authenticated user id is missing from token")
			return
		}

		bookingID, err := uuid.Parse(r.PathValue("bookingId"))
		if err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		booking, err := h.service.Cancel(ctx, bookingID, userID)
		if err != nil {
			switch {
			case errors.Is(err, bookingservice.ErrBookingNotFound):
				httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeBookingNotFound, "booking not found")
			case errors.Is(err, bookingservice.ErrForbiddenCancel):
				httpcommon.WriteForbidden(w, "cannot cancel another user's booking")
			default:
				httpcommon.WriteInternalError(w)
			}
			return
		}

		httpcommon.WriteJSON(w, http.StatusOK, cancelBookingResponse{Booking: toBookingResponse(*booking)})
	})
}
