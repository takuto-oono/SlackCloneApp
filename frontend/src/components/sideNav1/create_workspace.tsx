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
    <div className="CreateWorkspace">
      <form onSubmit={handleSubmit}>
        <h2>Create Workspace</h2>
        <label>ワークスペースの名前
          <input type="text" value={ name } name="name" onChange={(e) => nameChange(e)} maxLength={50} required/>
        </label><br />
        <input type="submit" value="作成" />
      </form>
    </div>
  );
}

export default CreateWorkspace;
