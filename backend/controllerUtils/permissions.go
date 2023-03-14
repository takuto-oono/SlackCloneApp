package controllerUtils

import (
	"fmt"

	"gorm.io/gorm"

	"backend/models"
)

func HasPermissionAddUserInWorkspace(userId uint32, workspaceId int) bool {
	wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId, userId)
	if err != nil {
		return false
	}
	return wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3
}

func HasPermissionRenamingWorkspaceName(workspaceId int, userId uint32) (bool, error) {
	wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId, userId)
	if err != nil {
		return false, err
	}
	return (wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3), nil
}

func HasPermissionDeletingUserFromWorkspace(workspaceId int, userId uint32) (bool, error) { wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId, userId)
	if err != nil {
		return false, err
	}
	return (wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3), nil	
}

func HasPermissionAddingUserInChannel(channelId int, userId uint32) bool {
	return models.IsAdminUserInChannel(channelId, userId)
}

func HasPermissionDeletingUserInChannel(userId uint32, workspaceId int, ch models.Channel) bool {
	if ch.IsPrivate {
		return models.IsExistCAUByChannelIdAndUserId(ch.ID, userId)
	}
	roleId, err := models.GetRoleIdByWorkspaceIdAndUserId(workspaceId, userId)
	if err != nil {
		return false
	}
	return roleId == 1 || roleId == 2 || roleId == 3
}

func HasPermissionDeletingChannel(wau models.WorkspaceAndUsers) bool { 
	return (wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3)	
}

func HasPermissionEditDM(dmId uint, userId uint32) bool {
	dm, err := models.GetDMById(dmId)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			fmt.Println(err.Error())
		}
		return false
	}
	return dm.SendUserId == userId

}
