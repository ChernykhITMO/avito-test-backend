package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
	"github.com/google/uuid"
)

type claimsContextKey struct{}

func Authenticate(manager *jwtplatform.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				httpcommon.WriteUnauthorized(w, "missing bearer token")
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			if token == header || token == "" {
				httpcommon.WriteUnauthorized(w, "invalid bearer token")
				return
			}

			claims, err := manager.Parse(token)
			if err != nil {
				httpcommon.WriteUnauthorized(w, "invalid bearer token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*jwtplatform.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*jwtplatform.Claims)
	return claims, ok
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}
