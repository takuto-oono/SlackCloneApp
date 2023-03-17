package controllerUtils

import (
	"fmt"

	"gorm.io/gorm"

	"backend/models"
)

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

func IsExistDMById(dmId uint) (bool, error) {
	_, err := models.GetDMById(db, dmId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// not found errorはfalse, nilを返すようにする
			return false, nil
		}
		return false, err
	}
	return true, nil
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
