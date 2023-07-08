import { ChannelComponent } from "../components/main/channel";
import { useRouter } from "next/router";

const Main: React.FC = () => {
  const router = useRouter()
  const { channelId } = router.query

  if (channelId) {
    return (
      <>
        <ChannelComponent channelID={Number(channelId)} />
      </>
    );
  } else {
    return <></>;
  }
};

export default Main;
