import { getToken } from "../utils/cookie";

export interface SendDMForm {
  receivedUserID: number;
  workspaceID: number;
  text: string;
}

export interface ResSendDM {
  id: number;
  text: string;
  sendUserId: number;
  dmLineId: number;
  createdAt: string;
  updatedAt: string;
}

export interface  Message {
  id: number;
  text: string;
  ChannelID: number;
  dmLineID: number;
  userID: number;
  threadID: number;
  createdAt: string;
  updatedAt: string;
}

const baseUrl = "http://localhost:8080/api";

// channel message API
export async function getMessagesFromChannel(channelID: number): Promise<Message[]|null> {
  const url: string = baseUrl + "/message/get_from_channel/" + channelID.toString();
  try {
    const res: Response = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    });
    if (res.status == 200) {
      const response = await res.json();
      const messages: Message[] = [];
      if (!response) {
        return messages
      }
      for (const r of response) {
        messages.push({
          id: r.id,
          text: r.text,
          ChannelID: r.channel_id,
          dmLineID: r.dm_line_id,
          userID: r.user_id,
          threadID: r.thread_id,
          createdAt: r.created_at,
          updatedAt: r.updated_at,
        })
      }
      return messages
    }
    console.log(res)
  } catch (e) {
    console.log(e)
  }
  return null
}

export async function sendMessage(text: string, channelID: number, mentionedUserIDs: number[]): Promise<Message|null> {
  const url: string = baseUrl + "/message/send";
  try {
    const res: Response = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        text: text,
        channel_id: channelID,
        mentioned_user_ids: mentionedUserIDs
      })
    });
    if (res.status == 200) {
      const response = await res.json();
      const message: Message = {
        id: response.id,
        text: response.text,
        ChannelID: response.channel_id,
        dmLineID: response.dm_line_id,
        userID: response.user_id,
        threadID: response.thread_id,
        createdAt: response.created_at,
        updatedAt: response.updated_at,
      }
      return message;
    }
    console.log(res);
  } catch (e) {
    console.log(e);
  }
  return null
}

// TODO readmessagebyuser

// TODO editmessage

// direct message API
// TODO getdmsinline

// TODO getdmlines

export async function sendDM(form: SendDMForm): Promise<ResSendDM|void> {
  const url = baseUrl + "dm/send";
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        received_user_id: form.receivedUserID,
        workspace_id: form.workspaceID,
        text: form.text,
      }),
    });
    return await res.json();
  } catch (err) {
    console.log(err);
  }
}

// TODO editdm

// TODO deletedm