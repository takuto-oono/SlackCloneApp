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
  const currentuser = {
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
      const tempUser = await res.json();

      return new Promise((resolve) => {
        const currentuser: currentUser = {
          token: tempUser.token,
          user_id: tempUser.user_id,
          username: tempUser.username
        };
        console.log(currentuser);
        resolve(currentuser);
      });


    } catch (err) {
      console.log(err);
  }

  return currentuser;
}
