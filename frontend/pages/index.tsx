import React from "react";
import LoginForm from "./component/login_form";
import Link from 'next/link'
export default function Home() {
  return (
      <main>
        <h2>Login</h2>
        < LoginForm />
        <p>---</p>
        <Link href="/component/workspace_index">
          WorkspaceIndex &gt;&gt;
        </Link>
      </main>
  )
}


