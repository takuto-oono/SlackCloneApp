export interface User {
    name: string;
    password: string;
}

const baseUrl = 'http://localhost:8080/api/user/'

export async function login(user: User): Promise<User> {
    const url = baseUrl + 'login'

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
        const currentuser = await res.json()
        console.log(currentuser)
        // Dom
        const Token = currentuser.token;
        console.log(Token)
        alert(Token)

    } catch (err) {
        console.log(err)
    }
    return user
}
