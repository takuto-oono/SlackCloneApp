import { createUrl, getFetcher, postFetcher } from './common'

export interface User {
  id: number
  name: string
}

export interface CurrentUser {
  token: string
  user_id: string
  username: string
}

export async function signUp(userName: string, password: string) {
  await postFetcher(
    createUrl('/user/signUp', new Array(0)),
    new Map<string, string | number>([
      ['name', userName],
      ['password', password],
    ]),
  )
}

export async function login(userName: string, password: string): Promise<CurrentUser> {
  const res = await postFetcher(
    createUrl('/user/login', new Array(0)),
    new Map<string, string | number>([
      ['name', userName],
      ['password', password],
    ]),
  )
  const user: CurrentUser = {
    token: res.token,
    user_id: res.user_id,
    username: res.username,
  }
  return user
}

export async function getAllUsers(): Promise<User[]> {
  const res = await getFetcher(createUrl('/user/all', []))
  let users: User[] = []
  console.log(res)
  for (const r of res) {
    users.push({
      id: r.id,
      name: r.name,
    })
  }
  return users
}
