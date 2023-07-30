import { useRecoilValue } from 'recoil'
import { UserInWorkspace } from '@src/fetchAPI/workspace'
import { usersInWState } from './atom'

interface SearchMemo {
  User: UserInWorkspace
  Point: number
}

export const getUserByUserID = (userID: number): UserInWorkspace => {
  const users: UserInWorkspace[] = useRecoilValue(usersInWState).filter((u) => u.id == userID)
  if (!users) throw new Error('not found user')
  if (users.length > 1) throw new Error('find users')
  return users[0]
}

export const getUserByUserName = (userName: string): UserInWorkspace => {
  const users: UserInWorkspace[] = useRecoilValue(usersInWState).filter((u) => u.name == userName)
  if (!users) throw new Error('not found user')
  if (users.length > 1) throw new Error('find users')
  return users[0]
}

export const getUsersByRoleID = (roleID: number): UserInWorkspace[] => {
  return useRecoilValue(usersInWState).filter((u) => u.roleId == roleID)
}

export const searchUserAlgo = (input: string): UserInWorkspace[] => {
  return new UserSearchAlgo(input).algoMain()
}

export class UserSearchAlgo {
  memos: SearchMemo[] = []
  input: string

  constructor(input: string) {
    this.input = input.toLowerCase()
    this.memos = useRecoilValue(usersInWState).map((user) => {
      return {
        User: user,
        Point: 0,
      }
    })
  }

  algoMain(): UserInWorkspace[] {
    if (this.input.length == 2) {
      this.searchInitial(5)
    }
    this.isSubstring(10)
    return this.createOutput()
  }

  createOutput(): UserInWorkspace[] {
    return this.memos
      .filter((user) => user.Point > 0)
      .sort((x, y) => y.Point - x.Point)
      .map((memo) => {
        return memo.User
      })
  }

  searchInitial(point: number): void {
    if (this.input.length != 2) return
    for (const memo of this.memos) {
      const words: string[] = memo.User.name.split(' ')
      if (words.length != 2) continue
      if (
        (words[0][0] == this.input[0] && words[1][0] == this.input[1]) ||
        (words[0][0] == this.input[1] && words[1][0] == this.input[0])
      ) {
        memo.Point += point
      }
    }
  }

  isSubstring(point: number): void {
    this.memos.forEach((memo, i) => {
      if (memo.User.name.includes(this.input)) this.memos[i].Point += point
    })
  }
}
