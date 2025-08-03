package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func zapLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			next.ServeHTTP(ww, r)
			duration := time.Since(start)

			reqID := middleware.GetReqID(r.Context())

			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("request_id", reqID),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Duration("duration", duration),
				zap.String("ip", r.RemoteAddr),
				zap.String("user-agent", r.UserAgent()),
				zap.String("error", http.StatusText(ww.Status())),
			)
		})
	}
}
