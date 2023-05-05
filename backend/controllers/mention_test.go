package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"

	"backend/models"
)

func getMessagesMentionedByUserTestFunc(workspaceID int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/mention/by_user/"+strconv.Itoa(workspaceID), nil)
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	threadRouter.ServeHTTP(rr, req)
	return rr
}

func TestGetMessagesMentionedByUser(t *testing.T) {
	userName1 := randomstring.EnglishFrequencyString(30)
	userName2 := randomstring.EnglishFrequencyString(30)
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

	rr = createWorkSpaceTestFunc(randomstring.EnglishFrequencyString(30), lr1.Token, lr1.UserId)
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	w := new(models.Workspace)
	json.Unmarshal(([]byte)(byteArray), w)

	rr = createChannelTestFunc(randomstring.EnglishFrequencyString(30), "", &isPrivate, lr1.Token, w.ID)
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	ch := new(models.Channel)
	json.Unmarshal(([]byte)(byteArray), ch)

	assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, lr2.UserId, lr1.Token).Code)

	assert.Equal(t, http.StatusOK, addUserInChannelTestFunc(ch.ID, lr2.UserId, lr1.Token).Code)

	cnt1, cnt2 := 0, 0
	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr1.Token, []uint32{lr2.UserId})
	cnt2++
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var m1 models.Message
	json.Unmarshal(([]byte)(byteArray), &m1)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr2.Token, []uint32{lr1.UserId})
	cnt1++
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var m2 models.Message
	json.Unmarshal(([]byte)(byteArray), &m2)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr1.Token, []uint32{lr2.UserId})
	cnt2++
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var m3 models.Message
	json.Unmarshal(([]byte)(byteArray), &m3)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr1.Token, []uint32{})
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var m4 models.Message
	json.Unmarshal(([]byte)(byteArray), &m4)

	rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr1.Token, lr2.UserId, w.ID, []uint32{lr2.UserId})
	cnt2++
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var m5 models.Message
	json.Unmarshal(([]byte)(byteArray), &m5)

	rr = sendMessageTestFunc(randomstring.EnglishFrequencyString(100), ch.ID, lr1.Token, []uint32{lr2.UserId})
	cnt2++
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var m6 models.Message
	json.Unmarshal(([]byte)(byteArray), &m6)

	rr = getMessagesMentionedByUserTestFunc(w.ID, lr1.Token)
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var res1 []models.Message
	json.Unmarshal(([]byte)(byteArray), &res1)

	rr = getMessagesMentionedByUserTestFunc(w.ID, lr2.Token)
	assert.Equal(t, http.StatusOK, rr.Code)
	byteArray, _ = io.ReadAll(rr.Body)
	var res2 []models.Message
	json.Unmarshal(([]byte)(byteArray), &res2)

	assert.Equal(t, cnt1, len(res1))
	assert.Equal(t, cnt2, len(res2))

	for i := 0; i < cnt1-1; i++ {
		assert.True(t, res1[i+1].CreatedAt.Before(res1[i].CreatedAt))
	}
	for i := 0; i < cnt2-1; i++ {
		assert.True(t, res2[i+1].CreatedAt.Before(res2[i].CreatedAt))
	}

	isExist := func(v models.Message, sl []models.Message) bool {
		for _, x := range sl {
			if x.ID == v.ID {
				return true
			}
		}
		return false
	}

	assert.True(t, isExist(m1, res2))
	assert.True(t, isExist(m2, res1))
	assert.True(t, isExist(m3, res2))
	assert.False(t, isExist(m4, res1))
	assert.False(t, isExist(m4, res2))
	assert.True(t, isExist(m5, res2))
	assert.True(t, isExist(m6, res2))
}
