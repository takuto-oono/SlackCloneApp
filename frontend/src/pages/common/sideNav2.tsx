import React from 'react'
import { Menu, MenuItem, ProSidebarProvider, Sidebar } from 'react-pro-sidebar'
import ShowJoinedChannels from '@src/components/channel/show_joined_channels'
import { workspaceIdState, workspacesState } from '@src/utils/atom'
import { useRecoilValue } from 'recoil'
import { AddUserInWorkspaceForm } from '@src/components/workspace/add_user_in_workspace_form'
import ShowContentsList from '@src/components/content/show_contents_list'

export default function SideNav2() {
  const workspaceId = useRecoilValue(workspaceIdState)
  const workspaces = useRecoilValue(workspacesState)
  const currentWorkspace = workspaces.find((workspace) => workspace.id === workspaceId)
  const workspaceName = currentWorkspace?.name

  if (workspaceId) {
    return (
      <div className='h-full' id='container'>
        <div className='bg-purple-200 h-full text-pink-700 border-pink-50'>
          <div>
            <ProSidebarProvider>
              <Sidebar>
                <div className='grid grid-cols-1 divide-y divide-inherit'>
                  <Menu className='pd-5 bg-purple-200 text-pink-800 text-lg'>
                    <MenuItem>
                      <button>{workspaceName}</button>
                    </MenuItem>
                  </Menu>
                  <div>
                    <ShowContentsList />
                  </div>
                  <div>
                    <ShowJoinedChannels />
                  </div>
                  <div>
                    <Menu className='bg-purple-200 text-pink-700'>
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
    )
  } else {
    return (
      <div className='h-full' id='container'>
        <div className='bg-purple-200 h-full text-pink-700 border-r-2 border-pink-50'>
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
