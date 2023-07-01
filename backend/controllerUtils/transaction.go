package controllerUtils

import (
	"gorm.io/gorm"

	"backend/models"
)

func initTX() (*gorm.DB, error) {
	tx := db.Begin()
	return tx, tx.Error
}

func SendMessageTX(m *models.Message, userID uint32, mentionedUserIDs []uint32) error {
	tx, err := initTX()
	if err != nil {
		return err
	}

	// message情報をDBに登録
	if err := m.Create(tx); err != nil {
		tx.Rollback()
		return err
	}

	// mentionの処理をする
	for _, userID := range mentionedUserIDs {
		men := models.NewMention(userID, m.ID)
		if err := men.Create(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// channelにいるuserをmessage_and_users tableに追加する
	caus, err := models.GetCAUsByChannelId(tx, m.ChannelId)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, cau := range caus {
		if cau.UserId == userID {
			continue
		}
		mau := models.NewMessageAndUser(m.ID, cau.UserId, false)
		if err := mau.Create(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
