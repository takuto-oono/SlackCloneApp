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

func IsExistCAUByChannelIdAndUserId(channelId int, userId uint32) bool {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE channel_id = ? AND user_id = ?", config.Config.ChannelsAndUserTableName)
	rows, err := DbConnection.Query(cmd, channelId, userId)
	if err != nil {
		return false
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt++
	}
	return cnt == 1
}

func IsAdminUserInChannel(channelId int, userId uint32) bool {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE channel_id = ? AND user_id = ? AND is_admin = ?", config.Config.ChannelsAndUserTableName)
	rows, err := DbConnection.Query(cmd, channelId, userId, true)
	if err != nil {
		return false
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		cnt++
	}
	return cnt == 1
}

func (cau *ChannelsAndUsers) DeleteUserFromChannel() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE channel_id = ? AND user_id = ?", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, cau.ChannelId, cau.UserId)
	return err
}