import SendDM from "@src/components/main/send_dm";
import { getToken } from "./cookie";
import { resetCookie, getUserId } from "./cookie";

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

const baseUrl = "http://localhost:8080/api/dm/";

export async function sendDM(form: SendDMForm): Promise<ResSendDM|void> {
  const url = baseUrl + "send";
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
