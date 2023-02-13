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

func addUserInChannelTestFunc(channelId, workspaceId int, userId uint32, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	cau := models.NewChannelsAndUses(channelId, userId, false)
	jsonInput, _ := json.Marshal(cau)
	req, _ := http.NewRequest("POST", "/api/channel/add_user/"+strconv.Itoa(workspaceId), bytes.NewBuffer(jsonInput))
	req.Header.Set("Authorization", jwtToken)
	channelRouter.ServeHTTP(rr, req)
	return rr
}

func TestCreateChannel(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. urlからパラメータが取得できない場合 400
	// 3. bodyにchannel nameがない場合 400
	// 4. requestしたuserが対象のworkspaceに所属していない場合 400
	// 5. すでに同じ名前のchannelが対象のworkspaceに存在している場合 400

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

func TestAddUserInChannel(t *testing.T) {
	/* if testing.Short() { */
	/* 	t.Skip("skipping test in short mode.") */
	/* } */

	// 1. 正常な場合 200
	// 2. bodyに不足がある場合(channel_id, user_id) 400
	// 3. リクエストしたuserがworkspaceに参加していない場合 400
	// 4. 追加されるuserがworkspaceに参加していない場合 400
	// 5. 対象のchannelがworkspaceに存在していない場合 400
	// 6. 追加されるuserが対象のchannelに既に存在してる場合 400
	// 7. リクエストしたuserにチャンネルの管理権限がない場合 400

	t.Run("1", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName1"
		addUserName := "testAddUserInChannelAddUserName1"
		workspaceName := "testAddUserInChannelWorkspaceName1"
		channelName := "testAddUserInChannelChannelName1"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, alr.UserId, rlr.Token).Code)

		rr = createChannelTestFunc(channelName, "des", true, rlr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		rr = addUserInChannelTestFunc(c.ID, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		cau := new(models.ChannelsAndUsers)
		json.Unmarshal(([]byte)(byteArray), cau)
		assert.Equal(t, c.ID, cau.ChannelId)
		assert.Equal(t, alr.UserId, cau.UserId)
		assert.ElementsMatch(t, false, cau.IsAdmin)
	})

	t.Run("2", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName2"
		addUserName := "testAddUserInChannelAddUserName2"
		workspaceName := "testAddUserInChannelWorkspaceName2"
		channelName := "testAddUserInChannelChannelName2"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, alr.UserId, rlr.Token).Code)

		rr = createChannelTestFunc(channelName, "des", true, rlr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		rr = addUserInChannelTestFunc(0, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found channel_id or user_id\"}", rr.Body.String())

		rr = addUserInChannelTestFunc(c.ID, w.ID, 0, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found channel_id or user_id\"}", rr.Body.String())
	})

	t.Run("3", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName3"
		addUserName := "testAddUserInChannelAddUserName3"
		createWorkspaceUserName := "testAddUserInChannelCreateWorkspaceUserName3"
		workspaceName := "testAddUserInChannelWorkspaceName3"
		channelName := "testAddUserInChannelChannelName3"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(createWorkspaceUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = loginTestFunc(createWorkspaceUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		clr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), clr)

		rr = createWorkSpaceTestFunc(workspaceName, clr.Token, clr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, alr.UserId, clr.Token).Code)

		rr = createChannelTestFunc(channelName, "des", true, clr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		rr = addUserInChannelTestFunc(c.ID, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not exist request user in workspace\"}", rr.Body.String())

	})

	t.Run("4", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName4"
		addUserName := "testAddUserInChannelAddUserName4"
		workspaceName := "testAddUserInChannelWorkspaceName4"
		channelName := "testAddUserInChannelChannelName4"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = createChannelTestFunc(channelName, "des", true, rlr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		rr = addUserInChannelTestFunc(c.ID, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not exist added user in workspace\"}", rr.Body.String())
	})

	t.Run("5", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName5"
		addUserName := "testAddUserInChannelAddUserName5"
		workspaceName := "testAddUserInChannelWorkspaceName5"
		channelName := "testAddUserInChannelChannelName5"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, alr.UserId, rlr.Token).Code)

		rr = createChannelTestFunc(channelName, "des", true, rlr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		rr = addUserInChannelTestFunc(-1, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not exist channel in workspace\"}", rr.Body.String())
	})

	t.Run("6", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName6"
		addUserName := "testAddUserInChannelAddUserName6"
		workspaceName := "testAddUserInChannelWorkspaceName6"
		channelName := "testAddUserInChannelChannelName6"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, alr.UserId, rlr.Token).Code)

		rr = createChannelTestFunc(channelName, "des", true, rlr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		assert.Equal(t, http.StatusOK, addUserInChannelTestFunc(c.ID, w.ID, alr.UserId, rlr.Token).Code)

		rr = addUserInChannelTestFunc(c.ID, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		assert.Equal(t, "{\"message\":\"already exist user in channel\"}", rr.Body.String())
	})

	t.Run("7", func(t *testing.T) {
		requestUserName := "testAddUserInChannelRequestUserName7"
		addUserName := "testAddUserInChannelAddUserName7"
		createChannelUserName := "testAddUserInChannelCreateUserName7"
		workspaceName := "testAddUserInChannelWorkspaceName7"
		channelName := "testAddUserInChannelChannelName7"

		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(createChannelUserName, "pass").Code)

		rr := loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = loginTestFunc(createChannelUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		clr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), clr)

		rr = createWorkSpaceTestFunc(workspaceName, rlr.Token, rlr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, alr.UserId, rlr.Token).Code)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, clr.UserId, rlr.Token).Code)

		rr = createChannelTestFunc(channelName, "des", true, clr.Token, w.ID)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		c := new(models.Channel)
		json.Unmarshal(([]byte)(byteArray), c)

		rr = addUserInChannelTestFunc(c.ID, w.ID, rlr.UserId, clr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		rr = addUserInChannelTestFunc(c.ID, w.ID, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"no permission adding user in channel\"}", rr.Body.String())
	})
}
