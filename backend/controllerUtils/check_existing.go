package controllerUtils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"backend/models"
)

func errHandling(err error) (bool, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func IsExistChannelAndUserInSameWorkspace(channelId int, userId uint32) (bool, error) {
	ch, err := models.GetChannelById(db, channelId)
	if err != nil {
		return false, err
	}
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, ch.WorkspaceId, userId)
	if err != nil {
		return false, err
	}
	return wau.WorkspaceId == ch.WorkspaceId, nil
}

func IsExistCAUByChannelIdAndUserId(channelId int, userId uint32) (bool, error) {
	cau, err := models.GetCAUByChannelIdAndUserId(db, channelId, userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return cau.ChannelId == channelId && cau.UserId == userId, nil
}

func IsExistUserSameUsernameAndPassword(userName, password string) bool {
	u, err := models.GetUserByNameAndPassword(db, userName, password)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return u.Name == userName && u.PassWord == password
}

func IsExistWorkspaceById(id int) bool {
	w, err := models.GetWorkspaceById(db, id)
	if err != nil {
		fmt.Println(err)
	}
	return w.ID == id
}

func IsExistWAUByWorkspaceIdAndUserId(workspaceId int, userId uint32) bool {
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, workspaceId, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return wau.WorkspaceId == workspaceId && wau.UserId == userId
}

func IsExistChannelByChannelIdAndWorkspaceId(channelId, workspaceId int) (bool, error) {
	_, err := models.GetChannelByIdAndWorkspaceId(db, channelId, workspaceId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsExistSameNameChannelInWorkspace(channelName string, workspaceId int) (bool, error) {
	chs, err := models.GetChannelsByWorkspaceId(db, workspaceId)
	if err != nil {
		return false, err
	}
	for _, ch := range chs {
		if ch.Name == channelName {
			return true, nil
		}
	}
	return false, nil
}

func IsExistMessageById(messageId uint) (bool, error) {
	_, err := models.GetMessageById(db, messageId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsExistUserInDL(userId uint32, dlId uint) (bool, error) {
	dl, err := models.GetDLById(db, dlId)
	if err != nil {
		return false, err
	}
	return dl.UserId1 == userId || dl.UserId2 == userId, nil
}

func IsExistThreadByMessageId(messageId uint) (bool, error) {
	_, err := models.GetThreadByParentMessageId(db, messageId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsExistTAMByThreadIdAndMessageId(threadId, messageId uint) (bool, error) {
	_, err := models.GetTAMByThreadIdAndMessageId(db, threadId, messageId)
	return errHandling(err)
}

func IsExistTAUByUserIdAndThreadId(tx *gorm.DB, userId uint32, threadId uint) (bool, error) {
	_, err := models.GetTAUByThreadIdAndUserId(tx, userId, threadId)
	return errHandling(err)
}

func IsExistThreadInWorkspace(tx *gorm.DB, thread models.Thread, channels []models.Channel, dls []models.DMLine) (bool, error) {
	parentMessage, err := models.GetMessageById(db, thread.ParentMessageId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if parentMessage.ChannelId != 0 && parentMessage.DMLineId == uint(0) {
		for _, ch := range channels {
			if parentMessage.ChannelId == ch.ID {
				return true, nil
			}
		}
	} else if parentMessage.ChannelId == 0 && parentMessage.DMLineId != uint(0) {
		for _, dl := range dls {
			if parentMessage.DMLineId == dl.ID {
				return true, nil
			}
		}
	}
	return false, nil
}

func IsExistUserInWorkspace(tx *gorm.DB, userID uint32, workspaceID int) (bool, error) {
	_, err := models.GetWAUByWorkspaceIdAndUserId(tx, workspaceID, userID)
	return errHandling(err)
}
