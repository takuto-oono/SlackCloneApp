import React, { useState } from "react";
import { postWorkspace } from "pages/fetchAPI/workspace";
import Link from "next/link";

function CreateWorkspace() {
  const [name, setName] = useState("");
  const nameChange = (e: any) => {
    setName(e.target.value);
  };

  const handleCreate = () => {
    console.log("create");
    let workspaceName = name;
    postWorkspace(workspaceName);
  };

  return (
    <div className="CreateWorkspace">
        <h2>Create Workspace</h2>
        <label>ワークスペースの名前
          <input type="text" value={ name } name="name" onChange={(e) => nameChange(e)} />
        </label><br />
        <button onClick={handleCreate} >作成</button>
        
        {/* テスト用 */}
        <Link href="/">
          <button>ログイン画面へ</button>
        </Link>
    </div>
  );
}

export default CreateWorkspace;
