package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	assert.Equal(t, http.StatusOK, rr.Code)

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

func TestAddUserWorkspace(t *testing.T) {
	// 1. jwtTokenがheaderに含まれていない場合 400
	// 2. jwtTokenからuserIdが取得できない場合 400
	// 3. requestのbodyの情報に不足がある場合 400
	// 4. 存在しないworkspaceNameだった場合 400
	// 5. 存在しないuserNameだった場合 400
	// 6. role >= 4の場合 400
	// 7. 正常な場合 200

	// 1
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/workspace/add_user", nil)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 2
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", nil)
	req.Header.Add("Authorization", "abc")
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 3
	rr = httptest.NewRecorder()
	input := UserInput{
		Name:     "addUserTest3Name",
		PassWord: "pass",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken := rr.Body.String()

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	rr = httptest.NewRecorder()
	auwi := AddUserWorkspaceInput{
		WorkspaceName: "",
	}
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: "name",
		AddUserName:   "",
		RoleId:        2,
	}
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 4
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "addUserWorkspaceTest4Name",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: "not exist workspace name",
		AddUserName:   "testUser",
		RoleId:        2,
	}
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 5
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "addUserWorkspaceTest5Name",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	wi := WorkspaceInput {
		Name: "testAddUserWorkspaceName5",
	}
	jsonInput, _ = json.Marshal(wi)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: "addUserWorkspaceTest5Name",
		AddUserName:   "testUser",
		RoleId:        2,
	}
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 6
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "addedUserWorkspaceTest6Name",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "addUserWorkspaceTest6Name",
		PassWord: "testPass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	wi = WorkspaceInput {
		Name: "testAddUserWorkspaceName6",
	}
	jsonInput, _ = json.Marshal(wi)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: "addUserWorkspaceTest6Name",
		AddUserName:   "addedUserWorkspaceTest6Name",
		RoleId:        4,
	}
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)


	// 7
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "addedUserInWorkspaceUser",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "testAddUserWorkspaceUser",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)
	jwtToken = rr.Body.String()
	fmt.Println(jwtToken)

	rr = httptest.NewRecorder()
	inputWorkspace := WorkspaceInput{
		Name: "testAddUserWorkspaceName",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: inputWorkspace.Name,
		AddUserName:   "addedUserInWorkspaceUser",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

}
