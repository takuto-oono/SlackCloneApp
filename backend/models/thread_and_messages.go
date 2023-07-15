package models

import (
	"time"

	"gorm.io/gorm"
)

type ThreadAndMessage struct {
	ThreadId  uint      `json:"id" gorm:"primaryKey"`
	MessageId uint      `json:"message_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
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

func GetTAMByThreadIdAndMessageId(tx *gorm.DB, threadId, messageId uint) (ThreadAndMessage, error) {
	var result ThreadAndMessage
	err := tx.Model(&ThreadAndMessage{}).Where("thread_id = ? AND message_id = ?", threadId, messageId).Take(&result).Error
	return result, err
}

func GetTAMsByThreadId(tx *gorm.DB, threadId uint) ([]ThreadAndMessage, error) {
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

func DeleteThreadAndMessagesTableRecords() {
	db.Exec("DELETE FROM thread_and_messages")
}
