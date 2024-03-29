package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"

	"backend/controllerUtils"
	"backend/models"
)

var threadRouter = SetupRouter1()

func postThreadTestFunc(text, jwtToken string, parentMessageId uint, mentionedUserIDs []uint32) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.PostThreadInput{
		Text:             text,
		ParentMessageId:  parentMessageId,
		MentionedUserIDs: mentionedUserIDs,
	})
	req, err := http.NewRequest("POST", "/api/thread/post", bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	threadRouter.ServeHTTP(rr, req)
	return rr
}

func getThreadsByUserTestFunc(workspaceID int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/thread/by_user/"+strconv.Itoa(workspaceID), nil)
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	threadRouter.ServeHTTP(rr, req)
	return rr
}

func TestPostThread(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1 正常な場合 channel 200
	// 2 正常な場合 dm 200
	// 3 channel メンションがある場合 200
	// 4 dm メンションがある場合 200

	t.Run("1 channel", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		isPrivate := false

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		lr := new(LoginResponse)
		byteArray, _ := io.ReadAll(rr.Body)
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

		rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)

		rr = postThreadTestFunc(randomstring.CookieFriendlyString(100), lr.Token, m.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m1)

		rr = postThreadTestFunc(randomstring.CookieFriendlyString(100), lr.Token, m.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m2)

		assert.Equal(t, m.ChannelId, m1.ChannelId)
		assert.Equal(t, m.ChannelId, m2.ChannelId)

		assert.Equal(t, m1.ThreadId, m2.ThreadId)

		assert.Equal(t, lr.UserId, m.UserId)
		assert.Equal(t, lr.UserId, m1.UserId)
		assert.Equal(t, lr.UserId, m2.UserId)
	})

	t.Run("2 dm", func(t *testing.T) {
		userName1 := randomstring.EnglishFrequencyString(30)
		userName2 := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName1, "pass").Code)
		rr := loginTestFunc(userName1, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		lr1 := new(LoginResponse)
		byteArray, _ := io.ReadAll(rr.Body)
		json.Unmarshal(([]byte)(byteArray), lr1)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName2, "pass").Code)
		rr = loginTestFunc(userName2, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		lr2 := new(LoginResponse)
		byteArray, _ = io.ReadAll(rr.Body)
		json.Unmarshal(([]byte)(byteArray), lr2)

		rr = createWorkSpaceTestFunc(workspaceName, lr1.Token, lr1.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, 4, lr2.UserId, lr1.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr1.Token, lr2.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm)

		rr = postThreadTestFunc(randomstring.EnglishFrequencyString(100), lr2.Token, dm.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm1)

		rr = postThreadTestFunc(randomstring.EnglishFrequencyString(100), lr1.Token, dm.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm2)

		assert.Equal(t, dm1.ThreadId, dm2.ThreadId)

		assert.Equal(t, dm.DMLineId, dm1.DMLineId)
		assert.Equal(t, dm.DMLineId, dm2.DMLineId)
	})
	t.Run("3 channel メンションがある場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		channelName := randomstring.EnglishFrequencyString(30)
		isPrivate := false

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		lr := new(LoginResponse)
		byteArray, _ := io.ReadAll(rr.Body)
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

		rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr.Token, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m)

		rr = postThreadTestFunc(randomstring.CookieFriendlyString(100), lr.Token, m.ID, []uint32{lr.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m1)

		rr = postThreadTestFunc(randomstring.CookieFriendlyString(100), lr.Token, m.ID, []uint32{lr.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		m2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), m2)

		assert.Equal(t, m.ChannelId, m1.ChannelId)
		assert.Equal(t, m.ChannelId, m2.ChannelId)

		assert.Equal(t, m1.ThreadId, m2.ThreadId)

		assert.Equal(t, lr.UserId, m.UserId)
		assert.Equal(t, lr.UserId, m1.UserId)
		assert.Equal(t, lr.UserId, m2.UserId)
	})

	t.Run("4 dm メンションがある場合", func(t *testing.T) {
		userName1 := randomstring.EnglishFrequencyString(30)
		userName2 := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName1, "pass").Code)
		rr := loginTestFunc(userName1, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		lr1 := new(LoginResponse)
		byteArray, _ := io.ReadAll(rr.Body)
		json.Unmarshal(([]byte)(byteArray), lr1)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName2, "pass").Code)
		rr = loginTestFunc(userName2, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		lr2 := new(LoginResponse)
		byteArray, _ = io.ReadAll(rr.Body)
		json.Unmarshal(([]byte)(byteArray), lr2)

		rr = createWorkSpaceTestFunc(workspaceName, lr1.Token, lr1.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, 4, lr2.UserId, lr1.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr1.Token, lr2.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm)

		rr = postThreadTestFunc(randomstring.EnglishFrequencyString(100), lr2.Token, dm.ID, []uint32{lr1.UserId, lr2.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm1)

		rr = postThreadTestFunc(randomstring.EnglishFrequencyString(100), lr1.Token, dm.ID, []uint32{lr1.UserId, lr2.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm2)

		assert.Equal(t, dm1.ThreadId, dm2.ThreadId)

		assert.Equal(t, dm.DMLineId, dm1.DMLineId)
		assert.Equal(t, dm.DMLineId, dm2.DMLineId)
	})
}

func TestGetThreadsByUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	userName1 := randomstring.EnglishFrequencyString(30)
	userName2 := randomstring.EnglishFrequencyString(30)
	workspaceName := randomstring.EnglishFrequencyString(30)
	isPrivate := false

	assert.Equal(t, http.StatusOK, signUpTestFunc(userName1, "pass").Code)
	rr := loginTestFunc(userName1, "pass")
	assert.Equal(t, http.StatusOK, rr.Code)
	lr1 := new(LoginResponse)
	byteArray, _ := io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), lr1)

	assert.Equal(t, http.StatusOK, signUpTestFunc(userName2, "pass").Code)
	rr = loginTestFunc(userName2, "pass")
	assert.Equal(t, http.StatusOK, rr.Code)
	lr2 := new(LoginResponse)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), lr2)

	rr = createWorkSpaceTestFunc(workspaceName, lr1.Token, lr1.UserId)
	assert.Equal(t, http.StatusOK, rr.Code)
	w := new(models.Workspace)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), w)

	rr = addUserWorkspaceTestFunc(w.ID, 4, lr2.UserId, lr1.Token)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = createChannelTestFunc(randomstring.EnglishFrequencyString(20), "", &isPrivate, lr1.Token, w.ID)
	assert.Equal(t, http.StatusOK, rr.Code)
	ch1 := new(models.Channel)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), ch1)

	rr = addUserInChannelTestFunc(ch1.ID, lr2.UserId, lr1.Token)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = createChannelTestFunc(randomstring.EnglishFrequencyString(20), "", &isPrivate, lr1.Token, w.ID)
	assert.Equal(t, http.StatusOK, rr.Code)
	ch2 := new(models.Channel)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), ch2)

	rr = addUserInChannelTestFunc(ch2.ID, lr2.UserId, lr1.Token)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(40), ch1.ID, lr1.Token, []uint32{})
	assert.Equal(t, http.StatusOK, rr.Code)
	pm1 := new(models.Message)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), pm1)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(40), ch2.ID, lr1.Token, []uint32{})
	assert.Equal(t, http.StatusOK, rr.Code)
	pm2 := new(models.Message)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), pm2)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(40), ch2.ID, lr1.Token, []uint32{})
	assert.Equal(t, http.StatusOK, rr.Code)
	pm3 := new(models.Message)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), pm3)

	rr = postThreadTestFunc(randomstring.EnglishFrequencyString(40), lr1.Token, pm1.ID, []uint32{})
	assert.Equal(t, http.StatusOK, rr.Code)
	m1 := new(models.Message)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), m1)

	rr = postThreadTestFunc(randomstring.EnglishFrequencyString(40), lr2.Token, pm2.ID, []uint32{})
	assert.Equal(t, http.StatusOK, rr.Code)
	m2 := new(models.Message)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), m2)

	var res1, res2 []ThreadInfo
	rr = getThreadsByUserTestFunc(w.ID, lr1.Token)
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), &res1)

	rr = getThreadsByUserTestFunc(w.ID, lr2.Token)
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	json.Unmarshal(([]byte)(byteArray), &res2)

	assert.Equal(t, 2, len(res1))
	assert.Equal(t, 1, len(res2))
}
