import React from "react";
import { ProSidebarProvider, Sidebar, Menu } from 'react-pro-sidebar';
import WorkspaceIndex from "@components/sideNav1/workspace_index"
import classes from '@styles/Home.module.css'
import { Outlet } from "react-router-dom";

export default function SideNav1() {
  return (
    <div style={{ display: 'flex', height: '100%', backgroundColor: "brown" }}>
      <div className={classes.container}>
        <div className={classes.item}>
            <ProSidebarProvider>
              <Sidebar>
                <Menu>
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
