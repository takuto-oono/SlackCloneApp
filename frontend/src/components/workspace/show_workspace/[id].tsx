import React, { useEffect, useState } from "react";
import { getChannelsByWorkspaceId, Channel } from 'src/fetchAPI/channel';
import { getToken } from "@src/fetchAPI/cookie";
import { useParams } from "react-router-dom";


function ShowWorkspace() {
  const [channelList, setChannelList] = useState<Channel[]>([]);
  const { id } = useParams<{ id: string }>();

  const list = channelList.map((item, index) => (
    <div key={index}>
      <p>channel_id:{item.id}</p>
      <p>name:{item.name}</p>
      <p>description:{item.description}</p>
      <p>is_private:{String(item.is_private)}</p>
      <p>is_archive:{String(item.is_archive)}</p>
      <p>---</p>
    </div>
  ));

  useEffect(() => {
    getChannelsByWorkspaceId(parseInt(id)).then((channels: Channel[]) => {
    if (Array.isArray(channels)) {
      setChannelList(channels)
    }
      console.log(channelList)
    });
  },[]);

    

  return (
    <div>
      <h2>Channel Index</h2>
      <p>workspace_id:{id}</p>
      <p>---</p>
      {list}
    </div>
  )
}

export default ShowWorkspace;
