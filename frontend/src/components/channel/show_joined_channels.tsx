import React, { useRef, useState } from 'react'
import { Menu, MenuItem, SubMenu } from 'react-pro-sidebar'
import Button from '@mui/material/Button'
import Popover from '@mui/material/Popover'
import CreateChannelForm from '@src/components/channel/create_channel'
import { useRecoilValue, useSetRecoilState } from 'recoil'
import { joinedChannelsState, usersInCState } from '@src/utils/atom'
import { useRouter } from 'next/router'
import { getUsersInChannel } from '@src/fetchAPI/channel'
import { UserInWorkspace } from '@src/fetchAPI/workspace'

function ShowJoinedChannels() {
  const [open, setOpen] = useState(false)
  const divRef = useRef(null)
  const joinedChannels = useRecoilValue(joinedChannelsState)
  const router = useRouter()
  const setUsersInC = useSetRecoilState(usersInCState)

  const handleClickOpen = () => {
    setOpen(true)
  }
  const handleClose = () => {
    setOpen(false)
  }

  const getChannelInfo = (channelId: number) => {
    getUsersInChannel(channelId).then((usersInC: UserInWorkspace[]) => {
      setUsersInC(usersInC)
    })
    router.push({
      pathname: `/main`,
      query: { channelId: channelId },
    })
  }

  const list = joinedChannels.map((channel, index) => (
    <div key={index}>
      <MenuItem className='bg-purple-200'>
        <button
          type='button'
          onClick={() => getChannelInfo(channel.id)}
          className='inline-block align-baseline text-sm'
        >
          <div className='truncate'>#{channel.name}</div>
        </button>
      </MenuItem>
    </div>
  ))

  return (
    <Menu className='pd-5 bg-purple-200 text-pink-700'>
      <SubMenu label='Channels' className='truncate w-36'>
        {list}
        <div ref={divRef} className='bg-purple-200'>
          <Button onClick={handleClickOpen}>
            <p className='bg-purple-200 text-pink-700'>+ チャンネルを追加</p>
          </Button>
          <Popover
            open={open}
            anchorEl={divRef.current}
            onClose={handleClose}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'left',
            }}
          >
            <CreateChannelForm />
            <button
              className='px-4 py-1 border w-full border-slate-400 hover:bg-slate-50'
              onClick={() =>
                router.push({
                  pathname: `/main`,
                  query: { contents: 'show_channels' },
                })
              }
            >
              <p className='text-black'>チャンネル一覧</p>
            </button>
          </Popover>
        </div>
      </SubMenu>
    </Menu>
  )
}

export default ShowJoinedChannels
