package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var workspaceRouter = SetupRouter()

func createWorkSpaceTestFunc(workspaceName, jwtToken string, userId uint32) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	inputWorkspace := controllerUtils.CreateWorkspaceInput{
		Name:          workspaceName,
		RequestUserId: userId,
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
	auwi := controllerUtils.AddUserInWorkspaceInput{
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

// func renameWorkSpaceNameTestFunc(workspaceId, workspacePrimaryOwnerId int, newWorkspaceName, jwtToken string) *httptest.ResponseRecorder {
// 	rr := httptest.NewRecorder()
// 	w := models.NewWorkspace(workspaceId, newWorkspaceName, uint32(workspacePrimaryOwnerId))
// 	jsonInput, _ := json.Marshal(w)
// 	req, _ := http.NewRequest("POST", "/api/workspace/rename", bytes.NewBuffer(jsonInput))
// 	req.Header.Add("Authorization", jwtToken)
// 	workspaceRouter.ServeHTTP(rr, req)
// 	return rr
// }

func deleteUserFromWorkspaceTestFunc(workspaceId int, userId uint32, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	jsonInput, _ := json.Marshal(controllerUtils.DeleteUserFromWorkSpaceInput{
		WorkspaceId: workspaceId,
		UserId:      userId,
	})
	req, _ := http.NewRequest("DELETE", "/api/workspace/delete_user", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func getWorkspacesByUserIdTestFunc(userId uint32, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/workspace/get_by_user", nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func GetUsersInWorkspaceTestFunc(workspaceId int, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/workspace/get_users/"+strconv.Itoa(workspaceId), nil)
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func TestCreateWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な状態(ログイン中のユーザーがworkspaceを作成する) 200
	// 2. jwtTokenから復元されるUserIdとbodyのprimaryOwnerUserIdが一致しない場合 400
	// 3. bodyにNameかPrimaryOwnerIdが含まれていない場合 400
	// 4. 既に同じ名前のworkspaceが存在する場合 409

	// 1
	t.Run("correctCase", func(t *testing.T) {
		userIds := []uint32{}
		UserNames := []string{}
		WorkSpaceNames := []string{}
		JwtTokens := []string{}

		for i := 0; i < 10; i++ {
			UserNames = append(UserNames, randomstring.EnglishFrequencyString(30))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, randomstring.EnglishFrequencyString(30))
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
			UserNames = append(UserNames, randomstring.EnglishFrequencyString(30))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, randomstring.EnglishFrequencyString(30))
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
			UserNames = append(UserNames, randomstring.EnglishFrequencyString(30))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, randomstring.EnglishFrequencyString(30))
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

	t.Run("4 同じ名前のworkspaceが存在している場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), &lr)

		assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId).Code)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusConflict, rr.Code)
		assert.Equal(t, "{\"message\":\"UNIQUE constraint failed: workspaces.name\"}", rr.Body.String())
	})
}

func TestAddUserInWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. requestのbodyの情報に不足がある場合 400
	// 3. 存在しないworkspaceIdだった場合 404
	// 4. requestしたユーザーがrole = 1 or role = 2 or role = 3でない場合 403
	// 5. 追加されるユーザーがrole = 1の場合 400
	// 6. 既に登録されているユーザーを追加する場合 409

	// 1
	t.Run("1", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		addUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
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
		ownerUserName := randomstring.EnglishFrequencyString(30)
		addUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

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
		assert.Equal(t, "{\"message\":\"role_id not found\"}", rr.Body.String())
	})

	t.Run("3", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		addUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
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

		rr = addUserWorkspaceTestFunc(rand.Int(), addUserRoleId, alr.UserId, olr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"workspace not found\"}", rr.Body.String())
	})

	t.Run("4", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		reqUserName := randomstring.EnglishFrequencyString(30)
		addUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
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
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"Unauthorized add user in workspace\"}", rr.Body.String())

	})

	t.Run("5", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		addUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
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
		ownerUserName := randomstring.EnglishFrequencyString(30)
		addUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
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
		assert.Equal(t, http.StatusConflict, rr.Code)
		assert.Equal(t, "{\"message\":\"user already exist in workspace\"}", rr.Body.String())
	})
}

func TestRenameWorkspaceName(t *testing.T) {
	// 1. 正常時 200
	// 2. headerに認証情報がない場合 400
	// 3. 認証が正常にできない場合 401
	// 4. 認証したユーザーがworkspaceに参加していない場合 404
	// 5. 認証したユーザーが対象のworkspaceでrole = 1 or role = 2 or role = 3のいずれか出ない場合 403
	// 6. 変更したいNameがすでに使用されていた場合 409

}

func TestDeleteUserFromWorkSpace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常時 200
	// 2. bodyにworkspaceId, userId, roleIdのいずれかが含まれていない場合 400
	// 3. requestしたuserのrole = 4の場合 403
	// 4. 削除されるユーザーのrole = 1の場合 400
	// 5. 該当するUserがいない場合 404

	// 1
	t.Run("1", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		deleteUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		deleteUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(deleteUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(deleteUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		dlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), dlr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, deleteUserRoleId, dlr.UserId, olr.Token).Code)

		rr = deleteUserFromWorkspaceTestFunc(w.ID, dlr.UserId, olr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		wau := new(models.WorkspaceAndUsers)
		json.Unmarshal(([]byte)(byteArray), wau)
		assert.Equal(t, dlr.UserId, wau.UserId)
		assert.Equal(t, w.ID, wau.WorkspaceId)
		assert.Equal(t, deleteUserRoleId, wau.RoleId)
	})

	// 2
	t.Run("2", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		deleteUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		deleteUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(deleteUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(deleteUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		dlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), dlr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, deleteUserRoleId, dlr.UserId, olr.Token).Code)

		rr = deleteUserFromWorkspaceTestFunc(w.ID, 0, olr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"user_id not found\"}", rr.Body.String())

		rr = deleteUserFromWorkspaceTestFunc(0, dlr.UserId, olr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"workspace_id not found\"}", rr.Body.String())
	})

	// 3
	t.Run("3", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		reqUserName := randomstring.EnglishFrequencyString(30)
		deleteUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		deleteUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(reqUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(deleteUserName, "pass").Code)

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

		rr = loginTestFunc(deleteUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		dlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), dlr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 4, rlr.UserId, olr.Token).Code)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, deleteUserRoleId, dlr.UserId, olr.Token).Code)

		rr = deleteUserFromWorkspaceTestFunc(w.ID, dlr.UserId, rlr.Token)
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Equal(t, "{\"message\":\"not permission\"}", rr.Body.String())
	})

	//4
	t.Run("4", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		reqUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(reqUserName, "pass").Code)

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

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, 2, rlr.UserId, olr.Token).Code)

		rr = deleteUserFromWorkspaceTestFunc(w.ID, olr.UserId, rlr.Token)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "{\"message\":\"not delete primary owner\"}", rr.Body.String())
	})

	// 5
	t.Run("5", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		deleteUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)
		deleteUserRoleId := 4

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		assert.Equal(t, http.StatusOK, signUpTestFunc(deleteUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		olr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), olr)

		rr = loginTestFunc(deleteUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		dlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), dlr)

		rr = createWorkSpaceTestFunc(workspaceName, olr.Token, olr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, deleteUserRoleId, dlr.UserId, olr.Token).Code)

		rr = deleteUserFromWorkspaceTestFunc(w.ID, 441553453, olr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"record not found\"}", rr.Body.String())

		rr = deleteUserFromWorkspaceTestFunc(5934759792, dlr.UserId, olr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"record not found\"}", rr.Body.String())
	})
}

func TestGetWorkspacesById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. workspaceが存在する場合 200
	// 2. workspaceが存在しない場合 200

	t.Run("1 workspaceが存在する場合", func(t *testing.T) {
		workspaceCount := 10
		userName := randomstring.EnglishFrequencyString(30)
		workspaceNames := make([]string, workspaceCount)
		workspaces := make([]models.Workspace, workspaceCount)
		for i := 0; i < workspaceCount; i++ {
			workspaceNames[i] = randomstring.EnglishFrequencyString(30)
		}

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		for i, workspaceName := range workspaceNames {
			rr := createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			w := new(models.Workspace)
			json.Unmarshal(([]byte)(byteArray), w)
			workspaces[i] = *w

		}

		rr = getWorkspacesByUserIdTestFunc(lr.UserId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		byteArray, _ = ioutil.ReadAll(rr.Body)
		ws := make([]models.Workspace, workspaceCount)
		json.Unmarshal(([]byte)(byteArray), &ws)

		for _, w := range ws {
			assert.Contains(t, workspaces, w)
		}

		for i := 0; i < workspaceCount; i++ {
			for j := i + 1; j < workspaceCount; j++ {
				assert.NotEqual(t, ws[i], ws[j])
			}
		}
	})

	t.Run("2 workspaceが存在しない場合", func(t *testing.T) {
		userName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(userName, "pass").Code)

		rr := loginTestFunc(userName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = getWorkspacesByUserIdTestFunc(lr.UserId, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)

		byteArray, _ = ioutil.ReadAll(rr.Body)
		ws := make([]models.Workspace, 5)
		json.Unmarshal(([]byte)(byteArray), &ws)
		assert.Equal(t, 0, len(ws))

	})
}

func TestGetUsersInWorkspace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 1. 正常な場合 200
	// 2. workspaceにいないuserからのアクセスだった場合 404

	t.Run("1 正常な場合", func(t *testing.T) {
		userCount := 10
		ownerUserName := randomstring.EnglishFrequencyString(30)
		userInfos := make([]controllerUtils.UserInfoInWorkspace, userCount)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		for i := 0; i < userCount; i++ {
			userInfos[i] = controllerUtils.UserInfoInWorkspace{
				ID:     0,
				Name:   randomstring.EnglishFrequencyString(30),
				RoleId: rand.Int()%3 + 2,
			}

		}

		for i, ui := range userInfos {
			assert.Equal(t, http.StatusOK, signUpTestFunc(ui.Name, "pass").Code)

			rr := loginTestFunc(ui.Name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			byteArray, _ := ioutil.ReadAll(rr.Body)
			ulr := new(LoginResponse)
			json.Unmarshal(([]byte)(byteArray), ulr)
			fmt.Println(ulr.UserId)
			userInfos[i].ID = ulr.UserId
		}

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		for _, ui := range userInfos {
			assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(w.ID, ui.RoleId, ui.ID, lr.Token).Code)
		}

		rr = GetUsersInWorkspaceTestFunc(w.ID, lr.Token)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		res := make([]controllerUtils.UserInfoInWorkspace, 0)
		json.Unmarshal(([]byte)(byteArray), &res)

		assert.Equal(t, userCount+1, len(res))
		userInfos = append(userInfos, controllerUtils.UserInfoInWorkspace{
			ID:     lr.UserId,
			Name:   lr.Username,
			RoleId: 1,
		})

		for _, ui := range userInfos {
			assert.Contains(t, res, ui)
		}
	})

	t.Run("2 workspaceにいないuserからのアクセスの場合", func(t *testing.T) {
		ownerUserName := randomstring.EnglishFrequencyString(30)
		requestUserName := randomstring.EnglishFrequencyString(30)
		workspaceName := randomstring.EnglishFrequencyString(30)

		assert.Equal(t, http.StatusOK, signUpTestFunc(ownerUserName, "pass").Code)
		assert.Equal(t, http.StatusOK, signUpTestFunc(requestUserName, "pass").Code)

		rr := loginTestFunc(ownerUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ := ioutil.ReadAll(rr.Body)
		lr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), lr)

		rr = loginTestFunc(requestUserName, "pass")
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		rlr := new(LoginResponse)
		json.Unmarshal(([]byte)(byteArray), rlr)

		rr = createWorkSpaceTestFunc(workspaceName, lr.Token, lr.UserId)
		assert.Equal(t, http.StatusOK, rr.Code)
		byteArray, _ = ioutil.ReadAll(rr.Body)
		w := new(models.Workspace)
		json.Unmarshal(([]byte)(byteArray), w)

		rr = GetUsersInWorkspaceTestFunc(w.ID, rlr.Token)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Equal(t, "{\"message\":\"user not found in workspace\"}", rr.Body.String())
	})
}
