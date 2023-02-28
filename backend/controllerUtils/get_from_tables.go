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

func GetUserInWorkspace(workspaceId int) ([]UserInfoInWorkspace, error) {
	// workspace内のuserの情報を配列にして返す
	// userぞれぞれの情報は返り値のstructを参照
	// アクセスしたuserの情報も含めて返す

	res := make([]UserInfoInWorkspace, 0)

	// workspaces_and_users tableからworkspace_idが等しいものをすべて取得する
	waus, err := models.GetWAUsByWorkspaceId(workspaceId)
	if err != nil {
		return res, err
	}

	// users tableの全情報を取得する
	users, err := models.GetUsers()
	if err != nil {
		return res, err
	}

	// 2つのデータからuser_idが等しいものの組み合わせを見つけて、res配列に追加する
	for _, wau := range waus {
		for _, user := range users {
			if user.ID == wau.UserId {
				res = append(res, UserInfoInWorkspace{
					ID:     user.ID,
					Name:   user.Name,
					RoleId: wau.RoleId,
				})
				break
			}
		}
	}
	return res, nil
}
