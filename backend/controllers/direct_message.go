package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/controllerUtils"
	"backend/models"
	"backend/utils"
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

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 2人のuserのdm_line_idを取得する(存在しなければ作成する)
	dl, err := models.GetDLByUserIdsAndWorkspaceId(db, userId, in.ReceiveUserId, in.WorkspaceId)
	if err != nil {
		ndm := models.NewDMLine(in.WorkspaceId, userId, in.ReceiveUserId)
		if err := ndm.Create(tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		dl, err = models.GetDLByUserIdsAndWorkspaceId(tx, userId, in.ReceiveUserId, in.WorkspaceId)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// direct_messages tableにデータを保存する
	dm := models.NewDirectMessage(in.Text, userId, dl.ID)
	if err := dm.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, dm)
}

func GetDMsInLine(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからdm_line_idを取得する
	dmLineId, err := utils.StringToUint(c.Param("dm_line_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// dm_lineの情報を取得する
	dl, err := models.GetDLById(db, dmLineId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// dm_lineにrequestしたuserが存在しているか確認
	if !(dl.UserId1 == userId || dl.UserId2 == userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "you don't access this page"})
		return
	}

	// direct_messages tableから情報を取得
	dms, err := models.GetAllDMsByDLId(db, dl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dms)
}

func EditDM(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからdm_idを取得
	dmId, err := utils.StringToUint(c.Param("dm_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyからtextを取得
	in, err := controllerUtils.InputAndValidateEditDM(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// dmが存在するかどうかを確認
	b, err := controllerUtils.IsExistDMById(dmId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusNotFound, gin.H{"message": "dm not found"})
		return
	}

	// 対象のdmがuserが送信したものかを確認
	if !controllerUtils.HasPermissionEditDM(dmId, userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "no permission"})
		return
	}

	// direct_messages tableをupdate
	dm, err := models.UpdateDM(db, dmId, in.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dm)
}

func DeleteDM(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからdm_idを取得
	dmId, err := utils.StringToUint(c.Param("dm_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// dmが存在するかどうかを確認
	b, err := controllerUtils.IsExistDMById(dmId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusNotFound, gin.H{"message": "dm not found"})
		return
	}

	// 対象のdmがuserが送信したものかを確認
	if !controllerUtils.HasPermissionEditDM(dmId, userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "no permission"})
		return
	}

	// direct_messages tableをdelete
	dm, err := models.GetDMById(db, dmId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if err = dm.DeleteDM(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dm)
}
