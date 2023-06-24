package controllers

import (
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
type Client struct {
	hub *Hub

	conn *websocket.Conn

	send chan []byte
}
