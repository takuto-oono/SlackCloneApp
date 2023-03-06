package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
	"backend/models"
)

func SendDM(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateSendDM(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// requestしたuserがworkspaceに所属しているか確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(in.WorkspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "send user not found in workspace"})
		return
	}
	// receive_userがworkspaceに所属しているか確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(in.WorkspaceId, in.ReceiveUserId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "receive user not found in workspace"})
		return
	}

	// 2人のuserのdm_line_idを取得する(存在しなければ作成する)
	dl, err := models.GetDLByUserIdsAndWorkspaceId(userId, in.ReceiveUserId, in.WorkspaceId)
	if err != nil {
		ndm := models.NewDMLine(in.WorkspaceId, userId, in.ReceiveUserId)
		if err := ndm.Create().Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		dl, err = models.GetDLByUserIdsAndWorkspaceId(userId, in.ReceiveUserId, in.WorkspaceId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// direct_messages tableにデータを保存する
	dm := models.NewDirectMessage(in.Text, userId, dl.ID)
	if err := dm.Create().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dm)
}
