package models

import (
	"gorm.io/gorm"
)

type Role struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
}

func NewRole(id int, name string) *Role {
	return &Role{ID: id, Name: name}
}

func (r *Role) Create(tx *gorm.DB) error {
	return tx.Create(r).Error
}

func GetRoleById(tx *gorm.DB, id int) (Role, error) {
	var result Role
	err := tx.Model(&Role{}).Where("id = ?", id).Take(&result).Error
	return result, err
}
