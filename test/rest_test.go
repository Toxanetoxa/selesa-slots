package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	httptransport "github.com/toxanetoxa/selesa-slots/internal/transport/http"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
)

func TestDepositWithdrawBalance(t *testing.T) {
	repo := wallet.NewMemoryWallet()
	svc := wallet.NewService(repo, nil)
	router := httptransport.NewRouter(
		httptransport.NewWalletHandler(svc, zap.NewNop()),
		nil,
		nil,
		zap.NewNop(),
	)

	ts := httptest.NewServer(router)
	defer ts.Close()

	type rq struct {
		UserID int64 `json:"user_id"`
		Amount int64 `json:"amount"`
	}

	// 1) deposit
	post(t, ts.URL+"/api/wallet/deposit", rq{1, 700})

	// 2) balance == 700
	require.Equal(t, int64(700), balance(t, ts.URL, 1))

	// 3) withdraw 100
	post(t, ts.URL+"/api/wallet/withdraw", rq{1, 100})

	// 4) balance == 600
	require.Equal(t, int64(600), balance(t, ts.URL, 1))
}

func post(t *testing.T, url string, body any) {
	raw, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewReader(raw))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func balance(t *testing.T, base string, uid int64) int64 {
	resp, err := http.Get(base + "/api/wallet/balance/" + strconv.FormatInt(uid, 10))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out struct{ Balance int64 }
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
	return out.Balance
}
