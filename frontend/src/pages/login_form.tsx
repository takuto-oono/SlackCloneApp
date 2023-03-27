import React, { useEffect } from "react";
import { useCookies } from "react-cookie";
import { LoginForm } from "@components/login";
import { Link } from 'react-router-dom';
import router from "next/router";

function Login() {
  const [cookies, setCookie, removeCookie] = useCookies(['token','user_id']);
  useEffect(() => {
    if (cookies.token) {
      console.log(cookies.token)
      router.push("/")
      
    }
  })
  return <div>
    <h2>Login</h2>
    < LoginForm />
    <br></br>
    <Link to="/signUp_form">
      <button>まだアカウントを持っていませんか？</button>
    </Link>
  </div>;
}

export default Login;
