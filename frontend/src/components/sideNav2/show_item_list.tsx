import React from "react";
import { Menu, MenuItem } from "react-pro-sidebar";
import { useRouter } from "next/router";

function ShowItemList() {
  const router = useRouter()

  return (
    <div>
      <Menu className="pd-5 bg-purple-200 text-pink-700">
        < ShowWorkspaceChannels />
      </Menu>
    </div>
  )
}

export default ShowItemList;

const ShowWorkspaceChannels = () => {
  const router = useRouter()

  return (
    <MenuItem>
      <button
        onClick={() => router.push({
          pathname: `/main`,
          query: { contents: "show_workspace_channels" },
        })}>
        <p>すべてのチャンネル</p></button>
    </MenuItem>
  )
}
