package handler

import (
	"context"
	"net/http"
	"strconv"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	"github.com/google/uuid"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

func hasRole(ctx context.Context, role authmodel.Role) bool {
	claims, ok := middleware.ClaimsFromContext(ctx)
	return ok && claims.Role == string(role)
}

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	return middleware.UserIDFromContext(ctx)
}

func parsePagination(r *http.Request) (int, int, bool) {
	page := defaultPage
	pageSize := defaultPageSize

	if raw := r.URL.Query().Get("page"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < defaultPage {
			return 0, 0, false
		}
		page = value
	}

	if raw := r.URL.Query().Get("pageSize"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < defaultPage || value > maxPageSize {
			return 0, 0, false
		}
		pageSize = value
	}

	return page, pageSize, true
}
