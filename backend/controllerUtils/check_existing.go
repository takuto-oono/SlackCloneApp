package controllerUtils

import (
	"backend/models"
)

func IsExistChannelAndUserInSameWorkspace(channelId int, userId uint32) (bool, error) {
	ch, err := models.GetChannelById(channelId)
	if err != nil {
		return false, err
	}
	wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(ch.WorkspaceId, userId)
	if err != nil {
		return false, err
	}
	return wau.WorkspaceId == ch.WorkspaceId, nil
}