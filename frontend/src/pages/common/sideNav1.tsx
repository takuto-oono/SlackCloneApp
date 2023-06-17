import React from "react";
import { ProSidebarProvider, Sidebar, Menu, MenuItem } from 'react-pro-sidebar';
import ShowWorkspaces from "@src/components/sideNav1/show_workspaces";
import Link from "next/link";

export default function SideNav1() {
  return (
    <div className="h-full" id="container">
      <div className="bg-purple-200 h-full text-pink-700 border-r-2 border-pink-50">
        <div>
          <ProSidebarProvider>
            <Sidebar>
              <Menu className="bg-purple-200 text-pink-700">
                < ShowWorkspaces />
                <MenuItem>
                  {/* <Link href="/create">
                    create
                  </Link> */}
                </MenuItem>
              </Menu>
            </Sidebar>
          </ProSidebarProvider>
        </div>
      </div>
    </div>
  );
}
