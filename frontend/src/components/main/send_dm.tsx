import { getUserId } from "@src/fetchAPI/cookie";
import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { sendDM, SendDMForm } from "src/fetchAPI/message";

const SendDM = () => {
  const [text, setText] = useState("");
  const workspaceID: number = Number(useParams<{ workspaceId: string }>().workspaceId)
  
  const doChangeText = (e: React.ChangeEvent<HTMLInputElement>) => {
    setText(e.target.value);
  };

  const doSubmit = () => {
    console.log("workspace id", workspaceID);
    let form: SendDMForm;
    form = {
      receivedUserID: Number(getUserId()),
      workspaceID: workspaceID,
      text: text,
    }
    const res = sendDM(form);
		console.log(res);
  };

  return (
    <div>
      <h2>DM</h2>
			<label htmlFor="text">
				send message
				<input
					type="text"
					name="text"
					onChange={doChangeText}
				/>
			</label>
			<br />

      <button type="submit" onClick={doSubmit}>
        送る
      </button>
    </div>
  );
};

export default SendDM;
