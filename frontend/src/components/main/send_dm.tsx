import { workspaceIdState } from "@src/utils/atom";
import { getUserId } from "@src/utils/cookie";
import React, { useState } from "react";
import { useRecoilValue } from "recoil";
import { sendDM, SendDMForm } from "src/fetchAPI/message";

const SendDM = () => {
  const [text, setText] = useState("");
  const workspaceID = useRecoilValue(workspaceIdState);

  const doChangeText = (e: React.ChangeEvent<HTMLInputElement>) => {
    setText(e.target.value);
  };

  const doSubmit = () => {
    let form: SendDMForm = {
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
