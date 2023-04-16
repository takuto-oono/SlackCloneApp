import React from "react";
import WorkspaceIndex from "src/components/workspace/workspace_index";
// import Link from 'next/link'
import { Link } from 'react-router-dom';

import { Outlet } from 'react-router-dom';


export default function IndexW() {
  return (
    <div>
    < WorkspaceIndex />
    <Link to="create">
      Create Workspace &gt;&gt;
    </Link><br></br>
    <Outlet />
    </div>
  );
}
