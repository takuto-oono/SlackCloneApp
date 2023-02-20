package controllerUtils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type SignUpAndLoginInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateWorkspaceInput struct {
	Name          string `json:"name"`
	RequestUserId uint32 `json:"user_id"`
}

type AddUserInWorkspaceInput struct {
	WorkspaceId int    `json:"workspace_id"`
	UserId      uint32 `json:"user_id"`
	RoleId      int    `json:"role_id"`
}

type RenameWorkspaceNameInput struct {
	UserId        uint32 `json:"user_id"`
	WorkspaceName string `json:"workspace_name"`
}

type DeleteUserFromWorkSpaceInput struct {
	WorkspaceId int    `json:"workspace_id"`
	UserId      uint32 `json:"user_id"`
}

type CreateChannelInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   *bool  `json:"is_private"`
	WorkspaceId int    `json:"workspace_id"`
}

type AddUserInChannelInput struct {
	ChannelId int    `json:"channel_id"`
	UserId    uint32 `json:"user_id"`
}

func InputSignUpAndLogin(c *gin.Context) (SignUpAndLoginInput, error) {
	var in SignUpAndLoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.Name == "" || in.Password == "" {
		return in, fmt.Errorf("not found name or password")
	}
	return in, nil
}

func InputAndValidateCreateWorkspace(c *gin.Context) (CreateWorkspaceInput, error) {
	var in CreateWorkspaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.Name == "" {
		return in, fmt.Errorf("not found name")
	}
	if in.RequestUserId == 0 {
		return in, fmt.Errorf("not found user_id")
	}
	return in, nil
}

func InputAndValidateAddUserInWorkspace(c *gin.Context) (AddUserInWorkspaceInput, error) {
	var in AddUserInWorkspaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("not found workspace_id")
	}

	if in.UserId == 0 {
		return in, fmt.Errorf("not found user_id")
	}

	if in.RoleId == 0 {
		return in, fmt.Errorf("not found role_id")
	}
	return in, nil
}

func InputAndValidateRenameWorkspace(c *gin.Context) (RenameWorkspaceNameInput, error) {
	var in RenameWorkspaceNameInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.UserId == 0 {
		return in, fmt.Errorf("not found user_id")
	}
	if in.WorkspaceName == "" {
		return in, fmt.Errorf("not found workspace_name")
	}
	return in, nil
}

func InputAndValidateDeleteUserFromWorkspace(c *gin.Context) (DeleteUserFromWorkSpaceInput, error) {
	var in DeleteUserFromWorkSpaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("not found workspace_id")
	}
	if in.UserId == 0 {
		return in, fmt.Errorf("not found user_id")
	}
	return in, nil
}

func InputAndValidateCreateChannel(c *gin.Context) (CreateChannelInput, error) {
	var in CreateChannelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("not found workspace_id")
	}
	if in.Name == "" {
		return in, fmt.Errorf("not found name")
	}
	if in.IsPrivate == nil {
		return in, fmt.Errorf("not found is_private")
	}
	return in, nil
}

func InputAndValidateAddUserInChannel(c *gin.Context) (AddUserInChannelInput, error) {
	var in AddUserInChannelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.ChannelId == 0 {
		return in, fmt.Errorf("not found channel_id")
	}
	if in.UserId == 0 {
		return in, fmt.Errorf("not found user_id")
	}
	return in, nil
}
