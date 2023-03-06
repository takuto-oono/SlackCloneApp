package models

import (
	"time"

	"gorm.io/gorm"
)

type DirectMessage struct {
	gorm.Model
	ID         uint      `json:"id" gorm:"primaryKey"`
	Text       string    `json:"text" gorm:"not null"`
	SendUserId uint32    `json:"send_user_id" gorm:"not null"`
	DMLineId   uint      `json:"dm_line_id" gorm:"not null; column: dm_line_id"`
	CreatedAt  time.Time `json:"create_at" gorm:"not null"`
	UpdatedAt  time.Time `json:"update_at" gorm:"not null"`
}

func NewDirectMessage(text string, sendUserId uint32, dmLineId uint) *DirectMessage {
	return &DirectMessage{
		Text:       text,
		SendUserId: sendUserId,
		DMLineId:   dmLineId,
	}
}

func (dm *DirectMessage) Create() *gorm.DB {
	return db.Create(dm)
}
