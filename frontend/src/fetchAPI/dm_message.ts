import { createUrl, deleteFetcher, getFetcher, patchFetcher, postFetcher } from './common'
import { Message } from './message'

export const getDmsInLine = async (dmLineID: number): Promise<Message[]> => {
  const res = await getFetcher(createUrl('/dm', [dmLineID]))
  let messages: Message[] = []
  for (const r of res) {
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

// TODO backendの仕様変更後に実装
// issueのurl: https://github.com/TO053037/SlackCloneApp/issues/241
export const getDmLinesByUserInWorkspace = async (workspaceID: number) => {}

export const sendDM = async (
  text: string,
  workspaceID: number,
  receivedUserID: number,
  mentionedUserIDs: number[],
): Promise<Message> => {
  const res = await postFetcher(
    createUrl('/dm/send', []),
    new Map<string, number | string | number[]>([
      ['text', text],
      ['received_user_id', receivedUserID],
      ['workspace_id', workspaceID],
      ['mentioned_user_ids', mentionedUserIDs],
    ]),
  )
  return {
    id: res.id,
    text: res.text,
    ChannelID: res.channel_id,
    dmLineID: res.dm_line_id,
    userID: res.user_id,
    threadID: res.thread_id,
    createdAt: res.created_at,
    updatedAt: res.updated_at,
  }
}

export const editDM = async (dmID: number, newText: string): Promise<Message> => {
  const res = await patchFetcher(
    createUrl('/dm', [dmID]),
    new Map<string, number | string>([['text', newText]]),
  )
  return {
    id: res.id,
    text: res.text,
    ChannelID: res.channel_id,
    dmLineID: res.dm_line_id,
    userID: res.user_id,
    threadID: res.thread_id,
    createdAt: res.created_at,
    updatedAt: res.updated_at,
  }
}

export const deleteDM = async (dmID: number): Promise<Message> => {
  const res = await deleteFetcher(createUrl('/dm', [dmID]))
  return {
    id: res.id,
    text: res.text,
    ChannelID: res.channel_id,
    dmLineID: res.dm_line_id,
    userID: res.user_id,
    threadID: res.thread_id,
    createdAt: res.created_at,
    updatedAt: res.updated_at,
  }
}
