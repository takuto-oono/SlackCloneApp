package controllerUtils

import "backend/models"

type UserInfoInWorkspace struct {
	ID     uint32 `json:"id"`
	Name   string `json:"name"`
	RoleId int    `json:"role_id"`
}

func GetWorkspacesByUserId(userId uint32) ([]models.Workspace, error) {
	// 引数で指定したuserIdのuserが所属しているworkspaceのstructを配列にして返す

	// 結果を保存する配列を作成
	workspaces := make([]models.Workspace, 0)

	// workspaces_and_users tableからuserIdが等しいものをすべて取得
	waus, err := models.GetWAUsByUserId(userId)
	if err != nil {
		return workspaces, err
	}

	// workspaces tableからidが等しいworkspaceの情報を取得する。
	for _, wau := range waus {
		workspace, err := models.GetWorkspaceById(wau.WorkspaceId)
		if err != nil {
			return workspaces, err
		}
		workspaces = append(workspaces, workspace)
	}
	return workspaces, err
}

func GetUserInWorkspaceUser(workspaceId int) ([]UserInfoInWorkspace, error) {
	// workspace内のuserの情報を配列にして返す
	// userぞれぞれの情報は返り値のstructを参照
	// アクセスしたuser
}
