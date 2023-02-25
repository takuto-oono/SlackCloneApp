package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
	"backend/models"
)

func SendMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateSendMessage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// message structを作成
	m := models.NewMessage(in.Text, in.ChannelId, userId)

	// userとchannelが同じworkspaceに存在しているかを確認
	if b, err := controllerUtils.IsExistChannelAndUserInSameWorkspace(m.ChannelId, userId); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not exist channel and user in same workspace"})
		return
	}

	// channelにuserが参加しているかを確認
	if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(m.ChannelId, userId); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not exist user in channel"})
		return
	}

	// message情報をDBに登録
	if err := m.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

func GetAllMessagesFromChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// path parameterからchannel_idを取得する
	channelId, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// channelにuserが所属していることを確認
	if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(channelId, userId); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not exist user in channel"})
		return
	}

	// DBからデータを取得
	messages, err := models.GetMessagesByChannelId(channelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}
