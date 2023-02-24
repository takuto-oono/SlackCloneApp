import React, { useState } from "react";
import { useCookies } from "react-cookie";
import { currentUser,login } from 'pages/fetchAPI/login'

function Cookie_test() {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [cookies, setCookie, removeCookie] = useCookies(['token']);

  const nameChange = (e: any) => {
    setName(e.target.value);
  };
  const passwordChange = (e: any) => {
    setPassword(e.target.value);
  };

  const handleDelete = () => {
    console.log("delete");
    removeCookie("token", {path: '/'});
  };
  const handleSubmit = () => {
    console.log("submit");
    let user = { name: name, password: password }
    login(user).then((currentuser: currentUser) => { 
      setCookie("token", currentuser.token);
    });    
  };

  return (
    <div className="App">
        <label>名前
          <input type="text" value={ name } name="name" onChange={(e) => nameChange(e)} />/ 
        </label>
        <label>パスワード
          <input type="password" value={ password } name="password" onChange={(e) => passwordChange(e)} />/ 
        </label>
        <button onClick={handleSubmit} >ログイン</button>
        <button onClick={handleDelete}>ログアウト</button>
    </div>
  );
}

export default Cookie_test;
