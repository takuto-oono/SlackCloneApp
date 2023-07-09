import React from "react";
import { ProSidebarProvider, Sidebar, Menu, MenuItem } from 'react-pro-sidebar';
import ShowWorkspaces from "@src/components/sideNav1/show_workspaces";
import { useRouter } from "next/router";
import { userChannelsState, loginUserState, workspaceIdState } from "@src/utils/atom";
import { useRecoilValue, useResetRecoilState } from "recoil";

export default function SideNav1() {
  const router = useRouter()
  const resetUserChannelsState = useResetRecoilState(userChannelsState);
  const resetWorkspaceIdState =  useResetRecoilState(workspaceIdState);
  const loginUser =  useRecoilValue(loginUserState);

  const exitWorkspace = () => {
    resetUserChannelsState();
    resetWorkspaceIdState();
    router.push('/create_workspace')
  }
  if (loginUser) {
    return (
      <div className="h-full" id="container">
        <div className="bg-purple-200 h-full text-pink-700 border-r-2 border-pink-50">
          <div>
            <ProSidebarProvider>
              <Sidebar>
                <Menu className="bg-purple-200 text-pink-700">
                  < ShowWorkspaces />
                  <MenuItem>
                    <button type="button" onClick={() => exitWorkspace() }>
                      <>create</>
                    </button>
                  </MenuItem>
                </Menu>
              </Sidebar>
            </ProSidebarProvider>
          </div>
        </div>
      </div>
    );
  } else {
    return (
      <div className="h-full" id="container">
        <div className="bg-purple-200 h-full text-pink-700 border-r-2 border-pink-50">
          <div>
            <ProSidebarProvider>
              <Sidebar>
                <></>
              </Sidebar>
            </ProSidebarProvider>
          </div>
        </div>
      </div>
    );
  }

  
}
