package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"backend/models"
	"backend/utils"
	"backend/controllerUtils"
)

var testRouter = SetupRouter1()
var testRouter2 = SetupRouter2()

type LoginResponse struct {
	Token    string `json:"token"`
	UserId   uint32 `json:"user_id"`
	Username string `json:"username"`
}

func Req(method, endpoint, jwtToken string, bodyInfo interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(bodyInfo)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	if jwtToken != "" {
		req.Header.Set("Authorization", jwtToken)
	}
	testRouter.ServeHTTP(rr, req)
	return rr
}

func signUpTestFuncV2(userName, pass string) (*httptest.ResponseRecorder, models.User) {
	rr := Req(http.MethodPost, "/api/user/signUp", "", controllerUtils.SignUpAndLoginInput{
		Name:     userName,
		Password: pass,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var u models.User
	json.Unmarshal(([]byte)(byteArray), &u)
	return rr, u
}

func loginTestFuncV2(name, password string) (*httptest.ResponseRecorder, LoginResponse) {
	rr := Req(http.MethodPost, "/api/user/login", "", controllerUtils.SignUpAndLoginInput{
		Name:     name,
		Password: password,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var lr LoginResponse
	json.Unmarshal(([]byte)(byteArray), &lr)
	return rr, lr
}

func createWorkspaceTestFuncV2(workspaceName, jwtToken string, userID uint32) (*httptest.ResponseRecorder, models.Workspace) {
	rr := Req(http.MethodPost, "/api/workspace/create", jwtToken, controllerUtils.CreateWorkspaceInput{
		Name:          workspaceName,
		RequestUserId: userID,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var w models.Workspace
	json.Unmarshal(([]byte)(byteArray), &w)
	return rr, w
}

func addUserInWorkspaceV2(workspaceID int, userID uint32, roleID int, jwtToken string) (*httptest.ResponseRecorder, models.WorkspaceAndUsers) {
	rr := Req(http.MethodPost, "/api/workspace/add_user", jwtToken, controllerUtils.AddUserInWorkspaceInput{
		WorkspaceId: workspaceID,
		UserId:      userID,
		RoleId:      roleID,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var wau models.WorkspaceAndUsers
	json.Unmarshal(([]byte)(byteArray), &wau)
	return rr, wau
}

func createChannelTestFuncV2(channelName, description string, isPrivate *bool, jwtToken string, workspaceID int) (*httptest.ResponseRecorder, models.Channel) {
	rr := Req(http.MethodPost, "/api/channel/create", jwtToken, controllerUtils.CreateChannelInput{
		Name:        channelName,
		Description: description,
		IsPrivate:   isPrivate,
		WorkspaceId: workspaceID,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var ch models.Channel
	json.Unmarshal(([]byte)(byteArray), &ch)
	return rr, ch
}

func addUserInChannelTestFuncV2(channelID int, userID uint32, jwtToken string) (*httptest.ResponseRecorder, models.ChannelsAndUsers) {
	rr := Req(http.MethodPost, "/api/channel/add_user", jwtToken, controllerUtils.AddUserInChannelInput{
		ChannelId: channelID,
		UserId:    userID,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var cau models.ChannelsAndUsers
	json.Unmarshal(([]byte)(byteArray), &cau)
	return rr, cau
}

func sendMessageTestFuncV2(text string, channelID int, jwtToken string, mentionedUserIDs []uint32) (*httptest.ResponseRecorder, models.Message) {
	rr := Req(http.MethodPost, "/api/message/send", jwtToken, controllerUtils.SendMessageInput{
		Text:             text,
		ChannelId:        channelID,
		MentionedUserIDs: mentionedUserIDs,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var m models.Message
	json.Unmarshal(([]byte)(byteArray), &m)
	return rr, m
}

func sendDMTestFuncV2(text string, jwtToken string, receiveUserID uint32, workspaceID int, mentionedUserIDs []uint32) (*httptest.ResponseRecorder, models.Message) {
	rr := Req(http.MethodPost, "/api/dm/send", jwtToken, controllerUtils.SendDMInput{
		Text:             text,
		ReceiveUserId:    receiveUserID,
		WorkspaceId:      workspaceID,
		MentionedUserIDs: mentionedUserIDs,
	})
	byteArray, _ := io.ReadAll(rr.Body)
	var m models.Message
	json.Unmarshal(([]byte)(byteArray), &m)
	return rr, m
}

func readMessageByUserTestFunc(messageID uint, jwtToken string) (*httptest.ResponseRecorder, models.MessageAndUser) {
	rr := Req(http.MethodPost, "/api/message/read_by_user/"+utils.UintToString(messageID), jwtToken, nil)
	byteArray, _ := io.ReadAll(rr.Body)
	var mau models.MessageAndUser
	json.Unmarshal(([]byte)(byteArray), &mau)
	return rr, mau
}

func getAllUsersInChannelTestFuncV2(channelID int, jwtToken string) (*httptest.ResponseRecorder, []UserResponse) {
	rr := Req(http.MethodGet, "/api/channel/all_user/"+strconv.Itoa(channelID), jwtToken, nil)
	byteArray, _ := io.ReadAll(rr.Body)
	var users []UserResponse
	json.Unmarshal(([]byte)(byteArray), &users)
	return rr, users
}

func getAllMessagesFromChannelTestFuncV2(channelID int, jwtToken string) (*httptest.ResponseRecorder, []models.Message) {
	rr := Req(http.MethodGet, "/api/message/get_from_channel/"+strconv.Itoa(channelID), jwtToken, nil)
	byteArray, _ := io.ReadAll(rr.Body)
	var messages []models.Message
	json.Unmarshal(([]byte)(byteArray), &messages)
	return rr, messages
}
