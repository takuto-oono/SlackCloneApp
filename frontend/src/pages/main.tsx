import { LoginForm } from "@src/components/main/user";
import { ChannelComponent } from "./main/channel";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";

const Main: React.FC = () => {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  const { channelID } = useParams<{ channelID: string }>();
  if (!cookies.token) {
    return (
			<>
        <LoginForm />
			</>
		);
  } else {
    if (channelID) {
		return (
			<>
        <ChannelComponent channelID={Number(channelID)} />
			</>
		);
    } else {
      return <></>;
    }
  }
};

export default Main;
