export interface User {
    name: string;
    password: string;
}

export interface currentUser {
  id: number;
  name: string;
  password: string;
}

const baseUrl = 'http://localhost:8080/api/user/';

export async function signUp(user: User): Promise<currentUser> {
  const url = baseUrl + 'signUp';
  const currentuser = {
    id: 0,
    name: "",
    password: ""
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
          id: tempUser.id,
          name: tempUser.name,
          password: tempUser.password
        };
        console.log(currentuser);
        resolve(currentuser);
      });


    } catch (err) {
      console.log(err);
  }

  return currentuser;
}
