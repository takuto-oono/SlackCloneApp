import React from "react";
import { Menu, ProSidebarProvider, Sidebar } from "react-pro-sidebar";
import ShowUserChannels from "@src/components/sideNav2/show_user_channels";
import { workspaceIdState } from "@src/utils/atom";
import { useRecoilValue } from "recoil";
import { AddUserInWorkspaceForm } from "@src/components/popUp/add_user_in_workspace_form";

export default function SideNav2() {
  const workspaceId = useRecoilValue(workspaceIdState);
  
  if (workspaceId) {
    return (
      <div className="h-full" id="container">
        <div className="bg-purple-200 h-full text-pink-700 border-pink-50" >
          <div>
            <ProSidebarProvider>
              <Sidebar>
                <div className="grid grid-cols-1 divide-y divide-inherit">
                  <div className="bg-purple-200 text-pink-700">
                    {/* ToDo: WorkspaceNameを表示する */}
                  </div>
                  <div>
                    <Menu className="pd-5 bg-purple-200 text-pink-700">
                      < ShowUserChannels />
                    </Menu>
                  </div>
                  <div>
                    <Menu className="bg-purple-200 text-pink-700">
                    {/* ToDo: ShowDMs */}
                      <AddUserInWorkspaceForm workspaceID={workspaceId} />
                    </Menu>
                  </div>                  
                </div>
              </Sidebar>
            </ProSidebarProvider>
          </div>
        </div>
      </div>
    );
  } else {
    return (
      <div className="h-full" id="container">
        <div className="bg-purple-200 h-full text-pink-700 border-r-2 border-pink-50" >
          <div>
            <ProSidebarProvider>
              <Sidebar>
                <></>
              </Sidebar>
            </ProSidebarProvider>
          </div>
        </div>
      </div>
    )
  }
}
