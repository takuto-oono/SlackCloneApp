import React from "react";
import { ProSidebarProvider, Sidebar, Menu, MenuItem } from 'react-pro-sidebar';

export default function SideNav1() {
  return (
    <div style={{ display: 'flex', height: '100%', backgroundColor: "brown" }}>
      <ProSidebarProvider>
        <Sidebar>
          <Menu>
            <MenuItem> Workspace 1 </MenuItem>
            <MenuItem> Workspace 2 </MenuItem>
            <MenuItem> Create Workspace</MenuItem>
          </Menu>
        </Sidebar>
      </ProSidebarProvider>
    </div>
  );
}
