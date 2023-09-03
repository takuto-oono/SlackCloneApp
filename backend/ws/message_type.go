package ws

import (
	"encoding/json"
	"fmt"

	"backend/controllerUtils"
)

var notFoundError error = fmt.Errorf("not found task func")

type WsMessage struct {
	messageType string      `json:"message_type"`
	body        interface{} `json:"body"`
}

func NewWsMessage(message []byte) *WsMessage {
	var wm WsMessage
	json.Unmarshal(([]byte)(message), &wm)
	for k, v := range messageTypeMap {
		if wm.messageType == v {
			wm.messageType = k
			break
		}
	}
	return &wm
}

func (wm *WsMessage) GetMessageType() string {
	return wm.messageType
}

func (wm *WsMessage) GetBody() interface{} {
	return wm.body
}

func (wm *WsMessage) TaskRun() error {
	switch wm.GetMessageType() {
	case "channelSendMessage":
		input, ok := wm.GetBody().(controllerUtils.SendMessageInput)
		if !ok {
			return notFoundError
		}
		return SendMessage(input)

	case "channelReadMessage":

	case "channelEditMessage":
	}

	return fmt.Errorf("%s not found messageType", wm.GetMessageType())
}
