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

func IsExistCAUByChannelIdAndUserId(channelId int, userId uint32) (bool, error) {
	cau, err := models.GetCAUByChannelIdAndUserId(channelId, userId)
	if err != nil {
		return false, err
	}
	return cau.ChannelId == channelId && cau.UserId == userId, nil
}
