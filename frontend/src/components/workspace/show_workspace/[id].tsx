import React, { useEffect, useState } from "react";
<<<<<<< HEAD:frontend/pages/component/show_workspace/[id].tsx
import { getChannelsByWorkspaceId, Channel } from 'pages/fetchAPI/channel';
import { useRouter } from "next/router";
import CreateChannel from '../create_channel';
import Link from 'next/link'
=======
import { getChannelsByWorkspaceId, Channel } from '@fetchAPI/channel';
import { useParams } from "react-router-dom";
>>>>>>> main:frontend/src/components/workspace/show_workspace/[id].tsx


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
      <CreateChannel workspace_id={Number(router.query.id)} />
    </div>
  )
}

export default ShowWorkspace;
