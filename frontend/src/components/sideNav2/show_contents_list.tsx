import React from "react";
import { Menu, MenuItem } from "react-pro-sidebar";
import { useRouter } from "next/router";
import { Channel, getChannelsInW } from "@src/fetchAPI/channel";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { workspaceChannelsState, workspaceIdState } from "@src/utils/atom";

function ShowContentsList() {
  return (
    <Menu className="pd-5 bg-purple-200 text-pink-700">
      < ShowWorkspaceChannels />
    </Menu>
  )
}

export default ShowContentsList;

const ShowWorkspaceChannels = () => {
  const router = useRouter();
  const workspaceId = useRecoilValue(workspaceIdState);
  const setWorkspaceChannels = useSetRecoilState(workspaceChannelsState);
  const handleClick = () => {
    getChannelsInW(workspaceId).then(
    (workspaceChannels: Channel[]) => {
      setWorkspaceChannels(workspaceChannels);
      }
    );
    router.push({
      pathname: `/main`,
      query: { contents: "show_workspace_channels" },
    })
  }

  return (
    <MenuItem>
      <button
        onClick={() => handleClick()}>
        <p>すべてのチャンネル</p></button>
    </MenuItem>
  )
}
