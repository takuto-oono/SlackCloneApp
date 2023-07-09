import { getToken } from "../utils/cookie";

export interface User {
    id: string;
    name: string;
    password: string;
}

const baseUrl = 'http://localhost:8080/'

export async function getUsers() {
    const url = baseUrl + 'users'
    let users: User[]
    try {
        const res = await fetch(url, {
            method: 'GET',
        })
        console.log(res)
        users = await res.json()
        console.log(users)
        return users
    } catch (err) {
        console.log(err)
    }
}

export async function getUserById(user: User): Promise<User> {
    const url = baseUrl + 'user/' + user.id
    try {
        const res = await fetch(url, {
            method: 'GET',
        })
        console.log(res)
        const user: User = await res.json()
        console.log(user)
    } catch (err) {
        console.log(err)
    }
    return user
}

export async function postUser(user: User): Promise<User> {
    const url = baseUrl + 'user'

    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: user.name,
                password: user.password,
            })
        })
        console.log(res)
        user = await res.json()
        console.log(user)
    } catch (err) {
        console.log(err)
    }
    return user
}

export async function updateUser(user: User) :Promise<User> {
    const url = baseUrl + 'user/' + user.id
    try {
        const res = await fetch(url, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: user.name,
                password: user.password,
            })
        })
        console.log(res)
        user = await res.json()
        console.log(user)
    } catch (err) {
        console.log(err)
    }
    return user
}

export async function deleteUser(user: User) :Promise<User> {
    const url = baseUrl + 'user/' + user.id
    try {
        const res = await fetch(url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: user.name,
                password: user.password,
            })
        })
        console.log(res)
        user = await res.json()
        console.log(user)
    } catch (err) {
        console.log(err)
    }
    return user
}

export async function getAllUsers(): Promise<User[]|null>{
  const url = baseUrl + 'api/user/all'
  try {
    const res = await fetch(url, {
      method: "GET",
      headers: {
        Authorization: getToken(),
      },
    })
    if (res.status === 200) {
      const users: User[] = [];
      const response = await res.json();
      for (const r of response) {
        users.push({
          id: r.id,
          name: r.name,
          password: "",
        })
      }
      return users
    }
    console.log(res);
  } catch(e) {
    console.log(e);
  }
  return null
}
