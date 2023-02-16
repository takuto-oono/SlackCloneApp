package utils

import (
	"backend/models"
)

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
