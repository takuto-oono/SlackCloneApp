package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Thread struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ParentMessageId uint      `json:"parent_message_id" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`
}

func NewThread(parentMessageId uint) *Thread {
	return &Thread{
		ParentMessageId: parentMessageId,
	}
}

func (th *Thread) Create(tx *gorm.DB) error {
	th.UpdatedAt = time.Now()
	return tx.Model(&Thread{}).Create(th).Error
}

func (th *Thread) EditUpdatedAt(tx *gorm.DB) error {
	return tx.Model(&Thread{}).Where("id = ?", th.ID).Update("updated_at", time.Now()).Error
}

func GetThreadById(tx *gorm.DB, id uint) (*Thread, error) {
	var result *Thread
	err := tx.Model(&Thread{}).Where("id = ?", id).Take(result).Error
	return result, err
}

func GetThreadByParentMessageId(tx *gorm.DB, parentMessageId uint) (*Thread, error) {
	var result Thread
	fmt.Println(1)
	err := tx.Model(&Thread{}).Where("parent_message_id = ?", parentMessageId).Take(&result).Error
	fmt.Println(err.Error())
	return &result, err
}
