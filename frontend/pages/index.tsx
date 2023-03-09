import React from "react";
import LoginForm from "./component/login_form";
import Link from 'next/link'
export default function Home() {
  return (
      <main>
        <h1>hello nextjs</h1>
        <h2>login</h2>
        < LoginForm />
        <Link href="/component/workspace_index">
          Go to WorkspaceIndex &gt;&gt;
        </Link>
      </main>
  )
}


