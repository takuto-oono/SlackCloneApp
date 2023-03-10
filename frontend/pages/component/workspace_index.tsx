import React, { useEffect, useState } from "react";
import { getWorkspaces, Workspace } from 'pages/fetchAPI/workspace'
import { getToken } from "@/pages/fetchAPI/cookie";
import Link from 'next/link'


function WorkspaceIndex() {
  const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);

  const list = workspaceList.map((item, index) => (
    <div key={index}>
      <Link href={{ pathname: "/component/show_workspace/"+item.id, query: item.id.toString}} as={"/component/show_workspace/"+item.id}>
          Workspace{item.id} &gt;&gt;
      </Link><br></br>
    </div>
  ));

  useEffect(() => {
    if (typeof getToken() !== 'undefined') {
      getWorkspaces().then((workspaces: Workspace[]) => {
      if (Array.isArray(workspaces)) {
        setWorkspaceList(workspaces)
      }
      console.log(workspaceList)
      });
    }
  },[]);

  return (
    <div className="App">
      <h2>Workspace Index</h2>
      <Link href="/component/create_workspace">
          Create Workspace &gt;&gt;
      </Link><br></br>
      <p>---</p>
      <div>
        {list}
      </div>
    </div>
  );
}

export default WorkspaceIndex;
