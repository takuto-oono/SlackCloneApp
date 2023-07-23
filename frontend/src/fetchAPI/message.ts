import { createUrl, getFetcher, patchFetcher, postFetcher } from './common'

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
