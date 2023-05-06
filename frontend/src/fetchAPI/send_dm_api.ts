import { getToken } from "./cookie";
import { resetCookie, getUserId } from "./cookie";

export interface SendDM {
  received_user_id: number;
  workspace_id: number;
  text: string;
}

export interface SendDMInfo {
  id: number;
  text: string;
  send_user_id: number;
  dm_line_id: number;
  created_at: string;
  updated_at: string;
}

const baseUrl = "http://localhost:8080/api/dm/";

export async function sendDM(current: SendDM) {
  const url = baseUrl + "send";
  let res_SendDM: SendDMInfo;
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        Authorization: getToken(),
      },
      body: JSON.stringify({
        received_user_id: getUserId(),
        workspace_id: current.workspace_id,
        text: current.text,
      }),
    });
    res_SendDM = await res.json();
  } catch (err) {
    console.log(err);
  }
  return;
}
