package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
	"backend/models"
)

func CreateWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	primaryOwnerId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateCreateWorkspace(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// requestをしたuserとbodyのprimaryOwnerIdが等しいか確認
	if in.RequestUserId != primaryOwnerId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not equal request user and primary owner id"})
		return
	}

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// dbに保存
	w := models.NewWorkspace(in.Name, in.RequestUserId)
	if err := w.Create(db); err != nil {
		if err.Error() == "UNIQUE constraint failed: workspaces.name" {
      tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
    tx.Rollback()
		return
	}

	// workspace_and_users tableにもuserを保存する
	wau := models.NewWorkspaceAndUsers(w.ID, w.PrimaryOwnerId, 1)
	err = wau.Create(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// general channelを作成する
	ch := models.NewChannel("general", "all users join", false, false, w.ID)
	if err := ch.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// general channelにuserを追加する
	cau := models.NewChannelsAndUses(ch.ID, primaryOwnerId, true)
	if err := cau.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tx.Commit()
	c.IndentedJSON(http.StatusOK, w)
}

func AddUserInWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を受け取る
	in, err := controllerUtils.InputAndValidateAddUserInWorkspace(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// WorkspaceAndUser structを作成
	wau := models.NewWorkspaceAndUsers(in.WorkspaceId, in.UserId, in.RoleId)

	// roleId = 1でないかを確認
	if wau.RoleId == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "can't add roleId = 1"})
		return
	}

	// workspaceが存在するか確認
	if !controllerUtils.IsExistWorkspaceById(wau.WorkspaceId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "workspace not found"})
		return
	}

	// userIdがそのworkspaceで追加する権限を持っているかを判定(roleId == 1 or roleId == 2 or roleId == 3)
	if !controllerUtils.HasPermissionAddUserInWorkspace(userId, wau.WorkspaceId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized add user in workspace"})
		return
	}

	// dbに保存する
	err = wau.Create(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, wau)
}

func RenameWorkspaceName(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// path parameterの値を取得
	workspaceId, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// requestのbodyの情報を取得
	in, err := controllerUtils.InputAndValidateRenameWorkspace(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// requestしているuserがそのworkspaceのrole = 1 or role = 2 or role = 3かどうかを判定
	b, err := controllerUtils.HasPermissionRenamingWorkspaceName(workspaceId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusForbidden, gin.H{"message": "no permission renaming name of workspace"})
		return
	}

	// データベースをupdate
	w, err := models.UpdateWorkspaceName(db, workspaceId, in.WorkspaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, w)
}

func DeleteUserFromWorkSpace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateDeleteUserFromWorkspace(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 削除されるuserがワークスペースに存在するかを確認
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, in.WorkspaceId, in.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// requestしたuserがそのworkspaceのrole = 1 or role = 2 or role = 3かどうかチェック
	b, err := controllerUtils.HasPermissionDeletingUserFromWorkspace(wau.WorkspaceId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permission"})
		return
	}

	// 削除されるユーザーがPrimaryOwnerすなわち role = 1でないかチェック
	if wau.RoleId == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not delete primary owner"})
		return
	}

	if err := wau.DeleteWorkspaceAndUser(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, wau)
}

func GetWorkspacesByUserId(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// workspace structの配列を取得する
	workspaces, err := controllerUtils.GetWorkspacesByUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, workspaces)
}

func GetUsersInWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	workspaceId, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// requestしたuserがworkspaceに存在しているか確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in workspace"})
		return
	}

	// userの情報を取得する
	res, err := controllerUtils.GetUserInWorkspace(workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
