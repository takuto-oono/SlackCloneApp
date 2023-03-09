import React, { useState } from "react";
import { getWorkspaces, Workspace } from 'pages/fetchAPI/workspace'
import { getToken } from "@/pages/fetchAPI/cookie";
import Link from 'next/link'
// import { Link } from "react-router-dom";
const query = {
  id: 1,
  name: "yakkun",
}

function WorkspaceIndex() {
  const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);

  const list = workspaceList.map((item, index) => (
    <div key={index}>
      <Link href={"/component/show_workspace/"+item.id}>
          Workspace{item.id} &gt;&gt;
      </Link><br></br>
      <p>{item.id}</p>
      <p>{item.name}</p>
      <p>{item.primary_owner_id}</p>
    </div>
  ));
  const handleGetWorkspaces = () => {
    if (typeof getToken() !== 'undefined') {
      getWorkspaces().then((workspaces: Workspace[]) => {
        console.log("workspaces")
        console.log(workspaces)
        setWorkspaceList(workspaces)
        console.log(workspaceList)
      });
    }
  }

  return (
    <div className="App">
      <Link href="/component/create_workspace">
          Create Workspace &gt;&gt;
      </Link><br></br>
      <button onClick={handleGetWorkspaces}>ワークスペース一覧を表示</button>
      <div>
        {list}
      </div>
    </div>
  );
}

export default WorkspaceIndex;
