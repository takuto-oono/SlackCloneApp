package controllers

import (
	"backend/controllerUtils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var messageRouter = SetupRouter()

func SendMessageTestFunc(text string, channelId int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.SendMessageInput{
		Text:      text,
		ChannelId: channelId,
	})
	req, err := http.NewRequest("POST", "/api/message/send", bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	messageRouter.ServeHTTP(rr, req)
	return rr
}

func TestSendMessage(t *testing.T) {
	// TODO test
	assert.Equal(t, "yes", "yes")
}
