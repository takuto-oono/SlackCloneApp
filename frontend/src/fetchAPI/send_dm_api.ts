import { getToken } from "./cookie";
import { resetCookie, getUserId } from "./cookie";

export interface SendDM {
  receivedUserId: number;
  workspaceId: number;
  text: string;
}

export interface SendDMInfo {
  id: number;
  text: string;
  sendUserId: number;
  dmLineId: number;
  createdAt: string;
  updatedAt: string;
}

const baseUrl = "http://localhost:8080/api/dm/";

export async function sendDM(current: SendDM) {
  const url = baseUrl + "send";
  let resSendDM: SendDMInfo;
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        receivedUserId: current.receivedUserId,
        workspaceId: current.workspaceId,
        text: current.text,
      }),
    });
    resSendDM = await res.json();
  } catch (err) {
    console.log(err);
  }
  return;
}
