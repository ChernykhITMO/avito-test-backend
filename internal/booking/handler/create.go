package handler

import (
	"errors"
	"net/http"
	"strings"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	bookingservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/service"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	"github.com/google/uuid"
)

// Create godoc
// @Summary Create booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createBookingRequest true "Create booking payload"
// @Success 201 {object} createBookingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bookings/create [post]
func (h *Handler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasRole(ctx, authmodel.RoleUser) {
			httpcommon.WriteForbidden(w, "booking creation requires user role")
			return
		}

		userID, ok := userIDFromContext(ctx)
		if !ok {
			httpcommon.WriteUnauthorized(w, "authenticated user id is missing from token")
			return
		}

		var req createBookingRequest
		if err := httpcommon.DecodeJSON(r, &req); err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if strings.TrimSpace(req.SlotID) == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		slotID, err := uuid.Parse(req.SlotID)
		if err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		booking, err := h.service.Create(ctx, bookingservice.CreateInput{
			SlotID:               slotID,
			UserID:               userID,
			CreateConferenceLink: req.CreateConferenceLink,
		})
		if err != nil {
			switch {
			case errors.Is(err, bookingservice.ErrSlotNotFound):
				httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeSlotNotFound, "slot not found")
			case errors.Is(err, bookingservice.ErrSlotInPast):
				httpcommon.WriteError(w, http.StatusBadRequest, httpcommon.CodeInvalidRequest, "slot is in the past")
			case errors.Is(err, bookingservice.ErrSlotAlreadyBooked):
				httpcommon.WriteError(w, http.StatusConflict, httpcommon.CodeSlotAlreadyBooked, "slot is already booked")
			default:
				httpcommon.WriteInternalError(w)
			}
			return
		}

		httpcommon.WriteJSON(w, http.StatusCreated, createBookingResponse{Booking: toBookingResponse(*booking)})
	})
}
