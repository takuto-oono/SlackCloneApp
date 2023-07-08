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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xyproto/randomstring"

	"backend/controllerUtils"
	"backend/models"
	"backend/utils"
)

var dmRouter = SetupRouter1()

func sendDMTestFunc(text, jwtToken string, receiveUserId uint32, workspaceId int, mentionedUserIDs []uint32) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.SendDMInput{
		Text:             text,
		WorkspaceId:      workspaceId,
		ReceiveUserId:    receiveUserId,
		MentionedUserIDs: mentionedUserIDs,
	})
	req, err := http.NewRequest("POST", "/api/dm/send", bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	dmRouter.ServeHTTP(rr, req)
	return rr
}

func getDMsInLineTestFunc(dlId uint, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/dm/"+strconv.Itoa(int(dlId)), nil)
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	dmRouter.ServeHTTP(rr, req)
	return rr
}

func getDMLinesTestFunc(workspaceId int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/dm/dm_lines/"+strconv.Itoa(workspaceId), nil)
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	dmRouter.ServeHTTP(rr, req)
	return rr
}

func editDMTestFunc(dmId uint, jwtToken, text string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.EditDMInput{
		Text: text,
	})
	req, err := http.NewRequest("PATCH", "/api/dm/"+strconv.Itoa(int(dmId)), bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	dmRouter.ServeHTTP(rr, req)
	return rr
}

func deleteDMTestFunc(dmId uint, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/dm/"+strconv.Itoa(int(dmId)), nil)
	req.Header.Set("Authorization", jwtToken)
	dmRouter.ServeHTTP(rr, req)
	return rr
}

func TestSendDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. 自分自身に送信する場合 200
	// 3. bodyに不足がある場合 400
	// 4. requestしたuserがworkspaceに存在しない場合 404 (workspaceが存在しない場合も含まれる)
	// 5. receiveするuserがworkspaceに存在しない場合 404
	// 6 メンションがある場合 200
	// 7 スケジュールされたメッセージがある場合 200

	t.Run("1 正常な場合", func(t *testing.T) {
		sendUserName := randomstring.EnglishFrequencyString(30)
		receiveUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		text1 := randomstring.EnglishFrequencyString(100)
		text2 := randomstring.EnglishFrequencyString(100)
		text3 := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(sendUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(receiveUserName, "pass").Code)

		rr := loginTestFunc(sendUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		slr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), slr)

		rr = loginTestFunc(receiveUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, slr.Token, slr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, slr.Token).Code)

		rr = sendDMTestFunc(text1, slr.Token, rlr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm1)
		assert.Equal(t, text1, dm1.Text)
		assert.NotEqual(t, uint(0), dm1.DMLineId)

		rr = sendDMTestFunc(text2, slr.Token, rlr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm2)
		assert.Equal(t, text2, dm2.Text)
		assert.NotEqual(t, uint(0), dm2.DMLineId)

		rr = sendDMTestFunc(text3, rlr.Token, slr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm3 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm3)
		assert.Equal(t, text3, dm3.Text)
		assert.NotEqual(t, uint(0), dm3.DMLineId)

		assert.Equal(t, dm1.DMLineId, dm2.DMLineId)
		assert.Equal(t, dm1.DMLineId, dm3.DMLineId)
	})

	t.Run("2 自分自身に送信する場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		text1 := randomstring.EnglishFrequencyString(100)
		text2 := randomstring.EnglishFrequencyString(100)

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

		rr = sendDMTestFunc(text1, lr.Token, lr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm1)
		assert.Equal(t, text1, dm1.Text)
		assert.NotEqual(t, uint(0), dm1.DMLineId)

		rr = sendDMTestFunc(text2, lr.Token, lr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm2)
		assert.Equal(t, text2, dm2.Text)
		assert.NotEqual(t, uint(0), dm2.DMLineId)

		assert.Equal(t, dm1.DMLineId, dm2.DMLineId)
	})

	t.Run("3 bodyに不足がある場合", func(t *testing.T) {
		sendUserName := randomstring.EnglishFrequencyString(30)
		receiveUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(sendUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(receiveUserName, "pass").Code)

		rr := loginTestFunc(sendUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		slr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), slr)

		rr = loginTestFunc(receiveUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, slr.Token, slr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, slr.Token).Code)

		rr = sendDMTestFunc("", slr.Token, rlr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"text not found\"}", rr.Body.String())

		rr = sendDMTestFunc(text, slr.Token, 0, w.ID, []uint32{})
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"received_user_id not found\"}", rr.Body.String())

		rr = sendDMTestFunc(text, slr.Token, rlr.UserId, 0, []uint32{})
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"workspace_id not found\"}", rr.Body.String())
	})

	t.Run("4 requestしたuserがworkspaceに存在しない場合", func(t *testing.T) {
		sendUserName := randomstring.EnglishFrequencyString(30)
		receiveUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(sendUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(receiveUserName, "pass").Code)

		rr := loginTestFunc(sendUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		slr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), slr)

		rr = loginTestFunc(receiveUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = sendDMTestFunc(text, slr.Token, rlr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"send user not found in workspace\"}", rr.Body.String())
	})

	t.Run("5 receiveするuserがworkspaceに存在しない場合", func(t *testing.T) {
		sendUserName := randomstring.EnglishFrequencyString(30)
		receiveUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		text := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(sendUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(receiveUserName, "pass").Code)

		rr := loginTestFunc(sendUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		slr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), slr)

		rr = loginTestFunc(receiveUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, slr.Token, slr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = sendDMTestFunc(text, slr.Token, rlr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"receive user not found in workspace\"}", rr.Body.String())
	})

	t.Run("6 メンションがある場合", func(t *testing.T) {
		sendUserName := randomstring.EnglishFrequencyString(30)
		receiveUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		text1 := randomstring.EnglishFrequencyString(100)
		text2 := randomstring.EnglishFrequencyString(100)
		text3 := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(sendUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(receiveUserName, "pass").Code)

		rr := loginTestFunc(sendUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		slr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), slr)

		rr = loginTestFunc(receiveUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, slr.Token, slr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, slr.Token).Code)

		rr = sendDMTestFunc(text1, slr.Token, rlr.UserId, w.ID, []uint32{rlr.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm1)
		assert.Equal(t, text1, dm1.Text)
		assert.NotEqual(t, uint(0), dm1.DMLineId)

		rr = sendDMTestFunc(text2, slr.Token, rlr.UserId, w.ID, []uint32{slr.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm2)
		assert.Equal(t, text2, dm2.Text)
		assert.NotEqual(t, uint(0), dm2.DMLineId)

		rr = sendDMTestFunc(text3, rlr.Token, slr.UserId, w.ID, []uint32{slr.UserId, rlr.UserId})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm3 := new(models.Message)
		json.Unmarshal(([]byte)(byteArray), dm3)
		assert.Equal(t, text3, dm3.Text)
		assert.NotEqual(t, uint(0), dm3.DMLineId)

		assert.Equal(t, dm1.DMLineId, dm2.DMLineId)
		assert.Equal(t, dm1.DMLineId, dm3.DMLineId)
	})

	t.Run("7 スケジュールされたメッセージがある場合", func(t *testing.T) {
		rr, user := signUpTestFuncV2(randomstring.EnglishFrequencyString(30), "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, lr := loginTestFuncV2(user.Name, user.PassWord)
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, w := createWorkspaceTestFuncV2(randomstring.EnglishFrequencyString(30), lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		rr, scheduledMessage := sendDMTestFuncV2(randomstring.EnglishFrequencyString(30), lr.Token, lr.UserId, w.ID, []uint32{}, time.Now().Add(time.Second*15))
		assert.Equal(t, http.StatusOK, rr.Code)

		rr, res := getDMsInLineTestFuncV2(scheduledMessage.DMLineId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 0, len(res))

		for i := 0; i < 5; i++ {
			rr, _ := sendDMTestFuncV2(randomstring.EnglishFrequencyString(30), lr.Token, lr.UserId, w.ID, []uint32{}, utils.CreateDefaultTime())
			assert.Equal(t, http.StatusOK, rr.Code)
			time.Sleep(1 * time.Second)
		}

		rr, res = getDMsInLineTestFuncV2(scheduledMessage.DMLineId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 5, len(res))
		time.Sleep(time.Second * 15)
		rr, res = getDMsInLineTestFuncV2(scheduledMessage.DMLineId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, 6, len(res))
		assert.Equal(t, scheduledMessage.ID, res[0].ID)
	})
}

func TestGetAllDMsByDLId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1 正常な場合 dmが存在している場合 200
	// 2 存在しないdm_lineだった場合 404
	// 3 requestしたuserが参加していないdm_lineの場合 403

	t.Run("1 正常な場合 dmが存在している場合", func(t *testing.T) {
		messageCount := 3
		username := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		dms := make([]models.Message, messageCount)
		text := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(username, "pass").Code)

		rr := loginTestFunc(username, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		for i := 0; i < messageCount; i++ {
			rr = sendDMTestFunc(text, lr.Token, lr.UserId, w.ID, []uint32{})
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ = io.ReadAll(rr.Body)
			var dm models.Message
			json.Unmarshal(([]byte)(byteArray), &dm)
			dms[i] = dm
		}

		rr = getDMsInLineTestFunc(dms[0].DMLineId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		res := make([]models.Message, messageCount)
		byteArray, _ = io.ReadAll(rr.Body)
		json.Unmarshal(([]byte)(byteArray), &res)
		assert.Equal(t, messageCount, len(res))

		for i := 0; i < messageCount-1; i++ {
			assert.True(t, res[i].CreatedAt.After(res[i+1].CreatedAt))
		}

		for _, dm := range res {
			assert.Equal(t, dms[0].DMLineId, dm.DMLineId)
			assert.Equal(t, text, dm.Text)
		}

	})

	t.Run("2 存在しないdm_lineだった場合", func(t *testing.T) {
		username := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(username, "pass").Code)

		rr := loginTestFunc(username, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = getDMsInLineTestFunc(uint(rand.Uint64()), lr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"record not found\"}", rr.Body.String())
	})

	t.Run("3 requestしたuserが参加していないdm_lineだった場合", func(t *testing.T) {
		messageCount := 3
		username := randomstring.EnglishFrequencyString(30)
		requestUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		dms := make([]models.Message, messageCount)
		text := randomstring.EnglishFrequencyString(100)

		assert.Equal(t, http.StatusOK, signUpTestFunc(username, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)

		rr := loginTestFunc(username, "pass")
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

		for i := 0; i < messageCount; i++ {
			rr = sendDMTestFunc(text, lr.Token, lr.UserId, w.ID, []uint32{})
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ = io.ReadAll(rr.Body)
			var dm models.Message
			json.Unmarshal(([]byte)(byteArray), &dm)
			dms[i] = dm
		}

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, lr.Token).Code)

		rr = getDMsInLineTestFunc(dms[0].DMLineId, rlr.Token)
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"you don't access this page\"}", rr.Body.String())
	})
}

func TestGetDMLines(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 データが存在する場合 200
	// 2. 正常な場合 データが存在しない場合 200
	// 3. userがworkspaceに存在していない場合 404

	t.Run("1 正常な場合 データが存在する場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		toUserNames := make([]string, 10)
		workspaceName := randomstring.EnglishFrequencyString(30)
		for i := 0; i < 10; i++ {
			toUserNames[i] = randomstring.EnglishFrequencyString(30)
		}

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		tlrs := make([]LoginResponse, 10)
		for i, userName := range toUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
			rr := loginTestFunc(userName, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := io.ReadAll(rr.Body)
			tlr := new(LoginResponse)
			json.Unmarshal(([]byte)(byteArray), tlr)
			tlrs[i] = *tlr
		}

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		for _, tlr := range tlrs {
			assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, tlr.UserId, lr.Token).Code)
		}

		for _, tlr := range tlrs {
			assert.Equal(t, http.StatusOK, sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, tlr.UserId, w.ID, []uint32{}).Code)
		}

		rr = getDMLinesTestFunc(w.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		response := new([]DMLineInfo)
		json.Unmarshal(([]byte)(byteArray), response)

		assert.Equal(t, 10, len(*response))

		for _, res := range *response {
			assert.Contains(t, toUserNames, res.ToName)
		}
	})

	t.Run("2 正常な場合 データが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

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

		rr = getDMLinesTestFunc(w.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		response := new([]DMLineInfo)
		json.Unmarshal(([]byte)(byteArray), response)

		assert.Equal(t, 0, len(*response))
	})

	t.Run("3 userがworkspaceに存在していない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		requestUserName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
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

		assert.Equal(t, http.StatusNotFound, getDMLinesTestFunc(w.ID, rlr.Token).Code)
	})

}

func TestEditDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. bodyにtextがなかった場合 400
	// 3. DMが存在しない場合 404
	// 4. DMが他のuserが送信したものだった場合 403

	t.Run("1 正常な場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		oldText := randomstring.EnglishFrequencyString(100)
		newText := randomstring.EnglishFrequencyString(100)

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

		rr = sendDMTestFunc(oldText, lr.Token, lr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.Message
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = editDMTestFunc(dm.ID, lr.Token, newText)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var res models.Message
		json.Unmarshal(([]byte)(byteArray), &res)

		assert.Equal(t, dm.ID, res.ID)
		assert.Equal(t, newText, res.Text)
		assert.Equal(t, dm.DMLineId, res.DMLineId)
		assert.Equal(t, dm.CreatedAt, res.CreatedAt)
		assert.True(t, res.UpdatedAt.After(dm.UpdatedAt))
	})

	t.Run("2 bodyにtextがなかった場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = editDMTestFunc(uint(rand.Uint64()), lr.Token, "")
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"text not found\"}", rr.Body.String())
	})

	t.Run("3 DMが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := io.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = editDMTestFunc(uint(rand.Uint64()), lr.Token, randomstring.EnglishFrequencyString(100))
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"dm not found\"}", rr.Body.String())
	})

	t.Run("DMが他のuserが送信したものだった場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		requestUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

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

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, lr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.Message
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = editDMTestFunc(dm.ID, rlr.Token, randomstring.EnglishFrequencyString(100))
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"no permission\"}", rr.Body.String())
	})
}

func TestDeleteDM(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. 対象のDMが存在しない場合 404
	// 3. 対象のDMは別のuserの場合 403

	t.Run("1 正常な場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

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

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, lr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.Message
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = deleteDMTestFunc(dm.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var res models.Message
		json.Unmarshal(([]byte)(byteArray), &res)
		assert.Equal(t, dm.ID, res.ID)
		assert.Equal(t, dm.Text, res.Text)
		assert.Equal(t, dm.DMLineId, res.DMLineId)
		assert.Equal(t, dm.CreatedAt, res.CreatedAt)
		assert.Equal(t, dm.UpdatedAt, res.UpdatedAt)

	})

	t.Run("2 対象のDMが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

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

		rr = deleteDMTestFunc(uint(rand.Uint64()), lr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"dm not found\"}", rr.Body.String())
	})

	t.Run("3 対象のDMは別のuserの場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		requestUserName := randomstring.EnglishFrequencyString(30)

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

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, lr.UserId, w.ID, []uint32{})
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.Message
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = deleteDMTestFunc(dm.ID, rlr.Token)
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"no permission\"}", rr.Body.String())
	})
}
