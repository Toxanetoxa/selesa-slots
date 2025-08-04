package http

import (
	"encoding/json"
	"github.com/toxanetoxa/selesa-slots/internal/game"
	"net/http"
	"time"
)

type GameHandler struct{ svc *game.Service }

func NewGameHandler(s *game.Service) *GameHandler { return &GameHandler{svc: s} }

func (h *GameHandler) Publish(w http.ResponseWriter, r *http.Request) {
	var req game.Event
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	req.Timestamp = time.Now().UTC()
	h.svc.Publish(req)
	w.WriteHeader(http.StatusAccepted)
}
