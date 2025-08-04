package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/toxanetoxa/selesa-slots/internal/leaderboard"
)

type LbHandler struct{ svc *leaderboard.Service }

func NewLBHandler(s *leaderboard.Service) *LbHandler { return &LbHandler{svc: s} }

func (h *LbHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req leaderboard.Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	req.UpdatedAt = time.Now().UTC()
	h.svc.Publish(req)
	w.WriteHeader(http.StatusAccepted)
}
