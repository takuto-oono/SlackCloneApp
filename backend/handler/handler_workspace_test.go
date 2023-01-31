package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var workspaceRouter = SetupRouter()

func TestCreateWorkspace(t *testing.T) {
	// 1. 正常な状態(ログイン中のユーザーがworkspaceを作成する) 200
	// 2. ログインしていない状態でのアクセス(headerにauthenticate情報がない) 400
	// 3. authenticate情報からUserIdが特定できない場合 400
	// 4. workspace nameがrequestのbodyから取得できない場合 400

	// 1
	rr := httptest.NewRecorder()
	name := "testCreateWorkspaceUserName"
	password := "pass"
	input := UserInput{
		Name:     name,
		PassWord: password,
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	router.ServeHTTP(rr, req)
	jwtToken := rr.Body.String()

	rr = httptest.NewRecorder()
	inputWorkspace := WorkspaceInput{
		Name: "testCreateWorkspaceOK",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)

	// 2
	rr = httptest.NewRecorder()
	inputWorkspace = WorkspaceInput{
		Name: "testCreateWorkspaceNotAuthenticate",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)

	// 3
	rr = httptest.NewRecorder()
	inputWorkspace = WorkspaceInput{
		Name: "testCreateWorkspaceCantAuthenticate",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", "testJWTToken")
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)

	// 4
	// body自体がない場合
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/workspace/create", nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)

	// bodyはあるが、Nameに何も指定がない場合
	rr = httptest.NewRecorder()
	inputWorkspace = WorkspaceInput{
		Name: "",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)

}
