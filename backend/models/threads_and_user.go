package models

import (
	"gorm.io/gorm"
)

type ThreadAndUser struct {
	UserId   uint32 `json:"user_id" gorm:"primaryKey"`
	ThreadId uint   `json:"thread_id" gorm:"primaryKey"`
}

func NewThreadAndUser(userId uint32, threadId uint) *ThreadAndUser {
	return &ThreadAndUser{
		UserId:   userId,
		ThreadId: threadId,
	}
}

func (tau *ThreadAndUser) Create(tx *gorm.DB) error {
	return tx.Model(&ThreadAndUser{}).Create(tau).Error
}

func GetTAUByThreadIdAndUserId(tx *gorm.DB, userId uint32, threadId uint) (ThreadAndUser, error) {
	var result ThreadAndUser
	err := tx.Model(&ThreadAndUser{}).Where("user_id = ? AND thread_id = ?", userId, threadId).Take(&result).Error
	return result, err
}

func GetTAUsByUserId(tx *gorm.DB, userId uint32) ([]ThreadAndUser, error) {
	var result []ThreadAndUser
	rows, err := tx.Model(&ThreadAndUser{}).Where("user_id = ?", userId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var tau ThreadAndUser
		if err := tx.ScanRows(rows, &tau); err != nil {
			return result, err
		}
		result = append(result, tau)
	}
	return result, nil
}

func GetTAUsByThreadId(tx *gorm.DB, threadId uint) ([]ThreadAndUser, error) {
	var result []ThreadAndUser
	rows, err := tx.Model(&ThreadAndUser{}).Where("thread_id = ?", threadId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var tau ThreadAndUser
		if err := tx.ScanRows(rows, &tau); err != nil {
			return result, err
		}
		result = append(result, tau)
	}
	return result, nil
}

func DeleteThreadAndUsersTableRecords() {
	db.Exec("DELETE FROM thread_and_users")
}
