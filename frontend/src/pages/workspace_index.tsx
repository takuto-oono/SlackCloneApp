import React from "react";
import WorkspaceIndex from "./component/workspace/workspace_index";
import Link from 'next/link'
export default function IndexW() {
  return (
      <main>
        < WorkspaceIndex />
        <Link href="/workspace_create">
            Create Workspace &gt;&gt;
        </Link><br></br>
      </main>
  )
}
