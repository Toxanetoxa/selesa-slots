package test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	httptransport "github.com/toxanetoxa/selesa-slots/internal/transport/http"
	wstransport "github.com/toxanetoxa/selesa-slots/internal/transport/ws"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
)

func TestWebSocketWalletEvent(t *testing.T) {
	hub := wstransport.NewHub()
	emitter := wstransport.NewEmitter(hub)

	repo := wallet.NewMemoryWallet()
	svc := wallet.NewService(repo, emitter)

	root := chi.NewRouter()
	root.Mount("/", httptransport.NewRouter(
		httptransport.NewWalletHandler(svc, zap.NewNop()),
		nil,
		nil,
		zap.NewNop(),
	))
	root.Handle("/ws", wstransport.Handler(hub, zap.NewNop()))

	ts := httptest.NewServer(root)
	defer ts.Close()

	dialURL := "ws" + ts.URL[4:] + "/ws?topics=wallet"
	conn, _, err := websocket.DefaultDialer.Dial(dialURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	// ждём событие
	require.NoError(t, svc.Deposit(context.TODO(), 5, 400))

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	require.NoError(t, err)

	var evt wallet.Event
	require.NoError(t, json.Unmarshal(data, &evt))
	require.Equal(t, wallet.UserID(5), evt.UserID)
	require.Equal(t, wallet.UserID(5), evt.UserID)
	require.Equal(t, wallet.Amount(400), evt.Amount)
}
