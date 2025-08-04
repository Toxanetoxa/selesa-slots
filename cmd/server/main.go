package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/toxanetoxa/selesa-slots/internal/game"
	"github.com/toxanetoxa/selesa-slots/internal/leaderboard"
	httptransport "github.com/toxanetoxa/selesa-slots/internal/transport/http"
	wstransport "github.com/toxanetoxa/selesa-slots/internal/transport/ws"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	log, _ := zap.NewDevelopment()
	defer log.Sync()

	hub := wstransport.NewHub()
	emitter := wstransport.NewEmitter(hub)

	repo := wallet.NewMemoryWallet()
	walletSvc := wallet.NewService(repo, emitter)
	gameSvc := game.NewService(hub)
	lbSvc := leaderboard.NewService(hub)

	hWallet := httptransport.NewWalletHandler(walletSvc, log)
	hGame := httptransport.NewGameHandler(gameSvc)
	hLb := httptransport.NewLBHandler(lbSvc)
	wsHandler := wstransport.Handler(hub, log)

	router := httptransport.NewRouter(hWallet, hGame, hLb, log)

	root := chi.NewRouter()
	root.Mount("/", router)
	root.Handle("/ws", wsHandler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      root,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("starting server", zap.String("addr", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
