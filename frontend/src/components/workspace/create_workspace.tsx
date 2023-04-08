import React, { useState } from "react";
<<<<<<< HEAD:frontend/pages/component/create_workspace.tsx
import { postWorkspace } from "pages/fetchAPI/workspace";
import Link from "next/link";

=======
import { postWorkspace } from '@fetchAPI/workspace';
>>>>>>> main:frontend/src/components/workspace/create_workspace.tsx
function CreateWorkspace() {
  const [name, setName] = useState("");
  const nameChange = (e: any) => {
    setName(e.target.value);
  };

  const handleCreate = () => {
    console.log("create");
    let workspaceName = name;
<<<<<<< HEAD:frontend/pages/component/create_workspace.tsx
    postWorkspace(workspaceName);
=======
    postWorkspace(workspaceName)
>>>>>>> main:frontend/src/components/workspace/create_workspace.tsx
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
