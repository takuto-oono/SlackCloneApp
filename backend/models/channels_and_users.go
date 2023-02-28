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
	cmd := fmt.Sprintf("INSERT INTO %s (channel_id, user_id, is_admin) VALUES ($1, $2, $3)", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, cau.ChannelId, cau.UserId, cau.IsAdmin)
	return err
}

func GetCAUByChannelIdAndUserId(channelId int, userId uint32) (ChannelsAndUsers, error) {
	cmd := fmt.Sprintf("SELECT channel_id, user_id, is_admin FROM %s WHERE channel_id = $1 AND user_id = $2", config.Config.ChannelsAndUserTableName)
	row := DbConnection.QueryRow(cmd, channelId, userId)
	var cau ChannelsAndUsers
	if err := row.Scan(&cau.ChannelId, &cau.UserId, &cau.IsAdmin); err != nil {
		return ChannelsAndUsers{}, err
	}
	return cau, nil
}

func IsExistCAUByChannelIdAndUserId(channelId int, userId uint32) bool {
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE channel_id = $1 AND user_id = $2", config.Config.ChannelsAndUserTableName)
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
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE channel_id = $1 AND user_id = $2 AND is_admin = $3", config.Config.ChannelsAndUserTableName)
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
	cmd := fmt.Sprintf("DELETE FROM %s WHERE channel_id = $1 AND user_id = $2", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, cau.ChannelId, cau.UserId)
	return err
}

func DeleteCAUByChannelId(channelId int) error {
	cmd := fmt.Sprintf("DELETE FROM %s WHERE channel_id = $1", config.Config.ChannelsAndUserTableName)
	_, err := DbConnection.Exec(cmd, channelId)
	return err
}

func GetCAUsByUserId(userId uint32) ([]ChannelsAndUsers, error) {
	caus := make([]ChannelsAndUsers, 0)
	cmd := fmt.Sprintf("SELECT channel_id, user_id, is_admin FROM %s WHERE user_id = ?", config.Config.ChannelsAndUserTableName)
	rows, err := DbConnection.Query(cmd, userId)
	if err != nil {
		return caus, err
	}
	defer rows.Close()
	for rows.Next() {
		var cau ChannelsAndUsers
		if err := rows.Scan(&cau.ChannelId, &cau.UserId, &cau.IsAdmin); err != nil {
			return caus, err
		}
		caus = append(caus, cau)
	}
	return caus, err
}
