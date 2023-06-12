import React from "react";
import { createRoot } from "react-dom/client";
import { RecoilRoot } from "recoil";
import Layout from "./common/layout";
import LoginPage from "./login_page";

export default function Home() {
  return <LoginPage />;
}
