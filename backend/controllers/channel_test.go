package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
)

var channelRouter = SetupRouter()

func createChannelTestFunc(name, description string, isPrivate bool, jwtToken string, workspaceId int) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ch := models.NewChannel(0, name, description, isPrivate, false)
	jsonInput, _ := json.Marshal(ch)
	var req *http.Request
	var err error
	if workspaceId == 0 {
		req, err = http.NewRequest("POST", "/api/channel/create/", bytes.NewBuffer(jsonInput))
	} else {
		req, err = http.NewRequest("POST", "/api/channel/create/"+strconv.Itoa(workspaceId), bytes.NewBuffer(jsonInput))

	}
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	channelRouter.ServeHTTP(rr, req)
	return rr
}

func TestCreateChannel(t *testing.T) {
	// 1. 正常な場合 200
	// 2. urlからパラメータが取得できない場合 400
	// 3. bodyにchannel nameがない場合 400
	// 4. requestしたuserが対象のworkspaceに所属していない場合 400
	// 5. すでに同じ名前のchannelが対象のworkspaceに存在している場合 400
	// 6.

	t.Run("1", func(t *testing.T) {
		userName := "testCreateChannelUserName1"
		workspaceName := "testCreateChannelWorkspaceName1"
		channelName := "testCreateChannelName1"
		description := "description1"
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, description, isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		ch := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), ch)
		assert.NotEqual(t, 0, ch.ID)
		assert.Equal(t, channelName, ch.Name)
		assert.Equal(t, isPrivate, ch.IsPrivate)
		assert.Equal(t, false, ch.IsArchive)
	})

	t.Run("2", func(t *testing.T) {
		userName := "testCreateChannelUserName2"
		workspaceName := "testCreateChannelWorkspaceName2"
		channelName := "testCreateChannelName2"
		description := "description2"
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, description, isPrivate, lr.Token, 0)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "404 page not found", rr.Body.String())
	})

	t.Run("3", func(t *testing.T) {
		userName := "testCreateChannelUserName3"
		workspaceName := "testCreateChannelWorkspaceName3"
		description := "description3"
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc("", description, isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found channel name\"}", rr.Body.String())
	})

	t.Run("4", func(t *testing.T) {
		userName := "testCreateChannelUserName4"
		createWorkspaceUserName := "testCreateChannelCreateWorkspaceUserName4"
		workspaceName := "testCreateChannelWorkspaceName4"
		channelName := "testCreateChannelName4"
		description := "description4"
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(createWorkspaceUserName, "pass").Code)

		rr := loginTestFunc(createWorkspaceUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		lr = new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createChannelTestFunc(channelName, description, isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found user in workspace\"}", rr.Body.String())
	})

	t.Run("5", func(t *testing.T) {
		userName := "testCreateChannelUserName5"
		workspaceName := "testCreateChannelWorkspaceName5"
		channelName := "testCreateChannelName5"
		description := "description5"
		isPrivate := true

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, description, isPrivate, lr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		rr = createChannelTestFunc(channelName, "", false, lr.Token, w.ID)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"already exist same name channel in workspace\"}", rr.Body.String())
	})
}
