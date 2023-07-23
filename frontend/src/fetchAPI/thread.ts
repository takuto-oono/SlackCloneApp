import { createUrl, getFetcher, postFetcher } from './common'
import { Message } from './message'

export interface Thread {
  id: number
  messages: Message[]
}

export const getThreadsByUser = async (workspaceID: number): Promise<Thread[]> => {
  const res = await getFetcher(createUrl('/thread/by_user', [workspaceID]))
  let threads: Thread[] = []
  for (const r of res) {
    let messages: Message[] = []
    for (const m of r.message) {
      messages.push({
        id: m.id,
        text: m.text,
        ChannelID: m.channel_id,
        dmLineID: m.dm_line_id,
        userID: m.user_id,
        threadID: m.thread_id,
        createdAt: m.created_at,
        updatedAt: m.updated_at,
      })
    }
    threads.push({
      id: r.id,
      messages: messages,
    })
  }
  return threads
}

export const postThread = async (text: string, parentMessageID: number): Promise<Message> => {
  const res = await postFetcher(
    createUrl('/thread/post', []),
    new Map<string, number | string>([
      ['text', text],
      ['parent_message_id', parentMessageID],
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
