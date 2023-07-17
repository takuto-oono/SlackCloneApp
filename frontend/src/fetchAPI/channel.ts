import { getToken } from "../utils/cookie";
import { UserInWorkspace } from "./workspace";

export interface Channel {
  id: number;
  name: string;
  description: string;
  is_private: boolean;
  is_archive: boolean;
  workspace_id: number;
}

export interface CurrentChannel {
  name: string;
  description: string;
  is_private: boolean;
  workspace_id: number;
}

const baseUrl = "http://localhost:8080/api/channel/";

export async function getJoinedChannelsInW(
  workspace_id: number
): Promise<Channel[]> {
  const url = baseUrl + "get_by_user_and_workspace/" + workspace_id;
  console.log(url);
  let res_channels = [
    {
      id: 0,
      name: "",
      description: "",
      is_private: false,
      is_archive: false,
      workspace_id: 0,
    },
  ];
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    console.log(res);
    res_channels = await res.json();
    return new Promise((resolve) => {
      resolve(res_channels);
    });
  } catch (err) {
    console.log("err");
    console.log(err);
  }
  return res_channels;
}

export async function getChannelsInW(
  workspace_id: number
): Promise<Channel[]> {
  const url = baseUrl + workspace_id;
  console.log(url);
  let res_channels = [
    {
      id: 0,
      name: "",
      description: "",
      is_private: false,
      is_archive: false,
      workspace_id: 0,
    },
  ];
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    console.log(res);
    res_channels = await res.json();
    return new Promise((resolve) => {
      resolve(res_channels);
    });
  } catch (err) {
    console.log("err");
    console.log(err);
  }
  return res_channels;
}

export async function getUsersInChannel(channelId: number): Promise<UserInWorkspace[]> {
  const url = baseUrl + "all_user/" + channelId;
  let resUsersInChannel: UserInWorkspace[];
  const usersInChannel = [{
    id: 0,
    name: "",
    roleid: 0,
  }];

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    
    if (res.status == 200) {
      resUsersInChannel = await res.json();
      return new Promise((resolve) => {
        const usersInChannel: UserInWorkspace[] = resUsersInChannel;
        resolve(usersInChannel);
      });
    }
    
  } catch (err) {
    console.log(err);
  }
  return usersInChannel;
}

export async function postChannel(current: CurrentChannel): Promise<number | undefined> {
  const url = baseUrl + "create";
  let channel: Channel;
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        name: current.name,
        description: current.description,
        is_private: current.is_private,
        workspace_id: current.workspace_id,
      }),
    });
    if (res.status == 200) {
      channel = await res.json();
      return channel.id;
    }
  } catch (err) {
    console.log(err);
  }
  return;
}

export async function addUserInChannel(chanelID: number, userID: number) {
  const url = baseUrl + "add_user";
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        channel_id: chanelID,
        user_id: userID,
      }),
    })
    if (res.status != 200) {
      console.log(res);
    }
  } catch (e) {
    console.log(e);
  }
}

// TODO deleteuserfromchannel
// TODO deletechannel