package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"not null"`
	Date      string    `json:"date" gorm:"not null"`
	ChannelId int       `json:"channel_id" gorm:"not null"`
	UserId    uint32    `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func NewMessage(text string, channelId int, userId uint32) *Message {
	return &Message{
		Text:      text,
		ChannelId: channelId,
		UserId:    userId,
	}
}

func (m *Message) Create(tx *gorm.DB) error {
	return tx.Model(&Message{}).Create(m).Error
}

func GetMessagesByChannelId(tx *gorm.DB, channelId int) ([]Message, error) {
	var result []Message
	rows, err := tx.Model(&Message{}).Where("channel_id = ?", channelId).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	var m Message
	for rows.Next() {
		if err := tx.ScanRows(rows, &m); err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}
