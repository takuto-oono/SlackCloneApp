import React from "react";
import { UserInWorkspace } from '@fetchAPI/workspace'
import { useRecoilValue, useSetRecoilState } from "recoil";
import { usersInCState, workspaceChannelsState } from "@src/utils/atom";
import { getUsersInChannel } from "@src/fetchAPI/channel";
import { useRouter } from "next/router";

function ShowWorkspaceChannels() {
  const router = useRouter();
  const setUsersInC = useSetRecoilState(usersInCState);
  const workspaceChannels = useRecoilValue(workspaceChannelsState);

  const getChannelInfo = (channelId: number) => {
    getUsersInChannel(channelId).then(
    (usersInC: UserInWorkspace[]) => {
      setUsersInC(usersInC);
    });
    router.push({
      pathname: `/main`,
      query: { channelId: channelId },
    })
  }

  const list = workspaceChannels.map((channel, index) => (
    <div key={index}>
      <button type="button" onClick={() => getChannelInfo(channel.id)} className="inline-block align-baseline text-sm text-pink-700" >
        <div className="truncate">
          {channel.name}
        </div>
      </button>
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

export default ShowWorkspaceChannels;
