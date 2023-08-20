package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func readPump(conn *websocket.Conn) {
	for {
		_, m, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(m[:]))
	}
}

func WsController(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conn)
	go readPump(conn)
}