package ws

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 5 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = pongWait * 9 / 10
	sendBuf    = 128
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	topics map[string]struct{}
	hub    *Hub
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func newClient(conn *websocket.Conn, hub *Hub, topics []string) *Client {
	tmap := make(map[string]struct{})
	for _, t := range topics {
		if t != "" {
			tmap[t] = struct{}{}
		}
	}
	return &Client{
		conn:   conn,
		send:   make(chan []byte, sendBuf),
		topics: tmap,
		hub:    hub,
	}
}

func (c *Client) remoteAddr() string {
	if ra := c.conn.RemoteAddr(); ra != nil {
		return ra.(*net.TCPAddr).String()
	}
	return ""
}
