import React, { useEffect } from "react";
import { useCookies } from "react-cookie";
import { LoginForm } from "@src/components/main/user";
import { useNavigate } from 'react-router-dom';

function Login() {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  const navigate = useNavigate();
  useEffect(() => {
    if (cookies.token) {
      navigate("workspace");
    }
  },[cookies.token])
  return (
    <div>
      < LoginForm />
    </div>
  );
}

export default Login;
