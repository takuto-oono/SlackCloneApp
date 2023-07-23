import React from "react";
import { Menu, MenuItem } from "react-pro-sidebar";
import { useRouter } from "next/router";
import { Channel, getChannelsInW } from "@src/fetchAPI/channel";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { channelsState, workspaceIdState } from "@src/utils/atom";

function ShowContentsList() {
  return (
    <Menu className="pd-5 bg-purple-200 text-pink-700">
      < ShowChannelsInWorkspace />
    </Menu>
  )
}

export default ShowContentsList;

const ShowChannelsInWorkspace = () => {
  const router = useRouter();
  const workspaceId = useRecoilValue(workspaceIdState);
  const setChannels = useSetRecoilState(channelsState);
  const handleClick = () => {
    getChannelsInW(workspaceId).then(
    (channels: Channel[]) => {
      setChannels(channels);
      }
    );
    router.push({
      pathname: `/main`,
      query: { contents: "show_channels" },
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
