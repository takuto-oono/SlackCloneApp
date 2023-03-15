import router from "next/router";
import { getToken } from "./cookie";

export interface Channel {
    id: number;
    name: string;
    description: string;
    is_private: boolean;
    is_archive: boolean;
    workspace_id: number;
}

const baseUrl = 'http://localhost:8080/api/channel/'

export async function getChannelsByWorkspaceId(workspace_id: number): Promise<Channel[]> {
  const url = baseUrl + 'get_by_user_and_workspace/' + workspace_id
  console.log(url);
  let res_channels: Channel[]
  const channels = [
    {
    id: 0,
    name: "",
    description: "",
    is_private: false,
    is_archive: false,
    workspace_id: 0
    }
  ]
  try {
      const res = await fetch(url, {
        method: 'GET',
        headers: {
              'Authorization': getToken(),
          },
      })
    console.log(res)
    if (res.status == 401) {
        console.log("redirect");
        router.replace('/')
      }
      res_channels = await res.json()
      return new Promise((resolve) => {
        const channels: Channel[] = res_channels;
        resolve(channels);
      });
  } catch (err) {
    console.log("err");
    console.log(err)
  }
  return channels;
}
