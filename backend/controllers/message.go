package controllers

import (
	"errors"
	"fmt"
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
	m := models.NewChannelMessage(in.Text, in.ChannelId, userId)

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

	if err = controllerUtils.SendMessageTX(m, userId, in.MentionedUserIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
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

	// channelにuserが所属していることを確認
	if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(channelId, userId); !b || err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found in channel"})
		return
	}

	// DBからデータを取得
	messages, err := models.GetMessagesByChannelId(db, channelId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
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
	dm := models.NewDMMessage(in.Text, dl.ID, userId)
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

	c.JSON(http.StatusOK, dms)
}

type DMLineInfo struct {
	ID     uint   `json:"id"`
	ToName string `json:"to_name"`
}

func GetDMLines(c *gin.Context) {
	var response []DMLineInfo

	setToName := func(userId1, userId2, requestUserId uint32) (string, error) {
		var toId uint32
		if userId1 == requestUserId && userId2 == requestUserId {
			toId = requestUserId
		} else if userId1 == requestUserId && userId2 != requestUserId {
			toId = userId2
		} else if userId1 != requestUserId && userId2 == requestUserId {
			toId = userId1
		} else {
			return "", fmt.Errorf("request_user not found")
		}
		u, err := models.GetUserById(db, toId)
		if err != nil {
			return "", err
		}
		return u.Name, nil
	}

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

	// 取得した情報をDMLineInfoに変換する
	for _, dl := range dls {
		toName, err := setToName(dl.UserId1, dl.UserId2, userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		response = append(response, DMLineInfo{
			ID:     dl.ID,
			ToName: toName,
		})
	}

	c.JSON(http.StatusOK, response)
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

func ChannelSocket(hub *Hub, ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	userID, err := Authenticate(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// urlからchannel_idを取得
	channelID, err := strconv.Atoi(ctx.Param("channel_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// channelの情報を取得
	ch, err := models.GetChannelById(db, channelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	go hub.run()
	fmt.Println("in channelSocket func")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	fmt.Println("upgrader setting")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	client := &Client{hub: hub, conn: conn, send: make(chan models.Message), channelID: ch.ID, userID: userID}
	fmt.Println("client: ", client)
	client.hub.register <- client
	fmt.Println(hub.clients)

	go client.readPump(userID)
	go client.writePump()
}
