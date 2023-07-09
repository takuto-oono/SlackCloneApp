import React from "react";
import LoginPage from "./login_page";
import { useRecoilValue } from "recoil";
import { loginUserState } from "@src/utils/atom";

export default function Home() {
  const loginUser =  useRecoilValue(loginUserState);

  if (!loginUser) {
    return <LoginPage />;
  } else {
    return <></>;
  }
  
}
