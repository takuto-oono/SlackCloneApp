import React from "react";
import { UserInWorkspace, getUsersInWorkspace } from '@fetchAPI/workspace'
import { MenuItem } from "react-pro-sidebar";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { joinedChannelsState, usersInWState, workspaceIdState, workspacesState } from "@src/utils/atom";
import { Channel, getJoinedChannelsInW } from "@src/fetchAPI/channel";
import { useRouter } from "next/router";

function ShowWorkspaces() {
  const router = useRouter();
  const setWorkspaceId = useSetRecoilState(workspaceIdState);
  const setUsersInW = useSetRecoilState(usersInWState);
  const setJoinedChannels = useSetRecoilState(joinedChannelsState);
  const workspaces = useRecoilValue(workspacesState);

  const getWorkspaceInfo = (workspaceId: number) => {
    setWorkspaceId(workspaceId);
    getUsersInWorkspace(workspaceId).then(
    (usersInW: UserInWorkspace[]) => {
      setUsersInW(usersInW);
      }
    );
    getJoinedChannelsInW(workspaceId).then(
    (joinedChannels: Channel[]) => {
      setJoinedChannels(joinedChannels);
      }
    );
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
