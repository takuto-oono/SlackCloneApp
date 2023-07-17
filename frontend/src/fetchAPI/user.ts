// TODO
// user interfaceを改善
// signup関数の改善
// login関数の改善

import router from 'next/router'

export interface User {
  name: string
  password: string
}

const baseUrl = 'http://localhost:8080/api/user/'

export async function signUp(user: User) {
  const url = baseUrl + 'signUp'
  try {
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: user.name,
        password: user.password,
      }),
    })
    if (res.status == 200) {
      console.log('redirect')
      router.replace('/')
    }
  } catch (err) {
    console.log(err)
  }
}

export interface currentUser {
  token: string
  user_id: string
  username: string
}

export async function login(user: User): Promise<currentUser> {
  const url = baseUrl + 'login'
  const currentUser = {
    token: '',
    user_id: '',
    username: '',
  }
  try {
    const res = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: user.name,
        password: user.password,
      }),
    })
    console.log(res)
    const User = await res.json()
    if (res.status == 200) {
      console.log('redirect')
      return new Promise((resolve) => {
        const currentUser: currentUser = {
          token: User.token,
          user_id: User.user_id,
          username: User.username,
        }
        resolve(currentUser)
      })
    } else {
      return new Promise((resolve) => {
        resolve(currentUser)
      })
    }
  } catch (err) {
    console.log(err)
  }

  return currentUser
}
