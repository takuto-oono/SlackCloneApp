import React from 'react'
import { UserInWorkspace } from '@fetchAPI/workspace'
import { useRecoilValue, useSetRecoilState } from 'recoil'
import { usersInCState, channelsState } from '@src/utils/atom'
import { getUsersInChannel } from '@src/fetchAPI/channel'
import { useRouter } from 'next/router'
import CreateChannelForm from './create_channel'

function ShowChannelsInWorkspace() {
  const router = useRouter()
  const setUsersInC = useSetRecoilState(usersInCState)
  const channels = useRecoilValue(channelsState)

  const getChannelInfo = (channelId: number) => {
    getUsersInChannel(channelId).then((usersInC: UserInWorkspace[]) => {
      setUsersInC(usersInC)
    })
    router.push({
      pathname: `/main`,
      query: { channelId: channelId },
    })
  }

  const list = channels.map((channel, index) => (
    <div key={index} className='border-t border-slate-400'>
      <button
        type='button'
        onClick={() => getChannelInfo(channel.id)}
        className='text-left my-2.5 text-pink-700'
      >
        <div>
          <p className='truncate text-lg'>#{channel.name}</p>
        </div>
      </button>
    </div>
  ))

  return (
    <div className='m-2.5'>
      <div className='my-2.5 flex'>
        <div className='text-lg'>すべてのチャンネル</div>
        <div className='pl-8 items-end'>
          <CreateChannelForm />
        </div>
      </div>

      {list}
    </div>
  )
}

export default ShowChannelsInWorkspace
