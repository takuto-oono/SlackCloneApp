package controllerUtils

import (
	"backend/models"

	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = models.ExportDB()
}
