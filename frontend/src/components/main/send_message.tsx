import { Channel } from "@src/fetchAPI/channel";
import { Message, sendMessage } from "@src/fetchAPI/message";
import React, { useState } from "react";

type Props = {
	channelID: number;
};

export const SendMessageComponent: React.FC<Props> = (props: Props) => {
	const [text, setText] = useState("");

	const changeText = (e: React.ChangeEvent<HTMLTextAreaElement>): void => {
		setText(e.target.value);
	};

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		const res: Message | null = await sendMessage(text, props.channelID, []);
		setText("");
	};

	return (
		<form className="my-12 mx-auto" onSubmit={handleSubmit}>
			<textarea
				className="box-border w-1/3 h-2/3 border-2 border-black"
				name="text"
				value={text}
				onChange={(e: React.ChangeEvent<HTMLTextAreaElement>): void => changeText(e)}
				required
			></textarea>
			<button className="border-2 border-black" type="submit">
				Send
			</button>
		</form>
	);
};
