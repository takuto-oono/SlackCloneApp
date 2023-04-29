import React, { useEffect, useState } from "react";
import { Link } from 'react-router-dom';
import { MenuItem, SubMenu } from "react-pro-sidebar";
import { getChannelsByWorkspaceId, Channel } from '@fetchAPI/channel';
import { useParams } from "react-router-dom";


function ChannelIndex() {
  const [channelList, setChannelList] = useState<Channel[]>([]);
  const { workspaceId } = useParams<{ workspaceId: string }>();

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
    console.log(workspaceId);
    getChannelsByWorkspaceId(parseInt(workspaceId)).then((channels: Channel[]) => {
    if (Array.isArray(channels)) {
      setChannelList(channels)
    }
      console.log(channelList)
    });
  },[]);

    

  return (
    <div>
      <SubMenu label="Channel Index">
        {list}
      </SubMenu>
      
    </div>
  )
}

export default ChannelIndex;
