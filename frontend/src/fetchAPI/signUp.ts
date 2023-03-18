import router from "next/router";

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
    if (res.status == 200) {
        console.log("redirect");
        router.replace('/')
      }
  } catch (err) {
    console.log(err);
  }
}
