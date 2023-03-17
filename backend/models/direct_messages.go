package models

import (
	"time"

	"gorm.io/gorm"
)

type DirectMessage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Text       string    `json:"text" gorm:"not null"`
	SendUserId uint32    `json:"send_user_id" gorm:"not null"`
	DMLineId   uint      `json:"dm_line_id" gorm:"not null; column:dm_line_id"`
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

func (dm *DirectMessage) Create(tx *gorm.DB) error {
	return tx.Create(dm).Error
}

func GetAllDMsByDLId(tx *gorm.DB, dmLineId uint) ([]DirectMessage, error) {
	var result []DirectMessage
	rows, err := tx.Model(&DirectMessage{}).Where("dm_line_id = ?", dmLineId).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var dm DirectMessage
		if err := tx.ScanRows(rows, &dm); err != nil {
			return result, err
		}
		result = append(result, dm)
	}
	return result, nil
}

func GetDMById(tx *gorm.DB, id uint) (DirectMessage, error) {
	var result DirectMessage
	err := tx.Model(&DirectMessage{}).Where("id = ?", id).Take(
		&result).Error
	return result, err
}

func UpdateDM(tx *gorm.DB, id uint, text string) (DirectMessage, error) {
	var result DirectMessage
	err := tx.Model(&DirectMessage{}).Where("id = ?", id).Update("text", text).Row().Scan(
		&result.ID,
		&result.Text,
		&result.SendUserId,
		&result.DMLineId,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	return result, err
}

func (dm DirectMessage) DeleteDM(tx *gorm.DB) error {
	return tx.Where("id = ?", dm.ID).Delete(&DirectMessage{}).Error
}
