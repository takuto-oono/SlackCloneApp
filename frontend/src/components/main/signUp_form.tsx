import React, { useRef, useState } from "react";
import { signUp } from "src/fetchAPI/signUp";
import { Link } from 'react-router-dom';

const SignUpForm = () => {
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  const nameChange = (e: any) => {
    setName(e.target.value);
  };

  const passwordChange = (e: any) => {
    setPassword(e.target.value);
  };

  const handleSubmit = () => {
    console.log("signup");
    let user = { name: name, password: password };
    signUp(user);
  };

  return (
    <div>
      <h2>SingUp</h2>
      <label htmlFor="name">
        名前
        <input type="text" value={name} name="name" onChange={nameChange} />
      </label><br />

      <label htmlFor="password">
        パスワード
        <input type="password" value={password} name="password" onChange={passwordChange} />
      </label><br />

      <button type="submit" onClick={handleSubmit}>
        作成
      </button><br />

      <Link to="/login_form">
        <button>既に作成してある方はこちらへ
        </button>
      </Link>
    </div>
  );
};

export default SignUpForm;
