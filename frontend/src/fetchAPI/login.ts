export interface User {
    name: string;
    password: string;
}

export interface currentUser {
  token: string;
  user_id: string;
  username: string;
}

const baseUrl = 'http://localhost:8080/api/user/';

export async function login(user: User): Promise<currentUser> {
  const url = baseUrl + 'login';
  const currentUser = {
    token: "",
    user_id: "",
    username: ""
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
            })
        })
      console.log(res);
      const User = await res.json();

      return new Promise((resolve) => {
        const currentUser: currentUser = {
          token: User.token,
          user_id: User.user_id,
          username: User.username
        };
        console.log(currentUser);
        resolve(currentUser);
      });


    } catch (err) {
      console.log(err);
  }

  return currentUser;
}
