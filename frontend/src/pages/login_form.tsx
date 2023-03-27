import React from "react";
import { LoginForm } from "@components/login";
import { Link } from 'react-router-dom';
function Login() {
  return <div>
    <h2>Login</h2>
    < LoginForm />
    <br></br>
    <Link to="/signUp_form">
      <button>まだアカウントを持っていませんか？</button>
    </Link>
  </div>;
}

export default Login;
