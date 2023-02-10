package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
)

var workspaceRouter = SetupRouter()

func createWorkSpaceTestFunc(workspaceName, jwtToken string, userId uint32) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	inputWorkspace := models.Workspace{
		Name:           workspaceName,
		PrimaryOwnerId: userId,
	}
	jsonInput, _ := json.Marshal(inputWorkspace)
	req, err := http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	if err != nil {
		return rr
	}
	req.Header.Set("Authorization", jwtToken)
	fmt.Println(req.Header.Get("Authorization"))
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func addUserWorkspaceTestFunc(workspaceId, roleId int, userId uint32, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	auwi := models.WorkspaceAndUsers{
		WorkspaceId: workspaceId,
		UserId:      userId,
		RoleId:      roleId,
	}
	jsonInput, _ := json.Marshal(auwi)
	req, _ := http.NewRequest("POST", "/api/workspace/add_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func renameWorkSpaceNameTestFunc(workspaceId, workspacePrimaryOwnerId int, newWorkspaceName, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	w := models.NewWorkspace(workspaceId, newWorkspaceName, uint32(workspacePrimaryOwnerId))
	jsonInput, _ := json.Marshal(w)
	req, _ := http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func deleteUserFromWorkspaceTestFunc(workspaceId, userId, roleId int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(models.WorkspaceAndUsers{
		WorkspaceId: workspaceId,
		UserId:      uint32(userId),
		RoleId:      roleId,
	})
	req, _ := http.NewRequest("POST", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func TestCreateWorkspace(t *testing.T) {
	// 1. 正常な状態(ログイン中のユーザーがworkspaceを作成する) 200
	// 2. jwtTokenから復元されるUserIdとbodyのprimaryOwnerUserIdが一致しない場合 400
	// 3. bodyにNameかPrimaryOwnerIdが含まれていない場合 400

	// 1
	t.Run("correctCase", func(t *testing.T) {
		userIds := []uint32{}
		UserNames := []string{}
		WorkSpaceNames := []string{}
		JwtTokens := []string{}

		for i := 0; i < 10; i++ {
			UserNames = append(UserNames, "createWorkSpaceTestUserName1"+strconv.Itoa(i))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, "createWorkspaceTestWorkspaceNames1"+strconv.Itoa(i))
		}

		for _, name := range UserNames {
			rr := signUpTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			rr = loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			lr := new(LoginResponse)
			json.Unmarshal(([]byte)(byteArray), lr)
			userIds = append(userIds, lr.UserId)
			JwtTokens = append(JwtTokens, lr.Token)
		}

		for i, workspaceName := range WorkSpaceNames {
			rr := createWorkSpaceTestFunc(workspaceName, JwtTokens[i%10], userIds[i%10])
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			w := new(models.Workspace)
			json.Unmarshal(([]byte)(byteArray), w)
			assert.Equal(t, workspaceName, w.Name)
			assert.Equal(t, userIds[i%10], w.PrimaryOwnerId)
			assert.NotEqual(t, 0, w.ID)
		}
	})

	// 2
	t.Run("2", func(t *testing.T) {
		userIds := []uint32{}
		UserNames := []string{}
		WorkSpaceNames := []string{}
		jwtTokens := []string{}

		for i := 0; i < 10; i++ {
			UserNames = append(UserNames, "createWorkSpaceTestUserName2"+strconv.Itoa(i))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, "createWorkspaceTestWorkspaceNames2"+strconv.Itoa(i))
		}

		for _, name := range UserNames {
			rr := signUpTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			u := new(models.User)
			json.Unmarshal(jsonBody, u)
			userIds = append(userIds, u.ID)
		}

		for i, name := range UserNames {
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			lr := new(LoginResponse)
			json.Unmarshal(jsonBody, lr)
			assert.Equal(t, userIds[i], lr.UserId)
			jwtTokens = append(jwtTokens, lr.Token)
		}

		for i, workspaceName := range WorkSpaceNames {
			assert.Equal(t, http.StatusBadRequest, createWorkSpaceTestFunc(workspaceName, jwtTokens[i%10], userIds[(i+1)%10]).Code)
		}
	})

	// 3
	t.Run("3", func(t *testing.T) {
		userIds := []uint32{}
		UserNames := []string{}
		WorkSpaceNames := []string{}
		jwtTokens := []string{}

		for i := 0; i < 10; i++ {
			UserNames = append(UserNames, "createWorkSpaceTestUserName3"+strconv.Itoa(i))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, "createWorkspaceTestWorkspaceNames3"+strconv.Itoa(i))
		}

		for i, name := range UserNames {
			rr := signUpTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			jsonBody := ([]byte)(byteArray)
			u := new(models.User)
			json.Unmarshal(jsonBody, u)
			userIds = append(userIds, u.ID)
			assert.Equal(t, UserNames[i], u.Name)

			rr = loginTestFunc(name, "pass")
			byteArray, _ = ioutil.ReadAll(rr.Body)
			jsonBody = ([]byte)(byteArray)
			lr := new(LoginResponse)
			json.Unmarshal(jsonBody, lr)
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, lr.Token)
		}

		for i, workspaceName := range WorkSpaceNames {
			var rr *httptest.ResponseRecorder
			if i%3 == 0 {
				rr = createWorkSpaceTestFunc("", jwtTokens[i%10], userIds[i%10])
			} else if i%3 == 1 {
				rr = createWorkSpaceTestFunc(workspaceName, jwtTokens[i%10], 0)
			} else {
				rr = createWorkSpaceTestFunc("", jwtTokens[i%10], 0)
			}
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}
	})
}

func TestAddUserInWorkspace(t *testing.T) {
	// 1. 正常な場合 200
	// 2. requestのbodyの情報に不足がある場合 400
	// 3. 存在しないworkspaceIdだった場合 400
	// 4. requestしたユーザーがrole = 1 or role = 2 or role = 3でない場合 400
	// 5. 追加されるユーザーがrole = 1の場合 400
	// 6. 既に登録されているユーザーを追加する場合 400

	// 1
	t.Run("1", func(t *testing.T) {
		ownerUserName := "AddUserInWorkspaceTestOwnerUser1"
		addUserName := "AddUserInWorkspaceTestUser1"
		workspaceName := "AddUserInWorkspaceTestWorkspace1"
		addUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, addUserRoleId, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		wau := new(models.WorkspaceAndUsers)
		json.Unmarshal(([]byte)(byteArray), wau)
		assert.Equal(t, alr.UserId, wau.UserId)
		assert.Equal(t, w.ID, wau.WorkspaceId)
		assert.Equal(t, addUserRoleId, wau.RoleId)
	})

	t.Run("2", func(t *testing.T) {
		ownerUserName := "AddUserInWorkspaceTestOwnerUser2"
		addUserName := "AddUserInWorkspaceTestUser2"
		workspaceName := "AddUserInWorkspaceTestWorkspace2"

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, 0, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"field empty\"}", rr.Body.String())
	})

	t.Run("3", func(t *testing.T) {
		ownerUserName := "AddUserInWorkspaceTestOwnerUser3"
		addUserName := "AddUserInWorkspaceTestUser3"
		workspaceName := "AddUserInWorkspaceTestWorkspace3"
		addUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(10000000000000000, addUserRoleId, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not found workspace\"}", rr.Body.String())
	})

	t.Run("4", func(t *testing.T) {
		ownerUserName := "AddUserInWorkspaceTestOwnerUser4"
		reqUserName := "reqUserInWorkspaceTestReqUser4"
		addUserName := "AddUserInWorkspaceTestUser4"
		workspaceName := "AddUserInWorkspaceTestWorkspace4"
		addUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(reqUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(reqUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, addUserRoleId, rlr.UserId, olr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		wau := new(models.WorkspaceAndUsers)
		json.Unmarshal(([]byte)(byteArray), wau)
		assert.Equal(t, rlr.UserId, wau.UserId)
		assert.Equal(t, w.ID, wau.WorkspaceId)
		assert.Equal(t, addUserRoleId, wau.RoleId)

		rr = addUserWorkspaceTestFunc(w.ID, addUserRoleId, alr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"Unauthorized add user in workspace\"}", rr.Body.String())

	})

	t.Run("5", func(t *testing.T) {
		ownerUserName := "AddUserInWorkspaceTestOwnerUser5"
		addUserName := "AddUserInWorkspaceTestUser5"
		workspaceName := "AddUserInWorkspaceTestWorkspace5"
		addUserRoleId := 1

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, addUserRoleId, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"can't add roleId = 1\"}", rr.Body.String())
	})

	t.Run("6", func(t *testing.T) {
		ownerUserName := "AddUserInWorkspaceTestOwnerUser6"
		addUserName := "AddUserInWorkspaceTestUser6"
		workspaceName := "AddUserInWorkspaceTestWorkspace6"
		addUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(addUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(addUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		alr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), alr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = addUserWorkspaceTestFunc(w.ID, addUserRoleId, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		wau := new(models.WorkspaceAndUsers)
		json.Unmarshal(([]byte)(byteArray), wau)
		assert.Equal(t, alr.UserId, wau.UserId)
		assert.Equal(t, w.ID, wau.WorkspaceId)
		assert.Equal(t, addUserRoleId, wau.RoleId)

		rr = addUserWorkspaceTestFunc(w.ID, 3, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"UNIQUE constraint failed: workspaces_and_users.workspace_id, workspaces_and_users.user_id\"}", rr.Body.String())

	})
}

func TestRenameWorkspaceName(t *testing.T) {
	// 1. 正常時 200
	// 2. headerに認証情報がない場合 400
	// 3. 認証が正常にできない場合 400
	// 4. 認証したユーザーがworkspaceに参加していない場合 400
	// 5. 認証したユーザーが対象のworkspaceでrole = 1 or role = 2 or role = 3のいずれか出ない場合 400
	// 6. 変更したいNameがすでに使用されていた場合 400

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

}
