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

func (dm *DirectMessage) Create() *gorm.DB {
	return db.Create(dm)
}

func GetAllDMsByDLId(dmLineId uint) ([]DirectMessage, error) {
	var result []DirectMessage
	rows, err := db.Model(&DirectMessage{}).Where("dm_line_id = ?", dmLineId).Order("created_at desc").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var dm DirectMessage
		db.ScanRows(rows, &dm)
		result = append(result, dm)
	}
	return result, nil
}

func GetDMById(id uint) (DirectMessage, error) {
	var result DirectMessage
	err := db.Model(&DirectMessage{}).Where("id = ?", id).First(
		&result).Error
	return result, err
}

func UpdateDM(id uint, text string) (DirectMessage, error) {
	var result DirectMessage
	err := db.Model(&DirectMessage{}).Where("id = ?", id).Update("text", text).Row().Scan(
		&result.ID,
		&result.Text,
		&result.SendUserId,
		&result.DMLineId,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	return result, err
}
