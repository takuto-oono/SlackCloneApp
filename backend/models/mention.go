package models

import (
	"time"

	"gorm.io/gorm"
)

type Mention struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint32    `json:"user_id" gorm:"not null column:user_id"`
	MessageID uint      `json:"message_id" gorm:"not null column:message_id"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

func NewMention(userID uint32, messageID uint) *Mention {
	return &Mention{
		UserID:    userID,
		MessageID: messageID,
	}
}

func (men *Mention) Create(tx *gorm.DB) error {
	return tx.Model(&Mention{}).Create(men).Error
}

func GetMentionsByMessageID(tx *gorm.DB, messageID uint) ([]Mention, error) {
	var result []Mention
	rows, err := tx.Model(&Mention{}).Where("message_id = ?", messageID).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var men Mention
		if err := tx.ScanRows(rows, &men); err != nil {	
			return result, err
		}
		result = append(result, men)
	}
	return result, nil
}

func GetMentionsByUserIDSortedByCreatedAt(tx *gorm.DB, userID uint32) ([]Mention, error) {
	var result []Mention
	rows, err := tx.Model(&Mention{}).Where("user_id = ?", userID).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var men Mention
		if err := tx.ScanRows(rows, &men); err != nil {
			return result, err
		}
		result = append(result, men)
	}
	return result, nil
}
