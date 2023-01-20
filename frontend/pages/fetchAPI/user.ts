interface User {
    id: string;
    name: string;
    password: string;
}

const baseUrl = 'http://localhost:8000/'

export async function getUsers() {
    const url = baseUrl + 'users'
    try {
        const res = await fetch(url, {
            method: 'GET',
        })
        console.log(res)
        const data = await res.json()
        console.log(data)
        return data
    } catch (err) {
        console.log(err)
    }
}

export async function getUserById() {
    const id = "1bd5adb749d24b60a364b5b864d49f10"
    const url = baseUrl + 'user/' + id
    try {
        const res = await fetch(url, {
            method: 'GET',
        })
        console.log(res)
        const data: User = await res.json()
        console.log(data)
        return data
    } catch (err) {
        console.log(err)
    }
}

export async function postUser() {
    const url = baseUrl + 'user'
    try {
        const res = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: "admin3",
                password: 'admin3',
            })
        })
        console.log(res)
        const data = await res.json()
        console.log(data)
        return data
    } catch (err) {
        console.log(err)
    }
}

export async function updateUser() {
    const url = baseUrl + 'user/' + '4c0bb0020bec47d9a8f9e1cc2db1822e'
    try {
        const res = await fetch(url, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: "joker",
                password: 'joker',
            })
        })
        console.log(res)
        const data = await res.json()
        console.log(data)
        return data
    } catch (err) {
        console.log(err)
    }
}

export async function deleteUser() {
    const url = baseUrl + 'user/' + '4c0bb0020bec47d9a8f9e1cc2db1822e'
    try {
        const res = await fetch(url, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: "joker",
                password: 'joker',
            })
        })
        console.log(res)
        const data = await res.json()
        console.log(data)
        return data
    } catch (err) {
        console.log(err)
    }
    
}

