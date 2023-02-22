package controllers

import (
	"backend/controllerUtils"
	"backend/models"
	"backend/utils"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

var messageRouter = SetupRouter()

func sendMessageTestFunc(text string, channelId int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.SendMessageInput{
		Text:      text,
		ChannelId: channelId,
	})
	req, err := http.NewRequest("POST", "/api/message/send", bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	messageRouter.ServeHTTP(rr, req)
	return rr
}

func getMessagesByChannelIdTestFunc(channelId int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/message/get_from_channel/"+strconv.Itoa(channelId), nil)
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	messageRouter.ServeHTTP(rr, req)
	return rr
}

func TestSendMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. bodyに不足がある場合 400
	// 3. userとchannelが同じworkspaceに存在していない場合 400
	// 4. channelにuserが存在しない場合 400

	t.Run("1 正常な場合", func(t *testing.T) {
		testCaseNum := "1"
		userName := "testSendMessageUserName" + testCaseNum
		workspaceName := "testSendMessageWorkspaceName" + testCaseNum
		channelName := "testSendMessageChannelName" + testCaseNum
		text := "testSendMessageText" + testCaseNum
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		rr = sendMessageTestFunc(text, ch.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)
		assert.NotEmpty(t, m.ID)
		assert.Equal(t, text, m.Text)
		assert.NotEmpty(t, m.Date)
		assert.Equal(t, ch.ID, m.ChannelId)
		assert.Equal(t, lr.UserId, m.UserId)
	})

	t.Run("2 bodyに不足がある場合", func(t *testing.T) {
		testCaseNum := "2"
		userName := "testSendMessageUserName" + testCaseNum
		workspaceName := "testSendMessageWorkspaceName" + testCaseNum
		channelName := "testSendMessageChannelName" + testCaseNum
		text := "testSendMessageText" + testCaseNum
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		rr = sendMessageTestFunc(text, 0, lr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found channel_id\"}", rr.Body.String())

		rr = sendMessageTestFunc("", ch.ID, lr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found text\"}", rr.Body.String())
	})

	t.Run("3 userとchannelが同じworkspaceに存在していない場合", func(t *testing.T) {
		testCaseNum := "3"
		userName := "testSendMessageUserName" + testCaseNum
		userName2 := "testSendMessageUserName" + testCaseNum + ".2"
		workspaceName := "testSendMessageWorkspaceName" + testCaseNum
		workspaceName2 := "testSendMessageWorkspaceName" + testCaseNum + ".2"
		channelName := "testSendMessageChannelName" + testCaseNum
		text := "testSendMessageText" + testCaseNum
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(userName2, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = loginTestFunc(userName2, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		lr2 := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr2)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createWorkSpaceTestFunc(workspaceName2, lr2.Token, lr2.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w2 := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w2)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		rr = sendMessageTestFunc(text, ch.ID, lr2.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not exist channel and user in same workspace\"}", rr.Body.String())
	})
	t.Run("4 channelにuserが存在しない場合", func(t *testing.T) {
		testCaseNum := "4"
		userName := "testSendMessageUserName" + testCaseNum
		userName2 := "testSendMessageUserName" + testCaseNum + ".2"
		workspaceName := "testSendMessageWorkspaceName" + testCaseNum
		channelName := "testSendMessageChannelName" + testCaseNum
		text := "testSendMessageText" + testCaseNum
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(userName2, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = loginTestFunc(userName2, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		lr2 := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr2)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, lr2.UserId, lr.Token).Code)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), ch)

		rr = sendMessageTestFunc(text, ch.ID, lr2.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not exist user in channel\"}", rr.Body.String())
	})
}

func TestGetAllMessagesFromChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. messageが存在する場合 200
	// 2. messageが存在しない場合 200
	// 3. userがchannelに所属していない場合 400

	t.Run("1 messageが存在する場合", func(t *testing.T) {
		testNum := "1"
		userName := "testGetAllMessageFromChannelUserName" + testNum
		workspaceName := "testGetAllMessageFromChannelWorkspaceName" + testNum
		channelName := "testGetAllMessageFromChannelName" + testNum
		isPrivate := true
		text := "testGetAllMessageFromChannelText" + testNum
		messageCount := 1000

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		for i := 0; i < messageCount; i++ {
			assert.Equal(t, http.StatusOK, sendMessageTestFunc(text, ch.ID, lr.Token).Code)
		}

		rr = getMessagesByChannelIdTestFunc(ch.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		messages := make([]models.Message, messageCount)
		json.Unmarshal(([]byte)(byteArray), &messages)
		assert.Equal(t, messageCount, len(messages))

		for i := 0; i < messageCount-1; i++ {
			d1, err1 := utils.TimeFromString(messages[i].Date)
			d2, err2 := utils.TimeFromString(messages[i+1].Date)
			assert.Empty(t, err1)
			assert.Empty(t, err2)
			assert.True(t, d2.Before(d1))
		}

		for _, m := range messages {
			assert.Equal(t, text, m.Text)
			assert.Equal(t, lr.UserId, m.UserId)
			assert.Equal(t, ch.ID, m.ChannelId)
		}
	})

	t.Run("2 messageが存在しない場合", func(t *testing.T) {
		testNum := "2"
		userName := "testGetAllMessageFromChannelUserName" + testNum
		workspaceName := "testGetAllMessageFromChannelWorkspaceName" + testNum
		channelName := "testGetAllMessageFromChannelName" + testNum
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		rr = getMessagesByChannelIdTestFunc(ch.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		messages := make([]models.Message, 0)
		json.Unmarshal(([]byte)(byteArray), &messages)
		assert.Equal(t, 0, len(messages))
	})

	t.Run("3 userがchannelに所属していない場合", func(t *testing.T) {
		testNum := "3"
		userName := "testGetAllMessageFromChannelUserName" + testNum
		userName2 := "testGetAllMessageFromChannelUserName" + testNum + ".2"
		workspaceName := "testGetAllMessageFromChannelWorkspaceName" + testNum
		channelName := "testGetAllMessageFromChannelName" + testNum
		isPrivate := true
		text := "testGetAllMessageFromChannelText" + testNum
		messageCount := 3

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(userName2, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = loginTestFunc(userName2, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		lr2 := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr2)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		for i := 0; i < messageCount; i++ {
			assert.Equal(t, http.StatusOK, sendMessageTestFunc(text, ch.ID, lr.Token).Code)
		}

		rr = getMessagesByChannelIdTestFunc(ch.ID, lr2.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not exist user in channel\"}", rr.Body.String())
	})
}
