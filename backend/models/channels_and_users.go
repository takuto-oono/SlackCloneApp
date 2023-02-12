package models

import (
	"fmt"

	"backend/config"
)

type ChannelsAndUsers struct {
	ChannelId int
	UserId    uint32
	IsAdmin   bool
}

func NewChannelsAndUses(channelId int, userId uint32, isAdmin bool) *ChannelsAndUsers {
	return &ChannelsAndUsers{
		ChannelId: channelId,
		UserId:    userId,
		IsAdmin:   isAdmin,
	}
}

func (cau *ChannelsAndUsers) CreateChannelAndUsers() error {
	cmd := fmt.Sprintf("INSERT INTO %s (channel_id, user_id, is_admin) VALUES (?, ?, ?)", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, cau.ChannelId, cau.UserId, cau.IsAdmin)
	return err
}
