package handler

import (
	"net/http"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

// List godoc
// @Summary List rooms
// @Tags Rooms
// @Produce json
// @Security BearerAuth
// @Success 200 {object} listRoomsResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms/list [get]
func (h *Handler) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasRole(ctx, authmodel.RoleAdmin) && !hasRole(ctx, authmodel.RoleUser) {
			httpcommon.WriteForbidden(w, "room list requires admin or user role")
			return
		}

		rooms, err := h.service.List(ctx)
		if err != nil {
			httpcommon.WriteInternalError(w)
			return
		}

		response := make([]roomResponse, 0, len(rooms))
		for _, room := range rooms {
			response = append(response, toRoomResponse(room))
		}

		httpcommon.WriteJSON(w, http.StatusOK, listRoomsResponse{Rooms: response})
	})
}
