import React, { useEffect, useState } from "react";
import { getWorkspaces, Workspace } from 'src/fetchAPI/workspace'
import { getToken } from 'src/fetchAPI/cookie'
import Link from 'next/link'
import router from "next/router";


function WorkspaceIndex() {
  const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);
  const list = workspaceList.map((item, index) => (
    <div key={index}>
      <Link href={{ pathname: "/workspace_show/" + item.id, query: { id: item.id, name: item.name, primary_owner_id: item.primary_owner_id } }} as={"/workspace_show/"+item.id}>
          {item.name} &gt;&gt;
      </Link><br></br>
    </div>
  ));


  useEffect(() => {
    getWorkspaces().then((workspaces: Workspace[]) => {
      if (!Array.isArray(workspaces)) {
        console.log("redirect");
        router.replace('/');
      } else {
        setWorkspaceList(workspaces)
        console.log(workspaceList)
      }
    });
  },[]);

  return (
    <div className="App">
      <h2>Workspace Index</h2>
      <br></br>
      <div>
        {list}
      </div><br></br>
    </div>
    );
}

export default WorkspaceIndex;
