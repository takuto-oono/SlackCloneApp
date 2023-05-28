import router from "next/router";
import { getToken,resetCookie,getUserId } from "./cookie";

export interface Workspace {
    id: number;
    name: string;
    primary_owner_id: number;
}

export interface UserInWorkspace {
  id: number;
  name: string;
  roleid: number;
}

const baseUrl = 'http://localhost:8080/api/workspace/'

export async function getWorkspaces(): Promise<Workspace[]> {
  const url = baseUrl + 'get_by_user'
  console.log("getWorkspaces")
  // console.log(getToken());
  let res_workspaces: Workspace[]
  const workspaces = [
    {
      id: 0,
      name: "",
      primary_owner_id: 0
    }
  ]
  
  try {
      const res = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': getToken(),
        },
      })
    // 認証エラーの時のみリダイレクトする
    if (res.status == 401) {
      resetCookie();
      console.log("redirect");
      router.push("/")
    }
      res_workspaces = await res.json()
      return new Promise((resolve) => {
      const workspaces: Workspace[] = res_workspaces;
      resolve(workspaces);
    });
  } catch (err) {
    console.log(err)
  }
  // console.log(workspaces);
  return workspaces;
}


export async function postWorkspace(workspaceName:string){
  const url = baseUrl + 'create'
  console.log(getUserId());
    try {
      const res = await fetch(url, {
        method: 'POST',
        headers: {
          'Authorization': getToken(),
        },
        body: JSON.stringify({
          name: workspaceName,
          user_id: getUserId(),
        })
      })
      console.log(res)
      if (res.status == 401) {
        resetCookie();
        console.log("redirect");
        router.push("/")
      }
    } catch (err) {
      console.log(err)
    }
  return
}

export async function getUsersInWorkspace(workspaceId: number): Promise<UserInWorkspace[]> {
  const url = baseUrl + "get_users/" + workspaceId;
  let resUsersInWorkspace: UserInWorkspace[];
  const usersInWorkspace = [{
    id: 0,
    name: "",
    roleid: 0,
  }];

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
    resUsersInWorkspace = await res.json();
    return new Promise((resolve) => {
      const usersInWorkspace: UserInWorkspace[] = resUsersInWorkspace;
      resolve(usersInWorkspace);
    });
  } catch (err) {
    console.log(err);
  }
  return usersInWorkspace;
}

export async function addUserInWorkspace(workspaceID: number, userID: number): Promise<void> {
  const url = baseUrl + "add_user";
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        workspace_id: workspaceID,
        user_id: userID,
        // とりあえず、4固定にしておく
        role_id: 4,
      }),
    })
    console.log(res);
  } catch(e) {
    console.log(e);
  }
}
