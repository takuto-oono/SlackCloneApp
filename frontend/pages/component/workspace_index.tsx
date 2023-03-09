import React, { useState } from "react";
import { getWorkspaces, Workspace } from 'pages/fetchAPI/workspace';

function WorkspaceIndex() {
  const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);
  const list = workspaceList.map((item, index) => (
    <div key={index}>
      <p>{item.id}</p>
      <p>{item.name}</p>
      <p>{item.primary_owner_id}</p>
    </div>
  ));
  
  const handleGetWorkspaces = () => {
        getWorkspaces().then((workspaces: Workspace[]) => {
        console.log("workspaces")
          console.log(workspaces)
          if (Array.isArray(workspaces)) {
            setWorkspaceList(workspaces)
          }
        console.log(workspaceList)
      });
  }
  
    return (
      <div className="App">
        <button onClick={handleGetWorkspaces}>ワークスペース一覧を表示</button>
        <div>
          {list}
        </div>
      </div>
    );
}

export default WorkspaceIndex;
