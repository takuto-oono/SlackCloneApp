package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/controllerUtils"
	"backend/models"
	"backend/utils"
)

func SendMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidateSendMessage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// message structを作成
	m := models.NewChannelMessage(in.Text, in.ChannelId, userId, in.ScheduleTime)

	// userとchannelが同じworkspaceに存在しているかを確認
	if b, err := controllerUtils.IsExistChannelAndUserInSameWorkspace(m.ChannelId, userId); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "channel and user not found in same workspace"})
		return
	}

	// channelにuserが参加しているかを確認
	if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(m.ChannelId, userId); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in channel"})
		return
	}

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// message情報をDBに登録
	if err := m.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// mentionの処理をする
	for _, userID := range in.MentionedUserIDs {
		men := models.NewMention(userID, m.ID)
		if err := men.Create(tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// channelにいるuserをmessage_and_users tableに追加する
	caus, err := models.GetCAUsByChannelId(tx, m.ChannelId)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	for _, cau := range caus {
		if cau.UserId == userId {
			continue
		}
		mau := models.NewMessageAndUser(m.ID, cau.UserId, false)
		if err := mau.Create(tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, m)
}

func GetAllMessagesFromChannel(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// path parameterからchannel_idを取得する
	channelId, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !controllerUtils.HasPermissionGetMessagesFromChannel(channelId, userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "not forbidden"})
		return
	}

	// DBからデータを取得
	messages, err := models.GetMessagesByChannelId(db, channelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	createResponse := func(messages []models.Message) []models.Message {
		messages = controllerUtils.FilterByFutureScheduleTimeOfMessages(messages)
		messages = controllerUtils.UpdateCreatedAt(messages)
		return controllerUtils.SortMessageByCreatedAt(messages)
	}

	c.JSON(http.StatusOK, createResponse(messages))
}

func EditMessage(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// path parameterからmessage_idを取得する
	messageId, err := utils.StringToUint(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// bodyからtextを取得する
	in, err := controllerUtils.InputAndValidateEditMessage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// 対象のmessageが存在するか確認する
	b, err := controllerUtils.IsExistMessageById(messageId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		c.JSON(http.StatusNotFound, gin.H{"message": "message not found"})
		return
	}

	// messageの編集権限があるかを確認
	if !controllerUtils.HasPermissionEditMessage(messageId, userId) {
		c.JSON(http.StatusForbidden, gin.H{"message": "no permission"})
		return
	}

	// DBを更新
	m, err := models.UpdateMessageText(db, messageId, in.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

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
	dm := models.NewDMMessage(in.Text, dl.ID, userId, in.ScheduleTime)
	if err := dm.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// mentionの処理をする
	for _, userID := range in.MentionedUserIDs {
		men := models.NewMention(userID, dm.ID)
		if err := men.Create(tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// message and users tableに保存する
	mau := models.NewMessageAndUser(dm.ID, in.ReceiveUserId, false)
	if err := mau.Create(tx); err != nil {
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
	dms, err := models.GetMessagesByDLId(db, dl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	createResponse := func(messages []models.Message) []models.Message {
		messages = controllerUtils.FilterByFutureScheduleTimeOfMessages(messages)
		messages = controllerUtils.UpdateCreatedAt(messages)
		return controllerUtils.SortMessageByCreatedAt(messages)
	}

	c.JSON(http.StatusOK, createResponse(dms))
}

func GetDMLines(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからworkspace_idを取得する
	workspaceId, err := strconv.Atoi(c.Param("workspace_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// userがworkspaceに存在しているかを確認
	if !controllerUtils.IsExistWAUByWorkspaceIdAndUserId(workspaceId, userId) {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in workspace"})
		return
	}

	// userが所属しているDMLineをすべて取得
	dls, err := models.GetDLsByUserIdAndWorkspaceId(db, userId, workspaceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dls)
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
	b, err := controllerUtils.IsExistMessageById(dmId)
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
	dm, err := models.UpdateMessageText(db, dmId, in.Text)
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
	b, err := controllerUtils.IsExistMessageById(dmId)
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
	dm, err := models.GetMessageById(db, dmId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if err = dm.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dm)
}

func ReadMessageByUser(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userID, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからdm_idを取得
	messageID, err := utils.StringToUint(c.Param("message_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// messageが存在するかどうかを確認
	b, err := controllerUtils.IsExistMessageById(messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if !b {
		c.JSON(http.StatusNotFound, gin.H{"message": "message not found"})
		return
	}

	// message_and_users tableにデータがあるか確認、あれば更新
	mau, err := models.GetMAUByMessageIDAndUserID(db, messageID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := mau.UpdateIsRead(db, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mau)
}
