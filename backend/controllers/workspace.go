package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/models"
)

func CreateWorkspace(c *gin.Context) {
	fmt.Println("in func")
	c.Header("Access-Control-Allow-Origin", "*")
	primaryOwnerId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// bodyの情報を取得
	var w models.Workspace
	if err := c.ShouldBindJSON(&w); err != nil {
		fmt.Println("err")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// workspaceのnameがあるか確認
	if w.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspace name"})
		return
	}

	// requestをしたuserとbodyのprimaryOwnerIdが等しいか確認
	if w.PrimaryOwnerId != primaryOwnerId || w.PrimaryOwnerId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not permission"})
		return
	}

	// はじめはidを0にしておく
	w.ID = 0

	// dbに保存
	if err := w.CreateWorkspace(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// workspace_and_users tableにもuserを保存する
	wau := models.NewWorkspaceAndUsers(w.ID, w.PrimaryOwnerId, 1)
	err = wau.Create()
	if err != nil {
		// TODO deleteWorkspaceを実行する
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, w)
}

func havePermissionAddUserInWorkspace(userId uint32, workspaceId int) bool {
	wau, err := models.GetWorkspaceAndUserByWorkspaceIdAndUserId(workspaceId, userId)
	if err != nil {
		return false
	}
	return wau.RoleId == 1 || wau.RoleId == 2 || wau.RoleId == 3
}

func AddUserInWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を受け取る
	var wau models.WorkspaceAndUsers
	if err := c.ShouldBindJSON(&wau); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if wau.WorkspaceId == 0 || wau.UserId == 0 || wau.RoleId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "field empty"})
		return
	}

	// roleId = 1でないかを確認
	if wau.RoleId == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't add roleId = 1"})
		return
	}

	// workspaceが存在するか確認
	if !models.IsExistWorkspaceById(wau.WorkspaceId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspace"})
		return
	}

	// userIdがそのworkspaceで追加する権限を持っているかを判定(roleId == 1 or roleId == 2 or roleId == 3)
	if !havePermissionAddUserInWorkspace(userId, wau.WorkspaceId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unauthorized add user in workspace"})
		return
	}

	// dbに保存する
	err = wau.Create()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, wau)
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
