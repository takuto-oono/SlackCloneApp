package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/models"
)

type WorkspaceInput struct {
	Name string `json:"name"`
}

type AddUserWorkspaceInput struct {
	WorkspaceName string `json:"workspace_name"`
	AddUserName   string `json:"add_user_name"`
	RoleId        int    `json:"role_id"`
}

func CreateWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	primaryOwnerId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var input WorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspace name"})
		return
	}
	w := models.NewWorkspace(0, input.Name, primaryOwnerId)
	if err := w.CreateWorkspace(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	wau := models.NewWorkspaceAndUsers(w.ID, w.PrimaryOwnerId, 1)
	err = wau.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, w)
}

func AddUserWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var input AddUserWorkspaceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if input.WorkspaceName == "" || input.AddUserName == "" || input.RoleId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "field empty"})
		return
	}

	// workspaceNameからworkspaceを取得
	w, err := models.GetWorkspaceByName(input.WorkspaceName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// userNameからuserIdを取得
	u, err := models.GetUserByName(input.AddUserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// userIdがそのworkspaceで追加する権限を持っているかを判定(roleId == 1 or roleId == 2 or roleId == 3)
	wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(w.ID, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if !(wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized add user in workspace"})
		return
	}
	// userをワークスペースに追加
	nwau := models.NewWorkspaceAndUsers(w.ID, u.ID, input.RoleId)
	err = nwau.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, nwau)
}

func RenameWorkspaceName(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// requestのbodyの情報を取得
	var w models.Workspace
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 必要な情報があるか確認
	if w.Name == "" || w.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id or name is empty"})
		return
	}

	// requestしているuserがそのworkspaceのrole = 1 or role = 2 or role = 3かどうかを判定
	wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(w.ID, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if !(wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not permission"})
		return
	}

	// データベースをupdate
	if err := w.RenameWorkspaceName(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, w)
}

func DeleteUserFromWorkSpace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	var wau models.WorkspaceAndUsers
	if err := c.ShouldBindJSON(&wau); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// wauにWorkspaceId, UserId, RoleIdの情報があるかを確認
	if wau.WorkspaceId == 0 || wau.UserId == 0 || wau.RoleId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspaceId or userId or roleId"})
		return
	}

	// requestしたuserがそのworkspaceのrole = 1 or role = 2 or role = 3かどうかチェック
	reqWau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(wau.WorkspaceId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if !(reqWau.RoleId == 1 || reqWau.RoleId == 2 || reqWau.RoleId == 3) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not permission"})
		return
	}

	// 削除されるユーザーがPrimaryOwnerすなわち role = 1でないかチェック
	if wau.RoleId == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not delete primary owner"})
		return
	}

	// wauがdbに存在するかチェック
	if !wau.IsExistWorkspaceAndUser() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspaceAndUser"})
		return
	}

	if err := wau.DeleteWorkspaceAndUser(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wau)
}
