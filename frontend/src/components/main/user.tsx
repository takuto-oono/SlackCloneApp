import React, { useState } from "react";
import { useCookies } from "react-cookie";
import { currentUser, login } from '@fetchAPI/login'
import { resetCookie } from "@src/fetchAPI/cookie";
import router from "next/router";
import Button from "@mui/material/Button";
import { Link } from "react-router-dom";




const LoginForm = () => {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [cookies, setCookie, removeCookie] = useCookies(['token','user_id']);

  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };
  const passwordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };
  const handleSubmit = () => {
    console.log("login");
    let user = { name: name, password: password }
    login(user).then((currentUser: currentUser) => { 
      if (currentUser.token) {
        setCookie("token", currentUser.token);
        setCookie("user_id", currentUser.user_id);
      }
    });
    
  };

  return (
    <div className="App">
      <form onSubmit={handleSubmit}>
        <h2>Login</h2>
        <label>名前
          <input type="text" value={ name } name="name" onChange={(e) => nameChange(e)} maxLength={80} required />
        </label><br />
        <label>パスワード
          <input type="password" value={ password } name="password" onChange={(e) => passwordChange(e)} minLength={6} maxLength={72} required />
        </label><br />
        <input type="submit" value="ログイン" />
      </form>
      <Link to="/signUp_form">
        <button>まだアカウントを持っていませんか？</button>
      </Link>
    </div>
  );
}

export { LoginForm };
  
const Logout = () => {
  const handleLogout = () => {
    console.log("logout");
    resetCookie();
    router.push("/");
  };

  return (
    <div>
      <Button variant="contained" color="secondary" onClick={handleLogout}>ログアウト</Button>
    </div>
  );
}

export { Logout };
