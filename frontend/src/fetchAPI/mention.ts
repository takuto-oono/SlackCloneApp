import { createUrl, getFetcher } from './common'
import { Message } from './message'

export const getMessagesMentionedByUser = async (workspaceID: number) => {
  const res = await getFetcher(createUrl('/mention/by_user', [workspaceID]))
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
