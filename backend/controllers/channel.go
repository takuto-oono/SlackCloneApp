package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/models"
)

func CreateChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// urlからパラメータを取得
	workspaceId, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	var ch models.Channel
	if err := c.ShouldBindJSON(&ch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// channel nameがbodyに含まれているか確認
	if ch.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found channel name"})
		return
	}

	// workspaceIdに対応するworkspaceが存在するか確認
	if !models.IsExistWorkspaceById(workspaceId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found workspace"})
		return
	}

	// 同じ名前のchannelが対応するworkspaceに存在しないか確認
	b, err := ch.IsExistSameNameChannelInWorkspace(workspaceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if b {
		c.JSON(http.StatusBadRequest, gin.H{"message": "already exist same name channel in workspace"})
		return
	}

	// userが対象のworkspaceに参加しているか確認
	wau := models.NewWorkspaceAndUsers(workspaceId, userId, 0)
	if !wau.IsExistWorkspaceAndUser() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not found user in workspace"})
		return
	}

	// channels tableに情報を保存
	if err := ch.CreateChannel(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// channels_and_workspaces tableに保存する情報を作成し保存
	caw := models.NewChannelsAndWorkspaces(ch.ID, workspaceId)
	if err := caw.CreateChannelsAndWorkspaces(); err != nil {
		// TODO DeleteChannelsAndWorkspaces funcを実行する
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// channels_and_users tableに保存する情報を作成し保存
	cau := models.NewChannelsAndUses(ch.ID, userId, true)
	if err := cau.CreateChannelAndUsers(); err != nil {
		// TODO DeleteChannel funcを実行する
		// TODO DeleteChannelsAndWorkspace funcを実行

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ch)
}
