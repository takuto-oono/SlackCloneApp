import React from "react";
import { UserInWorkspace, getUsersInWorkspace } from '@fetchAPI/workspace'
import { MenuItem } from "react-pro-sidebar";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { channelsState, usersInWState, workspaceIdState, workspacesState } from "@src/utils/atom";
import { Channel, getChannelsByWorkspaceId } from "@src/fetchAPI/channel";
import { AddUserInWorkspaceForm } from "../popUp/add_user_in_workspace_form";
import { useRouter } from "next/router";

function ShowWorkspaces() {
  const router = useRouter();
  const setWorkspaceId = useSetRecoilState(workspaceIdState);
  const setUsersInW = useSetRecoilState(usersInWState);
  const setChannels = useSetRecoilState(channelsState);
  const workspaces = useRecoilValue(workspacesState);

  const getWorkspaceInfo = (workspaceId: number) => {
    setWorkspaceId(workspaceId);
    getUsersInWorkspace(workspaceId).then(
    (usersInW: UserInWorkspace[]) => {
      setUsersInW(usersInW);
    });
    getChannelsByWorkspaceId(workspaceId).then(
    (channels: Channel[]) => {
      setChannels(channels);
    });
    router.push({
      query: { workspaceId: workspaceId },
    })
  }

  const list = workspaces.map((workspace, index) => (
    <div key={index}>
      <MenuItem>
        <button type="button" onClick={() => getWorkspaceInfo(workspace.id)} className="inline-block align-baseline text-sm text-pink-700" >
          <div className="truncate">
            {workspace.name}
          </div>
        </button>
			</MenuItem>
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
