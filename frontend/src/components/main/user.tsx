import React, { useState } from "react";
import { useCookies } from "react-cookie";
import { currentUser, login } from '@fetchAPI/login'
import { resetCookie } from "@src/fetchAPI/cookie";
import router from "next/router";
import Button from "@mui/material/Button";
import { Link, useNavigate } from "react-router-dom";
import { useSetRecoilState } from "recoil";
import { getWorkspaces, Workspace} from '@fetchAPI/workspace'
import { workspacesState } from "@src/utils/atom";

const LoginForm = () => {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  const setWorkspaces = useSetRecoilState(workspacesState);
  const navigate = useNavigate();

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
        getWorkspaces().then((workspaces: Workspace[]) => {
          setWorkspaces(workspaces);
          navigate("workspace");
        });
      }
    });
  };

  return (
    <div className="App">
      <form className="px-8 py-8" onSubmit={handleSubmit}>
        <p className="text-2xl p-1">Login</p>
        <div className="mb-4">
          <label  className="block mb-2 font-bold">名前</label>
          <input className="border border-black w-full py-2 px-3" type="text" value={ name } name="name" onChange={(e) => nameChange(e)} maxLength={80} required />
        </div>
        <div className="mb-6">
          <label className="block mb-2 font-bold">パスワード</label>
          <input className="border border-black w-full py-2 px-3" type="password" value={ password } name="password" onChange={(e) => passwordChange(e)} minLength={6} maxLength={72} required />
        </div>
        <div>
          <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">ログイン</button>
        </div>
        <div>
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
