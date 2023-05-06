import React, { useState } from "react";
import { sendDM, SendDM } from "src/fetchAPI/send_dm_api";

const SendDMForm = () => {
  const [received_user_id, setReceived_user_id] = useState(0);
  const [workspace_id, setWorkspace_id] = useState(0);
  const [text, setText] = useState("");
  const [form, setForm] = useState({
    received_user_id: 0,
    workspace_id: 0,
    text: "",
  });

  const doChangeReceiveUserId = (e: any) => {
    setForm((preSetting) => ({ ...preSetting, received_user_id: Number(e.target.value) }));
  };

  const doChangeWorkspaceId = (e: any) => {
    setForm((prevState) => ({ ...prevState, workspace_id: Number(e.target.value) }));
  };

  const doChangeText = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prevState) => ({ ...prevState, text: e.target.value }));
  };

  const doSubmit = () => {
    sendDM(form);
		console.log(form);
  };

  return (
    <div>
      <h2>DM</h2>
			<label htmlFor="text">
				received user_id
				<input
					type="number"
					name="receive_id"
					onChange={doChangeReceiveUserId}
				/>
			</label>
			<br />

			<label htmlFor="text">
				workspace ID
				<input
					type="number"
					name="workspace_id"
					onChange={doChangeWorkspaceId}
				/>
			</label>
			<br />

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

export default SendDMForm;
