import { getUserId } from '../utils/cookie'
import { createUrl, deleteFetcher, getFetcher, patchFetcher, postFetcher } from './common'

export interface Workspace {
  id: number
  name: string
  primary_owner_id: number
}

export interface UserInWorkspace {
  id: number
  name: string
  roleId: number
}

export async function getWorkspaces(): Promise<Workspace[]> {
  const res = await getFetcher(createUrl('/workspace/get_by_user', []))
  const workspaces: Workspace[] = []
  for (const r of res) {
    workspaces.push({
      id: r.id,
      name: r.name,
      primary_owner_id: r.primary_owner_id,
    })
  }
  return workspaces
}

export async function getUsersInWorkspace(workspaceId: number): Promise<UserInWorkspace[]> {
  const res = await getFetcher(createUrl('/workspace/get_users', [workspaceId]))
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

export async function postWorkspace(workspaceName: string): Promise<Workspace> {
  const userID = getUserId()
  if (!userID) {
    throw new Error('not found userID')
  }
  const res = await postFetcher(
    createUrl('/workspace/create', []),
    new Map<string, string | number>([
      ['name', workspaceName],
      ['user_id', userID],
    ]),
  )
  return {
    id: res.id,
    name: res.name,
    primary_owner_id: res.primary_owner_id,
  }
}

export async function addUserInWorkspace(
  workspaceID: number,
  userID: number,
): Promise<UserInWorkspace> {
  const res = await postFetcher(
    createUrl('/workspace/add_user', []),
    new Map<string, string | number>([
      ['workspace_id', workspaceID],
      ['user_id', userID],
      ['role_id', 4],
    ]),
  )
  return {
    id: res.id,
    name: res.name,
    roleId: res.role_id,
  }
}

export const renameWorkspaceName = async (
  workspaceID: number,
  newWorkspaceName: string,
): Promise<Workspace> => {
  const res = await patchFetcher(
    createUrl('/workspace/rename', [workspaceID]),
    new Map<string, string>([['workspace_name', newWorkspaceName]]),
  )
  return {
    id: res.id,
    name: res.name,
    primary_owner_id: res.primary_owner_id,
  }
}

export const deleteUserFromWorkspace = async (
  deletingUserID: number,
  workspaceID: number,
): Promise<UserInWorkspace> => {
  const res = await deleteFetcher(
    createUrl('/workspace/delete_user', []),
    new Map<string, number>([
      ['workspace_id', workspaceID],
      ['user_id', deletingUserID],
    ]),
  )
  return {
    id: res.id,
    name: res.name,
    roleId: res.role_id,
  }
}
