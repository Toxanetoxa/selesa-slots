package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/toxanetoxa/selesa-slots/internal/game"
	"github.com/toxanetoxa/selesa-slots/internal/leaderboard"
	grpcWallet "github.com/toxanetoxa/selesa-slots/internal/transport/grpc"
	httptransport "github.com/toxanetoxa/selesa-slots/internal/transport/http"
	wstransport "github.com/toxanetoxa/selesa-slots/internal/transport/ws"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      root,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	grpcSrv := grpcWallet.NewServer(walletSvc, log)
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatal("grpc listen", zap.Error(err))
	}

	var g errgroup.Group

	g.Go(func() error {
		log.Info("HTTP/WS up", zap.String("addr", httpSrv.Addr))
		return httpSrv.ListenAndServe()
	})

	g.Go(func() error {
		log.Info("gRPC up", zap.String("addr", lis.Addr().String()))
		return grpcSrv.Serve(lis)
	})

	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Info("shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		grpcSrv.GracefulStop()
		_ = httpSrv.Shutdown(ctx)
		return nil
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server error", zap.Error(err))
	}
}
