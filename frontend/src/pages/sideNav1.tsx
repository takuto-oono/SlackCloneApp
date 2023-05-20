import React from "react";
import { ProSidebarProvider, Sidebar, Menu, MenuItem } from 'react-pro-sidebar';
import { Link, Outlet } from "react-router-dom";
import ShowWorkspaces from "@src/components/sideNav1/show_workspaces";

export default function SideNav1() {
  return (
    
    <div className="h-full flex" id="container">
      {/* SideNav1表示用item */}
      <div className="bg-purple-200 h-full flex text-pink-700 border-r-2 border-pink-50"  id="item">
        <div>
          <ProSidebarProvider>
            <Sidebar>
              <Menu className="bg-purple-200 text-pink-700">
                < ShowWorkspaces />
                <MenuItem>
                  <Link to="create" >
                    create
                  </Link>
                </MenuItem>
              </Menu>
            </Sidebar>
          </ProSidebarProvider>
        </div>
      </div>
      {/* SideNav2,main,sub表示用item */}
      <div id="item">
        <Outlet />
      </div>
    </div>
  );
}
