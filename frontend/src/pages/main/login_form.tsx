import React, { useEffect } from "react";
import { useCookies } from "react-cookie";
import { LoginForm } from "@src/components/main/user";
import { useNavigate } from 'react-router-dom';
import { Workspace, getWorkspaces } from "@src/fetchAPI/workspace";
import { useSetRecoilState } from "recoil";
import { workspacesState } from "@src/utils/atom";


function Login() {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  const setWorkspaces = useSetRecoilState(workspacesState);
  const navigate = useNavigate();

  // useEffect削除（未）
  useEffect(()=>{
    if (cookies.token) {
      getWorkspaces().then((workspaces: Workspace[]) => {
        setWorkspaces(workspaces);
        navigate("workspace");
      });
    }
  }, [cookies.token])
  
  return (
    <div>
      < LoginForm />
    </div>
  );
}

export default Login;


