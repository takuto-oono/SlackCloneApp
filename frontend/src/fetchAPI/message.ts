import { createUrl, deleteFetcher, getFetcher, patchFetcher, postFetcher } from './common'

export interface SendDMForm {
  receivedUserID: number
  workspaceID: number
  text: string
}

export interface ResSendDM {
  id: number
  text: string
  sendUserId: number
  dmLineId: number
  createdAt: string
  updatedAt: string
}

export interface Message {
  id: number
  text: string
  ChannelID: number
  dmLineID: number
  userID: number
  threadID: number
  createdAt: string
  updatedAt: string
}

export interface MessageAndUser {
  messageID: number
  userID: number
  isRead: boolean
}

// channel message API
export async function getMessagesFromChannel(channelID: number): Promise<Message[]> {
  const res = await getFetcher(createUrl('/message/get_from_channel', [channelID]))
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

export const sendMessage = async (
  text: string,
  channelID: number,
  mentionedUserIDs: number[],
): Promise<Message> => {
  const res = await postFetcher(
    createUrl('/message/send', []),
    new Map<string, number | string | number[]>([
      ['text', text],
      ['channel_id', channelID],
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

export const readMessageByUser = async (messageID: number): Promise<MessageAndUser> => {
  const res = await postFetcher(createUrl('/message/read_by_user', [messageID]))
  return {
    messageID: res.message_id,
    userID: res.user_id,
    isRead: res.is_read,
  }
}

export const editMessage = async (messageID: number, newText: string): Promise<Message> => {
  const res = await patchFetcher(
    createUrl('/message/edit', [messageID]),
    new Map<string, string>([['text', newText]]),
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
// direct message API
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
