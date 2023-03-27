import React from "react";
import LoginForm from "@/components/login_form";
import Link from "next/link";
const Login= () => {
  return (
    <div>
      <h2>Login</h2>
      < LoginForm />
      <br></br>
      <Link href="signUp_form">
        <button>まだアカウントを持っていませんか？</button>
      </Link>
    </div>
  )
}

export default Login;
