package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text" gorm:"not null"`
	ChannelId int       `json:"channel_id"`
	DMLineId  uint      `json:"dm_line_id" gorm:"column:dm_line_id"`
	UserId    uint32    `json:"user_id" gorm:"not null"`
	ThreadId  uint      `json:"thread_id"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func NewChannelMessage(text string, channelId int, userId uint32) *Message {
	return &Message{
		Text:      text,
		ChannelId: channelId,
		UserId:    userId,
		DMLineId:  uint(0),
	}
}

func NewDMMessage(text string, dmLineId uint, userId uint32) *Message {
	return &Message{
		Text:      text,
		DMLineId:  dmLineId,
		UserId:    userId,
		ChannelId: 0,
	}
}

func (m *Message) Create(tx *gorm.DB) error {
	if m.ChannelId != 0 && m.DMLineId != uint(0) {
		return fmt.Errorf("channelId and dmLineId equal 0")
	}
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

func GetMessagesByDLId(tx *gorm.DB, dlId uint) ([]Message, error) {
	var result []Message
	rows, err := tx.Model(&Message{}).Where("dm_line_id = ?", dlId).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Message
		if err := tx.ScanRows(rows, &m); err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}

func GetMessagesByThreadId(tx *gorm.DB, threadId uint) ([]Message, error) {
	var result []Message
	rows, err := tx.Model(&Message{}).Where("thread_id = ?", threadId).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Message
		if err := tx.ScanRows(rows, &m); err != nil {
			return result, err
		}
		result = append(result, m)
	}
	return result, nil
}

func GetMessageById(tx *gorm.DB, id uint) (Message, error) {
	var result Message
	err := tx.Model(&Message{}).Where("id = ?", id).Take(&result).Error
	return result, err
}

func UpdateMessageText(tx *gorm.DB, id uint, text string) (Message, error) {
	var result Message
	err := tx.Model(&Message{}).Where("id = ?", id).Update("text", text).Take(&result).Error
	return result, err
}

func (m *Message) UpdateMessageThreadId(tx *gorm.DB, threadId uint) error {
	return tx.Model(&Message{}).Where("id = ?", m.ID).Update("thread_id", threadId).Take(m).Error
}

func (m Message) Delete(tx *gorm.DB) error {
	return tx.Where("id = ?", m.ID).Delete(&Message{}).Error
}
