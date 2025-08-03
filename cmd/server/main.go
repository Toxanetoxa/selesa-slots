package main

import (
	"context"
	"errors"
	httptransport "github.com/toxanetoxa/selesa-slots/internal/transport/http"
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

	repo := wallet.NewMemoryWallet()
	svc := wallet.NewService(repo, nil)
	h := httptransport.NewHandler(svc, log)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      httptransport.NewRouter(h, log),
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
