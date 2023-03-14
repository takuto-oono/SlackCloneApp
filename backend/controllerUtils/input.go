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

type SendMessageInput struct {
	Text      string `json:"text"`
	ChannelId int    `json:"channel_id"`
}

type SendDMInput struct {
	ReceiveUserId uint32 `json:"received_user_id"`
	Text          string `json:"text"`
	WorkspaceId   int    `json:"workspace_id"`
}

type EditDMInput struct {
	Text string `json:"text"`
}

func InputSignUpAndLogin(c *gin.Context) (SignUpAndLoginInput, error) {
	var in SignUpAndLoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.Name == "" || in.Password == "" {
		return in, fmt.Errorf("name or password not found")
	}
	return in, nil
}

func InputAndValidateCreateWorkspace(c *gin.Context) (CreateWorkspaceInput, error) {
	var in CreateWorkspaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.Name == "" {
		return in, fmt.Errorf("name not found")
	}
	if in.RequestUserId == 0 {
		return in, fmt.Errorf("user_id not found")
	}
	return in, nil
}

func InputAndValidateAddUserInWorkspace(c *gin.Context) (AddUserInWorkspaceInput, error) {
	var in AddUserInWorkspaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("workspace_id not found")
	}

	if in.UserId == 0 {
		return in, fmt.Errorf("user_id not found")
	}

	if in.RoleId == 0 {
		return in, fmt.Errorf("role_id not found")
	}
	return in, nil
}

func InputAndValidateRenameWorkspace(c *gin.Context) (RenameWorkspaceNameInput, error) {
	var in RenameWorkspaceNameInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.UserId == 0 {
		return in, fmt.Errorf("user_id not found")
	}
	if in.WorkspaceName == "" {
		return in, fmt.Errorf("workspace_name not found")
	}
	return in, nil
}

func InputAndValidateDeleteUserFromWorkspace(c *gin.Context) (DeleteUserFromWorkSpaceInput, error) {
	var in DeleteUserFromWorkSpaceInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("workspace_id not found")
	}
	if in.UserId == 0 {
		return in, fmt.Errorf("user_id not found")
	}
	return in, nil
}

func InputAndValidateCreateChannel(c *gin.Context) (CreateChannelInput, error) {
	var in CreateChannelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("workspace_id not found")
	}
	if in.Name == "" {
		return in, fmt.Errorf("name not found")
	}
	if in.IsPrivate == nil {
		return in, fmt.Errorf("is_private not found")
	}
	return in, nil
}

func InputAndValidateAddUserInChannel(c *gin.Context) (AddUserInChannelInput, error) {
	var in AddUserInChannelInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.ChannelId == 0 {
		return in, fmt.Errorf("channel_id not found")
	}
	if in.UserId == 0 {
		return in, fmt.Errorf("user_id not found")
	}
	return in, nil
}

func InputAndValidateSendMessage(c *gin.Context) (SendMessageInput, error) {
	var in SendMessageInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.ChannelId == 0 {
		return in, fmt.Errorf("channel_id not found")
	}
	if in.Text == "" {
		return in, fmt.Errorf("text not found")
	}
	return in, nil
}

func InputAndValidateSendDM(c *gin.Context) (SendDMInput, error) {
	var in SendDMInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.ReceiveUserId == 0 {
		return in, fmt.Errorf("received_user_id not found")
	}
	if in.WorkspaceId == 0 {
		return in, fmt.Errorf("workspace_id not found")
	}
	if in.Text == "" {
		return in, fmt.Errorf("text not found")
	}
	return in, nil
}

func InputAndValidateEditDM(c *gin.Context) (EditDMInput, error) {
	var in EditDMInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.Text == "" {
		return in, fmt.Errorf("text not found")
	}
	return in, nil
}