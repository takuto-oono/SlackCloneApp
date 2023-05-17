import React from "react";
import { Menu, ProSidebarProvider, Sidebar } from "react-pro-sidebar";
import classes from '@styles/Home.module.css'
import { Outlet } from "react-router-dom";
import ShowChannels from "@src/components/sideNav2/show_channels/[id]";

export default function SideNav2() {
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
                < ShowChannels />
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
