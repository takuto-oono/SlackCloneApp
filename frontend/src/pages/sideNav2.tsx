import React from "react";
import { Menu, ProSidebarProvider, Sidebar } from "react-pro-sidebar";
import classes from '@styles/Home.module.css'
import ChannelIndex from "@src/components/sideNav2/channel_index";
import { Outlet } from "react-router-dom";



export default function SideNav2() {
  return (
    <div style={{ display: 'flex', height: '100%', backgroundColor: "gray" }}>
      <div className={classes.item}>
            <ProSidebarProvider>
              <Sidebar>
                <Menu>
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
