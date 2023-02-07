package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"backend/models"
)

var workspaceRouter = SetupRouter()

func createWorkSpaceTestFunc(workspaceName, jwtToken string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	inputWorkspace := WorkspaceInput{
		Name: workspaceName,
	}
	jsonInput, _ := json.Marshal(inputWorkspace)
	req, _ := http.NewRequest("POST", "/api/workspace/create", bytes.NewBuffer(jsonInput))
	req.Header.Add("Authorization", jwtToken)
	workspaceRouter.ServeHTTP(rr, req)
	return rr
}

func addUserWorkspaceTestFunc(workspaceName, userName, jwtToken string, roleId int) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	auwi := AddUserWorkspaceInput{
		WorkspaceName: workspaceName,
		AddUserName:   userName,
		RoleId:        roleId,
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
	// 2. ログインしていない状態でのアクセス(headerにauthenticate情報がない) 400
	// 3. authenticate情報からUserIdが特定できない場合 400

	// 1
	t.Run("correctCase", func(t *testing.T) {
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
			JwtTokens = append(JwtTokens, rr.Body.String())
		}

		for i, workspaceName := range WorkSpaceNames {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceName, JwtTokens[i%10]).Code)
			fmt.Println(i)
		}
	})

	// 2
	t.Run("no login", func(t *testing.T) {
		UserNames := []string{}
		WorkSpaceNames := []string{}

		for i := 0; i < 10; i++ {
			UserNames = append(UserNames, "createWorkSpaceTestUserName2"+strconv.Itoa(i))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, "createWorkspaceTestWorkspaceNames2"+strconv.Itoa(i))
		}

		for _, name := range UserNames {
			rr := signUpTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
		}

		for _, workspaceName := range WorkSpaceNames {
			assert.Equal(t, http.StatusBadRequest, createWorkSpaceTestFunc(workspaceName, "").Code)
		}
	})

	// 3
	t.Run("no jwt token", func(t *testing.T) {
		UserNames := []string{}
		WorkSpaceNames := []string{}

		for i := 0; i < 10; i++ {
			UserNames = append(UserNames, "createWorkSpaceTestUserName3"+strconv.Itoa(i))
		}

		for i := 0; i < 100; i++ {
			WorkSpaceNames = append(WorkSpaceNames, "createWorkspaceTestWorkspaceNames3"+strconv.Itoa(i))
		}

		for _, name := range UserNames {
			rr := signUpTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			rr = loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
		}

		for _, workspaceName := range WorkSpaceNames {
			assert.Equal(t, http.StatusBadRequest, createWorkSpaceTestFunc(workspaceName, "abc").Code)
		}
	})
}

func TestAddUserWorkspace(t *testing.T) {
	// 1. 正常な場合 200
	// 2. jwtTokenがheaderに含まれていない場合 400
	// 3. jwtTokenからuserIdが取得できない場合 400
	// 4. requestのbodyの情報に不足がある場合 400
	// 5. 存在しないworkspaceNameかuserNameだった場合 400
	// 6. role = 4の場合 400

	// 1
	t.Run("test case 1", func(t *testing.T) {
		primaryUserNames := []string{}
		userNames := []string{}
		jwtTokens := []string{}
		workspaceNames := []string{}

		for i := 0; i < 10; i++ {
			primaryUserNames = append(primaryUserNames, "addUserWorkspacePrimaryUserName1"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceName1"+strconv.Itoa(i))
		}

		for _, name := range primaryUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, rr.Body.String())
		}

		for i := 0; i < 100; i++ {
			userNames = append(userNames, "addUserWorkspaceUserName1"+strconv.Itoa(i))
		}

		for _, name := range userNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}

		for i, jwtToken := range jwtTokens {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceNames[i], jwtToken).Code)
		}

		for i, jwtToken := range jwtTokens {
			for _, userName := range userNames {
				assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(workspaceNames[i], userName, jwtToken, 4).Code)
			}
		}
	})

	// 2
	t.Run("test case 2", func(t *testing.T) {
		primaryUserNames := []string{}
		userNames := []string{}
		jwtTokens := []string{}
		workspaceNames := []string{}

		for i := 0; i < 10; i++ {
			primaryUserNames = append(primaryUserNames, "addUserWorkspacePrimaryUserName2"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceName2"+strconv.Itoa(i))
		}

		for _, name := range primaryUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, rr.Body.String())
		}

		for i := 0; i < 100; i++ {
			userNames = append(userNames, "addUserWorkspaceUserName2"+strconv.Itoa(i))
		}

		for _, name := range userNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}

		for i, jwtToken := range jwtTokens {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceNames[i], jwtToken).Code)
		}

		for _, workspaceName := range workspaceNames {
			for _, userName := range userNames {
				assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc(workspaceName, userName, "", 4).Code)
			}
		}
	})

	// 3
	t.Run("test case 3", func(t *testing.T) {
		primaryUserNames := []string{}
		userNames := []string{}
		jwtTokens := []string{}
		workspaceNames := []string{}

		for i := 0; i < 10; i++ {
			primaryUserNames = append(primaryUserNames, "addUserWorkspacePrimaryUserName3"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceName3"+strconv.Itoa(i))
		}

		for _, name := range primaryUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, rr.Body.String())
		}

		for i := 0; i < 100; i++ {
			userNames = append(userNames, "addUserWorkspaceUserName3"+strconv.Itoa(i))
		}

		for _, name := range userNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}

		for i, jwtToken := range jwtTokens {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceNames[i], jwtToken).Code)
		}

		for _, workspaceName := range workspaceNames {
			for _, userName := range userNames {
				assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc(workspaceName, userName, "jwtToken", 4).Code)
			}
		}
	})

	// 4
	t.Run("test case 4", func(t *testing.T) {
		primaryUserNames := []string{}
		userNames := []string{}
		jwtTokens := []string{}
		workspaceNames := []string{}

		for i := 0; i < 10; i++ {
			primaryUserNames = append(primaryUserNames, "addUserWorkspacePrimaryUserName4"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceName4"+strconv.Itoa(i))
		}

		for _, name := range primaryUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, rr.Body.String())
		}

		for i := 0; i < 100; i++ {
			userNames = append(userNames, "addUserWorkspaceUserName4"+strconv.Itoa(i))
		}

		for _, name := range userNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}

		for i, jwtToken := range jwtTokens {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceNames[i], jwtToken).Code)
		}

		for i, jwtToken := range jwtTokens {
			for _, userName := range userNames {
				if i%3 == 0 {
					assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc(workspaceNames[i], userName, jwtToken, 0).Code)
				}
				if i%3 == 1 {
					assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc(workspaceNames[i], "", jwtToken, 4).Code)
				}
				if i%3 == 2 {
					assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc("", userName, jwtToken, 4).Code)
				}
			}
		}
	})

	// 5
	t.Run("test case 5", func(t *testing.T) {
		primaryUserNames := []string{}
		userNames := []string{}
		jwtTokens := []string{}
		workspaceNames := []string{}

		for i := 0; i < 10; i++ {
			primaryUserNames = append(primaryUserNames, "addUserWorkspacePrimaryUserName5"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceName5"+strconv.Itoa(i))
		}

		for _, name := range primaryUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, rr.Body.String())
		}

		for i := 0; i < 100; i++ {
			userNames = append(userNames, "addUserWorkspaceUserName5"+strconv.Itoa(i))
		}

		for _, name := range userNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}

		for i, jwtToken := range jwtTokens {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceNames[i], jwtToken).Code)
		}

		for i, jwtToken := range jwtTokens {
			for j, userName := range userNames {
				if j%2 == 0 {
					assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc(workspaceNames[i], "wrongUserName", jwtToken, 4).Code)
				}
				if j%2 == 1 {
					assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc("wrongWorkspaceName", userName, jwtToken, 4).Code)
				}
			}
		}
	})

	// 6
	t.Run("test case 6", func(t *testing.T) {
		primaryUserNames := []string{}
		userNames := []string{}
		workspaceUserNames := []string{}
		jwtTokens := []string{}
		primaryUserJwtTokens := []string{}
		workspaceNames := []string{}

		for i := 0; i < 10; i++ {
			primaryUserNames = append(primaryUserNames, "addUserWorkspacePrimaryUserName6"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceUserName6.1"+strconv.Itoa(i))
			workspaceNames = append(workspaceNames, "addUserWorkspaceName6"+strconv.Itoa(i))
		}

		for _, name := range primaryUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			primaryUserJwtTokens = append(primaryUserJwtTokens, rr.Body.String())
		}

		for _, name := range workspaceUserNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
			rr := loginTestFunc(name, "pass")
			assert.Equal(t, http.StatusOK, rr.Code)
			jwtTokens = append(jwtTokens, rr.Body.String())
		}

		for i := 0; i < 100; i++ {
			userNames = append(userNames, "addUserWorkspaceUserName6"+strconv.Itoa(i))
		}

		for _, name := range userNames {
			assert.Equal(t, http.StatusOK, signUpTestFunc(name, "pass").Code)
		}

		for i, jwtToken := range primaryUserJwtTokens {
			assert.Equal(t, http.StatusOK, createWorkSpaceTestFunc(workspaceNames[i], jwtToken).Code)
		}

		for i, jwtToken := range jwtTokens {
			assert.Equal(t, http.StatusOK, addUserWorkspaceTestFunc(workspaceNames[i], workspaceUserNames[i], jwtToken, 4).Code)

		}

		for i, jwtToken := range jwtTokens {
			for _, userName := range userNames {
				assert.Equal(t, http.StatusBadRequest, addUserWorkspaceTestFunc(workspaceNames[i], userName, jwtToken, 4).Code)
			}
		}
	})
}

func TestRenameWorkspaceName(t *testing.T) {
	// 1. 正常時 200
	// 2. headerに認証情報がない場合 400
	// 3. 認証が正常にできない場合 400
	// 4. 認証したユーザーがworkspaceに参加していない場合 400
	// 5. 認証したユーザーが対象のworkspaceでrole = 1 or role = 2 or role = 3のいずれか出ない場合 400
	// 6. 変更したいNameがすでに使用されていた場合 400

	// 1


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
