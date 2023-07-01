import { ShowMessagesComponent } from "@src/components/main/show_messages";
import React from "react";
import { Channel } from "@src/fetchAPI/channel";
import { SendMessageComponent } from "@src/components/main/send_message";

type Props = {
	channelID: number;
};

export const ChannelComponent: React.FC<Props> = (props: Props) => {
	return (
		<div className="h-full border-2 border-black w-full">
			<div className="h-5/6 border-2 border-black w-full overflow-y-scroll">
				<ShowMessagesComponent
					channelID={props.channelID}
				></ShowMessagesComponent>
			</div>
			<div className="absolute bottom-0 h-1/6 border-2 border-black w-full bg-white">
				<SendMessageComponent
					channelID={props.channelID}
				></SendMessageComponent>
			</div>
		</div>
	);
};
