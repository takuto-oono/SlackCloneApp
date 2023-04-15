package controllerUtils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"backend/models"
)

func CreateOrGetThreadByParentMessageId(tx *gorm.DB, parentMessageId uint) (*models.Thread, error) {
	th, err := models.GetThreadByParentMessageId(tx, parentMessageId)
	fmt.Println("1")

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			th = models.NewThread(parentMessageId)
			fmt.Println("2")
			return th, nil
		}
		fmt.Println("4")
		return th, err
	}
	fmt.Println("3")
	return th, nil
}
