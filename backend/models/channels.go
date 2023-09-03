package models

import (
	"gorm.io/gorm"
)

type Channel struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"not null"`
	IsPrivate   bool   `json:"is_private" gorm:"default:false"`
	IsArchive   bool   `json:"is_archive" gorm:"default:false"`
	WorkspaceId int    `json:"workspace_id" gorm:"not null"`
}

func NewChannel(name, description string, isPrivate, isArchive bool, workspaceId int) *Channel {
	return &Channel{
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		IsArchive:   isArchive,
		WorkspaceId: workspaceId,
	}
}

func (c *Channel) Create(tx *gorm.DB) error {
	return tx.Model(&Channel{}).Create(c).Error
}

func GetChannelById(tx *gorm.DB, channelId int) (Channel, error) {
	var result Channel
	err := tx.Model(&Channel{}).Where("id = ?", channelId).Take(&result).Error
	return result, err
}

func (c Channel) Delete(tx *gorm.DB) error {
	return tx.Where("id = ?", c.ID).Delete(&Channel{}).Error
}

func GetChannelByIdAndWorkspaceId(tx *gorm.DB, id, workspaceId int) (Channel, error) {
	var result Channel
	err := tx.Model(&Channel{}).Where("id = ? AND workspace_id = ?", id, workspaceId).Take(&result).Error
	return result, err
}

func GetChannelsByWorkspaceId(tx *gorm.DB, workspaceId int) ([]Channel, error) {
	var result []Channel
	rows, err := tx.Model(&Channel{}).Where("workspace_id = ?", workspaceId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Channel
		if err := tx.ScanRows(rows, &c); err != nil {
			return result, err
		}
		result = append(result, c)
	}
	return result, nil
}

func GetAllChannels(tx *gorm.DB) ([]Channel, error) {
	var result []Channel
	rows, err := tx.Model(&Channel{}).Rows()
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var c Channel
		if err := tx.ScanRows(rows, &c); err != nil {
			return result, err
		}
		result = append(result, c)
	}
	return result, nil
}

func GetAllChannelsByUserID(tx *gorm.DB, userID uint32) ([]Channel, error) {
	var result []Channel
	rows, err := tx.Model(&Channel{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var c Channel
		if err := tx.ScanRows(rows, &c); err != nil {
			return result, err
		}
		result = append(result, c)
	}
	return result, nil
}

func DeleteChannelsTableRecords() {
	db.Exec("DELETE FROM channels")
}
