package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "internal.platform.middleware.Logging"

			start := time.Now()
			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			attrs := []any{
				slog.String("op", op),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.Int("status", rw.status),
				slog.Duration("duration", time.Since(start)),
			}

			switch {
			case rw.status >= http.StatusInternalServerError:
				logger.Error("request handled", attrs...)
			case rw.status >= http.StatusBadRequest:
				logger.Warn("request handled", attrs...)
			default:
				logger.Info("request handled", attrs...)
			}
		})
	}
}
