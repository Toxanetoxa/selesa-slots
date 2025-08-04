package ws

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, //FIXME CORS для всех
}

func Handler(hub *Hub, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warn("upgrade failed", zap.Error(err))
			return
		}

		topics := strings.Split(r.URL.Query().Get("topics"), ",")
		cli := newClient(conn, hub, topics)
		hub.register <- cli

		log.Info("ws connected",
			zap.String("addr", cli.remoteAddr()),
			zap.Strings("topics", topics),
		)

		go cli.writePump()
		go cli.readPump()
	}
}
