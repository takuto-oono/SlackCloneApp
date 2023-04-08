package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint32 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	PassWord string `json:"password" gorm:"not null; column:password"`
}

func NewUser(id uint32, name, password string) *User {
	return &User{ID: id, Name: name, PassWord: password}
}

func (u *User) Create(tx *gorm.DB) error {
	return tx.Create(u).Error
}

func GetUserById(tx *gorm.DB, id uint32) (User, error) {
	var result User
	err := tx.Model(&User{}).Where("id = ?", id).Take(&result).Error
	return result, err
}

func GetUserByNameAndPassword(tx *gorm.DB, username, password string) (User, error) {
	var result User
	err := tx.Model(&User{}).Where("name = ? AND password = ?", username, password).Take(&result).Error
	return result, err
}

func GetUsers(tx *gorm.DB) ([]User, error) {
	var result []User
	rows, err := tx.Model(&User{}).Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		if err := tx.ScanRows(rows, &u); err != nil {
			return result, err
		}
		result = append(result, u)
	}
	return result, nil
}
