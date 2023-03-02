import React, { useState } from "react";
import { useCookies } from "react-cookie";
import { currentUser, login } from 'pages/fetchAPI/login'


function LoginForm() {

  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [userId, setUserId] = useState("");
  const [cookies, setCookie, removeCookie] = useCookies(['token']);
  

  const nameChange = (e: any) => {
    setName(e.target.value);
  };
  const passwordChange = (e: any) => {
    setPassword(e.target.value);
  };

  const handleLogout = () => {
    console.log("logout");
    removeCookie("token", {path: '/'});
  };
  const handleLogin = () => {
    console.log("login");
    let user = { name: name, password: password }
    login(user).then((currentUser: currentUser) => { 
      setCookie("token", currentUser.token);
      setUserId(currentUser.user_id);
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
