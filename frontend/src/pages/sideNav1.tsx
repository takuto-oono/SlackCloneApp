import React from "react";
import { ProSidebarProvider, Sidebar, Menu } from 'react-pro-sidebar';
import WorkspaceIndex from "@src/components/sideNav1/show_workspaces"
import classes from '@styles/Home.module.css'
import { Outlet } from "react-router-dom";

export default function SideNav1() {
  return (
    <div style={{ display: 'flex', height: '100%', backgroundColor: "rgb(63, 14, 64)" }}>
      <div className={classes.container}>
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
                </Menu>
              </Sidebar>
            </ProSidebarProvider>
          </div>
          <div className={classes.item}>
            <Outlet />
          </div>
        </div>
    </div>
  );
}
