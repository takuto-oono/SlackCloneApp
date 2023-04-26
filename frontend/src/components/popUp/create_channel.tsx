import { getToken } from "@fetchAPI/cookie";

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

export async function getChannelsByWorkspaceId(
  workspace_id: number
): Promise<Channel[]> {
  const url = baseUrl + "get_by_user_and_workspace/" + workspace_id;
  console.log(url);
  // let res_channels: Channel[]
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

export async function postChannel(current: CurrentChannel) {
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
    channel = await res.json();
  } catch (err) {
    console.log(err);
  }
  return;
}
