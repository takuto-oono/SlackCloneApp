import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";
import { postChannel } from "@src/components/popUp/create_channel";
import Link from "next/link";
import { getChannelsByWorkspaceId, Channel } from "@fetchAPI/channel";
import { useParams } from "react-router-dom";
import { UserInfo, getUsers } from "@fetchAPI/workspace";

function ShowWorkspace() {
  const [channelList, setChannelList] = useState<Channel[]>([]);
  const [userList, setUserList] = useState<UserInfo[]>([]);
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

  const list2 = userList.map((item , index) => (
    <div key={index}>
      <p>name:{item.name}</p>
      <p>---</p>
    </div>
  ));

  useEffect(() => {
    getChannelsByWorkspaceId(parseInt(id)).then((channels: Channel[]) => {
      if (Array.isArray(channels)) {
        setChannelList(channels);
      }
      console.log(channelList);
    });
    getUsers(parseInt(id)).then((userList: UserInfo[]) => {
      if (Array.isArray(userList)) {
        setUserList(userList);
      }
    });
  }, []);

  return (
    <div>
      <h2>Channel Index</h2>
      <p>workspace_id:{id}</p>
      <p>---</p>
      {list}
      {/* <CreateChannel workspace_id={Number(router.query.id)} /> */}
      {list2}
    </div>
  );
}

export default ShowWorkspace;
