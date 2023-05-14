package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"

	"backend/controllerUtils"
	"backend/models"
	"backend/utils"
)

var messageRouter = SetupRouter()

func sendMessageTestFunc(text string, channelId int, jwtToken string, mentionedUserIDs []uint32) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.SendMessageInput{
		Text:             text,
		ChannelId:        channelId,
		MentionedUserIDs: mentionedUserIDs,
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

func editMessageTestFunc(messageId uint, text, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.EditMessageInput{
		Text: text,
	})
	req, err := http.NewRequest("PATCH", "/api/message/edit/"+utils.UintToString(messageId), bytes.NewBuffer(jsonInput))
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
	// 3. userとchannelが同じworkspaceに存在していない場合 404
	// 4. channelにuserが存在しない場合 404
	// 5. mentionがある場合 200

	t.Run("1 正常な場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(30)
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

		rr = sendMessageTestFunc(text, ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)
		assert.NotEmpty(t, m.ID)
		assert.Equal(t, text, m.Text)
		assert.NotEmpty(t, m.CreatedAt)
		assert.Equal(t, ch.ID, m.ChannelId)
		assert.Equal(t, lr.UserId, m.UserId)
	})

	t.Run("2 bodyに不足がある場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(30)
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

		rr = sendMessageTestFunc(text, 0, lr.Token, []uint32{})
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"channel_id not found\"}", rr.Body.String())

		rr = sendMessageTestFunc("", ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"text not found\"}", rr.Body.String())
	})

	t.Run("3 userとchannelが同じworkspaceに存在していない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		userName2 := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		workspaceName2 := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(30)
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

		rr = sendMessageTestFunc(text, ch.ID, lr2.Token, []uint32{})
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"channel and user not found in same workspace\"}", rr.Body.String())
	})
	t.Run("4 channelにuserが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		userName2 := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(30)
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

		rr = sendMessageTestFunc(text, ch.ID, lr2.Token, []uint32{})
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"user not found in channel\"}", rr.Body.String())
	})

	t.Run("5 メンションがある場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(30)
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

		rr = sendMessageTestFunc(text, ch.ID, lr.Token, []uint32{lr.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)
		assert.NotEmpty(t, m.ID)
		assert.Equal(t, text, m.Text)
		assert.NotEmpty(t, m.CreatedAt)
		assert.Equal(t, ch.ID, m.ChannelId)
		assert.Equal(t, lr.UserId, m.UserId)
	})

}

func TestGetAllMessagesFromChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. messageが存在する場合 200
	// 2. messageが存在しない場合 200
	// 3. userがchannelに所属していない場合 404

	t.Run("1 messageが存在する場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		isPrivate := true
		text := randomstring.EnglishFrequencyString(30)
		messageCount := 10

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
			assert.Equal(t, http.StatusOK, sendMessageTestFunc(text, ch.ID, lr.Token, []uint32{}).Code)
		}

		rr = getMessagesByChannelIdTestFunc(ch.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		messages := make([]models.Message, messageCount)
		json.Unmarshal(([]byte)(byteArray), &messages)
		assert.Equal(t, messageCount, len(messages))

		for i := 0; i < messageCount-1; i++ {
			d1 := messages[i].CreatedAt
			d2 := messages[i+1].CreatedAt
			assert.True(t, d2.Before(d1))
		}

		for _, m := range messages {
			assert.Equal(t, text, m.Text)
			assert.Equal(t, lr.UserId, m.UserId)
			assert.Equal(t, ch.ID, m.ChannelId)
		}
	})

	t.Run("2 messageが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
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
		userName := randomstring.EnglishFrequencyString(30)
		userName2 := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		isPrivate := true
		text := randomstring.EnglishFrequencyString(30)
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
			assert.Equal(t, http.StatusOK, sendMessageTestFunc(text, ch.ID, lr.Token, []uint32{}).Code)
		}

		rr = getMessagesByChannelIdTestFunc(ch.ID, lr2.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"user not found in channel\"}", rr.Body.String())
	})
}

func TestEditMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. bodyにtextがない場合 400
	// 3. 対象のmessageが存在しない場合 404
	// 4. 別のuserが作成したmessageだった場合 403

	t.Run("1 正常な場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		oldText := randomstring.EnglishFrequencyString(100)
		newText := randomstring.EnglishFrequencyString(100)
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

		rr = sendMessageTestFunc(oldText, ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)

		rr = editMessageTestFunc(m.ID, newText, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		res := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), res)

		assert.Equal(t, m.ID, res.ID)
		assert.Equal(t, newText, res.Text)
		assert.Equal(t, m.ChannelId, res.ChannelId)
		assert.Equal(t, lr.UserId, res.UserId)
		assert.Equal(t, m.CreatedAt, res.CreatedAt)
		assert.True(t, m.UpdatedAt.Before(res.UpdatedAt))
	})

	t.Run("2 bodyにtextがない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		oldText := randomstring.EnglishFrequencyString(100)
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

		rr = sendMessageTestFunc(oldText, ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)

		rr = editMessageTestFunc(m.ID, "", lr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"text nof found\"}", rr.Body.String())
	})

	t.Run("3 対象のmessageが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		newText := randomstring.EnglishFrequencyString(100)
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

		rr = editMessageTestFunc(uint(rand.Uint32()), newText, lr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"message not found\"}", rr.Body.String())
	})

	t.Run("4 別のuserが作成したmessageだった場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		requestUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		oldText := randomstring.EnglishFrequencyString(100)
		newText := randomstring.EnglishFrequencyString(100)
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, lr.Token).Code)

		rr = createChannelTestFunc(channelName, "", &isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)

		assert.Equal(t, http.StatusOK, addUserInChannelTestFunc(ch.ID, rlr.UserId, lr.Token).Code)

		rr = sendMessageTestFunc(oldText, ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)

		rr = editMessageTestFunc(m.ID, newText, rlr.Token)
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"no permission\"}", rr.Body.String())
	})
}

func TestReadMessageByUser(t *testing.T) {
	// if testing.Short() {
	// 	t.Skip("skipping test in short mode.")
	// }

	// 1 channelのmessageの場合 200
	// 2 dmのmessageの場合 200

	t.Run("1 channelのmessageの場合 200", func(t *testing.T) {
		rr, u1 := signUpTestFuncV2(randomstring.EnglishFrequencyString(30), "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, u2 := signUpTestFuncV2(randomstring.EnglishFrequencyString(30), "pass")
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, lr1 := loginTestFuncV2(u1.Name, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, lr2 := loginTestFuncV2(u2.Name, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, w := createWorkspaceTestFuncV2(randomstring.EnglishFrequencyString(30), lr1.Token, lr1.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, _ = addUserInWorkspaceV2(w.ID, lr2.UserId, 4, lr1.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		isPrivate := false
		rr, ch := createChannelTestFuncV2(randomstring.EnglishFrequencyString(30), "", &isPrivate, lr1.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, _ = addUserInChannelTestFuncV2(ch.ID, lr2.UserId, lr1.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, m := sendMessageTestFuncV2(randomstring.EnglishFrequencyString(30), ch.ID, lr1.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, _ = readMessageByUserTestFunc(m.ID, lr1.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)

		rr, mau2 := readMessageByUserTestFunc(m.ID, lr2.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, m.ID, mau2.MessageID)
		assert.Equal(t, lr2.UserId, mau2.UserID)
		assert.True(t, mau2.IsRead)
	})

	t.Run("2 dmのmessageの場合 200", func(t *testing.T) {
		rr, u1 := signUpTestFuncV2(randomstring.EnglishFrequencyString(30), "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, u2 := signUpTestFuncV2(randomstring.EnglishFrequencyString(30), "pass")
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, lr1 := loginTestFuncV2(u1.Name, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, lr2 := loginTestFuncV2(u2.Name, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, w := createWorkspaceTestFuncV2(randomstring.EnglishFrequencyString(30), lr1.Token, lr1.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, _ = addUserInWorkspaceV2(w.ID, lr2.UserId, 4, lr1.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		isPrivate := false
		rr, ch := createChannelTestFuncV2(randomstring.EnglishFrequencyString(30), "", &isPrivate, lr1.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, _ = addUserInChannelTestFuncV2(ch.ID, lr2.UserId, lr1.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, m := sendDMTestFuncV2(randomstring.EnglishFrequencyString(30), lr1.Token, lr2.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, _ = readMessageByUserTestFunc(m.ID, lr1.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)

		rr, mau2 := readMessageByUserTestFunc(m.ID, lr2.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, m.ID, mau2.MessageID)
		assert.Equal(t, lr2.UserId, mau2.UserID)
		assert.True(t, mau2.IsRead)
	})
}
