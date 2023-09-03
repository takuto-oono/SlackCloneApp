package controllers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"
)

func TestWebsocket(t *testing.T) {
	rr, u := SignUpTestFuncV2(randomstring.EnglishFrequencyString(30), "abc123")
	assert.Equal(t, http.StatusOK, rr.Code)
	rr, _ = LoginTestFuncV2(u.Name, u.PassWord)
	assert.Equal(t, http.StatusOK, rr.Code)

	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/websocket/", nil)
	if err != nil {
		fmt.Println(err)
	}
	defer ws.Close()
}
