package handler

import (
	"net/http"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

// ListAll godoc
// @Summary List all bookings with pagination
// @Tags Bookings
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} listAllBookingsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bookings/list [get]
func (h *Handler) ListAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if !hasRole(ctx, authmodel.RoleAdmin) {
			httpcommon.WriteForbidden(w, "booking list requires admin role")
			return
		}

		page, pageSize, ok := parsePagination(r)
		if !ok {
			httpcommon.WriteInvalidRequest(w)
			return
		}

		bookings, total, err := h.service.ListAll(ctx, page, pageSize)
		if err != nil {
			httpcommon.WriteInternalError(w)
			return
		}

		response := make([]bookingResponse, 0, len(bookings))
		for _, booking := range bookings {
			response = append(response, toBookingResponse(booking))
		}

		httpcommon.WriteJSON(w, http.StatusOK, listAllBookingsResponse{
			Bookings: response,
			Pagination: paginationResponse{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		})
	})
}
