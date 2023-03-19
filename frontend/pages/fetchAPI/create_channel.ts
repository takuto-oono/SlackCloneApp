import { Channel, getChannelsByWorkspaceId } from "./channel";
import { getToken,getUserId } from "./cookie";

const baseUrl = "http://localhost:8080/api/channel/";

export async function getChannels(): Promise<Channel[]> {
  const url = baseUrl + "get_by_user";
  console.log("getToken()");
  console.log(getToken());
  let res_channels: Channel[];
  const channels = [
    {
			id: 0,
			name: "",
			description: "",
			is_private: false,
			is_archive: false,
			workspace_id: 0
    }
  ];

  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    res_channels = await res.json();

    return new Promise((resolve) => {
      const channels: Channel[] = res_channels;
      resolve(channels);
    });
  } catch (err) {
    console.log(err);
  }
  return channels;
}

export async function postChannel(channelName: string) {
  const url = baseUrl + "create";
  let channel: Channel;
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        name: channelName,
        user_id: getUserId(),
      }),
    });
    console.log(res);
    channel = await res.json();
    console.log(channel);
  } catch (err) {
    console.log(err);
  }
  return;
}
