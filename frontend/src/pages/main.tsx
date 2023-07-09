import ShowWorkspaceChannels from "@src/components/main/show_workspace_channels";
import { ChannelComponent } from "../components/main/channel";
import { useRouter } from "next/router";

const Main: React.FC = () => {
  const router = useRouter()
  const { channelId } = router.query
  const { contents } = router.query

  if (channelId) {
    return (
      <ChannelComponent channelID={Number(channelId)} />
    );
  } else if (contents == "show_workspace_channels") {
    return (
      <ShowWorkspaceChannels />
    );
  } else {
    return <></>;
  }
};

export default Main;
