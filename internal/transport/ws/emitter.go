package ws

import (
	//"github.com/toxanetoxa/selesa-slots/internal/game"
	//"github.com/toxanetoxa/selesa-slots/internal/leaderboard"
	"github.com/toxanetoxa/selesa-slots/internal/wallet"
)

type Emitter struct {
	hub *Hub
}

func NewEmitter(hub *Hub) *Emitter {
	return &Emitter{hub: hub}
}

func (e *Emitter) Publish(walEvt wallet.Event) {
	e.hub.Publish("wallet", walEvt)
	//e.hub.Publish("game", gameEvt)
	//e.hub.Publish("leaderboard", lbEvt)
}
