import React, { useRef, useState } from "react";
import { useCookies } from "react-cookie";
import { signUp } from 'pages/fetchAPI/signUp'
import Link from "next/link"

const SignUpForm = () => {
	const [name, setName] = useState("");
	const [password, setPassword] = useState("");
	const [cookies, setCookie, removeCookie] = useCookies(['token']);

	const nameCreate = (e: any) => {
    setName(e.target.value);
  };

  const passwordCreate = (e: any) => {
    setPassword(e.target.value);
  };

  const handleSubmit = () => {
    console.log("signup");
    let user = { name: name, password: password }
    signUp(user).then((currentUser: any) => { 
      setCookie("token", currentUser.token);
    });
  };

	return (
		<div>
			<h2>SingUp</h2>
			<label htmlFor="name">名前
				<input type="text" value={ name } name="name" onChange={(e) => nameCreate(e)} />
			</label><br />
				<label htmlFor="password">パスワード
			<input type="password" value={ password } name="password" onChange={(e) => passwordCreate(e)} />
			</label><br />
			<button type="submit" onClick={handleSubmit}>作成</button><br />
			<Link href="/">
				<button>既に作成してある方はこちらへ</button>
			</Link>
		</div>
	)
}

export default SignUpForm