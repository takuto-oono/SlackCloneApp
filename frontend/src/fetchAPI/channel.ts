import { createUrl, getFetcher, postFetcher } from './common'
import { UserInWorkspace } from './workspace'

// TODO 命名規則を守るように変更
export interface Channel {
  id: number
  name: string
  description: string
  is_private: boolean
  is_archive: boolean
  workspace_id: number
}

// TODO 命名規則を守るように変更
export interface CurrentChannel {
  name: string
  description: string
  is_private: boolean
  workspace_id: number
}

export interface UserInChannel {
  userID: number
  channelID: number
}

export const getJoinedChannelsInW = async (workspaceID: number): Promise<Channel[]> => {
  const res = await getFetcher(createUrl('/channel/get_by_user_and_workspace', [workspaceID]))
  let channels: Channel[] = []
  for (const r of res) {
    channels.push({
      id: r.id,
      name: r.name,
      description: r.description,
      is_private: r.is_private,
      is_archive: r.is_archive,
      workspace_id: r.workspace_id,
    })
  }
  return channels
}

export const getChannelsInW = async (workspaceID: number): Promise<Channel[]> => {
  const res = await getFetcher(createUrl('/channel', [workspaceID]))
  let channels: Channel[] = []
  for (const r of res) {
    channels.push({
      id: r.id,
      name: r.name,
      description: r.description,
      is_private: r.is_private,
      is_archive: r.is_archive,
      workspace_id: r.workspace_id,
    })
  }
  return channels
}

export const getUsersInChannel = async (channelID: number): Promise<UserInWorkspace[]> => {
  const res = await getFetcher(createUrl('/channel/all_user', [channelID]))
  let usersInWorkspace: UserInWorkspace[] = []
  for (const r of res) {
    usersInWorkspace.push({
      id: r.id,
      name: r.name,
      roleId: r.role_id,
    })
  }
  return usersInWorkspace
}

export const postChannel = async (
  channelName: string,
  description: string,
  isPrivate: boolean,
  workspaceID: number,
): Promise<Channel> => {
  const res = await postFetcher(
    createUrl('/channel/create', []),
    new Map<string, number | string | boolean>([
      ['name', channelName],
      ['description', description],
      ['is_private', isPrivate],
      ['workspace_id', workspaceID],
    ]),
  )
  return {
    id: res.id,
    name: res.name,
    description: res.description,
    is_private: res.is_private,
    is_archive: res.is_archive,
    workspace_id: res.workspace_id,
  }
}

export const addUserInChannel = async (
  channelID: number,
  userID: number,
): Promise<UserInChannel> => {
  const res = await postFetcher(
    createUrl('/channel/add_user', []),
    new Map<string, number>([
      ['channel_id', channelID],
      ['user_id', userID],
    ]),
  )
  return {
    channelID: res.channel_id,
    userID: res.user_id,
  }
}

// TODO deleteuserfromchannel
// backendの仕様変更待ち
// issueのurl: https://github.com/TO053037/SlackCloneApp/issues/238
export const deleteUserFromChannel = async () => {}

// TODO deletechannel
// backendの仕様変更待ち
// issueのurl: https://github.com/TO053037/SlackCloneApp/issues/239
export const deleteChannel = async () => {}
