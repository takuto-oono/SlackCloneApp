import React, { useRef, useState } from "react";
import { MenuItem, SubMenu } from "react-pro-sidebar";
import Button from "@mui/material/Button";
import Popover from "@mui/material/Popover";
import CreateChannelForm from "@src/components/popUp/create_channel";
import { useRecoilValue } from "recoil";
import { channelsState } from "@src/utils/atom";
import { useRouter } from "next/router";

function ShowChannels() {
  const [open, setOpen] = useState(false);
  const divRef = useRef(null);
  const channels = useRecoilValue(channelsState);
  const router = useRouter()


  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const list = channels.map((item, index) => (
    <div key={index}>
      <MenuItem className="bg-purple-200 text-pink-700">
        <button type="button" onClick={() => router.push(`/main/${item.id}`)} className="inline-block align-baseline text-sm text-pink-700" >
         <>{item.name}</>
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

