package ws

import (
	"gorm.io/gorm"

	"backend/models"
)

var db *gorm.DB
var messageTypeMap map[string]string

func createMessageTypeMap() map[string]string {
	messageTypeMap := make(map[string]string)

	messageTypeMap["channelSendMessage"] = "ws_channel_send_message"
	messageTypeMap["channelReadMessage"] = "ws_channel_read_message"
	messageTypeMap["channelEditMessage"] = "ws_channel_edit_message"

	return messageTypeMap
}

func init() {
	db = models.ExportDB()
	messageTypeMap = createMessageTypeMap()
}
