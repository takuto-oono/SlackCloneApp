package models

import (
	"gorm.io/gorm"
)

type WorkspaceAndUsers struct {
	WorkspaceId int    `json:"workspace_id" gorm:"primaryKey"`
	UserId      uint32 `json:"user_id" gorm:"primaryKey"`
	RoleId      int    `json:"role_id" gorm:"not null"`
}

func NewWorkspaceAndUsers(workspaceId int, userId uint32, roleId int) *WorkspaceAndUsers {
	return &WorkspaceAndUsers{
		WorkspaceId: workspaceId,
		UserId:      userId,
		RoleId:      roleId,
	}
}

func (wau *WorkspaceAndUsers) Create(tx *gorm.DB) error {
	return tx.Create(wau).Error
}

func GetWAUByWorkspaceIdAndUserId(tx *gorm.DB, workspaceId int, userId uint32) (WorkspaceAndUsers, error) {
	var result WorkspaceAndUsers
	err := tx.Model(&WorkspaceAndUsers{}).Where("workspace_id = ? AND user_id = ?", workspaceId, userId).Take(&result).Error
	return result, err
}

func (wau WorkspaceAndUsers) DeleteWorkspaceAndUser(tx *gorm.DB) error {
	return tx.Where("workspace_id = ? AND user_id = ?", wau.WorkspaceId, wau.UserId).Delete(&wau).Error
}

func GetWAUsByUserId(tx *gorm.DB, userId uint32) ([]WorkspaceAndUsers, error) {
	var result []WorkspaceAndUsers
	rows, err := tx.Model(&WorkspaceAndUsers{}).Where("user_id = ?", userId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var wau WorkspaceAndUsers
		if err := tx.ScanRows(rows, &wau); err != nil {
			return result, err
		}
		result = append(result, wau)
	}
	return result, nil
}

func GetWAUsByWorkspaceId(tx *gorm.DB, workspaceId int) ([]WorkspaceAndUsers, error) {
	var result []WorkspaceAndUsers
	rows, err := tx.Model(&WorkspaceAndUsers{}).Where("workspace_id = ?", workspaceId).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var wau WorkspaceAndUsers
		if err := tx.ScanRows(rows, &wau); err != nil {
			return result, err
		}
		result = append(result, wau)
	}
	return result, nil
}

func DeleteWorkspaceAndUsersTableRecords() {
	db.Exec("DELETE FROM workspace_and_users")
}
