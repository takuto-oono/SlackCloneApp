package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"

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

	fmt.Println("1")

	// threadの元になるmessageを取得
	parentMessage, err := models.GetMessageById(db, in.ParentMessageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("2")

	// トランザクションを宣言
	tx := db.Begin()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("3")

	// threadを取得 or 作成
	th, err := controllerUtils.CreateOrGetThreadByParentMessageId(tx, parentMessage.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("4")

	// thread tableのupdated_atを更新
	if err := th.EditUpdatedAt(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("5")

	// 新しいmessageを作成する
	var m *models.Message
	if parentMessage.ChannelId != 0 && parentMessage.DMLineId == uint(0) {
		// channelとuserが同じworkspaceに存在しているか確認
		if b, err := controllerUtils.IsExistChannelAndUserInSameWorkspace(parentMessage.ChannelId, userId); !b || err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "channel and user not found in same workspace"})
			return
		}

		// channelにuserが参加しているかを確認
		if b, err := controllerUtils.IsExistCAUByChannelIdAndUserId(parentMessage.ChannelId, userId); !b || err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found in channel"})
			return
		}

		// message structを作成
		m = models.NewChannelMessage(in.Text, parentMessage.ChannelId, userId)
	} else if parentMessage.ChannelId == 0 && parentMessage.DMLineId != uint(0) {
		// parentMessageが存在するDMLineにuserが参加しているかを確認
		b, err := controllerUtils.IsExistUserInDL(userId, parentMessage.DMLineId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !b {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found in dm_line"})
			return
		}

		// message structを作成
		m = models.NewDMMessage(in.Text, parentMessage.DMLineId, userId)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "wrong channel_id or dm_line_id"})
		return
	}
	fmt.Println("6")
	
	// messageをdbに登録
	m.ThreadId = th.ID
	if m.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("7")

	// thread_and_message tableにデータを保存
	tam := models.NewThreadAndMessage(th.ID, m.ID)
	if err := tam.Create(tx); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	fmt.Println("8")
	ptam := models.NewThreadAndMessage(th.ID, parentMessage.ID)
	if err := ptam.Create(tx); err != nil {
		if !errors.Is(err, sqlite3.ErrConstraintUnique) {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	fmt.Println("9")

	// threads_and_users tableにユーザーを追加
	tau := models.NewThreadAndUser(userId, th.ID)
	if err := tau.Create(tx); err != nil {
		fmt.Println(err)
		if !errors.As(err, sqlite3.ErrConstraintUnique) {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}
	fmt.Println("10")

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
