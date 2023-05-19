import React from "react";
import { Menu, ProSidebarProvider, Sidebar } from "react-pro-sidebar";
import { Outlet } from "react-router-dom";
import ShowChannels from "@src/components/sideNav2/show_channels/[id]";

export default function SideNav2() {
  return (
    <div className="h-full flex" id="container">
      {/* SideNav2表示用item */}
      <div className="bg-purple-200 h-full flex text-pink-700 border-r-2 border-pink-50" id="item">
        <div>
          <ProSidebarProvider>
            <Sidebar>
              <Menu className="bg-purple-200 text-pink-700">
                < ShowChannels />
              </Menu>
            </Sidebar>
          </ProSidebarProvider>
        </div>
      </div>
      {/* main,sub表示用item */}
      <div id="item">
        <Outlet />
      </div>
    </div>
  );
}
