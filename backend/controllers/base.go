package controllers

import (
	"gorm.io/gorm"

	"backend/models"
)

var db *gorm.DB

func init() {
	db = models.ExportDB()
}
