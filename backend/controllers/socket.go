package controllers

import (
	"fmt"

	"github.com/gorilla/websocket"

	"backend/controllerUtils"
	"backend/models"
	"backend/utils"
)

var upgrader = websocket.Upgrader{}

type Client struct {
	hub       *Hub
	channelID int
	conn      *websocket.Conn
	userID    uint32
	send      chan models.Message
}

func (c *Client) readPump(userID uint32) {
	for {
		fmt.Println("in readPump for")
		fmt.Println(c.conn);
		mt, message, err := c.conn.ReadMessage()
		fmt.Println("after readMessage");
		fmt.Println(c.conn);
		if err != nil {
			fmt.Println(err)
			fmt.Println("=================")
			fmt.Println(mt)
			fmt.Println(message)
			break
		}
		in := utils.SendMessageInputToByte(message)
		// message structを作成
		m := models.NewChannelMessage(in.Text, in.ChannelId, userID)
		fmt.Println("================")
		fmt.Println(m)
		fmt.Println("================")
		// userとchannelが同じworkspaceに存在しているかを確認
		if b, err := controllerUtils.IsExistChannelAndUserInSameWorkspace(m.ChannelId, userID); !b || err != nil {
			break
		}

		// channelにuserが参加しているかを確認
		if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(m.ChannelId, userID); !b || err != nil {
			break
		}

		if err = controllerUtils.SendMessageTX(m, userID, in.MentionedUserIDs); err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			
			w.Write(utils.ByteToStruct(message))
		}
	}
}
