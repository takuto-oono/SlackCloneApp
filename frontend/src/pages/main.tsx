import ShowChannelsInWorkspace from '@src/components/channel/show_channels'
import { ChannelComponent } from '../components/channel/channel'
import { useRouter } from 'next/router'

const Main: React.FC = () => {
  const router = useRouter()
  const { channelId } = router.query
  const { contents } = router.query

  if (channelId) {
    return <ChannelComponent channelID={Number(channelId)} />
  } else if (contents == 'show_channels') {
    return <ShowChannelsInWorkspace />
  } else {
    return <></>
  }
}

export default Main
