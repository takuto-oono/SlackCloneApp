package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
)

func GetMessagesMentionedByUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
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

	// userがworkspaceに存在しているかを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceID, userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permission"})
		return
	}

	// userがmentionされたメッセージを取得する
	messages, err := controllerUtils.GetMessagesMentionedByUserAndWorkspace(userId, workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	for _, mes := range messages {
		fmt.Println(mes)
	}
	c.JSON(http.StatusOK, messages)
}
