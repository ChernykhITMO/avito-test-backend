package middleware

import (
	"log/slog"
	"net/http"

	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
)

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(
						"panic recovered",
						slog.String("op", "internal.platform.middleware.Recover"),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
						slog.Any("panic", rec),
					)
					httpcommon.WriteInternalError(w)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
