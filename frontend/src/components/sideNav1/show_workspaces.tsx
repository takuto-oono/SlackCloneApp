import React from "react";
import { UserInWorkspace, getUsersInWorkspace} from '@fetchAPI/workspace'
import { Link } from 'react-router-dom';
import { MenuItem } from "react-pro-sidebar";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { channelsState, usersInWState, workspacesState } from "@src/utils/atom";
import { Channel, getChannelsByWorkspaceId } from "@src/fetchAPI/channel";

function ShowWorkspaces() {
  const setUsersInW = useSetRecoilState(usersInWState);
  const setChannels = useSetRecoilState(channelsState);
  const workspaces = useRecoilValue(workspacesState);
  const getWorkspaceInfo = (workspaceId: number) =>{
    getUsersInWorkspace(workspaceId).then(
      (usersInW: UserInWorkspace[]) => {
        setUsersInW(usersInW);
      });
    getChannelsByWorkspaceId(workspaceId).then(
      (channels: Channel[]) => {
        setChannels(channels);
      });
  }
  const list = workspaces.map((workspace, index) => (
    <div key={index}>
      <MenuItem>
        <Link to={`${workspace.id}`} onClick={() => getWorkspaceInfo(workspace.id)}>
          <span>{workspace.name}</span>
        </Link>
      </ MenuItem>
    </div>
  ));

  return (
    <div className="App">
      <div>
        {list}
      </div>      
    </div>
  );
}

export default ShowWorkspaces;
