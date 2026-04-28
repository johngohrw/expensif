package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// RecoverPanic recovers from panics in HTTP handlers, logs them, and returns a 500.
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("panic recovered",
					"error", fmt.Sprintf("%v", err),
					"path", r.URL.Path,
					"method", r.Method,
					"stack", string(debug.Stack()),
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// responseWriter captures the status code written by the handler.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLog logs every HTTP request with method, path, status, and duration.
func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", time.Since(start),
		)
	})
}
