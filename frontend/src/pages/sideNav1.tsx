import React from "react";
import { ProSidebarProvider, Sidebar, Menu, MenuItem } from 'react-pro-sidebar';
import WorkspaceIndex from "@src/components/sideNav1/show_workspaces"
import classes from '@styles/Home.module.css'
import { Link, Outlet } from "react-router-dom";

export default function SideNav1() {
  return (
    
    <div className={classes.container}>
      <div style={{ display: 'flex', height: '100%', backgroundColor: "rgb(63, 14, 64)" }}>
        <div className={classes.item}>
            <ProSidebarProvider>
              <Sidebar>
                <Menu
                  menuItemStyles={{
                    button: () => {
                        return {
                          color: 'rgb(153, 133, 156)',
                          backgroundColor: 'rgb(63, 14, 64)' 
                        };
                    },
                  }}
                >
                < WorkspaceIndex />
                <MenuItem>
                  <Link to="create" >
                    create
                  </Link>

                {/* テスト用 */}
                </MenuItem>
                </Menu>
              </Sidebar>
            </ProSidebarProvider>
          </div>
      </div>
      <div>
        <Outlet />
      </div>
    </div>
  );
}
