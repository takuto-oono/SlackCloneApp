import React from "react";
import { Menu, ProSidebarProvider, Sidebar } from "react-pro-sidebar";
import ShowChannels from "@src/components/sideNav2/show_channels/[id]";

export default function SideNav2() {
  return (
    <div className="h-full" id="container">
      <div className="bg-purple-200 h-full text-pink-700 border-r-2 border-pink-50" >
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
    </div>
  );
}
