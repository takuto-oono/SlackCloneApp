package models

import (
	"gorm.io/gorm"
)

type DMLine struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	WorkspaceId int    `json:"workspace_id" gorm:"not null"`
	UserId1     uint32 `json:"user_id_1" gorm:"not null; column:user_id_1"`
	UserId2     uint32 `json:"user_id_2" gorm:"not null; column:user_id_2"`
}

func NewDMLine(workspaceId int, userId1, userId2 uint32) *DMLine {
	return &DMLine{
		WorkspaceId: workspaceId,
		UserId1:     userId1,
		UserId2:     userId2,
	}
}

func (dl *DMLine) Create(tx *gorm.DB) error {
	if !(dl.UserId1 <= dl.UserId2) {
		dl.UserId1, dl.UserId2 = dl.UserId2, dl.UserId1
	}
	return tx.Create(dl).Error
}

func GetDLByUserIdsAndWorkspaceId(tx *gorm.DB, userId1, userId2 uint32, workspaceId int) (DMLine, error) {
	var result DMLine
	if !(userId1 <= userId2) {
		userId1, userId2 = userId2, userId1
	}

	err := tx.Model(&DMLine{}).Where("user_id_1 = ? AND user_id_2 = ? AND workspace_id = ?", userId1, userId2, workspaceId).Take(&result).Error
	return result, err
}

func GetDLById(tx *gorm.DB, id uint) (DMLine, error) {
	var result DMLine
	err := tx.Model(&DMLine{}).Where("id = ?", id).Take(&result).Error
	return result, err
}
