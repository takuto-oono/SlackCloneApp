import React, { useState } from "react";
import { signUp } from "src/fetchAPI/signUp";
import { Link } from 'react-router-dom';

const SignUpForm = () => {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  const passwordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement> ) => {
    e.preventDefault();
    console.log("signup");
    let user = { name: name, password: password };
    signUp(user);
  };

  return (
    <div>
      <form onSubmit={handleSubmit}>
      <h2>SingUp</h2>
      <label htmlFor="name">
        名前
        <input type="text" value={name} name="name" onChange={nameChange} maxLength={80} required />
      </label><br />
      <label htmlFor="password">
        パスワード
        <input type="password" value={password} name="password" onChange={passwordChange} minLength={6} maxLength={72} required/>
      </label><br />
      <input type="submit" value="作成" />
    </form>
      <Link to="/">
        <button>既に作成してある方はこちらへ
        </button>
      </Link>
    </div>
  );
};

export default SignUpForm;
