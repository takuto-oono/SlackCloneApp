package controllers

import (
	"backend/models"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	message    chan models.Message
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case m := <-h.message:
			for client := range h.clients {
				if m.ChannelId != client.channelID {
					continue
				}
				select {
				case client.send <- m:
				}

			}
		}
	}
}
