import React from "react";
import LoginForm from "src/components/login_form";
import Link from "next/link";
export default function Home() {
  return (
      <main>
        <h2>Login</h2>
        < LoginForm />
        <p>---</p>
        <Link href="/component/workspace_index">
          WorkspaceIndex &gt;&gt;
        </Link>
        <Link href="component/signUp_form">
          <button>まだアカウントを持っていませんか？</button>
        </Link>
      </main>
  )
}
