package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"

	"go.uber.org/zap"
	"net/http"
)

func NewRouter(
	walletH *WalletHandler,
	gameH *GameHandler,
	lbH *LbHandler,
	log *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(zapLogger(log))

	r.Route("/api/wallet", func(r chi.Router) {
		r.Post("/deposit", walletH.Deposit)
		r.Post("/withdraw", walletH.Withdraw)
		r.Get("/balance/{user_id}", walletH.Balance)
	})

	r.Route("/api/game", func(r chi.Router) {
		r.Post("/outcome", gameH.Publish)
	})
	r.Route("/api/leaderboard", func(r chi.Router) {
		r.Post("/update", lbH.Publish)
	})

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}
