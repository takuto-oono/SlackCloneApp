import useWebSocket, { ReadyState } from 'react-use-websocket'
import { getToken } from '../utils/cookie'

export const connectSocket = () => {
  useWebSocket(getSocketUrl(), {
    onOpen: () => {
      console.log('websocket open')
    },
    share: true,
    filter: () => false,
  })
}

const getSocketUrl = (): string => {
  return 'ws://' + getToken() + 'localhost:8000/websocket/'
}
