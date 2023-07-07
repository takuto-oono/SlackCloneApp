import React, { useRef, useState } from "react";
import { MenuItem, SubMenu } from "react-pro-sidebar";
import Button from "@mui/material/Button";
import Popover from "@mui/material/Popover";
import CreateChannelForm from "@src/components/popUp/create_channel";
import { useRecoilValue, useSetRecoilState } from "recoil";
import { channelsState, usersInCState, usersInWState } from "@src/utils/atom";
import { useRouter } from "next/router";
import { UserInChannel, getUsersInChannel } from "@src/fetchAPI/channel";

function ShowChannels() {
  const [open, setOpen] = useState(false);
  const divRef = useRef(null);
  const channels = useRecoilValue(channelsState);
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
      (usersInC: UserInChannel[]) => {
        setUsersInC(usersInC);
      });
    router.push(`/main/${channelId}`)
  }

  const list = channels.map((channel, index) => (
    <div key={index}>
      <MenuItem className="bg-purple-200 text-pink-700">
        <button type="button" onClick={() => getChannelInfo(channel.id)} className="inline-block align-baseline text-sm text-pink-700" >
         <>{channel.name}</>
        </button>
      </MenuItem>
    </div>
  ));

  return (
    <div>
      <SubMenu label="Channels">
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
            <Button><p className="text-black">チャンネル一覧</p></Button>
            {/* チャンネル一覧ページへの移動ボタンを設置する（未） */}
          </Popover>
        </div>
      </SubMenu>
    </div>
  )
}


export default ShowChannels;

