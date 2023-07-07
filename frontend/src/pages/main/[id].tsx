import { LoginForm } from "@src/components/main/user";
import { ChannelComponent } from "../../components/main/channel";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";
import { useRouter } from "next/router";

const Main: React.FC = () => {
  const [cookies, setCookie, removeCookie] = useCookies(['token', 'user_id']);
  const { channelID } = useParams<{ channelID: string }>();
  const router = useRouter()
  const { id, lang } = router.query

  if (id) {
    return (
      <>
        <ChannelComponent channelID={Number(id)} />
      </>
    );
  } else {
    return <>ttt</>;
  }
};

export default Main;
