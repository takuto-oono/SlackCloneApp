// import TestAPI1 from "./api/test_api_1";
// import { User, getUsers, getUserById, postUser, updateUser, deleteUser, } from "./fetchAPI/user";
import React from "react";
import LoginForm from "./component/login_form";
import WorkspaceIndex from "./component/workspace_index";
import CreateWorkspace from "./component/create_workspace";
import { CookiesProvider } from "react-cookie";
// import testAPI1 from "./api/test_api_1";
export default function Home() {
  return (
    <CookiesProvider>
      <main>
        <h1>hello nextjs</h1>
        <h2>login</h2>
        < LoginForm />
        < WorkspaceIndex />
        < CreateWorkspace />
      </main>
    </CookiesProvider>
  )
}


