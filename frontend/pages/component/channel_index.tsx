import React, { useState } from "react";
import { getChannelsByWorkspaceId, Channel } from 'pages/fetchAPI/channel'
import { getToken } from "@/pages/fetchAPI/cookie";

function ChannelIndex() {
  const [channelList, setChannelList] = useState<Channel[]>([]);
  const workspace_id = 1;
  const list = channelList.map((item, index) => (
    <div key={index}>
      <p>id:{item.id}</p>
      <p>name:{item.name}</p>
      <p>description:{item.description}</p>
      <p>is_private:{item.is_private}</p>
      <p>is_archive:{item.is_archive}</p>
      <p>workspace_id:{item.workspace_id}</p>
      <p>---</p>
    </div>
  ));
  const handleGetChannels = () => {
    console.log(getToken());
    if (typeof getToken() !== 'undefined') {
      getChannelsByWorkspaceId(workspace_id).then((channels: Channel[]) => {
        console.log("channels")
        console.log(channels)
        setChannelList(channels)
        console.log(channelList)
      });
    }
  }

  return (
    <div className="App">
      <button onClick={handleGetChannels}>チャンネル一覧を表示</button>
      <div>
        {list}
      </div>
    </div>
  );
}

export default ChannelIndex;
