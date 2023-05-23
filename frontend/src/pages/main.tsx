import { ChannelComponent } from "./main/channel";
import { Channel } from "@src/fetchAPI/channel";
import { useParams } from "react-router-dom";

const Main: React.FC = () => {
	const { channelID } = useParams<{ channelID: string }>();

	if (channelID) {
		return (
			<div className="w-full">
				<ChannelComponent channelID={Number(channelID)} />
			</div>
		);
	}

	return <></>;
};

export default Main;
