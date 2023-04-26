package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/controllerUtils"
	"backend/models"
)

func PostThread(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// bodyの情報を取得
	in, err := controllerUtils.InputAndValidatePostThread(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// threadの元になるmessageを取得
	parentMessage, err := models.GetMessageById(db, in.ParentMessageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// スレッド内のメッセージだった場合は返信できないようにする
	if parentMessage.ThreadId != uint(0) {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// threadを取得 or 作成
	th, err := controllerUtils.CreateOrGetThreadByParentMessageId(tx, parentMessage.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// thread tableのupdated_atを更新
	if err := th.EditUpdatedAt(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 新しいmessageを作成する
	var m *models.Message
	if parentMessage.ChannelId != 0 && parentMessage.DMLineId == uint(0) {
		// channelとuserが同じworkspaceに存在しているか確認
		if b, err := controllerUtils.IsExistChannelAndUserInSameWorkspace(parentMessage.ChannelId, userId); !b || err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"message": "channel and user not found in same workspace"})
			return
		}

		// channelにuserが参加しているかを確認
		if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(parentMessage.ChannelId, userId); !b || err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found in channel"})
			return
		}

		// message structを作成
		m = models.NewChannelMessage(in.Text, parentMessage.ChannelId, userId)
	} else if parentMessage.ChannelId == 0 && parentMessage.DMLineId != uint(0) {
		// parentMessageが存在するDMLineにuserが参加しているかを確認
		b, err := controllerUtils.IsExistUserInDL(userId, parentMessage.DMLineId)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !b {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found in dm_line"})
			return
		}

		// message structを作成
		m = models.NewDMMessage(in.Text, parentMessage.DMLineId, userId)
	} else {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "wrong channel_id or dm_line_id"})
		return
	}
	
	// messageをdbに登録
	m.ThreadId = th.ID
	if m.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// thread_and_message tableにデータを保存
	tam := models.NewThreadAndMessage(th.ID, m.ID)
	if err := tam.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	b, err := controllerUtils.IsExistTAMByThreadIdAndMessageId(th.ID, parentMessage.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		ptam := models.NewThreadAndMessage(th.ID, parentMessage.ID)
		if err := ptam.Create(tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// threads_and_users tableにユーザーを追加
	b, err = controllerUtils.IsExistTAUByUserIdAndThreadId(tx, userId, th.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if !b {
		tau := models.NewThreadAndUser(userId, th.ID)
		if err := tau.Create(tx); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, m)
}

func GetThreadsByUser(c *gin.Context) {
	type ThreadInfo struct {
		ID       uint
		Messages []models.Message
	}
	var res []ThreadInfo

	c.Header("Access-Control-Allow-Origin", "*")
	userId, err := Authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	// userが所属しているthreadを取得
	ths, err := controllerUtils.GetThreadsByUserSortedByEditedTime(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// responseを作成
	for _, th := range ths {
		messages, err := models.GetMessagesByThreadId(db, th.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		res = append(res, ThreadInfo{
			ID: th.ID,
			Messages: messages,
		})
	}

	c.JSON(http.StatusOK, res)
}
