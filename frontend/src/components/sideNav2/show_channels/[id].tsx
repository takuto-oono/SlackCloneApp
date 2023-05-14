import React, { useEffect, useRef, useState } from "react";
import { Link } from 'react-router-dom';
import { MenuItem, SubMenu } from "react-pro-sidebar";
import { getChannelsByWorkspaceId, Channel } from '@fetchAPI/channel';
import { useParams } from "react-router-dom";

import Button from "@mui/material/Button";
import Popover from "@mui/material/Popover";
import CreateChannelForm from "@src/components/popUp/create_channel";

function ChannelIndex() {
  const [open, setOpen] = useState(false);
  const [channelList, setChannelList] = useState<Channel[]>([]);
  const { workspaceId } = useParams<{ workspaceId: string }>();
  const divRef = useRef(null);

  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };

  const list = channelList.map((item, index) => (
    <div key={index}>
      <MenuItem>
        <Link to="tmp_main">
          <span>{item.name}</span>
        </Link>

      </MenuItem>
    </div>
  ));

  useEffect(() => {
    getChannelsByWorkspaceId(Number(workspaceId)).then((channels: Channel[]) => {
      setChannelList(channels)
    });
  },[workspaceId])

  return (
    <div>
      <SubMenu label="Channels">
        {list}
        <div ref={divRef} style={{ backgroundColor: "rgb(63, 14, 64)" }}>
          <Button onClick={handleClickOpen}>
          <p style={{ color: 'rgb(153, 133, 156)'}}>+ チャンネルを追加</p>
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
            <Button><p style={{ color: 'black'}}>チャンネル一覧</p></Button>
            {/* チャンネル一覧ページへの移動ボタンを設置する（未） */}
          </Popover>
        </div>
      </SubMenu>
    </div>
  )
}

export default ChannelIndex;
