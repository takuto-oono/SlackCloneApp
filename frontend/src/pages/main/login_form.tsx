import React, { useEffect } from "react";
import { useCookies } from "react-cookie";
import { LoginForm } from "@src/components/main/user";
import { Link,useNavigate } from 'react-router-dom';
import router from "next/router";

function Login() {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  const navigate = useNavigate();
  useEffect(() => {
    if (cookies.token) {
      // console.log(cookies.token)
      navigate("workspace");
      
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
