// import TestAPI1 from "./api/test_api_1";
// import { User, getUsers, getUserById, postUser, updateUser, deleteUser, } from "./fetchAPI/user";
import React from "react";
import LoginForm from "./component/login_form";
import { CookiesProvider } from "react-cookie";
import Link from 'next/link'

// import testAPI1 from "./api/test_api_1";
export default function Home() {
  return (
    <CookiesProvider>
      <main>
        <h1>hello nextjs</h1>
        <h2>login</h2>
        < LoginForm />
        <Link href="/component/workspace_index">
          Go to WorkspaceIndex &gt;&gt;
        </Link>
      </main>
    </CookiesProvider>
  )
}


