package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	scheduleservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/service"
	"github.com/google/uuid"
)

// Create godoc
// @Summary Create room schedule
// @Tags Schedules
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param roomId path string true "Room ID"
// @Param request body createScheduleRequest true "Create schedule payload"
// @Success 201 {object} createScheduleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /rooms/{roomId}/schedule/create [post]
func (h *Handler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasAdminRole(ctx) {
			httpcommon.WriteForbidden(w, "schedule creation requires admin role")
			return
		}

		roomID, err := uuid.Parse(r.PathValue("roomId"))
		if err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		var req createScheduleRequest
		if err := httpcommon.DecodeJSON(r, &req); err != nil {
			httpcommon.WriteInvalidRequest(w)
			return
		}
		if len(req.DaysOfWeek) == 0 || strings.TrimSpace(req.StartTime) == "" || strings.TrimSpace(req.EndTime) == "" {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		schedule, err := h.service.Create(ctx, scheduleservice.CreateInput{
			RoomID:     roomID,
			DaysOfWeek: req.DaysOfWeek,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
		})
		if err != nil {
			switch {
			case errors.Is(err, scheduleservice.ErrRoomNotFound):
				httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeRoomNotFound, "room not found")
			case errors.Is(err, scheduleservice.ErrInvalidSchedule):
				httpcommon.WriteInvalidRequest(w)
			case errors.Is(err, scheduleservice.ErrScheduleAlreadyExists):
				httpcommon.WriteError(w, http.StatusConflict, httpcommon.CodeScheduleExists, "schedule for this room already exists and cannot be changed")
			default:
				httpcommon.WriteInternalError(w)
			}
			return
		}

		httpcommon.WriteJSON(w, http.StatusCreated, createScheduleResponse{Schedule: toScheduleResponse(*schedule)})
	})
}

func hasAdminRole(ctx context.Context) bool {
	claims, ok := middleware.ClaimsFromContext(ctx)
	return ok && claims.Role == string(authmodel.RoleAdmin)
}
