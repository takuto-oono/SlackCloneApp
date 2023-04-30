import router from "next/router";
import { getToken,resetCookie,getUserId } from "./cookie";

export interface Workspace {
    id: number;
    name: string;
    primary_owner_id: number;
}

const baseUrl = 'http://localhost:8080/api/workspace/'

export async function getWorkspaces(): Promise<Workspace[]> {
  const url = baseUrl + 'get_by_user'
  console.log("getWorkspaces")
  console.log(getToken());
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
  console.log(workspaces);
  return workspaces;
}


export async function postWorkspace(workspaceName:string){
    const url = baseUrl + 'create'
  let workspace: Workspace
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
      workspace = await res.json()
      console.log(workspace)
      console.log("status-code")
      console.log(res.status);
      
      if (res.status == 401) {
        console.log("作成失敗");
        // router.push('/')
      } else if (res.status == 200) {
        console.log("redirect");
        // router.push('/')
      }
      
    } catch (err) {
      console.log(err)
    }
    return
}
