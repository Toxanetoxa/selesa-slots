package ws

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]struct{}
	byTopic    map[string]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
	publish    chan envelope
}

type envelope struct {
	topic string
	data  []byte
}

func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[*Client]struct{}),
		byTopic:    make(map[string]map[*Client]struct{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		publish:    make(chan envelope, 256),
	}

	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			for t := range c.topics {
				h.addToTopic(t, c)
			}
			h.mu.Unlock()

		case c := <-h.unregister:
			h.mu.Lock()
			h.removeClient(c)
			h.mu.Unlock()

		case e := <-h.publish:
			h.mu.RLock()
			for c := range h.byTopic[e.topic] {
				select {
				case c.send <- e.data:
				default: // переполнен — отрубить
					h.mu.RUnlock()
					h.unregister <- c
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) addToTopic(topic string, c *Client) {
	if h.byTopic[topic] == nil {
		h.byTopic[topic] = make(map[*Client]struct{})
	}
	h.byTopic[topic][c] = struct{}{}
}

func (h *Hub) removeClient(c *Client) {
	delete(h.clients, c)
	for t := range h.byTopic {
		delete(h.byTopic[t], c)
		if len(h.byTopic[t]) == 0 {
			delete(h.byTopic, t)
		}
	}
	close(c.send)
}

func (h *Hub) Publish(topic string, v any) {
	raw, _ := json.Marshal(v)
	h.publish <- envelope{topic: topic, data: raw}
}
