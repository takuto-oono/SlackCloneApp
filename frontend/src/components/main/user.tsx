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
  const handleSubmit = (e: React.FormEvent<HTMLFormElement> ) => {
    e.preventDefault();
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
      <form className="rounded px-8 pt-6 pb-8 mb-4" onSubmit={handleSubmit}>
        <p className="text-gray-900 text-2xl p-1">Login</p>
        <div className="mb-4">
          <label  className="block text-gray-700 text-sm font-bold mb-2">名前</label>
          <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" value={ name } name="name" onChange={(e) => nameChange(e)} maxLength={80} required />
        </div>
        <div className="mb-6">
          <label className="block text-gray-700 text-sm font-bold mb-2">パスワード</label>
          <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="password" value={ password } name="password" onChange={(e) => passwordChange(e)} minLength={6} maxLength={72} required />
        </div>
        <div className="items-center">
          <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">ログイン</button>
        </div>
        <div className="items-center">
          <Link to="/signUp_form">
              <button className="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800">まだアカウントを持っていませんか？</button>
          </Link>
        </div>
      </form>
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
