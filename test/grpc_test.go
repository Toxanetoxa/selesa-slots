package test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	grpcwallet "github.com/toxanetoxa/selesa-slots/internal/transport/grpc"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
	walletv1 "github.com/toxanetoxa/selesa-slots/pkg/api/wallet/v1"
)

func TestGetBalance(t *testing.T) {
	lis := bufconn.Listen(1024 * 64)

	// server
	repo := wallet.NewMemoryWallet()
	repo.Deposit(context.TODO(), 3, 900)
	svc := wallet.NewService(repo, nil)

	s := grpcwallet.NewServer(svc, zap.NewNop())
	go s.Serve(lis)
	t.Cleanup(s.GracefulStop)

	// client over bufconn
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "buf",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	require.NoError(t, err)
	defer conn.Close()

	client := walletv1.NewWalletServiceClient(conn)
	res, err := client.GetBalance(ctx, &walletv1.GetBalanceRequest{UserId: 3})
	require.NoError(t, err)
	require.Equal(t, int64(900), res.Balance)
}
