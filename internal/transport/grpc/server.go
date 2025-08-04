package grpcwallet

import (
	"context"

	"github.com/toxanetoxa/selesa-slots/internal/wallet"
	walletv1 "github.com/toxanetoxa/selesa-slots/pkg/api/wallet/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewServer(svc *wallet.Service, log *zap.Logger) *grpc.Server {
	gs := grpc.NewServer()
	walletv1.RegisterWalletServiceServer(gs, &handler{svc: svc})
	reflection.Register(gs)
	return gs
}

type handler struct {
	walletv1.UnimplementedWalletServiceServer
	svc *wallet.Service
}

func (h *handler) GetBalance(ctx context.Context,
	req *walletv1.GetBalanceRequest,
) (*walletv1.GetBalanceResponse, error) {

	bal, err := h.svc.GetBalance(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &walletv1.GetBalanceResponse{
		UserId:  req.GetUserId(),
		Balance: bal,
	}, nil
}
