package models

import (
	"gorm.io/gorm"
)

type Workspace struct {
	ID             int    `json:"id" gorm:"primaryKey"`
	Name           string `json:"name" gorm:"not null; unique"`
	PrimaryOwnerId uint32 `json:"primary_owner_id" gorm:"not null"`
}

func NewWorkspace(name string, primaryOwnerId uint32) *Workspace {
	return &Workspace{
		Name:           name,
		PrimaryOwnerId: primaryOwnerId,
	}
}

func (w *Workspace) Create(tx *gorm.DB) error {
	return tx.Create(w).Error
}

func GetWorkspaceById(tx *gorm.DB, id int) (Workspace, error) {
	var result Workspace
	err := tx.Model(&Workspace{}).Where("id = ?", id).Take(&result).Error
	return result, err
}

func UpdateWorkspaceName(tx *gorm.DB, id int, name string) (Workspace, error) {
	var result Workspace
	err := tx.Model(&Workspace{}).Where("id = ?", id).Update("name", name).Row().Scan(
		&result.ID,
		&result.Name,
		&result.PrimaryOwnerId,
	)
	return result, err
}

func DeleteWorkspacesTableRecords() {
	db.Exec("DELETE FROM workspaces")
}
