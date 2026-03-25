package handler

import (
	"context"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
)

func hasRole(ctx context.Context, role authmodel.Role) bool {
	claims, ok := middleware.ClaimsFromContext(ctx)
	return ok && claims.Role == string(role)
}
