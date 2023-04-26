import React from "react";
import { Menu, MenuItem, Sidebar, SubMenu } from "react-pro-sidebar";
import ShowWorkspace from "@src/components/sideNav2/workspace/show/[id]";



export default function SideNav2() {
  return (
    <div style={{ display: 'flex', height: '100%' ,backgroundColor: "gray"}}>
      <Sidebar>
        <Menu>
          <SubMenu label="Channel Index">
            <MenuItem> Channel 1 </MenuItem>
            <MenuItem> Channel 2 </MenuItem>
          </SubMenu>
          <SubMenu label="DM Index">
            <MenuItem> DM 1 </MenuItem>
            <MenuItem> DM 2 </MenuItem>
          </SubMenu>
        </Menu>
      </Sidebar>
    </div>
  );
}
