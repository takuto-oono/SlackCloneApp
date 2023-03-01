export interface Workspace {
    id: number;
    name: string;
    primary_owner_id: number;
}

export interface workspaceList {
    workspaces?: Workspace[]
}


//Cookieに保存されているjwtTokenを取り出す
function getToken() {
  if (typeof document !== 'undefined') {
    const jwtToken = document.cookie.split( '; ' )[ 0 ].split( '=' )[ 1 ];
    return jwtToken
  }
  else {
    const jwtToken = "";
    return jwtToken
  }
}

const baseUrl = 'http://localhost:8080/api/workspace/'

export async function getWorkspaces(): Promise<workspaceList> {
    const url = baseUrl + 'get_by_user'
  let res_workspaces: Workspace[]
  const workspace = {
    id: 0,
    name: "",
    primary_owner_id: 0
  }
  const workspace_list = {
    workspaces: [workspace]
  }
    try {
        const res = await fetch(url, {
          method: 'GET',
          headers: {
              'Authorization': getToken(),
          },
        })
        console.log(res)
        res_workspaces = await res.json()
        console.log("workspaces1")
        console.log(res_workspaces);
      
        return new Promise((resolve) => {
        const workspace_list: workspaceList = {
          workspaces: res_workspaces
        };
          console.log("workspaces2")
          console.log(workspace_list)
        resolve(workspace_list);
      });
    } catch (err) {
        console.log(err)
    }
  return workspace_list;
}
