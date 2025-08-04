package leaderboard

import wstransport "github.com/toxanetoxa/selesa-slots/internal/transport/ws"

type Service struct{ hub *wstransport.Hub }

func NewService(h *wstransport.Hub) *Service { return &Service{hub: h} }

func (s *Service) Publish(e Event) {
	s.hub.Publish("leaderboard", e)
}
