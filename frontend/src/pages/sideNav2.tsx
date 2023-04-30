import React from "react";
import { Menu, ProSidebarProvider, Sidebar } from "react-pro-sidebar";
import classes from '@styles/Home.module.css'
import ChannelIndex from "@src/components/sideNav2/show_channels/[id]";
import { Outlet } from "react-router-dom";

export default function SideNav2() {
  return (
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
                  < ChannelIndex />
                </Menu>
              </Sidebar>
            </ProSidebarProvider>
          </div>
          <div className={classes.item}>
            <Outlet />
          </div>
    </div>
  );
}
