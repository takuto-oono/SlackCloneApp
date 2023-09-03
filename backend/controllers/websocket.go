package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"backend/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsController(ctx *gin.Context, h *ws.Hub) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ctx.Header("Access-Control-Allow-Origin", "*")
	userID, err := Authenticate(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println(conn)

	client := ws.NewClient(h, conn, userID)
	go client.ReadPump(conn)

}
