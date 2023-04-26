package controllerUtils

import (
	"errors"

	"gorm.io/gorm"

	"backend/models"
)

func CreateOrGetThreadByParentMessageId(tx *gorm.DB, parentMessageId uint) (models.Thread, error) {
	th, err := models.GetThreadByParentMessageId(tx, parentMessageId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			nth := models.NewThread(parentMessageId)
			err = nth.Create(tx)
			return *nth, err
		}
		return th, err
	}
	return th, nil
}
