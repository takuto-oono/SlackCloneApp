import { getToken,getUserId } from "./cookie";

export interface Workspace {
    id: number;
    name: string;
    primary_owner_id: number;
}

type Authorization = string | undefined;



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
        
        // console.log(res)
        res_workspaces = await res.json()
        // console.log("workspaces1")
        // console.log(res_workspaces);
      
        return new Promise((resolve) => {
        const workspaces: Workspace[] = res_workspaces;
          // console.log("workspaces2")
          // console.log(workspaces)
        resolve(workspaces);
      });
    } catch (err) {
      console.log(err)
      
    }
  return workspaces;
}


export async function postWorkspace(workspaceName:string) {
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
    } catch (err) {
        console.log(err)
    }
    return
}
