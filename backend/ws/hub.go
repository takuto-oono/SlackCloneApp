package ws

import (
	"fmt"
	"os"

	"backend/models"
)

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client

	channelIDs []int
}

func NewHub() *Hub {
	getAllChannelIDs := func() []int {
		result := make([]int, 0)
		channels, err := models.GetAllChannels(db)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, ch := range channels {
			result = append(result, ch.ID)
		}
		return result
	}

	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		channelIDs: getAllChannelIDs(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
			}
		}
	}
}
