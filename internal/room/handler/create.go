package handler

import (
	"errors"
	"net/http"
	"strings"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	roomservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/service"
)

// Create godoc
// @Summary Create room
// @Tags Rooms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createRoomRequest true "Create room payload"
// @Success 201 {object} createRoomResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms/create [post]
func (h *Handler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasRole(ctx, authmodel.RoleAdmin) {
			httpcommon.WriteForbidden(w, "room creation requires admin role")
			return
		}

		var req createRoomRequest
		if err := httpcommon.DecodeJSON(r, &req); err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if strings.TrimSpace(req.Name) == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if req.Capacity != nil && *req.Capacity <= 0 {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		room, err := h.service.Create(ctx, roomservice.CreateInput{
			Name:        req.Name,
			Description: req.Description,
			Capacity:    req.Capacity,
		})
		if err != nil {
			if errors.Is(err, roomservice.ErrInvalidRoom) {
				httpcommon.WriteInvalidRequest(w)
				return
			}

			httpcommon.WriteInternalError(w)
			return
		}

		httpcommon.WriteJSON(w, http.StatusCreated, createRoomResponse{Room: toRoomResponse(*room)})
	})
}
