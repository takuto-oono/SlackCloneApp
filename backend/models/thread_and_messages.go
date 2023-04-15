package models

import (
	"gorm.io/gorm"
)

type ThreadAndMessage struct {
	ThreadId  uint `json:"id" gorm:"primaryKey"`
	MessageId uint `json:"message_id" gorm:"primaryKey"`
}

func NewThreadAndMessage(threadId, messageId uint) *ThreadAndMessage {
	return &ThreadAndMessage{
		ThreadId:  threadId,
		MessageId: messageId,
	}
}

func (tam *ThreadAndMessage) Create(tx *gorm.DB) error {
	return tx.Model(&ThreadAndMessage{}).Create(tam).Error
}

func GetTAMByThreadID(tx *gorm.DB, threadId uint) ([]ThreadAndMessage, error) {
	var result []ThreadAndMessage
	rows, err := tx.Model(&ThreadAndMessage{}).Where("thread_id = ?", threadId).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	var tam ThreadAndMessage
	for rows.Next() {
		if err := tx.ScanRows(rows, &tam); err != nil {
			return result, err
		}
		result = append(result, tam)
	}
	return result, nil
}

func (tam ThreadAndMessage) Delete(tx *gorm.DB) error {
	return tx.Where("thread_id = ? AND message_id = ?", tam.ThreadId, tam.MessageId).Delete(&ThreadAndMessage{}).Error
}
