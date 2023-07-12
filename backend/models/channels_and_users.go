package models

import (
	"gorm.io/gorm"
)

type ChannelsAndUsers struct {
	ChannelId int    `json:"channel_id" gorm:"primaryKey"`
	UserId    uint32 `json:"user_id" gorm:"primaryKey"`
	IsAdmin   bool   `json:"is_admin" gorm:"not null"`
}

func NewChannelsAndUses(channelId int, userId uint32, isAdmin bool) *ChannelsAndUsers {
	return &ChannelsAndUsers{
		ChannelId: channelId,
		UserId:    userId,
		IsAdmin:   isAdmin,
	}
}

func (cau *ChannelsAndUsers) Create(tx *gorm.DB) error {
	return tx.Model(&ChannelsAndUsers{}).Create(cau).Error
}

func GetCAUByChannelIdAndUserId(tx *gorm.DB, channelId int, userId uint32) (ChannelsAndUsers, error) {
	var result ChannelsAndUsers
	err := tx.Model(&ChannelsAndUsers{}).Where("channel_id = ? AND user_id = ?", channelId, userId).Take(&result).Error
	return result, err
}

func (cau ChannelsAndUsers) Delete(tx *gorm.DB) error {
	return tx.Where("channel_id = ? AND user_id = ?", cau.ChannelId, cau.UserId).Delete(&ChannelsAndUsers{}).Error
}

func DeleteCAUsByChannelId(tx *gorm.DB, channelId int) error {
	return tx.Where("channel_id = ?", channelId).Delete(&ChannelsAndUsers{}).Error
}

func GetCAUsByUserId(tx *gorm.DB, userId uint32) ([]ChannelsAndUsers, error) {
	var result []ChannelsAndUsers
	rows, err := tx.Model(&ChannelsAndUsers{}).Where("user_id = ?", userId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var cau ChannelsAndUsers
		if err := tx.ScanRows(rows, &cau); err != nil {
			return result, err
		}
		result = append(result, cau)
	}
	return result, nil
}

func GetCAUsByChannelId(tx *gorm.DB, channelId int) ([]ChannelsAndUsers, error) {
	var result []ChannelsAndUsers
	rows, err := tx.Model(&ChannelsAndUsers{}).Where("channel_id = ?", channelId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var cau ChannelsAndUsers
		if err := tx.ScanRows(rows, &cau); err != nil {
			return result, err
		}
		result = append(result, cau)
	}
	return result, nil
}

func DeleteChannelsAndUsersTableRecords() {
	db.Exec("DELETE FROM channels_and_users")
}
