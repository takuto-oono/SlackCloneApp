import router from "next/router";
import { getToken, resetCookie, getUserId } from "./cookie";

export interface UserInfo {
  id: number;
  name: string;
  roleid: number;
}

export interface Workspace {
  id: number;
  name: string;
  primary_owner_id: number;
}

const baseUrl = "http://localhost:8080/api/workspace/";

export async function getWorkspaces(): Promise<Workspace[]> {
  const url = baseUrl + "get_by_user";
  console.log("getWorkspaces");
  console.log(getToken());
  let res_workspaces: Workspace[];
  const workspaces = [
    {
      id: 0,
      name: "",
      primary_owner_id: 0,
    },
  ];

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    // 認証エラーの時のみリダイレクトする
    if (res.status == 401) {
      resetCookie();
      console.log("redirect");
      router.replace("/");
    }
    res_workspaces = await res.json();
    return new Promise((resolve) => {
      const workspaces: Workspace[] = res_workspaces;
      resolve(workspaces);
    });
  } catch (err) {
    console.log(err);
  }
  console.log(workspaces);
  return workspaces;
}

export async function postWorkspace(workspaceName: string) {
  const url = baseUrl + "create";
  let workspace: Workspace;
  console.log(getUserId());
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        name: workspaceName,
        user_id: getUserId(),
      }),
    });
    console.log(res);
    workspace = await res.json();
    console.log(workspace);
    console.log("status-code");
    console.log(res.status);

    if (res.status == 401) {
      console.log("作成失敗");
      router.push("/");
    } else if (res.status == 200) {
      console.log("redirect");
      router.push("/");
    }
  } catch (err) {
    console.log(err);
  }
  return;
}

export async function getUsers(workspace_id: number): Promise<UserInfo[]> {
  const url = baseUrl + "get_user_in_workspace" + workspace_id;
  console.log("getUserInWorkspace");
  let res_user_in_workspaces: UserInfo[];
  const user_in_workspaces = [
    {
      id: 0,
      name: "",
      roleid: 0,
    },
  ];

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    // 認証エラーの時のみリダイレクトする
    if (res.status == 401) {
      resetCookie();
      console.log("redirect");
      router.replace("/");
    }
    res_user_in_workspaces = await res.json();
    return new Promise((resolve) => {
      const user_in_workspaces: UserInfo[] = res_user_in_workspaces;
      resolve(user_in_workspaces);
    });
  } catch (err) {
    console.log(err);
  }
  console.log(user_in_workspaces);
  return user_in_workspaces;
}
