package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
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
	assert.Equal(t, http.StatusOK, rr.Code)

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
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 3
	rr = httptest.NewRecorder()
	inputWorkspace = WorkspaceInput{
		Name: "testCreateWorkspaceCantAuthenticate",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", "testJWTToken")
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 4
	// body自体がない場合
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/workspace/create", nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// bodyはあるが、Nameに何も指定がない場合
	rr = httptest.NewRecorder()
	inputWorkspace = WorkspaceInput{
		Name: "",
	}
	jsonInput, _ = json.Marshal(inputWorkspace)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
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
	wi := WorkspaceInput{
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
	wi = WorkspaceInput{
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

func TestRenameWorkspaceName(t *testing.T) {
	// 1. 正常時 200
	// 2. headerに認証情報がない場合 400
	// 3. 認証が正常にできない場合 400
	// 4. 認証したユーザーがworkspaceに参加していない場合 400
	// 5. 認証したユーザーが対象のworkspaceでrole = 1 or role = 2 or role = 3のいずれか出ない場合 400
	// 6. 変更したいNameがすでに使用されていた場合 400

	// 1
	rr := httptest.NewRecorder()
	input := UserInput{
		Name:     "renameWorkspaceNameUser1",
		PassWord: "pass",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken := rr.Body.String()

	rr = httptest.NewRecorder()
	workspaceInput := WorkspaceInput{
		Name: "renameWorkspaceNameOld1",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ := ioutil.ReadAll(rr.Body)
	jsonBody := ([]byte)(byteArray)
	w := new(models.Workspace)
	err := json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	w.Name = "renameWorkspaceNameNew1"
	jsonInput, _ = json.Marshal(w)
	req, _ = http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// 2
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "renameWorkspaceNameUser2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "renameWorkspaceNameOld2",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	w.Name = "renameWorkspaceNameNew3"
	jsonInput, _ = json.Marshal(w)
	req, _ = http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 3
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "renameWorkspaceNameUser3",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "renameWorkspaceNameOld3",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	w.Name = "renameWorkspaceNameNew3"
	jsonInput, _ = json.Marshal(w)
	req, _ = http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", "wrong jwt token")
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	//4
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "renameWorkspaceNameUser4",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	input2 := UserInput{
		Name:     "renameWorkspaceNameUser4.2",
		PassWord: "pass",
	}

	jsonInput, _ = json.Marshal(input2)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input2)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken2 := rr.Body.String()

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "renameWorkspaceNameOld4",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	w.Name = "renameWorkspaceNameNew4"
	jsonInput, _ = json.Marshal(w)
	req, _ = http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken2)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 5
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "renameWorkspaceNameUser5",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	input2 = UserInput{
		Name:     "renameWorkspaceNameUser5.2",
		PassWord: "pass",
	}

	jsonInput, _ = json.Marshal(input2)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input2)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken2 = rr.Body.String()

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "renameWorkspaceNameOld5",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	addUserWorkspaceInput := AddUserWorkspaceInput{
		WorkspaceName: w.Name,
		AddUserName:   input2.Name,
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(addUserWorkspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	w.Name = "renameWorkspaceNameNew5"
	jsonInput, _ = json.Marshal(w)
	req, _ = http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken2)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 6
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "renameWorkspaceNameUser6",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "renameWorkspaceNameNew6",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "renameWorkspaceNameOld6",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	w.Name = "renameWorkspaceNameNew6"
	jsonInput, _ = json.Marshal(w)
	req, _ = http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.NotEqual(t, http.StatusOK, rr.Code)
}

func TestDeleteUserFromWorkSpace(t *testing.T) {
	// 1. 正常時 200
	// 2. headerにJWTTokenがない場合 400
	// 3. JWTTokenでuserIdが復元できない場合 400
	// 4. bodyにworkspaceIdが0の場合 400
	// 5. bodyにuserIdがない場合 400
	// 6. bodyにroleIdがない場合 400
	// 7. requestしたuserのrole = 4の場合
	// 8. 削除されるユーザーのrole = 1の場合 400
	// 9. dbに一致する情報が存在しない場合 400

	// 1
	rr := httptest.NewRecorder()
	input := UserInput{
		Name:     "deleteUserFromWorkSpaceUser1",
		PassWord: "pass",
	}
	jsonInput, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken := rr.Body.String()

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkspaceNameUser1.2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput := WorkspaceInput{
		Name: "deleteUserFromWorkspace1",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ := ioutil.ReadAll(rr.Body)
	jsonBody := ([]byte)(byteArray)
	w := new(models.Workspace)
	err := json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	auwi := AddUserWorkspaceInput{
		WorkspaceName: workspaceInput.Name,
		AddUserName:   "deleteUserFromWorkspaceNameUser1.2",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	wau := new(models.WorkspaceAndUsers)
	err = json.Unmarshal(jsonBody, wau)
	assert.Empty(t, err)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: w.ID,
		UserId:      wau.UserId,
		RoleId:      wau.RoleId,
	})
	req, _ = http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// 2
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkSpaceUser2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkspaceNameUser2.2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "deleteUserFromWorkspace2",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: workspaceInput.Name,
		AddUserName:   "deleteUserFromWorkspaceNameUser2.2",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	wau = new(models.WorkspaceAndUsers)
	err = json.Unmarshal(jsonBody, wau)
	assert.Empty(t, err)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: w.ID,
		UserId:      wau.UserId,
		RoleId:      wau.RoleId,
	})
	req, _ = http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 3
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkSpaceUser3",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkspaceNameUser3.2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "deleteUserFromWorkspace3",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: workspaceInput.Name,
		AddUserName:   "deleteUserFromWorkspaceNameUser3.2",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	wau = new(models.WorkspaceAndUsers)
	err = json.Unmarshal(jsonBody, wau)
	assert.Empty(t, err)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: w.ID,
		UserId:      wau.UserId,
		RoleId:      wau.RoleId,
	})
	req, _ = http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", "abc")
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 4
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkSpaceUser4",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkspaceNameUser4.2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "deleteUserFromWorkspace4",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: workspaceInput.Name,
		AddUserName:   "deleteUserFromWorkspaceNameUser4.2",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	wau = new(models.WorkspaceAndUsers)
	err = json.Unmarshal(jsonBody, wau)
	assert.Empty(t, err)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: 0,
		UserId:      wau.UserId,
		RoleId:      wau.RoleId,
	})
	req, _ = http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 5
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkSpaceUser5",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkspaceNameUser5.2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "deleteUserFromWorkspace5",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: workspaceInput.Name,
		AddUserName:   "deleteUserFromWorkspaceNameUser5.2",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	wau = new(models.WorkspaceAndUsers)
	err = json.Unmarshal(jsonBody, wau)
	assert.Empty(t, err)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: w.ID,
		RoleId:      wau.RoleId,
	})
	req, _ = http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// 6
	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkSpaceUser6",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	jwtToken = rr.Body.String()

	rr = httptest.NewRecorder()
	input = UserInput{
		Name:     "deleteUserFromWorkspaceNameUser6.2",
		PassWord: "pass",
	}
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/signUp", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(input)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonInput))
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	workspaceInput = WorkspaceInput{
		Name: "deleteUserFromWorkspace6",
	}
	jsonInput, _ = json.Marshal(workspaceInput)
	req, _ = http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	w = new(models.Workspace)
	err = json.Unmarshal(jsonBody, w)
	assert.Empty(t, err)
	assert.Equal(t, workspaceInput.Name, w.Name)
	assert.NotEmpty(t, w.ID)
	assert.NotEmpty(t, w.PrimaryOwnerId)

	rr = httptest.NewRecorder()
	auwi = AddUserWorkspaceInput{
		WorkspaceName: workspaceInput.Name,
		AddUserName:   "deleteUserFromWorkspaceNameUser6.2",
		RoleId:        4,
	}
	jsonInput, _ = json.Marshal(auwi)
	req, _ = http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	byteArray, _ = ioutil.ReadAll(rr.Body)
	jsonBody = ([]byte)(byteArray)
	wau = new(models.WorkspaceAndUsers)
	err = json.Unmarshal(jsonBody, wau)
	assert.Empty(t, err)

	rr = httptest.NewRecorder()
	jsonInput, _ = json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: w.ID,
		UserId:      wau.UserId,
	})
	req, _ = http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

}
