import { ShowMessagesComponent } from "@src/components/main/show_messages";
import React from "react";
import { SendMessageComponent } from "@src/components/main/send_message";
import { useRecoilValue } from "recoil";
import { workspaceChannelsState } from "@src/utils/atom";

type Props = {
	channelID: number;
};

export const ChannelComponent: React.FC<Props> = (props: Props) => {
  const workspaceChannels = useRecoilValue(workspaceChannelsState);
  const currentChannel = workspaceChannels.find(channel => channel.id === props.channelID)
  const channelName = currentChannel?.name

	return (
    <div className="h-full border border-black w-full">
      <p className="text-xl p-2">{channelName}</p>
			<div className="h-5/6 border border-black w-full overflow-y-scroll">
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
