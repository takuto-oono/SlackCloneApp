// import TestAPI1 from "./api/test_api_1";
// import { User, getUsers, getUserById, postUser, updateUser, deleteUser, } from "./fetchAPI/user";
import React from "react";
import LoginForm from "./component/login_form";
import WorkspaceIndex from "./component/workspace_index";
import SignUpForm from "./component/signUp_form";
import { CookiesProvider } from "react-cookie";
import Link from "next/link";
// import testAPI1 from "./api/test_api_1";
export default function Home() {
  // console.log(getUsers())
  // let user: User = {
  //   id: "",
  //   name: "test docker",
  //   password: "test docker"
  // }
  // let userRes = postUser(user)
  // console.log(userRes)
  return (
    <CookiesProvider>
      <main>
        <h1>hello nextjs</h1>
        <h2>login</h2>
        <LoginForm />
        <WorkspaceIndex />
        <Link href="component/signUp_form">
          <button>まだアカウントを持っていませんか？</button>
        </Link>
      </main>
    </CookiesProvider>
  );
}
