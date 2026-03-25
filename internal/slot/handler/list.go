package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	slotservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/service"
	"github.com/google/uuid"
)

// ListAvailable godoc
// @Summary List available slots by room and date
// @Tags Slots
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param date query string true "Date in YYYY-MM-DD"
// @Success 200 {object} listSlotsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms/{roomId}/slots/list [get]
func (h *Handler) ListAvailable() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		roomID, err := uuid.Parse(r.PathValue("roomId"))
		if err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		dateParam := r.URL.Query().Get("date")
		if dateParam == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		date, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		slots, err := h.service.ListAvailable(ctx, roomID, date.UTC())
		if err != nil {
			switch {
			case errors.Is(err, slotservice.ErrRoomNotFound):
				httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeRoomNotFound, "room not found")
			default:
				httpcommon.WriteInternalError(w)
			}
			return
		}

		response := make([]slotResponse, 0, len(slots))
		for _, slot := range slots {
			response = append(response, toSlotResponse(slot))
		}

		httpcommon.WriteJSON(w, http.StatusOK, listSlotsResponse{Slots: response})
	})
}
