import React, { useEffect, useState } from "react";
import { getWorkspaces, Workspace, UserInWorkspace, getUsersInWorkspace} from '@fetchAPI/workspace'
import { Link } from 'react-router-dom';
import { MenuItem } from "react-pro-sidebar";
import { atom, useRecoilState } from "recoil";

export const usersInWState = atom<UserInWorkspace[]>({
  key: "usersInW",
  default: []
})

function WorkspaceIndex() {
  const [usersInW, setUsersInW] = useRecoilState(usersInWState);
  const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);
  
  const getWorkspaceInfo = (workspaceId: number) =>{
    getUsersInWorkspace(workspaceId).then(
      (usersInW: UserInWorkspace[]) => {
      setUsersInW(usersInW);
    });
    console.log(usersInW);
  }

  const list = workspaceList.map((workspace, index) => (
    <div key={index}>
      <MenuItem>
        <Link to={`${workspace.id}`} onClick={() => getWorkspaceInfo(workspace.id)}>
          <span>{workspace.name}</span>
        </Link>
      </ MenuItem>
    </div>
  ));

  useEffect(() => {
    getWorkspaces().then((workspaces: Workspace[]) => {
      setWorkspaceList(workspaces);
    });
  },[]);

  return (
    <div className="App">
      <div>
        {list}
      </div>      
    </div>
  );
}

export default WorkspaceIndex;
