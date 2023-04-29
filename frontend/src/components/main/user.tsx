import React, { useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { currentUser, login } from '@fetchAPI/login'
import { resetCookie } from "@src/fetchAPI/cookie";
import router from "next/router";


const LoginForm = () => {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [cookies, setCookie, removeCookie] = useCookies(['token','user_id']);

  const nameChange = (e: any) => {
    setName(e.target.value);
  };
  const passwordChange = (e: any) => {
    setPassword(e.target.value);
  };
  const handleLogin = () => {
    console.log("login");
    let user = { name: name, password: password }
    login(user).then((currentUser: currentUser) => { 
      if (currentUser.token != "" && currentUser.token) {
        setCookie("token", currentUser.token);
        setCookie("user_id", currentUser.user_id);
      }
    });
    
  };

  return (
    <div className="App">
        <label>名前
          <input type="text" value={ name } name="name" onChange={(e) => nameChange(e)} />
        </label><br />
        <label>パスワード
          <input type="password" value={ password } name="password" onChange={(e) => passwordChange(e)} />
        </label><br />
        <button onClick={handleLogin} >ログイン</button>
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
      <button onClick={handleLogout}>ログアウト</button>
    </div>
  );
}

export { Logout };
