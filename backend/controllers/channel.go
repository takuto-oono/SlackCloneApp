package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
	"backend/models"
)

func CreateChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateCreateChannel(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ch := models.NewChannel(in.Name, in.Description, *in.IsPrivate, false, in.WorkspaceId)

	// workspaceIdに対応するworkspaceが存在するか確認
	if !controllerUtils.IsExistWorkspaceById(ch.WorkspaceId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "workspace not found"})
		return
	}

	// 同じ名前のchannelが対応するworkspaceに存在しないか確認
	b, err := controllerUtils.IsExistSameNameChannelInWorkspace(ch.Name, ch.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if b {
		c.JSON(http.StatusConflict, gin.H{"message": "already exist same name channel in workspace"})
		return
	}

	// userが対象のworkspaceに参加しているか確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(ch.WorkspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in workspace"})
		return
	}

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// channels tableに情報を保存
	if err := ch.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// channels_and_users tableに保存する情報を作成し保存
	cau := models.NewChannelsAndUses(ch.ID, userId, true)
	if err := cau.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, ch)
}

func AddUserInChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateAddUserInChannel(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// とりあえず管理権限はなし
	cau := models.NewChannelsAndUses(in.ChannelId, in.UserId, false)

	// channelを取得
	ch, err := models.GetChannelById(db, cau.ChannelId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// リクエストしたuserがworkspaceに参加してるかを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(ch.WorkspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "request user not found in workspace"})
		return
	}

	// 追加されるuserがworkspaceに参加しているかを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(ch.WorkspaceId, cau.UserId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "added user not found in workspace"})
		return
	}

	// 対象のchannelがworkspace内に存在するかを確認
	b, err := controllerUtils.IsExistChannelByChannelIdAndWorkspaceId(cau.ChannelId, ch.WorkspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusNotFound, gin.H{"message": "channel not found in workspace"})
		return
	}

	// 追加されるユーザーが既に対象のchannelに存在していないかを確認
	b, err = controllerUtils.IsExistCAUByChannelIdAndUserId(cau.ChannelId, cau.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if b {
		c.JSON(http.StatusConflict, gin.H{"message": "already exist user in channel"})
		return
	}

	// リクエストしたuserにchannelの管理権限があるかを確認(結果的にリクエストしたuserがchannelに所属しているかも確認される)
	if !controllerUtils.HasPermissionAddingUserInChannel(cau.ChannelId, userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "no permission adding user in channel"})
		return
	}

	// channels_and_users tableに登録
	if err := cau.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cau)
}

func DeleteUserFromChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからパラメータを取得
	workspaceId, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyを取得
	var cau models.ChannelsAndUsers
	if err := c.ShouldBindJSON(&cau); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyに必要な情報があるかを確認
	if cau.UserId == 0 || cau.ChannelId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "user_id or channel_id not found"})
		return
	}

	// requestしたuserがworkspaceにいることを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "request user not found in workspace"})
		return
	}

	// deleteされるuserがworkspaceにいることを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceId, cau.UserId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in workspace"})
		return
	}

	// channelの情報を取得
	ch, err := models.GetChannelById(db, cau.ChannelId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// channelがworkspaceに存在することを確認
	if ch.WorkspaceId != workspaceId {
		c.JSON(http.StatusNotFound, gin.H{"message": "channel not found in workspace"})
		return
	}

	// channelのnameがgeneralでないことを確認
	if ch.Name == "general" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "don't delete general channel"})
		return
	}

	// channelがアーカイブされていないことを確認
	if ch.IsArchive {
		c.JSON(http.StatusBadRequest, gin.H{"message": "don't delete archived channel"})
		return
	}

	// deleteされるuserがchannelに存在することを確認
	b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(cau.ChannelId, cau.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in channel"})
		return
	}

	// deleteする権限があるかを確認
	if !controllerUtils.HasPermissionDeletingUserInChannel(userId, workspaceId, ch) {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permission deleting user in channel"})
		return
	}

	// channels_and_users tableから削除
	if err := cau.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cau)
}

func DeleteChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	var ch models.Channel
	if err := c.ShouldBindJSON(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報に不足がないか確認
	if ch.ID == 0 || ch.WorkspaceId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "channel_id or workspace_id not found"})
		return
	}

	// requestしたuserがworkspaceに参加しているかを確認
	wau, err := models.GetWAUByWorkspaceIdAndUserId(db, ch.WorkspaceId, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// deleteする権限があるかを確認
	if !controllerUtils.HasPermissionDeletingChannel(wau) {
		c.JSON(http.StatusForbidden, gin.H{"message": "no permission deleting channel"})
		return
	}

	// channelがworkspaceにあるかどうかを確認
	if _, err := models.GetChannelByIdAndWorkspaceId(db, ch.ID, ch.WorkspaceId); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// channels tableからデータを削除
	if err := ch.Delete(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// channel_and_users tableからデータを削除
	if err := models.DeleteCAUsByChannelId(tx, ch.ID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, ch)
}

func GetChannelsByUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからworkspace_idを取得
	workspaceId, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// userがworkspaceに存在しているかを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "request user not found in workspace"})
		return
	}

	// channelの配列を取得
	chs, err := controllerUtils.GetChannelsByUserIdAndWorkspaceId(userId, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chs)
}

func GetChannelsByWorkspace(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userID, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからworkspace_idを取得
	workspaceID, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// userがworkspaceに存在しているか確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceID, userID) {
		c.JSON(http.StatusNotFound, gin.H{"message": "workspace or user not found"})
		return
	}

	chs, err := controllerUtils.GetPublicAndJoinPrivateChannelInWorkspace(workspaceID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chs)
}

func GetAllUsersInChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userID, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからworkspace_idを取得
	channelID, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// userとchannelが同じworkspaceに存在するかを
	if b, err := controllerUtils.IsExistChannelAndUserInSameWorkspace(channelID, userID); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "channel and user not found in same workspace"})
		return
	}

	caus, err := models.GetCAUsByChannelId(db, channelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	res := make([]UserResponse, len(caus))
	for i, cau := range caus {
		u, err := models.GetUserById(db, cau.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		res[i] = UserResponse{
			ID:   u.ID,
			Name: u.Name,
		}
	}
	c.JSON(http.StatusOK, res)
}
