package models

import (
	"fmt"

	"backend/config"
)

type ChannelsAndUsers struct {
	ChannelId int    `json:"channel_id"`
	UserId    uint32 `json:"user_id"`
	IsAdmin   bool   `json:"is_admin"`
}

func NewChannelsAndUses(channelId int, userId uint32, isAdmin bool) *ChannelsAndUsers {
	return &ChannelsAndUsers{
		ChannelId: channelId,
		UserId:    userId,
		IsAdmin:   isAdmin,
	}
}

func (cau *ChannelsAndUsers) Create() error {
	cmd := fmt.Sprintf("INSERT INTO %s (channel_id, user_id, is_admin) VALUES (?, ?, ?)", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, cau.ChannelId, cau.UserId, cau.IsAdmin)
	return err
}

func GetCAUByChannelIdAndUserId(channelId int, userId uint32) (ChannelsAndUsers, error) {
	cmd := fmt.Sprintf("SELECT channel_id, user_id, is_admin FROM %s WHERE channel_id = ? AND user_id = ?", config.Config.ChannelsAndUserTableName)
	row := DbConnection.QueryRow(cmd, channelId, userId)
	var cau ChannelsAndUsers
	if err := row.Scan(&cau.ChannelId, &cau.UserId, &cau.IsAdmin); err != nil {
		return ChannelsAndUsers{}, err
	}
	return cau, nil
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

func (cau *ChannelsAndUsers) Delete() error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE channel_id = ? AND user_id = ?", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, cau.ChannelId, cau.UserId)
	return err
}

func DeleteCAUByChannelId(channelId int) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE channel_id = ?", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, channelId)
	return err
}
