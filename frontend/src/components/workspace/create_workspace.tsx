import React, { useState } from "react";
import { postWorkspace } from 'src/fetchAPI/workspace';
function CreateWorkspace() {
  const [name, setName] = useState("");
  const nameChange = (e: any) => {
    setName(e.target.value);
  };
  const handleCreate = () => {
    console.log("create");
    let workspaceName = name
    postWorkspace(workspaceName)
    };

  return (
    <div className="CreateWorkspace">
        <h2>Create Workspace</h2>
        <label>ワークスペースの名前
          <input type="text" value={ name } name="name" onChange={(e) => nameChange(e)} />
        </label><br />
        <button onClick={handleCreate} >作成</button>
    </div>

  );

}

export default CreateWorkspace;
