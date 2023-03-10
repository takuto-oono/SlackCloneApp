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
)

var dmRouter = SetupRouter()

func sendDMTestFunc(text, jwtToken string, receiveUserId uint32, workspaceId int) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.SendDMInput{
		Text:          text,
		WorkspaceId:   workspaceId,
		ReceiveUserId: receiveUserId,
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

		rr = sendDMTestFunc(text1, slr.Token, rlr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.DirectMessage)
		json.Unmarshal(([]byte)(byteArray), dm1)
		assert.Equal(t, text1, dm1.Text)
		assert.Equal(t, slr.UserId, dm1.SendUserId)
		assert.NotEqual(t, uint(0), dm1.DMLineId)

		rr = sendDMTestFunc(text2, slr.Token, rlr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.DirectMessage)
		json.Unmarshal(([]byte)(byteArray), dm2)
		assert.Equal(t, text2, dm2.Text)
		assert.Equal(t, slr.UserId, dm2.SendUserId)
		assert.NotEqual(t, uint(0), dm2.DMLineId)

		rr = sendDMTestFunc(text3, rlr.Token, slr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm3 := new(models.DirectMessage)
		json.Unmarshal(([]byte)(byteArray), dm3)
		assert.Equal(t, text3, dm3.Text)
		assert.Equal(t, rlr.UserId, dm3.SendUserId)
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

		rr = sendDMTestFunc(text1, lr.Token, lr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm1 := new(models.DirectMessage)
		json.Unmarshal(([]byte)(byteArray), dm1)
		assert.Equal(t, text1, dm1.Text)
		assert.Equal(t, lr.UserId, dm1.SendUserId)
		assert.NotEqual(t, uint(0), dm1.DMLineId)

		rr = sendDMTestFunc(text2, lr.Token, lr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		dm2 := new(models.DirectMessage)
		json.Unmarshal(([]byte)(byteArray), dm2)
		assert.Equal(t, text2, dm2.Text)
		assert.Equal(t, lr.UserId, dm2.SendUserId)
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

		rr = sendDMTestFunc("", slr.Token, rlr.UserId, w.ID)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"text not found\"}", rr.Body.String())

		rr = sendDMTestFunc(text, slr.Token, 0, w.ID)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"received_user_id not found\"}", rr.Body.String())

		rr = sendDMTestFunc(text, slr.Token, rlr.UserId, 0)
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

		rr = sendDMTestFunc(text, slr.Token, rlr.UserId, w.ID)
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

		rr = sendDMTestFunc(text, slr.Token, rlr.UserId, w.ID)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"receive user not found in workspace\"}", rr.Body.String())
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
		dms := make([]models.DirectMessage, messageCount)
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
			rr = sendDMTestFunc(text, lr.Token, lr.UserId, w.ID)
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ = io.ReadAll(rr.Body)
			var dm models.DirectMessage
			json.Unmarshal(([]byte)(byteArray), &dm)
			dms[i] = dm
		}

		rr = getDMsInLineTestFunc(dms[0].DMLineId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		res := make([]models.DirectMessage, messageCount)
		byteArray, _ = io.ReadAll(rr.Body)
		json.Unmarshal(([]byte)(byteArray), &res)
		assert.Equal(t, messageCount, len(res))

		for i := 0; i < messageCount-1; i++ {
			assert.True(t, res[i].CreatedAt.After(res[i+1].CreatedAt))
		}

		for _, dm := range res {
			assert.Equal(t, dms[0].DMLineId, dm.DMLineId)
			assert.Equal(t, text, dm.Text)
			assert.Equal(t, lr.UserId, dm.SendUserId)
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
		dms := make([]models.DirectMessage, messageCount)
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
			rr = sendDMTestFunc(text, lr.Token, lr.UserId, w.ID)
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ = io.ReadAll(rr.Body)
			var dm models.DirectMessage
			json.Unmarshal(([]byte)(byteArray), &dm)
			dms[i] = dm
		}

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, lr.Token).Code)

		rr = getDMsInLineTestFunc(dms[0].DMLineId, rlr.Token)
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"you don't access this page\"}", rr.Body.String())
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

		rr = sendDMTestFunc(oldText, lr.Token, lr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.DirectMessage
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = editDMTestFunc(dm.ID, lr.Token, newText)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var res models.DirectMessage
		json.Unmarshal(([]byte)(byteArray), &res)

		assert.Equal(t, dm.ID, res.ID)
		assert.Equal(t, newText, res.Text)
		assert.Equal(t, dm.SendUserId, res.SendUserId)
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

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, lr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.DirectMessage
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

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, lr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.DirectMessage
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = deleteDMTestFunc(dm.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)	
		byteArray, _ = io.ReadAll(rr.Body)
		var res models.DirectMessage
		json.Unmarshal(([]byte)(byteArray), &res)
		assert.Equal(t, dm.ID, res.ID)
		assert.Equal(t, dm.Text, res.Text)
		assert.Equal(t, dm.SendUserId, res.SendUserId)
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

		rr = sendDMTestFunc(randomstring.EnglishFrequencyString(100), lr.Token, lr.UserId, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = io.ReadAll(rr.Body)
		var dm models.DirectMessage
		json.Unmarshal(([]byte)(byteArray), &dm)

		rr = deleteDMTestFunc(dm.ID, rlr.Token)
		assert.Equal(t, http.StatusForbidden, rr.Code)	
		assert.Equal(t, "{\"message\":\"no permission\"}", rr.Body.String())
	})
}
