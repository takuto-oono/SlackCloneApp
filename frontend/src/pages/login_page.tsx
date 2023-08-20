import React, { useEffect } from 'react'
import { useCookies } from 'react-cookie'
import { LoginForm } from '@src/components/user/user'
import { Workspace, getWorkspaces } from '@src/fetchAPI/workspace'
import { useSetRecoilState } from 'recoil'
import { workspacesState } from '@src/utils/atom'
import { connectSocket, sendMessage } from '@src/utils/websocket'

function LoginPage() {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id'])
  const setWorkspaces = useSetRecoilState(workspacesState)
  connectSocket()
  sendMessage('testmess')

  // useEffect削除（未）
  useEffect(() => {
    if (cookies.token) {
      getWorkspaces().then((workspaces: Workspace[]) => {
        setWorkspaces(workspaces)
      })
    }
  }, [cookies.token])

  return (
    <div>
      <LoginForm />
    </div>
  )
}

export default LoginPage
