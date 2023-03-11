export interface User {
  name: string;
  password: string;
}

const baseUrl = "http://localhost:8080/api/user/";

export async function signUp(user: User) {
  const url = baseUrl + "signUp";
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: user.name,
        password: user.password,
      }),
    });
  } catch (err) {
    console.log(err);
  }
}
