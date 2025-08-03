package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Handler struct {
	svc *wallet.Service
	log *zap.Logger
}

func NewHandler(svc *wallet.Service, log *zap.Logger) *Handler {
	return &Handler{
		svc: svc,
		log: log,
	}
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	const op = `transport.http.handler.Deposit`

	var req amountReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Bind() != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Deposit(r.Context(), req.UserId, req.Amount); err != nil {
		HTTPError(w, err)
		return
	}

	h.respondBalance(w, r, req.UserId)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	const op = `transport.http.handler.Withdraw`

	var req amountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Bind() != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Withdraw(r.Context(), req.UserId, req.Amount); err != nil {
		HTTPError(w, err)
		return
	}

	h.respondBalance(w, r, req.UserId)
}

func (h *Handler) Balance(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	h.respondBalance(w, r, userID)
}

func (h *Handler) respondBalance(w http.ResponseWriter, r *http.Request, userID int64) {
	balance, err := h.svc.GetBalance(r.Context(), userID)
	if err != nil {
		HTTPError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(balanceResp{userID, balance})
	if err != nil {
		HTTPError(w, err)
		return
	}
}
