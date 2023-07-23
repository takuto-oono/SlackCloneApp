import { atom, useRecoilValue, useResetRecoilState } from 'recoil'
import { UserInWorkspace } from '@src/fetchAPI/workspace'
import { usersInWState } from './atom'

interface SearchMemo {
  User: UserInWorkspace
  Point: number
}

export class UserSearch {
  memos: SearchMemo[] = []
  input: string

  constructor(input: string) {
    this.input = input.toLowerCase()
    for (const userInW of useRecoilValue(usersInWState)) {
      this.memos.push({
        User: userInW,
        Point: 0,
      })
    }
  }

  Do(): UserInWorkspace[] {
    if (this.input.length == 2) {
      this.searchInitial()
    }
    this.searchContains()
    console.log(this.memos)
    return this.createOutput()
  }

  createOutput(): UserInWorkspace[] {
    const result: UserInWorkspace[] = []
    for (let i = 0; i < this.memos.length; i++) {
      for (let j = i + 1; j < this.memos.length; j++) {
        if (this.memos[i] < this.memos[j]) {
          this.memos[i], (this.memos[j] = this.memos[j]), this.memos[i]
        }
      }
    }
    for (const memo of this.memos) {
      if (memo.Point > 0) {
        result.push(memo.User)
      }
    }
    return result
  }

  searchInitial(): void {
    if (this.input.length != 2) {
      return
    }
    for (const memo of this.memos) {
      const words: string[] = memo.User.name.split(' ')
      if (words.length != 2) {
        continue
      }
      console.log(this.input[0], this.input[1])
      if (
        (words[0][0] == this.input[0] && words[1][0] == this.input[1]) ||
        (words[0][0] == this.input[1] && words[1][0] == this.input[0])
      ) {
        memo.Point += 5
      }
    }
  }

  searchContains(): void {
    for (const memo of this.memos) {
      for (const word of memo.User.name.split(' ')) {
        if (this.input.length > word.length) {
          continue
        }
        if (word.includes(this.input)) {
          memo.Point += 10
          break
        }
      }
    }
  }
}
