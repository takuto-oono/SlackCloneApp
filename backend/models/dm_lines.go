package models

import (
	"fmt"

	"gorm.io/gorm"
)

type DMLine struct {
	ID          uint   `json:"id" gorm:"unique"`
	WorkspaceId int    `json:"workspace_id" gorm:"not null"`
	UserId1     uint32 `json:"user_id_1" gorm:"primaryKey; column:user_id_1"`
	UserId2     uint32 `json:"user_id_2" gorm:"primaryKey; column:user_id_2"`
}

func NewDMLine(workspaceId int, userId1, userId2 uint32) *DMLine {
	return &DMLine{
		WorkspaceId: workspaceId,
		UserId1:     userId1,
		UserId2:     userId2,
	}
}

func (dl *DMLine) SetID() {
	rows, err := db.Model(&DMLine{}).Rows()
	if err != nil {
		fmt.Println(err)
	}
	cnt := 0
	defer rows.Close()
	for rows.Next() {
		cnt++
	}
	dl.ID = uint(cnt + 1)
}

func (dl *DMLine) Create() *gorm.DB {
	dl.SetID()
	if !(dl.UserId1 <= dl.UserId2) {
		dl.UserId1, dl.UserId2 = dl.UserId2, dl.UserId1
	}
	return db.Create(dl)
}

func GetDLByUserIdsAndWorkspaceId(userId1, userId2 uint32, workspaceId int) (DMLine, error) {
	var dm_line DMLine
	if !(userId1 <= userId2) {
		userId1, userId2 = userId2, userId1
	}
	result := db.First(&dm_line, "user_id_1 = ? AND user_id_2 = ? AND workspace_id = ?", userId1, userId2, workspaceId)
	return dm_line, result.Error
}

func GetDLById(id uint) (DMLine, error) {
	var dl DMLine
	result := db.First(&dl, "id = ?", id)
	return dl, result.Error
}
