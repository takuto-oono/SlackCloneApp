import React, { useEffect, useState } from "react";
import { getWorkspaces, Workspace } from '@fetchAPI/workspace'
import router from "next/router";
import { Link } from 'react-router-dom';
import { MenuItem } from "react-pro-sidebar";




function WorkspaceIndex() {
  const [workspaceList, setWorkspaceList] = useState<Workspace[]>([]);
  const list = workspaceList.map((item, index) => (
    <div key={index}>
      {/* workspaceオブジェクトも渡したい（未） */}
      <MenuItem>
        <Link to={`${item.id}`}>
          <span>{item.name}</span>
        </Link>
      </ MenuItem>
    </div>
  ));


  useEffect(() => {
    getWorkspaces().then((workspaces: Workspace[]) => {
      // if (Array.isArray(workspaces)) {
      //   setWorkspaceList(workspaces);
      // }
      // console.log(workspaces);
      setWorkspaceList(workspaces);
    });
  },[]);

  return (
    <div className="App">
      <div>
        {list}
      </div>      
    </div>
    );
}

export default WorkspaceIndex;
