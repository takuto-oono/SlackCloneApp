import router from "next/router";
import React, { useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { currentUser, login } from 'src/fetchAPI/login'

function LoginForm() {

  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [cookies, setCookie, removeCookie] = useCookies(['token','user_id']);
  

  const nameChange = (e: any) => {
    setName(e.target.value);
  };
  const passwordChange = (e: any) => {
    setPassword(e.target.value);
  };

  const handleLogout = () => {
    console.log("logout");
    removeCookie("token", {path: '/' });
    removeCookie("user_id", {path: '/'});
  };
  const handleLogin = () => {
    console.log("login");
    let user = { name: name, password: password }
    login(user).then((currentUser: currentUser) => { 
      setCookie("token", currentUser.token);
      setCookie("user_id", currentUser.user_id);
      router.replace('/workspace_index')
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
      <button onClick={handleLogout}>ログアウト</button>
    </div>
  );
}

export default LoginForm;
