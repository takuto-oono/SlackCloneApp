import { getMessagesFromChannel, Message } from "@src/fetchAPI/message";
import { UserInWorkspace } from "@src/fetchAPI/workspace";
import React, { useState, useEffect } from "react";
import { useRecoilValue } from "recoil";
import { usersInWState } from "@src/components/sideNav1/show_workspaces";

type Props = {
	channelID: number;
};

const showItem: React.FC<Message> = (
	message: Message,
	users: UserInWorkspace[]
) => {
	let userName: string = "";
	for (const user of users) {
		if (user.id == message.userID) {
			userName = user.name;
			break;
		}
	}

	return (
		<>
			<p>{userName}</p>
			<p>{message.createdAt}</p>
			<p className="break-words">{message.text}</p>
		</>
	);
};

export const ShowMessagesComponent: React.FC<Props> = (props: Props) => {
	const [messages, setMessages] = useState<Message[]>([]);
	const users: UserInWorkspace[] = useRecoilValue(usersInWState);

	useEffect(() => {
		const timer = setInterval(async () => {
			const messages: Message[] | null = await getMessagesFromChannel(
				props.channelID
			);
			if (messages != null) {
				setMessages(messages.reverse());
			}
		}, 1000);
		return () => {
			clearInterval(timer);
		};
	}, [props]);

	return (
		<div className="w-full">
			<ul>
				{messages.map((message: Message) => (
					<li
						className="border-2 border-black w-full mx-auto my-1"
						key={message.id}
					>
						{showItem(message, users)}
					</li>
				))}
			</ul>
		</div>
	);
};
