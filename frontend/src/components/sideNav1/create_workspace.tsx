import React, { useState } from "react";
import { postWorkspace } from '@fetchAPI/workspace';
import { useNavigate } from "react-router-dom";

function CreateWorkspace() {
  const [name, setName] = useState("");
  const navigate = useNavigate();
  const nameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  const handleSubmit = (e: React.FormEvent<HTMLFormElement> ) => {
    e.preventDefault();
    console.log("create");
    let workspaceName = name;
    postWorkspace(workspaceName);
    navigate("/workspace");
    // ワークスペースのリストを更新する(Todo)
  };
  return (
    <div className="rounded px-8 pt-6 pb-8 mb-4">
      <form onSubmit={handleSubmit}>
        <p className="text-gray-900 text-2xl p-1">Create Workspace</p>
        <div className="mb-4">
        <label  className="block text-gray-700 text-sm font-bold mb-2">ワークスペースの名前</label>
          <input className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" type="text" value={ name } name="name" onChange={(e) => nameChange(e)} maxLength={50} required/>
        </div>
        <div className="items-center">
          <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">作成</button>
        </div>
      </form>
    </div>
  );
}

export default CreateWorkspace;
