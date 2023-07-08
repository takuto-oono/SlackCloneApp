import React from "react";
import LoginPage from "./login_page";
import { useRecoilValue } from "recoil";
import { loginUserState } from "@src/components/main/user";

export default function Home() {
  const loginUser =  useRecoilValue(loginUserState);

  if (loginUser.length == 0) {
    return <LoginPage />;
  } else {
    return <></>;
  }
  
}
