package controllerUtils

import (
	"fmt"

	"gorm.io/gorm"

	"backend/models"
)

func HasPermissionAddUserInWorkspace(userId uint32, workspaceId int) bool {
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, workspaceId, userId)
	if err != nil {
		return false
	}
	return wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3
}

func HasPermissionRenamingWorkspaceName(workspaceId int, userId uint32) (bool, error) {
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, workspaceId, userId)
	if err != nil {
		return false, err
	}
	return (wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3), nil
}

func HasPermissionDeletingUserFromWorkspace(workspaceId int, userId uint32) (bool, error) { wau, err := models.GetWAUByWorkspaceIdAndUserId(db, workspaceId, userId)
	if err != nil {
		return false, err
	}
	return (wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3), nil	
}

func HasPermissionAddingUserInChannel(channelId int, userId uint32) bool {
	cau, err := models.GetCAUByChannelIdAndUserId(db, channelId, userId)
	if err != nil {
		return false
	}
	return cau.IsAdmin
}

func HasPermissionDeletingUserInChannel(userId uint32, workspaceId int, ch models.Channel) bool {
	if ch.IsPrivate {
		cau, err := models.GetCAUByChannelIdAndUserId(db, ch.ID, userId)
		if err != nil {
			return false
		}
		return cau.IsAdmin
	}
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, workspaceId, userId)
	if err != nil {
		return false
	}
	return wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3
}

func HasPermissionDeletingChannel(wau models.WorkspaceAndUsers) bool { 
	return (wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3)	
}

func HasPermissionEditDM(dmId uint, userId uint32) bool {
	dm, err := models.GetDMById(db, dmId)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			fmt.Println(err.Error())
		}
		return false
	}
	return dm.SendUserId == userId
}

func HasPermissionEditMessage(messageId int, userId uint32) bool {
	// 作成したuserと編集アクセスをしたuserが同じかを判定
	
	m, err := models.GetMessageById(db, messageId)
	if err != nil {
		return false
	}
	return m.UserId == userId
}
