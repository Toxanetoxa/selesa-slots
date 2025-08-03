package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"

	"go.uber.org/zap"
	"net/http"
)

func NewRouter(h *Handler, log *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(zapLogger(log))

	r.Route("/api/wallet", func(r chi.Router) {
		r.Post("/deposit", h.Deposit)
		r.Post("/withdraw", h.Withdraw)
		r.Get("/balance/{user_id}", h.Balance)
	})

	return r
}
