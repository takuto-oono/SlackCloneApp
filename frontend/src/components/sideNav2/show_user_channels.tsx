import React, { useRef, useState } from "react";
import { MenuItem, SubMenu } from "react-pro-sidebar";
import Button from "@mui/material/Button";
import Popover from "@mui/material/Popover";
import CreateChannelForm from "@src/components/popUp/create_channel";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { userChannelsState, usersInCState } from "@src/utils/atom";
import { useRouter } from "next/router";
import { getUsersInChannel } from "@src/fetchAPI/channel";
import { UserInWorkspace } from "@src/fetchAPI/workspace";

function ShowUserChannels() {
  const [open, setOpen] = useState(false);
  const divRef = useRef(null);
  const userChannels = useRecoilValue(userChannelsState);
  const router = useRouter()
  const setUsersInC = useSetRecoilState(usersInCState);

  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const getChannelInfo = (channelId: number) =>{
    getUsersInChannel(channelId).then(
      (usersInC: UserInWorkspace[]) => {
        setUsersInC(usersInC);
      });
    router.push({
      pathname: `/main`,
      query: { channelId: channelId },
    })
  }

  const list = userChannels.map((channel, index) => (
    <div key={index}>
      <MenuItem className="bg-purple-200 text-pink-700">
        <button type="button" onClick={() => getChannelInfo(channel.id)} className="inline-block align-baseline text-sm text-pink-700" >
          <div className="truncate">
            {channel.name}
          </div>
        </button>
      </MenuItem>
    </div>
  ));

  return (
    <div>
      <SubMenu label="Channels" className="truncate w-36">
        {list}
        <div ref={divRef}  className="bg-purple-200 text-pink-700">
          <Button onClick={handleClickOpen}>
            <p className="bg-purple-200 text-pink-700">+ チャンネルを追加</p>
          </Button>
          <Popover
              open={open}
              anchorEl={divRef.current}
              onClose={handleClose}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'left',
              }}
            >
            < CreateChannelForm />
            <Button
              onClick={() => router.push({
                pathname: `/main`,
                query: { contents: "show_workspace_channels" },
              })}>
              <p className="text-black">チャンネル一覧</p></Button>
          </Popover>
        </div>
      </SubMenu>
    </div>
  )
}

export default ShowUserChannels;

