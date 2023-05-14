package models

import (
	"gorm.io/gorm"
)

type MessageAndUser struct {
	MessageID uint   `json:"message_id" gorm:"primaryKey column:message_id"`
	UserID    uint32 `json:"user_id" gorm:"primaryKey column:user_id"`
	IsRead    bool   `json:"is_read" gorm:"not null default:false"`
}

func NewMessageAndUser(messageID uint, userID uint32, isRead bool) *MessageAndUser {
	return &MessageAndUser{
		MessageID: messageID,
		UserID:    userID,
		IsRead:    isRead,
	}
}

func (mau *MessageAndUser) Create(tx *gorm.DB) error {
	return tx.Create(mau).Error
}

func GetMAUByUserIDAndIsRead(tx *gorm.DB, userID uint32, isRead bool) ([]MessageAndUser, error) {
	var result []MessageAndUser
	rows, err := tx.Model(&MessageAndUser{}).Where("user_id = ? AND is_read = ?", userID, isRead).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var mau MessageAndUser
		if err := tx.ScanRows(rows, &mau); err != nil {
			return result, err
		}
		result = append(result, mau)
	}
	return result, nil
}

func (mau *MessageAndUser) UpdateIsRead(tx *gorm.DB, newIsRead bool) error {
	return tx.Model(&MessageAndUser{}).Where("message_id = ? AND user_id = ?", mau.MessageID, mau.UserID).Update("is_read", newIsRead).Take(mau).Error
}
