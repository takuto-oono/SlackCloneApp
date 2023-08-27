import useWebSocket, { ReadyState } from 'react-use-websocket'
import { getToken } from '../utils/cookie'

export const connectSocket = () => {
  const { sendMessage } = useWebSocket(getSocketUrl(), {
    onOpen: () => {
      console.log('websocket open')
    },
    share: true,
    filter: () => false,
  })
}

export const sendMessage = (message: any) => {
  const { sendJsonMessage, readyState } = useWebSocket(getSocketUrl(), {
    share: true,
  })

  if (readyState === ReadyState.OPEN) {
    sendJsonMessage({
      message,
      type: 'userevent',
    })
    console.log(message)
  }
}

const getSocketUrl = (): string => {
  return 'ws://localhost:8000/websocket/?token=' + getToken()
}
