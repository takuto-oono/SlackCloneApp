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
      <form className="rounded px-8 pt-6 pb-8 mb-4" onSubmit={handleSubmit}>
      <p className="text-gray-900 text-2xl p-1">SingUp</p>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2">名前</label>
        <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" value={name} name="name" onChange={nameChange} maxLength={80} required />
      </div>
      <div className="mb-6">
        <label className="block text-gray-700 text-sm font-bold mb-2">パスワード</label>
        <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="password" value={password} name="password" onChange={passwordChange} minLength={6} maxLength={72} required/>
      </div>
      <div className="items-center">
        <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">作成</button>
      </div>
      <div className="items-center">
        <Link to="/">
            <button className="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800">既に作成してある方はこちらへ</button>
        </Link>
      </div>
      </form>
    </div>
  );
};

export default SignUpForm;
