package ws

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"

	"backend/models"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn

	userID     uint32
	channelIDs []int
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uint32) *Client {
	getChannelIDsByUserID := func(userID uint32) []int {
		// userが所属しているすべてのworkspaceのあるchannelを取得
		result := make([]int, 0)
		waus, err := models.GetWAUsByUserId(db, userID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, wau := range waus {
			channels, err := models.GetChannelsByWorkspaceId(db, wau.WorkspaceId)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, ch := range channels {
				result = append(result, ch.ID)
			}
		}
		return result
	}

	return &Client{
		hub:  hub,
		conn: conn,

		userID:     userID,
		channelIDs: getChannelIDsByUserID(userID),
	}
}

func (c *Client) ReadPump(conn *websocket.Conn) {
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(m))
		if err := NewWsMessage(m).TaskRun(); err != nil {
			fmt.Println(err)
			break
		}
	}
}
