import { getToken } from "@/src/fetchAPI/cookie";
import router from "next/router";
import React, { useEffect, useState } from "react";
import { postWorkspace } from 'src/fetchAPI/workspace';
function CreateWorkspace() {
  const [name, setName] = useState("");
  const nameChange = (e: any) => {
    setName(e.target.value);
  };
  const handleCreate = () => {
    if (getToken() === undefined) {
      console.log("redirect");
      router.replace('/')
    } else {
      console.log("create");
      let workspaceName = name
      postWorkspace(workspaceName)
    }
  };
  useEffect(() => {
    if (getToken() === undefined) {
      console.log("redirect");
      router.replace('/')
    }
  });

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
