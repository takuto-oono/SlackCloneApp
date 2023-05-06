import React, { useEffect, useState } from "react";
import { Link } from 'react-router-dom';
import { MenuItem, SubMenu } from "react-pro-sidebar";
import { getChannelsByWorkspaceId, Channel } from '@fetchAPI/channel';
import { useParams } from "react-router-dom";
import { UserInfo, getUsers } from "@fetchAPI/workspace";

function ChannelIndex() {
  const [channelList, setChannelList] = useState<Channel[]>([]);
  const [userList, setUserList] = useState<UserInfo[]>([]);
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

  const usersList = userList.map((item , index) => (
    <div key={index}>
      <MenuItem>
        <span>{item.name}</span>
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
  },[workspaceId]);

  useEffect(() => {
    getUsers(parseInt(workspaceId)).then((userList: UserInfo[]) => {
      if (Array.isArray(userList)) {
        setUserList(userList);
        console.log(userList)
      }
    });
  },[workspaceId])

  return (
    <div>
      <SubMenu label="Channels">
        {list}
      </SubMenu>

      {/* テスト用 */}
      <SubMenu label="UsersInWorkspace">
        {usersList}
      </SubMenu>
    </div>
  )
}

export default ChannelIndex;
